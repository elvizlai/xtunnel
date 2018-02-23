package tunnel

import (
	"testing"
)

func TestBlank(t *testing.T) {
	secret := []byte("testsecret")
	c := NewCipher("blank", secret)
	testCiphter(t, c, "blank")
}

func TestRC4(t *testing.T) {
	secret := []byte("testsecret")
	c := NewCipher("rc4", secret)
	testCiphter(t, c, "rc4")
}

func TestRC4MD5(t *testing.T) {
	secret := []byte("testsecret")
	c := NewCipher("rc4-md5", secret)
	testCiphter(t, c, "rc4-md5")
}

func TestAES256CFB(t *testing.T) {
	secret := []byte("testsecret")
	c := NewCipher("aes256cfb", secret)
	testCiphter(t, c, "aes256cfb")
}

func TestSALSA20(t *testing.T) {
	secret := []byte("testsecret")
	c := NewCipher("salsa20", secret)
	testCiphter(t, c, "salsa20")
}

func TestCHACHA20(t *testing.T) {
	secret := []byte("testsecret")
	c := NewCipher("chacha20", secret)
	testCiphter(t, c, "chacha20")
}

const text = "Don't tell me the moon is shining; show me the glint of light on broken glass."

func testCiphter(t *testing.T, c *Cipher, msg string) {
	n := len(text)

	cipherBuf := make([]byte, n)
	originTxt := make([]byte, n)

	c.encrypt(cipherBuf, []byte(text))
	c.decrypt(originTxt, cipherBuf)

	if string(originTxt) != text {
		t.Error(msg, "encrypt then decrytp does not get original text")
	}
}

// benchmark

const size = 1000000

func BenchmarkBlank(b *testing.B) {
	benchmark(b, "blank", size)
}

func BenchmarkRC4(b *testing.B) {
	benchmark(b, "rc4", size)
}

func BenchmarkRC4MD5(b *testing.B) {
	benchmark(b, "rc4-md5", size)
}

func BenchmarkAES256CFB(b *testing.B) {
	benchmark(b, "aes256cfb", size)
}

func BenchmarkCHACHA20(b *testing.B) {
	benchmark(b, "chacha20", size)
}

func BenchmarkSALSA20(b *testing.B) {
	benchmark(b, "salsa20", size)
}

func benchmark(b *testing.B, cryptoMethod string, size int) {
	secret := []byte("testsecret")
	c := NewCipher(cryptoMethod, secret)
	b.N = size
	for i := 0; i < b.N; i++ {
		n := len(text)

		cipherBuf := make([]byte, n)
		originTxt := make([]byte, n)

		c.encrypt(cipherBuf, []byte(text))
		c.decrypt(originTxt, cipherBuf)
	}
}
