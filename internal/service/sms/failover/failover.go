package failover

import (
	"context"
	"errors"
	"fire/internal/service/sms"
	"log"
	"sync/atomic"
)

type SimpleFailOverSMSService struct {
	services []sms.Service
	idx      uint64
}

func (f SimpleFailOverSMSService) SendV1(ctx context.Context, tplId string, args []string, numbers ...string) error {
	idx := atomic.AddUint64(&f.idx, 1)
	length := uint64(len(f.services))

	for i := idx; i < length; i++ {
		svc := f.services[i%length]
		err := svc.Send(ctx, tplId, args, numbers...)
		switch err {
		case nil:
			return nil
		case context.Canceled, context.DeadlineExceeded:
			return err
		}
		log.Println(err)
	}
	return errors.New("send SMS failed for all services")
}

func (f *SimpleFailOverSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	for _, service := range f.services {
		err := service.Send(ctx, tplId, args, numbers...)
		if err == nil {
			return err
		}
		log.Println(err)

	}
	return errors.New("all SMS service failed to send")
}

func NewFailOverSMSService(services []sms.Service) *SimpleFailOverSMSService {
	return &SimpleFailOverSMSService{
		services: services,
	}
}

type TimeoutFailOverSMSService struct {
	services  []sms.Service
	idx       int32
	cnt       int32
	threshold int32
}

func (t TimeoutFailOverSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	idx := atomic.LoadInt32(&t.idx)
	cnt := atomic.LoadInt32(&t.cnt)

	if cnt >= t.threshold {
		newIdx := (idx + 1) % int32(len(t.services))
		if atomic.CompareAndSwapInt32(&t.idx, idx, newIdx) {
			atomic.StoreInt32(&t.cnt, 0)
		}
		idx = newIdx
	}
	svc := t.services[idx]
	err := svc.Send(ctx, tplId, args, numbers...)
	switch err {
	case nil:
		atomic.StoreInt32(&t.cnt, 0)
		return nil
	case context.DeadlineExceeded:
		atomic.AddInt32(&t.cnt, 1)
		return err
	}
	return err
}
