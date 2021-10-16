package httpapi

import (
	"context"
	"net/http"

	"github.com/hmoragrega/fastlane"
	"github.com/julienschmidt/httprouter"
)

type lister interface {
	ListOpenReviews(ctx context.Context) ([]fastlane.Review, error)
}

func ListOpenReviews(svc lister) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		mrs, err := svc.ListOpenReviews(r.Context())
		encode(w, mrs, err)
	}
}
