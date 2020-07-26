// +build linux

package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"

	"github.com/gotk3/gotk3/gtk"
	"github.com/keybase/go-keychain/secretservice"
	dbus "github.com/keybase/go.dbus"
)

// modes: NONE=providernoop, DBUS=providerdbus, FALLBACK=providerbuiltin
func initKeyChain(mode string) {
	Keychain = new(kc)

	switch mode {
	case ModeDBUS:
		p := &providerdbus{}
		err := p.ping()
		if err == nil {
			Keychain.provider = p
			Keychain.Mode = ModeDBUS
			return
		}
		//Env.Log.Errorf("failed to open dbus: %#v", err)

		fallthrough
	case ModeBUILTIN:
		Keychain.provider = &providerbuiltin{}
		Keychain.Mode = ModeBUILTIN
	case ModeNONE:
		Keychain.provider = &providernoop{}
		Keychain.Mode = ModeNONE
	}
}

var Keychain *kc

const (
	ModeNONE    = "NONE"
	ModeDBUS    = "DBUS"
	ModeBUILTIN = "BUILTIN"
)

type kc struct {
	Mode     string
	provider interface {
		Set(w *gtk.Window, keys map[string]string, pass string) (string, error)
		Get(w *gtk.Window, path string) (string, error)
	}
}

func (k *kc) Set(w *gtk.Window, keys map[string]string, pass string) (string, error) {
	return k.provider.Set(w, keys, pass)
}

func (k *kc) Get(w *gtk.Window, path string) (string, error) {
	return k.provider.Get(w, path)
}

type providerdbus struct {
	srv *secretservice.SecretService
}

func (p *providerdbus) Get(_ *gtk.Window, path string) (string, error) {
	session, err := p.srv.OpenSession(secretservice.AuthenticationDHAES)
	if err != nil {
		return "", err
	}
	//defer srv.CloseSession(session)

	secret, err := p.srv.GetSecret(dbus.ObjectPath(path), *session)
	if err != nil {
		return "", err
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

type providernoop struct{}

func (p *providernoop) Get(w *gtk.Window, path string) (string, error) {
	return path, nil
}

func (p *providernoop) Set(_ *gtk.Window, keys map[string]string, pass string) (string, error) {
	return pass, nil
}

type providerbuiltin struct {
	key string
}

func (p *providerbuiltin) Get(w *gtk.Window, path string) (string, error) {
	key, err := p.ask(w)
	if err != nil {
		return "", err
	}

	stored, err := p.decrypt([]byte(key), []byte(path))
	if err != nil {
		return "", err
	}

	return string(stored), nil
}

func (p *providerbuiltin) Set(w *gtk.Window, keys map[string]string, pass string) (string, error) {
	key, err := p.ask(w)
	if err != nil {
		return "", err
	}

	enc, err := p.encrypt([]byte(key), []byte(pass))
	if err != nil {
		return "", err
	}

	return string(enc), err
}

func (p *providerbuiltin) ask(w *gtk.Window) (string, error) {
	if p.key != "" {
		return p.key, nil
	}

	modal, err := gtk.DialogNewWithButtons(
		"Unlock keychain",
		nil,
		gtk.DIALOG_DESTROY_WITH_PARENT|gtk.DIALOG_MODAL,
		[]interface{}{"Unlock", gtk.RESPONSE_ACCEPT},
		[]interface{}{"Cancel", gtk.RESPONSE_CANCEL},
	)
	if err != nil {
		return "", err
	}
	modal.SetDefaultSize(250, 130)
	content, err := modal.GetContentArea()
	if err != nil {
		return "", err
	}

	label, err := gtk.LabelNew("Please enter the password to en/decrypt your saved connections")
	if err != nil {
		return "", err
	}
	entry, err := gtk.EntryNew()
	if err != nil {
		return "", err
	}

	entry.SetProperty("input-purpose", gtk.INPUT_PURPOSE_PASSWORD)
	entry.SetProperty("visibility", false)

	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	if err != nil {
		return "", err
	}

	box.PackStart(label, true, true, 0)
	box.PackStart(entry, true, true, 0)
	content.Add(box)
	content.ShowAll()

	defer modal.Destroy()
	resp := modal.Run()
	if resp != gtk.RESPONSE_ACCEPT {
		return "", nil
	}

	key, err := entry.GetText()
	if err != nil {
		return "", err
	}

	p.key = key
	return key, nil
}

func (p *providerbuiltin) encrypt(key, pass []byte) ([]byte, error) {
	tmp := sha256.Sum256(key)
	key = tmp[:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	b := base64.StdEncoding.EncodeToString(pass)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))

	return []byte(base64.StdEncoding.EncodeToString(ciphertext)), nil
}

func (p *providerbuiltin) decrypt(key, cryptedPass []byte) ([]byte, error) {
	var err error
	cryptedPass, err = base64.StdEncoding.DecodeString(string(cryptedPass))
	if err != nil {
		return nil, err
	}
	tmp := sha256.Sum256(key)
	key = tmp[:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(cryptedPass) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := cryptedPass[:aes.BlockSize]
	cryptedPass = cryptedPass[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(cryptedPass, cryptedPass)

	return base64.StdEncoding.DecodeString(string(cryptedPass))
}
