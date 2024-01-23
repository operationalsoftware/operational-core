(function () {
  const pagination = me();
  const paginationList = me("ul", pagination);
  const prevPageBtn = me("li.pagination__btn--left", paginationList);
  const nextPageBtn = me("li.pagination__btn--right", paginationList);

  let currentPageNumber = parseInt(pagination.dataset.currentPage);
  const totalPages = parseInt(pagination.dataset.totalPages);

  function updatePaginationState() {
    if (currentPageNumber === 1) {
      prevPageBtn.classList.add("pagination__btn--disabled");
    } else {
      prevPageBtn.classList.remove("pagination__btn--disabled");
    }

    if (currentPageNumber === totalPages) {
      nextPageBtn.classList.add("pagination__btn--disabled");
    } else {
      nextPageBtn.classList.remove("pagination__btn--disabled");
    }
  }

  function initPagination() {
    updatePaginationState();
    addUrlParams([
      { name: "page", value: currentPageNumber },
      { name: "pageSize", value: 100 },
    ]);
  }

  function goToNextPage() {
    const nextPageNumber = currentPageNumber + 1;
    const pages = paginationList.children;

    if (nextPageNumber === totalPages) {
      nextPageBtn.classList.add("pagination__btn--disabled");
    }

    if (nextPageNumber <= totalPages) {
      pages[currentPageNumber].classList.remove("pagination__item--current");
      pages[nextPageNumber].classList.add("pagination__item--current");
      currentPageNumber++;
      updatePaginationState();
      addUrlParams([
        { name: "page", value: currentPageNumber },
        { name: "pageSize", value: 100 },
      ]);
    }
  }

  function goToPrevPage() {
    const prevPageNumber = currentPageNumber - 1;
    const pages = paginationList.children;

    if (prevPageNumber === 1) {
      prevPageBtn.classList.add("pagination__btn--disabled");
    }

    if (prevPageNumber > 0) {
      pages[currentPageNumber].classList.remove("pagination__item--current");
      pages[prevPageNumber].classList.add("pagination__item--current");
      currentPageNumber--;
      updatePaginationState();
      addUrlParams([
        { name: "page", value: currentPageNumber },
        { name: "pageSize", value: 100 },
      ]);
    }
  }

  // Event Listeners
  prevPageBtn.on("click", (e) => {
    e.preventDefault();
    goToPrevPage();
  });

  nextPageBtn.on("click", (e) => {
    e.preventDefault();
    goToNextPage();
  });

  paginationList.on("click", (e) => {
    e.preventDefault();

    if (e.target.tagName === "A") {
      const pages = paginationList.children;
      const pageNumber = parseInt(e.target.textContent);

      pages[currentPageNumber].classList.remove("pagination__item--current");
      pages[pageNumber].classList.add("pagination__item--current");
      currentPageNumber = pageNumber;
      updatePaginationState();
      addUrlParams([
        { name: "page", value: currentPageNumber },
        { name: "pageSize", value: 100 },
      ]);
    }
  });

  initPagination();
})();
