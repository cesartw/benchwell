package mysql

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/user"
	"strconv"
	"strings"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type Endpoint struct {
	Host string
	Port int
}

func (endpoint *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
}

type SSHtunnel struct {
	Local  *Endpoint
	Server *Endpoint
	Remote *Endpoint

	Config *ssh.ClientConfig
}

func (tunnel *SSHtunnel) Start() error {
	listener, err := net.Listen("tcp", tunnel.Local.String())
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go tunnel.forward(conn)
	}
}

func (tunnel *SSHtunnel) forward(localConn net.Conn) {
	serverConn, err := ssh.Dial("tcp", tunnel.Server.String(), tunnel.Config)
	if err != nil {
		fmt.Printf("Server dial error: %s\n", err)
		return
	}

	remoteConn, err := serverConn.Dial("tcp", tunnel.Remote.String())
	if err != nil {
		fmt.Printf("Remote dial error: %s\n", err)
		return
	}

	copyConn := func(writer, reader net.Conn) {
		defer writer.Close()
		defer reader.Close()

		_, err := io.Copy(writer, reader)
		if err != nil {
			fmt.Printf("io.Copy error: %s", err)
		}
	}

	go copyConn(localConn, remoteConn)
	go copyConn(remoteConn, localConn)
}

func SSHAgent(socket string) (ssh.AuthMethod, error) {
	sshAgent, err := net.Dial("unix", os.Getenv(socket))
	if err != nil {
		return nil, err
	}

	client := agent.NewClient(sshAgent)
	return ssh.PublicKeysCallback(client.Signers), nil
}

func PrivateKeyFile(file, phrase string) (ssh.AuthMethod, error) {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParsePrivateKeyWithPassphrase(buffer, []byte(phrase))
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(key), nil
}

func init() {
	mysql.RegisterDial("ssh", func(addr string) (net.Conn, error) {
		var (
			localHost string
			localPort int
			sshHost   string
			sshSocket string
			sshPort   int
		)

		parts := strings.Split(addr, ",")
		if len(parts) != 2 {
			return nil, errors.New("wrong method address format")
		}

		local := parts[0]  // localhost:3306
		remote := parts[1] // localhost;SSH_AUTH_SOCK

		parts = strings.Split(local, ":")
		switch len(parts) {
		case 1:
			localHost = parts[0]
			localPort = 3306
		case 2:
			localHost = parts[0]
			localPort, _ = strconv.Atoi(parts[1])
		default:
			return nil, errors.New("wrong local address format")
		}

		parts = strings.Split(remote, ";")
		switch len(parts) {
		case 2:
			sshHost = parts[0]
			sshPort = 22
			sshSocket = parts[1]
		case 3:
			sshHost = parts[0]
			sshPort, _ = strconv.Atoi(parts[1])
			sshSocket = parts[2]
		default:
			return nil, errors.New("wrong remote address format")
		}

		localEndpoint := &Endpoint{
			Host: "localhost",
			Port: 9999,
		}

		serverEndpoint := &Endpoint{
			Host: sshHost,
			Port: sshPort,
		}

		remoteEndpoint := &Endpoint{
			Host: localHost,
			Port: localPort,
		}

		u, err := user.Current()
		if err != nil {
			return nil, err
		}

		auth, err := SSHAgent(sshSocket)
		if err != nil {
			return nil, err
		}

		sshConfig := &ssh.ClientConfig{
			User: u.Username,
			Auth: []ssh.AuthMethod{
				auth,
			},
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
		}

		tunnel := &SSHtunnel{
			Config: sshConfig,
			Local:  localEndpoint,
			Server: serverEndpoint,
			Remote: remoteEndpoint,
		}

		go tunnel.Start()

		return net.Dial("tcp", "localhost:9999")
	})
}
