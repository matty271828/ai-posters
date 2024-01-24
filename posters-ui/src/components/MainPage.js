import React, { useState, useEffect } from 'react';
import TextField from '@mui/material/TextField';
import Button from '@mui/material/Button';
import CircularProgress from '@mui/material/CircularProgress';
import Box from '@mui/material/Box';
import Grid from '@mui/material/Grid';
import Container from '@mui/material/Container';
import generateImage from '../services/imageService';
import handleDownload from '../services/handleDownload';

function MainPage() {
    const [prompt, setPrompt] = useState('');
    const [loading, setLoading] = useState(false);
    const [images, setImages] = useState([]);

    const handleGenerate = async () => {
        setLoading(true);
        await generateImage(prompt, setImages);
        setLoading(false);
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
            {images.length > 0 && (
                <Button 
                    variant="contained" 
                    color="secondary" 
                    onClick={() => handleDownload(images)}>
                    Download
                </Button>
            )}
        </Container>
    );
}

export default MainPage;

