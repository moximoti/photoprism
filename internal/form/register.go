package form

type Register struct {
	Email    string `json:"email"`
	FullName string `json:"fullname"`
	UserName string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
	IdToken  string `json:"idToken"`
}
