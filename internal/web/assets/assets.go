package assets

import "embed"

//go:embed favicon.ico style.css script.js
var All embed.FS
