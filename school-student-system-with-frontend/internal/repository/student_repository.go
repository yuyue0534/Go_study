package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"school-student-system/internal/model"
)

type StudentRepository struct {
	db *sql.DB
}

func NewStudentRepository(db *sql.DB) *StudentRepository {
	return &StudentRepository{db: db}
}

func (r *StudentRepository) ClassExists(ctx context.Context, classID int64) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(1) FROM classes WHERE id = ?`, classID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("query class existence: %w", err)
	}
	return count > 0, nil
}

func (r *StudentRepository) StudentNoExists(ctx context.Context, studentNo string) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(1) FROM students WHERE student_no = ?`, studentNo).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("query student number existence: %w", err)
	}
	return count > 0, nil
}

func (r *StudentRepository) Create(ctx context.Context, student model.Student) (int64, error) {
	result, err := r.db.ExecContext(ctx, `
        INSERT INTO students (student_no, name, class_id, phone, email, address, status)
        VALUES (?, ?, ?, ?, ?, ?, ?)
    `, student.StudentNo, student.Name, student.ClassID, student.Phone, student.Email, student.Address, student.Status)
	if err != nil {
		return 0, fmt.Errorf("insert student: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("read inserted student id: %w", err)
	}
	return id, nil
}

func (r *StudentRepository) GetByID(ctx context.Context, id int64) (*model.StudentDetail, error) {
	return r.queryOne(ctx, `
        SELECT
            s.id,
            s.student_no,
            s.name,
            s.class_id,
            s.phone,
            s.email,
            s.address,
            s.status,
            s.created_at,
            s.updated_at,
            c.class_code,
            c.class_name,
            c.grade_year,
            m.id,
            m.major_code,
            m.major_name
        FROM students s
        JOIN classes c ON c.id = s.class_id
        JOIN majors m ON m.id = c.major_id
        WHERE s.id = ?
    `, id)
}

func (r *StudentRepository) GetByStudentNo(ctx context.Context, studentNo string) (*model.StudentDetail, error) {
	return r.queryOne(ctx, `
        SELECT
            s.id,
            s.student_no,
            s.name,
            s.class_id,
            s.phone,
            s.email,
            s.address,
            s.status,
            s.created_at,
            s.updated_at,
            c.class_code,
            c.class_name,
            c.grade_year,
            m.id,
            m.major_code,
            m.major_name
        FROM students s
        JOIN classes c ON c.id = s.class_id
        JOIN majors m ON m.id = c.major_id
        WHERE s.student_no = ?
    `, studentNo)
}

func (r *StudentRepository) Update(ctx context.Context, id int64, student model.Student) (bool, error) {
	result, err := r.db.ExecContext(ctx, `
        UPDATE students
        SET
            name = ?,
            class_id = ?,
            phone = ?,
            email = ?,
            address = ?,
            status = ?,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = ?
    `, student.Name, student.ClassID, student.Phone, student.Email, student.Address, student.Status, id)
	if err != nil {
		return false, fmt.Errorf("update student: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("read updated row count: %w", err)
	}
	return affected > 0, nil
}

func (r *StudentRepository) SoftDelete(ctx context.Context, id int64) (bool, error) {
	result, err := r.db.ExecContext(ctx, `
        UPDATE students
        SET status = 0, updated_at = CURRENT_TIMESTAMP
        WHERE id = ? AND status <> 0
    `, id)
	if err != nil {
		return false, fmt.Errorf("soft delete student: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("read deleted row count: %w", err)
	}
	return affected > 0, nil
}

func (r *StudentRepository) List(ctx context.Context, query model.StudentListQuery) (*model.StudentPage, error) {
	conditions, args := buildStudentConditions(query)
	whereClause := "WHERE " + strings.Join(conditions, " AND ")

	var total int64
	countSQL := `
        SELECT COUNT(1)
        FROM students s
        JOIN classes c ON c.id = s.class_id
        JOIN majors m ON m.id = c.major_id
    ` + whereClause

	if err := r.db.QueryRowContext(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, fmt.Errorf("count students: %w", err)
	}

	offset := (query.Page - 1) * query.PageSize
	listArgs := append(append([]any{}, args...), query.PageSize, offset)

	listSQL := `
        SELECT
            s.id,
            s.student_no,
            s.name,
            s.class_id,
            s.phone,
            s.email,
            s.address,
            s.status,
            s.created_at,
            s.updated_at,
            c.class_code,
            c.class_name,
            c.grade_year,
            m.id,
            m.major_code,
            m.major_name
        FROM students s
        JOIN classes c ON c.id = s.class_id
        JOIN majors m ON m.id = c.major_id
    ` + whereClause + `
        ORDER BY s.id DESC
        LIMIT ? OFFSET ?
    `

	rows, err := r.db.QueryContext(ctx, listSQL, listArgs...)
	if err != nil {
		return nil, fmt.Errorf("list students: %w", err)
	}
	defer rows.Close()

	students := make([]model.StudentDetail, 0)
	for rows.Next() {
		student, err := scanStudentDetail(rows)
		if err != nil {
			return nil, err
		}
		students = append(students, *student)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate student rows: %w", err)
	}

	return &model.StudentPage{
		List:     students,
		Page:     query.Page,
		PageSize: query.PageSize,
		Total:    total,
	}, nil
}

func buildStudentConditions(query model.StudentListQuery) ([]string, []any) {
	conditions := make([]string, 0, 8)
	args := make([]any, 0, 8)

	if query.StudentNo != "" {
		conditions = append(conditions, "s.student_no = ?")
		args = append(args, query.StudentNo)
	}
	if query.Name != "" {
		conditions = append(conditions, "s.name LIKE ?")
		args = append(args, "%"+query.Name+"%")
	}
	if query.ClassName != "" {
		conditions = append(conditions, "c.class_name LIKE ?")
		args = append(args, "%"+query.ClassName+"%")
	}
	if query.MajorName != "" {
		conditions = append(conditions, "m.major_name LIKE ?")
		args = append(args, "%"+query.MajorName+"%")
	}
	if query.GradeYear > 0 {
		conditions = append(conditions, "c.grade_year = ?")
		args = append(args, query.GradeYear)
	}
	if query.ClassID > 0 {
		conditions = append(conditions, "c.id = ?")
		args = append(args, query.ClassID)
	}
	if query.MajorID > 0 {
		conditions = append(conditions, "m.id = ?")
		args = append(args, query.MajorID)
	}

	if query.Status != nil {
		conditions = append(conditions, "s.status = ?")
		args = append(args, *query.Status)
	} else {
		conditions = append(conditions, "s.status = 1")
	}

	if len(conditions) == 0 {
		conditions = append(conditions, "1 = 1")
	}
	return conditions, args
}

func (r *StudentRepository) queryOne(ctx context.Context, query string, arg any) (*model.StudentDetail, error) {
	row := r.db.QueryRowContext(ctx, query, arg)
	student, err := scanStudentDetail(row)
	if err != nil {
		return nil, err
	}
	return student, nil
}

type studentScanner interface {
	Scan(dest ...any) error
}

func scanStudentDetail(scanner studentScanner) (*model.StudentDetail, error) {
	var student model.StudentDetail
	err := scanner.Scan(
		&student.ID,
		&student.StudentNo,
		&student.Name,
		&student.ClassID,
		&student.Phone,
		&student.Email,
		&student.Address,
		&student.Status,
		&student.CreatedAt,
		&student.UpdatedAt,
		&student.ClassCode,
		&student.ClassName,
		&student.GradeYear,
		&student.MajorID,
		&student.MajorCode,
		&student.MajorName,
	)
	if err != nil {
		return nil, err
	}
	return &student, nil
}
