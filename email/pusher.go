package email

import (
	"context"
	"fmt"
	"github.com/aiechoic/services/message/queue"
	"time"
)

type Pusher struct {
	mq queue.Queue[Msg]
}

func NewPusher(mq queue.Queue[Msg]) *Pusher {
	return &Pusher{
		mq: mq,
	}
}

func (p *Pusher) Push(ctx context.Context, email string, data map[string]interface{}, expire time.Duration) error {
	if email == "" {
		return fmt.Errorf("email is empty")
	}
	msg := Msg{
		Email: email,
		Data:  data,
	}
	err := p.mq.Push(ctx, &msg, expire)
	if err != nil {
		return fmt.Errorf("failed to push message: %w", err)
	}
	return nil
}
