package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"yadro-go-course/internal/adapters/web/handlers"
	"yadro-go-course/internal/core/ports"
	"yadro-go-course/internal/core/services"
)

type credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func Login(userSrv *services.UserService, tokenMaxTime time.Duration) handlers.ErrHandleFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var cr credentials

		err := json.NewDecoder(r.Body).Decode(&cr)
		if err != nil {
			return errors.Join(handlers.ErrBadRequest, err)
		}

		if cr.Login == "" || len(cr.Password) == 0 {
			return errors.Join(handlers.ErrBadRequest, handlers.ErrNoLoginOrPassword)
		}

		u, err := userSrv.UserLogin(r.Context(), cr.Login)
		if err != nil {
			if errors.Is(err, ports.ErrNotFound) {
				return errors.Join(handlers.ErrNotFound, err)
			}

			return errors.Join(handlers.ErrInternal, err)
		}

		err = bcrypt.CompareHashAndPassword(u.Password, []byte(cr.Password))
		if err != nil {
			return errors.Join(handlers.ErrForbidden, err)
		}

		claims := customClaims{UserID: u.ID, IsAdmin: u.IsAdmin, RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenMaxTime))}}
		token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

		tokenStr, err := token.SignedString(jwtSecret)
		if err != nil {
			return errors.Join(handlers.ErrInternal, err)
		}

		w.Header().Set("Authorization", tokenStr)

		return nil
	}
}
