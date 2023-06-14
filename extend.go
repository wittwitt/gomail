package gomail

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
)

const (
	GOMAIL_SERVER       = "GOMAIL_SERVER" // addr:port, if 465, use ssl
	GOMAIL_ACCOUNT      = "GOMAIL_ACCOUNT"
	GOMAIL_ACCOUNT_NAME = "GOMAIL_ACCOUNT_NAME"
	GOMAIL_PASS         = "GOMAIL_PASS"
)

var gCfg *Cfg

type Cfg struct {
	Host        string
	Port        int
	Account     string
	Pass        string
	AccountName string
}

func init() {
	InitEnv()
}

func InitEnv() {
	cfg, err := NewCfgFromEnv()
	if err != nil {
		log.Printf("NewSmtpClientFromEnv: %v", err)
		return
	}
	gCfg = cfg
}

func NewCfgFromEnv() (*Cfg, error) {
	gomailServerEnv := os.Getenv(GOMAIL_SERVER)
	if gomailServerEnv == "" {
		return nil, fmt.Errorf("GOMAIL_SERVER not set")
	}

	host, portStr, err := net.SplitHostPort(gomailServerEnv)
	if err != nil {
		return nil, fmt.Errorf("GOMAIL_SERVER: %v, %v", gomailServerEnv, err)
	}
	portI64, _ := strconv.ParseInt(portStr, 10, 64)

	cfg := &Cfg{
		Host:        host,
		Port:        int(portI64),
		Account:     os.Getenv(GOMAIL_ACCOUNT),
		AccountName: os.Getenv(GOMAIL_ACCOUNT_NAME),
		Pass:        os.Getenv(GOMAIL_PASS),
	}
	return cfg, nil
}

func SendOnce(ctx context.Context, to []string, cc []string, subject string, body string) error {
	if gCfg == nil {
		return fmt.Errorf("no smtp client ")
	}

	msg := NewMessage()
	msg.SetAddressHeader("From", gCfg.Account, gCfg.AccountName)
	msg.SetHeader("To", to...)
	for _, ccAddr := range cc {
		msg.SetAddressHeader("Cc", ccAddr, ccAddr)
	}
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	cnn := NewDialer(gCfg.Host, gCfg.Port, gCfg.Account, gCfg.Pass)
	return cnn.DialAndSend(msg)
}
