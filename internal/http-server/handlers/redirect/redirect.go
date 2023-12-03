package redirect

import (
	"errors"
	resp "github.com/MirToykin/url-shortner/internal/lib/api/response"
	"github.com/MirToykin/url-shortner/internal/lib/logger/sl"
	"github.com/MirToykin/url-shortner/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
	"net/http"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLGetter
type URLGetter interface {
	GetUrl(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.redirect.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		log.Info("params from url: ", slog.String("alias", alias))

		url, err := urlGetter.GetUrl(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info(storage.ErrURLNotFound.Error(), slog.String("alias", alias))

			render.JSON(w, r, resp.Error(storage.ErrURLNotFound.Error()))

			return
		}

		if err != nil {
			log.Info("failed to get url", slog.String("alias", alias), sl.Err(err))

			render.JSON(w, r, resp.Error("Internal error"))

			return
		}

		log.Info("url retrieved", slog.String("alias", alias), slog.String("url", url))
		http.Redirect(w, r, url, http.StatusFound)
	}
}
