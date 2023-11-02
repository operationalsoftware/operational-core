function generateMe(startEl) {
  return (selector) => {
    return me(selector, startEl);
  };
}

me().on("keydown", (e) => {
  if (e.key === "Escape") {
    const openModal = me("dialog.modal.open");
    if (openModal) {
      halt(e);
      closeModal(openModal);
    }
  }
});

/*
// Popconfirm logic
  const confirmButton = document.getElementById('confirmButton');
  const popConfirm = document.getElementById('popconfirm');
  const confirmYes = document.getElementById('confirmYes');
  const confirmNo = document.getElementById('confirmNo');

  confirmButton.addEventListener('click', () => {
    const buttonRect = confirmButton.getBoundingClientRect();
    const bodyRect = document.body.getBoundingClientRect();

    const top = buttonRect.top - bodyRect.top - popConfirm.offsetHeight - 10;
    const left = buttonRect.left - bodyRect.left - (popConfirm.offsetWidth - buttonRect.width) / 2;

    popConfirm.style.top = `${top}px`;
    popConfirm.style.left = `${left}px`;

    popConfirm.classList.add('show');
  });

  confirmYes.addEventListener('click', () => {
    showToast('You clicked yes');
    popConfirm.classList.remove('show');
  });

  confirmNo.addEventListener('click', () => {
    showToast('You clicked no');
    popConfirm.classList.remove('show');
  });
*/

function openModal(el) {
  me("body").styles({ overflow: "hidden" });
  el.classRemove("hidden");
  el.classAdd("open");
  el.showModal();
}

function closeModal(el) {
  el.classRemove("open");
  el.classAdd("hidden");
  setTimeout(() => {
    me("body").styles({ overflow: "auto" });
    el.close();
  }, 250);
}
