package views

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _Assets02d14f41a13038998aa3a65ac4573c54fa0d9201 = "(function () {\n  const localMe = generateMe(me());\n  const modal = localMe(\"dialog.modal\");\n\n  // Modal Logic\n  localMe(\"#open-modal\").on(\"click\", () => {\n    openModal(modal);\n  });\n\n  localMe(\"#close-btn\").on(\"click\", () => {\n    closeModal(modal);\n  });\n})();\n"
var _Assetse7f942629d571e16ac35ae6ea69a0d512337b2b5 = "me {\n  display: flex;\n  justify-content: center;\n}\n\nme form {\n  max-width: 400px;\n}\n"
var _Assetsb9637c000a10ff4204b41d99efe26b3b7726d123 = "(function () {\n  const form = me(\"form\");\n\n  const reviver = (key, value) => {\n    // Check if the value is a string and represents an integer\n    if (typeof value === \"string\" && /^\\d+$/.test(value)) {\n      return parseInt(value, 10);\n    }\n\n    // If not, return the original value\n    return value;\n  };\n\n  form.on(\"submit\", (e) => {\n    e.preventDefault();\n    // capture the form data\n    const formData = new FormData(form);\n    // convert the form data into a JSON object\n    const data = Object.fromEntries(formData);\n    // Get the multi-select values\n    const multiSelectValues = JSON.parse(data[\"multi-select\"], reviver);\n    const multiSearchSelectValues = JSON.parse(\n      data[\"multi-search-select\"],\n      reviver\n    );\n    // Remove the multi-select from the form data\n    delete data[\"multi-select\"];\n    delete data[\"multi-search-select\"];\n    // Add the multi-select values to the form data\n    data[\"multi-select\"] = multiSelectValues;\n    data[\"multi-search-select\"] = multiSearchSelectValues;\n    // Log the form data\n    console.log(data);\n  });\n})();\n"
var _Assets99175b3d1c70661560390b724da3c17ec5fe3b80 = "me {\n  display: flex;\n  flex-direction: column;\n  align-items: center;\n  justify-content: center;\n  font-family: \"Courier New\", Courier, monospace;\n}\n\nme h1 {\n  margin-bottom: var(--spacing-xl);\n  background: linear-gradient(\n    90deg,\n    var(--primary-color) 0%,\n    var(--secondary-color) 100%\n  );\n  -webkit-background-clip: text;\n  -webkit-text-fill-color: transparent;\n}\n\nme .card {\n  width: 100%;\n  max-width: 400px;\n  margin: 0 auto;\n  padding: var(--spacing-xl);\n  border-radius: var(--border-radius);\n  box-shadow: 0 0 10px rgba(0, 0, 0, 0.2);\n}\n\nme .card .row {\n  display: flex;\n  align-items: center;\n}\n\nme .card .img-container {\n  width: 100px;\n  height: 100px;\n  border-radius: 50%;\n  padding: var(--spacing);\n  overflow: hidden;\n  margin-right: var(--spacing-xl);\n}\n\n:root[data-theme=\"light\"] me .card .img-container {\n  background-color: black;\n  box-shadow: 0 0 10px rgba(0, 0, 0, 0.2);\n}\n\nme .card .img-container img {\n  width: 100%;\n  height: 100%;\n  object-fit: contain;\n}\n\nme .card .text-404 {\n  font-size: var(--font-size-xl);\n  text-align: center;\n  color: var(--primary-color);\n}\n\nme .card .home-btn {\n  margin-top: var(--spacing);\n  display: flex;\n  align-items: center;\n  justify-content: center;\n  text-decoration: none;\n}\n"
var _Assets6850ef0400e0c476ced907d92d4d46801c8e7e7e = "me .add-button-container {\n  display: flex;\n  justify-content: flex-end;\n}\n\nme .table-container {\n  margin-top: var(--spacing-lg);\n}\n"
var _Assetsfe9c4a63690d43799c11c81c0fd46347b4198a48 = "me {\n  padding-top: var(--spacing-xl);\n  display: flex;\n  justify-content: center;\n}\n\nme .container {\n  width: 100%;\n  max-width: 400px;\n  display: flex;\n  flex-direction: column;\n  align-items: center;\n}\n\n"
var _Assets0679cb2ba52be1dce6444f34ae990b5bcde2fe29 = "me .edit-button-container {\n  display: flex;\n  justify-content: flex-end;\n}\n\nme .properties-grid {\n  display: inline-grid;\n  grid-template-columns: repeat(2, auto);\n  grid-gap: var(--spacing-lg);\n}\n\n"
var _Assets1d5cd818f95e939cff881702655b2b1cc6e3825e = "me {\n  display: flex;\n  justify-content: center;\n}\n\nme form {\n  max-width: 400px;\n}\n"

// Assets returns go-assets FileSystem
var Assets = assets.NewFileSystem(map[string][]string{"/": []string{"Users.css", "Login.css", "User.css", "Index.js", "AddUser.css", "EditUser.css", "Form.js", "404.css"}}, map[string]*assets.File{
	"/User.css": &assets.File{
		Path:     "/User.css",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1706218940, 1706218940377227250),
		Data:     []byte(_Assets0679cb2ba52be1dce6444f34ae990b5bcde2fe29),
	}, "/EditUser.css": &assets.File{
		Path:     "/EditUser.css",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1706218940, 1706218940377227250),
		Data:     []byte(_Assets1d5cd818f95e939cff881702655b2b1cc6e3825e),
	}, "/AddUser.css": &assets.File{
		Path:     "/AddUser.css",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1706218940, 1706218940377227250),
		Data:     []byte(_Assetse7f942629d571e16ac35ae6ea69a0d512337b2b5),
	}, "/Form.js": &assets.File{
		Path:     "/Form.js",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1706218940, 1706218940377227250),
		Data:     []byte(_Assetsb9637c000a10ff4204b41d99efe26b3b7726d123),
	}, "/404.css": &assets.File{
		Path:     "/404.css",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1706218940, 1706218940377227250),
		Data:     []byte(_Assets99175b3d1c70661560390b724da3c17ec5fe3b80),
	}, "/": &assets.File{
		Path:     "/",
		FileMode: 0x800001fd,
		Mtime:    time.Unix(1707132617, 1707132617866661574),
		Data:     nil,
	}, "/Users.css": &assets.File{
		Path:     "/Users.css",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1706218940, 1706218940377227250),
		Data:     []byte(_Assets6850ef0400e0c476ced907d92d4d46801c8e7e7e),
	}, "/Login.css": &assets.File{
		Path:     "/Login.css",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1706218940, 1706218940377227250),
		Data:     []byte(_Assetsfe9c4a63690d43799c11c81c0fd46347b4198a48),
	}, "/Index.js": &assets.File{
		Path:     "/Index.js",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1706218940, 1706218940377227250),
		Data:     []byte(_Assets02d14f41a13038998aa3a65ac4573c54fa0d9201),
	}}, "")
