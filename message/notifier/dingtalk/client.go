package dingtalk

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Client struct {
	accessToken string
	secret      string
	lastSent    time.Time
	minInterval time.Duration
	mu          sync.Mutex
}

func NewClient() *Client {
	return &Client{
		minInterval: 5 * time.Second,
	}
}

func (d *Client) UpdateConfig(cfg *Config) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.accessToken = cfg.AccessToken
	d.secret = cfg.Secret
}

func (d *Client) Notify(message string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	now := time.Now()
	if now.Sub(d.lastSent) < d.minInterval {
		return fmt.Errorf("sending messages too frequently, please wait")
	}

	timestamp := now.UnixNano() / 1e6
	signature := d.sign(timestamp)

	apiURL := fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s&timestamp=%d&sign=%s", d.accessToken, timestamp, signature)

	msg := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"title": "Notification",
			"text":  message,
		},
		"at": map[string]interface{}{
			"isAtAll": true,
		},
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	d.lastSent = now
	return nil
}

func (d *Client) sign(timestamp int64) string {
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, d.secret)
	h := hmac.New(sha256.New, []byte(d.secret))
	h.Write([]byte(stringToSign))
	signData := h.Sum(nil)
	return url.QueryEscape(base64.StdEncoding.EncodeToString(signData))
}
