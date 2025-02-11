package redis

import (
	"errors"

	"github.com/boj/redistore"
	"github.com/gin-contrib/sessions"
	"github.com/gomodule/redigo/redis"
)

type Store interface {
	sessions.Store
}

// size: maximum number of idle connections.
// network: tcp or udp
// address: host:port
// password: redis-password
// Keys are defined in pairs to allow key rotation, but the common case is to set a single
// authentication key and optionally an encryption key.
//
// The first key in a pair is used for authentication and the second for encryption. The
// encryption key can be set to nil or omitted in the last pair, but the authentication key
// is required in all pairs.
//
// It is recommended to use an authentication key with 32 or 64 bytes. The encryption key,
// if set, must be either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256 modes.
func NewStore(size int, network, address, password string, keyPairs ...[]byte) (Store, error) {
	s, err := redistore.NewRediStore(size, network, address, password, keyPairs...)
	if err != nil {
		return nil, err
	}
	return &store{s}, nil
}

// NewStoreWithDB - like NewStore but accepts `DB` parameter to select
// redis DB instead of using the default one ("0")
//
// Ref: https://godoc.org/github.com/boj/redistore#NewRediStoreWithDB
func NewStoreWithDB(size int, network, address, password, DB string, keyPairs ...[]byte) (Store, error) {
	s, err := redistore.NewRediStoreWithDB(size, network, address, password, DB, keyPairs...)
	if err != nil {
		return nil, err
	}
	return &store{s}, nil
}

// NewStoreWithPool instantiates a RediStore with a *redis.Pool passed in.
//
// Ref: https://godoc.org/github.com/boj/redistore#NewRediStoreWithPool
func NewStoreWithPool(pool *redis.Pool, keyPairs ...[]byte) (Store, error) {
	s, err := redistore.NewRediStoreWithPool(pool, keyPairs...)
	if err != nil {
		return nil, err
	}
	return &store{s}, nil
}

type store struct {
	*redistore.RediStore
}

// GetRedisStore get the actual woking store.
// Ref: https://godoc.org/github.com/boj/redistore#RediStore
func GetRedisStore(s Store) (err error, rediStore *redistore.RediStore) {
	realStore, ok := s.(*store)
	if !ok {
		err = errors.New("unable to get the redis store: Store isn't *store")
		return
	}

	rediStore = realStore.RediStore
	return
}

// SetKeyPrefix sets the key prefix in the redis database.
func SetKeyPrefix(s Store, prefix string) error {
	err, rediStore := GetRedisStore(s)
	if err != nil {
		return err
	}

	rediStore.SetKeyPrefix(prefix)
	return nil
}

// SetMaxLength sets the max length for a new session in the redis database.
// For more information check https://github.com/boj/redistore/blob/cd5dcc76aeff9ba06b0a924829fe24fd69cdd517/redistore.go#L101
func SetMaxLength(s Store, length int) error {
	err, rediStore := GetRedisStore(s)
	if err != nil {
		return err
	}

	rediStore.SetMaxLength(length)
	return nil
}

func (c *store) Options(options sessions.Options) {
	c.RediStore.Options = options.ToGorillaOptions()
}
