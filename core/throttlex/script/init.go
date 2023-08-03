package script

import (
	"context"
	_ "embed"

	"github.com/redis/go-redis/v9"
)

//go:embed restrictions-allow.lua
var restrictions_allow_script string
var RestrictionsAllowScript *redis.Script

//go:embed limit_send_script.lua
var limit_send_script string
var LimitSendScript *redis.Script

func init() {
	// * restrictionsScript
	RestrictionsAllowScript = redis.NewScript(restrictions_allow_script)
	// * limitSendScript
	LimitSendScript = redis.NewScript(limit_send_script)
}

func GetLimitSendScript(ctx context.Context, rds *redis.Client) *redis.Script {
	result := redis.NewScript(limit_send_script).Load(ctx, rds)
	str, err := result.Result()
	if err != nil {
		panic(err)
	}
	_ = str
	return LimitSendScript
}
