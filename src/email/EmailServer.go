package email

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
)

type emailServer struct {
	tlsConfig *tls.Config
	auth      smtp.Auth
	conn      *tls.Conn
}

func autenticate(account, password, host string) (smtp.Auth, *tls.Config) {
	auth := smtp.PlainAuth("", account, password, host)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	return auth, tlsConfig
}

func makeConn(serverName string, tlsConfig *tls.Config) (*tls.Conn, error) {
	conn, err := tls.Dial("tcp", serverName, tlsConfig)
	if err != nil {
		return conn, err
	}

	return conn, nil
}

func newClient(conn *tls.Conn, host string, auth smtp.Auth) (*smtp.Client, error) {
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if err := client.Auth(auth); err != nil {
		log.Fatalf("error al autenticar el cliente")
	}

	return client, nil
}

func NewEmailServer(account, password, host, serverName string) emailServer {
	smptAuth, tlsConfig := autenticate(account, password, host)
	Conn, err := makeConn(serverName, tlsConfig)
	if err != nil {
		log.Fatal(err)
	}

	return emailServer{
		tlsConfig: tlsConfig,
		auth:      smptAuth,
		conn:      Conn,
	}
}
