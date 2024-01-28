// handleDownload.js 

const handleDownload = (image) => {
    downloadImage(image);
};

const downloadImage = (image) => {
    const link = document.createElement('a');
    link.href = image;
    link.download = `downloaded_image.png`; // Naming the file
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
};

export default handleDownload;

