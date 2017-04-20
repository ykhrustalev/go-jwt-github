package dao

import (
	"github.com/go-pg/pg"
	"golang.org/x/oauth2"
)

type GithubToken struct {
	tableName struct{} `sql:"github_tokens"`

	Id     int64        `sql:"id"`
	UserId int64        `sql:"github_user_id"`
	Value  oauth2.Token `sql:"value"`
}

func SaveGithubToken(tx *pg.Tx, token *GithubToken) error {
	_, err := tx.Model(token).
		OnConflict("(id) DO UPDATE").
		Set("github_user_id = ?github_user_id, value = ?value").
		Insert()
	return err
}

func GetGithubTokenForUser(db *pg.DB, userId int64) (*GithubToken, error) {
	var token GithubToken
	err := db.Model(&token).
		Where("github_user_id = ?", userId).
		Limit(1).
		Select()

	if err == nil {
		return &token, nil
	} else if err == pg.ErrNoRows {
		return nil, nil
	} else {
		return nil, err
	}
}

func DeleteGithubTokensForUser(tx *pg.Tx, userId int64) error {
	_, err := tx.Model(&GithubToken{}).
		Where("github_user_id = ?", userId).
		Delete()
	return err
}
