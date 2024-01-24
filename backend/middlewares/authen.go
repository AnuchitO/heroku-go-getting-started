package middlewares

import (
	"context"
	"encoding/json"
	"net/http"

	"strings"
	"time"

	"gitdev.devops.krungthai.com/aster/ariskill/app"
	"gitdev.devops.krungthai.com/aster/ariskill/app/user"
	authen "gitdev.devops.krungthai.com/aster/ariskill/authen"
	"gitdev.devops.krungthai.com/aster/ariskill/config"
	"gitdev.devops.krungthai.com/aster/ariskill/errs"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"google.golang.org/api/idtoken"
)

const (
	ERROR_AUTH_UNAUTHORIZED        = "Unauthorized!"
	ERROR_AUTH_NOT_SUPPORTED       = "Only Bearer authorization is supported."
	ERROR_AUTH_EMPTY               = "Bearer authorization can't be left blank!"
	ERROR_GET_DATA                 = "Failed to transfer data from json to bson"
	ERROR_USER_UPSERT_FAILED       = "Failed to upsert data"
	ERROR_USER_TIMESTAMPING_FAILED = "Failed to timestamp"
	ERROR_DEV_AUTH_EXPIRED         = "Even in dev mode, token can't be expired!"
	ERROR_DEV_AUTH_AUD_MISMATCH    = "idtoken: audience provided does not match aud claim in the JWT"
	ERROR_WRONG_DOMAIN             = "Only @arise.tech email is allowed"
)

type userStorageFunc func(ctx context.Context, email string) (*user.User, error)

const (
	contextProfileID = "profileID"
	contextEmail     = "email"
	contextRole      = "role"
)

func ValidateGoogleIdToken(userStorage userStorageFunc, googleOidc config.GoogleOidc, clock app.Clock) gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/swagger") {
			c.Next()
			return
		}

		authorization := c.GetHeader("Authorization")

		token, err := authorizationToken(authorization)
		if err != nil {
			c.AbortWithStatusJSON(authen.AuthResponseError(err))
			return
		}

		var finalParsedIdToken *ParsedIdToken
		// Bypass the real validation process if the corresponding ENV is set
		if googleOidc.IsDevMode {
			var err error
			finalParsedIdToken, err = devModeToken(token, googleOidc, clock)
			if err != nil {
				c.AbortWithStatusJSON(authen.AuthResponseError(err))
				return
			}
		} else { // Real Validation Process
			ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
			defer cancel()
			googleIdToken, err := idtoken.Validate(ctx, token, googleOidc.ClientId)
			if err != nil {
				c.AbortWithStatusJSON(authen.AuthResponseError(errs.NewUnauthorizedError(err.Error())))
				return
			}
			parsedIdToken := NewParsedIdToken(googleIdToken)
			finalParsedIdToken = parsedIdToken
		}

		// Check that the domain part of email authenticated must be "arise.tech"
		if !strings.HasSuffix(finalParsedIdToken.Email, "@arise.tech") {
			c.AbortWithStatusJSON(authen.AuthResponseError(errs.NewUnauthorizedError(ERROR_WRONG_DOMAIN)))
			return
		}

		currentUser, err := tokenUser(finalParsedIdToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()
		user, err := userStorage(ctx, currentUser.Email)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.Set(contextProfileID, currentUser.ID)
		c.Set(contextEmail, currentUser.Email)
		c.Set(contextRole, user.JobRole)

		c.Next()
	}
}

func devModeToken(token string, googleOidc config.GoogleOidc, clock app.Clock) (*ParsedIdToken, error) {
	unsafeToken, _, err := new(jwt.Parser).ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}
	unsafeParsedToken := new(ParsedIdToken)
	bs, _ := json.Marshal(unsafeToken.Claims)
	if err := json.Unmarshal(bs, unsafeParsedToken); err != nil {
		return nil, err
	}
	if clock.Now().Unix() > int64(unsafeParsedToken.Expires) {
		return nil, errs.NewUnauthorizedError(ERROR_DEV_AUTH_EXPIRED)
	}
	if googleOidc.ClientId != "" && googleOidc.ClientId != unsafeParsedToken.Audience {
		return nil, errs.NewUnauthorizedError(ERROR_DEV_AUTH_AUD_MISMATCH)
	}
	return unsafeParsedToken, nil
}

func authorizationToken(authorization string) (token string, err error) {
	if authorization == "" {
		return "", errs.NewUnauthorizedError(ERROR_AUTH_UNAUTHORIZED)
	}

	authorizationParts := strings.Fields(authorization)
	if len(authorizationParts) != 2 {
		return "", errs.NewUnauthorizedError(ERROR_AUTH_EMPTY)
	}

	if authorizationParts[0] != "Bearer" {
		return "", errs.NewUnauthorizedError(ERROR_AUTH_NOT_SUPPORTED)
	}

	if authorizationParts[1] == "" {
		return "", errs.NewUnauthorizedError(ERROR_AUTH_EMPTY)
	}

	return authorizationParts[1], nil
}

func tokenUser(finalParsedIdToken *ParsedIdToken) (user.User, error) {
	var currentUser user.User
	bs, err := json.Marshal(finalParsedIdToken)
	if err != nil {
		return user.User{}, errors.WithMessage(err, ERROR_GET_DATA)
	}

	err = json.Unmarshal(bs, &currentUser)
	if err != nil {
		return user.User{}, errors.WithMessage(err, ERROR_GET_DATA)
	}

	return currentUser, nil
}
