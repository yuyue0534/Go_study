package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"school-student-system/internal/dto"
	"school-student-system/internal/model"
	"school-student-system/internal/repository"
)

type StudentService struct {
	repo *repository.StudentRepository
}

func NewStudentService(repo *repository.StudentRepository) *StudentService {
	return &StudentService{repo: repo}
}

func (s *StudentService) Create(ctx context.Context, req dto.CreateStudentRequest) (*model.StudentDetail, error) {
	req = normalizeCreateRequest(req)
	if err := validateCreateRequest(req); err != nil {
		return nil, err
	}

	classExists, err := s.repo.ClassExists(ctx, req.ClassID)
	if err != nil {
		return nil, Internal("failed to validate class", err)
	}
	if !classExists {
		return nil, BadRequest("class_id does not exist")
	}

	studentNoExists, err := s.repo.StudentNoExists(ctx, req.StudentNo)
	if err != nil {
		return nil, Internal("failed to validate student number", err)
	}
	if studentNoExists {
		return nil, Conflict("student_no already exists")
	}

	id, err := s.repo.Create(ctx, model.Student{
		StudentNo: req.StudentNo,
		Name:      req.Name,
		ClassID:   req.ClassID,
		Phone:     req.Phone,
		Email:     req.Email,
		Address:   req.Address,
		Status:    1,
	})
	if err != nil {
		return nil, Internal("failed to create student", err)
	}

	student, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, Internal("failed to load created student", err)
	}
	return student, nil
}

func (s *StudentService) GetByID(ctx context.Context, id int64) (*model.StudentDetail, error) {
	if id <= 0 {
		return nil, BadRequest("student id must be positive")
	}

	student, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, NotFound("student not found")
		}
		return nil, Internal("failed to query student", err)
	}
	return student, nil
}

func (s *StudentService) GetByStudentNo(ctx context.Context, studentNo string) (*model.StudentDetail, error) {
	studentNo = strings.TrimSpace(studentNo)
	if studentNo == "" {
		return nil, BadRequest("student_no is required")
	}

	student, err := s.repo.GetByStudentNo(ctx, studentNo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, NotFound("student not found")
		}
		return nil, Internal("failed to query student", err)
	}
	return student, nil
}

func (s *StudentService) Update(ctx context.Context, id int64, req dto.UpdateStudentRequest) (*model.StudentDetail, error) {
	if id <= 0 {
		return nil, BadRequest("student id must be positive")
	}

	current, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, NotFound("student not found")
		}
		return nil, Internal("failed to query student before update", err)
	}

	req = normalizeUpdateRequest(req)
	if err := validateUpdateRequest(req); err != nil {
		return nil, err
	}

	classExists, err := s.repo.ClassExists(ctx, req.ClassID)
	if err != nil {
		return nil, Internal("failed to validate class", err)
	}
	if !classExists {
		return nil, BadRequest("class_id does not exist")
	}

	status := current.Status
	if req.Status != nil {
		status = *req.Status
	}

	updated, err := s.repo.Update(ctx, id, model.Student{
		Name:    req.Name,
		ClassID: req.ClassID,
		Phone:   req.Phone,
		Email:   req.Email,
		Address: req.Address,
		Status:  status,
	})
	if err != nil {
		return nil, Internal("failed to update student", err)
	}
	if !updated {
		return nil, NotFound("student not found")
	}

	student, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, Internal("failed to load updated student", err)
	}
	return student, nil
}

func (s *StudentService) SoftDelete(ctx context.Context, id int64) error {
	if id <= 0 {
		return BadRequest("student id must be positive")
	}

	deleted, err := s.repo.SoftDelete(ctx, id)
	if err != nil {
		return Internal("failed to delete student", err)
	}
	if !deleted {
		student, lookupErr := s.repo.GetByID(ctx, id)
		if lookupErr != nil {
			if errors.Is(lookupErr, sql.ErrNoRows) {
				return NotFound("student not found")
			}
			return Internal("failed to verify delete target", lookupErr)
		}
		if student.Status == 0 {
			return Conflict("student is already inactive")
		}
		return NotFound("student not found")
	}
	return nil
}

func (s *StudentService) List(ctx context.Context, query model.StudentListQuery) (*model.StudentPage, error) {
	query.StudentNo = strings.TrimSpace(query.StudentNo)
	query.Name = strings.TrimSpace(query.Name)
	query.ClassName = strings.TrimSpace(query.ClassName)
	query.MajorName = strings.TrimSpace(query.MajorName)

	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 20
	}
	if query.PageSize > 100 {
		query.PageSize = 100
	}
	if query.Status != nil && (*query.Status != 0 && *query.Status != 1) {
		return nil, BadRequest("status must be 0 or 1")
	}
	if query.GradeYear < 0 {
		return nil, BadRequest("grade_year must not be negative")
	}
	if query.ClassID < 0 {
		return nil, BadRequest("class_id must not be negative")
	}
	if query.MajorID < 0 {
		return nil, BadRequest("major_id must not be negative")
	}

	page, err := s.repo.List(ctx, query)
	if err != nil {
		return nil, Internal("failed to list students", err)
	}
	return page, nil
}

func normalizeCreateRequest(req dto.CreateStudentRequest) dto.CreateStudentRequest {
	req.StudentNo = strings.TrimSpace(req.StudentNo)
	req.Name = strings.TrimSpace(req.Name)
	req.Phone = strings.TrimSpace(req.Phone)
	req.Email = strings.TrimSpace(req.Email)
	req.Address = strings.TrimSpace(req.Address)
	return req
}

func normalizeUpdateRequest(req dto.UpdateStudentRequest) dto.UpdateStudentRequest {
	req.Name = strings.TrimSpace(req.Name)
	req.Phone = strings.TrimSpace(req.Phone)
	req.Email = strings.TrimSpace(req.Email)
	req.Address = strings.TrimSpace(req.Address)
	return req
}

func validateCreateRequest(req dto.CreateStudentRequest) error {
	if req.StudentNo == "" {
		return BadRequest("student_no is required")
	}
	if len([]rune(req.StudentNo)) > 32 {
		return BadRequest("student_no must not exceed 32 characters")
	}
	if req.Name == "" {
		return BadRequest("name is required")
	}
	if len([]rune(req.Name)) > 64 {
		return BadRequest("name must not exceed 64 characters")
	}
	if req.ClassID <= 0 {
		return BadRequest("class_id must be positive")
	}
	if len([]rune(req.Phone)) > 32 {
		return BadRequest("phone must not exceed 32 characters")
	}
	if len([]rune(req.Email)) > 128 {
		return BadRequest("email must not exceed 128 characters")
	}
	if len([]rune(req.Address)) > 255 {
		return BadRequest("address must not exceed 255 characters")
	}
	return nil
}

func validateUpdateRequest(req dto.UpdateStudentRequest) error {
	if req.Name == "" {
		return BadRequest("name is required")
	}
	if len([]rune(req.Name)) > 64 {
		return BadRequest("name must not exceed 64 characters")
	}
	if req.ClassID <= 0 {
		return BadRequest("class_id must be positive")
	}
	if len([]rune(req.Phone)) > 32 {
		return BadRequest("phone must not exceed 32 characters")
	}
	if len([]rune(req.Email)) > 128 {
		return BadRequest("email must not exceed 128 characters")
	}
	if len([]rune(req.Address)) > 255 {
		return BadRequest("address must not exceed 255 characters")
	}
	if req.Status != nil && (*req.Status != 0 && *req.Status != 1) {
		return BadRequest("status must be 0 or 1")
	}
	return nil
}
