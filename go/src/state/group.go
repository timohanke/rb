package state

import (
	"bls"
	"common2"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"log"
)

// Group -- encodes all data of a group as recorded on the blockchain
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

// NewGroup -- create a new Group struct with list of members and empty pubkey
// NewGroup --
func NewGroup(addresses []common.Address, k uint16) Group {
	return Group{addresses, bls.Pubkey{}, k}
}

// SetPubkey -- set the group's pubkey and threshold
// SetPubkey --
func (g *Group) SetPubkey(pub bls.Pubkey, k uint16) {
	g.pub = pub
	g.threshold = k
}

// Address - return the group address 
// Address --
func (g Group) Address() (a common.Address) {
	// hash of all member addresses
	d := sha3.NewKeccak256()
	var err error
        for _, addr := range common2.SortAddresses(g.members) {
                _, err = d.Write(addr[:])
		if err != nil { log.Fatalln("Error when calling Keccak256") }
        }

	var h common.Hash
        d.Sum(h[:0])
	return common.BytesToAddress(h[:]) 
}

// Pubkey -- 
// Pubkey --
func (g Group) Pubkey() bls.Pubkey {
	return g.pub
}

// Members --
// Members --
func (g Group) Members() []common.Address {
	return g.members
}

// Threshold --
// Threshold --
func (g Group) Threshold() int {
	return int(g.threshold)
}

// Size --
// Size --
func (g Group) Size() int {
	return len(g.members)
}

// Log --
// Log --
func (g Group) Log() {
	fmt.Println("    members: ", len(g.members))
	for _, m := range g.members {
		fmt.Printf("      address: % x\n", m)
	}
	fmt.Printf("    addr: % x\n", g.pub.Address())
	fmt.Println("    threshold: ", g.threshold)
}

// String --
// String --
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

// isValid --
/* LATER: check if pubkey is correctly signed by signatories */
// isValid --
func (g Group) isValid() bool {
	return true
}
