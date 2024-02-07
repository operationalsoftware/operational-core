(function () {
  const table = me("table");
  const tableHeads = any("table thead th", table);
  const tableRows = me("table tbody tr", table);

  /**
   * Return a structured sort object from the url query params
   *
   * @returns {{key: string, sort: string}[]}
   */
  function getSortFromUrl() {
    let sortParamsStr = new URLSearchParams(window.location.search).get("sort");
    if (!sortParamsStr) return [];
    let sortParams = [];
    const sortParamsTokens = sortParamsStr.split("_");
    sortParamsTokens.forEach((token) => {
      const [key, sort] = token.split("-");
      if (!key || !sort) return;
      // case insensitive regex to check if sort is asc or desc
      if (!/^(asc|desc)$/i.test(sort)) return;
      sortParams.push({ key, sort: sort.toUpperCase() });
    });
    return sortParams;
  }

  /**
   * Update the url query params with the sort params
   * @param {{key: string, sort: string}[]} sortParams
   * @returns {void}
   */
  function updateSortInUrl(sortParams) {
    let url = new URL(window.location.href);
    const sortParamsStr = sortParams
      .map((param) => `${param.key}-${param.sort}`)
      .join("_");

    url.searchParams.set("sort", sortParamsStr);

    window.history.pushState({ path: url.href }, "", url.href);
  }

  document.body.addEventListener("htmx:beforeRequest", (ev) => {
    console.log(ev);
    const target = ev.detail.elt;
    const key = target.getAttribute("data-key");
    const currentSortDirection = target.getAttribute("data-sort");
    const sortParams = getSortFromUrl();
    const sortIndex = sortParams.findIndex((param) => param.key === key);
    if (sortIndex === -1) {
      const newSortParams = [...sortParams, { key, sort: "ASC" }];
      updateSortInUrl(newSortParams);
      return;
    } else {
      // Allow the user to remove the sort if it is currently descending and it is the last sort
      if (
        currentSortDirection === "DESC" &&
        sortParams.length - 1 === sortIndex
      ) {
        const newSortParams = sortParams.slice(0, -1);
        updateSortInUrl(newSortParams);
        return;
      } else if (currentSortDirection === "DESC") {
        const newSortParams = [...sortParams];
        newSortParams[sortIndex].sort = "ASC";
        updateSortInUrl(newSortParams);
        return;
      } else {
        const newSortParams = [...sortParams];
        newSortParams[sortIndex].sort = "DESC";
        updateSortInUrl(newSortParams);
        return;
      }
    }
  });
})();
