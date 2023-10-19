package layout

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _Assets8466ed20315f3bc808e945b0b594f15aad744e3c = "me {\n\tbackground: var(--background-color-grey);\n\tpadding: var(--spacing-xs) var(--spacing-sm);\n\tfont-size: var(--font-size-sm);\n}\n\nme ol {\n\tdisplay: flex;\n}\n\nme ol li {\n\tcolor: var(--text-color-light) !important;\n}\n\n/* me ol li a:visted { */\n/* \tcolor: var(--text-color-light) !important; */\n/* } */\n\nme ol li .separator {\n\tcolor: var(--text-color-light);\n\tpadding: 0 var(--spacing-sm);\n}\n"
var _Assetsff9ed4d1c0d5f6a8ed200a512336471201c6127f = "/* body */\nme {\n  display: flex;\n  flex-direction: column;\n  min-height: 100vh;\n  background-color: var(--background-color);\n}\n\nme main {\n  flex-grow: 1;\n}\n"
var _Assetsf2962b763da34337e494322f4ff6ea85a3cb804b = "me {\n  height: 30px;\n  display: flex;\n  align-items: center;\n  justify-content: center;\n  color: var(--text-color-light);\n  font-size: var(--font-size-sm);\n}\n\nme sup {\n  font-size: var(--font-size-xs);\n}\n"
var _Assets4ced01a70d4c76f05f76b3ea328abe2356b805de = "me {\n  height: 60px;\n  display: flex;\n  align-items: center;\n  justify-content: space-between;\n  background-color: var(--primary-color);\n}\n\nme .logo-container {\n  height: 100%;\n  padding: var(--spacing);\n}\n\nme .logo-container img {\n  height: 100%;\n  width: auto;\n}\n"

// Assets returns go-assets FileSystem
var Assets = assets.NewFileSystem(map[string][]string{"/": []string{"layout.css", "footer.css", "navbar.css", "breadcrumbs.css"}}, map[string]*assets.File{
	"/footer.css": &assets.File{
		Path:     "/footer.css",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1697643194, 1697643194260570599),
		Data:     []byte(_Assetsf2962b763da34337e494322f4ff6ea85a3cb804b),
	}, "/navbar.css": &assets.File{
		Path:     "/navbar.css",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1697643194, 1697643194260570599),
		Data:     []byte(_Assets4ced01a70d4c76f05f76b3ea328abe2356b805de),
	}, "/breadcrumbs.css": &assets.File{
		Path:     "/breadcrumbs.css",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1697643194, 1697643194260570599),
		Data:     []byte(_Assets8466ed20315f3bc808e945b0b594f15aad744e3c),
	}, "/": &assets.File{
		Path:     "/",
		FileMode: 0x800001fd,
		Mtime:    time.Unix(1697643194, 1697643194260570599),
		Data:     nil,
	}, "/layout.css": &assets.File{
		Path:     "/layout.css",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1697643194, 1697643194260570599),
		Data:     []byte(_Assetsff9ed4d1c0d5f6a8ed200a512336471201c6127f),
	}}, "")
