package components

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _Assets7a5e77e4cb373f48812845f452fe0498ca567170 = "/* BUTTON */\nme,\nme.primary {\n  border: var(--button-border-width) solid var(--primary-color);\n  border-radius: var(--border-radius);\n  color: var(--primary-color);\n  background-color: var(--background-color-content);\n  padding: var(--spacing-md) var(--spacing-lg);\n  transition: color 0.2s ease-in-out, background-color 0.2s ease-in-out;\n  cursor: pointer;\n}\n\nme[data-loading=\"true\"] .loading-spinner {\n  margin-right: var(--spacing-sm);\n}\n\nme.secondary {\n  border: var(--button-border-width) solid var(--secondary-color);\n  color: var(--secondary-color);\n}\n\nme.danger {\n  border: var(--button-border-width) solid var(--danger-color);\n  color: var(--danger-color);\n}\n\nme.success {\n  border: var(--button-border-width) solid var(--success-color);\n  color: var(--success-color);\n}\n\nme.warning {\n  border: var(--button-border-width) solid var(--warning-color);\n  color: var(--warning-color);\n}\n\nme:hover,\nme.primary:hover {\n  background-color: var(--primary-color);\n  color: var(--text-color-contrast);\n}\n\nme.secondary:hover {\n  background-color: var(--secondary-color);\n  color: var(--text-color-contrast);\n}\n\nme.danger:hover {\n  background-color: var(--danger-color);\n  color: var(--text-color-contrast);\n}\n\nme.success:hover {\n  background-color: var(--success-color);\n  color: var(--text-color-contrast);\n}\n\nme.warning:hover {\n  background-color: var(--warning-color);\n  color: var(--text-color-contrast);\n}\n\nme:disabled,\nme:disabled:hover {\n  background-color: var(--disabled-color);\n  border-color: var(--disabled-color);\n  color: var(--text-color-light);\n  cursor: not-allowed;\n}\n\nme.small {\n  padding: var(--spacing-sm) var(--spacing-md);\n}\n\nme.large {\n  padding: var(--spacing-lg) var(--spacing-xl);\n}\n\nme[data-icon=\"true\"],\nme[data-loading=\"true\"] {\n  display: flex;\n  justify-content: center;\n  align-items: center;\n}\n\nme .icon {\n  width: 20px;\n  height: 20px;\n  vertical-align: middle;\n}\n"
var _Assets092c1faf43196e909b1075f10477cb34abc7605c = "/* Statistic Component */\nme {\n  max-width: 300px;\n  display: flex;\n  flex-direction: column;\n  align-items: center;\n  flex: 1;\n  padding: var(--spacing-xl);\n  border-radius: var(--border-radius);\n  background-color: var(--background-color-content);\n  color: var(--text-color);\n  transition: border-color 0.2s ease-in-out, box-shadow 0.2s ease-in-out;\n  border: 1.5px solid var(--border-color);\n  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);\n}\n\nme:hover {\n  border-color: var(--border-hover-color);\n  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);\n}\n\nme .stat-element {\n  display: flex;\n  flex-direction: column;\n  align-items: flex-start;\n  margin-bottom: var(--spacing-lg);\n}\n\nme .stat-heading {\n  font-size: var(--font-size-md);\n  color: var(--text-color-light);\n  margin-bottom: var(--spacing-xs);\n}\n\nme .stat-value {\n  display: flex;\n  justify-content: center;\n  align-items: center;\n}\n\nme .stat-value .icon {\n  width: 30px;\n  height: 30px;\n  margin-right: var(--spacing-xs);\n}\n\nme .stat-value span {\n  font-size: var(--font-size-lg);\n  font-weight: bold;\n  text-align: left;\n}\n"
var _Assets65ccf34f5260672bf1cf91e4b34ddfc7451bd453 = "/* Modal */\nme {\n  display: flex;\n  flex-direction: column;\n  width: 500px;\n  max-width: 100%;\n  min-height: 250px;\n  padding: var(--spacing-lg);\n  background: var(--background-color-content);\n  color: var(--text-color);\n  border-radius: var(--border-radius);\n  transform-origin: center;\n  pointer-events: none;\n  transition: opacity 0.2s ease-in-out, transform 0.35s ease-in-out;\n  border: none;\n}\n\nme:focus-within {\n  outline: none;\n}\n\n.open {\n  opacity: 1;\n  transform: scale(1);\n}\n\n.hidden {\n  opacity: 0;\n  transform: scale(0);\n}\n\nme[open] {\n  pointer-events: inherit;\n}\n\nme::backdrop {\n  background: rgba(0, 0, 0, 0.7);\n  backdrop-filter: blur(5px);\n}\n\nme .modal-header {\n  margin-bottom: var(--spacing-lg);\n  display: flex;\n  justify-content: space-between;\n  align-items: center;\n}\n\nme .close-btn {\n  cursor: pointer;\n  display: flex;\n  justify-content: center;\n  align-items: center;\n  width: 30px;\n  height: 30px;\n  border-radius: 50%;\n  transition: background-color 0.2s ease-in-out, color 0.2s ease-in-out;\n}\n\nme .close-btn .icon {\n  width: 20px;\n  height: 20px;\n  fill: black;\n}\n\nme .close-btn:hover {\n  background-color: var(--border-hover-color);\n}\n\nme .close-btn:hover .icon {\n  fill: white;\n}\n\nme .modal-content {\n  margin-bottom: var(--spacing-lg);\n  flex: 1;\n}\n"
var _Assets8bd9efe3e913f3cfc8c2405a1ef3fab1c7e66a1f = "/* Tooltip */\nme {\n  display: inline-block;\n  position: relative; /* making the .tooltip span a container for the tooltip text */\n}\n\nme::before {\n  z-index: 1;\n  content: attr(data-content);\n  position: absolute;\n\n  top: 50%;\n  transform: translateY(-50%);\n\n  /* move to right */\n  right: initial;\n  left: 100%;\n  margin-left: var(--spacing-md);\n\n  /* basic styles */\n  width: 200px;\n  padding: var(--spacing);\n  border-radius: var(--border-radius);\n  background: #000;\n  color: #fff;\n  text-align: center;\n  visibility: hidden;\n  opacity: 0;\n  transition: opacity 0.2s ease-in-out, visibility 0.2s ease-in-out;\n  pointer-events: none;\n}\n\nme.left::before {\n  left: initial;\n  margin: initial;\n  right: 100%;\n  margin-right: var(--spacing-md);\n}\n\nme.top::before {\n  top: initial;\n  bottom: 100%;\n  margin-bottom: var(--spacing-md);\n  transform: translateX(-50%);\n}\n\nme.bottom::before {\n  top: 100%;\n  bottom: initial;\n  margin-top: var(--spacing-md);\n  transform: translateX(-50%);\n}\n\nme::after {\n  content: \"\";\n  position: absolute;\n\n  /* position tooltip correctly */\n  left: 100%;\n  margin-left: -5px;\n\n  /* vertically center */\n  top: 50%;\n  transform: translateY(-50%);\n\n  /* the arrow */\n  border: 10px solid #000;\n  border-color: transparent black transparent transparent;\n\n  visibility: hidden;\n  opacity: 0;\n  transition: opacity 100ms ease-in-out, visibility 100ms ease-in-out;\n  pointer-events: none;\n}\n\nme.left::after {\n  left: initial;\n  right: 100%;\n  margin-right: -5px;\n  border-color: transparent transparent transparent black;\n}\n\nme.top::after {\n  top: initial;\n  bottom: 100%;\n  margin-bottom: -5px;\n  transform: translateX(-50%);\n  border-color: black transparent transparent transparent;\n}\n\nme.bottom::after {\n  top: 100%;\n  bottom: initial;\n  margin-top: -5px;\n  transform: translateX(-50%);\n  border-color: transparent transparent black transparent;\n}\n\nme:hover::before,\nme:hover::after {\n  display: inline-block;\n  visibility: visible;\n  opacity: 1;\n}\n"
var _Assets835ab29723abda1b6c53883e64bdb0a0df16131c = ""
var _Assets0fe6beb7fb3b7b64cd71ef16e3f8c7fd4c6954c4 = "/* Loading spinner */\nme {\n  width: 35px;\n  height: 35px;\n  mask: radial-gradient(farthest-side, #0000 calc(100% - 4px), #000 0);\n  -webkit-mask: radial-gradient(farthest-side, #0000 calc(100% - 4px), #000 0);\n  animation: s3 1s infinite linear;\n  border-radius: 50%;\n  background: conic-gradient(#0000 10%, var(--primary-color));\n}\n\nme.sm {\n  width: 20px;\n  height: 20px;\n  mask: radial-gradient(farthest-side, #0000 calc(100% - 2px), #000 0);\n  -webkit-mask: radial-gradient(farthest-side, #0000 calc(100% - 2px), #000 0);\n}\n\nme.lg {\n  width: 50px;\n  height: 50px;\n  mask: radial-gradient(farthest-side, #0000 calc(100% - 8px), #000 0);\n  -webkit-mask: radial-gradient(farthest-side, #0000 calc(100% - 8px), #000 0);\n}\n\nme.xl {\n  width: 75px;\n  height: 75px;\n  mask: radial-gradient(farthest-side, #0000 calc(100% - 12px), #000 0);\n  -webkit-mask: radial-gradient(farthest-side, #0000 calc(100% - 12px), #000 0);\n}\n\n@keyframes s3 {\n  to {\n    transform: rotate(1turn);\n  }\n}\n"
var _Assetsb9e42c3084e4aef44f97be2ecbeff94dfb111569 = "me().on(\"mouseenter\", (ev) => {\n  const tooltip = ev.target;\n  const tooltipWidth = ev.target.offsetWidth;\n  const tooltipHeight = ev.target.offsetHeight;\n  const rect = tooltip.getBoundingClientRect();\n  // if (rect.left < tooltipWidth / 2) {\n  //   tooltip.styles({ left: \"0\" });\n  // } else if (window.innerWidth - rect.right < tooltipWidth / 2) {\n  //   tooltip.styles({ left: \"auto\", right: \"0\" });\n  // }\n  // if (rect.top < tooltipHeight) {\n  //   tooltip.style.top = \"0\";\n  // }\n});\n"
var _Assetsa7729f94837aa0817093c5f73c6a846500ab4f72 = "me {\n  background-color: var(--background-color-content);\n  border-radius: var(--border-radius-md);\n  border: var(--card-border-width) solid var(--border-color);\n  padding: var(--spacing-md);\n}\n"
var _Assets3ba904d5d26b984cac0326a2c9eb7591076cc34b = "// Modal Logic\nfunction closeModal() {\n  me(\".modal\").classRemove(\"open\");\n  me(\".modal\").classAdd(\"hidden\");\n  setTimeout(() => {\n    me(\".modal\").styles({ overflow: \"auto\" });\n    me(\".modal\").close();\n  }, 250);\n}\n\nme(\"#open-modal\").on(\"click\", () => {\n  me(\"body\").styles({ overflow: \"hidden\" });\n  me(\".modal\").classRemove(\"hidden\");\n  me(\".modal\").classAdd(\"open\");\n  me(\".modal\").showModal();\n});\n\nme(\"#close-btn\").on(\"click\", () => {\n  closeModal();\n});\n\nme(\"body\").on(\"keydown\", (e) => {\n  halt(e);\n  if (e.key === \"Escape\") {\n    closeModal();\n  }\n});\n"

// Assets returns go-assets FileSystem
var Assets = assets.NewFileSystem(map[string][]string{"/": []string{"Button.css", "Statistic.css", "placeholder.asset", "LoadingSpinner.css", "Tooltip.js", "Modal.css", "Tooltip.css", "Card.css", "Modal.js"}}, map[string]*assets.File{
	"/": &assets.File{
		Path:     "/",
		FileMode: 0x800001fd,
		Mtime:    time.Unix(1698158560, 1698158560915022808),
		Data:     nil,
	}, "/Button.css": &assets.File{
		Path:     "/Button.css",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1697648034, 1697648034021572542),
		Data:     []byte(_Assets7a5e77e4cb373f48812845f452fe0498ca567170),
	}, "/Statistic.css": &assets.File{
		Path:     "/Statistic.css",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1698156105, 1698156105762748615),
		Data:     []byte(_Assets092c1faf43196e909b1075f10477cb34abc7605c),
	}, "/Modal.css": &assets.File{
		Path:     "/Modal.css",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1698250367, 1698250367836375662),
		Data:     []byte(_Assets65ccf34f5260672bf1cf91e4b34ddfc7451bd453),
	}, "/Tooltip.css": &assets.File{
		Path:     "/Tooltip.css",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1698159550, 1698159550113092816),
		Data:     []byte(_Assets8bd9efe3e913f3cfc8c2405a1ef3fab1c7e66a1f),
	}, "/Modal.js": &assets.File{
		Path:     "/Modal.js",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1698244563, 1698244563881526065),
		Data:     []byte(_Assets3ba904d5d26b984cac0326a2c9eb7591076cc34b),
	}, "/placeholder.asset": &assets.File{
		Path:     "/placeholder.asset",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1697643194, 1697643194260570599),
		Data:     []byte(_Assets835ab29723abda1b6c53883e64bdb0a0df16131c),
	}, "/LoadingSpinner.css": &assets.File{
		Path:     "/LoadingSpinner.css",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1698249131, 1698249131628670280),
		Data:     []byte(_Assets0fe6beb7fb3b7b64cd71ef16e3f8c7fd4c6954c4),
	}, "/Tooltip.js": &assets.File{
		Path:     "/Tooltip.js",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1698158448, 1698158448282967169),
		Data:     []byte(_Assetsb9e42c3084e4aef44f97be2ecbeff94dfb111569),
	}, "/Card.css": &assets.File{
		Path:     "/Card.css",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1697736645, 1697736645190108638),
		Data:     []byte(_Assetsa7729f94837aa0817093c5f73c6a846500ab4f72),
	}}, "")
