package tunnel

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rc4"
	"encoding/binary"
	"log"

	"github.com/codahale/chacha20"
	"golang.org/x/crypto/salsa20/salsa"
)

type DecOrEnc int

const (
	Decrypt DecOrEnc = iota
	Encrypt
)

type Cipher struct {
	enc  cipher.Stream
	dec  cipher.Stream
	key  []byte
	info *cipherInfo
}

type cipherInfo struct {
	keyLen    int
	ivLen     int
	newStream func(key, iv []byte, doe DecOrEnc) (cipher.Stream, error)
}

func newStream(block cipher.Block, err error, key, iv []byte, doe DecOrEnc) (cipher.Stream, error) {
	if err != nil {
		return nil, err
	}
	if doe == Encrypt {
		return cipher.NewCFBEncrypter(block, iv), nil
	}
	return cipher.NewCFBDecrypter(block, iv), nil
}

// Initializes the block cipher with CFB mode, returns IV.
func (c *Cipher) initEncrypt() (iv []byte, err error) {
	iv = c.key[:c.info.ivLen]
	c.enc, err = c.info.newStream(c.key, iv, Encrypt)
	if err != nil {
		return nil, err
	}
	return
}

func (c *Cipher) initDecrypt(iv []byte) (err error) {
	c.dec, err = c.info.newStream(c.key, iv, Decrypt)
	return
}

func (c *Cipher) encrypt(dst, src []byte) {
	c.enc.XORKeyStream(dst, src)
}

func (c *Cipher) decrypt(dst, src []byte) {
	c.dec.XORKeyStream(dst, src)
}

var cipherMethod = map[string]*cipherInfo{
	"blank":     {0, 0, nil},
	"rc4":       {16, 0, nil},
	"rc4-md5":   {16, 16, newRC4MD5Stream},
	"aes256cfb": {32, 16, newAESStream},
	"chacha20":  {32, 8, newChaCha20Stream},
	"salsa20":   {32, 8, newSalsa20Stream},
}

func secretToKey(secret []byte, size int) []byte {
	// size mod 16 must be 0
	h := md5.New()
	buf := make([]byte, size)
	count := size / md5.Size
	// repeatly fill the key with the secret
	for i := 0; i < count; i++ {
		h.Write(secret)
		copy(buf[md5.Size*i:md5.Size*(i+1)-1], h.Sum(nil))
	}
	return buf
}

func newBlankCipher() (enc, dec cipher.Stream) {
	enc = &blankStreamCipher{}
	return enc, enc
}

func newRC4Cipher(key []byte) (enc, dec cipher.Stream) {
	rc4Enc, err := rc4.NewCipher(key)
	if err != nil {
		//return
		log.Fatal(err)
	}
	// create a copy, as RC4 encrypt and decrypt uses the same keystream
	rc4Dec := *rc4Enc
	return rc4Enc, &rc4Dec
}

func newRC4MD5Stream(key, iv []byte, _ DecOrEnc) (cipher.Stream, error) {
	h := md5.New()
	h.Write(key)
	h.Write(iv)
	rc4key := h.Sum(nil)

	return rc4.NewCipher(rc4key)
}

func newAESStream(key, iv []byte, doe DecOrEnc) (cipher.Stream, error) {
	block, err := aes.NewCipher(key)
	return newStream(block, err, key, iv, doe)
}

func newChaCha20Stream(key, iv []byte, _ DecOrEnc) (cipher.Stream, error) {
	return chacha20.New(key, iv)
}

func newSalsa20Stream(key, iv []byte, _ DecOrEnc) (cipher.Stream, error) {
	var c salsaStreamCipher
	copy(c.nonce[:], iv[:8])
	copy(c.key[:], key[:32])
	return &c, nil
}

type blankStreamCipher struct {
}

func (c *blankStreamCipher) XORKeyStream(dst, src []byte) {
	copy(dst, src)
	return
}

type salsaStreamCipher struct {
	nonce   [8]byte
	key     [32]byte
	counter int
}

func (c *salsaStreamCipher) XORKeyStream(dst, src []byte) {
	var buf []byte
	padLen := c.counter % 64
	dataSize := len(src) + padLen
	if cap(dst) >= dataSize {
		buf = dst[:dataSize]
	} else if leakyBufSize >= dataSize {
		buf = leakyBuf.Get()
		defer leakyBuf.Put(buf)
		buf = buf[:dataSize]
	} else {
		buf = make([]byte, dataSize)
	}

	var subNonce [16]byte
	copy(subNonce[:], c.nonce[:])
	binary.LittleEndian.PutUint64(subNonce[len(c.nonce):], uint64(c.counter/64))

	// It's difficult to avoid data copy here. src or dst maybe slice from
	// Conn.Read/Write, which can't have padding.
	copy(buf[padLen:], src[:])
	salsa.XORKeyStream(buf, buf, &subNonce, &c.key)
	copy(dst, buf[padLen:])

	c.counter += len(src)
}

func NewCipher(cryptoMethod string, secret []byte) *Cipher {
	mi, ok := cipherMethod[cryptoMethod]
	if !ok {
		log.Fatalf("unsupported crypto method %s", cryptoMethod)
	}

	key := secretToKey(secret, mi.keyLen)

	c := &Cipher{key: key, info: mi}

	if mi.newStream == nil {
		if cryptoMethod == "blank" {
			c.enc, c.dec = newBlankCipher()
		} else {
			c.enc, c.dec = newRC4Cipher(key)
		}
	} else {
		iv, err := c.initEncrypt()
		if err != nil {
			log.Fatal(err)
		}

		if err := c.initDecrypt(iv); err != nil {
			log.Fatal(err)
		}
	}

	return c
}
