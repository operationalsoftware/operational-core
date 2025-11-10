function deleteMetric(e) {
  e.preventDefault();
  const targetBtn = e.currentTarget;
  const metricId = targetBtn.dataset.id;

  confirmUpdate = confirm(`Are you sure you want to delete this Metric?`);

  if (confirmUpdate) {
    fetch(`/services/metrics/${metricId}`, {
      method: "DELETE",
    }).then((res) => {
      if (res.ok) {
        window.location.href = "/services/metrics";
      } else {
        alert("Failed to delete service metric");
      }
    });
  }
}
