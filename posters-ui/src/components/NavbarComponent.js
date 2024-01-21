import React from 'react';
import AppBar from '@mui/material/AppBar';
import Toolbar from '@mui/material/Toolbar';
import Typography from '@mui/material/Typography';

function Navbar() {
    return (
        <AppBar position="static">
            <Toolbar>
                <Typography variant="h6">
                    MindBrush
                </Typography>
            </Toolbar>
        </AppBar>
    );
}

export default Navbar;
