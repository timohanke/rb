package state

import (
	"fmt"
	"bls"
	"github.com/ethereum/go-ethereum/common"
)

type Node struct {
	pub bls.Pubkey
	pop bls.Pop
	// LATER: proof of creation (stake, etc.)
}

func NodeFromSeckey(sec bls.Seckey) Node {
	pub := bls.PubkeyFromSeckey(sec)
	return Node{pub, bls.GeneratePop(sec,pub)}
}

func (n Node) Address() common.Address {
	return n.pub.Address()
}

func (n Node) hasPop() bool {
	// LATER: check proof of creation 	
	return bls.VerifyPop(n.pub, n.pop)
}

func (n Node) VerifySigned(r bls.Rand, sig bls.Signature) bool {
	// LATER: Verify signed group key (Pubkey), master pubkey ([]Pubkey), secret key share (Seckey) against this node 
	return true
}

func (n Node) Log() {
	fmt.Printf("    pub: % x\n", n.pub.Address())
//	fmt.Printf("  Seckey: % x\n", p.sec.Bytes())
	fmt.Println("    pop: ", n.pop)
}

func (n Node) String() string {
	a := n.pub.Address()
	return fmt.Sprintf("Node: (addr)%x (pub)%s", string(a[:2]), n.pub.String()[:8])
}
