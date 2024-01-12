package layout

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _Assetsff9ed4d1c0d5f6a8ed200a512336471201c6127f = "/* body */\nme {\n  display: flex;\n  flex-direction: column;\n  min-height: 100vh;\n  background-color: var(--background-color);\n}\n\nme main {\n  flex-grow: 1;\n}\n"
var _Assetsf2962b763da34337e494322f4ff6ea85a3cb804b = "me {\n  height: 30px;\n  display: flex;\n  align-items: center;\n  justify-content: center;\n  color: var(--text-color-light);\n  font-size: var(--font-size-sm);\n}\n\nme sup {\n  font-size: var(--font-size-xs);\n}\n"
var _Assetsac50f5ae9fbaaf0afb9e8e49488ce346bbb994ed = "function generateMe(startEl) {\n  return (selector) => {\n    return me(selector, startEl);\n  };\n}\n\nfunction setAriaAttribute(el) {\n  el.setAttribute(\n    \"aria-expanded\",\n    el.getAttribute(\"aria-expanded\") === \"true\" ? \"false\" : \"true\"\n  );\n}\n\nme().on(\"keydown\", (e) => {\n  if (e.key === \"Escape\") {\n    const openModal = me(\"dialog.modal.open\");\n    if (openModal) {\n      halt(e);\n      closeModal(openModal);\n    }\n  }\n});\n\nfunction createIcon(iconString) {\n  const domParser = new DOMParser();\n  const iconDoc = domParser.parseFromString(iconString, \"image/svg+xml\");\n  return iconDoc.documentElement;\n}\n\nfunction addUrlParams(url, params) {\n  let queryParams = new URLSearchParams(url.search);\n\n  params.forEach((param) => {\n    queryParams.set(param.name, param.value);\n  });\n  url.search = queryParams.toString();\n\n  // Use pushState to update the browser URL without reloading the page\n  window.history.pushState({ path: url.href }, \"\", url.href);\n}\n\nfunction openModal(el) {\n  me(\"body\").styles({ overflow: \"hidden\" });\n  el.classRemove(\"hidden\");\n  el.classAdd(\"open\");\n  el.showModal();\n}\n\nfunction closeModal(el) {\n  el.classRemove(\"open\");\n  el.classAdd(\"hidden\");\n  setTimeout(() => {\n    me(\"body\").styles({ overflow: \"auto\" });\n    el.close();\n  }, 250);\n}\n"
var _Assets4ced01a70d4c76f05f76b3ea328abe2356b805de = "me {\n  height: 60px;\n  background-color: var(--primary-color);\n}\n\nme .nav_links-container {\n  width: 100%;\n  height: 100%;\n  display: flex;\n  justify-content: space-between;\n  align-items: center;\n  padding: var(--spacing);\n}\n\nme .nav_links {\n  display: flex;\n  justify-content: space-between;\n  align-items: center;\n  height: 100%;\n}\n\nme .logo-container {\n  height: 100%;\n  padding: var(--spacing);\n}\n\nme .logo-container img {\n  height: 100%;\n}\n"
var _Assets8466ed20315f3bc808e945b0b594f15aad744e3c = "me {\n  background: var(--background-color-grey);\n  padding: var(--spacing-xs) var(--spacing-sm);\n  font-size: var(--font-size-sm);\n}\n\nme ol {\n  margin: 0;\n  padding: 0;\n  list-style: none;\n}\n\nme ol li {\n  display: inline-block;\n  margin-right: var(--spacing-xs);\n}\n\nme ol li:last-child {\n  margin-right: 0;\n}\n\nme ol li a {\n  color: var(--color-grey);\n  text-decoration: none;\n}\n"

// Assets returns go-assets FileSystem
var Assets = assets.NewFileSystem(map[string][]string{"/": []string{"layout.css", "footer.css", "global.js", "navbar.css", "breadcrumbs.css"}}, map[string]*assets.File{
	"/": &assets.File{
		Path:     "/",
		FileMode: 0x800001fd,
		Mtime:    time.Unix(1698772547, 1698772547473844933),
		Data:     nil,
	}, "/layout.css": &assets.File{
		Path:     "/layout.css",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1697643194, 1697643194260570599),
		Data:     []byte(_Assetsff9ed4d1c0d5f6a8ed200a512336471201c6127f),
	}, "/footer.css": &assets.File{
		Path:     "/footer.css",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1697643194, 1697643194260570599),
		Data:     []byte(_Assetsf2962b763da34337e494322f4ff6ea85a3cb804b),
	}, "/global.js": &assets.File{
		Path:     "/global.js",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1704910432, 1704910432597589262),
		Data:     []byte(_Assetsac50f5ae9fbaaf0afb9e8e49488ce346bbb994ed),
	}, "/navbar.css": &assets.File{
		Path:     "/navbar.css",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1704738544, 1704738544314758505),
		Data:     []byte(_Assets4ced01a70d4c76f05f76b3ea328abe2356b805de),
	}, "/breadcrumbs.css": &assets.File{
		Path:     "/breadcrumbs.css",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1704822351, 1704822351767073331),
		Data:     []byte(_Assets8466ed20315f3bc808e945b0b594f15aad744e3c),
	}}, "")
