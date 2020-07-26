// +build linux

package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

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

	enc, hash, err := p.encrypt([]byte(key), []byte(pass))
	if err != nil {
		return "", err
	}

	return string(enc) + "/" + hash, err
}

func (p *providerbuiltin) ask(w *gtk.Window) (string, error) {
	if p.key != "" {
		return p.key, nil
	}

	modal, err := gtk.DialogNewWithButtons(
		"Unlock keychain",
		w,
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

func (p *providerbuiltin) encrypt(key, pass []byte) ([]byte, string, error) {
	tmp := sha256.Sum256(key)
	key = tmp[:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, "", err
	}

	b := base64.StdEncoding.EncodeToString(pass)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, "", err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))

	sum := sha256.Sum256([]byte(pass))
	hash := fmt.Sprintf("%x", sum)

	return []byte(base64.StdEncoding.EncodeToString(ciphertext)), hash, nil
}

func (p *providerbuiltin) decrypt(key, cryptedPass []byte) ([]byte, error) {
	parts := strings.Split(string(cryptedPass), "/")
	cryptedPass = []byte(strings.Join(parts[:len(parts)-1], "/"))
	storedHash := parts[len(parts)-1]

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
	pass, err := base64.StdEncoding.DecodeString(string(cryptedPass))
	if err != nil {
		return nil, errors.New("wrong key")
	}

	sum := sha256.Sum256([]byte(pass))
	hash := fmt.Sprintf("%x", sum)
	if storedHash != hash {
		return nil, errors.New("wrong key")
	}
	return pass, nil
}
