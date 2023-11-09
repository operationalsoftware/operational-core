const localMe = generateMe(me());

const openPopconfirm = localMe(".popconfirm-trigger");
const popConfirmContent = localMe(".popconfirm-content");
const popConfirmYes = localMe(".popconfirm-yes");
const popConfirmNo = localMe(".popconfirm-no");

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

  popConfirmContent.styles({
    left: `${left}px`,
    top: `${top}px`,
  });
}

openPopconfirm.on("click", () => {
  updatePopconfirmPosition();
  popConfirmContent.classRemove("hide");
  popConfirmContent.classAdd("show");
});

popConfirmYes.on("click", () => {
  popConfirmContent.classRemove("show");
  popConfirmContent.classAdd("hide");
});

popConfirmNo.on("click", () => {
  popConfirmContent.classRemove("show");
  popConfirmContent.classAdd("hide");
});

window.addEventListener("resize", () => {
  updatePopconfirmPosition();
});
