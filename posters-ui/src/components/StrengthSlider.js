import React from 'react';
import Box from '@mui/material/Box';
import Slider from '@mui/material/Slider';
import Typography from '@mui/material/Typography';

// Props are destructured for clarity and convenience
const StrengthSlider = ({ strength, setStrength }) => {
    return (
        <Box my={2} style={{ paddingTop: '20px' }}>
            <Typography id="strength-slider" gutterBottom>
                Slide to vary influence of uploaded image.
            </Typography>
            <Slider
                value={strength}
                onChange={(event, newValue) => setStrength(newValue)}
                aria-labelledby="strength-slider"
                valueLabelDisplay="auto"
                step={0.01}
                min={0}
                max={1}
            />
        </Box>
    );
};

export default StrengthSlider;
