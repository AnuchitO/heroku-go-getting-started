package squad

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"gitdev.devops.krungthai.com/aster/ariskill/app"
	"gitdev.devops.krungthai.com/aster/ariskill/app/user"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

func TestSquadHandlerGetAll(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("should return 200 and all squads", func(t *testing.T) {
		id := "5e201c51e09c2c084c88a790"  // #nosec
		id2 := "5e201c51e09c2c084c88a791" // #nosec
		idObj, _ := primitive.ObjectIDFromHex(id)
		idObj2, _ := primitive.ObjectIDFromHex(id2)
		now := time.Now()
		nowFormatted := now.Format(time.RFC3339Nano)
		mockStorage := &mockSquadStorage{
			squad: []*Squad{
				{
					Id:            idObj,
					Name:          "Aster",
					TeamleadMail:  "john.d@arise.tech",
					Description:   "Squad in Blockchain team krub",
					SkillsRatings: []SkillRatings{},
					CreatedAt:     now,
				}, {
					Id:            idObj2,
					Name:          "Next",
					TeamleadMail:  "john.d@arise.tech",
					Description:   "Squad in core team krub",
					SkillsRatings: []SkillRatings{},
					CreatedAt:     now,
				},
			},
		}
		mockStorage.ExpectToCall("GetAll")
		handler := NewSquadHandler(mockStorage)

		engine := gin.New()
		engine.GET("/squads", app.NewGinHandler(handler.GetAll, zap.NewNop()))
		rec := httptest.NewRecorder()
		body := strings.NewReader("")
		req, _ := http.NewRequest(http.MethodGet, "/squads", body)
		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()

		want := fmt.Sprintf(`{
			"status": "success",
			"message": "",
			"data": [
				{
					"id": "%s",
					"name": "Aster",
					"teamleadMail": "john.d@arise.tech",
					"desc": "Squad in Blockchain team krub",
					"createdAt": "%s",
					"skillsRatings": []
				}, {
					"id": "%s",
					"name": "Next",
					"teamleadMail": "john.d@arise.tech",
					"desc": "Squad in core team krub",
					"createdAt": "%s",
					"skillsRatings": []
				}
			]
		}`, id, nowFormatted, id2, nowFormatted)

		assert.Equal(t, 200, rec.Code)
		assert.JSONEq(t, want, resp)
		mockStorage.Verify(t)
	})

	t.Run("should return 404 when database is empty", func(t *testing.T) {
		mockStorage := &mockSquadStorage{
			squad: nil,
			err:   squadNotFoundError,
		}
		mockStorage.ExpectToCall("GetAll")
		handler := NewSquadHandler(mockStorage)

		engine := gin.New()
		engine.GET("/squads", app.NewGinHandler(handler.GetAll, zap.NewNop()))
		rec := httptest.NewRecorder()
		body := strings.NewReader("")
		req, _ := http.NewRequest(http.MethodGet, "/squads", body)
		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()

		want := fmt.Sprintf(`{
			"status":"error",
			"message":"%s"
		}`, squadNotFoundError.Error())

		assert.Equal(t, 404, rec.Code)
		assert.JSONEq(t, want, resp)
		mockStorage.Verify(t)
	})

	t.Run("should return filtered squads when have filter", func(t *testing.T) {
		id := "5e201c51e09c2c084c88a790"  // #nosec
		id2 := "5e201c51e09c2c084c88a791" // #nosec
		idObj, _ := primitive.ObjectIDFromHex(id)
		idObj2, _ := primitive.ObjectIDFromHex(id2)
		now := time.Now()
		nowFormatted := now.Format(time.RFC3339Nano)
		mockSquad := []*Squad{
			{
				Id:            idObj,
				Name:          "Aster",
				TeamleadMail:  "john.d@arise.tech",
				Description:   "Squad in Blockchain team krub",
				SkillsRatings: []SkillRatings{},
				CreatedAt:     now,
			}, {
				Id:            idObj2,
				Name:          "Next",
				TeamleadMail:  "john.d@arise.tech",
				Description:   "Squad in core team krub",
				SkillsRatings: []SkillRatings{},
				CreatedAt:     now,
			},
		}

		mockStorage := &mockSquadStorage{squad: mockSquad}
		mockStorage.ExpectToCall("GetByFilter")
		handler := NewSquadHandler(mockStorage)

		engine := gin.New()
		engine.GET("/squads", app.NewGinHandler(handler.GetAll, zap.NewNop()))
		rec := httptest.NewRecorder()
		body := strings.NewReader(`{"name": "Aster"}`)
		req, _ := http.NewRequest(http.MethodGet, "/squads", body)
		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()

		want := fmt.Sprintf(`{
			"status": "success",
			"message": "",
			"data": [
				{
					"id": "%s",
					"name": "Aster",
					"teamleadMail": "john.d@arise.tech",
					"desc": "Squad in Blockchain team krub",
					"createdAt": "%s",
					"skillsRatings": []
				}
			]
		}`, id, nowFormatted)

		assert.Equal(t, 200, rec.Code)
		assert.JSONEq(t, want, resp)
		mockStorage.Verify(t)
	})

	t.Run("should return 404 when filtered squads is empty", func(t *testing.T) {
		mockStorage := &mockSquadStorage{squad: nil, err: squadNotFoundError}
		mockStorage.ExpectToCall("GetByFilter")
		handler := NewSquadHandler(mockStorage)

		engine := gin.New()
		engine.GET("/squads", app.NewGinHandler(handler.GetAll, zap.NewNop()))
		rec := httptest.NewRecorder()
		body := strings.NewReader(`{"name": "Connext"}`)
		req, _ := http.NewRequest(http.MethodGet, "/squads", body)
		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()

		want := fmt.Sprintf(`{
			"status":"error",
			"message":"%s"
		}`, squadNotFoundError.Error())

		assert.Equal(t, 404, rec.Code)
		assert.JSONEq(t, want, resp)
		mockStorage.Verify(t)
	})
}

func TestSquadHandlerGetOneByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("should return 200 and squad when found a squad match with squadID", func(t *testing.T) {
		squadId := mockObjectId(10)
		mockStorage := &mockSquadStorage{
			squad: []*Squad{
				{
					Id:            squadId.objectId,
					Name:          "Aster",
					TeamleadMail:  "john.d@arise.tech",
					Description:   "Squad in Blockchain team krub",
					SkillsRatings: []SkillRatings{},
				},
			},
			err: nil,
		}
		mockStorage.ExpectToCall("GetOneByID")
		handler := NewSquadHandler(mockStorage)

		engine := gin.New()
		engine.GET("/squads/:squadID", app.NewGinHandler(handler.GetOneByID, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/squads/"+squadId.hexId, nil)

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := `{
			"status": "success",
			"message": "",
			"data": {
				"id": "000000000000000000000010",
				"name": "Aster",
				"teamleadMail": "john.d@arise.tech",
				"desc": "Squad in Blockchain team krub",
				"createdAt": "0001-01-01T00:00:00Z",
				"skillsRatings": []
			}
		}`
		assert.Equal(t, 200, rec.Code)
		assert.JSONEq(t, want, resp)
		mockStorage.Verify(t)
	})

	t.Run("should return 400 and data is nil when invalid Id", func(t *testing.T) {
		squadId := "55"
		mockStorage := &mockSquadStorage{
			err: invalidIdError,
		}
		mockStorage.ExpectToCall("GetOneByID")
		handler := NewSquadHandler(mockStorage)

		engine := gin.New()
		engine.GET("/squads/:squadID", app.NewGinHandler(handler.GetOneByID, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/squads/"+squadId, nil)

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := fmt.Sprintf(`{
			"status":"error",
			"message":"%s"
		}`, invalidIdError.Error())
		assert.Equal(t, 400, rec.Code)
		assert.JSONEq(t, want, resp)
		mockStorage.Verify(t)
	})

	t.Run("should return 404 and data is nil when cannot found match Squad", func(t *testing.T) {
		squadId := mockObjectId(10)
		mockStorage := &mockSquadStorage{
			err: squadNotFoundError,
		}
		mockStorage.ExpectToCall("GetOneByID")
		handler := NewSquadHandler(mockStorage)

		engine := gin.New()
		engine.GET("/squads/:squadID", app.NewGinHandler(handler.GetOneByID, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/squads/"+squadId.hexId, nil)

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := fmt.Sprintf(`{
			"status":"error",
			"message":"%s"
		}`, squadNotFoundError.Error())
		assert.Equal(t, 404, rec.Code)
		assert.JSONEq(t, want, resp)
		mockStorage.Verify(t)
	})

	t.Run("should return 500 internal server error from storage", func(t *testing.T) {
		squadId := mockObjectId(10)
		mockStorage := &mockSquadStorage{
			err: errors.New("error from storage"),
		}
		mockStorage.ExpectToCall("GetOneByID")
		handler := NewSquadHandler(mockStorage)

		engine := gin.New()
		engine.GET("/squads/:squadID", app.NewGinHandler(handler.GetOneByID, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/squads/"+squadId.hexId, nil)

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := `{
			"status":"error",
			"message":"error from storage"
		}`
		assert.Equal(t, 500, rec.Code)
		assert.JSONEq(t, want, resp)
		mockStorage.Verify(t)
	})
}

func TestSquadHandlerInsertOne(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("should return 200 and insert squad success", func(t *testing.T) {
		squadId := mockObjectId(100)
		skillIds := []mockingObjectId{
			mockObjectId(10),
			mockObjectId(11),
			mockObjectId(12),
		}
		userIds := []GoogleUserId{
			mockGoogleUserId(20),
			mockGoogleUserId(21),
			mockGoogleUserId(22),
		}

		body := fmt.Sprintf(`{
			"name": "Aster",
			"teamleadMail": "john.d@arise.tech",
			"desc": "Squad in Blockchain team krub",
			"skillsRatings": [
				{
					"skid": "%s",
					"ratings": [
						{
							"uid": "%s",
							"score": 5
						}
					]
				},
				{
					"skid": "%s",
					"ratings": [
						{
							"uid": "%s",
							"score": 5
						}
					]
				},
				{
					"skid": "%s",
					"ratings": [
						{
							"uid": "%s",
							"score": 5
						}
					]
				}
			]
		}`, skillIds[0].hexId, userIds[0], skillIds[1].hexId, userIds[1], skillIds[2].hexId, userIds[2])

		mockStorage := &mockSquadStorage{
			squad: []*Squad{
				{
					Id:           squadId.objectId,
					Name:         "Aster",
					TeamleadMail: "john.d@arise.tech",
					Description:  "Squad in Blockchain team krub",
					SkillsRatings: []SkillRatings{
						{
							SkillId: skillIds[0].objectId,
							Ratings: []SkillRating{
								{
									UserId: userIds[0],
									Score:  5,
								},
							},
						},
						{
							SkillId: skillIds[1].objectId,
							Ratings: []SkillRating{
								{
									UserId: userIds[1],
									Score:  5,
								},
							},
						},
						{
							SkillId: skillIds[2].objectId,
							Ratings: []SkillRating{
								{
									UserId: userIds[2],
									Score:  5,
								},
							},
						},
					},
				},
			},
			err: nil,
		}
		mockStorage.ExpectToCall("InsertOne")
		handler := NewSquadHandler(mockStorage)

		engine := gin.New()
		engine.POST("/squads", app.NewGinHandler(handler.InsertOneByID, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/squads", strings.NewReader(body))

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := fmt.Sprintf(`{
			"status": "success",
			"message": "",
			"data": {
				"id": "%s",
				"name": "Aster",
				"teamleadMail": "john.d@arise.tech",
				"desc": "Squad in Blockchain team krub",
				"createdAt": "0001-01-01T00:00:00Z",
				"skillsRatings": [
					{
						"skid": "%s",
						"ratings": [
							{
								"uid": "%s",
								"score": 5
							}
						]
					},
					{
						"skid": "%s",
						"ratings": [
							{
								"uid": "%s",
								"score": 5
							}
						]
					},
					{
						"skid": "%s",
						"ratings": [
							{
								"uid": "%s",
								"score": 5
							}
						]
					}
				]
			}
		}`, squadId.hexId, skillIds[0].hexId, userIds[0], skillIds[1].hexId, userIds[1], skillIds[2].hexId, userIds[2])

		assert.Equal(t, 200, rec.Code)
		assert.JSONEq(t, want, resp)
	})

	t.Run("should return 400 and invalid input error", func(t *testing.T) {
		body := ""

		mockStorage := &mockSquadStorage{}
		mockStorage.ExpectToCall("InsertOne")
		handler := NewSquadHandler(mockStorage)

		engine := gin.New()
		engine.POST("/squads", app.NewGinHandler(handler.InsertOneByID, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/squads", strings.NewReader(body))

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := fmt.Sprintf(`{"status":"error","message":"%s"}`, invalidSquadInputError.Error())
		assert.Equal(t, 400, rec.Code)
		assert.JSONEq(t, want, resp)
	})

	t.Run("should return 400 and error amount of skills must be at least 3", func(t *testing.T) {
		skillId := mockObjectId(10)
		userId := mockObjectId(20)
		body := fmt.Sprintf(`{
			"name": "Aster",
			"teamleadMail": "john.d@arise.tech",
			"desc": "Squad in Blockchain team krub",
			"skillsRatings": [
				{
					"skid": "%s",
					"ratings": [
						{
							"uid": "%s",
							"score": 5
						}
					]
				}
			]
		}`, skillId.hexId, userId.hexId)

		mockStorage := &mockSquadStorage{}
		mockStorage.ExpectToCall("InsertOne")
		handler := NewSquadHandler(mockStorage)

		engine := gin.New()
		engine.POST("/squads", app.NewGinHandler(handler.InsertOneByID, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/squads", strings.NewReader(body))

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := fmt.Sprintf(`{
			"status":"error",
			"message":"%s"
		}`, missingSkillsRatingsError.Error())
		assert.Equal(t, 400, rec.Code)
		assert.JSONEq(t, want, resp)
	})

	t.Run("should return 400 and error skill ratings should have at least 1 rating from user", func(t *testing.T) {
		skillIds := []mockingObjectId{
			mockObjectId(10),
			mockObjectId(11),
			mockObjectId(12),
		}
		body := fmt.Sprintf(`{
			"name": "Aster",
			"teamleadMail": "john.d@arise.tech",
			"desc": "Squad in Blockchain team krub",
			"skillsRatings": [
				{
					"skid": "%s",
					"ratings": []
				},
				{
					"skid": "%s",
					"ratings": []
				},
				{
					"skid": "%s",
					"ratings": []
				}
			]
		}`, skillIds[0].hexId, skillIds[1].hexId, skillIds[2].hexId)

		mockStorage := &mockSquadStorage{
			err: errors.New("error from storage"),
		}
		mockStorage.ExpectToCall("InsertOne")
		handler := NewSquadHandler(mockStorage)

		engine := gin.New()
		engine.POST("/squads", app.NewGinHandler(handler.InsertOneByID, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/squads", strings.NewReader(body))

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := fmt.Sprintf(`{
			"status":"error",
			"message":"%s"}
		`, missingRatingsError.Error())
		assert.Equal(t, 400, rec.Code)
		assert.JSONEq(t, want, resp)
	})

	t.Run("should return 500 and error from storage", func(t *testing.T) {
		skillIds := []mockingObjectId{
			mockObjectId(10),
			mockObjectId(11),
			mockObjectId(12),
		}
		userIds := []mockingObjectId{
			mockObjectId(20),
			mockObjectId(21),
			mockObjectId(22),
		}
		body := fmt.Sprintf(`{
			"name": "Aster",
			"teamleadMail": "john.d@arise.tech",
			"desc": "Squad in Blockchain team krub",
			"skillsRatings": [
				{
					"skid": "%s",
					"ratings": [
						{
							"uid": "%s",
							"score": 5
						}
					]
				},
				{
					"skid": "%s",
					"ratings": [
						{
							"uid": "%s",
							"score": 5
						}
					]
				},
				{
					"skid": "%s",
					"ratings": [
						{
							"uid": "%s",
							"score": 5
						}
					]
				}
			]
		}`, skillIds[0].hexId, userIds[0].hexId, skillIds[1].hexId, userIds[1].hexId, skillIds[2].hexId, userIds[2].hexId)

		mockStorage := &mockSquadStorage{
			err: errors.New("error from storage"),
		}
		mockStorage.ExpectToCall("InsertOneByID")
		handler := NewSquadHandler(mockStorage)

		engine := gin.New()
		engine.POST("/squads", app.NewGinHandler(handler.InsertOneByID, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/squads", strings.NewReader(body))

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := `{
			"status":"error",
			"message":"error from storage"
		}`
		assert.Equal(t, 500, rec.Code)
		assert.JSONEq(t, want, resp)
		mockStorage.Verify(t)
	})
}

func TestSquadHandlerDeleteByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("should return 200 and success delete", func(t *testing.T) {
		squadId := mockObjectId(100)
		mockStorage := &mockSquadStorage{
			err: nil,
		}
		mockStorage.ExpectToCall("DeleteByID")
		handler := NewSquadHandler(mockStorage)

		engine := gin.New()
		engine.DELETE("/squad/:squadID", app.NewGinHandler(handler.DeleteByID, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/squad/"+squadId.hexId, nil)

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := `{
			"status":"success",
			"message":""
		}`
		assert.Equal(t, 200, rec.Code)
		assert.JSONEq(t, want, resp)
		mockStorage.Verify(t)
	})

	t.Run("should return 400 when invalid Id", func(t *testing.T) {
		squadId := "55"
		mockStorage := &mockSquadStorage{
			err: invalidIdError,
		}
		mockStorage.ExpectToCall("DeleteByID")
		handler := NewSquadHandler(mockStorage)

		engine := gin.New()
		engine.DELETE("/squad/:squadID", app.NewGinHandler(handler.DeleteByID, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/squad/"+squadId, nil)

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := fmt.Sprintf(`{
			"status":"error",
			"message":"%s"
		}`, invalidIdError.Error())
		assert.Equal(t, 400, rec.Code)
		assert.JSONEq(t, want, resp)
		mockStorage.Verify(t)
	})

	t.Run("should return 404 when not found match squad", func(t *testing.T) {
		squadId := mockObjectId(100)
		mockStorage := &mockSquadStorage{
			err: squadNotFoundError,
		}
		mockStorage.ExpectToCall("DeleteByID")
		handler := NewSquadHandler(mockStorage)

		engine := gin.New()
		engine.DELETE("/squad/:squadID", app.NewGinHandler(handler.DeleteByID, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/squad/"+squadId.hexId, nil)

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := fmt.Sprintf(`{
			"status":"error",
			"message":"%s"
		}`, squadNotFoundError.Error())
		assert.Equal(t, 404, rec.Code)
		assert.JSONEq(t, want, resp)
		mockStorage.Verify(t)
	})

	t.Run("should return 500 intenal server error", func(t *testing.T) {
		squadId := mockObjectId(100)
		mockStorage := &mockSquadStorage{
			err: errors.New("error from storage"),
		}
		mockStorage.ExpectToCall("DeleteByID")
		handler := NewSquadHandler(mockStorage)

		engine := gin.New()
		engine.DELETE("/squad/:squadID", app.NewGinHandler(handler.DeleteByID, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/squad/"+squadId.hexId, nil)

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := `{
			"status":"error",
			"message":"error from storage"
		}`
		assert.Equal(t, 500, rec.Code)
		assert.JSONEq(t, want, resp)
		mockStorage.Verify(t)
	})
}

func TestSquadHandlerGetAvgSkillRatingByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("should return 200 and squad average skill rating when found a squad match with squadID", func(t *testing.T) {
		squadId := mockObjectId(100)
		skillIds := []mockingObjectId{
			mockObjectId(10),
			mockObjectId(11),
			mockObjectId(12),
		}
		userIds := []GoogleUserId{
			mockGoogleUserId(20),
			mockGoogleUserId(21),
			mockGoogleUserId(22),
		}
		mockStorage := &mockSquadStorage{
			squad: []*Squad{
				{
					Id:           squadId.objectId,
					Name:         "Aster",
					Description:  "Squad in Blockchain team krub",
					TeamleadMail: "john.d@arise.tech",
					SkillsRatings: []SkillRatings{
						{
							SkillId: skillIds[0].objectId,
							Ratings: []SkillRating{
								{
									UserId: userIds[0],
									Score:  5,
								},
							},
						},
						{
							SkillId: skillIds[1].objectId,
							Ratings: []SkillRating{
								{
									UserId: userIds[1],
									Score:  5,
								},
							},
						},
						{
							SkillId: skillIds[2].objectId,
							Ratings: []SkillRating{
								{
									UserId: userIds[2],
									Score:  5,
								},
							},
						},
					},
				},
			},
			err: nil,
		}
		mockStorage.ExpectToCall("GetOneByID")
		handler := NewSquadHandler(mockStorage)

		engine := gin.New()
		engine.GET("/squads/:squadID/skills-require-avg", app.NewGinHandler(handler.GetAvgSkillRatingByID, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/squads/"+squadId.hexId+"/skills-require-avg", nil)

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := fmt.Sprintf(`{
			"status": "success",
			"message": "",
			"data": {
				"squadId": "%s",
				"averagesSkill": [
					{
						"skid": "%s",
						"average": 5
					},
					{
						"skid": "%s",
						"average": 5
					},
					{
						"skid": "%s",
						"average": 5
					}
				]
			}
		}`, squadId.hexId, skillIds[0].hexId, skillIds[1].hexId, skillIds[2].hexId)

		assert.Equal(t, 200, rec.Code)
		assert.JSONEq(t, want, resp)
		mockStorage.Verify(t)
	})

	t.Run("should return 400 and data is nil when Id is invalid", func(t *testing.T) {
		squadId := "55"
		mockStorage := &mockSquadStorage{
			err: invalidIdError,
		}
		mockStorage.ExpectToCall("GetOneByID")
		handler := NewSquadHandler(mockStorage)

		engine := gin.New()
		engine.GET("/squads/:squadID/skills-require-avg", app.NewGinHandler(handler.GetAvgSkillRatingByID, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/squads/"+squadId+"/skills-require-avg", nil)

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := fmt.Sprintf(`{
			"status":"error",
			"message":"%s"
		}`, invalidIdError.Error())
		assert.Equal(t, 400, rec.Code)
		assert.JSONEq(t, want, resp)
		mockStorage.Verify(t)
	})

	t.Run("should return 404 and data is nil when cannot found match Squad", func(t *testing.T) {
		squadId := mockObjectId(100)
		mockStorage := &mockSquadStorage{
			err: squadNotFoundError,
		}
		mockStorage.ExpectToCall("GetOneByID")
		handler := NewSquadHandler(mockStorage)

		engine := gin.New()
		engine.GET("/squads/:squadID/skills-require-avg", app.NewGinHandler(handler.GetAvgSkillRatingByID, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/squads/"+squadId.hexId+"/skills-require-avg", nil)

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := fmt.Sprintf(`{
			"status":"error",
			"message":"%s"
		}`, squadNotFoundError.Error())
		assert.Equal(t, 404, rec.Code)
		assert.JSONEq(t, want, resp)
		mockStorage.Verify(t)
	})

	t.Run("should return 500 and when found another error from storage", func(t *testing.T) {
		squadId := mockObjectId(100)
		mockStorage := &mockSquadStorage{
			err: errors.New("error from storage"),
		}
		mockStorage.ExpectToCall("GetOneByID")
		handler := NewSquadHandler(mockStorage)

		engine := gin.New()
		engine.GET("/squads/:squadID/skills-require-avg", app.NewGinHandler(handler.GetAvgSkillRatingByID, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/squads/"+squadId.hexId+"/skills-require-avg", nil)

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := `{
			"status":"error",
			"message":"error from storage"
		}`
		assert.Equal(t, 500, rec.Code)
		assert.JSONEq(t, want, resp)
		mockStorage.Verify(t)
	})
}

func TestSquadHandlerCalculateSquadMemberAveragePerSkill(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("should return 200 and success calculate", func(t *testing.T) {
		squadId := mockObjectId(100)
		skillIds := []mockingObjectId{
			mockObjectId(10),
			mockObjectId(11),
			mockObjectId(12),
		}
		userIds := []GoogleUserId{
			mockGoogleUserId(20),
			mockGoogleUserId(21),
			mockGoogleUserId(22),
		}
		usersSkillScores := [][]int{
			{5, 5, 5},
			{5, 6, 7},
			{5, 7, 9},
		}
		averageSkillScores := []float64{5, 6, 7}
		mockStorage := &mockSquadStorage{
			squad: []*Squad{
				{
					Id:           squadId.objectId,
					Name:         "Aster",
					Description:  "Squad in Blockchain team krub",
					TeamleadMail: "john.d@arise.tech",
					SkillsRatings: []SkillRatings{
						{
							SkillId: skillIds[0].objectId,
							Ratings: []SkillRating{
								{
									UserId: userIds[0],
									Score:  usersSkillScores[0][0],
								},
								{
									UserId: userIds[1],
									Score:  usersSkillScores[1][0],
								},
								{
									UserId: userIds[2],
									Score:  usersSkillScores[2][0],
								},
							},
						},
						{
							SkillId: skillIds[1].objectId,
							Ratings: []SkillRating{
								{
									UserId: userIds[0],
									Score:  usersSkillScores[0][1],
								},
								{
									UserId: userIds[1],
									Score:  usersSkillScores[1][1],
								},
								{
									UserId: userIds[2],
									Score:  usersSkillScores[2][1],
								},
							},
						},
						{
							SkillId: skillIds[2].objectId,
							Ratings: []SkillRating{
								{
									UserId: userIds[0],
									Score:  usersSkillScores[0][2],
								},
								{
									UserId: userIds[1],
									Score:  usersSkillScores[1][2],
								},
								{
									UserId: userIds[2],
									Score:  usersSkillScores[2][2],
								},
							},
						},
					},
				},
			},
		}
		mockStorage.ExpectToCall("GetOneByID")

		mockStorage.users = []user.User{
			{
				ID: string(userIds[0]),
				TechnicalSkill: []user.MySkill{
					{
						SkillID: skillIds[0].objectId,
						Score:   usersSkillScores[0][0],
					},
					{
						SkillID: skillIds[1].objectId,
						Score:   usersSkillScores[0][1],
					},
					{
						SkillID: skillIds[2].objectId,
						Score:   usersSkillScores[0][2],
					},
				},
			},
			{
				ID: string(userIds[1]),
				TechnicalSkill: []user.MySkill{
					{
						SkillID: skillIds[0].objectId,
						Score:   usersSkillScores[1][0],
					},
					{
						SkillID: skillIds[1].objectId,
						Score:   usersSkillScores[1][1],
					},
					{
						SkillID: skillIds[2].objectId,
						Score:   usersSkillScores[1][2],
					},
				},
			},
			{
				ID: string(userIds[2]),
				TechnicalSkill: []user.MySkill{
					{
						SkillID: skillIds[0].objectId,
						Score:   usersSkillScores[2][0],
					},
					{
						SkillID: skillIds[1].objectId,
						Score:   usersSkillScores[2][1],
					},
					{
						SkillID: skillIds[2].objectId,
						Score:   usersSkillScores[2][2],
					},
				},
			},
		}
		mockStorage.ExpectToCall("GetAllBySquadId")
		handler := NewSquadHandler(mockStorage)

		engine := gin.New()
		engine.GET("/squads/:squadID/member-skills-avg", app.NewGinHandler(handler.CalculateSquadMemberAveragePerSkill, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/squads/"+squadId.objectId.String()+"/member-skills-avg", nil)
		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := fmt.Sprintf(`{
			"status": "success",
			"message": "",
			"data": {
				"squadId": "%s",
				"averagesSkill": [
					{
						"skid": "%s",
						"average": %f
					},
					{
						"skid": "%s",
						"average": %f
					},
					{
						"skid": "%s",
						"average": %f
					}
				]
			}
		}`, squadId.hexId, skillIds[0].hexId, averageSkillScores[0], skillIds[1].hexId, averageSkillScores[1], skillIds[2].hexId, averageSkillScores[2])

		assert.Equal(t, 200, rec.Code)
		assert.JSONEq(t, want, resp)

		mockStorage.Verify(t)
	})
}

func TestSquadHandlerUpdateOneByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("should return 200 with the updated squad", func(t *testing.T) {
		squadId := mockObjectId(20)
		mockStorage := &mockSquadStorage{
			squad: []*Squad{
				{
					Id:            squadId.objectId,
					Name:          "Aster",
					TeamleadMail:  "john.d@arise.tech",
					Description:   "Squad in Blockchain team krub",
					SkillsRatings: []SkillRatings{},
				},
			},
		}
		mockStorage.ExpectToCall("GetOneByID")
		mockStorage.ExpectToCall("UpdateOneByID")
		handler := NewSquadHandler(mockStorage)

		engine := gin.New()
		engine.PUT("/squads/:squadID", app.NewGinHandler(handler.UpdateOneByID, zap.NewNop()))
		rec := httptest.NewRecorder()
		body := `{
			"name": "Aster",
			"desc": "This is modified."
		}`
		req, _ := http.NewRequest(http.MethodPut, "/squads/"+squadId.hexId, strings.NewReader(body))
		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := fmt.Sprintf(`{
			"status": "success",
			"message": "",
			"data": {
				"id": "%s",
				"name": "Aster",
				"desc": "This is modified.",
				"teamleadMail":"john.d@arise.tech",
				"createdAt": "0001-01-01T00:00:00Z",
				"skillsRatings": []
			}
		}`, squadId.hexId)

		assert.Equal(t, 200, rec.Code)
		assert.JSONEq(t, want, resp)
		mockStorage.Verify(t)
	})

	t.Run("should return 400 if got invalid input", func(t *testing.T) {
		squadId := mockObjectId(20)
		mockStorage := &mockSquadStorage{
			squad: []*Squad{
				{
					Id:            squadId.objectId,
					Name:          "Aster",
					TeamleadMail:  "john.d@arise.tech",
					Description:   "Squad in Blockchain team krub",
					SkillsRatings: []SkillRatings{},
				},
			},
		}
		mockStorage.ExpectToCall("GetOneByID")

		handler := NewSquadHandler(mockStorage)

		engine := gin.New()
		engine.PUT("/squads/:squadID", app.NewGinHandler(handler.UpdateOneByID, zap.NewNop()))
		rec := httptest.NewRecorder()
		body := ``
		req, _ := http.NewRequest(http.MethodPut, "/squads/"+squadId.hexId, strings.NewReader(body))
		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := fmt.Sprintf(`{
			"status": "error",
			"message": "%s"
		}`, invalidSquadInputError.Error())

		assert.Equal(t, 400, rec.Code)
		assert.JSONEq(t, want, resp)
		mockStorage.Verify(t)
	})

	t.Run("should return 404 if update non-existing squad", func(t *testing.T) {
		squadId := mockObjectId(20)
		mockStorage := &mockSquadStorage{
			err: squadNotFoundError,
		}
		mockStorage.ExpectToCall("GetOneByID")

		handler := NewSquadHandler(mockStorage)

		engine := gin.New()
		engine.PUT("/squads/:squadID", app.NewGinHandler(handler.UpdateOneByID, zap.NewNop()))
		rec := httptest.NewRecorder()
		body := ``
		req, _ := http.NewRequest(http.MethodPut, "/squads/"+squadId.hexId, strings.NewReader(body))
		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := fmt.Sprintf(`{
			"status": "error",
			"message": "%s"
		}`, squadNotFoundError.Error())
		assert.Equal(t, 404, rec.Code)
		assert.JSONEq(t, want, resp)
		mockStorage.Verify(t)
	})
}
