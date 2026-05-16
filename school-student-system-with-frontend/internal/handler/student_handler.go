package handler

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"school-student-system/internal/dto"
	"school-student-system/internal/model"
	"school-student-system/internal/service"
)

type StudentHandler struct {
	service *service.StudentService
}

func NewStudentHandler(service *service.StudentService) *StudentHandler {
	return &StudentHandler{service: service}
}

func (h *StudentHandler) Create(c *gin.Context) {
	var req dto.CreateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, service.BadRequest("invalid JSON body"))
		return
	}

	student, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		fail(c, err)
		return
	}
	created(c, student)
}

func (h *StudentHandler) GetByID(c *gin.Context) {
	id, err := parsePositiveInt64(c.Param("id"), "student id")
	if err != nil {
		fail(c, err)
		return
	}

	student, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, student)
}

func (h *StudentHandler) GetByStudentNo(c *gin.Context) {
	studentNo := strings.TrimSpace(c.Param("student_no"))
	student, err := h.service.GetByStudentNo(c.Request.Context(), studentNo)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, student)
}

func (h *StudentHandler) Update(c *gin.Context) {
	id, err := parsePositiveInt64(c.Param("id"), "student id")
	if err != nil {
		fail(c, err)
		return
	}

	var req dto.UpdateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, service.BadRequest("invalid JSON body"))
		return
	}

	student, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, student)
}

func (h *StudentHandler) Delete(c *gin.Context) {
	id, err := parsePositiveInt64(c.Param("id"), "student id")
	if err != nil {
		fail(c, err)
		return
	}

	if err := h.service.SoftDelete(c.Request.Context(), id); err != nil {
		fail(c, err)
		return
	}
	deleted(c)
}

func (h *StudentHandler) List(c *gin.Context) {
	query, err := parseListQuery(c)
	if err != nil {
		fail(c, err)
		return
	}

	page, err := h.service.List(c.Request.Context(), query)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, page)
}

func parseListQuery(c *gin.Context) (model.StudentListQuery, error) {
	query := model.StudentListQuery{
		StudentNo: strings.TrimSpace(c.Query("student_no")),
		Name:      strings.TrimSpace(c.Query("name")),
		ClassName: strings.TrimSpace(c.Query("class_name")),
		MajorName: strings.TrimSpace(c.Query("major_name")),
		Page:      1,
		PageSize:  20,
	}

	if raw := strings.TrimSpace(c.Query("grade_year")); raw != "" {
		value, err := strconv.Atoi(raw)
		if err != nil {
			return query, service.BadRequest("grade_year must be an integer")
		}
		query.GradeYear = value
	}

	if raw := strings.TrimSpace(c.Query("class_id")); raw != "" {
		value, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return query, service.BadRequest("class_id must be an integer")
		}
		query.ClassID = value
	}

	if raw := strings.TrimSpace(c.Query("major_id")); raw != "" {
		value, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return query, service.BadRequest("major_id must be an integer")
		}
		query.MajorID = value
	}

	if raw := strings.TrimSpace(c.Query("status")); raw != "" {
		value, err := strconv.Atoi(raw)
		if err != nil {
			return query, service.BadRequest("status must be an integer")
		}
		query.Status = &value
	}

	if raw := strings.TrimSpace(c.Query("page")); raw != "" {
		value, err := strconv.Atoi(raw)
		if err != nil {
			return query, service.BadRequest("page must be an integer")
		}
		query.Page = value
	}

	if raw := strings.TrimSpace(c.Query("page_size")); raw != "" {
		value, err := strconv.Atoi(raw)
		if err != nil {
			return query, service.BadRequest("page_size must be an integer")
		}
		query.PageSize = value
	}

	return query, nil
}

func parsePositiveInt64(raw string, fieldName string) (int64, error) {
	value, err := strconv.ParseInt(strings.TrimSpace(raw), 10, 64)
	if err != nil || value <= 0 {
		return 0, service.BadRequest(fieldName + " must be a positive integer")
	}
	return value, nil
}
