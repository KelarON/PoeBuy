package utils

import (
	"compress/flate"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/zstd"
)

// ReadEncodedResponse Decode any response using one of the encoding methods supported by the Path Of Exile site and return it as a byte slice
func ReadEncodedResponse(res *http.Response) ([]byte, error) {
	var encodedReader io.Reader
	switch res.Header.Get("content-encoding") {
	case "gzip":
		gzipReader, err := gzip.NewReader(res.Body)
		encodedReader = gzipReader
		if err != nil {
			return nil, fmt.Errorf("error decoding gzip response: %v", err)
		}
		defer gzipReader.Close()
	case "br":
		encodedReader = brotli.NewReader(res.Body)
	case "deflate":
		flateReader := flate.NewReader(res.Body)
		encodedReader = flateReader
		defer flateReader.Close()
	case "zstd":
		zstdReader, err := zstd.NewReader(res.Body)
		encodedReader = zstdReader
		if err != nil {
			return nil, fmt.Errorf("error decoding zstd response: %v", err)
		}
		defer zstdReader.Close()
	case "":
		encodedReader = res.Body
	default:
		return nil, fmt.Errorf("unknown encoding type of response, please contact developers to fix this:")
	}

	bt, err := io.ReadAll(encodedReader)
	if err != nil {
		return nil, fmt.Errorf("error reading response bytes: %v", err)
	}

	return bt, nil
}
