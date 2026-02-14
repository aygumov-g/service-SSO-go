package login

type loginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
