// +build linux

package config

import (
	"errors"
	"time"

	"github.com/gotk3/gotk3/gtk"
	"github.com/keybase/go-keychain/secretservice"
	dbus "github.com/keybase/go.dbus"
)

type providerdbus struct {
	srv *secretservice.SecretService
}

func (p *providerdbus) Get(_ *gtk.Window, path string) (string, error) {
	t := time.NewTimer(3 * time.Second)

	c := make(chan error, 1)

	var secret []byte
	go func() {
		var err error
		session, err := p.srv.OpenSession(secretservice.AuthenticationDHAES)
		if err != nil {
			c <- err
			return
		}

		secret, err = p.srv.GetSecret(dbus.ObjectPath(path), *session)
		if err != nil {
			c <- err
			return
		}
		c <- nil
	}()

	select {
	case err := <-c:
		if err != nil {
			return "", err
		}
	case <-t.C:
		return "", errors.New("DBUS timeout")
	}

	return string(secret), nil
}

func (p *providerdbus) Set(_ *gtk.Window, keys map[string]string, pass string) (string, error) {
	session, err := p.srv.OpenSession(secretservice.AuthenticationDHAES)
	if err != nil {
		return "", err
	}
	//defer p.srv.CloseSession(session)

	err = p.srv.Unlock([]dbus.ObjectPath{secretservice.DefaultCollection})
	if err != nil {
		return "", err
	}

	secret, err := session.NewSecret([]byte(pass))
	if err != nil {
		return "", err
	}

	item, err := p.srv.CreateItem(
		secretservice.DefaultCollection,
		secretservice.NewSecretProperties("benchwell", keys),
		secret,
		secretservice.ReplaceBehaviorReplace,
	)

	return string(item), nil
}

func (p *providerdbus) ping() (err error) {
	p.srv, err = secretservice.NewService()
	if err != nil {
		return err
	}

	return p.srv.Unlock([]dbus.ObjectPath{secretservice.DefaultCollection})
}
