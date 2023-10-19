package components

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _Assetsa7729f94837aa0817093c5f73c6a846500ab4f72 = "me {\n  background-color: var(--background-color-content);\n  border-radius: var(--border-radius-md);\n  border: var(--card-border-width) solid var(--border-color);\n  padding: var(--spacing-md);\n}\n"
var _Assets7a5e77e4cb373f48812845f452fe0498ca567170 = "/* BUTTON */\nme,\nme.primary {\n  border: var(--button-border-width) solid var(--primary-color);\n  border-radius: var(--border-radius);\n  color: var(--primary-color);\n  background-color: var(--background-color-content);\n  padding: var(--spacing-md) var(--spacing-lg);\n  transition: color 0.2s ease-in-out, background-color 0.2s ease-in-out;\n  cursor: pointer;\n}\n\nme[data-loading=\"true\"] .loading-spinner {\n  margin-right: var(--spacing-sm);\n}\n\nme.secondary {\n  border: var(--button-border-width) solid var(--secondary-color);\n  color: var(--secondary-color);\n}\n\nme.danger {\n  border: var(--button-border-width) solid var(--danger-color);\n  color: var(--danger-color);\n}\n\nme.success {\n  border: var(--button-border-width) solid var(--success-color);\n  color: var(--success-color);\n}\n\nme.warning {\n  border: var(--button-border-width) solid var(--warning-color);\n  color: var(--warning-color);\n}\n\nme:hover,\nme.primary:hover {\n  background-color: var(--primary-color);\n  color: var(--text-color-contrast);\n}\n\nme.secondary:hover {\n  background-color: var(--secondary-color);\n  color: var(--text-color-contrast);\n}\n\nme.danger:hover {\n  background-color: var(--danger-color);\n  color: var(--text-color-contrast);\n}\n\nme.success:hover {\n  background-color: var(--success-color);\n  color: var(--text-color-contrast);\n}\n\nme.warning:hover {\n  background-color: var(--warning-color);\n  color: var(--text-color-contrast);\n}\n\nme:disabled,\nme:disabled:hover {\n  background-color: var(--disabled-color);\n  border-color: var(--disabled-color);\n  color: var(--text-color-light);\n  cursor: not-allowed;\n}\n\nme.small {\n  padding: var(--spacing-sm) var(--spacing-md);\n}\n\nme.large {\n  padding: var(--spacing-lg) var(--spacing-xl);\n}\n\nme[data-icon=\"true\"],\nme[data-loading=\"true\"] {\n  display: flex;\n  justify-content: center;\n  align-items: center;\n}\n\nme .icon {\n  width: 20px;\n  height: 20px;\n  vertical-align: middle;\n}\n"
var _Assets835ab29723abda1b6c53883e64bdb0a0df16131c = ""
var _Assets0fe6beb7fb3b7b64cd71ef16e3f8c7fd4c6954c4 = "/* Loading spinner */\nme {\n  width: 50px;\n  height: 50px;\n  border-radius: 50%;\n  background: conic-gradient(#0000 10%, var(--primary-color));\n  mask: radial-gradient(farthest-side, #0000 calc(100% - 8px), #000 0);\n  -webkit-mask: radial-gradient(farthest-side, #0000 calc(100% - 8px), #000 0);\n  animation: s3 1s infinite linear;\n}\n\nme.xs {\n  width: 20px;\n  height: 20px;\n  mask: radial-gradient(farthest-side, #0000 calc(100% - 2px), #000 0);\n  -webkit-mask: radial-gradient(farthest-side, #0000 calc(100% - 2px), #000 0);\n}\n\nme.sm {\n  width: 35px;\n  height: 35px;\n  mask: radial-gradient(farthest-side, #0000 calc(100% - 4px), #000 0);\n  -webkit-mask: radial-gradient(farthest-side, #0000 calc(100% - 4px), #000 0);\n}\n\nme.lg {\n  width: 75px;\n  height: 75px;\n  mask: radial-gradient(farthest-side, #0000 calc(100% - 12px), #000 0);\n  -webkit-mask: radial-gradient(farthest-side, #0000 calc(100% - 12px), #000 0);\n}\n\n@keyframes s3 {\n  to {\n    transform: rotate(1turn);\n  }\n}\n"

// Assets returns go-assets FileSystem
var Assets = assets.NewFileSystem(map[string][]string{"/": []string{"Button.css", "placeholder.asset", "LoadingSpinner.css", "Card.css"}}, map[string]*assets.File{
	"/placeholder.asset": &assets.File{
		Path:     "/placeholder.asset",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1697643194, 1697643194260570599),
		Data:     []byte(_Assets835ab29723abda1b6c53883e64bdb0a0df16131c),
	}, "/LoadingSpinner.css": &assets.File{
		Path:     "/LoadingSpinner.css",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1697736664, 1697736664752156735),
		Data:     []byte(_Assets0fe6beb7fb3b7b64cd71ef16e3f8c7fd4c6954c4),
	}, "/Card.css": &assets.File{
		Path:     "/Card.css",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1697736645, 1697736645190108638),
		Data:     []byte(_Assetsa7729f94837aa0817093c5f73c6a846500ab4f72),
	}, "/": &assets.File{
		Path:     "/",
		FileMode: 0x800001fd,
		Mtime:    time.Unix(1697736495, 1697736495051303064),
		Data:     nil,
	}, "/Button.css": &assets.File{
		Path:     "/Button.css",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1697648034, 1697648034021572542),
		Data:     []byte(_Assets7a5e77e4cb373f48812845f452fe0498ca567170),
	}}, "")
