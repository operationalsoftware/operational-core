(function () {
  const gallery = document.currentScript.closest(".gallery");
  if (!gallery) return;

  const lightbox = document.createElement("div");
  lightbox.className = "lightbox";
  lightbox.style.display = "none";
  document.body.appendChild(lightbox);

  const img = document.createElement("img");
  lightbox.appendChild(img);

  const prevBtn = document.createElement("button");
  prevBtn.className = "nav-btn prev";
  prevBtn.textContent = "‹";
  lightbox.appendChild(prevBtn);

  const nextBtn = document.createElement("button");
  nextBtn.className = "nav-btn next";
  nextBtn.textContent = "›";
  lightbox.appendChild(nextBtn);

  const closeBtn = document.createElement("button");
  closeBtn.className = "nav-btn close";
  closeBtn.textContent = "✕";
  lightbox.appendChild(closeBtn);

  let currentImages = [];
  let currentIndex = 0;

  const images = gallery.querySelectorAll(".gallery-item img");

  function openLightbox(index) {
    currentImages = Array.from(images);
    currentIndex = index;
    img.src = currentImages[currentIndex].src;
    lightbox.style.display = "flex";
    document.addEventListener("keydown", handleKeydown);
  }

  function closeLightbox() {
    lightbox.style.display = "none";
    document.removeEventListener("keydown", handleKeydown);
  }

  function handleKeydown(e) {
    if (e.key === "Escape") {
      closeLightbox();
    }
  }

  images.forEach((image, index) => {
    image.addEventListener("click", () => openLightbox(index));
  });

  prevBtn.addEventListener("click", (e) => {
    e.stopPropagation();
    if (currentImages.length === 0) return;
    currentIndex =
      (currentIndex - 1 + currentImages.length) % currentImages.length;
    img.src = currentImages[currentIndex].src;
  });

  nextBtn.addEventListener("click", (e) => {
    e.stopPropagation();
    if (currentImages.length === 0) return;
    currentIndex = (currentIndex + 1) % currentImages.length;
    img.src = currentImages[currentIndex].src;
  });

  closeBtn.addEventListener("click", (e) => {
    e.stopPropagation();
    closeLightbox();
  });

  lightbox.addEventListener("click", (e) => {
    if (e.target === lightbox) {
      closeLightbox();
    }
  });
})();

