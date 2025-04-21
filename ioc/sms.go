package ioc

import (
	"fire/internal/service/sms"
	"fire/internal/service/sms/local"
)

func InitSMS() sms.Service {
	return local.NewService()
}
