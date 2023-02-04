package handlers

import (
	"Forum"
	"Forum/pkg/service"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	authURL      = "https://accounts.google.com/o/oauth2/auth"
	tokenURL     = "https://oauth2.googleapis.com/token"
	clientID     = "1095628612590-6f670bssurhsgdv33glilg582skv530a.apps.googleusercontent.com"
	clientSecret = "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"
	signUpURI    = "http://localhost:8080/signUp/callback"
	signInURI    = "http://localhost:8080/signIn/callback"
)

var googleConfig = &Forum.OauthCfg{
	ClientID:     clientID,
	ClientSecret: clientSecret,
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	},
}

func (h *Handler) googleRequest(w http.ResponseWriter, r *http.Request, redirectURL string) {
	URL, err := url.Parse(authURL)
	if err != nil {
		log.Printf("Parse: %s", err)
	}

	parameters := url.Values{}
	parameters.Add("client_id", googleConfig.ClientID)
	parameters.Add("redirect_uri", redirectURL)
	parameters.Add("scope", strings.Join(googleConfig.Scopes, " "))
	parameters.Add("response_type", "code")
	URL.RawQuery = parameters.Encode()
	url := URL.String()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) googleSignUp(w http.ResponseWriter, r *http.Request) {
	h.googleRequest(w, r, signUpURI)
}

func (h *Handler) googleSignIn(w http.ResponseWriter, r *http.Request) {
	h.googleRequest(w, r, signInURI)
}

func (h *Handler) callbackGoogleSignIn(w http.ResponseWriter, r *http.Request) {
	user, err := getUserInfoFromGoogle(r, googleConfig, signInURI)
	if err != nil {
		h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
		return
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

func (h *Handler) callbackGoogleSignUp(w http.ResponseWriter, r *http.Request) {
	user, err := getUserInfoFromGoogle(r, googleConfig, signUpURI)
	if err != nil {
		h.HandleErrorPage(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
		return
	}
	err = h.services.Authorization.CreateUser(*user)
	if err != nil {
		if errors.Is(err, service.ErrCheckInvalid) {
			h.HandleErrorPage(w, http.StatusBadRequest, service.ErrCheckInvalid)
			return
		} else if errors.Is(err, sql.ErrNoRows) {
			h.HandleErrorPage(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
			return
		}
		h.HandleErrorPage(w, http.StatusBadRequest, errors.New(http.StatusText(http.StatusBadRequest)))
		return
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

func getUserInfoFromGoogle(r *http.Request, cfg *Forum.OauthCfg, redirectUrl string) (*Forum.User, error) {
	code := r.FormValue("code")
	token, err := getGoogleAccessToken(cfg, code, redirectUrl)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo?access_token="+url.QueryEscape(token), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userinfo *Forum.Userinfo
	if err := json.NewDecoder(resp.Body).Decode(&userinfo); err != nil {
		return nil, err
	}
	if len(strings.TrimSpace(userinfo.Email)) == 0 || len(strings.TrimSpace(userinfo.Name)) == 0 {
		return nil, errors.New("something went wrong with get user information")
	}
	user := &Forum.User{
		Username: strings.Join(strings.Split(userinfo.Name, " "), ""),
		Email:    userinfo.Email,
		Password: "ForumAuthentication!1",
	}
	return user, nil
}

func getGoogleAccessToken(cfg *Forum.OauthCfg, code, redirectUrl string) (string, error) {
	v := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {redirectUrl},
		"client_id":     {cfg.ClientID},
		"client_secret": {cfg.ClientSecret},
	}
	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(v.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var token Forum.AuthToken
	if err := json.NewDecoder(res.Body).Decode(&token); err != nil {
		return "", err
	}

	return token.AccessToken, nil
}
