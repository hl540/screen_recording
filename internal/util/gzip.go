package util

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

// Gzip gzip压缩
func Gzip(data []byte) ([]byte, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write([]byte(data)); err != nil {
		return nil, err
	}
	gz.Close()
	if err := gz.Flush(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// Ugzip gzip解压
func Ugzip(data []byte) ([]byte, error) {
	b := bytes.NewReader(data)
	gz, err := gzip.NewReader(b)
	if err != nil {
		return nil, err
	}
	defer gz.Close()
	s, err := ioutil.ReadAll(gz)
	if err != nil {
		return nil, err
	}
	return s, nil
}
