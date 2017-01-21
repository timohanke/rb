package sim

import (
	"fmt"
	"bls"
	"state"
	"github.com/ethereum/go-ethereum/common"
//	"github.com/davecgh/go-spew/spew"
)

/// simulator section
// Process
type ProcessSimulator struct {
	sec     bls.Seckey
	reginfo state.Node
	rseed   bls.Rand
	// rseed is the seed used for the internal randomness of the process, it did not seed the secret key  
	shares_source    map[common.Address]bls.SeckeyMap
	// TODO change to bls.SeckeyMap
	shares_combined  map[common.Address]bls.Seckey
}

func NewProcessSimulator(sec bls.Seckey, seed bls.Rand) (p ProcessSimulator) {
	p.sec = sec
	p.reginfo = state.NodeFromSeckey(sec)
	p.rseed = seed
	p.shares_source   = make(map[common.Address]bls.SeckeyMap)
	p.shares_combined = make(map[common.Address]bls.Seckey)
	return 
}

/* Makes the process simulator DETERMINISTIC by setting rseed to the process' own address.
   Since the address is public the process' behaviour becomes predictable from the outside.
   This will benefit testing. */
func NewProcessSimulatorDet(sec bls.Seckey) ProcessSimulator{
	node := state.NodeFromSeckey(sec)
	// assign temporary variable to make the value addressable
	tmp := node.Address() 
	return NewProcessSimulator(sec, bls.RandFromBytes(tmp[:]))
}

func (p *ProcessSimulator) Address() common.Address {
	return p.reginfo.Address()
}
			
func (p *ProcessSimulator) SetGroupShare(g state.Group, source common.Address, share bls.Seckey, vvec []bls.Pubkey) {
	addr := g.Address()
//	fmt.Printf("Setting source share: (proc)%.4x (grp)%.2x (src)%.4x (sec)%.4s\n", p.Address(), addr, source, share.String())
	// verify share
	if Vvec {
		if bls.SharePubkey(vvec, p.Address().Big()).String() != bls.PubkeyFromSeckey(share).String() {
			fmt.Println("Error: Received secret share does not match committed verification vector")
		}
	}

	// if key source does not exist yet then make a bls.SeckeyMap
	_, exists := p.shares_source[addr]
	if !exists {
	 	p.shares_source[addr] = bls.SeckeyMap{}
	}
	// store source share
	p.shares_source[addr][source] = share
	return
}

func (p *ProcessSimulator) AggregateGroupShares(g state.Group) {
	addr := g.Address()
	vlist := make([]bls.Seckey, len(p.shares_source[addr]))
	i := 0
	for _, sec := range p.shares_source[addr] {
		vlist[i] = sec
		i++
	}
	p.shares_combined[addr] = bls.AggregateSeckeys(vlist)
	return
}

// temp: 
func (p *ProcessSimulator) GetAggregatedGroupShare(g state.Group) bls.Seckey {
	return p.shares_combined[g.Address()]
}

func (p *ProcessSimulator) GetSeckeyForGroup(g state.Group) (sec bls.Seckey) {
	// provide own secret for the group setup (function of internal seed and group address)
	addr := g.Address()
	gseed := p.rseed.DerivedRand(addr[:])
	sec = bls.SeckeyFromRand(gseed.Deri(0))
//	fmt.Printf("sec for group: %s\n", sec.String())
	return 
}

func (p *ProcessSimulator) GetSeckeySharesForGroup(g state.Group) (bls.SeckeyMap, []bls.Pubkey) {
	// take own secret for the group setup (function of internal seed and group address) and split it up in shares for all group members
	// from the process seed (rseed) and derive a per-group seed based on the group's address
	addr := g.Address()
	gseed := p.rseed.DerivedRand(addr[:])
	// from the per-group seed derive a vector of k seckeys as the master seckey where k is the threshold 
	// the master seckey defines a polynomial of degree k-1
	k := g.Threshold()
	msec := make([]bls.Seckey, k)
	vvec := make([]bls.Pubkey, k)
	for i:=0; i<int(k); i++ {
		msec[i] = bls.SeckeyFromRand(gseed.Deri(i))
		vvec[i] = bls.PubkeyFromSeckey(msec[i])
	}
	shares := bls.SeckeyMap{}
	for _, m := range g.Members() {
		shares[m] = bls.ShareSeckeyByAddr(msec, &m)
	}
	return shares, vvec
}

func (p *ProcessSimulator) SignForGroup(g state.Group, msg []byte) bls.Signature {
	sec := p.shares_combined[g.Address()]
//	fmt.Printf("sign for group: (grp)%.2x (sec)%x\n", g.Address(), sec.String())
	return bls.Sign(sec, msg)
}

func (p *ProcessSimulator) Sign(msg []byte) bls.Signature {
	return bls.Sign(p.sec, msg)
}

func (p *ProcessSimulator) Log() {
	fmt.Printf("Process simulator: % x\n", p.reginfo.Address())
	fmt.Printf("  Seckey: % x\n", p.sec.Bytes())
	fmt.Printf("  rseed: % x\n", p.rseed)
	p.reginfo.Log()
}

func (p *ProcessSimulator) String() string {
	return fmt.Sprintf("Proc: (sec)%s (seed)%x %s", p.sec.String()[:4], p.rseed.String()[:2], p.reginfo.String())
}
