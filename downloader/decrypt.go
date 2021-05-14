package downloader

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/tls"
	"encoding/binary"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gan-of-culture/go-hentai-scraper/config"
	"github.com/grafov/m3u8"
)

// thanks to https://github.com/canhlinh/hlsdl for creating to beautiful way to decrypt
// I only did some small modifications to the original

const (
	syncByte = uint8(71) //0x47
)

func decryptAES128(crypted, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, iv[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = pkcs5UnPadding(origData)
	return origData, nil
}

func pkcs5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}

// Decrypt descryps a segment
func decrypt(segment *m3u8.MediaSegment, fileName string) ([]byte, error) {

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	if segment.Key != nil {
		key, iv, err := getKey(segment)
		if err != nil {
			return nil, err
		}
		data, err = decryptAES128(data, key, iv)
		if err != nil {
			return nil, err
		}
	}

	for j := 0; j < len(data); j++ {
		if data[j] == syncByte {
			data = data[j:]
			break
		}
	}

	return data, nil
}

func getKey(segment *m3u8.MediaSegment) (key []byte, iv []byte, err error) {
	req, err := http.NewRequest(http.MethodGet, segment.Key.URI, nil)
	if err != nil {
		return nil, nil, errors.New("Request can't be created")
	}

	for k, v := range config.FakeHeaders {
		req.Header.Set(k, v)
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Referer", segment.Key.URI)

	client := &http.Client{
		Transport: &http.Transport{
			DisableCompression:  true,
			TLSHandshakeTimeout: 10 * time.Second,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			DisableKeepAlives:   true,
		},
		Timeout: 15 * time.Minute,
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, nil, errors.New("Failed to get descryption key")
	}

	if res.ContentLength > 64 {
		return nil, nil, errors.New("Invalid key content length for segment")
	}

	key, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}

	iv = []byte(segment.Key.IV)
	if len(iv) == 0 {
		iv = defaultIV(segment.SeqId)
	}
	return
}

func defaultIV(seqID uint64) []byte {
	buf := make([]byte, 16)
	binary.BigEndian.PutUint64(buf[8:], seqID)
	return buf
}
