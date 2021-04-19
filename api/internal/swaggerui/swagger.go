package swaggerui

import (
	"net/http"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	swagger "github.com/sei-ri/pingpong.io/api/web/ui/swagger/data"
)

func Serve(mux *http.ServeMux) {
	fileServer := http.FileServer(&assetfs.AssetFS{
		Asset:    swagger.Asset,
		AssetDir: swagger.AssetDir,
		Prefix:   "third_party/swagger-ui",
	})
	prefix := "/swagger-ui/"
	mux.Handle(prefix, http.StripPrefix(prefix, fileServer))
}
