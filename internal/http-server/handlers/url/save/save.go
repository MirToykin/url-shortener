package save

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	resp "gitlab.com/mt65/url-shortner/internal/lib/api/response"
	"gitlab.com/mt65/url-shortner/internal/lib/logger/sl"
	"gitlab.com/mt65/url-shortner/internal/lib/random"
	"gitlab.com/mt65/url-shortner/internal/storage"
	"golang.org/x/exp/slog"
	"net/http"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLSaver
type URLSaver interface {
	SaveUrl(urlToSave, alias string) (int64, error)
	CheckIfAliasExists(alias string) error
}

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

// TODO move to config
const aliasLength = 6
const maxAttemptsToGenerateAlias = 10

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		alias := req.Alias
		if alias == "" {
			for i := 0; i < maxAttemptsToGenerateAlias; i++ {
				alias = random.NewRandomString(aliasLength)
				err = urlSaver.CheckIfAliasExists(alias)
				if err == nil {
					break
				}

				log.Info("generate alias failed, retrying...", slog.Int("attempt", i+1), sl.Err(err))
			}

			if err != nil {
				msg := "failed to generate alias"
				log.Error(msg, sl.Err(err))

				render.JSON(w, r, resp.Error(msg))

				return
			}
		}

		id, err := urlSaver.SaveUrl(req.URL, req.Alias)
		if errors.Is(err, storage.ErrAliasExists) {
			log.Info("alias exists", slog.String("alias", req.Alias))

			render.JSON(w, r, resp.Error("url already exists"))

			return
		}

		if err != nil {
			log.Info("failed to add url", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to add url"))

			return
		}

		log.Info("url added", slog.Int64("id", id), slog.String("alias", req.Alias), slog.String("url", req.URL))
		responseOk(w, r, alias)
	}
}

func responseOk(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Alias:    alias,
	})
}
