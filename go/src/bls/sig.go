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

func CTest() {
	C.FooInit()
	return
}
*/

// Debugging counters
var sig_gen_calls, sig_verify_calls, sig_agg_calls, sig_agg_len, sig_recover_calls, sig_recover_len int

func SignatureCtrs() string {
	return fmt.Sprintf("(sig:gen,ver,rec) %d,%d,%d/%d", sig_gen_calls, sig_verify_calls, sig_recover_calls, sig_recover_len)
}

/// Crypto 
// types
type Signature struct {
	value []byte
}

type SignatureMap map[common.Address]Signature

// Conversion
func (sig Signature) Rand() Rand {
	return RandFromBytes(sig.value)
}

func (sig Signature) String() string {
	return string(sig.value)
}

// Signing 
func Sign(sec Seckey, msg []byte) (sig Signature) {
	sig_gen_calls += 1
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

// Verification
func VerifySig(pub Pubkey, msg []byte, sig Signature) bool {
	sig_verify_calls += 1
	// call bls_tool
	stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
	cmd := exec.Command("bls_tool.exe","verify")
	str := sig.String() + "\n" + pub.String() + "\n" + base64.StdEncoding.EncodeToString(msg)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = strings.NewReader(str), stdout, stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("bls_tool.exe verify: %v\n%s", err, string(stderr.Bytes()))
	}
	val, _ := strconv.Atoi(strings.TrimRight(stdout.String(),"\n")) 
	if val > 0 {
		return true
	} else {
		return false
	}
}

func VerifyAggregateSig(pubs []Pubkey, msg []byte, asig Signature) bool {
	return VerifySig(AggregatePubkeys(pubs), msg, asig)
}

func BatchVerify(pubs []Pubkey, msg []byte, sigs []Signature) bool {
	return VerifyAggregateSig(pubs, msg, AggregateSigs(sigs)) 
}
	
// Aggregate multiple into one by summing up
func AggregateSigs(sigs []Signature) (sig Signature) {
	sig_agg_calls += 1
	sig_agg_len += len(sigs)
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

// Recover master from shares through Lagrange interpolation
func RecoverSignature(sigs []Signature, ids []*big.Int) (sig Signature) {
	sig_recover_calls += 1
	sig_recover_len += len(sigs)
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
