package sqlite

import (
	"database/sql"

	"github.com/AbhishekSinghDev/student-management/internal/config"
	"github.com/AbhishekSinghDev/student-management/internal/types"
	_ "modernc.org/sqlite"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	// force connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	email TEXT,
	age INTEGER
	)`)

	if err != nil {
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, nil
}

func (s Sqlite) CreateStudent(name string, email string, age int) (int64, error) {
	statement, err := s.Db.Prepare(`INSERT INTO students (name, email, age) VALUES (?, ?, ?)`)
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	result, err := statement.Exec(name, email, age)
	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastId, nil
}

func (s Sqlite) GetStudentById(id int64) (types.Student, error) {
	statement, err := s.Db.Prepare((`SELECT * FROM students WHERE id = ? LIMIT 1`))
	if err != nil {
		return types.Student{}, err
	}
	defer statement.Close()

	var student types.Student


	err = statement.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		return types.Student{}, err
	}


	return student, nil
}
