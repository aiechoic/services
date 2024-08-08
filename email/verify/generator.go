package verify

import (
	"context"
	"github.com/aiechoic/services/email"
	"github.com/aiechoic/services/internal/random"
	"github.com/aiechoic/services/rate"
	"time"
)

type Generator struct {
	limiter     *rate.Limiter[string]
	storage     rate.Storage[string]
	pusher      *email.Pusher
	opts        *GeneratorOptions
	randomChars []rune
}

type GeneratorOptions struct {
	RandomChars     string
	CodeLength      int
	CacheExpire     time.Duration
	RateLimitEvery  time.Duration
	RateLimitAllowN int
}

func NewGenerator(storage rate.Storage[string], pusher *email.Pusher) *Generator {
	return &Generator{
		pusher:  pusher,
		storage: storage,
	}
}

func (cm *Generator) UpdateConfig(opts *GeneratorOptions) error {
	if opts.RandomChars == "" {
		opts.RandomChars = "0123456789"
	}
	if opts.CodeLength <= 0 {
		opts.CodeLength = 6
	}
	cm.opts = opts
	cm.randomChars = []rune(opts.RandomChars)
	cm.limiter = rate.NewRateLimiter[string](cm.storage, opts.RateLimitEvery, opts.RateLimitAllowN)
	return nil
}

func (cm *Generator) generateCode() string {
	return random.StringWithCharset(cm.opts.CodeLength, cm.randomChars)
}

func (cm *Generator) GetWaitTime(email string) time.Duration {
	return cm.limiter.GetWaiteTime(email)
}

func (cm *Generator) GenerateCode(email string) (code string, err error) {
	code = cm.generateCode()
	err = cm.limiter.Set(email, &code, cm.opts.CacheExpire)
	if err != nil {
		return "", err
	}
	err = cm.pusher.Push(
		context.Background(),
		email,
		map[string]interface{}{
			"code":     code,
			"expireIn": int64(cm.opts.CacheExpire / time.Minute),
		},
		cm.opts.CacheExpire,
	)
	return code, nil
}

func (cm *Generator) VerifyCode(email, code string) (bool, error) {
	c, err := cm.limiter.Get(email)
	if err != nil {
		return false, err
	}
	if c == nil {
		return false, nil
	}
	if *c != code {
		err = cm.limiter.DelData(email)
		return false, err
	}
	err = cm.limiter.Del(email)
	if err != nil {
		return false, err
	}
	return true, nil
}
