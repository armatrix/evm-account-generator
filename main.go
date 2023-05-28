package main

import (
	"crypto/ecdsa"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {

	var wg sync.WaitGroup
	ch := make(chan struct{}, 15)
	for i := 0; i < 100000000; i++ {
		ch <- struct{}{}
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			address, publicKey, privateKey := genAddress()
			if condition(address) {
				fmt.Println("Address:", address)
				fmt.Println("Public Key:", publicKey)
				fmt.Println("SAVE BUT DO NOT SHARE THIS (Private Key):", privateKey)
			}
			<-ch
		}(i)
	}
	wg.Wait()
	fmt.Println("job done")
}

func genAddress() (accAddr, pubKey, prvKey string) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	return address, hexutil.Encode(publicKeyBytes), hexutil.Encode(privateKeyBytes)
}

func condition(addr string) bool {
	lower := strings.ToLower(addr)
	prefix := strings.HasPrefix(lower, "000000")
	suffix := strings.HasSuffix(lower, "000000") || true
	return prefix && suffix
}
