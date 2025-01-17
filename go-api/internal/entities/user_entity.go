package entities

type User struct {
	EntityBase
	UserName    string `dynamodbav:"UserName"`
	FirstName   string `dynamodbav:"FirstName"`
	LastName    string `dynamodbav:"LastName"`
	Street      string `dynamodbav:"Street"`
	City        string `dynamodbav:"City"`
	Zip         string `dynamodbav:"Zip"`
	Country     string `dynamodbav:"Country"`
	PhoneNumber string `dynamodbav:"PhoneNumber"`
	Email       string `dynamodbav:"Email"`
}
