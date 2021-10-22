package main	


//when encrypting the datastore don't forget to set the origin of the file to the beginning 
//AES 256 in GCM mode (supposedly it's safe!
//return ecrypted data... may change
//MAKE SURE THE KEY IS 32 BYTES!!!!
func Encrypt(key []byte, nonce, []byte, file io.ReadWriter)([]byte,error){
	var txt []byte
	_, err := io.ReadFull(file, txt)
	if err != nil{
		return nil,err
	}
	block, err := aes.NewCipher(key)
	if err != nil{
		return nil,err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil{
		return nil,err
	}
	
	cipher := gcm.Seal(nil,nonce, txt, nil)

	return cipher, nil
}

func Decrypt(key []byte, file io.ReadWriter)([]byte, error){
	var cipher []byte
	_, err := io.ReadFull(file, cipher)
	if err != nil{
		return nil,err
	}
	
	block, err := aes.NewCipher(key)
	if err != nil{
		return nil,err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil{
		return nil,err
	}

	nonce := cipher[:gcm.NonceSize()]
	cipher = cipher[gcm.NonceSie():]

	txt, err := gcm.Open(nil, nonce, cipher, nil)
	if err != nil{
		return nil,err
	}
	return txt, nil
}

func makeNonce(a cipher.AEAD)([]byte, error){
	n := make([]btye, a.NonceSize())
	if _, err := rand.Read(n); err != nil{
		return nil, err
	}
	return n
}
