package get

import (
	"database/sql"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	resp "github.com/iamvkosarev/go-shared-utils/api/response"
	"github.com/iamvkosarev/go-shared-utils/logger/sl"
	"github.com/iamvkosarev/wishlist/back/internal/model"
	"log/slog"
	"net/http"
	"strconv"
)

const WishlistIDParam = "wishlist_id"

type Response struct {
	resp.Response
	Wishlist model.Wishlist `json:"wishlist"`
}

type WishlistProvider interface {
	GetWishlist(wishlistID int64) (model.Wishlist, error)
}

func NewWishlistHandler(log *slog.Logger, wishlistProvider WishlistProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.wishlist.get.NewWishlistHandler"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		wishlistID := chi.URLParam(r, WishlistIDParam)

		id, err := strconv.ParseInt(wishlistID, 10, 64)
		if err != nil {
			log.Error("failed to parse id parameter", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to parse id parameter from url"))
			return
		}

		wishlist, err := wishlistProvider.GetWishlist(id)
		if errors.Is(err, sql.ErrNoRows) {
			log.Error("failed to find wishlist", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to find wishlist"))
			return
		}
		if err != nil {
			log.Error("failed to get wishlist", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to get wishlist"))
			return
		}

		//TODO: Add owner and display type to prevent showing wrong user data
		log.Info("wishlist got", slog.Any("wishlist", wishlist))
		responseOK(r, w, wishlist)
	}
}

func responseOK(r *http.Request, w http.ResponseWriter, wishlist model.Wishlist) {
	render.JSON(
		w, r, Response{
			Response: resp.Ok(),
			Wishlist: wishlist,
		},
	)
}
