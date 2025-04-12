package tencent

import (
	"context"
	"fmt"
	"github.com/ecodeclub/ekit/slice"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type Service struct {
	client   *sms.Client
	appId    *string
	signName *string
}

func NewService(client *sms.Client, appId, signName string) *Service {
	return &Service{
		client:   client,
		appId:    &appId,
		signName: &signName,
	}
}

func (svc *Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	request := sms.NewSendSmsRequest()
	request.SetContext(ctx)
	request.SmsSdkAppId = svc.appId
	request.SignName = svc.signName
	request.TemplateId = &tplId
	request.TemplateParamSet = svc.toPtrSlice(args)
	request.PhoneNumberSet = svc.toPtrSlice(numbers)
	response, err := svc.client.SendSms(request)
	if err != nil {
		return err
	}

	for _, status := range response.Response.SendStatusSet {
		if status.Code == nil || *status.Code != "Ok" {
			return fmt.Errorf("send sms failed code: %s msg: %s", *status.Code, *status.Message)
		}
	}
	return nil
}

func (svc *Service) toPtrSlice(data []string) []*string {
	return slice.Map[string, *string](data, func(idx int, src string) *string {
		return &src
	})
}
