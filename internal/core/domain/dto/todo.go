package dto

type Todo struct {
	Id          string `json:"id" gorm:"primaryKey;"`
	Topic       string `json:"topic"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type TodoInputSave struct {
	Topic       string `json:"topic" validate:"required"`
	Description string `json:"description" validate:"required"`
	Status      string `json:"status" validate:"required"`
}

type TodoInputUpdateStatus struct {
	Id     string `json:"id" validate:"required"`
	Status string `json:"status" validate:"required"`
}
type TodoInputDelete struct {
	Id string `json:"id" validate:"required"`
}
