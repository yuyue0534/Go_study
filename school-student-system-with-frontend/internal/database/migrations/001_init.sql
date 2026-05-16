PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS majors (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    major_code TEXT NOT NULL UNIQUE,
    major_name TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS classes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    class_code TEXT NOT NULL UNIQUE,
    class_name TEXT NOT NULL,
    grade_year INTEGER NOT NULL,
    major_id INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (major_id) REFERENCES majors(id)
);

CREATE TABLE IF NOT EXISTS students (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    student_no TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    class_id INTEGER NOT NULL,
    phone TEXT NOT NULL DEFAULT '',
    email TEXT NOT NULL DEFAULT '',
    address TEXT NOT NULL DEFAULT '',
    status INTEGER NOT NULL DEFAULT 1 CHECK (status IN (0, 1)),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (class_id) REFERENCES classes(id)
);

CREATE INDEX IF NOT EXISTS idx_classes_major_id ON classes(major_id);
CREATE INDEX IF NOT EXISTS idx_classes_grade_year ON classes(grade_year);
CREATE INDEX IF NOT EXISTS idx_students_class_id ON students(class_id);
CREATE INDEX IF NOT EXISTS idx_students_status ON students(status);

INSERT OR IGNORE INTO majors (id, major_code, major_name) VALUES
    (1, 'CS', 'Computer Science'),
    (2, 'AGRI', 'Smart Agriculture'),
    (3, 'MGT', 'School Management');

INSERT OR IGNORE INTO classes (id, class_code, class_name, grade_year, major_id) VALUES
    (1, 'CS-2024-1', 'Computer Science Class 1', 2024, 1),
    (2, 'AGRI-2024-1', 'Smart Agriculture Class 1', 2024, 2),
    (3, 'MGT-2025-1', 'Management Class 1', 2025, 3);
