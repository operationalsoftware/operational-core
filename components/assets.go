package components

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _Assets835ab29723abda1b6c53883e64bdb0a0df16131c = ""

// Assets returns go-assets FileSystem
var Assets = assets.NewFileSystem(map[string][]string{"/": []string{"placeholder.asset"}}, map[string]*assets.File{
	"/": &assets.File{
		Path:     "/",
		FileMode: 0x800001fd,
		Mtime:    time.Unix(1697187495, 1697187495687636959),
		Data:     nil,
	}, "/placeholder.asset": &assets.File{
		Path:     "/placeholder.asset",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1697180096, 1697180096563144516),
		Data:     []byte(_Assets835ab29723abda1b6c53883e64bdb0a0df16131c),
	}}, "")
