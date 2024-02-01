import React, { useState, useEffect } from 'react';
import TextField from '@mui/material/TextField';
import Button from '@mui/material/Button';
import CircularProgress from '@mui/material/CircularProgress';
import Box from '@mui/material/Box';
import Grid from '@mui/material/Grid';
import Container from '@mui/material/Container';
import { generateImage2Image, generateTextToImage } from '../services/imageService';
import handleDownload from '../services/handleDownload';
import Slider from '@mui/material/Slider';
import Typography from '@mui/material/Typography';

function MainPage() {
    const [prompt, setPrompt] = useState('');
    const [loading, setLoading] = useState(false);
    const [images, setImages] = useState([]);
    const [uploadedImage, setUploadedImage] = useState(null);
    const [strength, setStrength] = useState(0.35); // Default value

    const toBase64 = (file) => new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.readAsDataURL(file);
        reader.onload = () => resolve(reader.result);
        reader.onerror = error => reject(error);
    });

    const handleImageUpload = async (event) => {
        const file = event.target.files[0];
        if (file) {
            try {
                const base64String = await toBase64(file);
                setUploadedImage(base64String);
            } catch (error) {
                console.error('Error converting file to base64:', error);
            }
        }
    };
    
    const removeUploadedImage = () => {
        setUploadedImage(null);
    };

    const handleGenerate = async () => {
        setLoading(true);
        if (uploadedImage) {
            await generateImage2Image(uploadedImage, prompt, setImages, strength);
        } else {
            await generateTextToImage(prompt, setImages);
        }
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
            <Grid container spacing={2} alignItems="center" justify="space-between">
                <Grid item xs={4} style={{ textAlign: 'left' }}>
                    <Button 
                        variant="contained" 
                        color="primary" 
                        onClick={handleGenerate}
                        disabled={loading}
                    >
                        Generate
                    </Button>
                    {loading && <CircularProgress style={{ marginLeft: '10px' }} />}
                </Grid>
                <Grid item xs={4}>
                    <Box my={2}>
                        <Typography id="strength-slider" gutterBottom>
                            Strength
                        </Typography>
                        <Slider
                            value={strength}
                            onChange={(event, newValue) => setStrength(newValue)}
                            aria-labelledby="strength-slider"
                            valueLabelDisplay="auto"
                            step={0.01}
                            marks
                            min={0}
                            max={1}
                        />
                    </Box>
                </Grid>
                <Grid item xs={4} style={{ textAlign: 'right' }}>
                    <div style={{ position: 'relative', display: 'inline-block' }}> {/* Positioned container */}
                        <input
                            accept="image/*"
                            style={{ display: 'none' }}
                            id="raised-button-file"
                            type="file"
                            onChange={handleImageUpload}
                        />
                        <label htmlFor="raised-button-file">
                            <Button variant="contained" color="secondary" component="span">
                                Upload Image
                            </Button>
                        </label>
                        {uploadedImage && (
                            <div 
                                onClick={removeUploadedImage}
                                style={{
                                    position: 'absolute',
                                    top: '0', // Adjust as necessary
                                    right: '0', // Adjust as necessary
                                    width: '20px',
                                    height: '20px',
                                    borderRadius: '50%',
                                    backgroundColor: '#f50057',
                                    color: 'white',
                                    display: 'flex',
                                    alignItems: 'center',
                                    justifyContent: 'center',
                                    cursor: 'pointer',
                                    boxShadow: '0 2px 4px rgba(0,0,0,0.2)'
                                }}
                            >
                                X
                            </div>
                        )}
                    </div>
                </Grid>
            </Grid>
            </Box>
            <Grid container spacing={2}>
                {images.map((image, index) => (
                    <Grid item xs={12} md={6} key={index}>
                        <img src={image} alt="Generated" style={{ width: '100%' }} />
                        <Button 
                            variant="contained" 
                            color="secondary" 
                            onClick={() => handleDownload(image)}
                            style={{ margin: '10px 0' }}
                        >
                            Download Image {index + 1}
                        </Button>
                    </Grid>
                ))}
            </Grid>
        </Container>
    );
    
}

export default MainPage;

