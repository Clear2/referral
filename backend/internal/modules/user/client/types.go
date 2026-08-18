package client

type UserCreateOneDto struct {
	Name  string `json:"name" binding:"required,min=1,max=100"`
	Email string `json:"email" binding:"required,email,max=254"`
}

type UserGetByIdDto struct {
	ID int `uri:"id" binding:"required,min=1"`
}
