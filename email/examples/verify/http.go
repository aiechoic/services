package main

import (
	"context"
	"fmt"
	"github.com/aiechoic/services/email"
	"github.com/aiechoic/services/email/verify"
	"github.com/aiechoic/services/ioc"
	"github.com/aiechoic/services/ioc/healthy"
	"log"
	"net/http"
	"time"
)

var indexHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Email Verification</title>
    <style>
        body {
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
            margin: 0;
            font-family: Arial, sans-serif;
        }
        form {
            display: flex;
            flex-direction: column;
            align-items: center;
            border: 1px solid #ccc;
            padding: 20px;
            border-radius: 10px;
            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
        }
        label, input, button {
            margin: 10px 0;
        }
        button {
            padding: 10px 20px;
            cursor: pointer;
        }
    </style>
</head>
<body>
    <form id="emailForm" action="/send" method="get">
        <label for="email">Email:</label>
        <input type="email" id="email" name="email" required>
        <label for="code">Verification Code:</label>
        <input type="text" id="code" name="code">
        <button type="button" onclick="sendCode()">Send Code</button>
        <button type="button" onclick="verifyCode()">Verify Code</button>
    </form>

    <script>
        function sendCode() {
            const form = document.getElementById('emailForm');
            form.action = '/send';
            form.submit();
        }

        function verifyCode() {
            const form = document.getElementById('emailForm');
            form.action = '/verify';
            form.submit();
        }
    </script>
</body>
</html>
`

func main() {
	addr := ":8866"
	c := ioc.NewContainer()
	defer c.Close()

	err := c.LoadConfig("./configs", ioc.ConfigEnvTest)
	if err != nil {
		panic(err)
	}

	generator := verify.GetGenerator(c)

	http.HandleFunc("/verify", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		email := r.URL.Query().Get("email")
		if code == "" || email == "" {
			w.Write([]byte("code and email are required"))
			return
		}
		ok, err := generator.VerifyCode(email, code)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		if ok {
			w.Write([]byte("ok"))
		} else {
			w.Write([]byte("fail"))
		}
	})

	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")
		if email == "" {
			w.Write([]byte("email is required"))
			return
		}
		wait := generator.GetWaitTime(email)
		if wait > 0 {
			w.Write([]byte("wait for " + wait.String()))
			return
		}
		// the code will push to redis queue
		_, err := generator.GenerateCode(email)
		if err != nil {
			w.Write([]byte(err.Error()))
		} else {
			wait = generator.GetWaitTime(email)
			msg := fmt.Sprintf("success, the code is pushed to message queue\n")
			msg += fmt.Sprintf("wait for %s", wait.String())
			w.Write([]byte(msg))
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(indexHTML))
	})

	// start checking email message queue, and send email, you can run this in another program
	go RunSender(c)

	// start health check
	go RunHealthCheck(c)

	fmt.Printf("server started at http://localhost%s\n\n", addr)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("http server error: %v", err)
	}
}

func RunSender(c *ioc.Container) {
	sender := email.GetSender(c)

	ctx := context.Background()

	// start checking email message queue, and send email
	sender.Run(ctx, func(err error) {
		log.Println(err)
	})
}

func RunHealthCheck(c *ioc.Container) {
	ticker := 20 * time.Second
	checkTimeout := 5 * time.Second
	c.RunHealthCheck(ticker, checkTimeout, func(errs []*healthy.Error) {
		if len(errs) > 0 {
			for _, err := range errs {
				log.Printf("health check error: %v\n", err)
			}
		} else {
			log.Printf("health check OK\n")
		}
	})
}
