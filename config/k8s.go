//go:build k8s

package config

var Config = config{
	DB: DBConfig{
		DSN: "root:root@tcp(fire-mysql:3306)/fire",
	},
	Redis: RedisConfig{
		Addr: "fire-redis:6379",
	},
}
