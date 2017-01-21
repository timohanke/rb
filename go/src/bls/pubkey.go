package bls

/* TODO
 - reduce dependency on ethereum code
*/

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
var pub_gen_calls, pub_agg_calls, pub_agg_len, pub_share_calls, pub_share_len int

func PubkeyCtrs() string {
	return fmt.Sprintf("(pub:gen,shr,agg) %d,%d/%d,%d/%d", pub_gen_calls, pub_share_calls, pub_share_len, pub_agg_calls, pub_agg_len)
}
//var pubkey_ctr uint16 = 0

// types
type Pubkey struct {
//	trace       []byte
	value       []byte
//	y, x1, x2   big.Int
}

type PubkeyMap map[common.Address]Pubkey

// hash & id
func (pub Pubkey) Hash() common.Hash {
	return crypto.Keccak256Hash(pub.value)
}

func (pub Pubkey) Address() common.Address {
	h := pub.Hash()
	return common.BytesToAddress(h[:])
//	return Trace2Address(pub.trace)
//	pubBytes := []byte("pubkey")
}

func (pub Pubkey) String() string {
	return string(pub.value)
//	a := pub.Address()
//	return fmt.Sprintf("%x", pub.Address())
}

// Generation
func PubkeyFromSeckey(sec Seckey) (pub Pubkey) {
	pub_gen_calls += 1
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

// Aggregate multiple into one by summing up
func AggregatePubkeys(pubs []Pubkey) (pub Pubkey) {
	pub_agg_calls += 1
	pub_agg_len += len(pubs)
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

// Derive shares from master through polynomial substitution 
func SharePubkey(mpub []Pubkey, i *big.Int) (pub Pubkey) {
	pub_share_calls += 1
	pub_share_len += len(mpub)
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

// Recover master from shares through Lagrange interpolation
func RecoverPubkey(pubs []Pubkey, ids []big.Int) Pubkey {
	return Pubkey{}
}
