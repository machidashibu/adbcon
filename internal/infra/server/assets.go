package server

import (
	_ "embed"
)

//go:embed assets/cert.pem
var serverCert []byte

//go:embed assets/key.pem
var serverKey []byte

//go:embed assets/html/main.html
var htmlMain []byte

//go:embed assets/css/styles.css
var cssStyles []byte

//go:embed assets/js/main.js
var jsMain []byte

//go:embed assets/js/conponents.js
var jsConponents []byte

var assetsTable AssetsMap = AssetsMap{
	"/gui":           AssetInfo{Ext: ".html", Data: htmlMain},
	"/styles.css":    AssetInfo{Ext: ".css", Data: cssStyles},
	"/main.js":       AssetInfo{Ext: ".js", Data: jsMain},
	"/conponents.js": AssetInfo{Ext: ".js", Data: jsConponents},
}
