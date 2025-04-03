//go:build !k8s

package config

var Config = config{
	DB: DBConfig{
		DSN: "root@tcp(localhost:3306)/fire",
	},
	Redis: RedisConfig{
		Addr: "localhost:6379",
	},
}
