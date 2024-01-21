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
            const endpoint = "http://localhost:8080/api/generate-image"; // Your API endpoint
    
            // POST request with the prompt as payload
            const response = await fetch(endpoint, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ prompt: prompt }), // 'prompt' is your state variable
            });
    
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
    
            // Assuming the response is a binary image
            const imageBlob = await response.blob();
            const imageUrl = URL.createObjectURL(imageBlob);
    
            // Update the images state with the new image URL
            setImages([imageUrl]);
        } catch (error) {
            console.error('Error fetching data: ', error);
            // Handle errors as needed (e.g., show an error message to the user)
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

