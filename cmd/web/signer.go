package main

import (
	"fmt"
	"strings"
	"time"

	goalone "github.com/bwmarrin/go-alone"
)

// Normalde activation link gönderdiğimizde => http://site.com/?email=abc@mail.com&plan=gold gibi olur.
// Ama böyle yapınca parametreler kötü niyetli kişiler tarafından brute force maruz kalabilir.
// Bu yüzden bu script ile gönderilen aktivasyon linkinde "TAMPER PROOF URL" elde edilmesi amaçlanmıştır.
// http://site.com/?params=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9....

// normalde bunu daha karmaşık bir string olarak yazarız ve kod içerisinde tutmayız.
// environment variable içerisinde tutulabilir.
const secret = "abc123abc123abc123"

var secretKey []byte

// NewURLSigner creates a new signer
func NewURLSigner() {
	secretKey = []byte(secret)
}

// GenerateTokenFromString generates a signed token
func GenerateTokenFromString(data string) string {
	var urlToSign string

	s := goalone.New(secretKey, goalone.Timestamp)
	if strings.Contains(data, "?") {
		urlToSign = fmt.Sprintf("%s&hash=", data)
	} else {
		urlToSign = fmt.Sprintf("%s?hash=", data)
	}

	tokenBytes := s.Sign([]byte(urlToSign))
	token := string(tokenBytes)

	return token
}

// VerifyToken verifies a signed token
func VerifyToken(token string) bool {
	s := goalone.New(secretKey, goalone.Timestamp)
	_, err := s.Unsign([]byte(token))

	// signature is not valid. Token was tampered with, forged, or maybe it's
	// not even a token at all! Either way, it's not safe to use it. (err != nil)
	// valid has (err == nil)
	return err == nil
}

// Expired checks to see if a token has expired
func Expired(token string, minutesUntilExpire int) bool {
	s := goalone.New(secretKey, goalone.Timestamp)
	ts := s.Parse([]byte(token))

	// time.Duration(seconds)*time.Second
	return time.Since(ts.Timestamp) > time.Duration(minutesUntilExpire)*time.Minute
}
