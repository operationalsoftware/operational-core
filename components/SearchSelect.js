(function () {
  const searchSelect = me(".search-select");
  const hiddenInput = me("input[type='hidden']", searchSelect);
  const searchInput = me("input[type='search']", searchSelect);
  const selectedValues = me(".selected-values", searchSelect);
  const options = me(".options", searchSelect);

  searchInput.addEventListener("keyup", () => {
    const defaultWidth = 50;
    const newWidth = defaultWidth + searchInput.value.length;

    searchInput.styles({
      width: `${newWidth + 1}px`,
    });
  });
})();
