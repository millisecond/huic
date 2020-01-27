package huiccrypt

import (
	"crypto/cipher"
)

func AESEncrypt(gcm cipher.AEAD, plain []byte) (encrypted []byte, nonce []byte, err error) {
	// generate a new nonce for every Seal invocation
	nonce, err = generateNonce(gcm)
	if err != nil {
		return
	}
	encrypted = gcm.Seal(encrypted, nonce, plain, nil)
	return
}

func AESDecrypt(gcm cipher.AEAD, encrypted []byte, nonce []byte) (plain []byte, err error) {
	return gcm.Open(plain, nonce, encrypted, nil)
}
