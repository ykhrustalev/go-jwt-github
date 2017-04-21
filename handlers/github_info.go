package handlers

import (
	"errors"
	"github.com/go-pg/pg"
	"github.com/ykhrustalev/exploregithub/dao"
	"github.com/ykhrustalev/exploregithub/githubapi"
	"github.com/ykhrustalev/exploregithub/jsonhttp"
	"github.com/ykhrustalev/exploregithub/jwtutils"
	"net/http"
	"strings"
)

func CreateGithubUserHandler(db *pg.DB, githubAuth *githubapi.Auth) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		claims := req.Context().Value("claims").(*jwtutils.Claims)

		githubToken, err := dao.GetGithubTokenForUser(db, claims.UserId)
		if err != nil {
			jsonhttp.ErrorResponse500(w, "can't get basic info", err)
			return
		}

		if githubToken == nil {
			jsonhttp.ErrorResponse(w, "can't request github", errors.New("missing token"), 401)
			return
		}

		client := githubAuth.ApiClient(&githubToken.Value)

		user, err := client.User()
		if err != nil {
			// TODO: move to client method to wrap that it is expired
			if strings.Contains(err.Error(), "token expired") {
				jsonhttp.ErrorResponse(w, "token expired", err, 401)
				return
			}
			jsonhttp.ErrorResponse500(w, "failed to get user info %s", err)
			return
		}

		jsonhttp.Response200(w, user)

	}
}
