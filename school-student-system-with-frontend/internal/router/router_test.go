package router_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"school-student-system/internal/database"
	"school-student-system/internal/handler"
	"school-student-system/internal/repository"
	"school-student-system/internal/router"
	"school-student-system/internal/service"
)

type apiResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type studentDetail struct {
	ID        int64  `json:"id"`
	StudentNo string `json:"student_no"`
	Name      string `json:"name"`
	Status    int    `json:"status"`
}

type studentPage struct {
	Total int64 `json:"total"`
}

func TestStudentLifecycleAndSearch(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine, db := setupTestRouter(t)
	defer db.Close()

	createPayload := map[string]any{
		"student_no": "20260001",
		"name":       "Zhang San",
		"class_id":   1,
		"phone":      "13800000000",
		"email":      "zhangsan@example.com",
		"address":    "Dormitory A-101",
	}
	createResp := performJSONRequest(t, engine, http.MethodPost, "/api/v1/students", createPayload, http.StatusCreated)

	var created studentDetail
	mustUnmarshalData(t, createResp, &created)
	if created.ID <= 0 {
		t.Fatalf("expected created ID > 0, got %d", created.ID)
	}
	if created.StudentNo != "20260001" {
		t.Fatalf("unexpected student no: %s", created.StudentNo)
	}

	byNoResp := performJSONRequest(t, engine, http.MethodGet, "/api/v1/students/by-no/20260001", nil, http.StatusOK)
	var byNo studentDetail
	mustUnmarshalData(t, byNoResp, &byNo)
	if byNo.ID != created.ID {
		t.Fatalf("expected same student ID, got %d and %d", created.ID, byNo.ID)
	}

	searchResp := performJSONRequest(t, engine, http.MethodGet, "/api/v1/students?name=Zhang&class_name=Computer&major_name=Science&grade_year=2024&page=1&page_size=10", nil, http.StatusOK)
	var page studentPage
	mustUnmarshalData(t, searchResp, &page)
	if page.Total != 1 {
		t.Fatalf("expected search total 1, got %d", page.Total)
	}

	updatePayload := map[string]any{
		"name":     "Zhang San Updated",
		"class_id": 2,
		"phone":    "13900000000",
		"email":    "zhangsan_updated@example.com",
		"address":  "Dormitory B-202",
		"status":   1,
	}
	updateResp := performJSONRequest(t, engine, http.MethodPut, "/api/v1/students/1", updatePayload, http.StatusOK)
	var updated studentDetail
	mustUnmarshalData(t, updateResp, &updated)
	if updated.Name != "Zhang San Updated" {
		t.Fatalf("expected updated name, got %s", updated.Name)
	}

	performJSONRequest(t, engine, http.MethodDelete, "/api/v1/students/1", nil, http.StatusOK)

	activeResp := performJSONRequest(t, engine, http.MethodGet, "/api/v1/students", nil, http.StatusOK)
	var activePage studentPage
	mustUnmarshalData(t, activeResp, &activePage)
	if activePage.Total != 0 {
		t.Fatalf("expected active total 0 after soft delete, got %d", activePage.Total)
	}

	inactiveResp := performJSONRequest(t, engine, http.MethodGet, "/api/v1/students?status=0", nil, http.StatusOK)
	var inactivePage studentPage
	mustUnmarshalData(t, inactiveResp, &inactivePage)
	if inactivePage.Total != 1 {
		t.Fatalf("expected inactive total 1 after soft delete, got %d", inactivePage.Total)
	}
}

func TestCORSPreflight(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine, db := setupTestRouter(t)
	defer db.Close()

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/students", nil)
	req.Header.Set("Origin", "http://127.0.0.1:5500")
	req.Header.Set("Access-Control-Request-Method", http.MethodPost)
	req.Header.Set("Access-Control-Request-Headers", "Content-Type")

	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected HTTP %d, got %d", http.StatusNoContent, recorder.Code)
	}
	if got := recorder.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Fatalf("expected wildcard CORS origin, got %q", got)
	}
	if got := recorder.Header().Get("Access-Control-Allow-Methods"); got == "" {
		t.Fatal("expected Access-Control-Allow-Methods header")
	}
	if got := recorder.Header().Get("Access-Control-Allow-Headers"); got == "" {
		t.Fatal("expected Access-Control-Allow-Headers header")
	}
}

func setupTestRouter(t *testing.T) (http.Handler, *sql.DB) {
	t.Helper()

	db, err := database.Open(":memory:")
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	if err := database.Migrate(db); err != nil {
		db.Close()
		t.Fatalf("migrate test database: %v", err)
	}

	repo := repository.NewStudentRepository(db)
	svc := service.NewStudentService(repo)
	h := handler.NewStudentHandler(svc)
	return router.New(h), db
}

func performJSONRequest(t *testing.T, engine http.Handler, method string, path string, payload any, expectedStatus int) apiResponse {
	t.Helper()

	var body bytes.Buffer
	if payload != nil {
		if err := json.NewEncoder(&body).Encode(payload); err != nil {
			t.Fatalf("encode payload: %v", err)
		}
	}

	req := httptest.NewRequest(method, path, &body)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, req)

	if recorder.Code != expectedStatus {
		t.Fatalf("expected HTTP %d, got %d, body=%s", expectedStatus, recorder.Code, recorder.Body.String())
	}

	var resp apiResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v, body=%s", err, recorder.Body.String())
	}
	if resp.Code != 0 {
		t.Fatalf("expected response code 0, got %d (%s)", resp.Code, resp.Message)
	}
	return resp
}

func mustUnmarshalData(t *testing.T, resp apiResponse, target any) {
	t.Helper()
	if err := json.Unmarshal(resp.Data, target); err != nil {
		t.Fatalf("decode response data: %v, raw=%s", err, string(resp.Data))
	}
}
