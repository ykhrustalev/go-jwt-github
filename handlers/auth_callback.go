package handlers

import (
	"errors"
	"fmt"
	"github.com/denkyl08/negroniJWT"
	"github.com/go-pg/pg"
	"github.com/ykhrustalev/exploregithub/dao"
	"github.com/ykhrustalev/exploregithub/githubapi"
	"github.com/ykhrustalev/exploregithub/jsonhttp"
	"github.com/ykhrustalev/exploregithub/jwtutils"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"time"
	"sync/atomic"
	"sync"
)

func requestBasicData(code string, githubAuth *githubapi.Auth) (*oauth2.Token, *githubapi.User, []githubapi.UserEmail, error) {
	oauthToken, err := githubAuth.ExchangeToken(code)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to exchange tokens %s", err)
	}
	client := githubAuth.ApiClient(oauthToken)

	user, err := client.User()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get user info %s", err)
	}

	emails, err := client.UserEmails()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get user emails %s", err)
	}

	return oauthToken, user, emails, nil
}

func saveOauthToken(tx *pg.Tx, oauthToken *oauth2.Token, userId int64) (*dao.GithubToken, error) {
	err := dao.DeleteGithubTokensForUser(tx, userId)
	if err != nil {
		log.Println("failed to wipe old oauth tokens", err)
		return nil, err
	}

	userToken := &dao.GithubToken{
		UserId: userId,
		Value:  *oauthToken,
	}
	err = dao.SaveGithubToken(tx, userToken)
	if err != nil {
		log.Println("failed to save GithubToken", err)
		return nil, err
	}
	return userToken, nil
}

func generateJwtToken(tx *pg.Tx, userId int64) (*dao.JwtToken, error) {
	// create empty jwt oauthToken record
	record := &dao.JwtToken{
		UserId: userId,
		Value:  "empty",
	}
	err := dao.SaveJwtToken(tx, record)
	if err != nil {
		log.Println("failed to save JwtToken", err)
		return nil, err
	}

	claims := jwtutils.Claims{
		UserId:  userId,
		TokenId: record.Id,
	}
	claimsMap, err := claims.ToMap()
	if err != nil {
		log.Println("failed to generate claims map", err)
		return nil, err
	}

	// generate JWT , TODO: define time
	jwtToken, err := negroniJWT.GenerateToken(claimsMap, time.Now().Add(3*24*time.Hour))
	if err != nil {
		log.Println("failed to generate jwt", err)
		return nil, err
	}

	record.Value = jwtToken
	err = dao.SaveJwtToken(tx, record)
	if err != nil {
		log.Println("failed to save JwtToken 2nd time", err)
		return nil, err
	}

	return record, nil
}

type SessionStore struct {
	mw sync.RWMutex
	Sessions map[string]string
}

func NewSessionStore() *SessionStore {
	return &SessionStore{Sessions: make(map[string]string)}
}

func(store *SessionStore) Save(session string) {
	store.mw.Lock()
	defer store.mw.Unlock()

	store.Sessions[session] = ""

	go func() {
		time.Sleep(60 * time.Second)
		store.Delete(session)
	}()
}
func(store *SessionStore) Delete(session string) {
	store.mw.Lock()
	defer store.mw.Unlock()

	delete(store.Sessions, session)
}

func CreateAuthCallbackHandler(db *pg.DB, githubAuth *githubapi.Auth) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		code := req.FormValue("code")
		if code == "" {
			jsonhttp.ErrorResponse(w, "can't obtain code", errors.New("no code"), http.StatusBadRequest)
			return
		}

		oauthToken, user, emails, err := requestBasicData(code, githubAuth)
		if err != nil {
			jsonhttp.ErrorResponse500(w, "can't get basic info", err)
			return
		}

		var jwtToken *dao.JwtToken

		err = db.RunInTransaction(func(tx *pg.Tx) error {
			userDetails := &dao.GithubUserDetails{
				GithubId:   int64(user.Id),
				UserInfo:   *user,
				UserEmails: emails,
			}

			err = dao.SaveGithubUserDetails(tx, userDetails)
			if err != nil {
				log.Println("failed to save GithubUserDetails", err)
				return err
			}

			saveOauthToken(tx, oauthToken, userDetails.Id)

			jwtToken, err = generateJwtToken(tx, userDetails.Id)
			if err != nil {
				log.Println("failed to generate JWT record", err)
				return err
			}

			return nil
		})

		if err != nil {
			jsonhttp.ErrorResponse500(w, "aborting auth handler", err)
			return
		}


		w.WriteHeader(http.StatusOK)
		//lp := filepath.Join("templates", "auth_callback_success.html")
		//
		//tmpl := template.Must(template.ParseFiles(lp))
		//
		//tmpl.ExecuteTemplate(w, "layout", nil)

		body := `
		<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Login success</title>
</head>
<body>

<script>
window.opener.HandlePopupResult('123');
    //window.opener.document.getElementById("token").innerHTML = '123';
</script>

</body>
</html>

		`

		w.Write([]byte(body))

		//jsonhttp.Response200(w, struct{ Token string }{jwtToken.Value})
	}
}
