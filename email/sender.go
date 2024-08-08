package email

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/aiechoic/services/message/queue"
	"html/template"
	"net/smtp"
	"sync"
)

type Msg struct {
	Email string
	Data  map[string]interface{}
}

type Sender struct {
	opts *SenderConfig
	tpl  *template.Template
	mq   queue.Queue[Msg]
	mu   sync.Mutex
}

func NewSender(mq queue.Queue[Msg]) *Sender {
	return &Sender{
		mq: mq,
	}
}

func (s *Sender) UpdateConfig(opts *SenderConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, err := template.ParseFiles(opts.Template)
	if err != nil {
		return err
	}
	s.opts = opts
	s.tpl = t
	return nil
}

func (s *Sender) Run(ctx context.Context, handler func(err error)) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			pop, err := s.mq.Pop(ctx)
			if err != nil {
				handler(fmt.Errorf("failed to pop message: %w", err))
				continue
			}
			if pop == nil {
				continue
			}
			err = s.send(pop.Email, s.opts.Subject, s.tpl, pop.Data)
			if err != nil {
				handler(fmt.Errorf("failed to send email to %s: %w", pop.Email, err))
			}
		}
	}

}

func (s *Sender) send(to, subject string, tpl *template.Template, data map[string]interface{}) error {
	buf := bytes.NewBuffer(nil)
	err := tpl.Execute(buf, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}
	return SendEmail(s.opts.Title, s.opts.From, s.opts.Password, s.opts.Host, s.opts.Port, to, subject, buf.String())
}

func SendEmail(title, from, pass, host, port, to, subject, body string) error {
	m := "From: " + title + " <" + from + ">\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n" +
		"MIME-Version: 1.0\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\n\n" +
		body

	auth := smtp.PlainAuth("", from, pass, host)

	// 设置 TLS 配置
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	conn, err := tls.Dial("tcp", host+":"+port, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to dial: %w", err)
	}

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// 验证身份
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("failed to auth: %w", err)
	}

	// 设置发件人和收件人
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("failed to set mail: %w", err)
	}
	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set rcpt: %w", err)
	}

	// 获取写入邮件数据的写入器
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get writer: %w", err)
	}

	_, err = writer.Write([]byte(m))
	if err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	err = client.Quit()
	if err != nil {
		return fmt.Errorf("failed to quit: %w", err)
	}
	return nil
}
