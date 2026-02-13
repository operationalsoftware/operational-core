"use strict";

self.addEventListener("install", (event) => {
  event.waitUntil(self.skipWaiting());
});

self.addEventListener("activate", (event) => {
  event.waitUntil(self.clients.claim());
});

function broadcastTrayRefresh() {
  return self.clients
    .matchAll({ type: "window", includeUncontrolled: true })
    .then((clientList) => {
      clientList.forEach((client) => {
        client.postMessage({ type: "notifications:refresh" });
      });
    });
}

function closeAllNotifications() {
  if (!self.registration.getNotifications) {
    return Promise.resolve();
  }
  return self.registration.getNotifications().then((list) => {
    list.forEach((notification) => notification.close());
  });
}

self.addEventListener("push", (event) => {
  let payload = {};

  if (event.data) {
    try {
      payload = event.data.json();
    } catch (err) {
      payload = { title: "Notification", body: event.data.text() };
    }
  }

  if (payload.type === "tray_refresh") {
    event.waitUntil(broadcastTrayRefresh());
    return;
  }

  if (payload.type === "notification_read") {
    event.waitUntil(
      Promise.all([closeAllNotifications(), broadcastTrayRefresh()])
    );
    return;
  }

  const title = payload.title || "Notification";
  const options = {
    body: payload.body || "",
    icon: "/static/img/logo.png",
    badge: "/static/img/logo.png",
    data: {
      url: payload.url || "/",
    },
  };

  event.waitUntil(self.registration.showNotification(title, options));
});

self.addEventListener("notificationclick", (event) => {
  event.notification.close();

  const url = (event.notification.data && event.notification.data.url) || "/";

  event.waitUntil(
    self.clients
      .matchAll({ type: "window", includeUncontrolled: true })
      .then((clientList) => {
        for (const client of clientList) {
          if (client.url === url && "focus" in client) {
            return client.focus();
          }
        }
        if (self.clients.openWindow) {
          return self.clients.openWindow(url);
        }
        return null;
      })
  );
});
