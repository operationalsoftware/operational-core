(function () {
  const rangeInput = document.querySelector("input[type='range']");
  const showValue = document.querySelector(".show-range-value");

  showValue.textContent = rangeInput.value;

  rangeInput.addEventListener("input", (e) => {
    showValue.textContent = e.target.value;
  });

  // on(
  //   any("input[type='range']"),
  //   "input",
  //   (e) => (me(".show-range-value").textContent = e.target.value)
  // );
})();
