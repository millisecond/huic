package huiccrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

// Note - NOT RFC4122 compliant, but endeavoring to have zero dependencies
// From https://stackoverflow.com/questions/15130321/is-there-a-method-to-generate-a-uuid-with-go-language
// Discussion about how this should be more random on reasonable hardware
func PseudoUUID() (uuid string, err error) {
	b := make([]byte, 16)
	_, err = rand.Read(b)
	if err != nil {
		return
	}
	uuid = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return
}

func AESKey() (key []byte, err error) {
	key = make([]byte, 32)
	_, err = rand.Read(key)
	return
}

func GCMForKey(key []byte) (gcm cipher.AEAD, err error) {
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	// gcm or Galois/Counter Mode, is a mode of operation
	// for symmetric key cryptographic block ciphers
	// - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	gcm, err = cipher.NewGCM(blockCipher)
	return
}

// Do NOT re-use a nonce across encryptions!
func generateNonce(gcm cipher.AEAD) (nonce []byte, err error) {
	nonce = make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	return
}
