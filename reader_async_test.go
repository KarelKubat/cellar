package cellar

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReader_ScanAsync(t *testing.T) {
	db, err := New(dbDir, WithNoFileLock, WithMetaDB(newBoltMetaDB()))
	require.NoError(t, err)

	// Write some testdata
	testdata := "TestReader_ScanAsync"
	_, err = db.Append([]byte(testdata))
	require.NoError(t, err)

	db.Flush()
	reader := db.Reader()

	var passed bool
	vals, errchan := reader.ScanAsync(context.Background(), 1)
	for {
		select {
		case <-errchan:
			require.NoError(t, <-errchan)
		case v, ok := <-vals:
			if !ok {
				break
			}
			if string(v.Data) == testdata {
				passed = true
				break
			}
		}
		if passed {
			break
		}
	}
	assert.True(t, passed)

}
