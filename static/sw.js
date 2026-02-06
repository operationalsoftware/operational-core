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

function closeTaggedNotifications(tag) {
  if (!tag) {
    return Promise.resolve(0);
  }
  if (!self.registration.getNotifications) {
    return Promise.resolve(0);
  }
  return self.registration.getNotifications({ tag }).then((list) => {
    list.forEach((notification) => notification.close());
    return list.length;
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
    const tag = payload.notificationId
      ? `notification:${payload.notificationId}`
      : "";
    event.waitUntil(
      closeTaggedNotifications(tag)
        .then((count) => {
          if (count === 0) {
            // Fallback: clear all app notifications when tag-based lookup isn't supported.
            return closeAllNotifications();
          }
          return null;
        })
        .then(() => broadcastTrayRefresh())
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

  if (payload.notificationId) {
    options.tag = `notification:${payload.notificationId}`;
  }

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
