package main

// go:generate go-bindata-assetfs -pkg bindata -o bindata/bindata.go resources/...
// go:generate env GOOS=linux GOARCH=arm go build -o RPi-Cast-ARM
// go:generate ssh pi@192.168.179.65 "pkill RPi-Cast-ARM; exit 0"
// go:generate scp RPi-Cast-ARM  pi@192.168.179.65:/home/pi/RPi-Cast/
// go:generate ssh pi@192.168.179.65 "cd /home/pi/RPi-Cast/; ./RPi-Cast-ARM"

//go:generate go-bindata-assetfs -pkg bindata -o bindata/bindata.go resources/...
//go:generate go build -o RPi-Cast

import (
	"net/http"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/go-chi/chi"
	"github.com/i-PUSH/RPi-Cast/bindata"
	"github.com/i-PUSH/RPi-Cast/src/controllers/chrome"
	"github.com/i-PUSH/RPi-Cast/src/controllers/omxplayer"
	"github.com/i-PUSH/RPi-Cast/src/utils"
)

func main() {
	router := chi.NewRouter()

	// serve static files
	router.Mount("/public/",
		http.StripPrefix("/public",
			http.FileServer(&assetfs.AssetFS{Asset: bindata.Asset, AssetDir: bindata.AssetDir, AssetInfo: bindata.AssetInfo, Prefix: "resources/public"}),
		),
	)

	// setup web controllers
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		utils.ExecTemplate(w, "resources/templates/index.html", nil)
	})

	chrome.Init(router)
	omxplayer.Init(router)

	// check and generate https certs
	utils.CheckHTTPS("127.0.0.1:8443")

	// start the server (HTTPS) on port 8443
	println("RPi-Cast listening on port 8443!")
	http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", router)
}
