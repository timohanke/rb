package bls

import (
	"fmt"
	"math/big"
	"strings"
	"strconv"
	"bytes"
	"os/exec"
	"github.com/ethereum/go-ethereum/common"
	"encoding/base64"
)

// Cgo test
/* Trial to link C++, didn't work  
// #include <bls/include/test.hpp>

// #cgo LDFLAGS: -L/root/beacon/go/src/bls/bls2/lib -L/root/beacon/go/src/bls/mcl/lib  -lbls -lmcl
// #include <bls2/include/bls.h>
import "C"

// CTest --
func CTest() {
	C.FooInit()
	return
}
*/

// Debugging counters
var sigGenCalls, sigVerifyCalls, sigAggCalls, sigAggLen, sigRecoverCalls, sigRecoverLen int

// SignatureCtrs --
func SignatureCtrs() string {
	return fmt.Sprintf("(sig:gen,ver,rec) %d,%d,%d/%d", sigGenCalls, sigVerifyCalls, sigRecoverCalls, sigRecoverLen)
}

/// Crypto 
// types

// Signature --
type Signature struct {
	value []byte
}

// SignatureMap --
type SignatureMap map[common.Address]Signature

// Conversion

// Rand --
func (sig Signature) Rand() Rand {
	return RandFromBytes(sig.value)
}

// String --
func (sig Signature) String() string {
	return string(sig.value)
}

// Sign --
// Signing 
func Sign(sec Seckey, msg []byte) (sig Signature) {
	sigGenCalls++
	// call bls_tool
	stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
	cmd := exec.Command("bls_tool.exe","sign")
	str := sec.String() + "\n" + base64.StdEncoding.EncodeToString(msg)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = strings.NewReader(str), stdout, stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("bls_tool.exe sign: %v\n%s", err, string(stderr.Bytes()))
	}
	sig.value = bytes.TrimRight(stdout.Bytes(),"\n")
	return 
}

// VerifySig --
// Verification
func VerifySig(pub Pubkey, msg []byte, sig Signature) bool {
	sigVerifyCalls++
	// call bls_tool
	stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
	cmd := exec.Command("bls_tool.exe","verify")
	str := sig.String() + "\n" + pub.String() + "\n" + base64.StdEncoding.EncodeToString(msg)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = strings.NewReader(str), stdout, stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("bls_tool.exe verify: %v\n%s", err, string(stderr.Bytes()))
	}
	val, _ := strconv.Atoi(strings.TrimRight(stdout.String(),"\n")) 
	return val > 0 
}

// VerifyAggregateSig --
func VerifyAggregateSig(pubs []Pubkey, msg []byte, asig Signature) bool {
	return VerifySig(AggregatePubkeys(pubs), msg, asig)
}

// BatchVerify --
func BatchVerify(pubs []Pubkey, msg []byte, sigs []Signature) bool {
	return VerifyAggregateSig(pubs, msg, AggregateSigs(sigs)) 
}
	
// AggregateSigs --
// Aggregate multiple into one by summing up
func AggregateSigs(sigs []Signature) (sig Signature) {
	sigAggCalls++
	sigAggLen += len(sigs)
	// call bls_tool
	stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
	cmd := exec.Command("bls_tool.exe","aggregate-sig")
	str := strconv.Itoa(len(sigs)) 
	for _, s := range sigs {
		str += "\n" + s.String() 
	}
	cmd.Stdin, cmd.Stdout, cmd.Stderr = strings.NewReader(str), stdout, stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("bls_tool.exe aggregate-sig: %v\n%s", err, string(stderr.Bytes()))
	}
	sig.value = bytes.TrimRight(stdout.Bytes(),"\n")
	return 
}

// RecoverSignature --
// Recover master from shares through Lagrange interpolation
func RecoverSignature(sigs []Signature, ids []*big.Int) (sig Signature) {
	sigRecoverCalls++
	sigRecoverLen += len(sigs)
	// call bls_tool
	stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
	cmd := exec.Command("bls_tool.exe","recover-sig")
	str := strconv.Itoa(len(sigs)) 
	for i, s := range sigs {
		str += "\n" + ids[i].String() 
		str += "\n" + s.String() 
	}
	cmd.Stdin, cmd.Stdout, cmd.Stderr = strings.NewReader(str), stdout, stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("bls_tool.exe recover-sig: %v\n%s", err, string(stderr.Bytes()))
	}
	sig.value = bytes.TrimRight(stdout.Bytes(),"\n")
	return 
}

// RecoverSignatureByMap --
func RecoverSignatureByMap(m SignatureMap, k int) (sec Signature) {
	ids := make([]*big.Int, k)
	sigs := make([]Signature, k)
	i := 0
	for a, s := range m {
		ids[i] = a.Big()
		sigs[i] = s
		i++
		if i >= k {
			break
		}
	}
	return RecoverSignature(sigs, ids)
}
