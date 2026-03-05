package models

import "time"

type User struct {
	Id             string     `json:"id"         bson:"_id"             example:"some uuid"`
	Name           string     `json:"name"        bson:"name"            example:"John Doe"`
	Role           string     `json:"user_role"   bson:"role"            example:"SYSTEM_ADMIN,ORGANIZATION_ADMIN,INSTRUCTOR,STUDENT"`
	Phone          string     `json:"phone"       bson:"phone"`
	Email          string     `json:"email"       bson:"email"`
	Password       string     `json:"-"           bson:"password"`
	Active         int        `json:"active"      bson:"active"`
	ExpiresAt      *time.Time `json:"expires_at"  bson:"expires_at"`
	CreatedAt      *time.Time `json:"created_at"  bson:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"  bson:"updated_at"`
	DeletedAt      *time.Time `json:"-"           bson:"deleted_at"`
}

type UserCreateModel struct {
	Name           string `json:"name" example:"John Doe"`
	Role           string `json:"user_role"`
	Phone          string `json:"phone"`
	Email          string `json:"email"`
}

type GetUserListModel struct {
	Count int     `json:"count"`
	Users []*User `json:"users"`
}

type UpdateUserProfileModel struct {
	Id             string `json:"id" example:"some uuid"`
	Name           string `json:"name" example:"John Doe"`
	Role           string `json:"user_role"`
	Phone          string `json:"phone"`
	Email          string `json:"email"`
}

