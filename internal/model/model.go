package model

import "github.com/google/uuid"

// структура записи одного элемента данных
type Data struct {
	ID         uuid.UUID
	Data       int
	InProgress int
}
