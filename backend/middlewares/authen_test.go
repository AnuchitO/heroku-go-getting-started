package middlewares

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitdev.devops.krungthai.com/aster/ariskill/app"
	"gitdev.devops.krungthai.com/aster/ariskill/app/profile"
	"gitdev.devops.krungthai.com/aster/ariskill/config"
	"github.com/gin-gonic/gin"
)

func TestValidateGoogleIdToken(t *testing.T) {
	t.Run("set profileID, email and role into context", func(t *testing.T) {
		var fakeUserStorageFunc = func(ctx context.Context, email string) (*profile.User, error) {
			return &profile.User{
				JobRole: "fullstack",
			}, nil
		}
		var googleOidc = config.GoogleOidc{
			IsDevMode: true,
		}
		middleware := ValidateGoogleIdToken(fakeUserStorageFunc, googleOidc, app.RealClock{})

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(http.MethodPost, "/users", nil)
		ctx.Request.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsImtpZCI6IjZmNzI1NDEwMWY1NmU0MWNmMzVjOTkyNmRlODRhMmQ1NTJiNGM2ZjEiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJhenAiOiIzMDY5NTMwMjQ1NTktamR1OW9kZGxmbTQ3YmdvMTU2dGNsMjE3YmE5dGRqOGwuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJhdWQiOiIzMDY5NTMwMjQ1NTktamR1OW9kZGxmbTQ3YmdvMTU2dGNsMjE3YmE5dGRqOGwuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJzdWIiOiI5OTk5OTk5OTk5OTk5OTk5OTk5OTEiLCJoZCI6ImFyaXNlLnRlY2giLCJlbWFpbCI6ImFtYWxyaWMubEBhcmlzZS50ZWNoIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImF0X2hhc2giOiJUTW9SOFVFRjQ4Q1lBNmQwYTVUUm5RIiwibmFtZSI6IkFtYWxyaWMgTG9ja3dvb2QiLCJwaWN0dXJlIjoiaHR0cHM6Ly9saDMuZ29vZ2xldXNlcmNvbnRlbnQuY29tL2EvQUNnOG9jSWxNdjFNS1BYTV8za01Kc3lNbTR6U2RnaEliM1FqVGVOdHFKd0prOFlSRWc9czk2LWMiLCJnaXZlbl9uYW1lIjoiQW1hbHJpYyIsImZhbWlseV9uYW1lIjoiTG9ja3dvb2QiLCJsb2NhbGUiOiJ0aCIsImlhdCI6MTY5NTg3NjExNCwiZXhwIjozNjk1ODc5NzE0fQ.1vlq0gu4ZZvQbOLQChshW0tOlWdeHO69-eF2C4pg6-U")
		middleware(ctx)

		expectedContextProfileID := "999999999999999999991"
		expectedContextEmail := "amalric.l@arise.tech"
		expectedRole := "fullstack"

		if profileID, ok := ctx.Get(ContextProfileID); ok {
			if expectedContextProfileID != profileID {
				t.Errorf("in case of everythings is ok, expect profileID is %q but get %v\n", expectedContextProfileID, profileID)
			}
		}
		if email, ok := ctx.Get(ContextEmail); ok {
			if expectedContextEmail != email {
				t.Errorf("in case of everythings is ok, expect email is %q but get %v\n", expectedContextEmail, email)
			}
		}
		if role, ok := ctx.Get(ContextRole); ok {
			if expectedRole != role {
				t.Errorf("in case of everythings is ok, expect role is %q but get %v\n", expectedRole, role)
			}
		}
	})
}
