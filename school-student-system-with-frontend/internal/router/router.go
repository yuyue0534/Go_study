package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"school-student-system/internal/handler"
	"school-student-system/internal/middleware"
)

func New(studentHandler *handler.StudentHandler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), middleware.CORS())

	api := r.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"status": "ok"}})
		})

		students := api.Group("/students")
		{
			students.POST("", studentHandler.Create)
			students.GET("", studentHandler.List)
			students.GET("/by-no/:student_no", studentHandler.GetByStudentNo)
			students.GET("/:id", studentHandler.GetByID)
			students.PUT("/:id", studentHandler.Update)
			students.DELETE("/:id", studentHandler.Delete)
		}
	}

	return r
}
