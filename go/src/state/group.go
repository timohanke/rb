package state

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"bls"
	"common2"
)

type Group struct {
	members    []common.Address
	// group pubkey
	pub        bls.Pubkey
	threshold  uint16

	// LATER:
	// signatories []Address // signed off the group pubkey, members active in the group setup
	// sig Signature // aggregated signature under the pubkey
	// status: inactive -> setup -> registered/active
	// registration time
}

// Create a new Group struct with list of members and empty pubkey
func NewGroup(addresses []common.Address, k uint16) Group {
	return Group{addresses, bls.Pubkey{}, k}
}

func (g *Group) SetPubkey(pub bls.Pubkey, k uint16) {
	g.pub = pub
	g.threshold = k
}

func (g Group) Address() (a common.Address) {
	// hash of all member addresses
	d := sha3.NewKeccak256()
        for _, addr := range common2.SortAddresses(g.members) {
                d.Write(addr[:])
        }

	var h common.Hash
        d.Sum(h[:0])
	return common.BytesToAddress(h[:]) 
}

func (g Group) Pubkey() bls.Pubkey {
	return g.pub
}

func (g Group) Members() []common.Address {
	return g.members
}

func (g Group) Threshold() int {
	return int(g.threshold)
}

func (g Group) Size() int {
	return len(g.members)
}

func (g Group) Log() {
	fmt.Println("    members: ", len(g.members))
	for _, m := range g.members {
		fmt.Printf("      address: % x\n", m)
	}
	fmt.Printf("    addr: % x\n", g.pub.Address())
	fmt.Println("    threshold: ", g.threshold)
}
	
func (g Group) String() string {
	a := g.Address()
	mem := "["
	for i, m := range g.members {
		if i>0 {
			mem += ","
		}
		mem += fmt.Sprintf("%x", string(m[:2]))
	}
	mem += "]"
	return fmt.Sprintf("GrpR: (addr)%x (pub)%.8s (n)%d (k)%d (mem)%s", a[:2], g.pub.String(), len(g.members), g.threshold, mem) 
}


/* LATER: check if pubkey is correctly signed by signatories */
func (g Group) isValid() bool {
	return true
}
