//go:build windows

package wslenv

import (
	"strings"

	"golang.org/x/sys/windows/registry"
)

const (
	lxssKey   = "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Lxss"
	defEnvKey = "DefaultEnvironment"
)

// as defined in wslservices.exe
var defEnvValue = [...]string{
	"HOSTTYPE=x86_64",
	"LANG=en_US.UTF-8",
	"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games:/usr/local/games",
	"TERM=xterm-256color",
}

type LxssKey struct {
	r registry.Key
}

func Open(uuid string) (*LxssKey, error) {
	r, err := registry.OpenKey(registry.CURRENT_USER, lxssKey+"\\{"+uuid+"}",
		registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		return nil, err
	}
	return &LxssKey{r}, nil
}

func (k *LxssKey) Close() error {
	return k.r.Close()
}

func (k *LxssKey) GetEnv(name string) (string, error) {
	vals, valtype, err := k.r.GetStringsValue(defEnvKey)
	if err != nil && err != registry.ErrNotExist {
		return "", err
	}
	if valtype != registry.NONE && valtype != registry.MULTI_SZ {
		return "", registry.ErrUnexpectedType
	}
	if len(vals) == 0 {
		vals = defEnvValue[:]
	}
	for _, val := range vals {
		pair := strings.SplitN(val, "=", 2)
		if pair[0] == name {
			return pair[1], nil
		}
	}
	return "", registry.ErrNotExist
}

func (k *LxssKey) SetEnv(name string, value string) error {
	vals, valtype, err := k.r.GetStringsValue(defEnvKey)
	if err != nil && err != registry.ErrNotExist {
		return err
	}
	if valtype != registry.NONE && valtype != registry.MULTI_SZ {
		return registry.ErrUnexpectedType
	}
	if len(vals) == 0 {
		vals = defEnvValue[:]
	}
	found := false
	for i, val := range vals {
		pair := strings.SplitN(val, "=", 2)
		if pair[0] == name {
			vals[i] = name + "=" + value
			found = true
			break
		}
	}
	if !found {
		vals = append(vals, name+"="+value)
	}
	return k.r.SetStringsValue(defEnvKey, vals)
}

func (k *LxssKey) UnsetEnv(name string) error {
	vals, valtype, err := k.r.GetStringsValue(defEnvKey)
	if err != nil && err != registry.ErrNotExist {
		return err
	}
	if valtype != registry.NONE && valtype != registry.MULTI_SZ {
		return registry.ErrUnexpectedType
	}
	if len(vals) == 0 {
		return nil
	}
	items := make([]string, 0, len(vals))
	for _, val := range vals {
		pair := strings.SplitN(val, "=", 2)
		if pair[0] != name {
			items = append(items, val)
		}
	}
	return k.r.SetStringsValue(defEnvKey, items)
}
