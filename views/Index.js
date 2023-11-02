(function () {
  const localMe = generateMe(me());
  const modal = localMe("dialog.modal");
  const openPopconfirm = localMe(".popconfirm button");
  const popConfirm = localMe(".popconfirm-content");

  console.log(popConfirm, openPopconfirm);

  openPopconfirm.on("click", () => {
    // const buttonRect = openPopconfirm.getBoundingClientRect();
    // const bodyRect = document.body.getBoundingClientRect();
    // const top = buttonRect.top - bodyRect.top - popConfirm.offsetHeight - 10;
    // let popConfirmLeft = Math.max(
    //   buttonRect.left -
    //     bodyRect.left -
    //     (popConfirm.offsetWidth - buttonRect.width) / 2,
    //   10 // Minimum left position to ensure it's not too close to the edge
    // );

    // // Check if the popConfirm overflows on the right side
    // if (popConfirmLeft + popConfirm.offsetWidth > bodyRect.width) {
    //   popConfirmLeft = bodyRect.width - popConfirm.offsetWidth - 10;
    // }

    // popConfirm.styles({ top: `${top}px`, left: `${popConfirmLeft}px` });
    popConfirm.classAdd("show");
  });

  // Modal Logic
  localMe("#open-modal").on("click", () => {
    openModal(modal);
  });

  localMe("#close-btn").on("click", () => {
    closeModal(modal);
  });
})();
