package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
	"sync"
)

// IntToHex converts an int64 to a byte array
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

const (
	targetBits = 24
)

var (
	maxNonce = math.MaxInt64
)

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{b, target}

	return pow
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)
	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		//fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")

	return nonce, hash[:]
}

// Validate validates block's PoW
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}

func (pow *ProofOfWork) RunParallel(num int) (int, []byte) {

	nonceChan := make(chan int, num)
	resultChan := make(chan int)
	doneChan := make(chan struct{})

	wg := sync.WaitGroup{}
	for i := 0; i < num; i++ {
		go pow.runWorker(wg, nonceChan, resultChan)
	}

	go pow.seedNonces(wg, nonceChan, doneChan)

	nonce := <-resultChan
	doneChan <- struct{}{}

	data := pow.prepareData(nonce)
	hash := sha256.Sum256(data)

	wg.Wait()
	return nonce, hash[:]
}

func (pow *ProofOfWork) seedNonces(wg sync.WaitGroup, nonceChan chan<- int, done <-chan struct{}) {
	wg.Add(1)
	defer wg.Done()

	nonce := 0
	for nonce < maxNonce {
		select {
		case <-done:
			goto DONE
		default:
			nonceChan <- nonce
			nonce++
		}

		if nonce%10000 == 0 {
			fmt.Printf("\r%x", nonce)
		}
	}
	fmt.Println()

DONE:
	close(nonceChan)
	return
}

func (pow *ProofOfWork) runWorker(wg sync.WaitGroup, nonceChan <-chan int, resultChan chan<- int) {
	wg.Add(1)
	defer wg.Done()

	var hashInt big.Int
	var hash [32]byte

	for nonce := range nonceChan {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		//fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			resultChan <- nonce
			return
		}
	}
}
