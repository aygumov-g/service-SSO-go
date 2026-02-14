package register

type registerRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
