package model

import (
	"context"
	"fmt"

	"github.com/portafolioLP/db"
)

var Ctx = context.Background()

func RedisRefreshToken(user User, token string) {
	rdb := db.Rdb()

	err := rdb.Set(Ctx, user.Username, token, 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(Ctx, user.Username).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(user.Username, val)
}

func ExistUserRedisToken(username string, token string) bool {
	rdb := db.Rdb()
	val, err := rdb.Get(Ctx, username).Result()
	if err != nil || token != val {
		return true
	} else {
		return false
	}
}
