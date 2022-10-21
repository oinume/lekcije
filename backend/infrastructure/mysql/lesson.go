package mysql

import (
	"database/sql"

	"github.com/oinume/lekcije/backend/domain/repository"
)

type lessonRepository struct {
	db *sql.DB
}

func NewLessonRepository(db *sql.DB) repository.Lesson {
	return &lessonRepository{db: db}
}
