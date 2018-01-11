package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
)

type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}}
	block.SetHash()
	return block
}

func NewGenesisBlock(msg string) *Block {
	return NewBlock(fmt.Sprintf("Genesis Block: %s", msg), []byte{})
}

func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)

	b.Hash = hash[:]
}

func (b Block) String() string {
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "Prev. hash: %x\n", b.PrevBlockHash)
	fmt.Fprintf(buf, "Data: %s\n", b.Data)
	fmt.Fprintf(buf, "Hash: %x\n", b.Hash)
	return buf.String()
}

type Blockchain struct {
	blocks []*Block
}

func NewBlockchain(genesis string) *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock(genesis)}}
}

func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}

func main() {
	bc := NewBlockchain("One must still have chaos in oneself to be able to give birth to a dancing star.")

	bc.AddBlock("Send 1 BTC to Ivan")
	bc.AddBlock("Send 2 more BTC to Ivan")

	for _, block := range bc.blocks {
		fmt.Println(block)
	}
}
