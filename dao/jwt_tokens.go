package dao

import (
	"github.com/go-pg/pg"
)

type JwtToken struct {
	tableName struct{} `sql:"jwt_tokens"`

	Id     int64  `sql:"id"`
	UserId int64  `sql:"github_user_id"`
	Value  string `sql:"value"`
}

func SaveJwtToken(tx *pg.Tx, token *JwtToken) error {
	_, err := tx.Model(token).
		OnConflict("(id) DO UPDATE").
		Set("github_user_id = ?github_user_id, value = ?value").
		Insert()
	return err
}
func DeleteJwtToken(db *pg.DB, tokenId int64) error {
	_, err := db.Model(&JwtToken{}).Where("id = ?", tokenId).Delete()
	return err
}

