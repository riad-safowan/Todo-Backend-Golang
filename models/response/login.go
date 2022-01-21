package response

import "time"

type LoginResponse struct {
	// ID            primitive.ObjectID `bson:"_id"`
	First_name    *string   `json:"first_name" validate:"required,min=2,max=100"`
	Last_name     *string   `json:"last_name" validate:"required"`
	Email         *string   `json:"email" validate:"email,required"`
	Phone_number  *string   `json:"phone_number" validate:"required"`
	Access_token  *string   `json:"access_token"`
	Refresh_token *string   `json:"refresh_token"`
	Created_at    time.Time `json:"created_at"`
	Updated_at    time.Time `json:"updated_at"`
	User_id       string    `json:"user_id"`
}