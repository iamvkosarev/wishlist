package create

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/iamvkosarev/go-shared-utils/slog/slog_discard"
	"github.com/iamvkosarev/wishlist/back/internal/http-server/handlers/wishlist/create/mocks"
	"github.com/iamvkosarev/wishlist/back/internal/model"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateHandler(t *testing.T) {
	tests := []struct {
		testName     string
		ownerID      int64
		wishlistName string
		description  string
		displayType  int
		respError    string
		mockError    error
	}{
		{
			testName:     "wrong display type",
			ownerID:      int64(1),
			wishlistName: "Test",
			description:  "",
			displayType:  -1,
			mockError:    nil,
			respError:    model.ErrorInvalidDisplayType.Error(),
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.testName, func(t *testing.T) {
				wishlistSaverMock := mocks.NewWishlistSaver(t)
				displayType, err := model.IntToDisplayType(tt.displayType)
				if err != nil {
					require.Equal(t, tt.respError, err.Error())
					return
				}

				wishlistSaverMock.On(
					"SaveWishlist",
					tt.ownerID,
					tt.wishlistName,
					tt.description,
					displayType,
				).Return(
					int64(0),
					tt.mockError,
				).Once()

				handler := NewWishlistHandler(slog_discard.NewDiscardLogger(), wishlistSaverMock)

				input := fmt.Sprintf(
					`{"owner_id": %v, "name": "%s", "description" : "%s", "display_type": %v}`,
					tt.ownerID, tt.wishlistName, tt.description, displayType,
				)

				req, err := http.NewRequest(http.MethodPost, "/wishlist", bytes.NewBufferString(input))
				require.NoError(t, err)

				rr := httptest.NewRecorder()
				handler.ServeHTTP(rr, req)

				require.Equal(t, rr.Code, http.StatusOK)

				body := rr.Body.String()

				var resp Response
				require.NoError(t, json.Unmarshal([]byte(body), &resp))

				require.Equal(t, tt.respError, resp.Error)
			},
		)
	}
}
