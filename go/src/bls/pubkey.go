package bls

import (
	"fmt"
	"math/big"
	"os/exec"
	"strings"
	"strconv"
	"bytes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

/// Crypto 
// Debugging counters
var pubGenCalls, pubAggCalls, pubAggLen, pubShareCalls, pubShareLen int

// PubkeyCtrs --
func PubkeyCtrs() string {
	return fmt.Sprintf("(pub:gen,shr,agg) %d,%d/%d,%d/%d", pubGenCalls, pubShareCalls, pubShareLen, pubAggCalls, pubAggLen)
}
//var pubkey_ctr uint16 = 0

// types

// Pubkey -
type Pubkey struct {
//	trace       []byte
	value       []byte
//	y, x1, x2   big.Int
}

// PubkeyMap --
type PubkeyMap map[common.Address]Pubkey

// Hash -- hash & id
func (pub Pubkey) Hash() common.Hash {
	return crypto.Keccak256Hash(pub.value)
}

// Address --
func (pub Pubkey) Address() common.Address {
	h := pub.Hash()
	return common.BytesToAddress(h[:])
//	return Trace2Address(pub.trace)
//	pubBytes := []byte("pubkey")
}

// String --
func (pub Pubkey) String() string {
	return string(pub.value)
//	a := pub.Address()
//	return fmt.Sprintf("%x", pub.Address())
}

// Generation

// PubkeyFromSeckey --
func PubkeyFromSeckey(sec Seckey) (pub Pubkey) {
	pubGenCalls++
//	pubkey_ctr++
//	pub.trace = append(sec.Bytes(),CtrBytes(pubkey_ctr)...)
	// call bls_tool
	stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
	cmd := exec.Command("bls_tool.exe","pubkey")
	cmd.Stdin, cmd.Stdout, cmd.Stderr = strings.NewReader(sec.String()), stdout, stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("bls_tool.exe pubkey: %v\n%s", err, string(stderr.Bytes()))
	}
	pub.value = bytes.TrimRight(stdout.Bytes(),"\n")
	return 
}

// AggregatePubkeys --
// Aggregate multiple into one by summing up
func AggregatePubkeys(pubs []Pubkey) (pub Pubkey) {
	pubAggCalls++
	pubAggLen += len(pubs)
	// call bls_tool
	stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
	cmd := exec.Command("bls_tool.exe","aggregate-pub")
	str := strconv.Itoa(len(pubs)) 
	for _, p := range pubs {
		str += "\n" + p.String() 
	}
	cmd.Stdin, cmd.Stdout, cmd.Stderr = strings.NewReader(str), stdout, stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("bls_tool.exe aggregate-pub: %v\n%s", err, string(stderr.Bytes()))
	}
	pub.value = bytes.TrimRight(stdout.Bytes(),"\n")
	return 
}

// SharePubkey --
// Derive shares from master through polynomial substitution 
func SharePubkey(mpub []Pubkey, i *big.Int) (pub Pubkey) {
	pubShareCalls++
	pubShareLen += len(mpub)
	// call bls_tool
	stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
	cmd := exec.Command("bls_tool.exe","share-pub")
	str := strconv.Itoa(len(mpub)) + "\n"
	for _, p := range mpub {
		str += p.String() + "\n"
	}
	str += i.String()
	cmd.Stdin, cmd.Stdout, cmd.Stderr = strings.NewReader(str), stdout, stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("bls_tool.exe share-pub: %v\n%s", err, string(stderr.Bytes()))
	}
	pub.value = bytes.TrimRight(stdout.Bytes(),"\n")
	return 
}

// RecoverPubkey --
// Recover master from shares through Lagrange interpolation
func RecoverPubkey(pubs []Pubkey, ids []big.Int) Pubkey {
	return Pubkey{}
}
