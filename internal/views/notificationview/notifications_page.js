"use strict";
// block scoping
{
  const buttonEl = document.querySelector("[data-push-toggle]");

  function initPush() {
    if (!buttonEl) {
      return;
    }

    const textEl = buttonEl.querySelector(".notifications-push-button-text");
    const vapidPublicKey = buttonEl.dataset.vapidPublicKey || "";

    function setButtonState(label, disabled, isOn, isError) {
      if (textEl) {
        textEl.textContent = label || "";
      }
      buttonEl.disabled = Boolean(disabled);
      buttonEl.setAttribute("aria-pressed", isOn ? "true" : "false");
      buttonEl.classList.toggle("success", Boolean(isOn));
      buttonEl.classList.toggle("danger", Boolean(isError));
    }

    function setError(message) {
      setButtonState(message, true, false, true);
    }

    function setOff() {
      setButtonState("Notifications: Off", false, false, false);
    }

    function setOn() {
      setButtonState("Notifications: On", false, true, false);
    }

    function setBusy(message, isOn) {
      setButtonState(message, true, isOn, false);
    }

    function ensureText(label) {
      if (!textEl) {
        return;
      }
      textEl.textContent = label || "";
    }

    function isSupported() {
      return (
        "serviceWorker" in navigator &&
        "PushManager" in window &&
        "Notification" in window
      );
    }

    function urlBase64ToUint8Array(base64String) {
      const padding = "=".repeat((4 - (base64String.length % 4)) % 4);
      const base64 = (base64String + padding)
        .replace(/-/g, "+")
        .replace(/_/g, "/");
      const rawData = window.atob(base64);
      const outputArray = new Uint8Array(rawData.length);

      for (let i = 0; i < rawData.length; i++) {
        outputArray[i] = rawData.charCodeAt(i);
      }

      return outputArray;
    }

    async function saveSubscription(subscription) {
      const res = await fetch("/notifications/subscriptions", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(subscription),
      });

      if (!res.ok) {
        throw new Error(`Failed to save subscription: ${res.status}`);
      }
    }

    async function registerServiceWorker() {
      return navigator.serviceWorker.register("/static/sw.js");
    }

    async function getSubscription() {
      const registration = await registerServiceWorker();
      return registration.pushManager.getSubscription();
    }

    async function ensureSubscription() {
      const registration = await registerServiceWorker();
      let subscription = await registration.pushManager.getSubscription();

      if (!subscription) {
        subscription = await registration.pushManager.subscribe({
          userVisibleOnly: true,
          applicationServerKey: urlBase64ToUint8Array(vapidPublicKey),
        });
      }

      await saveSubscription(subscription);
    }

    async function removeSubscription() {
      const registration = await registerServiceWorker();
      const subscription = await registration.pushManager.getSubscription();
      if (subscription) {
        await subscription.unsubscribe();
      }
    }

    async function updatePermissionUI() {
      if (!isSupported() || vapidPublicKey === "") {
        setError("Push unsupported");
        return;
      }

      const permission = Notification.permission;
      if (permission === "denied") {
        setError("Notifications blocked");
        return;
      }

      let subscription = null;
      try {
        subscription = await getSubscription();
      } catch (err) {
        console.error(err);
      }

      if (subscription) {
        setOn();
      } else {
        setOff();
      }
    }

    buttonEl.addEventListener("click", async () => {
      if (!isSupported() || vapidPublicKey === "") {
        setError("Push unsupported");
        return;
      }

      try {
        const subscription = await getSubscription();
        const isEnabled = Boolean(subscription);

        if (!isEnabled) {
          setBusy("Enabling...", false);
          let permission = Notification.permission;
          if (permission !== "granted") {
            permission = await Notification.requestPermission();
          }
          if (permission !== "granted") {
            if (permission === "denied") {
              setError("Notifications blocked");
            } else {
              setOff();
            }
            return;
          }

          await ensureSubscription();
          setOn();
        } else {
          setBusy("Disabling...", true);
          await removeSubscription();
          setOff();
        }
      } catch (err) {
        console.error(err);
        setButtonState("Update failed", false, true, true);
      }
    });

    updatePermissionUI();
  }

  initPush();
}
