// Auto-refresh every 30 seconds
setInterval(function () {
  fetch("/analytics/stats")
    .then((response) => response.json())
    .then((data) => {
      // Update timestamp
      document.getElementById("last-updated").textContent =
        new Date().toLocaleTimeString();
      // You could add more dynamic updates here
    })
    .catch((error) => console.error("Error fetching stats:", error));
}, 30000);
