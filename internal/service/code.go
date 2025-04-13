package service

import (
	"context"
	"errors"
	"fire/internal/repository"
	"fire/internal/service/sms"
	"fmt"
	"math/rand"
)

var (
	ErrCodeSendTooMany = repository.ErrCodeVerifyTooMany
)

type CodeService struct {
	repo repository.CodeRepository
	sms  sms.Service
}

func NewCodeService(repo repository.CodeRepository, smsSvc sms.Service) *CodeService {
	return &CodeService{
		repo: repo,
		sms:  smsSvc,
	}
}

func (svc *CodeService) Send(ctx context.Context, biz, phone string) error {
	code := svc.generate()
	err := svc.repo.Set(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	codeTmpId := "123456"
	return svc.sms.Send(ctx, codeTmpId, []string{code}, phone)
}

func (svc *CodeService) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	ok, err := svc.repo.Verify(ctx, biz, phone, code)
	if errors.Is(err, repository.ErrCodeVerifyTooMany) {
		return false, nil
	}
	return ok, err
}

func (svc *CodeService) generate() string {
	code := rand.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}
