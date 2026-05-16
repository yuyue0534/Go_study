package dto

type CreateStudentRequest struct {
	StudentNo string `json:"student_no"`
	Name      string `json:"name"`
	ClassID   int64  `json:"class_id"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Address   string `json:"address"`
}

type UpdateStudentRequest struct {
	Name    string `json:"name"`
	ClassID int64  `json:"class_id"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Address string `json:"address"`
	Status  *int   `json:"status"`
}
