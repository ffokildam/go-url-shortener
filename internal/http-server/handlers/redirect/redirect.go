package redirect

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	resp "url-short/internal/lib/api/response"
	"url-short/internal/lib/logger/sl"
	"url-short/internal/storage"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.redirect.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		log.Info("retrieved alias", slog.String("alias", alias))
		if alias == "" {
			log.Error("empty alias")

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		resURL, err := urlGetter.GetURL(alias)
		fmt.Println(resURL)

		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", alias)

			render.JSON(w, r, resp.Error("url not found"))
		}

		if err != nil {
			log.Error("failed to get url", sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}
		log.Info("retrieved url", slog.String("url", resURL))

		http.Redirect(w, r, resURL, http.StatusFound)
	}
}
