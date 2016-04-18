package redis

import (
	"testing"

	"errors"

	"github.com/stretchr/testify/assert"
	"github.com/vsco/dcdr/config"
)

type MockRedis struct {
	Response interface{}
	Error    error
}

func (r *MockRedis) Close() error {
	return nil
}

func (r *MockRedis) Err() error {
	return nil
}

func (r *MockRedis) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	return r.Response, r.Error
}

func (r *MockRedis) Send(commandName string, args ...interface{}) error {
	return nil
}

func (r *MockRedis) Flush() error {
	return nil
}

func (r *MockRedis) Receive() (reply interface{}, err error) {
	return nil, nil
}

func TestGet(t *testing.T) {
	r := &Store{
		cfg: config.TestConfig(),
		conn: &MockRedis{
			Response: "abcde",
		},
	}

	kv, err := r.Get("foo")

	assert.NoError(t, err)
	assert.Equal(t, "foo", kv.Key)
	assert.Equal(t, []byte("abcde"), kv.Bytes)
}

func TestGetNilError(t *testing.T) {
	r := &Store{
		cfg: config.TestConfig(),
		conn: &MockRedis{
			Error: errors.New("nil returned"),
		},
	}

	kv, err := r.Get("foo")

	assert.Nil(t, kv)
	assert.NoError(t, err)
}

func TestFetchKeys(t *testing.T) {
	r := &Store{
		cfg: config.TestConfig(),
		conn: &MockRedis{
			Response: []interface{}{[]byte("k")},
		},
	}

	bts, err := r.fetchKeys("k")

	assert.NoError(t, err)
	assert.Equal(t, []byte("k"), bts[0])
}

func TestFetch(t *testing.T) {
	r := &Store{
		cfg: config.TestConfig(),
		conn: &MockRedis{
			Response: []interface{}{[]byte("x")},
		},
	}

	bts, err := r.fetch("k")

	assert.NoError(t, err)
	assert.Equal(t, "x", bts[0].Key)
	assert.Equal(t, []byte("x"), bts[0].Bytes)
}
