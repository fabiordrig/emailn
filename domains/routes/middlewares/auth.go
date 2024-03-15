package middlewares

import (
	"context"
	"emailn/domains/campaign"
	"net/http"
	"strings"

	oidc "github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-chi/render"
	jwtGo "github.com/golang-jwt/jwt/v5"
)

type Handler struct {
	CampaignService campaign.Service
}

type RouterFunc func(w http.ResponseWriter, r *http.Request) (interface{}, int, error)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("Authorization")

		if token == "" {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, map[string]string{
				"error": "Unauthorized",
			})
			return
		}

		token = strings.Replace(token, "Bearer ", "", 1)

		provider, err := oidc.NewProvider(r.Context(), "http://localhost:8080/realms/provider")

		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{
				"error": "Internal Server Error",
			})
			return
		}

		_, err = provider.Verifier(&oidc.Config{ClientID: "emailn"}).Verify(r.Context(), token)

		if err != nil {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, map[string]string{
				"error": "Unauthorized",
			})
			return
		}

		value, _ := jwtGo.Parse(token, nil)
		claims := value.Claims.(jwtGo.MapClaims)
		email := claims["email"]

		ctx := context.WithValue(r.Context(), "email", email)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
