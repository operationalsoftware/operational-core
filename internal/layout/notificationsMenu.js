"use strict";
// block scoping
{
  const buttonEl = document.getElementById("navbar-notifications-menu-button");
  const panelEl = document.getElementById("navbar-notifications-menu");

  if (!buttonEl || !panelEl) {
    // Navbar notifications button is not available on this page.
  } else {
    const badgeEl = buttonEl.querySelector(".notifications-badge");
    const defaultLabel = "Notifications";
    let hasLoaded = false;
    let isLoading = false;

    function setBadge(count) {
      if (!badgeEl) {
        return;
      }

      if (count > 0) {
        const label = count > 99 ? "99+" : count.toString();
        badgeEl.textContent = label;
        badgeEl.classList.add("show");
        buttonEl.setAttribute("aria-label", `${defaultLabel} (${count} unread)`);
      } else {
        badgeEl.textContent = "";
        badgeEl.classList.remove("show");
        buttonEl.setAttribute("aria-label", defaultLabel);
      }
    }

    async function loadTray(force) {
      if (isLoading || (hasLoaded && !force)) {
        return;
      }

      isLoading = true;

      try {
        const res = await fetch("/notifications/tray", {
          headers: { "X-Requested-With": "fetch" },
        });

        if (!res.ok) {
          throw new Error(`Tray request failed: ${res.status}`);
        }

        const html = await res.text();
        panelEl.innerHTML = html;

        const trayEl = panelEl.querySelector(".notifications-tray");
        if (trayEl) {
          const count = Number.parseInt(
            trayEl.dataset.unreadCount || "0",
            10
          );
          if (!Number.isNaN(count)) {
            setBadge(count);
          }
        }

        hasLoaded = true;
      } catch (err) {
        panelEl.innerHTML =
          '<div class="notifications-tray-error">Unable to load notifications.</div>';
        setBadge(0);
        console.error(err);
      } finally {
        isLoading = false;
      }
    }

    buttonEl.addEventListener("click", () => {
      panelEl.classList.toggle("show");
      if (panelEl.classList.contains("show")) {
        loadTray(false);
      }
    });

    // Add click event listener to the document to close the panel on click outside
    document.addEventListener("click", (event) => {
      if (!buttonEl.contains(event.target) && !panelEl.contains(event.target)) {
        panelEl.classList.remove("show");
      }
    });

    if (document.readyState === "loading") {
      document.addEventListener("DOMContentLoaded", () => loadTray(false));
    } else {
      loadTray(false);
    }
  }
}
