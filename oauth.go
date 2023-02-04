package Forum

type OauthCfg struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}
type AuthToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type,omitempty"`
}

type Userinfo struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type GithubUserInfo struct {
	Username string `json:"login"`
	Email    string `json:"email"`
}
