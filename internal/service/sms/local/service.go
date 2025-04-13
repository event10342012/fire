package local

import (
	"context"
	"log"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (svc *Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	log.Println("code:", tplId, args[0], numbers)
	return nil
}
