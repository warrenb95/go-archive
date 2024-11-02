package compresser_test

import (
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/warrenb95/go-archive/internal/compresser"
)

func TestCompress(t *testing.T) {
	sourceFile, err := os.Open("test_data/test.txt")
	defer func() {
		err := sourceFile.Close()
		require.NoError(t, err, "closing source file")
	}()
	require.NoError(t, err, "opening test.txt")

	filepath, err := compresser.Compress(sourceFile)
	require.NoError(t, err, "compressing source file")

	compressedFile, err := os.Open(filepath)
	defer func() {
		err := compressedFile.Close()
		require.NoError(t, err, "closing compressed file")
	}()
	require.NoErrorf(t, err, "opening file: %s", filepath)

	buf := new(bytes.Buffer)
	gr, err := gzip.NewReader(compressedFile)
	require.NoError(t, err, "new gzip reader")
	defer func() {
		err := gr.Close()
		require.NoError(t, err, "closing gzip reader")
	}()

	_, err = io.Copy(buf, gr)
	require.NoError(t, err, "copying gzip reader")

	// TODO: fix issue with sourceBytes being empty
	// re-read the file?
	sourceBytes, err := io.ReadAll(sourceFile)
	require.NoError(t, err, "read all source file")

	assert.Equal(t, sourceBytes, buf.Bytes(), "buffer bytes")
}
