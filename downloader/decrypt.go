package downloader

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"io"
	"os"
)

// thanks to https://github.com/canhlinh/hlsdl for creating to beautiful way to decrypt
// I only did some small modifications to the original

/*const (
	syncByte = uint8(71) //0x47
)*/

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

/*func pkcs5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}*/

func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}

// Decrypt descryps a segment
func decrypt(key []byte, fileName string) ([]byte, error) {

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	iv := defaultIV(uint64(0))
	data, err = decryptAES128(data, key, iv)
	if err != nil {
		return nil, err
	}

	// sync byte is part of the Transport Stream (ts) header which is 4 bytes
	// it works by ditching all data that is infront of the sync byte
	// this is not something you would want for anything besides Transport Streams
	// it's also not really necessary here since this is a complete download so there is no syncing
	/*for j := 0; j < len(data); j++ {
		if data[j] == syncByte {
			data = data[j:]
			break
		}
	}*/

	return data, nil
}

func defaultIV(seqID uint64) []byte {
	buf := make([]byte, 16)
	binary.BigEndian.PutUint64(buf[8:], seqID)
	return buf
}
