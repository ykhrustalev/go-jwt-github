package githubapi

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"net/http"
)

type Auth struct {
	clientId     string
	clientSecret string
}

func NewAuth(clientId, clientSecret string) *Auth {
	return &Auth{
		clientId:     clientId,
		clientSecret: clientSecret,
	}
}

func (auth *Auth) GetRedirectUrl() string {
	return auth.GetOAuthConf().AuthCodeURL("state", oauth2.AccessTypeOffline)
}

func (auth *Auth) GetOAuthConf() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     auth.clientId,
		ClientSecret: auth.clientSecret,
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}
}

func (auth *Auth) ApiClient(token *oauth2.Token) *Api {
	ctx := context.Background()
	return &Api{auth.GetOAuthConf().Client(ctx, token)}
}

func (auth *Auth) ExchangeToken(code string) (*oauth2.Token, error) {
	ctx := context.Background()
	return auth.GetOAuthConf().Exchange(ctx, code)
}

type Api struct {
	client *http.Client
}

func (api *Api) NewGetRequest(path string) *http.Request {
	// no error here to expect
	req, _ := http.NewRequest("GET", fmt.Sprintf("https://api.github.com%s", path), nil)
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	return req
}

type User struct {
	Login             string `json:"login"`
	Id                int    `json:"id"`
	AvatarUrl         string `json:"avatar_url"`
	GravatarId        string `json:"gravatar_id"`
	Url               string `json:"url"`
	HtmlUrl           string `json:"html_url"`
	FollowersUrl      string `json:"followers_url"`
	FollowingUrl      string `json:"following_url"`
	GistsUrl          string `json:"gists_url"`
	StarredUrl        string `json:"starred_url"`
	SubscriptionsUrl  string `json:"subscriptions_url"`
	OrganizationsUrl  string `json:"organizations_url"`
	ReposUrl          string `json:"repos_url"`
	EventsUrl         string `json:"events_url"`
	ReceivedEventsUrl string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
	Name              string `json:"name"`
	Company           string `json:"company"`
	Blog              string `json:"blog"`
	Location          string `json:"location"`
	Email             string `json:"email"`
	Hireable          bool   `json:"hireable"`
	Bio               string `json:"bio"`
	PublicRepos       int    `json:"public_repos"`
	PublicGists       int    `json:"public_gists"`
	Followers         int    `json:"followers"`
	Following         int    `json:"following"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}

type UserEmail struct {
	Email    string `json:"email"`
	Verified bool   `json:"verified"`
	Primary  bool   `json:"primary"`
}

func (api *Api) doGetRequest(req *http.Request, obj interface{}) (*http.Response, error) {
	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(obj)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
func (api *Api) User() (*User, error) {
	var user User
	_, err := api.doGetRequest(api.NewGetRequest("/user"), &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (api *Api) UserEmails() ([]UserEmail, error) {
	var emails []UserEmail
	_, err := api.doGetRequest(api.NewGetRequest("/user/emails"), &emails)
	if err != nil {
		return nil, err
	}

	return emails, nil
}
