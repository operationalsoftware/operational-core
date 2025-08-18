(function() {
  const toast = document.querySelector(".toast");
  if (toast) {
    toast.style.opacity = "1";
    toast.style.pointerEvents = "auto";
    setTimeout(() => {
      toast.style.opacity = "0";
    }, 4000);
  }
})();

