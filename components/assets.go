package components

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _Assets7a5e77e4cb373f48812845f452fe0498ca567170 = "/* BUTTON */\nme,\nme.primary {\n  border: var(--button-border-width) solid var(--primary-color);\n  border-radius: var(--border-radius);\n  color: var(--primary-color);\n  background-color: var(--background-color-content);\n  padding: var(--spacing-md) var(--spacing-lg);\n  transition: color 0.2s ease-in-out, background-color 0.2s ease-in-out;\n  cursor: pointer;\n}\n\nme[data-loading=\"true\"] .loading-spinner {\n  margin-right: var(--spacing-sm);\n}\n\nme.secondary {\n  border: var(--button-border-width) solid var(--secondary-color);\n  color: var(--secondary-color);\n}\n\nme.danger {\n  border: var(--button-border-width) solid var(--danger-color);\n  color: var(--danger-color);\n}\n\nme.success {\n  border: var(--button-border-width) solid var(--success-color);\n  color: var(--success-color);\n}\n\nme.warning {\n  border: var(--button-border-width) solid var(--warning-color);\n  color: var(--warning-color);\n}\n\nme:hover,\nme.primary:hover {\n  background-color: var(--primary-color);\n  color: var(--text-color-contrast);\n}\n\nme.secondary:hover {\n  background-color: var(--secondary-color);\n  color: var(--text-color-contrast);\n}\n\nme.danger:hover {\n  background-color: var(--danger-color);\n  color: var(--text-color-contrast);\n}\n\nme.success:hover {\n  background-color: var(--success-color);\n  color: var(--text-color-contrast);\n}\n\nme.warning:hover {\n  background-color: var(--warning-color);\n  color: var(--text-color-contrast);\n}\n\nme:disabled,\nme:disabled:hover {\n  background-color: var(--disabled-color);\n  border-color: var(--disabled-color);\n  color: var(--text-color-light);\n  cursor: not-allowed;\n}\n\nme.small {\n  padding: var(--spacing-sm) var(--spacing-md);\n}\n\nme.large {\n  padding: var(--spacing-lg) var(--spacing-xl);\n}\n\nme[data-icon=\"true\"],\nme[data-loading=\"true\"] {\n  display: flex;\n  justify-content: center;\n  align-items: center;\n}\n\nme .icon {\n  width: 20px;\n  height: 20px;\n  vertical-align: middle;\n}\n"
var _Assets092c1faf43196e909b1075f10477cb34abc7605c = "/* Statistic Component */\nme {\n  max-width: 300px;\n  display: flex;\n  flex-direction: column;\n  align-items: center;\n  flex: 1;\n  padding: var(--spacing-xl);\n  border-radius: var(--border-radius);\n  background-color: var(--background-color-content);\n  color: var(--text-color);\n  transition: border-color 0.2s ease-in-out, box-shadow 0.2s ease-in-out;\n  border: 1.5px solid var(--border-color);\n  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);\n}\n\nme:hover {\n  border-color: var(--border-hover-color);\n  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);\n}\n\nme .stat-element {\n  display: flex;\n  flex-direction: column;\n  align-items: flex-start;\n  margin-bottom: var(--spacing-lg);\n}\n\nme .stat-heading {\n  font-size: var(--font-size-md);\n  color: var(--text-color-light);\n  margin-bottom: var(--spacing-xs);\n}\n\nme .stat-value {\n  display: flex;\n  justify-content: center;\n  align-items: center;\n}\n\nme .stat-value .icon {\n  width: 30px;\n  height: 30px;\n  margin-right: var(--spacing-xs);\n}\n\nme .stat-value span {\n  font-size: var(--font-size-lg);\n  font-weight: bold;\n  text-align: left;\n}\n"
var _Assets3ba904d5d26b984cac0326a2c9eb7591076cc34b = "(function () {\n  // Modal Logic\n  function closeModal() {\n    me(\".modal\").classRemove(\"open\");\n    me(\".modal\").classAdd(\"hidden\");\n    setTimeout(() => {\n      me(\"body\").styles({ overflow: \"auto\" });\n      me(\".modal\").close();\n    }, 250);\n  }\n\n  me(\"#open-modal\").on(\"click\", () => {\n    me(\"body\").styles({ overflow: \"hidden\" });\n    me(\".modal\").classRemove(\"hidden\");\n    me(\".modal\").classAdd(\"open\");\n    me(\".modal\").showModal();\n  });\n\n  me(\"#close-btn\").on(\"click\", () => {\n    closeModal();\n  });\n\n  me(\"body\").on(\"keydown\", (e) => {\n    halt(e);\n    if (e.key === \"Escape\") {\n      closeModal();\n    }\n  });\n})();\n"
var _Assets8350e932c9320eb7d70887a1b20e026f5d21b75b = "me {\n  display: flex;\n  flex-direction: column;\n  gap: var(--spacing);\n  margin-bottom: var(--spacing-lg);\n}\n\nme input {\n  border: 1.5px solid var(--border-color-dark);\n  border-radius: var(--border-radius);\n  background-color: var(--background-color-content);\n  padding: var(--spacing);\n  font-family: var(--font-family);\n  transition: border-color 0.2s ease-in-out, border-width 0.2s ease-in-out,\n    box-shadow 0.2s ease-in-out;\n  position: relative;\n  z-index: 1;\n}\n\nme input:hover {\n  border-color: var(--border-hover-color);\n}\n\nme input:focus-visible {\n  outline: none;\n  border-color: var(--border-hover-color);\n  box-shadow: 0 0 0.1px 0.5px var(--border-hover-color);\n}\n\nme input.sm {\n  font-size: var(--font-size-xs);\n}\n\nme input.md {\n  font-size: var(--font-size-sm);\n}\n\nme input.lg {\n  font-size: var(--font-size-md);\n}\n\nme input:disabled {\n  background-color: var(--disabled-color);\n  cursor: not-allowed;\n}\n"
var _Assets0fe6beb7fb3b7b64cd71ef16e3f8c7fd4c6954c4 = "/* Loading spinner */\nme {\n  width: 35px;\n  height: 35px;\n  mask: radial-gradient(farthest-side, #0000 calc(100% - 4px), #000 0);\n  -webkit-mask: radial-gradient(farthest-side, #0000 calc(100% - 4px), #000 0);\n  animation: s3 1s infinite linear;\n  border-radius: 50%;\n  background: conic-gradient(#0000 10%, var(--primary-color));\n}\n\nme.sm {\n  width: 20px;\n  height: 20px;\n  mask: radial-gradient(farthest-side, #0000 calc(100% - 2px), #000 0);\n  -webkit-mask: radial-gradient(farthest-side, #0000 calc(100% - 2px), #000 0);\n}\n\nme.lg {\n  width: 50px;\n  height: 50px;\n  mask: radial-gradient(farthest-side, #0000 calc(100% - 8px), #000 0);\n  -webkit-mask: radial-gradient(farthest-side, #0000 calc(100% - 8px), #000 0);\n}\n\nme.xl {\n  width: 75px;\n  height: 75px;\n  mask: radial-gradient(farthest-side, #0000 calc(100% - 12px), #000 0);\n  -webkit-mask: radial-gradient(farthest-side, #0000 calc(100% - 12px), #000 0);\n}\n\n@keyframes s3 {\n  to {\n    transform: rotate(1turn);\n  }\n}\n"
var _Assetsb9e42c3084e4aef44f97be2ecbeff94dfb111569 = "(function () {\n  me().on(\"mouseenter\", (ev) => {\n    const tooltip = ev.target;\n    const tooltipWidth = ev.target.offsetWidth;\n    const tooltipHeight = ev.target.offsetHeight;\n    const rect = tooltip.getBoundingClientRect();\n    // if (rect.left < tooltipWidth / 2) {\n    //   tooltip.styles({ left: \"0\" });\n    // } else if (window.innerWidth - rect.right < tooltipWidth / 2) {\n    //   tooltip.styles({ left: \"auto\", right: \"0\" });\n    // }\n    // if (rect.top < tooltipHeight) {\n    //   tooltip.style.top = \"0\";\n    // }\n  });\n})();\n"
var _Assets65ccf34f5260672bf1cf91e4b34ddfc7451bd453 = "/* Modal */\nme {\n  display: flex;\n  flex-direction: column;\n  width: 500px;\n  max-width: 100%;\n  min-height: 250px;\n  padding: var(--spacing-lg);\n  background: var(--background-color-content);\n  color: var(--text-color);\n  border-radius: var(--border-radius);\n  transform-origin: center;\n  pointer-events: none;\n  transition: opacity 0.2s ease-in-out, transform 0.35s ease-in-out;\n  border: none;\n}\n\nme:focus-within {\n  outline: none;\n}\n\n.open {\n  opacity: 1;\n  transform: scale(1);\n}\n\n.hidden {\n  opacity: 0;\n  transform: scale(0);\n}\n\nme[open] {\n  pointer-events: inherit;\n}\n\nme::backdrop {\n  background: rgba(0, 0, 0, 0.7);\n  backdrop-filter: blur(5px);\n}\n\nme .modal-header {\n  margin-bottom: var(--spacing-lg);\n  display: flex;\n  justify-content: space-between;\n  align-items: center;\n}\n\nme .close-btn {\n  cursor: pointer;\n  display: flex;\n  justify-content: center;\n  align-items: center;\n  width: 30px;\n  height: 30px;\n  border-radius: 50%;\n  transition: background-color 0.2s ease-in-out, color 0.2s ease-in-out;\n}\n\nme .close-btn .icon {\n  width: 20px;\n  height: 20px;\n  fill: black;\n}\n\nme .close-btn:hover {\n  background-color: var(--border-hover-color);\n}\n\nme .close-btn:hover .icon {\n  fill: white;\n}\n\nme .modal-content {\n  margin-bottom: var(--spacing-lg);\n  flex: 1;\n}\n"
var _Assets1eb0157acf7a8aaff323f61ae64000a30cae3278 = "/* Input helper */\nme {\n  margin: var(--spacing-sm);\n  font-size: var(--font-size-xs);\n}\n\nme.success {\n  color: var(--success-color);\n}\n\nme.warning {\n  color: var(--warning-color);\n}\n\nme.error {\n  color: var(--danger-color);\n}\n"
var _Assets900dec90ce34bda078405e1223b753edbd1dcb11 = "me {\n  border-radius: var(--border-radius-lg);\n  margin: var(--spacing) auto;\n  position: relative;\n}\n\nme .progress-bar {\n  width: 100%;\n  height: 12px;\n  border-radius: var(--border-radius-lg);\n  background-color: var(--border-color);\n  position: relative;\n  overflow: hidden;\n}\n\nme .progress-label {\n  position: absolute;\n  top: 50%;\n  left: 50%;\n  transform: translate(-50%, -50%);\n  font-size: var(--font-size-xs);\n  font-weight: bold;\n  z-index: 1;\n}\n\nme .progress {\n  position: absolute;\n  top: 0;\n  left: 0;\n  width: 0;\n  height: 100%;\n  background-color: var(--primary-color);\n  transition: width 0.2s ease-in-out;\n  border-radius: inherit;\n}\n\nme[data-type=\"success\"] .progress {\n  background-color: var(--success-color);\n}\n\nme[data-type=\"warning\"] .progress {\n  background-color: var(--warning-color);\n}\n\nme[data-type=\"danger\"] .progress {\n  background-color: var(--danger-color);\n}\n"
var _Assetsc8049e54d8c450e18295624c50106bf34a5bd072 = "me {\n  display: flex;\n  flex-direction: column;\n  gap: var(--spacing);\n  margin-bottom: var(--spacing-lg);\n}\n\nme textarea {\n  border: 1.5px solid var(--border-color-dark);\n  border-radius: var(--border-radius);\n  background-color: var(--background-color-content);\n  padding: var(--spacing);\n  font-family: var(--font-family);\n  resize: vertical;\n  transition: border-color 0.2s ease-in-out, border-width 0.2s ease-in-out,\n    box-shadow 0.2s ease-in-out;\n}\n\nme textarea:disabled {\n  background-color: var(--disabled-color);\n  cursor: not-allowed;\n}\n\nme textarea:focus-visible {\n  outline: none;\n  border-color: var(--primary-color);\n}\n"
var _Assets0b54eedfbd58bc6889720d70aa1d610bd11c4491 = "me {\n  display: flex;\n}\n\nme label {\n  margin-right: var(--spacing-sm);\n}\n\nme input[type=\"radio\"] {\n  /* Add if not using autoprefixer */\n  -webkit-appearance: none;\n  appearance: none;\n  /* For iOS < 15 to remove gradient background */\n  background-color: #fff;\n  /* Not removed via appearance */\n  margin: 0;\n  font: inherit;\n  color: var(--border-hover-color);\n  width: 1.15em;\n  height: 1.15em;\n  border: 0.15em solid currentColor;\n  border-radius: 50%;\n  display: grid;\n  place-content: center;\n}\n\nme input[type=\"radio\"]::before {\n  content: \"\";\n  width: 0.65em;\n  height: 0.65em;\n  border-radius: 50%;\n  transform: scale(0);\n  transition: 200ms transform ease-in-out;\n  transform-origin: center;\n  box-shadow: inset 1em 1em var(--border-hover-color);\n}\n\nme input[type=\"radio\"]:checked::before {\n  transform: scale(1) translateX(0.009em);\n}\n\nme input[type=\"radio\"]:focus {\n  outline: none;\n  border-color: var(--primary-color);\n}\n"
var _Assets835ab29723abda1b6c53883e64bdb0a0df16131c = ""
var _Assetsa7729f94837aa0817093c5f73c6a846500ab4f72 = "me {\n  background-color: var(--background-color-content);\n  border-radius: var(--border-radius-md);\n  border: var(--card-border-width) solid var(--border-color);\n  padding: var(--spacing-md);\n}\n"
var _Assets8bd9efe3e913f3cfc8c2405a1ef3fab1c7e66a1f = "/* Tooltip */\nme {\n  display: inline-block;\n  position: relative; /* making the .tooltip span a container for the tooltip text */\n}\n\nme::before {\n  z-index: 1;\n  content: attr(data-content);\n  position: absolute;\n\n  top: 50%;\n  transform: translateY(-50%);\n\n  /* move to right */\n  right: initial;\n  left: 100%;\n  margin-left: var(--spacing-md);\n\n  /* basic styles */\n  width: 200px;\n  padding: var(--spacing);\n  border-radius: var(--border-radius);\n  background: #000;\n  color: #fff;\n  text-align: center;\n  visibility: hidden;\n  opacity: 0;\n  transition: opacity 0.2s ease-in-out, visibility 0.2s ease-in-out;\n  pointer-events: none;\n}\n\nme.left::before {\n  left: initial;\n  margin: initial;\n  right: 100%;\n  margin-right: var(--spacing-md);\n}\n\nme.top::before {\n  top: initial;\n  bottom: 100%;\n  margin-bottom: var(--spacing-md);\n  transform: translateX(-50%);\n}\n\nme.bottom::before {\n  top: 100%;\n  bottom: initial;\n  margin-top: var(--spacing-md);\n  transform: translateX(-50%);\n}\n\nme::after {\n  content: \"\";\n  position: absolute;\n\n  /* position tooltip correctly */\n  left: 100%;\n  margin-left: -5px;\n\n  /* vertically center */\n  top: 50%;\n  transform: translateY(-50%);\n\n  /* the arrow */\n  border: 10px solid #000;\n  border-color: transparent black transparent transparent;\n\n  visibility: hidden;\n  opacity: 0;\n  transition: opacity 100ms ease-in-out, visibility 100ms ease-in-out;\n  pointer-events: none;\n}\n\nme.left::after {\n  left: initial;\n  right: 100%;\n  margin-right: -5px;\n  border-color: transparent transparent transparent black;\n}\n\nme.top::after {\n  top: initial;\n  bottom: 100%;\n  margin-bottom: -5px;\n  transform: translateX(-50%);\n  border-color: black transparent transparent transparent;\n}\n\nme.bottom::after {\n  top: 100%;\n  bottom: initial;\n  margin-top: -5px;\n  transform: translateX(-50%);\n  border-color: transparent transparent black transparent;\n}\n\nme:hover::before,\nme:hover::after {\n  display: inline-block;\n  visibility: visible;\n  opacity: 1;\n}\n"
var _Assetsa0084176b3c3ed34696e0c93f712ea4c10c94315 = "console.log(me().id);\nme().attribute(\"data-percentage\");\n"

// Assets returns go-assets FileSystem
var Assets = assets.NewFileSystem(map[string][]string{"/": []string{"Button.css", "Statistic.css", "Radio.css", "placeholder.asset", "LoadingSpinner.css", "Tooltip.js", "Modal.css", "Tooltip.css", "InputHelper.css", "Progress.css", "Card.css", "Progress.js", "Modal.js", "Textarea.css", "Input.css"}}, map[string]*assets.File{
	"/Button.css": &assets.File{
		Path:     "/Button.css",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1697648034, 1697648034021572542),
		Data:     []byte(_Assets7a5e77e4cb373f48812845f452fe0498ca567170),
	}, "/Statistic.css": &assets.File{
		Path:     "/Statistic.css",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1698156105, 1698156105762748615),
		Data:     []byte(_Assets092c1faf43196e909b1075f10477cb34abc7605c),
	}, "/Modal.js": &assets.File{
		Path:     "/Modal.js",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1698682295, 1698682295069148208),
		Data:     []byte(_Assets3ba904d5d26b984cac0326a2c9eb7591076cc34b),
	}, "/Input.css": &assets.File{
		Path:     "/Input.css",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1698683098, 1698683098620153921),
		Data:     []byte(_Assets8350e932c9320eb7d70887a1b20e026f5d21b75b),
	}, "/LoadingSpinner.css": &assets.File{
		Path:     "/LoadingSpinner.css",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1698249131, 1698249131628670280),
		Data:     []byte(_Assets0fe6beb7fb3b7b64cd71ef16e3f8c7fd4c6954c4),
	}, "/Tooltip.js": &assets.File{
		Path:     "/Tooltip.js",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1698682280, 1698682280505539554),
		Data:     []byte(_Assetsb9e42c3084e4aef44f97be2ecbeff94dfb111569),
	}, "/Modal.css": &assets.File{
		Path:     "/Modal.css",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1698250367, 1698250367836375662),
		Data:     []byte(_Assets65ccf34f5260672bf1cf91e4b34ddfc7451bd453),
	}, "/InputHelper.css": &assets.File{
		Path:     "/InputHelper.css",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1698676036, 1698676036060755427),
		Data:     []byte(_Assets1eb0157acf7a8aaff323f61ae64000a30cae3278),
	}, "/Progress.css": &assets.File{
		Path:     "/Progress.css",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1698257317, 1698257317247722623),
		Data:     []byte(_Assets900dec90ce34bda078405e1223b753edbd1dcb11),
	}, "/Textarea.css": &assets.File{
		Path:     "/Textarea.css",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1698684688, 1698684688078846389),
		Data:     []byte(_Assetsc8049e54d8c450e18295624c50106bf34a5bd072),
	}, "/Radio.css": &assets.File{
		Path:     "/Radio.css",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1698687352, 1698687352416475496),
		Data:     []byte(_Assets0b54eedfbd58bc6889720d70aa1d610bd11c4491),
	}, "/placeholder.asset": &assets.File{
		Path:     "/placeholder.asset",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1697643194, 1697643194260570599),
		Data:     []byte(_Assets835ab29723abda1b6c53883e64bdb0a0df16131c),
	}, "/Card.css": &assets.File{
		Path:     "/Card.css",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1697736645, 1697736645190108638),
		Data:     []byte(_Assetsa7729f94837aa0817093c5f73c6a846500ab4f72),
	}, "/": &assets.File{
		Path:     "/",
		FileMode: 0x800001fd,
		Mtime:    time.Unix(1698686254, 1698686254005370038),
		Data:     nil,
	}, "/Tooltip.css": &assets.File{
		Path:     "/Tooltip.css",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1698159550, 1698159550113092816),
		Data:     []byte(_Assets8bd9efe3e913f3cfc8c2405a1ef3fab1c7e66a1f),
	}, "/Progress.js": &assets.File{
		Path:     "/Progress.js",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1698690370, 1698690370636077486),
		Data:     []byte(_Assetsa0084176b3c3ed34696e0c93f712ea4c10c94315),
	}}, "")
