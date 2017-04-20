package handlers

import (
	"github.com/go-pg/pg"
	"github.com/ykhrustalev/exploregithub/dao"
	"github.com/ykhrustalev/exploregithub/jsonhttp"
	"github.com/ykhrustalev/exploregithub/jwtutils"
	"net/http"
)

func CreateLogoutHandler(db *pg.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		claims, err := jwtutils.GetClaims(req)
		if err != nil {
			jsonhttp.ErrorResponse(w, "invalid token", err, 401)
			return
		}

		err = dao.DeleteJwtToken(db, claims.TokenId)
		if err != nil {
			jsonhttp.ErrorResponse500(w, "failed to save JwtToken", err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
