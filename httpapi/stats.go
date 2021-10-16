package httpapi

import (
	"net/http"

	"github.com/hmoragrega/fastlane"
	"github.com/julienschmidt/httprouter"
)

type stater interface {
	Stats() fastlane.Stats
}

func GetStats(svc stater) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		encodeData(w, svc.Stats())
	}
}
