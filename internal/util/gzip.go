package util

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
)

// Gzip gzip压缩
func Gzip(data []byte) ([]byte, error) {
	return data, nil
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	defer gz.Close()
	if _, err := gz.Write([]byte("hello")); err != nil {
		return nil, err
	}
	if err := gz.Flush(); err != nil {
		return nil, err
	}
	return []byte(base64.StdEncoding.EncodeToString(b.Bytes())), nil
}

// Ugzip gzip解压
func Ugzip(data []byte) ([]byte, error) {
	return data, nil
	dataByte, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, err
	}
	b := bytes.NewReader(dataByte)
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
