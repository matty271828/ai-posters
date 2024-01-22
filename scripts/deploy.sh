#!/bin/bash
# Check if the server IP argument is provided
if [ $# -ne 1 ]; then
    echo "Usage: $0 <server_ip>"
    exit 1
fi

#------------------------------- Configure initial set up if not done yet --------------------
# Set the remote address
REMOTE_SERVER="root@$1"
REMOTE_PATH="/root/ai-posters"
NGINX_CONFIG_LOCAL="configs/nginx/sites-available/ai-posters.conf"
NGINX_CONFIG_REMOTE="/etc/nginx/sites-available/ai-posters"

# Check if Nginx is installed, if not install it
ssh ${REMOTE_SERVER} "which nginx || sudo apt-get install nginx"

# TODO: Add SSL configuration when domain name has been purchased.

# source the local environment variables
source .env

#------------------------------- Build and transfer application files -------------------------
# Step 1: Build the Linux binary
echo "Building the Linux binary..."
GOOS=linux GOARCH=amd64 go build -o ai-posters-binary ./cmd
if [ $? -ne 0 ]; then
    echo "Error building the binary. Exiting."
    exit 1
fi

# Step 2: SSH into the remote server and stop the service
echo "Stopping the remote service..."
ssh ${REMOTE_SERVER} "sudo systemctl stop ai-poster"
if [ $? -ne 0 ]; then
    echo "Error stopping the service. Exiting."
    exit 1
fi

# Function to update and check the environment variable in the service file
update_env_var() {
    VAR_NAME="$1"
    VAR_VALUE="$2"
    
    # Check and potentially update the environment variable on the remote server
    ssh ${REMOTE_SERVER} "
      # Ensure the variable is in the [Service] section
      if grep -q 'Environment=\"${VAR_NAME}=' /etc/systemd/system/ai-poster.service; then
        # If it exists but is different, then replace
        if ! grep -q 'Environment=\"${VAR_NAME}=${VAR_VALUE}\"' /etc/systemd/system/ai-poster.service; then
          sed -i '/Environment=\"${VAR_NAME}=/c\Environment=\"${VAR_NAME}=${VAR_VALUE}\"' /etc/systemd/system/ai-poster.service
          echo "${VAR_NAME} updated"
        fi
      else
        # If it doesn't exist, insert it after [Service]
        sed -i '/\[Service\]/a Environment=\"${VAR_NAME}=${VAR_VALUE}\"' /etc/systemd/system/ai-poster.service
        echo "${VAR_NAME} added"
      fi
    "
}

# Update the environment variables
update_env_var "STABILITY_API_KEY" "$STABILITY_API_KEY"

# Step 3: Push the binary to the remote server
echo "Transferring the binary..."
scp bin/ai-poster-binary ${REMOTE_SERVER}:${REMOTE_PATH}
if [ $? -ne 0 ]; then
    echo "Error transferring the binary. Exiting."
    exit 1
fi

# Step 4: Push the UI to the remote server (assuming you have a build directory)
echo "Building the React application..."
cd posters-ui  # Navigate to your React app directory
npm run build  # Create a production build
cd ..  # Go back to the root directory

# Check if the build was successful
if [ ! -d "posters-ui/build" ]; then
    echo "React build failed. Exiting."
    exit 1
fi

echo "Transferring the UI..."
scp -r posters-ui/build ${REMOTE_SERVER}:${REMOTE_PATH}/ui
if [ $? -ne 0 ]; then
    echo "Error transferring the UI. Exiting."
    exit 1
fi

# Transfer images for the UI
echo "Transferring assets..."
scp -r assets ${REMOTE_SERVER}:${REMOTE_PATH}/assets

# Step 6: Push the Nginx config to the remote server
echo "Transferring the Nginx configuration..."
scp ${NGINX_CONFIG_LOCAL} ${REMOTE_SERVER}:${NGINX_CONFIG_REMOTE}
if [ $? -ne 0 ]; then
    echo "Error transferring the Nginx configuration. Exiting."
    exit 1
fi

#------------------------------- Reload services -------------------------
# TODO: handle lets encrypt certificate when domain name has been purchased.

# SSH into the remote server and reload Nginx
echo "Reloading Nginx on the remote server..."
ssh ${REMOTE_SERVER} "sudo nginx -t && sudo systemctl reload nginx"
if [ $? -ne 0 ]; then
    echo "Error reloading Nginx. Exiting."
    exit 1
fi

# Step 9: Ensure the Nginx configuration is linked to sites-enabled
echo "Ensuring the Nginx configuration is linked to sites-enabled..."
ssh ${REMOTE_SERVER} "sudo ln -sf ${NGINX_CONFIG_REMOTE} /etc/nginx/sites-enabled/"
if [ $? -ne 0 ]; then
    echo "Error linking the Nginx configuration. Exiting."
    exit 1
fi

# Step 10: SSH into the remote server and start the service
echo "Starting the remote service..."
ssh ${REMOTE_SERVER} "systemctl daemon-reload && sudo systemctl start ai-poster"
if [ $? -ne 0 ]; then
    echo "Error starting the service. Exiting."
    exit 1
fi

echo "Deployment complete!"