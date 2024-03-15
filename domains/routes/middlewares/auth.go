package middlewares

import (
	"context"
	"emailn/domains/campaign"
	"net/http"
	"os"
	"strings"

	oidc "github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-chi/render"
	jwtGo "github.com/golang-jwt/jwt/v5"
)

type Handler struct {
	CampaignService campaign.Service
}

type TokenVerifier interface {
	VerifyToken(ctx context.Context, token string) (*TokenClaims, error)
}

type TokenClaims struct {
	Email string
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

		keyCloakURL := os.Getenv("KEY_CLOAK_URL")
		keyCloakClientID := os.Getenv("KEY_CLOAK_CLIENT_ID")
		provider, err := oidc.NewProvider(r.Context(), keyCloakURL)

		if err != nil {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, map[string]string{
				"error": "Invalid token",
			})
			return
		}

		_, err = provider.Verifier(&oidc.Config{ClientID: keyCloakClientID}).Verify(r.Context(), token)

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
