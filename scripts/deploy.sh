#!/bin/bash
# Check if the server IP argument is provided
if [ $# -ne 1 ]; then
    echo "Usage: $0 <server_ip>"
    exit 1
fi

#------------------------------- Configure initial set up if not done yet --------------------
# Set the remote address
REMOTE_SERVER="root@$1"
REMOTE_PATH="/var/www/ai-posters"
NGINX_CONFIG_LOCAL="configs/nginx/sites-available/ai-posters.conf"
NGINX_CONFIG_REMOTE="/etc/nginx/sites-available/ai-posters"

# Check if Nginx is installed, if not install it
ssh ${REMOTE_SERVER} "which nginx || sudo apt-get install nginx"

# Check if SSL certificates exist
SSL_CERTIFICATE_PATH="/etc/letsencrypt/live/mindbrush.art/fullchain.pem"
SSL_SETUP=$(ssh ${REMOTE_SERVER} "[ -f ${SSL_CERTIFICATE_PATH} ] && echo '1' || echo '0'")

if [ "$SSL_SETUP" -eq "0" ]; then
    # Certificates are not set up, check if Certbot is installed
    echo "Checking if Certbot is installed..."
    CERTBOT_INSTALLED=$(ssh ${REMOTE_SERVER} "which certbot && echo '1' || echo '0'")
    
    if [ "$CERTBOT_INSTALLED" -eq "0" ]; then
        # Install Certbot and its Nginx plugin
        echo "Certbot not found. Installing Certbot..."
        ssh ${REMOTE_SERVER} "sudo apt-get update && sudo apt-get install -y certbot python3-certbot-nginx"
        echo "Certbot installation completed."
    else
        echo "Certbot is already installed."
    fi

    # Certificates are not set up, run Certbot
    read -p "Enter your email address for Let's Encrypt: " EMAIL_ADDRESS
    echo "Running Certbot..."
    ssh ${REMOTE_SERVER} "sudo certbot --nginx -d mindbrush.art -d www.mindbrush.art --email $EMAIL_ADDRESS --non-interactive --agree-tos"
    echo "Certbot setup completed."
else
    echo "SSL certificates found. Skipping Certbot setup."
fi


# source the local environment variables
source .env

#------------------------------- Build and transfer application files -------------------------
# Step 1: Build the Linux binary
echo "Building the Linux binary..."
GOOS=linux GOARCH=amd64 go build -o bin/ai-posters-binary ./cmd
if [ $? -ne 0 ]; then
    echo "Error building the binary. Exiting."
    exit 1
fi

# Step 2: SSH into the remote server and stop the service
echo "Stopping the remote service..."
ssh ${REMOTE_SERVER} "if sudo systemctl is-active --quiet ai-posters; then 
    sudo systemctl stop ai-posters
    if [ $? -ne 0 ]; then
        echo 'Error stopping the service. Exiting.'
        exit 1
    fi
else 
    echo 'Service not found or already stopped.'
fi"

# Ensure the bin directory exists on the remote server
echo "Ensuring bin directory exists on the remote server..."
ssh ${REMOTE_SERVER} "mkdir -p ${REMOTE_PATH}/bin"

# Step 3: Push the binary to the remote server
echo "Transferring the binary..."
scp bin/ai-posters-binary ${REMOTE_SERVER}:${REMOTE_PATH}/bin/
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
scp -r posters-ui/build/ ${REMOTE_SERVER}:${REMOTE_PATH}/ui
if [ $? -ne 0 ]; then
    echo "Error transferring the UI. Exiting."
    exit 1
fi

ssh ${REMOTE_SERVER} << 'EOF'
# Clear out the existing static directory
rm -rf /var/www/ai-posters/ui/static/*

# Move new files from build directory
if [ -d /var/www/ai-posters/ui/build ]; then
    mv /var/www/ai-posters/ui/build/* /var/www/ai-posters/ui/
    rm -rf /var/www/ai-posters/ui/build
fi
chown -R www-data:www-data /var/www/ai-posters/ui
chmod -R 755 /var/www/ai-posters/ui
EOF

echo "UI transferred and directory structure adjusted."

# Transfer updated assets
echo "Transferring updated assets..."
rsync -avz --progress assets/ ${REMOTE_SERVER}:${REMOTE_PATH}/assets

# Set permissions for Nginx
echo "Setting permissions for Nginx..."
ssh ${REMOTE_SERVER} "sudo chown -R www-data:www-data ${REMOTE_PATH} && sudo chmod -R 755 ${REMOTE_PATH}"
if [ $? -ne 0 ]; then
    echo "Error setting permissions for Nginx. Exiting."
    exit 1
fi

# Step 6: Transfer the systemd service file to the remote server
echo "Transferring the ai-posters service file..."
SERVICE_FILE="configs/systemd/system/ai-posters.service"
REMOTE_SERVICE_FILE="/etc/systemd/system/ai-posters.service"

# Replace placeholder in the service file and transfer
TEMP_SERVICE_FILE=$(mktemp)
sed "s/PLACEHOLDER_API_KEY/${STABILITY_API_KEY}/" ${SERVICE_FILE} > ${TEMP_SERVICE_FILE}
scp ${TEMP_SERVICE_FILE} ${REMOTE_SERVER}:${REMOTE_SERVICE_FILE}
rm ${TEMP_SERVICE_FILE}

if [ $? -ne 0 ]; then
    echo "Error transferring the ai-posters service file. Exiting."
    exit 1
fi

echo "ai-posters service file transferred successfully."

# Step 6: Push the Nginx config to the remote server
echo "Transferring the Nginx configuration..."
ssh ${REMOTE_SERVER} "sudo rm /etc/nginx/sites-enabled/default"
scp ${NGINX_CONFIG_LOCAL} ${REMOTE_SERVER}:${NGINX_CONFIG_REMOTE}
if [ $? -ne 0 ]; then
    echo "Error transferring the Nginx configuration. Exiting."
    exit 1
fi

echo "Ensuring the Nginx configuration is linked to sites-enabled..."
ssh ${REMOTE_SERVER} "sudo ln -sf ${NGINX_CONFIG_REMOTE} /etc/nginx/sites-enabled/"
if [ $? -ne 0 ]; then
    echo "Error linking the Nginx configuration. Exiting."
    exit 1
fi

#------------------------------- Reload services -------------------------
# SSH into the remote server and reload Nginx
echo "Reloading Nginx on the remote server..."
ssh ${REMOTE_SERVER} "sudo nginx -t && sudo systemctl reload nginx"
if [ $? -ne 0 ]; then
    echo "Error reloading Nginx. Exiting."
    exit 1
fi


# Step 10: SSH into the remote server and start the service
echo "Starting the remote service..."
ssh ${REMOTE_SERVER} "systemctl daemon-reload && sudo systemctl start ai-posters"
if [ $? -ne 0 ]; then
    echo "Error starting the service. Exiting."
    exit 1
fi

echo "Deployment complete!"