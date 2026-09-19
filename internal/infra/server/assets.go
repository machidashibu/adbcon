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

var assetsTable AssetsMap = AssetsMap{
	"/gui":        AssetsInfo{Ext: ".html", Data: htmlMain},
	"/styles.css": AssetsInfo{Ext: ".css", Data: cssStyles},
}
