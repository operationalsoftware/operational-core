# 🔐 Register Your App in Microsoft Entra ID (Azure AD)

This guide walks you through registering your web app in Microsoft Entra ID (Azure Active Directory), adding redirect URIs for both staging and production, and collecting your Client ID, Tenant ID, and Client Secret.

---

## ⚙️ 1. Open Microsoft Entra ID

1. Go to the [Azure Portal](https://portal.azure.com/).
2. From the left menu, select **Microsoft Entra ID** (formerly Azure Active Directory).
3. Under **Manage**, click **App registrations**.

---

## 🧩 2. Create a New App Registration

1. Click **+ New registration**.
2. Enter a **Name** for your app (e.g. `Company App`).
3. Under **Supported account types**, choose:  
   ➤ **Accounts in this organizational directory only (Single tenant)**
4. Under **Redirect URI**, select **Web** and add both URIs:
   ```
   https://staging.yourcompany.com/auth/microsoft/callback
   https://app.yourcompany.com/auth/microsoft/callback
   ```
5. Click **Register**.

---

## 🪪 3. Collect Your IDs

After the app is registered, you’ll be redirected to its **Overview** page.

Copy and save the following:

- **Application (client) ID** → `CLIENT_ID`
- **Directory (tenant) ID** → `TENANT_ID`

---

## 🔑 4. Create a Client Secret

1. From the left menu, under **Manage**, click **Certificates & secrets**.
2. In the **Client secrets** section, click **+ New client secret**.
3. Enter a **name** (e.g. `App Secret`) and select an **expiry period**.
4. Click **Add**.
5. **Copy the secret value immediately** — it’s displayed **only once**.

This is your `CLIENT_SECRET`.

---

## ✅ 5. Set Your Environment Variables

Save the following values in your environment variables:

```bash
MS_OAUTH_CLIENT_ID=<your client id>
MS_OAUTH_SECRET=<your client secret>
MS_OAUTH_TENANT_ID=<your tenant id>
```
