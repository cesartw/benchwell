// +build linux

package config

import (
	"github.com/keybase/go-keychain/secretservice"
	dbus "github.com/keybase/go.dbus"
)

var Keychain = new(kc)

type kc struct{}

func (k *kc) Set(keys map[string]string, pass string) (string, error) {
	srv, err := secretservice.NewService()
	if err != nil {
		return "", err
	}

	session, err := srv.OpenSession(secretservice.AuthenticationDHAES)
	if err != nil {
		return "", err
	}
	defer srv.CloseSession(session)

	err = srv.Unlock([]dbus.ObjectPath{secretservice.DefaultCollection})
	if err != nil {
		return "", err
	}

	secret, err := session.NewSecret([]byte(pass))
	if err != nil {
		return "", err
	}

	item, err := srv.CreateItem(
		secretservice.DefaultCollection,
		secretservice.NewSecretProperties("sqlaid", keys),
		secret,
		secretservice.ReplaceBehaviorReplace,
	)

	return string(item), nil
}

func (k *kc) Get(path string) (string, error) {
	srv, err := secretservice.NewService()
	if err != nil {
		return "", err
	}

	session, err := srv.OpenSession(secretservice.AuthenticationDHAES)
	if err != nil {
		return "", err
	}
	defer srv.CloseSession(session)

	secret, err := srv.GetSecret(dbus.ObjectPath(path), *session)
	if err != nil {
		return "", err
	}

	return string(secret), nil
}
