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

	sourceBuf := new(bytes.Buffer)
	sourceReader := io.TeeReader(sourceFile, sourceBuf)
	require.NoError(t, err, "read all source file")

	fp := "out.gzip"
	err = compresser.Compress(fp, sourceReader)
	require.NoError(t, err, "compressing source file")

	compressedFile, err := os.Open(fp)
	defer func() {
		err := compressedFile.Close()
		require.NoError(t, err, "closing compressed file")
	}()
	require.NoErrorf(t, err, "opening file: %s", fp)

	decompressedBuf := new(bytes.Buffer)
	gr, err := gzip.NewReader(compressedFile)
	require.NoError(t, err, "new gzip reader")
	defer func() {
		err := gr.Close()
		require.NoError(t, err, "closing gzip reader")
	}()

	_, err = io.Copy(decompressedBuf, gr)
	require.NoError(t, err, "copying gzip reader")

	assert.Equal(t, sourceBuf.Bytes(), decompressedBuf.Bytes(), "buffer bytes")

	err = os.Remove(fp)
	require.NoError(t, err, "removing compressed file")
}
