package main

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"golang.org/x/crypto/ed25519"
)

var numThreads int
var saveKeys bool
var start time.Time // Globale Variable für den Startzeitpunkt

func init() {
	flag.IntVar(&numThreads, "t", 4, "number of threads")
	flag.BoolVar(&saveKeys, "s", false, "save keys to files")
	flag.Parse()
}

func search(target string, wg *sync.WaitGroup) {
	defer wg.Done()

	needle, err := hex.DecodeString(target)
	if err != nil {
		fmt.Println("Decoding failed:", err)
		return
	}

	for {
		publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			fmt.Println("Error generating keys:", err)
			return
		}

		if len(publicKey) >= len(needle) && compare(publicKey[:len(needle)], needle) {
			fmt.Println("public key:", hex.EncodeToString(publicKey))
                        fmt.Println("secret key:", hex.EncodeToString(privateKey))

			if saveKeys {
				err = ioutil.WriteFile("seckey", []byte(hex.EncodeToString(privateKey)), 0600)
				if err != nil {
					fmt.Println("Error writing secret key to file:", err)
					return
				}

				err = ioutil.WriteFile("pubkey", []byte(hex.EncodeToString(publicKey)), 0600)
				if err != nil {
					fmt.Println("Error writing public key to file:", err)
					return
				}
			}

			elapsed := time.Since(start)
			fmt.Printf("Time elapsed: %s\n", elapsed)

			os.Exit(0)
		}
	}
}

func compare(a, b []byte) bool {
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func main() {
	if flag.NArg() != 1 {
		fmt.Printf("usage: %s -t [num threads] -s [save keys] <search string>\n", os.Args[0])
		os.Exit(1)
	}

	target := flag.Arg(0)

	fmt.Printf("searching for key with public part starting with '%s' on %d threads\n", target, numThreads)

	var wg sync.WaitGroup
	wg.Add(numThreads)

	start = time.Now()

	for i := 0; i < numThreads; i++ {
		go search(target, &wg)
	}

	wg.Wait()

}
