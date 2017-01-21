package bls

import (
	"strconv"
	"math/big"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

/// Rand
const RandLength = 32 
type Rand [RandLength]byte

// Constructors
func RandFromBytes(b []byte) (r Rand) {
	h := crypto.Keccak256Hash(b)
        copy(r[:RandLength], h[:])
	return
}

// Getters
func (r Rand) Bytes() []byte {
	return r[:]
}

func (r Rand) String() string {
	return string(r[:])
}

// Derived Randomness hierarchically
func (r Rand) DerivedRand(idx []byte) Rand {
	// Keccak is not susceptible to length-extension-attacks, so we can use it as-is to implement an HMAC
	return RandFromBytes(crypto.Keccak256(r.Bytes(), idx))
}

// Shortcuts to the derivation function
// ... by string
func (seed Rand) Ders(s ...string) Rand {
	r := seed
	for _, si := range s {
		r = r.DerivedRand([]byte(si))
	}
	return r
}

// ... by int
func (seed Rand) Deri(i int) Rand {
	return seed.Ders(strconv.Itoa(i))
}

// Convert to a random integer from the interval [0,n-1]. 
func (r Rand) Modulo(n int) int {
	// modulo len(groups) with big.Ints (Mod method works on pointers)
	var b big.Int
	b.Mod(common.Bytes2Big(r.Bytes()), big.NewInt(int64(n)))
        return int(b.Int64())
}

// Convert to a random permutation
func (r Rand) RandomPerm(n int, k int) []int {
	// modulo len(groups) with big.Ints (Mod method works on pointers)
	l := make([]int, n)
	for i,_ := range l {
		l[i] = i
	}
	for i := 0; i < k; i++ {
		j := r.Deri(i).Modulo(n-i) + i
		l[i], l[j] = l[j], l[i]
	}
	return l[:k]	
}
