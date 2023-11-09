(function () {
  const rangeInput = me("input[type='range']");
  const showValue = me(".show-range-value");

  showValue.textContent = rangeInput.value;

  rangeInput.on("input", (e) => {
    showValue.textContent = e.target.value;
  });

  // on(
  //   any("input[type='range']"),
  //   "input",
  //   (e) => (me(".show-range-value").textContent = e.target.value)
  // );
})();
