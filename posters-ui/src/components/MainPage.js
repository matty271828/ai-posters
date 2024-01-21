import React, { useState, useEffect } from 'react';
import TextField from '@mui/material/TextField';
import Button from '@mui/material/Button';
import CircularProgress from '@mui/material/CircularProgress';
import Box from '@mui/material/Box';
import Grid from '@mui/material/Grid';
import Container from '@mui/material/Container';

function MainPage() {
    const [prompt, setPrompt] = useState('');
    const [loading, setLoading] = useState(false);
    const [images, setImages] = useState([]);

    const handleGenerate = async () => {
        setLoading(true);
        try {
            // First Endpoint: Generate Image
            const generateEndpoint = "http://localhost:8080/api/generate-image";
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
            const frameEndpoint = "http://localhost:8080/api/frame-image?frameSize=small";
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
        } finally {
            setLoading(false);
        }
    };
       

    useEffect(() => {
        // Add resize listener
        window.addEventListener('resize', handleResize);
        return () => {
            window.removeEventListener('resize', handleResize);
        };
    }, []);

    const handleResize = () => {
        // Adjust layout based on screen width
    };

    return (
        <Container>
            <Box my={4}>
                <TextField 
                    label="Enter a prompt" 
                    variant="outlined"
                    value={prompt}
                    onChange={(e) => setPrompt(e.target.value)}
                    fullWidth
                    margin="normal"
                />
                <Button 
                    variant="contained" 
                    color="primary" 
                    onClick={handleGenerate}
                    disabled={loading}
                >
                    Generate
                </Button>
                {loading && <CircularProgress />}
            </Box>
            <Grid container spacing={2}>
                {images.map((image, index) => (
                    <Grid item xs={12} md={6} key={index}>
                        <img src={image} alt="Generated" style={{ width: '100%' }} />
                    </Grid>
                ))}
            </Grid>
        </Container>
    );
}

export default MainPage;

