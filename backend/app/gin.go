package app

import (
	gcontext "context"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type context struct {
	*gin.Context
	logger    *zap.Logger
	validator *validator.Validate
}

type Context interface {
	Bind(v any) error
	Validate(v any) ([]ErrorField, error)
	OK(v any)
	BadRequest(err error)
	StoreError(err error)
	InternalServerError(err error)
	NotFound(err error)
	JSON(code int, v any)
	Ctx() gcontext.Context
	GetString(key string) string
	ShouldBindJSON(v any) error
	Param(key string) string
	Query(key string) string
}

func NewContext(c *gin.Context, logger *zap.Logger) Context {
	validate := validator.New()
	return &context{Context: c, logger: logger, validator: validate}
}

func (c *context) Bind(v any) error {
	return c.Context.ShouldBindJSON(v)
}

func (c *context) Validate(v any) ([]ErrorField, error) {
	if err := c.validator.Struct(v); err != nil {
		var fields []ErrorField
		for _, v := range err.(validator.ValidationErrors) {
			errField := ErrorField{
				Value: v.Param(),
				Field: v.Field(),
				Tag:   v.Tag(),
			}
			fields = append(fields, errField)
		}
		return fields, err
	}
	return nil, nil
}

func (c *context) ValidateFail(field []ErrorField, err error) {
	c.Context.JSON(http.StatusBadRequest, Response{
		Status:  Fail,
		Message: err.Error(),
		Data:    field,
	})
}

func (c *context) OK(v any) {
	c.Context.JSON(http.StatusOK, Response{
		Status: Success,
		Data:   v,
	})
}

func (c *context) BadRequest(err error) {
	c.logger.Error(err.Error())
	c.Context.JSON(http.StatusBadRequest, Response{
		Status:  Fail,
		Message: err.Error(),
	})
}

func (c *context) StoreError(err error) {
	c.logger.Error(err.Error())
	c.Context.JSON(storeErrorStatus, Response{
		Status:  Fail,
		Message: err.Error(),
	})
}

func (c *context) InternalServerError(err error) {
	c.logger.Error(err.Error())
	c.Context.JSON(http.StatusInternalServerError, Response{
		Status:  Fail,
		Message: err.Error(),
	})
}

func (c *context) NotFound(err error) {
	c.logger.Error(err.Error())
	c.Context.JSON(http.StatusNotFound, Response{
		Status:  Fail,
		Message: err.Error(),
	})
}

func (c *context) JSON(code int, v any) {
	c.Context.JSON(code, v)
}

func (c *context) ShouldBindJSON(v any) error {
	return c.Context.ShouldBindJSON(v)
}

func (c *context) Ctx() gcontext.Context {
	return c.Context.Request.Context()
}

func (c *context) GetString(key string) string {
	return c.Context.GetString(key)
}

func (c *context) Param(key string) string {
	return c.Context.Param(key)
}

func (c *context) Query(key string) string {
	return c.Context.Query(key)
}

func NewGinHandler(handler func(Context), logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler(NewContext(c, logger.With(zap.String("transaction-id", c.Request.Header.Get("transaction-id")))))
	}
}

func NewSwaggerHandler() func(Context) {
	return func(c Context) {
		ginSwagger.WrapHandler(swaggerFiles.Handler)(c.(*context).Context)
	}
}

type Router struct {
	*gin.Engine
	logger *zap.Logger
}

func NewRouter(logger *zap.Logger) *Router {
	r := gin.Default()

	config := cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"X-Requested-With", "Authorization", "Origin", "Content-Length", "Content-Type", "TransactionID"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}

	r.Use(cors.New(config))

	return &Router{Engine: r, logger: logger}
}

func (r *Router) GET(path string, handler func(Context)) {
	r.Engine.GET(path, NewGinHandler(handler, r.logger))
}

func (r *Router) POST(path string, handler func(Context)) {
	r.Engine.POST(path, NewGinHandler(handler, r.logger))
}

func (r *Router) PUT(path string, handler func(Context)) {
	r.Engine.PUT(path, NewGinHandler(handler, r.logger))
}

func (r *Router) DELETE(path string, handler func(Context)) {
	r.Engine.DELETE(path, NewGinHandler(handler, r.logger))
}
