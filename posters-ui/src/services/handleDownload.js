// handleDownload.js 

const handleDownload = async (imageURL) => {
    console.log(`Initiating upscaling for image: ${imageURL} with target dimensions: 2048x2048`);
    try {
        const imageBase64 = await imageToBase64(imageURL);
        const upscaleEndpoint = `${process.env.REACT_APP_API_BASE_URL}/api/upscale-image?height=2048`;

        console.log(`Sending request to upscaling endpoint.`);
        const response = await fetch(upscaleEndpoint, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ seedBase64: imageBase64 }), // Send as base64
        });

        if (!response.ok) {
            console.error(`Upscaling failed with status: ${response.status}`);
            throw new Error(`Network response was not ok for upscale image. Status: ${response.status}`);
        }

        // Assuming the upscaled image is returned in the response
        const blob = await response.blob();
        const upscaledImageUrl = URL.createObjectURL(blob);
        console.log(`Upscaling successful. Downloading upscaled image.`);

        downloadImage(upscaledImageUrl);
    } catch (error) {
        console.error('Error during image upscaling or download:', error);
    }
};

const imageToBase64 = async (url) => {
    const response = await fetch(url);
    const blob = await response.blob();
    return new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.readAsDataURL(blob);
        reader.onloadend = () => resolve(reader.result);
        reader.onerror = reject;
    });
};

const downloadImage = (image) => {
    console.log(`Creating download link for image: ${image}`);
    const link = document.createElement('a');
    link.href = image;
    link.download = `upscaled_image.png`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    console.log(`Download initiated for: ${image}`);
};


export default handleDownload;

