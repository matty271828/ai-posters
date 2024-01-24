// handleDownload.js 

const handleDownload = (images) => {
    images.forEach((image, index) => {
        const link = document.createElement('a');
        link.href = image;
        link.download = `downloaded_image_${index}.png`; // This names the files
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
    });
};

export default handleDownload;