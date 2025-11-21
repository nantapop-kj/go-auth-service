package dto

import "mime/multipart"

type RegisterRequest struct {
	Email        string                `form:"email" validate:"required,email"`
	Password     string                `form:"password" validate:"required,min=8"`
	FirstName    string                `form:"first_name" validate:"required"`
	LastName     string                `form:"last_name" validate:"required"`
	ProfileImage *multipart.FileHeader `form:"profile_image"`
}
