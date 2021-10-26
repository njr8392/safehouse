package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"
)

//when encrypting the datastore don't forget to set the origin of the file to the beginning
//AES 256 in GCM mode (supposedly it's safe!
//return ecrypted data... may change
//MAKE SURE THE KEY IS 32 BYTES!!!!
//maybe change file to a read writer????
func Encrypt(key []byte, file *os.File) ([]byte, error) {
	size, _ := FileSize(file)
	txt := make([]byte, size)
	n, err := io.ReadFull(file, txt)
	fmt.Printf("read %d bytes from store", n)
	if err != nil {
		return nil, err
	}
	fmt.Println(txt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce, err := makeNonce(gcm)
	if err != nil {
		return nil, err
	}

	cipher := gcm.Seal(nil, nonce, txt, nil)
	cipher = append(nonce, cipher...)

	return cipher, nil
}

func Decrypt(key []byte, file *os.File) ([]byte, error) {
	size, _ := FileSize(file)
	buf := make([]byte, size)
	_, err := io.ReadFull(file, buf)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := buf[:gcm.NonceSize()]
	buf = buf[gcm.NonceSize():]

	txt, err := gcm.Open(nil, nonce, buf, nil)
	if err != nil {
		return nil, err
	}
	return txt, nil
}

func makeNonce(a cipher.AEAD) ([]byte, error) {
	n := make([]byte, a.NonceSize())
	if _, err := rand.Read(n); err != nil {
		return nil, err
	}
	return n, nil
}

func FileSize(f *os.File) (int64, error) {
	stat, err := f.Stat()
	if err != nil {
		return 0, err
	}

	return stat.Size(), nil
}
