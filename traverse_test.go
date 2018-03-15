package maxminddb

import (
	"fmt"
	"testing"

	"github.com/gotestyourself/gotestyourself/assert"
	is "github.com/gotestyourself/gotestyourself/assert/cmp"
)

func TestNetworks(t *testing.T) {
	for _, recordSize := range []uint{24, 28, 32} {
		for _, ipVersion := range []uint{4, 6} {
			fileName := fmt.Sprintf("test-data/test-data/MaxMind-DB-test-ipv%d-%d.mmdb", ipVersion, recordSize)
			reader, err := Open(fileName)
			assert.NilError(t, err, "unexpected error while opening database: %v", err)
			defer reader.Close()

			n := reader.Networks()
			for n.Next() {
				record := struct {
					IP string `maxminddb:"ip"`
				}{}
				network, err := n.Network(&record)
				assert.Check(t, err)
				assert.Check(t, is.Equal(record.IP, network.IP.String()), "expected %s got %s", record.IP, network.IP.String())

			}
			assert.Check(t, n.Err())
		}
	}
}

func TestNetworksWithInvalidSearchTree(t *testing.T) {
	reader, err := Open("test-data/test-data/MaxMind-DB-test-broken-search-tree-24.mmdb")
	assert.NilError(t, err, "unexpected error while opening database: %v", err)
	defer reader.Close()

	n := reader.Networks()
	for n.Next() {
		var record interface{}
		_, err := n.Network(&record)
		assert.Check(t, err)
	}
	assert.Check(t, n.Err() != nil, "no error received when traversing an broken search tree")
	assert.Check(t, is.Equal(n.Err().Error(), "invalid search tree at 128.128.128.128/32"))
}
