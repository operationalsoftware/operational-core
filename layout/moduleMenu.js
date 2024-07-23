'use strict';
// block scoping
{

  const buttonEl = document.getElementById("navbar-module-menu-button");
  const panelEl = document.getElementById("navbar-module-menu");

  buttonEl.addEventListener('click', showPanel);

  function showPanel() {
    panelEl.classList.toggle("show");
  }

  // Add click event listener to the document to close the panel on click outside
  document.addEventListener('click', closePanel);

  function closePanel(event) {
    // Check if the click is outside both the button and the panel
    if (!buttonEl.contains(event.target) && !panelEl.contains(event.target)) {
      panelEl.classList.remove("show");
    }
  }
}

