//go:build js

package inkwell

import (
	"fmt"
	"syscall/js"
)

// LocalStore keeps each key as one window.localStorage entry, prefixed with
// the app's name so two games on one origin do not read each other's slots.
// A browser with storage disabled makes every Read ErrNotStored and every
// Write an error, and the game keeps running without saves.
type LocalStore struct {
	Prefix string
}

func defaultStore(app string) Store { return &LocalStore{Prefix: app + "."} }

func (l *LocalStore) storage() (js.Value, error) {
	storage := js.Global().Get("localStorage")
	if storage.IsUndefined() || storage.IsNull() {
		return js.Value{}, fmt.Errorf("inkwell: localStorage unavailable")
	}
	return storage, nil
}

func (l *LocalStore) Read(key string) ([]byte, error) {
	storage, err := l.storage()
	if err != nil {
		return nil, ErrNotStored
	}
	value := storage.Call("getItem", l.Prefix+key)
	if value.IsNull() || value.IsUndefined() {
		return nil, ErrNotStored
	}
	return []byte(value.String()), nil
}

func (l *LocalStore) Write(key string, data []byte) (err error) {
	storage, err := l.storage()
	if err != nil {
		return err
	}
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("inkwell: localStorage write %s: %v", key, recovered)
		}
	}()
	storage.Call("setItem", l.Prefix+key, string(data))
	return nil
}
