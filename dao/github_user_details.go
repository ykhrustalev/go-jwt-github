package dao

import (
	"github.com/go-pg/pg"
	"github.com/ykhrustalev/exploregithub/githubapi"
)

type GithubUserDetails struct {
	tableName struct{} `sql:"github_user_details"`

	Id         int64                 `sql:"id"`
	GithubId   int64                 `sql:"github_id"`
	UserInfo   githubapi.User        `sql:"info"`
	UserEmails []githubapi.UserEmail `sql:"emails"`
}

func SaveGithubUserDetails(tx *pg.Tx, userDetails *GithubUserDetails) error {
	_, err := tx.Model(userDetails).
		OnConflict("(github_id) DO UPDATE").
		Set("info = ?info, emails = ?emails").
		Insert()
	return err
}
