# School Student Information System

A lightweight student information management backend built with **Go + Gin + SQLite**.

## Completed phases

### Phase 1: Base skeleton
- Layered project structure: `handler / service / repository`
- SQLite initialization and schema migration
- Unified JSON response format
- Health check endpoint
- Environment-based configuration

### Phase 2: Student profile management
- Create student
- Query student detail by ID
- Query student detail by student number
- Update student profile
- Soft delete student by setting `status = 0`

### Phase 3: Combined retrieval
- Exact lookup by student number
- Fuzzy search by student name
- Fuzzy search by class name
- Fuzzy search by major name
- Filter by grade year
- Filter by status
- Optional exact filter by class ID / major ID
- Pagination

---

## Tech stack
- Go
- Gin
- SQLite
- `database/sql`
- `modernc.org/sqlite`

---

## Project startup

### 1. Enter the project
```bash
cd school-student-system
```

### 2. Download dependencies
```bash
go mod tidy
```

### 3. Start the server
```bash
go run ./cmd/server
```

Default service address:
```text
http://127.0.0.1:8080
```

Default database path:
```text
data/student.db
```

---

## Optional environment variables

```bash
export STUDENT_SYS_HOST=127.0.0.1
export STUDENT_SYS_PORT=8080
export STUDENT_SYS_DB_PATH=./data/student.db
```

On Windows PowerShell:
```powershell
$env:STUDENT_SYS_HOST="127.0.0.1"
$env:STUDENT_SYS_PORT="8080"
$env:STUDENT_SYS_DB_PATH="./data/student.db"
```

---

## Seed dictionary data

The migration initializes a few majors and classes so that student CRUD can be tested immediately.

### Majors
| id | code | name |
|---:|---|---|
| 1 | CS | Computer Science |
| 2 | AGRI | Smart Agriculture |
| 3 | MGT | School Management |

### Classes
| id | code | name | grade year | major id |
|---:|---|---|---:|---:|
| 1 | CS-2024-1 | Computer Science Class 1 | 2024 | 1 |
| 2 | AGRI-2024-1 | Smart Agriculture Class 1 | 2024 | 2 |
| 3 | MGT-2025-1 | Management Class 1 | 2025 | 3 |

---

## API overview

### Health check
```http
GET /api/v1/health
```

### Create student
```http
POST /api/v1/students
Content-Type: application/json
```

Request body:
```json
{
  "student_no": "20260001",
  "name": "Zhang San",
  "class_id": 1,
  "phone": "13800000000",
  "email": "zhangsan@example.com",
  "address": "Dormitory A-101"
}
```

### Query by ID
```http
GET /api/v1/students/1
```

### Query by student number
```http
GET /api/v1/students/by-no/20260001
```

### Update student
```http
PUT /api/v1/students/1
Content-Type: application/json
```

Request body:
```json
{
  "name": "Zhang San Updated",
  "class_id": 2,
  "phone": "13900000000",
  "email": "zhangsan_updated@example.com",
  "address": "Dormitory B-202",
  "status": 1
}
```

### Soft delete student
```http
DELETE /api/v1/students/1
```

### Combined retrieval
```http
GET /api/v1/students?name=Zhang&class_name=Computer&major_name=Science&grade_year=2024&page=1&page_size=20
```

Supported query parameters:

| Parameter | Type | Description |
|---|---|---|
| `student_no` | string | Exact match |
| `name` | string | Fuzzy match |
| `class_name` | string | Fuzzy match |
| `major_name` | string | Fuzzy match |
| `grade_year` | int | Exact filter |
| `class_id` | int | Exact filter |
| `major_id` | int | Exact filter |
| `status` | int | `1=active`, `0=inactive` |
| `page` | int | Default: 1 |
| `page_size` | int | Default: 20, max: 100 |

When `status` is omitted, the API returns active students only.

---

## Unified response structure

### Success
```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

### Error
```json
{
  "code": 400,
  "message": "invalid student request"
}
```

---

## Run tests
```bash
go test ./...
```

The test suite covers:
- Create student
- Query by student number
- Combined search
- Update student
- Soft delete student
- Status-based retrieval after deletion

---

## Notes on data boundaries

This project intentionally separates:
- `majors`
- `classes`
- `students`

This makes the student module easier to search and keeps future extensions cleaner, such as:
- scores
- courses
- health records


---

## Frontend debugging page

A complete static frontend page is included under:

```text
frontend/index.html
```

It covers all backend features completed in the first three phases:
- Health check
- Student create
- Student update
- Student soft delete
- Student detail by ID
- Exact student number lookup
- Combined student retrieval with pagination
- Request/response preview panel for debugging

### Open directly on Windows

```text
scripts/open_frontend_windows.bat
```

### Open directly on macOS / Linux

```bash
./scripts/open_frontend_unix.sh
```

The page uses Bootstrap 5.3.x through a CDN and defaults to:

```text
http://127.0.0.1:8080/api/v1
```

You can change the API base address directly on the page.

### CORS support

The backend now enables development CORS responses for non-credentialed local frontend requests:
- `Access-Control-Allow-Origin: *`
- `GET, POST, PUT, DELETE, OPTIONS`
- JSON request headers
- Preflight `OPTIONS` handling

This allows the frontend page to call the Go backend when opened locally or served from another local port.
