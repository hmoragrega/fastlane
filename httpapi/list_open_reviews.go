package httpapi

import (
	"net/http"

	"github.com/hmoragrega/fastlane"
	"github.com/julienschmidt/httprouter"
)

type lister interface {
	OpenReviews() []fastlane.Review
}

func ListOpenReviews(svc lister) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		encodeData(w, svc.OpenReviews())
	}
}
