//go:build !k8s

package config

var Config = config{
	DB: DBConfig{
		DSN: "host=localhost user=leochen dbname=postgres search_path=fire port=5432 sslmode=disable TimeZone=Asia/Taipei",
	},
	Redis: RedisConfig{
		Addr: "localhost:6379",
	},
}
