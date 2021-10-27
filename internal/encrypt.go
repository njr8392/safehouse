package internal

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"os"
)

//when encrypting the datastore don't forget to set the origin of the file to the beginning
//AES 256 in GCM mode
//all passwords will be hashed with sha256 to gurantee that the key size will be 32 bytes
//maybe change file to a read writer????
func Encrypt(key []byte, file *os.File) ([]byte, error) {
	size, _ := fileSize(file)

	txt := make([]byte, size)
	_, err := io.ReadFull(file, txt)
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
	nonce, err := makeNonce(gcm)
	if err != nil {
		return nil, err
	}

	cipher := gcm.Seal(nil, nonce, txt, nil)
	cipher = append(nonce, cipher...)

	return cipher, nil
}

func Decrypt(key []byte, file *os.File) ([]byte, error) {
	size, _ := fileSize(file)
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

func fileSize(f *os.File) (int64, error) {
	stat, err := f.Stat()
	if err != nil {
		return 0, err
	}

	return stat.Size(), nil
}

//helper function to convert a sha256 hash into a []byte.
//aes.NewCipher takes an []byte as an argument
func CopySha256(hash [32]byte) []byte {
	var b []byte
	for _, h := range hash {
		b = append(b, h)
	}
	return b
}
