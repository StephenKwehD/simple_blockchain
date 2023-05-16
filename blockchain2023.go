package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

type MerkleTree struct {
	transactions []string
	root         string
}

func NewMerkleTree(transactions []string) *MerkleTree {
	tree := &MerkleTree{
		transactions: transactions,
		root:         "",
	}
	tree.buildMerkleTree()
	return tree
}

func (mt *MerkleTree) buildMerkleTree() {
	if len(mt.transactions) == 0 {
		return
	}

	hashedTransactions := make([]string, len(mt.transactions))
	for i, tx := range mt.transactions {
		hashedTransactions[i] = mt.hashTransaction(tx)
	}

	for len(hashedTransactions) > 1 {
		newTransactions := []string{}
		for i := 0; i < len(hashedTransactions)-1; i += 2 {
			combinedHash := mt.hashTransaction(hashedTransactions[i] + hashedTransactions[i+1])
			newTransactions = append(newTransactions, combinedHash)
		}
		if len(hashedTransactions)%2 != 0 {
			newTransactions = append(newTransactions, hashedTransactions[len(hashedTransactions)-1])
		}
		hashedTransactions = newTransactions
	}

	mt.root = hashedTransactions[0]
}

func (mt *MerkleTree) hashTransaction(transaction string) string {
	h := sha256.New()
	h.Write([]byte(transaction))
	return hex.EncodeToString(h.Sum(nil))
}

type Block struct {
	index        int
	timestamp    int64
	transactions []string
	previousHash string
	hash         string
	merkleRoot   string
	nonce        int
}

func NewBlock(index int, timestamp int64, transactions []string, previousHash string) *Block {
	block := &Block{
		index:        index,
		timestamp:    timestamp,
		transactions: transactions,
		previousHash: previousHash,
		hash:         "",
		merkleRoot:   "",
		nonce:        0,
	}
	block.calculateHash()
	block.calculateMerkleRoot()
	return block
}

func (b *Block) calculateHash() {
	for {
		blockString := fmt.Sprintf(`{"index":%d,"timestamp":%d,"transactions":%v,"previous_hash":"%s","nonce":%d}`, b.index, b.timestamp, b.transactions, b.previousHash, b.nonce)
		h := sha256.New()
		h.Write([]byte(blockString))
		blockHash := hex.EncodeToString(h.Sum(nil))
		if blockHash[:4] == "0000" {
			b.hash = blockHash
			break
		}
		b.nonce++
	}
}

func (b *Block) calculateMerkleRoot() {
	mt := NewMerkleTree(b.transactions)
	b.merkleRoot = mt.root
}

type Blockchain struct {
	chain               []*Block
	currentTransactions []string
}

func NewBlockchain() *Blockchain {
	blockchain := &Blockchain{
		chain:               []*Block{NewGenesisBlock()},
		currentTransactions: []string{},
	}
	return blockchain
}

func NewGenesisBlock() *Block {
	return NewBlock(0, time.Now().Unix(), []string{}, "0")
}

func (bc *Blockchain) GetLatestBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *Blockchain) AddTransaction(transaction string) {
	bc.currentTransactions = append(bc.currentTransactions, transaction)
}

func (bc *Blockchain) CreateNewBlock() {
	newBlock := NewBlock(len(bc.chain), time.Now().Unix(), bc.currentTransactions, bc.GetLatestBlock().hash)
	bc.currentTransactions = []string{}
	bc.AddBlock(newBlock)
}

func (bc *Blockchain) AddBlock(newBlock *Block) {
	newBlock.previousHash = bc.GetLatestBlock().hash
	newBlock.calculateHash()
	bc.chain = append(bc.chain, newBlock)
}

func (bc *Blockchain) IsChainValid() bool {
	for i := 1; i < len(bc.chain); i++ {
		currentBlock := bc.chain[i]
		previousBlock := bc.chain[i-1]
		if currentBlock.hash != "" && currentBlock.hash != currentBlock.hash {
			return false
		}
		if currentBlock.previousHash != previousBlock.hash {
			return false
		}
	}
	return true
}

func main() {
	blockchain := NewBlockchain()

	for {
		fmt.Println("Enter student information (or type 'q' to quit):")
		var stdID, stdName, stdSurname, stdGender, stdBirthYear string
		fmt.Print("Student ID: ")
		fmt.Scanln(&stdID)
		if stdID == "q" {
			break
		}
		fmt.Print("Student Name: ")
		fmt.Scanln(&stdName)
		fmt.Print("Student Surname: ")
		fmt.Scanln(&stdSurname)
		fmt.Print("Student Gender: ")
		fmt.Scanln(&stdGender)
		fmt.Print("Student Birth Year: ")
		fmt.Scanln(&stdBirthYear)

		data := map[string]string{
			"std_id":         stdID,
			"std_name":       stdName,
			"std_surname":    stdSurname,
			"std_gender":     stdGender,
			"std_birth_year": stdBirthYear,
		}

		transactions, _ := json.Marshal(data)
		blockchain.AddTransaction(string(transactions))

		var addAnother string
		fmt.Print("Add another student information to the current block? (y/n): ")
		fmt.Scanln(&addAnother)
		if addAnother == "n" {
			blockchain.CreateNewBlock()
			fmt.Println("New block created.")
			fmt.Println()
		}
	}

	for _, block := range blockchain.chain {
		fmt.Println("Block #:", block.index)
		fmt.Println("Timestamp:", block.timestamp)
		fmt.Println("Transactions:", block.transactions)
		fmt.Println("Previous Hash:", block.previousHash)
		fmt.Println("Hash:", block.hash)
		fmt.Println("Merkle Root:", block.merkleRoot)
		fmt.Println("Nonce:", block.nonce)
		fmt.Println("--------------------------------------")
	}

	fmt.Println("Blockchain Validation:", blockchain.IsChainValid())
}
