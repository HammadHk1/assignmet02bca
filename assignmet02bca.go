// package maain

package assignment02bca

import (
	"bufio"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"math/rand"

	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

// bootStraping
type BootStrapNode struct {
	port       string
	ip_address string
}

var nodes []BootStrapNode

// func main() {
// 	var i int

// 	go server()

// 	var blocks = []string{"A", "B", "C", "D", "E"}

// 	for i = 0; i < 5; i += 1 {
// 		time.Sleep(2 * time.Second)
// 		client(blocks[i])
// 	}

// 	// add new block
// 	addnewNode("N")
// }

func addnewNode(a string) {
	var hashtable = []Hashable{}
	var i int
	var newBlock = Block(a)
	hashtable = append(hashtable, newBlock)
	newBlock.mine(0)
	printTree(buildTree(hashtable)[0].(Node))

	println("Perform BootStraping........")
	println("Add new block in Blockchain")
	println()

	println("Available Nodes to be connected ")
	for i = 0; i < 5; i += 1 {
		println("Portno: ", nodes[i].port, " ip-address: ", nodes[i].ip_address)

	}
}
func Client(a string) {
	println("Client-> ", a, "connected")
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		//handle error
	}
	var recvdBlock Block
	dec := gob.NewDecoder(conn)
	err = dec.Decode(&recvdBlock)
	if err != nil {

		//handle error

	}
	fmt.Println(recvdBlock.hash())
	min := 5000
	max := 6000

	var node BootStrapNode
	node.port = strconv.Itoa(rand.Intn(max-min) + min)
	node.ip_address = "127.0.0.1"
	nodes = append(nodes, node)

}
func Server() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn)

	}

}
func handleConnection(c net.Conn) {

	println()
	log.Println("A client has connected", c.RemoteAddr())
	println()

	fmt.Print("Tree 1:\n")
	var i int

	// hashtable(Blochain) is declared and initialized
	var hashtable = []Hashable{}
	var blocks = []string{"A", "B", "C", "D", "E"}

	for i = 0; i < 5; i += 1 {

		var newBlock = Block(blocks[i])
		hashtable = append(hashtable, newBlock)
		newBlock.mine(0)

	}
	printTree(buildTree(hashtable)[0].(Node))

	gobEncoder := gob.NewEncoder(c)
	err := gobEncoder.Encode(buildTree(hashtable)[0].hash())
	if err != nil {

		log.Println(err)

	}

}

// func main() {
// 	fmt.Print("Tree 1:\n")
// 	var i int

// 	// hashtable(Blochain) is declared and initialized
// 	var hashtable = []Hashable{}
// 	var blocks = []string{"A", "B", "C", "D", "E"}

// 	for i = 0; i < 5; i += 1 {

// 		var newBlock = Block(blocks[i])
// 		hashtable = append(hashtable, newBlock)
// 		newBlock.mine(0)

// 	}
// 	printTree(buildTree(hashtable)[0].(Node))

// }

func (b *Block) mine(difficulty int) {
	for !strings.HasPrefix(b.hash().String(), strings.Repeat("0", difficulty)) {
		CalculateHash(b.hash().String())
	}
}

func buildTree(blocks []Hashable) []Hashable {
	var nodes []Hashable
	var i int
	for i = 0; i < len(blocks); i += 2 {
		if i+1 < len(blocks) {
			nodes = append(nodes, Node{left: blocks[i], right: blocks[i+1]})
		} else {
			nodes = append(nodes, Node{left: blocks[i], right: EmptyBlock{}})
		}
	}
	if len(nodes) > 1 {
		return buildTree(nodes)
	} else if len(nodes) == 1 {
		return nodes
	} else {
		panic(" ")
	}
}

type Block string

func (b Block) hash() Hash {
	return hash([]byte(b)[:])
}

type Hashable interface {
	hash() Hash
}

type Hash [20]byte

func (h Hash) String() string {
	return hex.EncodeToString(h[:])
}

type EmptyBlock struct {
}

func (_ EmptyBlock) hash() Hash {
	return [20]byte{}
}

type Node struct {
	left  Hashable
	right Hashable
}

func (n Node) hash() Hash {
	var l, r [sha1.Size]byte
	l = n.left.hash()
	r = n.right.hash()
	return hash(append(l[:], r[:]...))
}

func hash(data []byte) Hash {
	return sha1.Sum(data)
}

func printTree(node Node) {
	printNode(node, 0)
}

func printNode(node Node, level int) {
	fmt.Printf("(%d) %s %s\n", level, strings.Repeat(" ", level), node.hash())
	if l, ok := node.left.(Node); ok {
		printNode(l, level+1)
	} else if l, ok := node.left.(Block); ok {
		fmt.Printf("(%d) %s %s (data: %s)\n", level+1, strings.Repeat(" ", level+1), l.hash(), l)
	}
	if r, ok := node.right.(Node); ok {
		printNode(r, level+1)
	} else if r, ok := node.right.(Block); ok {
		fmt.Printf("(%d) %s %s (data: %s)\n", level+1, strings.Repeat(" ", level+1), r.hash(), r)
	}
}

type block struct {
	x           int
	hash        string
	prev_hash   string
	transaction string
}
type BlockChain struct {
	//list of blocks
	list []*block
}

func NewBlock(transaction string, nonce int, previousHash string, blockchain *BlockChain) *block {
	block1 := new(block)
	block1.transaction = transaction
	block1.x = nonce
	block1.prev_hash = previousHash
	block1.hash = CalculateHash(block1.transaction + strconv.Itoa(block1.x) + block1.prev_hash)
	blockchain.list = append(blockchain.list, block1)
	return block1
}

func DisplayBlocks(blockchain *BlockChain) {
	for i, a := range blockchain.list {
		fmt.Printf("%s BLOCK %d %s\n", strings.Repeat("=", 25), i+1, strings.Repeat("=", 25))
		fmt.Printf(" TRANSACTION: %s \n NONCE VALUE: %d \n HASH OF PREVIOUS BLOCK : %s \n HASH OF CURRENT BLOCK %s \n \n ", a.transaction, a.x, a.prev_hash, a.hash)
	}

}

// changeBlock will get blochain and index of the block
func ChangeBlock(blockchain *BlockChain, index int) {

	var chainLength int
	chainLength = len(blockchain.list)
	if index < chainLength {

		fmt.Println("Your current transaction is as follows \n")
		fmt.Printf("%s", blockchain.list[index].transaction)
		scan := bufio.NewScanner(os.Stdin)
		fmt.Println("enter new transaction: \n")
		scan.Scan()
		text := scan.Text()
		blockchain.list[index].transaction = text
		fmt.Println("Block details are changed")

	}
}

func VerifyChain(blockchain *BlockChain) {

	var verify = false
	for _, num := range blockchain.list {

		Hash := CalculateHash(num.transaction + strconv.Itoa(num.x) + num.prev_hash)
		if Hash != num.hash {
			verify = true
			break
		}

	}

	if verify == false {
		fmt.Println("verification complete, no changes detected")
	} else {
		fmt.Println("change detected")
	}
}

func CalculateHash(stringToHash string) string {

	return fmt.Sprintf("%x", sha256.Sum256([]byte(stringToHash)))
}
