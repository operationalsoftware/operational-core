package layout

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _Assetsf2962b763da34337e494322f4ff6ea85a3cb804b = "me {\n  height: 30px;\n  display: flex;\n  align-items: center;\n  justify-content: center;\n  color: var(--text-color-light);\n  font-size: var(--font-size-sm);\n}\n\nme sup {\n  font-size: var(--font-size-xs);\n}\n"
var _Assetsac50f5ae9fbaaf0afb9e8e49488ce346bbb994ed = "function generateMe(startEl) {\n  return (selector) => {\n    return me(selector, startEl);\n  };\n}\n\nfunction setAriaAttribute(el) {\n  el.setAttribute(\n    \"aria-expanded\",\n    el.getAttribute(\"aria-expanded\") === \"true\" ? \"false\" : \"true\"\n  );\n}\n\nme().on(\"keydown\", (e) => {\n  if (e.key === \"Escape\") {\n    const openModal = me(\"dialog.modal.open\");\n    if (openModal) {\n      halt(e);\n      closeModal(openModal);\n    }\n  }\n});\n\nfunction createIcon(iconString) {\n  const domParser = new DOMParser();\n  const iconDoc = domParser.parseFromString(iconString, \"image/svg+xml\");\n  return iconDoc.documentElement;\n}\n\nfunction getUrlParams(specificParams) {\n  const url = window.location.href;\n  const params = {};\n  const queryString = url ? url.split(\"?\")[1] : window.location.search.slice(1);\n\n  if (queryString) {\n    const keyValuePairs = queryString.split(\"&\");\n\n    keyValuePairs.forEach((keyValue) => {\n      const [key, value] = keyValue.split(\"=\");\n      params[key] = decodeURIComponent(value.replace(/\\+/g, \" \"));\n    });\n\n    // Filter specific parameters if specified\n    if (specificParams && specificParams.length > 0) {\n      const filteredParams = {};\n      specificParams.forEach((param) => {\n        if (params.hasOwnProperty(param)) {\n          filteredParams[param] = params[param];\n        }\n      });\n      return filteredParams;\n    }\n  }\n\n  return params;\n}\n\nfunction addUrlParams(params) {\n  let url = new URL(window.location.href);\n  let queryParams = new URLSearchParams(window.location.search);\n\n  params.forEach((param) => {\n    queryParams.set(param.name, param.value);\n  });\n  url.search = queryParams.toString();\n\n  // Use pushState to update the browser URL without reloading the page\n  window.history.pushState({ path: url.href }, \"\", url.href);\n}\n\nfunction removeUrlParams(params) {\n  let url = new URL(window.location.href);\n  let queryParams = new URLSearchParams(window.location.search);\n\n  params.forEach((param) => {\n    queryParams.delete(param);\n  });\n  url.search = queryParams.toString();\n\n  // Use pushState to update the browser URL without reloading the page\n  window.history.pushState({ path: url.href }, \"\", url.href);\n}\n\nfunction openModal(el) {\n  me(\"body\").styles({ overflow: \"hidden\" });\n  el.classRemove(\"hidden\");\n  el.classAdd(\"open\");\n  el.showModal();\n}\n\nfunction closeModal(el) {\n  el.classRemove(\"open\");\n  el.classAdd(\"hidden\");\n  setTimeout(() => {\n    me(\"body\").styles({ overflow: \"auto\" });\n    el.close();\n  }, 250);\n}\n\n(function () {\n  // find if theme cookie exists\n  const cookies = document.cookie;\n  const themeCookie = cookies.split(\";\").find((cookie) => {\n    return cookie.trim().startsWith(\"theme=\");\n  });\n\n  if (themeCookie) {\n    const theme = themeCookie.split(\"=\")[1];\n    document.documentElement.setAttribute(\"data-theme\", theme);\n  } else {\n    // if theme cookie doesn't exist, set it to dark\n    let theme = window.matchMedia(\"(prefers-color-scheme: dark)\").matches\n      ? \"dark\"\n      : \"light\";\n    document.documentElement.setAttribute(\"data-theme\", theme);\n    document.cookie = `theme=${theme};path=/;max-age=31536000`;\n  }\n})();\n"
var _Assetsff9ed4d1c0d5f6a8ed200a512336471201c6127f = "/* body */\nme {\n  display: flex;\n  flex-direction: column;\n  min-height: 100vh;\n  background-color: var(--background-color);\n}\n\nme main {\n  flex-grow: 1;\n}\n\nme main.main-padding {\n  padding: var(--spacing-lg);\n}\n"
var _Assets8466ed20315f3bc808e945b0b594f15aad744e3c = "me {\n  background: var(--background-color-grey);\n  padding: var(--spacing-sm);\n  font-size: var(--font-size-sm);\n}\n\nme ol {\n  margin: 0;\n  padding: 0;\n  list-style: none;\n  display: flex;\n  justify-content: flex-start;\n  align-items: center;\n}\n\nme ol li {\n  margin-right: var(--spacing-xs);\n  display: flex;\n  align-items: center;\n}\n\nme ol li .icon {\n  width: 20px;\n  height: 20px;\n  fill: var(--text-color-light);\n  margin-right: var(--spacing-xs);\n}\n\nme ol li .divider {\n  padding: var(--spacing-sm);\n  font-size: var(--font-size-lg);\n  color: var(--text-color-light);\n}\n\nme ol li div,\nme ol li a {\n  display: flex;\n  align-items: center;\n  color: var(--text-color-light);\n  padding: var(--spacing-xs) var(--spacing-sm);\n}\n\nme ol li a {\n  text-decoration: none;\n  border-radius: var(--border-radius-sm);\n  transition: background-color 0.2s ease-in-out;\n}\n\nme ol li a:hover {\n  color: var(--text-color);\n  background-color: var(--background-color-content);\n}\n"
var _Assets4ced01a70d4c76f05f76b3ea328abe2356b805de = "me {\n  height: 60px;\n  background-color: var(--primary-color);\n}\n\nme .nav_links-container {\n  width: 100%;\n  height: 100%;\n  display: flex;\n  justify-content: space-between;\n  align-items: center;\n  padding: var(--spacing);\n}\n\nme .nav_links {\n  display: flex;\n  justify-content: space-between;\n  align-items: center;\n  height: 100%;\n}\n\nme .logo-container {\n  height: 100%;\n  padding: var(--spacing);\n}\n\nme .logo-container img {\n  height: 100%;\n}\n"

// Assets returns go-assets FileSystem
var Assets = assets.NewFileSystem(map[string][]string{"/": []string{"global.js", "layout.css", "breadcrumbs.css", "navbar.css", "footer.css"}}, map[string]*assets.File{
	"/navbar.css": &assets.File{
		Path:     "/navbar.css",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1706218940, 1706218940373227300),
		Data:     []byte(_Assets4ced01a70d4c76f05f76b3ea328abe2356b805de),
	}, "/footer.css": &assets.File{
		Path:     "/footer.css",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1706218940, 1706218940373227300),
		Data:     []byte(_Assetsf2962b763da34337e494322f4ff6ea85a3cb804b),
	}, "/": &assets.File{
		Path:     "/",
		FileMode: 0x800001fd,
		Mtime:    time.Unix(1707132617, 1707132617862661505),
		Data:     nil,
	}, "/global.js": &assets.File{
		Path:     "/global.js",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1706218940, 1706218940373227300),
		Data:     []byte(_Assetsac50f5ae9fbaaf0afb9e8e49488ce346bbb994ed),
	}, "/layout.css": &assets.File{
		Path:     "/layout.css",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1706218940, 1706218940373227300),
		Data:     []byte(_Assetsff9ed4d1c0d5f6a8ed200a512336471201c6127f),
	}, "/breadcrumbs.css": &assets.File{
		Path:     "/breadcrumbs.css",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1706218940, 1706218940373227300),
		Data:     []byte(_Assets8466ed20315f3bc808e945b0b594f15aad744e3c),
	}}, "")
