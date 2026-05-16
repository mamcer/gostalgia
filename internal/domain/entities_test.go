package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntityTableNames(t *testing.T) {
	assert.Equal(t, "nfile", NFile{}.TableName())
	assert.Equal(t, "ndirectory", NDirectory{}.TableName())
	assert.Equal(t, "nfilenode", NFileNode{}.TableName())
	assert.Equal(t, "nscan", NScan{}.TableName())
	assert.Equal(t, "ntag", NTag{}.TableName())
}
