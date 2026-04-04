package token

type Input struct {
	Code         string
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

type Output struct {
	AccessToken  string
	RefreshToken string
}
