package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Task struct {
	ID           primitive.ObjectID `bson:"_id"`
	Task_name    *string            `json:"task_name" validate:"required"`
	Is_important *bool              `json:"is_important" validate:"required"`
	Is_completed *bool              `json:"is_completed" validate:"required"`
	Time         *int64             `json:"created_at" validate:"required"`
	Task_id      string             `json:"task_id"`
}

type Tasks struct {
	UserId string `bson:"user_id"`
	Tasks  []Task `bson:"tasks" validate:"required"`
}
