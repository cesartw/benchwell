// +build linux

package config

import (
	"github.com/gotk3/gotk3/gtk"
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

type providernoop struct{}

func (p *providernoop) Get(w *gtk.Window, path string) (string, error) {
	return path, nil
}

func (p *providernoop) Set(_ *gtk.Window, keys map[string]string, pass string) (string, error) {
	return pass, nil
}
