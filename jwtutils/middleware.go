package jwtutils

import (
	"context"
	"errors"
	"fmt"
	"github.com/denkyl08/negroniJWT"
	"github.com/go-pg/pg"
	"github.com/urfave/negroni"
	"github.com/ykhrustalev/exploregithub/dao"
	"github.com/ykhrustalev/exploregithub/jsonhttp"
	"net/http"
)

func getClaims(req *http.Request) (*Claims, error) {
	claimsMap, ok := negroniJWT.Get(req)
	if !ok {
		return nil, errors.New("missing token in request")
	}

	claims, err := NewClaimsFromMap(claimsMap)
	if err != nil {
		return nil, fmt.Errorf("malformated token data, %s", err)
	}

	return claims, nil
}

func CreateAuthMiddleware(db *pg.DB) negroni.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		claims, err := getClaims(req)
		if err != nil {
			jsonhttp.ErrorResponse(w, "invalid token", err, 401)
			return
		}

		token, err := dao.GetJwtToken(db, claims.TokenId)
		if err != nil {
			jsonhttp.ErrorResponse500(w, "failed to get", err)
			return
		}

		if token == nil {
			jsonhttp.ErrorResponse(w, "no token", nil, 401)
			return
		}

		ctx := context.WithValue(req.Context(), "claims", claims)
		req = req.WithContext(ctx)
		next(w, req)
	}

}
