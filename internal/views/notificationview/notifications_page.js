"use strict";
// block scoping
{
  const buttonEl = document.querySelector("[data-push-toggle]");
  const clearTestSentParam = () => {
    const url = new URL(window.location.href);
    if (url.searchParams.get("TestSent") === "1") {
      url.searchParams.delete("TestSent");
      window.history.replaceState({}, "", url.toString());
    }
  };

  clearTestSentParam();

  if (!buttonEl) {
    // No push toggle on this page.
  } else {

  const textEl = buttonEl.querySelector(".notifications-push-button-text");
  const vapidPublicKey = buttonEl.dataset.vapidPublicKey || "";
  const swPath = "/static/sw.js";

  const setState = (label, variant, disabled) => {
    if (textEl) {
      textEl.textContent = label;
    }
    buttonEl.disabled = Boolean(disabled);
    buttonEl.classList.toggle("success", variant === "success");
    buttonEl.classList.toggle("danger", variant === "danger");
  };

  const isSupported = () =>
    "serviceWorker" in navigator &&
    "PushManager" in window &&
    "Notification" in window;

  const urlBase64ToUint8Array = (base64String) => {
    const padding = "=".repeat((4 - (base64String.length % 4)) % 4);
    const base64 = (base64String + padding)
      .replace(/-/g, "+")
      .replace(/_/g, "/");
    const rawData = window.atob(base64);
    return Uint8Array.from(rawData, (char) => char.charCodeAt(0));
  };

  const getRegistration = () => navigator.serviceWorker.register(swPath);

  const getSubscription = async () => {
    const registration = await getRegistration();
    return registration.pushManager.getSubscription();
  };

  const saveSubscription = async (subscription) => {
    const res = await fetch("/notifications/subscriptions", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(subscription),
    });
    if (!res.ok) {
      throw new Error(`Failed to save subscription: ${res.status}`);
    }
  };

  const deleteSubscription = async () => {
    const res = await fetch("/notifications/subscriptions/delete", {
      method: "POST",
    });
    if (!res.ok) {
      throw new Error(`Failed to delete subscription: ${res.status}`);
    }
  };

  const ensureSubscription = async () => {
    const registration = await getRegistration();
    let subscription = await registration.pushManager.getSubscription();
    if (!subscription) {
      subscription = await registration.pushManager.subscribe({
        userVisibleOnly: true,
        applicationServerKey: urlBase64ToUint8Array(vapidPublicKey),
      });
    }
    await saveSubscription(subscription);
  };

  const removeSubscription = async () => {
    const registration = await getRegistration();
    const subscription = await registration.pushManager.getSubscription();
    if (subscription) {
      await subscription.unsubscribe();
    }
    await deleteSubscription();
  };

  const syncState = async () => {
    if (!isSupported() || vapidPublicKey === "") {
      setState("Push unsupported", "danger", true);
      return;
    }
    if (Notification.permission === "denied") {
      setState("Notifications blocked", "danger", true);
      return;
    }
    const subscription = await getSubscription().catch(() => null);
    setState(
      subscription ? "Notifications: On" : "Notifications: Off",
      subscription ? "success" : "",
      false
    );
  };

    buttonEl.addEventListener("click", async () => {
    if (!isSupported() || vapidPublicKey === "") {
      setState("Push unsupported", "danger", true);
      return;
    }
    try {
      const subscription = await getSubscription();
      if (!subscription) {
        setState("Enabling...", "", true);
        let permission = Notification.permission;
        if (permission !== "granted") {
          permission = await Notification.requestPermission();
        }
        if (permission !== "granted") {
          setState(permission === "denied" ? "Notifications blocked" : "Notifications: Off", permission === "denied" ? "danger" : "", false);
          return;
        }
        await ensureSubscription();
        setState("Notifications: On", "success", false);
      } else {
        setState("Disabling...", "", true);
        await removeSubscription();
        setState("Notifications: Off", "", false);
      }
    } catch (err) {
      console.error(err);
      setState("Update failed", "danger", false);
    }
  });

    syncState();
  }
}
