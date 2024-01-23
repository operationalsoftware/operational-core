(function () {
  const ARROW_UP = createIcon(
    `<svg class="icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M13,20H11V8L5.5,13.5L4.08,12.08L12,4.16L19.92,12.08L18.5,13.5L13,8V20Z" /></svg>`
  );
  const ARROW_DOWN = createIcon(
    `<svg class="icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M11,4H13V16L18.5,10.5L19.92,11.92L12,19.84L4.08,11.92L5.5,10.5L11,16V4Z" /></svg>`
  );
  const DEFAULT_ICON = createIcon(
    `<svg class="icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M8,19H11V23H13V19H16L12,15L8,19M16,5H13V1H11V5H8L12,9L16,5M4,11V13H20V11H4Z" /></svg>`
  );

  const table = me("table");
  const tableHeads = any("table thead th", table);
  const tableRows = me("table tbody tr", table);

  tableHeads.on("click", (e) => {
    if (e.target.matches(".icon")) {
      const parent = e.target.parentElement;
      const key = parent.getAttribute("data-key");
      const sortType = parent.getAttribute("data-sort");

      if (sortType === "") {
        parent.setAttribute("data-sort", "asc");
        parent.appendChild(ARROW_UP.cloneNode(true));
        addUrlParams([{ name: key, value: "asc" }]);
      } else if (sortType === "asc") {
        parent.setAttribute("data-sort", "desc");
        parent.appendChild(ARROW_DOWN.cloneNode(true));
        addUrlParams([{ name: key, value: "desc" }]);
      } else if (sortType === "desc") {
        parent.setAttribute("data-sort", "");
        removeUrlParams([key]);
        parent.appendChild(DEFAULT_ICON.cloneNode(true));
      }
      parent.removeChild(e.target);
    }
  });
})();
