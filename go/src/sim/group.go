package sim

import (
	"time"
	"fmt"
	"bls"
	"state"
	"github.com/ethereum/go-ethereum/common"
//	"github.com/davecgh/go-spew/spew"
)

// Group
type GroupSimulator struct {
	// TODO remove sec
	sec      bls.Seckey
	reginfo  state.Group
	proclist []*ProcessSimulator
	procmap  map[common.Address]*ProcessSimulator
}

func NewGroupSimulator(members []*ProcessSimulator, k uint16) GroupSimulator {
	m := len(members)
	// collect all members' addresses in a Group struct with empty Pubkey
	addresses := make([]common.Address, m) 
	pmap := make(map[common.Address]*ProcessSimulator)
	for i, p := range members {
		addresses[i] = p.Address()
		pmap[p.Address()] = p
	}
	g := state.NewGroup(addresses, k)

	// optional double-check of the group secret
	var secs []bls.Seckey
	if Double_check {
		secs = make([]bls.Seckey, m) 
	}

	pubs := make([]bls.Pubkey, m)
	// get all members' contribution to the group secret
	for i, p := range members {
		// get secret shares for all other processes
		shares, vvec := p.GetSeckeySharesForGroup(g)
		pubs[i] = vvec[0]
		// send shares out to all other individual processes
		for _, q := range members {
			q.SetGroupShare(g, p.Address(), shares[q.Address()], vvec)
		}
		// optional double-check of the group secret
		if Double_check {
			secs[i] = p.GetSeckeyForGroup(g)
			recovered := bls.RecoverSeckeyByMap(shares, int(k))
			if secs[i].String() != recovered.String() {
				fmt.Println("Error: recovered seckey share (ByMap) does not match.")
			} 
		}
	}

	// build group pubkey
	pub := bls.AggregatePubkeys(pubs)
	
        // set group pubkey in Group struct
	g.SetPubkey(pub, k)

	// tell each process to aggregate their shares
	// processes need their aggregated shares for signing later
	for _, q := range members {
		q.AggregateGroupShares(g)
	}

	if Double_check {
		// fetch the combined shares from each process into a SeckeyMap
		// (every process does this so we wouldn't need to in the group simulator)
		agg_shares := bls.SeckeyMap{}
		for _, p := range members {
			agg_shares[p.Address()] = p.GetAggregatedGroupShare(g)
		}

		// recover the combined group secret from combined shares 
		// choose k random shares, combine and compare
		sec := bls.RecoverSeckeyByMap(agg_shares, int(k))
		pub_dup := bls.PubkeyFromSeckey(sec)

		// optional double-check: aggregate all contributions into the group secret and compare
		sec_dup := bls.AggregateSeckeys(secs)
		if sec.String() != sec_dup.String() {
			fmt.Println("Error: recovered aggregated seckey does not match.")
		}

		if pub.String() != pub_dup.String() {
			fmt.Println("Error: recovered aggregated pubkey does not match.")
		}
	}
 
	// TODO remove seckey
	return GroupSimulator{bls.Seckey{}, g, members, pmap}
}

func (g GroupSimulator) Sign(msg []byte) bls.Signature {
	sigmap := make(map[common.Address]bls.Signature)
	// get signature share from each process
	t0 := time.Now()
	for _, p := range g.proclist {
		sigmap[p.Address()] = p.SignForGroup(g.reginfo, msg)
	}
	delta1 := time.Since(t0)
	t1 := time.Now()
	sig1 := bls.RecoverSignatureByMap(sigmap, g.reginfo.Threshold())
	delta2 := time.Since(t1)
	fmt.Printf("Time for group signatures with %d shares: %v (%vms / share) + %v (recovery).\n", len(g.proclist), delta1, (delta1.Nanoseconds()/1000000)/int64(len(g.proclist)), delta2)
 
	// optional verification
	/*
	sig2 := bls.Sign(g.sec, msg)
	if sig1.String() != sig2.String(
		fmt.Println("Error in Group sign: Recovered signature does not match.")
	}
        */
	return sig1
}

func (g GroupSimulator) Address() common.Address {
	return g.reginfo.Address()
}

func (g GroupSimulator) Log() {
	fmt.Printf("Group simulator: % x\n", g.reginfo.Address())
	fmt.Printf("  Seckey: % x\n", g.sec.Bytes())
	g.reginfo.Log()
}

func (g *GroupSimulator) String() string {
	return fmt.Sprintf("GrpP: (sec)%s %s", g.sec.String()[:4], g.reginfo.String())
}
