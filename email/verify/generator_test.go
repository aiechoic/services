package verify_test

import (
	"github.com/aiechoic/services/email"
	"github.com/aiechoic/services/email/verify"
	"github.com/aiechoic/services/message/queue"
	"github.com/aiechoic/services/rate"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGenerateAndVerifyCode(t *testing.T) {
	opts := &verify.GeneratorOptions{
		RandomChars:     "0123456789",
		CodeLength:      6,
		CacheExpire:     1 * time.Minute,
		RateLimitEvery:  2 * time.Second,
		RateLimitAllowN: 1,
	}

	pusher := email.NewPusher(queue.NewMemoryQueue[email.Msg]())
	storage := rate.NewMemoryStorage[string](opts.CacheExpire)

	cm := verify.NewGenerator(storage, pusher)
	_ = cm.UpdateConfig(opts)
	emailAddr := "test@example.com"

	// test code generation
	code, err := cm.GenerateCode(emailAddr)
	assert.Nil(t, err)
	assert.NotEqual(t, "", code)

	// test re-generate code wait time
	waitSeconds := int64(cm.GetWaitTime(emailAddr).Seconds())
	expectedWaitSeconds := int64((opts.RateLimitEvery / time.Duration(opts.RateLimitAllowN)).Seconds())
	assert.Equal(t, expectedWaitSeconds, waitSeconds)

	// test code verification
	isValid, err := cm.VerifyCode(emailAddr, code)
	assert.Nil(t, err)
	assert.True(t, isValid)

	// test invalid code verification(one email can only verify once)
	isValid, err = cm.VerifyCode(emailAddr, code)
	assert.Nil(t, err)
	assert.False(t, isValid)

	// test rate limit
	_, _ = cm.GenerateCode(emailAddr)
	code, err = cm.GenerateCode(emailAddr)
	assert.ErrorIs(t, err, rate.ErrLimitExceed)

	// wait until rate limit is reset
	time.Sleep(opts.RateLimitEvery)
	code, err = cm.GenerateCode(emailAddr)
	assert.Nil(t, err)

	// test invalid code verification
	isValid, err = cm.VerifyCode(emailAddr, "invalid_code")
	assert.Nil(t, err)
	assert.False(t, isValid)
}

func TestCodeExpiration(t *testing.T) {

	opts := &verify.GeneratorOptions{
		RandomChars:     "0123456789",
		CodeLength:      6,
		CacheExpire:     2 * time.Second,
		RateLimitEvery:  1 * time.Second,
		RateLimitAllowN: 1,
	}

	pusher := email.NewPusher(queue.NewMemoryQueue[email.Msg]())
	storage := rate.NewMemoryStorage[string](opts.CacheExpire)
	cm := verify.NewGenerator(storage, pusher)
	_ = cm.UpdateConfig(opts)
	emailAddr := "test_expiration@example.com"

	code, _ := cm.GenerateCode(emailAddr)

	// wait until code expires
	time.Sleep(3 * time.Second)

	// the code should be invalid after expiration
	isValid, err := cm.VerifyCode(emailAddr, code)
	assert.Nil(t, err)
	assert.False(t, isValid)
}
