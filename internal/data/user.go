package data

// add type enum
type User struct {
	Id                int    `json:"id"`
	Email             string `json:"email"`
	PasswordHash      string `json:"-"`
	Type              string `json:"type"`
	ImageId           string `json:"imageId"`
	LastReadMessageId int    `json:"-"`
}

type RegisterUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Type     string `json:"type"`
	ImageId  string `json:"imageId"`
}

type RegisterUserResponse struct {
	Email   string `json:"email"`
	Type    string `json:"type"`
	ImageId string `json:"imageId"`
}
