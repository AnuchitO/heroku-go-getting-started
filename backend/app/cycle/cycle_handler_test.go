package cycle

import (
	"encoding/json"
	"errors"
	"fmt"
	_ "io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"gitdev.devops.krungthai.com/aster/ariskill/app"
	"gitdev.devops.krungthai.com/aster/ariskill/app/user"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type MockStorageA struct {
	Storage
}

func TestCycleHandlerGetOneByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should return 200 and cycle when found a cycle match with cycleID", func(t *testing.T) {
		cycleId := mockObjectId(10)
		startDate := time.Now()
		startDateFormatted := startDate.Format(time.RFC3339Nano)
		endDate := time.Now().AddDate(1, 0, 0)
		endDateFormatted := endDate.Format(time.RFC3339Nano)
		mockStorage := &mockCycleStorage{
			newCycles: []*NewCycle{
				{
					ID:             cycleId.objectId,
					AriserMail:     "ariser@arise.tech",
					TeamLeaderMail: "teamleader@arise.tech",
					StartDate:      startDate,
					EndDate:        endDate,
					HardSkills:     []HardSkill{},
					Status:         "pending",
					Comment:        "",
				},
			},
			err: nil,
		}

		mockStorage.ExpectToCall("GetNewByID")
		handler := NewCycleHandler(mockStorage)

		engine := gin.New()
		engine.GET("/cycles/:id", app.NewGinHandler(handler.GetOneByID, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/cycles/"+cycleId.hexId, nil)

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := fmt.Sprintf(`{
			"status": "success",
			"message": "",
			"data": {
				"id": "000000000000000000000010",
				"givenName": "Ariser",
				"familyName": "by krungthai",
				"jobRole": "full-stack",
				"level": "junior",
				"ariserMail": "ariser@arise.tech",
				"teamLeaderMail": "teamleader@arise.tech",
				"startDate": "%s",
				"endDate": "%s",
				"status": "pending",
				"state": "",
				"comment": ""
			}
		}`, startDateFormatted, endDateFormatted)
		assert.Equal(t, 200, rec.Code)
		assert.JSONEq(t, want, resp)
		mockStorage.Verify(t)
	})

	t.Run("should return 400 and data is nil when invalid Id", func(t *testing.T) {
		cycleId := "55"
		mockStorage := &mockCycleStorage{
			err: invalidRequestError,
		}
		mockStorage.ExpectToCall("GetNewByID")

		handler := NewCycleHandler(mockStorage)

		engine := gin.New()
		engine.GET("/cycles/:id", app.NewGinHandler(handler.GetOneByID, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/cycles/"+cycleId, nil)

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := fmt.Sprintf(`{
			"status":"error",
			"message":"%s"
		}`, invalidRequestError.Error())
		assert.Equal(t, 400, rec.Code)
		assert.JSONEq(t, want, resp)
		mockStorage.Verify(t)
	})

	t.Run("should return 404 and data is nil when cannot found match Cycle", func(t *testing.T) {
		cycleId := mockObjectId(10)
		mockStorage := &mockCycleStorage{
			err: cycleNotFoundError,
		}
		mockStorage.ExpectToCall("GetNewByID")

		handler := NewCycleHandler(mockStorage)

		engine := gin.New()
		engine.GET("/cycles/:id", app.NewGinHandler(handler.GetOneByID, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/cycles/"+cycleId.hexId, nil)

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := fmt.Sprintf(`{
			"status":"error",
			"message":"%s"
		}`, cycleNotFoundError.Error())
		assert.Equal(t, 404, rec.Code)
		assert.JSONEq(t, want, resp)
		mockStorage.Verify(t)
	})

	t.Run("should return 500 internal server error from storage", func(t *testing.T) {
		cycleId := mockObjectId(10)
		mockStorage := &mockCycleStorage{
			err: errors.New("error from storage"),
		}
		mockStorage.ExpectToCall("GetNewByID")

		handler := NewCycleHandler(mockStorage)

		engine := gin.New()
		engine.GET("/cycles/:id", app.NewGinHandler(handler.GetOneByID, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/cycles/"+cycleId.hexId, nil)

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

func TestGetAllFromEmailWithMockContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cycles := randomCycles(20)
	expectCycles := getExpectCycles(cycles)

	testCases := []struct {
		name          string
		status        string
		page          int
		email         string
		cyclesReturn  []*Cycle
		err           error
		checkResponse func(t *testing.T, res interface{}, code int, storage *mockCycleStorage)
	}{
		{
			name:         "Should return 200 Status OK when input correct request",
			status:       StatusAll,
			page:         1,
			email:        expectCycles[0].ReceiverMail,
			cyclesReturn: expectCycles,
			err:          nil,
			checkResponse: func(t *testing.T, res interface{}, code int, storage *mockCycleStorage) {
				storage.Verify(t)
				require.Equal(t, http.StatusOK, code)
				require.Len(t, res, len(expectCycles))
			},
		},
		{
			name:         "Should return 450 Status StoreError when db is error",
			status:       StatusAll,
			page:         1,
			email:        expectCycles[0].ReceiverMail,
			cyclesReturn: nil,
			err:          dbConnectNotFound,
			checkResponse: func(t *testing.T, res interface{}, code int, storage *mockCycleStorage) {
				storage.Verify(t)
				require.Equal(t, 450, code)
				require.Equal(t, dbConnectNotFound, res)
			},
		},
		{
			name:         "Should return 400 Status BadRequest when page is not valid",
			status:       StatusAll,
			page:         -1,
			email:        expectCycles[0].ReceiverMail,
			cyclesReturn: nil,
			err:          invalidRequestError,
			checkResponse: func(t *testing.T, res interface{}, code int, storage *mockCycleStorage) {
				require.Equal(t, http.StatusBadRequest, code)
				require.Equal(t, invalidRequestError, res)
			},
		},
		{
			name:         "Should return 400 Status BadRequest when email is not valid",
			status:       StatusAll,
			page:         1,
			email:        "not-valid",
			cyclesReturn: nil,
			err:          invalidRequestError,
			checkResponse: func(t *testing.T, res interface{}, code int, storage *mockCycleStorage) {
				require.Equal(t, http.StatusBadRequest, code)
				require.Equal(t, invalidRequestError, res)
			},
		},
		{
			name:         "Should return 400 Status BadRequest when status is not valid",
			status:       "test",
			page:         1,
			email:        expectCycles[0].ReceiverMail,
			cyclesReturn: nil,
			err:          invalidRequestError,
			checkResponse: func(t *testing.T, res interface{}, code int, storage *mockCycleStorage) {
				require.Equal(t, http.StatusBadRequest, code)
				require.Equal(t, invalidRequestError, res)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			mockStorage := &mockCycleStorage{
				cycles:       cycles,
				cyclesReturn: tc.cyclesReturn,
				err:          tc.err,
			}
			mockStorage.ExpectToCall("GetAllFromEmail")
			handler := NewCycleHandler(mockStorage)
			context := app.MockAppContext{
				Params: map[string]string{
					"status": tc.status,
					"page":   strconv.Itoa(tc.page),
				},
				Email: tc.email,
			}
			handler.GetAllFromReceiverEmail(&context)
			// log.Println(context.ResponseData)
			tc.checkResponse(t, context.ResponseData, context.ResponseCode, mockStorage)
		})
	}
}

// type CycleInput struct {
// 	ReceiverMail string    `json:"receiverMail" bson:"receiver_mail"`
// 	StartDate    time.Time `json:"startDate" bson:"start_date"`
// 	EndDate      time.Time `json:"endDate" bson:"end_date"`
// 	// QuantitiveSkill >= 1 skill
// 	QuantitativeSkill []QuantitativeSkill `json:"quantitativeSkill" bson:"quantitative_skill"`
// 	IntuitiveSkill    []IntuitiveSkill    `json:"intuitiveSkill" bson:"intuitive_skill"`
// 	Comment           string              `json:"comment" bson:"comment"`
// }

//	func fmtJsonBody(body map[string]any) *bytes.Buffer {
//		rsl, _ := json.Marshal(body)
//		return bytes.NewBuffer(rsl)
//	}
func TestInsertOneCycle(t *testing.T) {
	gin.SetMode("test")
	t.Run("should return 200 when include is true", func(t *testing.T) {
		reqBody :=
			`{
				"receiverMail":"prawith.a@arise.tech",
				"startDate":"2023-11-01T00:00:00.000Z",
				"endDate":"2023-11-30T00:00:00.000Z",
				"quantitativeSkill":[
					{
							"id":"64e17f43e098346113ae4f53",
							"personalScore":3,
							"goalScore":5,
							"finalScore":3,
							"comment":""
						},
						{
							"id":"64e17f43e098346113ae4f5c",
							"personalScore":2,
							"goalScore":3,
							"finalScore":2,
							"comment":""
						}
				],
				"intuitiveSkill":[],
				"comment":""
			}`

		mock := &mockCycleStorage{
			cycle: &Cycle{},
		}
		handler := NewCycleHandler(mock)

		engine := gin.New()
		engine.POST("/cycles", app.NewGinHandler(handler.InsertOne, zap.NewNop()))
		rec := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodPost, "/cycles", strings.NewReader(reqBody))
		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := `{
			"status": "success",
			"message": "",
			"data": {
				"comment":"",
				"endDate":"0001-01-01T00:00:00Z",
				"id":"000000000000000000000000",
				"intuitiveSkill":[],
				"quantitativeSkill":[],
				"receiverMail":"",
				"senderMail":"",
				"startDate":"0001-01-01T00:00:00Z",
				"status":"",
				"state":""
			}
		}`
		// fmt.Println(resp)
		assert.Equal(t, 200, rec.Code)
		assert.JSONEq(t, want, resp)
	})

	t.Run("Return 400 When Input in cycle is invalid", func(t *testing.T) {
		reqBody :=
			`{
				"startDate":"2023-11-01T00:00:00.000Z",
				"endDate":"2023-11-30T00:00:00.000Z",
				"quantitativeSkill":[
					{
							"id":"64e17f43e098346113ae4f53",
							"personalScore":3,
							"goalScore":5,
							"finalScore":3,
							"comment":""
						},
						{
							"id":"64e17f43e098346113ae4f5c",
							"personalScore":2,
							"goalScore":3,
							"finalScore":2,
							"comment":""
						}
				],
				"intuitiveSkill":[],
				"comment":""
			}`

		mock := &mockCycleStorage{
			cycle: &Cycle{},
		}
		handler := NewCycleHandler(mock)

		engine := gin.New()
		engine.POST("/cycles", app.NewGinHandler(handler.InsertOne, zap.NewNop()))
		rec := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodPost, "/cycles", strings.NewReader(reqBody))
		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := fmt.Sprintf(`{
			"status": "error",
			"message": "%s"
		}`, invalidInsertOneInputError)
		// fmt.Println(resp)
		assert.Equal(t, 400, rec.Code)
		assert.JSONEq(t, want, resp)
	})

	t.Run("should return 400 and message When input of Cycle is incorrect", func(t *testing.T) {
		reqBody :=
			`{
				"receiverMail":"prawith.a@arise.tech",
				"startDate":"2023-11-01T00:00:00.000Z",
				"endDate":"2023-11-30T00:00:00.000Z",
				"quantitativeSkill":[],
				"intuitiveSkill":[],
				"comment":""
			}`

		mock := &mockCycleStorage{
			cycle: &Cycle{},
		}
		handler := NewCycleHandler(mock)

		engine := gin.New()
		engine.POST("/cycles", app.NewGinHandler(handler.InsertOne, zap.NewNop()))
		rec := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodPost, "/cycles", strings.NewReader(reqBody))
		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := fmt.Sprintf(`{
			"status": "error",
			"message": "%s"
		}`, numberOfSkillError)
		// fmt.Println(resp)
		assert.Equal(t, 400, rec.Code)
		assert.JSONEq(t, want, resp)
	})

	t.Run("should return 500 when storage is failed", func(t *testing.T) {
		reqBody :=
			`{
				"receiverMail":"prawith.a@arise.tech",
				"startDate":"2023-11-01T00:00:00.000Z",
				"endDate":"2023-11-30T00:00:00.000Z",
				"quantitativeSkill":[
					{
							"id":"64e17f43e098346113ae4f53",
							"personalScore":3,
							"goalScore":5,
							"finalScore":3,
							"comment":""
						},
						{
							"id":"64e17f43e098346113ae4f5c",
							"personalScore":2,
							"goalScore":3,
							"finalScore":2,
							"comment":""
						}
				],
				"intuitiveSkill":[],
				"comment":""
			}`

		mock := &mockCycleStorage{
			err: errors.New("Error from mongo"),
		}
		handler := NewCycleHandler(mock)

		engine := gin.New()
		engine.POST("/cycles", app.NewGinHandler(handler.InsertOne, zap.NewNop()))
		rec := httptest.NewRecorder()
		// log.Println("463", rec)
		req, _ := http.NewRequest(http.MethodPost, "/cycles", strings.NewReader(reqBody))
		// log.Println("465", err)
		engine.ServeHTTP(rec, req)
		// log.Println(err, req, rec)
		resp := rec.Body.String()
		want := `{
			"status": "error",
			"message": "Error from mongo"
		}`

		// log.Println(rec.Code)
		assert.Equal(t, 500, rec.Code)
		assert.JSONEq(t, want, resp)
	})
}

func TestCycleHandlerDeleteByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("should return 200 and success delete", func(t *testing.T) {
		cycleId := mockObjectId(24)
		mockStorage := &mockCycleStorage{
			err: nil,
		}
		mockStorage.ExpectToCall("DeleteByID")
		handler := NewCycleHandler(mockStorage)

		engine := gin.New()
		engine.DELETE("/cycles/:id", app.NewGinHandler(handler.DeleteByID, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/cycles/"+cycleId.hexId, nil)

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := `{
			"status": "success",
			"message": ""
		}`
		assert.Equal(t, 200, rec.Code)
		assert.JSONEq(t, want, resp)
		mockStorage.Verify(t)
	})

	t.Run("should return 400 when invalid Id", func(t *testing.T) {
		cycleId := "dkls"
		mockStorage := &mockCycleStorage{
			err: invalidRequestError,
		}
		mockStorage.ExpectToCall("DeleteByID")

		handler := NewCycleHandler(mockStorage)

		engine := gin.New()
		engine.DELETE("/cycles/:id", app.NewGinHandler(handler.DeleteByID, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/cycles/"+cycleId, nil)

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := fmt.Sprintf(`{
				"status":"error",
				"message":"%s"
			}`, invalidRequestError.Error())
		assert.Equal(t, 400, rec.Code)
		assert.JSONEq(t, want, resp)
		mockStorage.Verify(t)
	})

	t.Run("should return 404 when not found match squad", func(t *testing.T) {
		cycleId := mockObjectId(24)
		mockStorage := &mockCycleStorage{
			err: cycleNotFoundError,
		}
		mockStorage.ExpectToCall("DeleteByID")

		handler := NewCycleHandler(mockStorage)

		engine := gin.New()
		engine.DELETE("/cycles/:id", app.NewGinHandler(handler.DeleteByID, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/cycles/"+cycleId.hexId, nil)

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := fmt.Sprintf(`{
			"status":"error",
			"message":"%s"
		}`, cycleNotFoundError.Error())
		assert.Equal(t, 404, rec.Code)
		assert.JSONEq(t, want, resp)
		mockStorage.Verify(t)
	})

	t.Run("should return 500 intenal server error", func(t *testing.T) {
		cycleId := mockObjectId(100)
		mockStorage := &mockCycleStorage{
			err: errors.New("error from storage"),
		}
		mockStorage.ExpectToCall("DeleteByID")

		handler := NewCycleHandler(mockStorage)

		engine := gin.New()
		engine.DELETE("/cycles/:id", app.NewGinHandler(handler.DeleteByID, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/cycles/"+cycleId.hexId, nil)

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

func randomCycles(n int) []*Cycle {
	var cycles []*Cycle

	recieverEmail := "ariser@arise.dev"
	for i := 0; i < n; i++ {
		cycle := Cycle{
			ID:                primitive.NewObjectID(),
			SenderMail:        "ariser@arise.dev",
			ReceiverMail:      recieverEmail,
			StartDate:         time.Now(),
			EndDate:           time.Now().Add(24 * 5 * time.Hour),
			QuantitativeSkill: []QuantitativeSkill{},
			IntuitiveSkill:    []IntuitiveSkill{},
			Status:            randomCycleStatus(),
		}

		cycles = append(cycles, &cycle)
	}

	return cycles
}

func randomCycleStatus() string {
	return StatusPending
}

func getExpectCycles(cycles []*Cycle) []*Cycle {
	email := cycles[0].ReceiverMail // random email from cycles
	var expectCycles []*Cycle
	for _, cy := range cycles {
		if cy.ReceiverMail == email {
			expectCycles = append(expectCycles, cy)
		}
	}

	return expectCycles
}

func TestUpdateById(t *testing.T) {
	t.Run("Return HTTP status 200 when update data by ID is working", func(t *testing.T) {
		cycle := randomCycles(1)
		reqBody := &Cycle{
			ID:                cycle[0].ID,
			SenderMail:        "updatetestById@arise.tech",
			ReceiverMail:      cycle[0].ReceiverMail,
			StartDate:         cycle[0].StartDate,
			EndDate:           cycle[0].EndDate,
			QuantitativeSkill: cycle[0].QuantitativeSkill,
			IntuitiveSkill:    cycle[0].IntuitiveSkill,
			Status:            cycle[0].Status,
			Comment:           cycle[0].Comment,
		}
		reqJson, err := json.Marshal(reqBody)
		require.NoError(t, err)

		mockStorage := &mockCycleStorage{
			cycles: cycle,
			err:    nil,
		}
		mockStorage.ExpectToCall("UpdateByID")
		handler := NewCycleHandler(mockStorage)
		engine := gin.New()

		engine.POST("/cycles/:id", app.NewGinHandler(handler.UpdateByID, zap.NewNop()))
		rec := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodPost, "/cycles/"+cycle[0].ID.Hex(), strings.NewReader(string(reqJson)))

		engine.ServeHTTP(rec, req)

		assert.Equal(t, 200, rec.Code)
		mockStorage.Verify(t)
	})
	t.Run("Return HTTP status 400 when user put invalid ID field or missing body", func(t *testing.T) {
		cycleId := "dkls"

		mockStorage := &mockCycleStorage{
			err: invalidRequestError,
		}
		mockStorage.ExpectToCall("UpdateByID")
		handler := NewCycleHandler(mockStorage)
		engine := gin.New()

		engine.POST("/cycles/:id", app.NewGinHandler(handler.UpdateByID, zap.NewNop()))
		rec := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodPost, "/cycles/"+cycleId, nil)

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := `{
			"status":"error",
			"message":"invalid request"
		}`
		assert.Equal(t, 400, rec.Code)
		assert.JSONEq(t, want, resp)
	})

	t.Run("Return HTTP status 500, Internal server error", func(t *testing.T) {
		cycle := randomCycles(1)
		reqBody := &Cycle{
			ID:                cycle[0].ID,
			SenderMail:        "updatetestById@arise.tech",
			ReceiverMail:      cycle[0].ReceiverMail,
			StartDate:         cycle[0].StartDate,
			EndDate:           cycle[0].EndDate,
			QuantitativeSkill: cycle[0].QuantitativeSkill,
			IntuitiveSkill:    cycle[0].IntuitiveSkill,
			Status:            cycle[0].Status,
			Comment:           cycle[0].Comment,
		}
		reqJson, err := json.Marshal(reqBody)
		require.NoError(t, err)

		mockStorage := &mockCycleStorage{
			err: cycleNotFoundError,
		}
		mockStorage.ExpectToCall("UpdateByID")
		handler := NewCycleHandler(mockStorage)
		engine := gin.New()

		engine.POST("/cycles/:id", app.NewGinHandler(handler.UpdateByID, zap.NewNop()))
		rec := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodPost, "/cycles/"+"123", strings.NewReader(string(reqJson)))

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := fmt.Sprintf(`{
			"status":"error",
			"message":"%s"
		}`, cycleNotFoundError.Error())
		assert.Equal(t, 500, rec.Code)
		assert.JSONEq(t, want, resp)
	})
}

func TestUpdateUserFinalScore(t *testing.T) {
	cycles := randomCycles(10)
	expectedCycles := getExpectCycles(cycles)

	testCases := []struct {
		name          string
		dbErr         error
		id            string
		checkResponse func(t *testing.T, res interface{}, code int, mockStorage *mockCycleStorage)
	}{
		{
			name:  "Should return 200 Status OK when input correct request",
			dbErr: nil,
			id:    expectedCycles[0].ID.Hex(),
			checkResponse: func(t *testing.T, res interface{}, code int, mockStorage *mockCycleStorage) {
				require.Equal(t, http.StatusOK, code)
				mockStorage.Verify(t)
			},
		},
		{
			name:  "Should return 400 Status BadRequest when cannot convert id to objectId",
			dbErr: nil,
			id:    "not-valid",
			checkResponse: func(t *testing.T, res interface{}, code int, mockStorage *mockCycleStorage) {
				require.Equal(t, http.StatusBadRequest, code)
			},
		},
		{
			name:  "Should return 450 Status StoreError when db is error",
			dbErr: mongo.ErrNilDocument,
			id:    expectedCycles[0].ID.Hex(),
			checkResponse: func(t *testing.T, res interface{}, code int, mockStorage *mockCycleStorage) {
				require.Equal(t, 450, code)
				mockStorage.Verify(t)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			mockStorage := &mockCycleStorage{
				cycles: cycles,
				err:    tc.dbErr,
			}
			mockStorage.ExpectToCall("UpdateUserFinalScore")

			handler := NewCycleHandler(mockStorage)

			engine := gin.New()
			engine.POST("/cycles/update/:id", app.NewGinHandler(handler.UpdateUserFinalScore, zap.NewNop()))
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/cycles/update/"+tc.id, nil)

			engine.ServeHTTP(rec, req)

			tc.checkResponse(t, rec.Body.String(), rec.Code, mockStorage)
		})
	}
}
func TestGetAllFromUserEmail(t *testing.T) {
	t.Run("Return HTTP status 200 when update data by ID is working", func(t *testing.T) {
		mockStorage := &mockCycleStorage{}
		mockStorage.ExpectToCall("GetFromUserEmail")
		handler := NewCycleHandler(mockStorage)
		engine := gin.New()

		engine.GET("/cycles/email/user", app.NewGinHandler(handler.GetAllFromUserEmail, zap.NewNop()))
		rec := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, "/cycles/email/user", nil)

		engine.ServeHTTP(rec, req)

		assert.Equal(t, 200, rec.Code)
		mockStorage.Verify(t)
	})
}

// func TestGetCycleProgess(t *testing.T) {
// 	t.Run("Return HTTP status 200 when update data by ID is working", func(t *testing.T) {

// 		mockStorage := &mockCycleStorage{}
// 		mockStorage.ExpectToCall("GetFromUserEmail")
// 		handler := NewCycleHandler(mockStorage)
// 		engine := gin.New()

// 		engine.GET("/cycles/email/user", app.NewGinHandler(handler.GetAllFromUserEmail, zap.NewNop()))
// 		rec := httptest.NewRecorder()

// 		req, _ := http.NewRequest(http.MethodGet, "/cycles/email/user", nil)

// 		engine.ServeHTTP(rec, req)

//			assert.Equal(t, 200, rec.Code)
//			mockStorage.Verify(t)
//		})
//	}

func TestUpdateNewCycle(t *testing.T) {
	testCases := []struct {
		name             string
		reqBody          string
		expectedStatus   int
		expectedResponse string
		storageError     error
	}{
		{
			name: "Return http status 200 when add data to latest cycle in progress by email",
			reqBody: `{
				"hardSkills": [
					{
						"name":"HTML",
						"description": "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
						"skillLevels":[
							{
								"level": 1,
								"levelDescription":"Lorem ipsum dolor sit amet"
							}
						],
						"personalScore": 2,
						"goalScore": 3
					},
					{
						"name":"CSS",
						"description": "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
						"skillLevels":[
							{
								"level": 1,
								"levelDescription":"Lorem ipsum dolor sit amet"
							}
						],
						"personalScore": 1,
						"goalScore": 2
					}
				]

			}`,
			expectedStatus:   200,
			expectedResponse: `{"status": "success","message": "","data": {}}`,
			storageError:     nil,
		},
		{
			name: "Return http status 400 when goal-score exceed personal-score more than 1",
			reqBody: `{
				"hardSkills": [
					{
						"name":"HTML",
						"description": "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
						"skillLevels":[
							{
								"level": 1,
								"levelDescription":"Lorem ipsum dolor sit amet"
							}
						],
						"personalScore": 2,
						"goalScore": 5
					},
					{
						"name":"CSS",
						"description": "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
						"skillLevels":[
							{
								"level": 1,
								"levelDescription":"Lorem ipsum dolor sit amet"
							}
						],
						"personalScore": 1,
						"goalScore": 4
					}
				]
			}`,
			expectedStatus:   400,
			expectedResponse: `{"status": "error","message": "miss match goal-score"}`,
			storageError:     nil,
		},
		{
			name:             "Return http status 400 when JSON binding fails",
			reqBody:          `invalid json`,
			expectedStatus:   400,
			expectedResponse: `{"status":"error","message":"Error marshaling JSON"}`,
			storageError:     nil,
		},
		{
			name: "Return http status 450 Status StoreError when db is error",
			reqBody: `{
				"hardSkills": [
					{
						"name":"HTML",
						"description": "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
						"skillLevels":[
							{
								"level": 1,
								"levelDescription":"Lorem ipsum dolor sit amet"
							}
						],
						"personalScore": 2,
						"goalScore": 3
					},
					{
						"name":"CSS",
						"description": "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
						"skillLevels":[
							{
								"level": 1,
								"levelDescription":"Lorem ipsum dolor sit amet"
							}
						],
						"personalScore": 1,
						"goalScore": 2
					}
				]
			}`,
			expectedStatus:   450,
			expectedResponse: `{"message":"document is nil", "status":"error"}`,
			storageError:     mongo.ErrNilDocument,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storedCycle := NewCycle{
				ID:             primitive.NewObjectID(),
				TeamLeaderMail: "test.t@arise.tech",
				AriserMail:     "test.a@ariser.tech",
				StartDate:      time.Now(),
				EndDate:        time.Now().Add(48 * time.Hour),
				Status:         "In Progress",
			}
			type getBody struct {
				HardSkill []user.MyHardSkill `json:"hardSkills"`
			}
			var hardSkills getBody
			if err := json.Unmarshal([]byte(tc.reqBody), &hardSkills); err != nil {
				hardSkills = getBody{}
			}

			userData := user.User{
				ID:         primitive.NewObjectID().String(),
				Email:      "test.a@ariser.tech",
				EmployeeID: "",
				HardSkills: hardSkills.HardSkill,
			}

			mockStorage := &mockNewCycleStorage{
				newCycle: &storedCycle,
				user:     userData,
				err:      tc.storageError,
			}
			mockHandler := NewCycleHandler(mockStorage)

			engine := gin.New()
			engine.PUT("/cycles/goal", app.NewGinHandler(
				mockHandler.UpdateHardSkillsByEmail, zap.NewNop(),
			))

			rec := httptest.NewRecorder()

			req, _ := http.NewRequest(http.MethodPut, "/cycles/goal", strings.NewReader(tc.reqBody))
			engine.ServeHTTP(rec, req)

			resp := rec.Body.String()
			// fmt.Printf("expect:%v\n", tc.expectedResponse)
			// fmt.Printf("actual:%v", rec)
			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.JSONEq(t, resp, tc.expectedResponse)
		})
	}
}

func TestGetLatestCycleFromUserEmail(t *testing.T) {
	id := primitive.NewObjectID()
	idHex := id.Hex()

	startTime := time.Now()
	endTime := time.Now().Add(48 * time.Hour)

	formattedStartTime := startTime.Format("2006-01-02T15:04:05.999999-07:00")
	formattedEndTime := endTime.Format("2006-01-02T15:04:05.999999-07:00")

	testCases := []struct {
		name             string
		reqBody          string
		expectedStatus   int
		expectedResponse string
		storageError     error
	}{
		{
			name:             "Return http status 200 and latest cycle by email when email is input",
			expectedStatus:   200,
			expectedResponse: fmt.Sprintf(`{"status":"success","message":"","data":{"id":"%v","teamLeaderMail":"test.t@arise.tech","ariserMail":"test.a@ariser.tech","comment":"","startDate":"%v","endDate":"%v", "state":"","status":"In Progress","hardSkills":null}}`, idHex, formattedStartTime, formattedEndTime),
			storageError:     nil,
		},
		{
			name:             "Return http status 400 when no document found",
			expectedStatus:   400,
			expectedResponse: `{"message":"no documents in result", "status":"error"}`,
			storageError:     mongo.ErrNoDocuments,
		},
	}
	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			storedCycle := NewCycle{
				ID:             id,
				TeamLeaderMail: "test.t@arise.tech",
				AriserMail:     "test.a@ariser.tech",
				StartDate:      startTime,
				Comment:        "",
				EndDate:        endTime,
				Status:         "In Progress",
			}

			mockStorage := &mockNewCycleStorage{
				newCycle: &storedCycle,
				err:      v.storageError,
			}

			mockHandler := NewCycleHandler(mockStorage)
			engine := gin.New()
			engine.GET("/cycles/email/lastest", app.NewGinHandler(
				mockHandler.GetLatestCycleFromUserEmail, zap.NewNop(),
			))
			rec := httptest.NewRecorder()

			req, _ := http.NewRequest(http.MethodGet, "/cycles/email/lastest", nil)

			engine.ServeHTTP(rec, req)

			actual := rec.Body.String()

			assert.Equal(t, v.expectedStatus, rec.Code)
			assert.JSONEq(t, v.expectedResponse, actual)
		})
	}
}
