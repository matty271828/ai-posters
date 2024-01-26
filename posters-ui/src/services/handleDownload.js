// handleDownload.js 

const handleDownload = (images) => {
    if (images.length === 0) return;

    // Immediately download the first image
    downloadImage(images[0], 0);

    // For the rest of the images, introduce a delay
    const delay = 100; // milliseconds
    images.slice(1).forEach((image, index) => {
        setTimeout(() => {
            downloadImage(image, index + 1);
        }, delay * (index + 1)); // Incremental delay for each subsequent image
    });
};

const downloadImage = (image, index) => {
    const link = document.createElement('a');
    link.href = image;
    link.download = `downloaded_image_${index}.png`; // Naming the file
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
};

export default handleDownload;
