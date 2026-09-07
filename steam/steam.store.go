package steam

import (
	"fmt"
	"log"

	inkwell "git.teletypegames.org/engines/inkwell"
	"github.com/bogdannbv/go-steamworks"
)

// cloudStore is an inkwell.Store over ISteamRemoteStorage for the keys the
// game sends to the cloud, and over the platform's own Store for every
// other key. A key is a file, <key>.json, so the cloud holds the same
// documents the disk would.
//
// It degrades toward the disk: a cloud read that finds nothing falls back to
// the local document, which is how a save made before the game had Steam
// reaches the load screen and gets uploaded on the next write; and a cloud
// write that fails is written locally instead, logged, and not reported as
// an error, because losing the player's save over a quota is worse than
// keeping it on one machine.
type cloudStore struct {
	cloud map[string]bool
	local inkwell.Store
}

func newCloudStore(keys []string, local inkwell.Store) *cloudStore {
	cloud := make(map[string]bool, len(keys))
	for _, key := range keys {
		cloud[key] = true
	}
	return &cloudStore{cloud: cloud, local: local}
}

func cloudFile(key string) string { return key + ".json" }

func (store *cloudStore) Read(key string) ([]byte, error) {
	if !store.cloud[key] {
		return store.local.Read(key)
	}
	remote := steamworks.SteamRemoteStorage()
	size := remote.GetFileSize(cloudFile(key))
	if size <= 0 {
		return store.local.Read(key)
	}
	data := make([]byte, size)
	if read := remote.FileRead(cloudFile(key), data); read != size {
		return nil, fmt.Errorf("steam: cloud read %s: %d of %d bytes", key, read, size)
	}
	return data, nil
}

func (store *cloudStore) Write(key string, data []byte) error {
	if !store.cloud[key] || len(data) == 0 {
		return store.local.Write(key, data)
	}
	if !steamworks.SteamRemoteStorage().FileWrite(cloudFile(key), data) {
		log.Printf("steam: cloud write %s failed; written locally", key)
		return store.local.Write(key, data)
	}
	return nil
}
