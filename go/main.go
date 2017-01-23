package main

// #include <test.h>
// int fortytwo() { return 42; }
//import "C"

import (
	"fmt"
	"bls"
	"sim"
	"blscgo"
//	"github.com/davecgh/go-spew/spew"
)

func main() {
	// init Cgo
	blscgo.Init()

	// select seed
	
	// seed, groupSize, threshold, nProcesses, nGroups
	// check this out for a bug at #37:
	mysim := sim.NewBlockchainSimulator(bls.RandFromBytes([]byte("XXX")), 3, 2, 5, 2)
	l := 20 
	fmt.Println("--- Genesis block ")
	fmt.Printf("%d: %s",mysim.Length(),mysim.Tip().String(true))
	fmt.Printf("--- Blockchain states: (l)%d\n", l)
	for i:=0; i<l; i++ {
		mysim.Advance(1, false)
		fmt.Printf("%3d: %s\n",mysim.Length(),mysim.Tip().String(false))
	}
	bls.PrintCtrs()
	fmt.Println("--- Info")
	fmt.Println("Seckey calls should be:    m*n/m*n^2, m*n^2/m*n^2*k") 
	fmt.Println("Pubkey calls should be:    N+m*n+m*n^2, m*n^2/m*n^2*k, m/m*n") 
	// pubkey generation: N is process generation, m*n is vvec generation, m*n^2 is rhs of vvec verification
	// pubkey sharing: m*n^2/m*n^2*k is lhs of vvec verification 
	// pubkey aggregation: m/m*n is generation of group pubkey from member shares
	fmt.Println("Signature calls should be: N+l*n, N, l/l*k") 

//	mysim.Log()
/*
	// test generation
	sec :=	bls.SeckeyFromRand(bls.RandFromBytes([]byte("seed")))
	fmt.Println(sec.Bytes())
	spew.Dump(sec)
	pub := bls.PubkeyFromSeckey(sec)
	spew.Dump(pub)

	// test sig + verify
	sig := bls.Sign(sec, []byte("hi"))
	spew.Dump(sig)
	val := bls.VerifySig(pub, []byte("hi"), sig)
	spew.Dump(val)

	// test share derivations for seckeys against the same derivation for pubkeys
	pub1 := bls.SharePubkey([]bls.Pubkey{pub, pub}, big.NewInt(1))
	spew.Dump(pub1)
	pub2 := bls.SharePubkey([]bls.Pubkey{pub, pub}, big.NewInt(2))
	spew.Dump(pub2)

	sec1 := bls.ShareSeckey([]bls.Seckey{sec, sec}, big.NewInt(1))
	sec2 := bls.ShareSeckey([]bls.Seckey{sec, sec}, big.NewInt(2))
	sec1pub := bls.PubkeyFromSeckey(sec1)
	spew.Dump(sec1pub)
	sec2pub := bls.PubkeyFromSeckey(sec2)
	spew.Dump(sec2pub)
	if pub1.String() != sec1pub.String() {
		fmt.Println("Error: sec1pub does not match.")
	}
	if pub2.String() != sec2pub.String() {
		fmt.Println("Error: sec2pub does not match.")
	}

	// test seckey recovery from shares
	recovered := bls.RecoverSeckey([]bls.Seckey{sec1, sec2}, []*big.Int{big.NewInt(1), big.NewInt(2)})
	spew.Dump(recovered)
	if recovered.String() != sec.String() {
		fmt.Println("Error: recovered seckey does not match.")
	} else {
		fmt.Println("Ok: recovered seckey matches.")
	}
	
	// test VerifySig again
//	val = bls.VerifySig(sec1pub, []byte("hi"), sig)
//	spew.Dump(val)

	sum := bls.AggregateSeckeys([]bls.Seckey{sec,sec})
	spew.Dump(sum)
	if sec1.String() != sum.String() {
		fmt.Println("Error: sec1 does not aggregate seckey.")
	}

	// test aggregate pubkeys here
	twopub := bls.AggregatePubkeys([]bls.Pubkey{pub,pub})
	sumpub := bls.PubkeyFromSeckey(sum)
	if twopub.String() != sumpub.String() {
		fmt.Println("Error: aggregated pubkey does not match.")
	}

	// aggregate sig
	twosig := bls.AggregateSigs([]bls.Signature{sig,sig})
	val = bls.VerifySig(twopub, []byte("hi"), twosig)
	spew.Dump("aggregate sig: ",val)

	// recover sig
	sig1 := bls.Sign(sec1, []byte("hi"))
	spew.Dump("sig1: ", sig1)
	val = bls.VerifySig(pub1, []byte("hi"), sig1)
	if !val {
		fmt.Println("Error: sig1 did not verify.")
	}
	sig2 := bls.Sign(sec2, []byte("hi"))
	spew.Dump("sig2: ", sig2)
	val = bls.VerifySig(pub2, []byte("hi"), sig2)
	if !val {
		fmt.Println("Error: sig2 did not verify.")
	}
	sig0 := bls.RecoverSignature([]bls.Signature{sig1, sig2}, []*big.Int{big.NewInt(1), big.NewInt(2)})
	val = bls.VerifySig(pub, []byte("hi"), sig0)
	if !val {
		fmt.Println("Error: sig0 did not verify.")
	} else {
		fmt.Println("Ok: sig0 verifies.")
	}
	spew.Dump(sig0)
	if sig0.String() != sig.String() {
		fmt.Println("Error: recovered sig does not match.")
	}
*/
}
