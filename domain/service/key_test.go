package service

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"testing"
)

func TestA(t *testing.T) {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	fmt.Println("D:",privateKey.D)
	fmt.Println("E:",privateKey.E)
	fmt.Println("N:",privateKey.N)
	fmt.Println("primes : ",privateKey.Primes)
	fmt.Println("N:",privateKey.PublicKey.N)
	fmt.Println("E:",privateKey.PublicKey.E)
}
