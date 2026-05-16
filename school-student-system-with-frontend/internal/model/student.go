package model

type Student struct {
	ID        int64  `json:"id"`
	StudentNo string `json:"student_no"`
	Name      string `json:"name"`
	ClassID   int64  `json:"class_id"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Address   string `json:"address"`
	Status    int    `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type StudentDetail struct {
	Student
	ClassCode string `json:"class_code"`
	ClassName string `json:"class_name"`
	GradeYear int    `json:"grade_year"`
	MajorID   int64  `json:"major_id"`
	MajorCode string `json:"major_code"`
	MajorName string `json:"major_name"`
}

type StudentListQuery struct {
	StudentNo string
	Name      string
	ClassName string
	MajorName string
	GradeYear int
	ClassID   int64
	MajorID   int64
	Status    *int
	Page      int
	PageSize  int
}

type StudentPage struct {
	List     []StudentDetail `json:"list"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
	Total    int64           `json:"total"`
}
