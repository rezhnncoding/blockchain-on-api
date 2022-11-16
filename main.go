package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"

	"time"
)

type Block struct {
	Pos       int
	Data      BookCheckOut
	TimeStamp string
	Hash      string
	PrevHash  string
}

type Book struct {
	Id          string `json:"Id"`
	Title       string `json:"Title"`
	Author      string `json:"Author"`
	PublishDate string `json:"PublishDate"`
	Isbn        string `json:"isbn"`
}

type BookCheckOut struct {
	BookId       string `json:"BookId"`
	User         string `json:"User"`
	CheckOutDate string `json:"CheckOutDate"`
	IsGenesis    bool   `json:"IsGenesis"`
}

type BlockChain struct {
	blocks []*Block
}

var blockchain *BlockChain

func CreateBlock(prevblock *Block, checkoutitem BookCheckOut) *Block {
	block := &Block{}
	block.Pos = prevblock.Pos + 1
	block.PrevHash = prevblock.Hash
	block.TimeStamp = time.Now().String()
	block.GenerateHash()
	return block
}

func (b *Block) GenerateHash() {
	bytes, _ := json.Marshal(b.Data)
	data := string(b.Pos) + b.TimeStamp + string(bytes) + b.PrevHash
	hash := sha256.New()
	hash.Write([]byte(data))
	b.Hash = hex.EncodeToString(hash.Sum(nil))
}

func (bc *BlockChain) AddBlock(data BookCheckOut) {
	prevblock := bc.blocks[len(bc.blocks)-1]
	block := CreateBlock(prevblock, data)
	if ValidBlock(block, prevblock) {
		bc.blocks = append(bc.blocks, block)
	}

}

func NewBook(w http.ResponseWriter, r *http.Request) {
	var book Book

	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("could not create:%v", err)
		w.Write([]byte("could not create new book"))
		return
	}
	h := md5.New()
	io.WriteString(h, book.Isbn+book.PublishDate)
	book.Id = fmt.Sprintf("%x", h.Sum(nil))

	res, err := json.MarshalIndent(book, "", "")
	if err != nil {
		w.WriteHeader(500)
		log.Printf("could not marshal payload:%v", err)
		w.Write([]byte("could not save book data"))
		return
	}
	w.WriteHeader(200)
	w.Write(res)
}

func GetBlockChain(w http.ResponseWriter, r *http.Request) {
	jbytes, err := json.MarshalIndent(blockchain.blocks, "", "")

	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(err)
		return
	}
	io.WriteString(w, string(jbytes))
}

func ValidBlock(block, prevblock *Block) bool {
	if prevblock.Hash != block.PrevHash {
		return false
	}
	if !block.validatehash(block.Hash) {
		return false
	}
	if prevblock.Pos+1 != block.Pos {
		return false
	}
	return true
}

func (b *Block) validatehash(hash string) bool {
	b.GenerateHash()
	if b.Hash != hash {
		return false
	}
	return true
}

func WriteBlock(w http.ResponseWriter, r *http.Request) {
	var checkoutitem BookCheckOut
	if err := json.NewDecoder(r.Body).Decode(&checkoutitem); err != nil {
		w.WriteHeader(500)
		log.Printf("could not write block: %v", err)
		w.Write([]byte("cold not write block"))
	}
	blockchain.AddBlock(checkoutitem)
}

func GenesisBlock() *Block {
	return CreateBlock(&Block{}, BookCheckOut{IsGenesis: true})
}

func NewBlockChain() *BlockChain {
	return &BlockChain{[]*Block{GenesisBlock()}}
}

func main() {

	blockchain = NewBlockChain()
	r := mux.NewRouter()
	r.HandleFunc("/", GetBlockChain).Methods("GET")
	r.HandleFunc("/", WriteBlock).Methods("POST")
	r.HandleFunc("/new", NewBook).Methods("POST")

	go func() {
		for _, block := range blockchain.blocks {
			fmt.Printf("prev.hash:%x\n", block.PrevHash)
			bytes, _ := json.MarshalIndent(block.Data, "", "")
			fmt.Printf("data:%v\n", string(bytes))
			fmt.Printf("hash:%x,\n", block.Hash)
			fmt.Println()
		}
	}()

	log.Println("listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))

}
