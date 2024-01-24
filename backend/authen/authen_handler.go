package authen

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gitdev.devops.krungthai.com/aster/ariskill/app"
	"gitdev.devops.krungthai.com/aster/ariskill/config"
	"gitdev.devops.krungthai.com/aster/ariskill/errs"
)

var ErrUnexpected = errors.New("Unexpected error")

type AuthenHandler interface {
	ExchangeForTokens(c app.Context)
}

type authenHandler struct {
	client     *http.Client
	googleOidc config.GoogleOidc
}

type AuthenRequest struct {
	ExchangeType string `json:"type"`
	Value        string `json:"value"`
}
type AuthResponseData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewAuthenHandler(client *http.Client, googleOidc config.GoogleOidc) *authenHandler {
	return &authenHandler{
		client:     client,
		googleOidc: googleOidc,
	}
}

func AuthResponseError(err error) (code int, obj any) {
	switch e := err.(type) {
	case errs.AppError:
		code = e.Code
		obj = AuthResponseData{
			Code:    code,
			Message: e.Message,
			Data:    nil,
		}
	case error:
		code = http.StatusInternalServerError
		obj = AuthResponseData{
			Code:    code,
			Message: ErrUnexpected.Error(),
			Data:    nil,
		}
	}
	return code, obj
}

// /auth/token
func (s *authenHandler) ExchangeForTokens(c app.Context) {
	reqJson := &AuthenRequest{}
	if err := c.ShouldBindJSON(reqJson); err != nil {
		c.InternalServerError(err) // TODO: Improve error code
		return
	}
	// Prepare request body for Google token API
	beReqData := url.Values{}
	beReqData.Set("client_id", s.googleOidc.ClientId)         // TODO: For now, use your own client id & Refactor to env
	beReqData.Set("client_secret", s.googleOidc.ClientSecret) // TODO: For now, use your own client secret & Refactor to env

	if reqJson.ExchangeType != "auth_code" && reqJson.ExchangeType != "refresh_token" {
		c.JSON(http.StatusBadRequest, "Wrong exchange type!") // TODO: Refactor to central error
		return
	}

	if reqJson.ExchangeType == "auth_code" {
		beReqData.Set("code", reqJson.Value)
		beReqData.Set("redirect_uri", s.googleOidc.RedirectUri)
		beReqData.Set("grant_type", "authorization_code")
	}
	if reqJson.ExchangeType == "refresh_token" {
		beReqData.Set("refresh_token", reqJson.Value)
		beReqData.Set("grant_type", "refresh_token")
	}
	ctx, cancel := context.WithTimeout(c.Ctx(), 10*time.Second)
	defer cancel()
	beReq, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://oauth2.googleapis.com/token", strings.NewReader(beReqData.Encode()))
	beReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		c.InternalServerError(err) // TODO: Improve error code
		return
	}
	beResp, err := s.client.Do(beReq)
	if err != nil {
		c.InternalServerError(err) // TODO: Improve error code
		return
	}
	defer beResp.Body.Close()
	tokenResultJsonString, err := io.ReadAll(beResp.Body)
	if err != nil {
		c.InternalServerError(err) // TODO: Improve error code
		return
	}
	var tokenResultMap map[string]interface{}
	if err := json.Unmarshal(tokenResultJsonString, &tokenResultMap); err != nil {
		c.InternalServerError(err) // TODO: Improve error code
		return
	}
	delete(tokenResultMap, "access_token")    // Remove "access_token" for security
	c.JSON(beResp.StatusCode, tokenResultMap) // Full-proxied, including HTTP status
}
