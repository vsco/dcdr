package redis

import (
	"fmt"

	"errors"

	"strings"

	"github.com/garyburd/redigo/redis"
	"github.com/vsco/dcdr/cli/api/stores"
	"github.com/vsco/dcdr/cli/printer"
	"github.com/vsco/dcdr/config"
)

const DefaultAddr = ":6379"

var (
	ErrRedisConfig = errors.New("missing redis config address")
	ErrIndexRange  = errors.New("K/V index out of range")
)

type Store struct {
	cfg     *config.Config
	conn    redis.Conn
	psConn  redis.Conn
	psc     redis.PubSubConn
	pubChan string
	cb      func(kvb stores.KVBytes)
}

func New(cfg *config.Config) (*Store, error) {
	if cfg.Redis.Address != "" {
		c, err := redis.Dial("tcp", cfg.Redis.Address)
		p, err := redis.Dial("tcp", cfg.Redis.Address)

		if err != nil {
			return nil, err
		}

		s := &Store{
			cfg:     cfg,
			conn:    c,
			psConn:  p,
			pubChan: fmt.Sprintf("%s/*", cfg.Namespace),
		}

		return s, nil
	}

	return nil, ErrRedisConfig
}

func (s *Store) Get(key string) (*stores.KVByte, error) {
	bts, err := redis.Bytes(s.conn.Do("GET", key))

	if err != nil {
		// redigo is super cool and makes you check error strings,
		// this is fine evidently.
		if strings.Contains(err.Error(), "nil returned") {
			return nil, nil
		}

		return nil, err
	}

	return toKVByte(key, bts), nil
}

func (s *Store) Set(key string, bts []byte) error {
	_, err := s.conn.Do("SET", key, bts)

	if err != nil {
		return err
	}

	return s.publish(key, bts)
}

func (s *Store) Delete(key string) error {
	_, err := s.conn.Do("DEL", key)

	if err != nil {
		return err
	}

	return s.publish(key, []byte(""))
}

// List redis requires two operations in order to obtain the values
// from a key prefix. First we append the '*' needed by the `KEYS`
// lookup. Then use the return values as a variadic argument for the
// `MGET` request. See `fetchKeys` and `fetchVals` for implementation.
func (s *Store) List(prefix string) (stores.KVBytes, error) {
	kvb, err := s.fetch(prefix)

	return kvb, err
}

func (s *Store) Watch() error {
	s.psc = redis.PubSubConn{Conn: s.psConn}
	s.psc.PSubscribe(s.cfg.Namespace + "/*")
	defer s.psc.Close()

	s.UpdateKeys()
	s.messageHandler()

	return nil
}

func (s *Store) UpdateKeys() {
	kvbs, err := s.List(s.cfg.Namespace)

	if err != nil {
		printer.LogErrf("pubsub error: %v", err)
	}

	s.Updated(kvbs)
}

func (s *Store) Register(cb func(kvb stores.KVBytes)) {
	s.cb = cb
}

func (s *Store) Updated(kvs interface{}) {
	s.cb(kvs.(stores.KVBytes))
}

func (s *Store) messageHandler() {
	for {
		switch n := s.psc.Receive().(type) {
		case redis.PMessage:
			s.UpdateKeys()
		case error:
			printer.LogErrf("watch error: %v", n.Error())
			return
		}
	}
}

func (s *Store) Close() {
	s.conn.Close()
	s.psConn.Close()
}

// fetch redis needs a KEYS and a MGET in order to grab all the K/Vs
// for a prefix. This function merges them and returns the KVBytes.
func (s *Store) fetch(prefix string) (stores.KVBytes, error) {
	kvb := make(stores.KVBytes, 0)
	keys, err := s.fetchKeys(prefix)

	if err != nil {
		return kvb, err
	}

	if len(keys) == 0 {
		return kvb, nil
	}

	vals, err := s.fetchVals(keys)

	if err != nil {
		return kvb, err
	}

	for i, key := range keys {
		if len(vals) >= i {
			kk := fmt.Sprintf("%s", key)
			kvb = append(kvb, toKVByte(kk, vals[i]))
		} else {
			return kvb, ErrIndexRange
		}
	}

	return kvb, nil
}

func (s *Store) fetchKeys(prefix string) ([]interface{}, error) {
	var strs []interface{}
	keys, err := s.keys(prefix)

	if err != nil {
		return strs, err
	}

	for _, k := range keys {
		strs = append(strs, k)
	}

	return strs, nil
}

func (s *Store) fetchVals(keys []interface{}) ([][]byte, error) {
	vals := make([][]byte, len(keys))
	keys, err := s.mget(keys)

	if err != nil {
		return vals, err
	}

	for i, v := range keys {
		vals[i] = v.([]byte)
	}

	return vals, nil
}

func (s *Store) publish(k string, bts []byte) error {
	_, err := s.conn.Do("PUBLISH", k, bts)

	return err
}

func (s *Store) mget(keys []interface{}) ([]interface{}, error) {
	return redis.Values(s.conn.Do("MGET", keys...))
}

func (s *Store) keys(prefix string) ([]interface{}, error) {
	return redis.Values(s.conn.Do("KEYS", fmt.Sprintf("%s*", prefix)))
}

func toKVByte(key string, bts []byte) *stores.KVByte {
	return &stores.KVByte{
		Key:   key,
		Bytes: bts,
	}
}
