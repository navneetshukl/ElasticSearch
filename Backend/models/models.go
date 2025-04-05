package models

import "context"

type User struct {
	ID    string `json:"id" faker:"uuid_hyphenated"` // or generate manually
	Name  string `json:"name" faker:"name"`
	Email string `json:"email" faker:"email"`
	Phone string `json:"phone" faker:"phone_number"` // Use string for phone numbers
}


type ElasticSrvUseCase interface {
	InsertToElastic(ctx context.Context) error
	GetData(ctx context.Context, query string) ([]User, error)
}
