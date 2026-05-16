package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"school-student-system/internal/service"
)

type responseEnvelope struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func ok(c *gin.Context, data any) {
	c.JSON(http.StatusOK, responseEnvelope{Code: 0, Message: "success", Data: data})
}

func created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, responseEnvelope{Code: 0, Message: "success", Data: data})
}

func deleted(c *gin.Context) {
	c.JSON(http.StatusOK, responseEnvelope{Code: 0, Message: "success"})
}

func fail(c *gin.Context, err error) {
	var appErr *service.AppError
	if errors.As(err, &appErr) {
		c.JSON(appErr.HTTPStatus, responseEnvelope{Code: appErr.Code, Message: appErr.Message})
		return
	}

	c.JSON(http.StatusInternalServerError, responseEnvelope{Code: 500, Message: "internal server error"})
}
