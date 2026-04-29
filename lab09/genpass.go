package main

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"os"
	"strings"
)

const itoa64 = "./0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func encode64(v uint32, n int) string {
	var b strings.Builder
	for i := 0; i < n; i++ {
		b.WriteByte(itoa64[v&0x3f])
		v >>= 6
	}
	return b.String()
}

func randomSalt() string {
	raw := make([]byte, 6)
	rand.Read(raw)
	return encode64(uint32(raw[0])<<16|uint32(raw[1])<<8|uint32(raw[2]), 4) +
		encode64(uint32(raw[3])<<16|uint32(raw[4])<<8|uint32(raw[5]), 4)
}

func apr1(password, salt string) string {
	final := md5.New()
	ctx := md5.New()
	ctx.Write([]byte(password))
	ctx.Write([]byte("$apr1$"))
	ctx.Write([]byte(salt))
	final.Write([]byte(password))
	final.Write([]byte("$apr1$"))
	final.Write([]byte(salt))

	plen := len(password)
	for plen > 0 {
		n := 16
		if plen < 16 {
			n = plen
		}
		ctx.Write(final.Sum(nil)[:n])
		plen -= 16
	}

	for i := len(password); i > 0; i >>= 1 {
		if i&1 != 0 {
			ctx.Write([]byte{0})
		} else {
			ctx.Write([]byte{password[0]})
		}
	}

	result := ctx.Sum(nil)
	for i := 0; i < 1000; i++ {
		c2 := md5.New()
		if i&1 != 0 {
			c2.Write([]byte(password))
		} else {
			c2.Write(result)
		}
		if i%3 != 0 {
			c2.Write([]byte(salt))
		}
		if i%7 != 0 {
			c2.Write([]byte(password))
		}
		if i&1 != 0 {
			c2.Write(result)
		} else {
			c2.Write([]byte(password))
		}
		result = c2.Sum(nil)
	}

	var hash strings.Builder
	hash.WriteString("$apr1$")
	hash.WriteString(salt)
	hash.WriteByte('$')
	hash.WriteString(encode64((uint32(result[0])<<16|uint32(result[6])<<8|uint32(result[12])), 4))
	hash.WriteString(encode64((uint32(result[1])<<16|uint32(result[7])<<8|uint32(result[13])), 4))
	hash.WriteString(encode64((uint32(result[2])<<16|uint32(result[8])<<8|uint32(result[14])), 4))
	hash.WriteString(encode64((uint32(result[3])<<16|uint32(result[9])<<8|uint32(result[15])), 4))
	hash.WriteString(encode64((uint32(result[4])<<16|uint32(result[10])<<8|uint32(result[5])), 4))
	hash.WriteString(encode64(uint32(result[11]), 2))
	return hash.String()
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run genpass.go <user> <pass>")
		os.Exit(1)
	}
	user, pass := os.Args[1], os.Args[2]
	hash := apr1(pass, randomSalt())
	os.WriteFile("htpasswd.txt", []byte(fmt.Sprintf("%s:%s\n", user, hash)), 0644)
	fmt.Printf("htpasswd.txt created: %s\n", user)
}
