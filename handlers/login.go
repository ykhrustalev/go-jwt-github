package handlers

import (
	"github.com/denkyl08/negroniJWT"
	"github.com/go-pg/pg"
	"github.com/ykhrustalev/exploregithub/dao"
	"github.com/ykhrustalev/exploregithub/githubapi"
	"github.com/ykhrustalev/exploregithub/jsonhttp"
	"github.com/ykhrustalev/exploregithub/jwtutils"
	"log"
	"net/http"
)

func CreateLoginHandler(db *pg.DB, githubAuth *githubapi.Auth) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// user with valid token should not get a new
		// TODO: ???
		claimsMap, ok := negroniJWT.Get(req)
		if ok {
			claims, err := jwtutils.NewClaimsFromMap(claimsMap)
			if err != nil {
				log.Println("user tries to login having invalid toke already")
			} else {
				err = dao.DeleteJwtToken(db, claims.TokenId)
				if err != nil {
					jsonhttp.ErrorResponse500(w, "failed to delete old token", err)
					return
				}
			}
		}

		http.Redirect(w, req, githubAuth.GetRedirectUrl(), http.StatusTemporaryRedirect)
	}
}
