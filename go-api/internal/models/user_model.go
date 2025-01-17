package models

type UserModel struct {
	BaseModel
	UserName    string `json:"userName"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Street      string `json:"street"`
	City        string `json:"city"`
	Zip         string `json:"zip"`
	Country     string `json:"country"`
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email"`
} //@name User

type SaveUserModel struct {
	UserName    string `json:"userName" validate:"required,uniqueUserName"`
	FirstName   string `json:"firstName" validate:"required"`
	LastName    string `json:"lastName"`
	Street      string `json:"street"`
	City        string `json:"city"`
	Zip         string `json:"zip"`
	Country     string `json:"country"`
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email"`
} //@name SaveUser
