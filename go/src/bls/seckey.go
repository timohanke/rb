package bls

import (
	"blscgo"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"log"
	"math/big"
)

/// Crypto 
// Logging counters
var sec_agg_calls, sec_agg_len, sec_share_calls, sec_share_len, sec_recover_calls, sec_recover_len int 

func SeckeyCtrs() string {
	return fmt.Sprintf("(sec:agg,shr)     %d/%d,%d/%d", sec_agg_calls, sec_agg_len, sec_share_calls, sec_share_len)
}

// Constants
var R big.Int = Decimal2Big("16798108731015832284940804142231733909759579603404752749028378864165570215949")

// types
type Seckey struct {
	secret *big.Int
}

type SeckeyMap map[common.Address]Seckey

func (sec Seckey) Bytes() []byte {
	// big endian
	return sec.secret.Bytes()
}

func (sec Seckey) String() string {
	// big endian
	return sec.secret.String()
}

func (sec Seckey) BigInt() *big.Int {
	return sec.secret
}

func (sec Seckey) Hex() string {
	return fmt.Sprintf("0x%x", sec.secret)
}

func (sec Seckey) SecretKey() (sk *blscgo.SecretKey) {
    sk = new(blscgo.SecretKey)
    err := sk.SetStr(sec.String())
    if err != nil { log.Fatalln("Error in SecretKey conversion from blscgo.") }
    return
}

// Constructors
func SeckeyFromBytes(b []byte) (sec Seckey) {
	// the secret has to be cut off at 31 bytes to make it smaller than the constant R
	// R has 254 bits
	// TODO mask only the two highest bits with zeros
	if len(b) > 31 {
		b = b[:31]
	}
	i := Bytes2Big(b)
	sec.secret = &i
	return
}
	
func SeckeyFromRand(seed Rand) Seckey {
	return SeckeyFromBytes(seed.Bytes()) 
}

func SeckeyFromBigInt(b *big.Int) (sec Seckey) {
	sec.secret = b
	return
}

func SeckeyFromInt(i int64) (sec Seckey) {
	sec.secret = big.NewInt(i)
	return
}
// Aggregate multiple seckeys into one by summing up, using native big.Ints
func AggregateSeckeys(secs []Seckey) (sec Seckey) {
	sec_agg_calls += 1
	sec_agg_len += len(secs)
	sec.secret = big.NewInt(0)
	for _, s := range secs {
		sec.secret.Add(sec.secret, s.secret)
	}
	sec.secret.Mod(sec.secret, &R)
	return 
}

// Derive shares from master through polynomial substitution 
// TODO make this function use PolynomialSubstitution
func ShareSeckey(msec []Seckey, x *big.Int) (sec Seckey) {
	sec_share_calls += 1
	sec_share_len += len(msec)
	sec.secret = big.NewInt(0)
	// degree of polynomial, need k >= 1, i.e. len(msec) >= 2
	k := len(msec)-1
	// msec = c_0, c_1, ..., c_k
	// evaluate polynomial f(x) with coefficients c0, ..., ck
	sec.secret.Set(msec[k].secret)
	for j:=k-1; j>=0; j-- {
		sec.secret.Mul(sec.secret, x)
		//sec.secret.Mod(&sec.secret, &R)
		sec.secret.Add(sec.secret, msec[j].secret)
		sec.secret.Mod(sec.secret, &R)
	}
	return 
}

func ShareSeckeyByAddr(msec []Seckey, addr *common.Address) (sec Seckey) {
	return ShareSeckey(msec, addr.Big())
}

// Recover master from shares through Lagrange interpolation
func RecoverSeckey(secs []Seckey, ids []*big.Int) (sec Seckey) {
	sec_recover_calls += 1
	sec_recover_len += len(secs)
	sec.secret = big.NewInt(0) 
	k := len(secs)
	// need len(ids) = k > 0 
	for i:=0; i<k; i++ {
		// compute delta_i depending on ids only
		var delta, num, den, diff *big.Int = big.NewInt(1), big.NewInt(1), big.NewInt(1), big.NewInt(0)
		for j:=0; j<k; j++ {
			if (j != i) {
				num.Mul(num, ids[j])
				num.Mod(num, &R)
				diff.Sub(ids[j], ids[i])
				den.Mul(den, diff)
				den.Mod(den, &R)
			}
		}
		// delta = num / den
		den.ModInverse(den, &R)
		delta.Mul(num, den)
		delta.Mod(delta, &R)
		// apply delta to secs[i] 
		delta.Mul(delta, secs[i].secret)
		// skip reducing delta modulo R here
		sec.secret.Add(sec.secret, delta)
		sec.secret.Mod(sec.secret, &R)
	}
	return 
}

func RecoverSeckeyByMap(m SeckeyMap, k int) (sec Seckey) {
	ids := make([]*big.Int, k)
	secs := make([]Seckey, k)
	i := 0
	for a, s := range m {
		ids[i] = a.Big()
		secs[i] = s
		i++
		if i >= k {
			break
		}
	}
	return RecoverSeckey(secs, ids)
}

/*
func RecoverSeckeyByAddr(secs []Seckey, addrs []common.Address) (sec Seckey) {
	ids := make([]*big.Int, len(addrs))
	for i, a := range addrs {
		ids[i] = a.Big()
	}
	return RecoverSeckey(secs, ids)
}
*/

