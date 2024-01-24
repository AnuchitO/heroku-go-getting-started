package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"gitdev.devops.krungthai.com/aster/ariskill/app"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type mockStorage struct {
	UserStorage
	users []User

	methodsToCall map[string]bool
	err           error
	userReturn    *User
}

func (ms *mockStorage) ExpectToCall(methodName string) {
	if ms.methodsToCall == nil {
		ms.methodsToCall = make(map[string]bool)
	}
	ms.methodsToCall[methodName] = false
}

func (m *mockStorage) GetAll(ctx context.Context) ([]User, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.users, nil
}

func (m *mockStorage) GetOneById(ctx context.Context, id string) (*User, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.userReturn, nil
}

func (m *mockStorage) GetHardSkillById(ctx context.Context, name string) ([]User, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.users, nil
}

func TestGetAllUsers(t *testing.T) {
	t.Run("should return 200 and users", func(t *testing.T) {
		date, _ := time.Parse("2006-01-02T15:04:05.999999-07:00", "2006-01-02T15:04:05.999999-07:00")
		users := []User{
			{
				ID:        "999999999999999999991",
				Email:     "ariskill@arise.tech",
				FirstName: "Ariser1",
				LastName:  "skills",
				CreatedAt: date,
				UpdatedAt: date,
				CreatedBy: "system",
				UpdatedBy: "someone",
				AboutMe:   "",
			},
		}
		mock := &mockStorage{
			users: users,
		}
		handler := NewUserHandler(mock)

		engine := gin.New()
		engine.GET("/users", app.NewGinHandler(handler.GetAllUsers, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/users", nil)

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()

		want := `{
			"status": "success",
			"message": "",
			"data": [
				{
					"aboutMe":         "",
					"createdAt":       "2006-01-02T15:04:05.999999-07:00",
					"createdBy":       "system",
					"email":            "ariskill@arise.tech",
					"employeeId":      "",
					"familyName":      "skills",
					"givenName":       "Ariser1",
					"hardSkills":		null,
					"jobRole":          "",
					"level":			"",
					"socialMedias":    null,
					"softSkills":      null,
					"squadId":         null,
					"sub":              "999999999999999999991",
					"tags":             null,
					"technicalSkills": null,
					"hardSkills": 	   null,
					"updatedAt":       "2006-01-02T15:04:05.999999-07:00",
					"updatedBy":       "someone"
				}
			]
		}`

		assert.Equal(t, 200, rec.Code)
		assert.JSONEq(t, want, resp)
	})

	t.Run("UnexpectedError", func(t *testing.T) {
		mock := &mockStorage{
			err: errors.New("UnexpectedError"),
		}
		handler := NewUserHandler(mock)

		engine := gin.New()
		engine.GET("/users", app.NewGinHandler(handler.GetAllUsers, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/users", nil)

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()

		want := `{
			"status": "error",
			"message": "UnexpectedError"
		}`

		assert.Equal(t, 500, rec.Code)
		assert.JSONEq(t, want, resp)
	})

	t.Run("User not found", func(t *testing.T) {
		mock := &mockStorage{
			err: errors.New("User not found"),
		}
		handler := NewUserHandler(mock)

		engine := gin.New()
		engine.GET("/users", app.NewGinHandler(handler.GetAllUsers, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/users", nil)

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()

		want := `{
			"status": "error",
			"message": "User not found"
		}`

		assert.Equal(t, 500, rec.Code)
		assert.JSONEq(t, want, resp)
	})
}

func randomUsers(n int) (users []User) {
	for i := 0; i < n; i++ {
		user := User{
			ID:         primitive.NewObjectID().Hex(),
			Email:      "ariser@arise.dev",
			EmployeeID: strconv.Itoa(i),
			FirstName:  "ariser",
			LastName:   "by krungthai",
			MySquad: []MySquad{
				{
					SquadID: primitive.NewObjectID(),
					Role:    "member",
				},
			},
		}
		users = append(users, user)
	}

	return
}

func TestGetUsersData(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("should return 200 and users Name and Email", func(t *testing.T) {
		users := randomUsers(20)
		mockStorage := &mockStorage{
			users: users,
			err:   nil,
		}
		mockStorage.ExpectToCall("GetAll")
		handler := NewUserHandler(mockStorage)

		engine := gin.New()
		engine.GET("/users/email", app.NewGinHandler(handler.GetUsersData, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/users/email", nil)

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()

		var userResponses []GetEmailNameResponse
		for _, item := range users {
			userResponses = append(userResponses, GetEmailNameResponse{
				Email: item.Email,
				Name:  item.FirstName + " " + item.LastName,
			})
		}

		jsonString, _ := json.Marshal(userResponses)
		want := fmt.Sprintf(`{
			"status": "success",
			"message": "",
			"data": %s
		}`, jsonString)

		assert.Equal(t, 200, rec.Code)
		assert.JSONEq(t, want, resp)
	})

	t.Run("Should return 500, Internal Server error", func(t *testing.T) {
		mockStorage := &mockStorage{
			err: errors.New("fjslafjs"),
		}
		mockStorage.ExpectToCall("GetAll")
		handler := NewUserHandler(mockStorage)

		engine := gin.New()
		engine.GET("/users/email", app.NewGinHandler(handler.GetUsersData, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/users/email", nil)

		engine.ServeHTTP(rec, req)
		resp := rec.Body.String()
		want := `{
			"status": "error",
			"message": "fjslafjs"
		}`
		assert.Equal(t, 500, rec.Code)
		assert.JSONEq(t, want, resp)
	})
}

func TestGetHardSkillById(t *testing.T) {
	t.Run("should return 200 and hardSkills", func(t *testing.T) {
		date, _ := time.Parse("2006-01-02T15:04:05.999999-07:00", "2006-01-02T15:04:05.999999-07:00")
		users := []User{
			{
				ID:        "999999999999999999991",
				Email:     "ariskill@arise.tech",
				FirstName: "Ariser1",
				LastName:  "skills",
				CreatedAt: date,
				UpdatedAt: date,
				CreatedBy: "system",
				UpdatedBy: "someone",
				AboutMe:   "",
			},
		}
		mock := &mockStorage{
			users: users,
		}
		handler := NewUserHandler(mock)

		engine := gin.New()
		engine.GET("/users/hardskills", app.NewGinHandler(handler.GetAllUsers, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/users/hardskills", nil)

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()

		want := `{
			"status": "success",
			"message": "",
			"data": [
				{
					"aboutMe": "",
					"createdAt": "2006-01-02T15:04:05.999999-07:00",
					"createdBy": "system",
					"email": "ariskill@arise.tech",
					"employeeId": "",
					"familyName": "skills",
					"givenName": "Ariser1",
					"hardSkills": null,
					"jobRole": "",
					"level": "",
					"socialMedias": null,
					"softSkills": null,
					"squadId": null,
					"sub": "999999999999999999991",
					"tags": null,
					"technicalSkills": null,
					"updatedAt": "2006-01-02T15:04:05.999999-07:00",
					"updatedBy": "someone"
				}
			]
		}`

		assert.Equal(t, 200, rec.Code)
		assert.JSONEq(t, want, resp)
	})
}
