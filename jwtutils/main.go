package jwtutils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/denkyl08/negroniJWT"
	"net/http"
)

type Claims struct {
	UserId  int64 `json:"userId"`
	TokenId int64
}

func (claims *Claims) ToMap() (map[string]interface{}, error) {
	data, err := json.Marshal(claims)
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
	m["data"] = string(data)
	return m, nil
}

func NewClaimsFromMap(raw map[string]interface{}) (*Claims, error) {
	data, ok := raw["data"].(string)
	if !ok {
		return nil, errors.New("missing data key in claims map")
	}

	var claims Claims
	err := json.Unmarshal([]byte(data), &claims)
	if err != nil {
		return nil, err
	}
	return &claims, nil
}

func GetClaims(req *http.Request) (*Claims, error) {
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
