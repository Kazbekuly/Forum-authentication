package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"Forum/pkg/service"

	"Forum"
)

const (
	ghAuthURL            = "https://github.com/login/oauth/authorize"
	ghTokenURL           = "https://github.com/login/oauth/access_token"
	ghClientIDSignUp     = "18f8079501b63ded82a4"
	ghClientSecretSignUp = "18c396444e9cf0334d89c15638b28899950bea91"
	ghClientIDSignIn     = "bb18a6b886610d4a60b0"
	ghClientSecretSignIn = "dfc1a5ae1aa7fcdcbc963a900be2e9f02c886176"
)

var (
	ghSignInCfg = &Forum.OauthCfg{
		ClientID:     ghClientIDSignIn,
		ClientSecret: ghClientSecretSignIn,
		RedirectURL:  "http://localhost:8080/signIn/github/callback",
		Scopes:       []string{"user:email"},
	}
	ghSignUpCfg = &Forum.OauthCfg{
		ClientID:     ghClientIDSignUp,
		ClientSecret: ghClientSecretSignUp,
		RedirectURL:  "http://localhost:8080/signUp/github/callback",
		Scopes:       []string{"user:email"},
	}
)

func requestToGithub(w http.ResponseWriter, r *http.Request, cfg Forum.OauthCfg) {
	URL, err := url.Parse(ghAuthURL)
	if err != nil {
		log.Printf("Parse: %s", err)
	}

	parameters := url.Values{}
	parameters.Add("client_id", cfg.ClientID)
	parameters.Add("redirect_uri", cfg.RedirectURL)
	parameters.Add("scope", strings.Join(cfg.Scopes, " "))

	URL.RawQuery = parameters.Encode()
	url := URL.String()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) githubSignIn(w http.ResponseWriter, r *http.Request) {
	requestToGithub(w, r, *ghSignInCfg)
}

func (h *Handler) githubSignUp(w http.ResponseWriter, r *http.Request) {
	requestToGithub(w, r, *ghSignUpCfg)
}

func (h *Handler) signInCallbackGithub(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	user, err := h.userFromGithubInfo(code, ghSignInCfg)
	if err != nil {
		if errors.Is(err, service.ErrorEmail) {
			h.HandleErrorPage(w, http.StatusBadRequest, service.ErrorEmail)
			return
		}
		h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
		return
	}
	token, err := h.services.Authorization.GenerateToken(user.Username, user.Password, true)
	if err != nil {
		if errors.Is(err, service.ErrorWrongPassword) {
			h.HandleErrorPage(w, http.StatusBadRequest, service.ErrorWrongPassword)
			return
		} else if errors.Is(err, service.ErrCheckInvalid) {
			h.HandleErrorPage(w, http.StatusBadRequest, service.ErrCheckInvalid)
			return
		}
		h.HandleErrorPage(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "session_token",
		Value: token.AuthToken,
		Path:  "/",
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) signUpCallbackGithub(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	user, err := h.userFromGithubInfo(code, ghSignUpCfg)
	if err != nil {
		if errors.Is(err, service.ErrorEmail) {
			h.HandleErrorPage(w, http.StatusBadRequest, service.ErrorEmail)
			return
		}
		h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
		return
	}
	if err := h.services.Authorization.CreateUser(*user); err != nil {
		if err != nil {
			if errors.Is(err, service.ErrCheckInvalid) {
				h.HandleErrorPage(w, http.StatusBadRequest, service.ErrCheckInvalid)
				return
			} else if errors.Is(err, sql.ErrNoRows) {
				h.HandleErrorPage(w, http.StatusUnauthorized, sql.ErrNoRows)
				return
			}
			h.HandleErrorPage(w, http.StatusBadRequest, errors.New(http.StatusText(http.StatusBadRequest)))
			return
		}
	}
	token, err := h.services.Authorization.GenerateToken(user.Username, user.Password, true)
	if err != nil {
		if errors.Is(err, service.ErrorWrongPassword) {
			h.HandleErrorPage(w, http.StatusBadRequest, service.ErrorWrongPassword)
			return
		}
		h.HandleErrorPage(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "session_token",
		Value: token.AuthToken,
		Path:  "/",
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) userFromGithubInfo(code string, cfg *Forum.OauthCfg) (*Forum.User, error) {
	accessToken, err := githubAccessToken(cfg, code)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil,
	)
	if err != nil {
		return nil, err
	}

	authHeaderValue := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", authHeaderValue)
	req.Header.Set("accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var users *Forum.GithubUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, err
	}

	email, err := emailFromGithub(code, accessToken, cfg)
	if err != nil {
		return nil, service.ErrorEmail
	}

	user := &Forum.User{
		Username: users.Username,
		Email:    email,
		Password: "ForumAuthentication01!01",
	}

	return user, nil
}

func emailFromGithub(code, accessToken string, cfg *Forum.OauthCfg) (string, error) {
	req, err := http.NewRequest(
		"GET",
		"https://api.github.com/user/emails",
		nil,
	)
	if err != nil {
		return "", err
	}
	authHeaderValue := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", authHeaderValue)
	req.Header.Set("accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data []*Forum.GithubUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	if data[0].Email == "" {
		return "", service.ErrorEmail
	}

	return data[0].Email, nil
}

func githubAccessToken(cfg *Forum.OauthCfg, code string) (string, error) {
	v := url.Values{
		"code":          {code},
		"client_id":     {cfg.ClientID},
		"client_secret": {cfg.ClientSecret},
	}

	req, err := http.NewRequest("POST", ghTokenURL, strings.NewReader(v.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var token Forum.AuthToken
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return "", err
	}

	return token.AccessToken, nil
}
