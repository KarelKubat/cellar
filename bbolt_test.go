package cellar

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"

	"go.etcd.io/bbolt"
)

var mu = &sync.Mutex{}
var asked int

func newBoltMetaDB() (db *BoltMetaDB) {
	mu.Lock()
	defer mu.Unlock()
	asked++

	blt, err := bolt.Open(fmt.Sprintf("testdata/meta%d.bolt", asked), 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		panic(err)
	}
	db = &BoltMetaDB{DB: blt}
	err = db.Init()
	if err != nil {
		panic(err)
	}
	return db
}

func TestBoltMetaDB_AddChunk_ListChunk(t *testing.T) {
	db := newBoltMetaDB()

	err := db.AddChunk(0, &ChunkDto{})
	require.NoError(t, err)

	err = db.AddChunk(10, &ChunkDto{})
	require.NoError(t, err)

	chunks, err := db.ListChunks()
	require.NoError(t, err)
	assert.True(t, len(chunks) == 2)
}
