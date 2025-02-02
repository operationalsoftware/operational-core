(function () {
  document.addEventListener("mouseenter", (ev) => {
    const tooltip = ev.target;
    const tooltipWidth = ev.target.offsetWidth;
    const tooltipHeight = ev.target.offsetHeight;
    const rect = tooltip.getBoundingClientRect();
    // if (rect.left < tooltipWidth / 2) {
    //   tooltip.styles({ left: "0" });
    // } else if (window.innerWidth - rect.right < tooltipWidth / 2) {
    //   tooltip.styles({ left: "auto", right: "0" });
    // }
    // if (rect.top < tooltipHeight) {
    //   tooltip.style.top = "0";
    // }
  });
})();
