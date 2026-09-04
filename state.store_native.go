//go:build !js

package inkwell

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// FileStore keeps each key as <dir>/<key>.json. A write goes to a temporary
// file first and is renamed into place, so a crash mid-write leaves the
// previous document rather than half of the new one.
type FileStore struct {
	Dir string
}

// NewFileStore is the store for an app on this machine: os.UserConfigDir
// with the app's name under it — ~/Library/Application Support/<app> on
// macOS, ~/.config/<app> on Linux, %AppData%\<app> on Windows. When the
// user's directory cannot be found it falls back to <app> under the working
// directory, so the game still saves somewhere.
func NewFileStore(app string) *FileStore {
	base, err := os.UserConfigDir()
	if err != nil {
		base = "."
	}
	return &FileStore{Dir: filepath.Join(base, app)}
}

func defaultStore(app string) Store { return NewFileStore(app) }

func (f *FileStore) path(key string) string {
	return filepath.Join(f.Dir, key+".json")
}

func (f *FileStore) Read(key string) ([]byte, error) {
	data, err := os.ReadFile(f.path(key))
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrNotStored
	}
	if err != nil {
		return nil, fmt.Errorf("inkwell: store read %s: %w", key, err)
	}
	return data, nil
}

func (f *FileStore) Write(key string, data []byte) error {
	if err := os.MkdirAll(f.Dir, 0o755); err != nil {
		return fmt.Errorf("inkwell: store mkdir: %w", err)
	}
	temp := f.path(key) + ".tmp"
	if err := os.WriteFile(temp, data, 0o644); err != nil {
		return fmt.Errorf("inkwell: store write %s: %w", key, err)
	}
	if err := os.Rename(temp, f.path(key)); err != nil {
		return fmt.Errorf("inkwell: store rename %s: %w", key, err)
	}
	return nil
}
