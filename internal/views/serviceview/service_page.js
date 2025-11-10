function updateService(e) {
  e.preventDefault();
  const targetBtn = e.currentTarget;
  const resourceId = targetBtn.dataset.id;
  const serviceId = targetBtn.dataset.serviceId;
  const serviceAction = targetBtn.dataset.action;

  confirmUpdate = confirm(
    `Are you sure you want to ${serviceAction} this Service?`
  );

  if (confirmUpdate) {
    fetch(
      `/services/${serviceId}/resource/${resourceId}/${serviceAction}/update`,
      {
        method: "PUT",
      }
    ).then((res) => {
      if (res.ok) {
        window.location.href = `/services/${serviceId}`;
      } else {
        alert("Failed to update Resource Service");
      }
    });
  }
}
