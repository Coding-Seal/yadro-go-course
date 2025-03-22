package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"yadro-go-course/pkg/http-util"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"yadro-go-course/internal/core/ports"
	"yadro-go-course/internal/core/services"
)

type credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func Login(userSrv *services.UserService, tokenMaxTime time.Duration) http_util.ErrHandleFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var cr credentials

		err := json.NewDecoder(r.Body).Decode(&cr)
		if err != nil {
			return errors.Join(http_util.ErrBadRequest, err)
		}

		if cr.Login == "" || len(cr.Password) == 0 {
			return errors.Join(http_util.ErrBadRequest, http_util.ErrNoLoginOrPassword)
		}

		u, err := userSrv.UserLogin(r.Context(), cr.Login)
		if err != nil {
			if errors.Is(err, ports.ErrNotFound) {
				return errors.Join(http_util.ErrNotFound, err)
			}

			return errors.Join(http_util.ErrInternal, err)
		}

		err = bcrypt.CompareHashAndPassword(u.Password, []byte(cr.Password))
		if err != nil {
			return errors.Join(http_util.ErrForbidden, err)
		}

		claims := CustomClaims{UserID: u.ID, IsAdmin: u.IsAdmin, RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenMaxTime))}}
		token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

		tokenStr, err := token.SignedString(jwtSecret)
		if err != nil {
			return errors.Join(http_util.ErrInternal, err)
		}

		w.Header().Set("Authorization", tokenStr)

		return nil
	}
}
