"use strict";
// block scoping
{
  const buttonEl = document.getElementById("notifications-push-button");
  const testButtonEl = document.getElementById("notifications-push-test");
  const statusEl = document.getElementById("notifications-push-status");

  function initPush() {
    if (!buttonEl) {
      return;
    }

    const vapidPublicKey = buttonEl.dataset.vapidPublicKey || "";

    function setStatus(message, isError) {
      if (!statusEl) {
        return;
      }

      statusEl.textContent = message || "";
      statusEl.classList.toggle("error", Boolean(isError));
    }

    function setButton(label, disabled) {
      buttonEl.textContent = label;
      buttonEl.disabled = disabled;
    }

    function setTestButton(label, disabled) {
      if (!testButtonEl) {
        return;
      }
      testButtonEl.textContent = label;
      testButtonEl.disabled = disabled;
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

    async function ensureSubscription() {
      const registration = await navigator.serviceWorker.register("/static/sw.js");
      let subscription = await registration.pushManager.getSubscription();

      if (!subscription) {
        subscription = await registration.pushManager.subscribe({
          userVisibleOnly: true,
          applicationServerKey: urlBase64ToUint8Array(vapidPublicKey),
        });
      }

      await saveSubscription(subscription);
    }

    async function updatePermissionUI() {
      if (!isSupported() || vapidPublicKey === "") {
        setButton("Push not supported", true);
        setTestButton("Send test notification", true);
        setStatus("Push notifications are unavailable.", true);
        return;
      }

      const permission = Notification.permission;
      if (permission === "denied") {
        setButton("Notifications blocked", true);
        setTestButton("Send test notification", true);
        setStatus("Enable notifications in your browser settings.", true);
        return;
      }

      if (permission === "granted") {
        setButton("Notifications enabled", true);
        setTestButton("Send test notification", false);
        setStatus("Enabled.");
        try {
          await ensureSubscription();
        } catch (err) {
          console.error(err);
          setStatus("Failed to register notifications.", true);
        }
        return;
      }

      setButton("Enable notifications", false);
      setTestButton("Send test notification", true);
      setStatus("");
    }

    buttonEl.addEventListener("click", async () => {
      if (!isSupported() || vapidPublicKey === "") {
        setStatus("Push notifications are unavailable.", true);
        return;
      }

      setButton("Enabling...", true);
      setTestButton("Send test notification", true);
      try {
        const permission = await Notification.requestPermission();
        if (permission !== "granted") {
          setButton("Notifications blocked", true);
          setTestButton("Send test notification", true);
          setStatus("Enable notifications in your browser settings.", true);
          return;
        }

        await ensureSubscription();
        setButton("Notifications enabled", true);
        setTestButton("Send test notification", false);
        setStatus("Enabled.");
      } catch (err) {
        console.error(err);
        setButton("Enable notifications", false);
        setTestButton("Send test notification", true);
        setStatus("Failed to register notifications.", true);
      }
    });

    if (testButtonEl) {
      testButtonEl.addEventListener("click", async () => {
        setTestButton("Sending...", true);
        setStatus("");
        try {
          const res = await fetch("/notifications/push-test-self", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
              title: "Test notification",
              body: "Hello from Operational Core.",
              url: "/notifications",
            }),
          });
          if (!res.ok) {
            throw new Error(`Failed to send test: ${res.status}`);
          }
          setStatus("Test notification sent.");
        } catch (err) {
          console.error(err);
          setStatus("Failed to send test notification.", true);
        } finally {
          setTestButton("Send test notification", false);
        }
      });
    }

    updatePermissionUI();
  }

  initPush();
}
