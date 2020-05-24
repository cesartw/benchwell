// +build windows

package config

var Keychain = new(kc)

type kc struct{}

func (k *kc) Set(keys map[string]string, pass string) (string, error) {
	return pass, nil
}

func (k *kc) Get(path string) (string, error) {
	return path, nil
}
