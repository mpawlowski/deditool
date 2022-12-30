package extractor

import (
	"io/fs"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const WORK_DIR = `work`

func cleanupWorkDir(t *testing.T) {
	if strings.HasPrefix(WORK_DIR, "/") {
		assert.Fail(t, "unable to cleanup, WORK_DIR was an absolute directory")
		return
	}
	err := os.RemoveAll(WORK_DIR)
	assert.Nil(t, err)
}

func TestExportParser(t *testing.T) {
	defer cleanupWorkDir(t)

	extractor := NewExtractor()
	assert.NotNil(t, extractor)

	err := extractor.Extract("testdata/server-7.atflaclabs.com-20221230.r2z", WORK_DIR, ".")
	assert.Nil(t, err)

	dirs, err := os.ReadDir(WORK_DIR)
	assert.Nil(t, err)

	for _, d := range dirs {
		info, err := d.Info()
		assert.Nil(t, err)
		assert.Equal(t, fs.ModePerm, info.Mode(), "Files should be mode.Perm")
	}
}
