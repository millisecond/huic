package huiccrypt

import (
	"testing"

	"github.com/facebookgo/ensure"
)

func TestRoundtripEncrypt(t *testing.T) {
	plain := make([]byte, 1024)

	key, err := AESKey()
	ensure.Nil(t, err)
	ensure.NotNil(t, key)
	ensure.True(t, len(key) > 0)

	gcm, err := GCMForKey(key)
	ensure.Nil(t, err)
	ensure.NotNil(t, gcm)

	encrypted, nonce, err := AESEncrypt(gcm, plain)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, len(nonce), gcm.NonceSize())

	decrypted, err := AESDecrypt(gcm, encrypted, nonce)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, decrypted, plain)

	//re-encrypt to make sure nonce varies
	_, nonce2, err := AESEncrypt(gcm, plain)
	ensure.Nil(t, err)
	ensure.NotDeepEqual(t, nonce2, nonce)
}
