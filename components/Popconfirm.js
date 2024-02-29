const popconfirm = document.querySelector(".popconfirm");
const openPopconfirm = popconfirm.querySelector(".popconfirm-trigger");
const popConfirmContent = popconfirm.querySelector(".popconfirm-content");
const popConfirmYes = popconfirm.querySelector(".popconfirm-yes");
const popConfirmNo = popconfirm.querySelector(".popconfirm-no");

function updatePopconfirmPosition() {
  const buttonRect = openPopconfirm.getBoundingClientRect();
  const bodyRect = document.body.getBoundingClientRect();
  const popConfirmWidth = popConfirmContent.offsetWidth;

  let left = buttonRect.left - bodyRect.left;
  const top =
    buttonRect.top -
    bodyRect.top -
    (popConfirmContent.offsetHeight * 3) / 2 -
    10;
  let arrowLeftPosition = popConfirmWidth / 2 + buttonRect.width / 2;

  // Check if popconfirm is too close to the left edge
  if (buttonRect.left < 10) {
    left = popConfirmWidth / 2;
    arrowLeftPosition = buttonRect.width / 2;
  }

  // Check if popconfirm is too close to the right edge
  const rightEdgeDistance =
    bodyRect.width - (buttonRect.left + buttonRect.width);
  if (rightEdgeDistance < popConfirmWidth / 2) {
    left = buttonRect.left - bodyRect.left - buttonRect.width / 2;
    arrowLeftPosition = popConfirmWidth - buttonRect.width / 2;
  }

  popConfirmContent.style.setProperty(
    "--arrow-left-position",
    `${arrowLeftPosition}px`
  );

  popConfirmContent.style.left = `${left}px`;
  popConfirmContent.style.top = `${top}px`;
}

if (openPopconfirm) {
  openPopconfirm.addEventListener("click", () => {
    updatePopconfirmPosition();
    popConfirmContent.classList.remove("hide");
    popConfirmContent.classList.add("show");
  });
}

popConfirmYes.addEventListener("click", () => {
  popConfirmContent.classList.remove("show");
  popConfirmContent.classList.add("hide");
});

popConfirmNo.addEventListener("click", () => {
  popConfirmContent.classList.remove("show");
  popConfirmContent.classList.add("hide");
});

window.addEventListener("resize", () => {
  updatePopconfirmPosition();
});
