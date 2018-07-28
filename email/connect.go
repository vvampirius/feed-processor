package email

import (
	"crypto/tls"
	"net/smtp"
	"net"
)

func ConnectTLS(hostPort, username, password string) (*tls.Conn, *smtp.Client, error) {
	host, _, err := net.SplitHostPort(hostPort)
	if err != nil { return nil, nil, err }

	conn, err := tls.Dial("tcp", hostPort, nil)
	if err != nil { return nil, nil, err }

	c, err := smtp.NewClient(conn, host)
	if err != nil { return nil, nil, err }

	if err = c.Auth(smtp.PlainAuth("", username, password, host)); err != nil {
		return nil, nil, err
	}

	return conn, c, nil
}