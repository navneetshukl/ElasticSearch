package models

type User struct {
	ID    string `json:"id" faker:"id"`
	Name  string `json:"name" faker:"name"`
	Email string `json:"email" faker:"email"`
	Phone int64  `json:"phone" faker:"phone"`
}

type ElasticSrvUseCase interface{}
