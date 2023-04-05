package database

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

type Users struct {
	Id int `bun:",pk,autoincrement"`
	Name string
	Email string
}

func (u Users) String() string {
	return fmt.Sprintf("User<%d %s %v", u.Id, u.Name, u.Email)
}

func GetAllUsers(ctx context.Context, db *bun.DB) (interface{}, error) {
	var users []Users
	err := db.NewSelect().Model(&users).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}