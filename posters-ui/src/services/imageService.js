// imageService.js

const generateImage = async (prompt, setImages) => {
    try {
        // First Endpoint: Generate Image
        const generateEndpoint = `${process.env.REACT_APP_API_BASE_URL}/api/generate-image`;
        const generateResponse = await fetch(generateEndpoint, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ prompt: prompt }),
        });

        if (!generateResponse.ok) {
            throw new Error('Network response was not ok for generate image');
        }

        const imageBlob = await generateResponse.blob();
        const generatedImageUrl = URL.createObjectURL(imageBlob);

        // Temporary Path for the Generated Image - Adjust as needed
        const imagePath = "./assets/out/v1_txt2img_0.png";

        // Second Endpoint: Frame Image
        const frameEndpoint = `${process.env.REACT_APP_API_BASE_URL}/api/frame-image?frameSize=small`;
        const frameResponse = await fetch(frameEndpoint, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ imagePath: imagePath }),
        });

        if (!frameResponse.ok) {
            throw new Error('Network response was not ok for frame image');
        }

        const frameImageBlob = await frameResponse.blob();
        const framedImageUrl = URL.createObjectURL(frameImageBlob);

        // Update the images state with the new image URLs
        setImages([generatedImageUrl, framedImageUrl]);
    } catch (error) {
        console.error('Error fetching data: ', error);
        // Handle errors as needed
    }
};

export default generateImage;
