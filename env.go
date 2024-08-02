//go:build windows

package wslenv

import (
	"os"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// https://devblogs.microsoft.com/commandline/share-environment-vars-between-wsl-and-windows/
const (
	WSLENV      = "WSLENV"
	ENVIRONMENT = "Environment"
)

func open() (registry.Key, error) {
	return registry.OpenKey(registry.CURRENT_USER, ENVIRONMENT,
		registry.QUERY_VALUE|registry.SET_VALUE)
}

func modify(r registry.Key, key string, flags string, remove bool) error {
	value, _, err := r.GetStringValue(WSLENV)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	found := false
	parts := strings.Split(value, ":")
	items := make([]string, 0, len(parts)+1)
	for _, part := range parts {
		if part == "" {
			continue
		}
		item := strings.SplitN(part, "/", 2)
		if item[0] != key {
			items = append(items, part)
		} else if !remove {
			items = append(items, key+"/"+flags)
			found = true
		}
	}
	if !found && !remove {
		items = append(items, key+"/"+flags)
	}
	value = strings.Join(items, ":")

	if value == "" {
		if os.IsNotExist(err) {
			return nil
		}
		return r.DeleteValue(WSLENV)
	}
	return r.SetStringValue(WSLENV, value)
}

func Publishenv(key string, flags string, publish bool) error {
	r, err := open()
	if err != nil {
		return err
	}
	defer r.Close()
	return modify(r, key, flags, false)
}

func Unpublishenv(key string, flags string, publish bool) error {
	r, err := open()
	if err != nil {
		return err
	}
	defer r.Close()
	return modify(r, key, "", true)
}

func Setenv(key string, value string, flags string, publish bool) error {
	r, err := open()
	if err != nil {
		return err
	}
	defer r.Close()
	err = r.SetStringValue(key, value)
	if err != nil || !publish {
		return err
	}
	return modify(r, key, flags, false)
}

func Unsetenv(key string, publish bool) error {
	r, err := open()
	if err != nil {
		return err
	}
	defer r.Close()
	err = r.DeleteValue(key)
	if err != nil || !publish {
		return err
	}
	return modify(r, key, "", true)
}
