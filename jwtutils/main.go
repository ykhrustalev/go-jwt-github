package jwtutils

import (
	"encoding/json"
	"errors"
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
