package etcd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vsco/dcdr/config"
)

func TestEtcd(t *testing.T) {
	s, err := DefaultETCDlStore(config.DefaultConfig())

	assert.NoError(t, err)

	fmt.Printf("%+v", s)
}
