package create

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	resp "github.com/iamvkosarev/go-shared-utils/api/response"
	"github.com/iamvkosarev/go-shared-utils/api/user"
	"github.com/iamvkosarev/go-shared-utils/logger/sl"
	"github.com/iamvkosarev/wishlist/back/internal/model"
	"github.com/iamvkosarev/wishlist/back/internal/storage"
	"io"
	"log/slog"
	"net/http"
)

type Request struct {
	OwnerID     int64  `json:"owner_id"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description,omitempty"`
	DisplayType int    `json:"display_type"`
}

type Response struct {
	resp.Response
	WishlistID int64 `json:"wishlist_id"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=WishlistSaver
type WishlistSaver interface {
	SaveWishlist(
		ownerID int64,
		name string,
		description string,
		displayType model.DisplayType,
	) (int64, error)
}

func NewWishlistHandler(log *slog.Logger, wishlistSaver WishlistSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.wishlist.create.NewWishlistHandler"

		log = log.With(
			slog.String("op", op),
			slog.String(
				"request_id",
				middleware.GetReqID(r.Context()),
			),
		)

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
			render.JSON(
				w, r, resp.Error("empty request"),
			)
			return
		}
		if err != nil {
			log.Error("failed to parse request", sl.Err(err))
			render.JSON(
				w, r, resp.Error("failed to parse request"),
			)
			return
		}

		log.Info("request body decoded", slog.Any("req", req))

		if err := validator.New().Struct(req); err != nil {

			var validationErr validator.ValidationErrors
			errors.As(err, &validationErr)

			log.Error("failed to validate request", sl.Err(err))
			render.JSON(
				w, r, resp.ValidateErrors(validationErr),
			)
			return
		}
		displayType, err := model.IntToDisplayType(req.DisplayType)
		if err != nil {
			log.Error("failed to parse display_type", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to parse display_type"))
			return
		}

		userID, err := user.GetUserId(r)

		if err != nil {
			log.Error("failed to get user id", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to get user id"))
			return
		}

		wishlistID, err := wishlistSaver.SaveWishlist(userID, req.Name, req.Description, displayType)
		if errors.Is(err, storage.ErrorWishlistExists) {
			log.Error("wishlist already exists", sl.Err(err))
			render.JSON(w, r, resp.Error("wishlist already exists"))
			return
		}
		if err != nil {
			log.Error("failed to create wishlist", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to create wishlist"))
			return
		}

		log.Info("wishlist added", slog.Int64("wishlist_id", wishlistID))
		responseOK(w, r, wishlistID)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, wishlistID int64) {
	render.JSON(
		w, r, Response{
			Response:   resp.Ok(),
			WishlistID: wishlistID,
		},
	)
}
