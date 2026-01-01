package main

import (
	"crypto/tls"
	"fmt"

	"github.com/nats-io/nats.go"
)

func main() {
	// 在 Go 裡面寫這個，它絕對能跳過 TLS 驗證
	nc, _ := nats.Connect("nats://100.81.97.99:4222",
		nats.UserInfo("oktopususer", "oktopuspw"),
		nats.Secure(&tls.Config{InsecureSkipVerify: true}))

	// 訂閱所有訊息並印出來
	nc.Subscribe(">", func(m *nats.Msg) {
		fmt.Printf("主題: %s | 內容: %s\n", m.Subject, string(m.Data))
	})

	// 讓程式別結束
	select {}
}
