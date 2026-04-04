package change_password

type request struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
