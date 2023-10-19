package views

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _Assets636a8d4dc640bc67422b3e0c56cbbf6a4ee4f60d = ""

// Assets returns go-assets FileSystem
var Assets = assets.NewFileSystem(map[string][]string{"/": []string{"placeholder.js"}}, map[string]*assets.File{
	"/": &assets.File{
		Path:     "/",
		FileMode: 0x800001fd,
		Mtime:    time.Unix(1697731258, 1697731258178357163),
		Data:     nil,
	}, "/placeholder.js": &assets.File{
		Path:     "/placeholder.js",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1697643194, 1697643194308569895),
		Data:     []byte(_Assets636a8d4dc640bc67422b3e0c56cbbf6a4ee4f60d),
	}}, "")
