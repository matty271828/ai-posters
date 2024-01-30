// imageService.js

export const generateTextToImage = async (prompt, setImages) => {
    // Function body
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

export const generateImage2Image = async (base64Image, prompt, setImages) => {
    try {
        // Endpoint: Generate Image from Image
        const generateImage2ImageEndpoint = `${process.env.REACT_APP_API_BASE_URL}/api/generate-image2image`;
        const generateImage2ImageResponse = await fetch(generateImage2ImageEndpoint, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                prompt: prompt,
                seedBase64: base64Image, // Sending Base64-encoded image
                strength: '0.5'
            }),
        });

        if (!generateImage2ImageResponse.ok) {
            throw new Error('Network response was not ok for generate image2image');
        }

        const imageBlob = await generateImage2ImageResponse.blob();
        const generatedImage2ImageUrl = URL.createObjectURL(imageBlob);

        // Assuming you have a separate endpoint to frame the image
        // For example, you send the generated image URL or the image blob to this endpoint
        const frameEndpoint = `${process.env.REACT_APP_API_BASE_URL}/api/frame-image?frameSize=small`;
        const frameResponse = await fetch(frameEndpoint, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            // Adjust the body according to how your frame-image endpoint expects the data
            body: JSON.stringify({ imagePath: generatedImage2ImageUrl }), // or use imageBlob
        });

        if (!frameResponse.ok) {
            throw new Error('Network response was not ok for frame image');
        }

        const frameImageBlob = await frameResponse.blob();
        const framedImageUrl = URL.createObjectURL(frameImageBlob);

        // Update the images state with the new image URLs
        setImages([generatedImage2ImageUrl, framedImageUrl]);
    } catch (error) {
        console.error('Error:', error);
    }
};