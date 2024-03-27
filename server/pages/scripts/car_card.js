const leftButton = document.querySelector('.carousel__button--left');
const rightButton = document.querySelector('.carousel__button--right');
const carouselImages = document.querySelector('.carousel__images');
const fullscreenLeftButton = document.querySelector('.fullscreen__button--left');
const fullscreenRightButton = document.querySelector('.fullscreen__button--right');

fullscreenLeftButton.addEventListener('click', scrollLeft1);
fullscreenRightButton.addEventListener('click', scrollRight1);

let expandedImage = null;
let intervalId = null;
let fullscreenImageIndex = null;
let isFullscreen = false; 

leftButton.addEventListener('click', scrollLeft1);
rightButton.addEventListener('click', scrollRight1);

function adjustButtons() {
    if (!isFullscreen) {
        const imagesCount = carouselImages.children.length;
        if (imagesCount < 4) {
            leftButton.disabled = true;
            rightButton.disabled = true;
        }
    }
}

adjustButtons();

startAutoscroll();

carouselImages.addEventListener('mouseover', stopAutoscroll);
carouselImages.addEventListener('mouseleave', startAutoscroll);

function expandImage(image) {
    expandedImage = image.src;

    const fullscreenContainer = document.querySelector('.fullscreen-container');
    const fullscreenImage = document.querySelector('.fullscreen-image');

    fullscreenImage.src = expandedImage;
    fullscreenContainer.style.display = 'flex';

    const images = Array.from(carouselImages.children);
    fullscreenImageIndex = images.findIndex(img => img.src === expandedImage);

    stopAutoscroll();

    isFullscreen = true; 
    document.body.classList.add('no-scroll');
}


const fullscreenContainer = document.querySelector('.fullscreen-container');
const fullscreenImage = document.querySelector('.fullscreen-image');

fullscreenContainer.addEventListener('click', handleClick);

function handleClick(event) {

    const zoomButton = document.querySelector('.zoom-button');

    if (event.target !== zoomButton) {
        const screenWidth = window.innerWidth;
        const clickX = event.clientX;
        const clickY = event.clientY;
        const imageTopOffset = fullscreenImage.getBoundingClientRect().top;
        const imageBottomOffset = fullscreenImage.getBoundingClientRect().bottom;
        const imageLeftOffset = fullscreenImage.getBoundingClientRect().left;
        const imageRightOffset = fullscreenImage.getBoundingClientRect().right;
        const containerWidth = fullscreenContainer.offsetWidth;

        const isLeftArea = clickX <= 100 && clickX >= 0 && clickY >= 0 && clickY <= window.innerHeight;
        const isRightArea = clickX >= containerWidth - 100 && clickX <= containerWidth && clickY >= 0 && clickY <= window.innerHeight;
        const isImageArea = clickY >= imageTopOffset && clickY <= imageBottomOffset && clickX >= imageLeftOffset && clickX <= imageRightOffset;

        if (!isLeftArea && !isRightArea && !isImageArea) {
            closeFullscreen();
        }

        if (clickX <= 100 && event.target !== fullscreenLeftButton) {
          scrollLeft1();
        }
        if (clickX >= containerWidth - 100 && event.target !== fullscreenRightButton) {
            scrollRight1();
        }

    }
}


function closeFullscreen() {
    if (isFullscreen) {
        const fullscreenContainer = document.querySelector('.fullscreen-container');
        fullscreenContainer.style.display = 'none';

        isFullscreen = false;

        startAutoscroll();
        document.body.classList.remove('no-scroll');
    }
}

function scrollLeft1() {
  if (isFullscreen) {
    const images = Array.from(carouselImages.children);
    fullscreenImageIndex = (fullscreenImageIndex - 1 + images.length) % images.length;
    fullscreenImage.src = images[fullscreenImageIndex].src;
  } else {
    var scrollPosition = carouselImages.scrollLeft;
    var carouselWidth = carouselImages.scrollWidth - carouselImages.clientWidth;

    if (scrollPosition === 0) {
        carouselImages.scroll(carouselWidth, 0);
    } else {
        carouselImages.scrollBy(-200, 0);
    }

    clearInterval(intervalId); 
    intervalId = setInterval(scrollRight1, 5000); 
  }
}

function scrollRight1() {
  if (isFullscreen) {
    const images = Array.from(carouselImages.children);
    fullscreenImageIndex = (fullscreenImageIndex + 1) % images.length;
    fullscreenImage.src = images[fullscreenImageIndex].src;
  } else {
    var scrollPosition = carouselImages.scrollLeft;
    var carouselWidth = carouselImages.scrollWidth - carouselImages.clientWidth;

    if (scrollPosition === carouselWidth) {
        carouselImages.scroll(0, 0);
    } else {
        carouselImages.scrollBy(200, 0);
    }

    clearInterval(intervalId);
    intervalId = setInterval(scrollRight1, 5000); 
  }

}

function startAutoscroll() {
    if (!isFullscreen) { 
        intervalId = setInterval(scrollRight1, 3000);
    }
}

function stopAutoscroll() {
    if (!isFullscreen) { 
        clearInterval(intervalId);
    }
}




let zoomLevel = 1.0; 
const zoomStep = 0.2; 
const maxZoom = 2.0; 
const minZoom = 0.5; 

function handleMouseWheelFullscreen(event) {
    event.preventDefault();

    if (event.ctrlKey) {
        if (event.deltaY < 0 && zoomLevel < maxZoom) {
            zoomLevel = Math.min(zoomLevel + zoomStep, maxZoom);
            zoomIn(zoomLevel);
        } else if (event.deltaY > 0 && zoomLevel > minZoom) {
            zoomLevel = Math.max(zoomLevel - zoomStep, minZoom);
            zoomOut(zoomLevel);
        }
    } else {
        if (event.deltaY < 0) {
            scrollLeft1();
        } else {
            scrollRight1();
        }
    }
}

function zoomIn(zoomLevel) {
    fullscreenImage.style.transform = `scale(${zoomLevel})`;
}

function zoomOut(zoomLevel) {
    fullscreenImage.style.transform = `scale(${zoomLevel})`;
}

fullscreenContainer.addEventListener('wheel', handleMouseWheelFullscreen);


document.addEventListener('keydown', function(event) {
    if (isFullscreen) {
        if (event.key === 'ArrowLeft' || event.key === 'ArrowDown') {
          scrollLeft1();
        } else if (event.key === 'ArrowRight' || event.key === 'ArrowUp') {
            scrollRight1();
        } else if (event.key === 'Escape') {
            closeFullscreen();
        }
    } else {
        if (event.key === 'ArrowLeft') {
            scrollLeft1();
        } else if (event.key === 'ArrowRight') {
            scrollRight1();
        } else if (event.key === 'ArrowUp') {
            scrollRight1();
        } else if (event.key === 'ArrowDown') {
            scrollLeft1();
        }
    }
});



function zoomImage() {
    var button = document.getElementsByClassName('zoom-button')[0];
    if (button.innerHTML === '🔍 +') {
        button.innerHTML = '🔍 -';
    } else {
        button.innerHTML = '🔍 +';
    }

    if (fullscreenImage.style.transform === 'scale(1.2)') {
        fullscreenImage.style.transform = 'scale(1)';
    } else {
        fullscreenImage.style.transform = 'scale(1.2)';
    }
}

fullscreenImage.addEventListener('click', scrollRight1);



document.addEventListener("DOMContentLoaded", function() {
    const expandButton = document.getElementById("expand");
    const collapseButton = document.getElementById("collapse");
    const tables = document.querySelectorAll(".tbl");
    const headingTwo = document.querySelector(".heading.two");
    const smallHeadings = document.querySelectorAll(".smallHeading");

    expandButton.addEventListener("click", function() {
        expandButton.style.display = "none";
        collapseButton.style.display = "inline";
        tables.forEach(function(table) {
            table.style.display = "table";
        });

        smallHeadings.forEach(function(smallHeading) {
            smallHeading.style.display = "inline";
        });
        headingTwo.style.display = "inline";
    });

    collapseButton.addEventListener("click", function() {
        expandButton.style.display = "inline";
        collapseButton.style.display = "none";
        tables.forEach(function(table) {
            table.style.display = "none";
        });

        smallHeadings.forEach(function(smallHeading) {
            smallHeading.style.display = "none";
        });
  
        headingTwo.style.display = "none";
    });
});