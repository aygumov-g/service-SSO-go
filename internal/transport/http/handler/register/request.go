package register

type request struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
