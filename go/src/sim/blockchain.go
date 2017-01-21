package sim

import (
	"fmt"
	"bls"
	"state"
	"github.com/ethereum/go-ethereum/common"
//	"github.com/davecgh/go-spew/spew"
)

type BlockchainSimulator struct {
	groupSize  uint16
	threshold  uint16
	seed       bls.Rand
	proc    []ProcessSimulator
	group   []GroupSimulator
	grpmap  map[common.Address]*GroupSimulator
	chain   []state.State
}

// optional double-check
var Double_check bool = false
var Vvec         bool = false

// is this function needed?
/*
func (sim *BlockchainSimulator) Init(seed bls.Rand, groupSize uint16, threshold uint16) {
	// leave the arrays empty
	sim = &BlockchainSimulator{seed: seed, groupSize: groupSize, threshold: threshold}
}
*/

func (sim *BlockchainSimulator) InitProcs(n uint) {
	sim.proc = make([]ProcessSimulator, n)
	rsec := sim.seed.Ders("InitProcs_sec")
	rseed := sim.seed.Ders("InitProcs_seed")
	for i := 0; i < int(n); i++ {
		sim.proc[i] = NewProcessSimulator(bls.SeckeyFromRand(rsec.Deri(i)), rseed.Deri(i))
		fmt.Println(sim.proc[i].String())
	}
}
	
func (sim *BlockchainSimulator) InitGroups(n uint16) {
	sim.group = make([]GroupSimulator, n)
	sim.grpmap = make(map[common.Address]*GroupSimulator)
	r := sim.seed.Ders("InitGroups")
	// build a temporary state datastructure from processes
	/* s := state.NewState()
	for _, p := range sim.proc {
		s.AddNode(p.reginfo)
	} */
	// create n groups
	for i := 0; i < int(n); i++ {
		// choose members based on r 
		/* groupinfo := s.NewRandomGroup(r.Deri(i), sim.groupSize)
	        groupinfo.Log() */
		// LATER: replace the following using groupinfo 
		indices := r.Deri(i).RandomPerm(len(sim.proc),int(sim.groupSize))
		members := make([]*ProcessSimulator, sim.groupSize)
		for j, idx := range indices {
			members[j] = &(sim.proc[idx])
		}
		sim.group[i] = NewGroupSimulator(members, sim.threshold)
		sim.grpmap[sim.group[i].Address()] = &sim.group[i]
		fmt.Println(sim.group[i].String())
	}
}

func NewBlockchainSimulator(seed bls.Rand, groupSize uint16, threshold uint16, nProcesses uint, nGroups uint16) BlockchainSimulator {
	sim := BlockchainSimulator{seed: seed, groupSize: groupSize, threshold: threshold}
	sim.Log()

	// Start the processes first
	fmt.Printf("--- Process setup: (N)%d\n", nProcesses)
	sim.InitProcs(nProcesses)

	// Start the groups
	fmt.Printf("--- Group setup: (m)%d\n", nGroups)
	sim.InitGroups(nGroups)

	// Build the genesis block 
	genesis := state.NewState()
	for _, p := range sim.proc {
		genesis.AddNode(p.reginfo)
		// this includes verification of proof-of-possession
	}
	for _, g := range sim.group {
		genesis.AddGroup(g.reginfo)
	}
	// the sig field remains empty because the genesis block is not signed

	// print op counts
	bls.PrintCtrs()

	// Build the chain with 1 block
	sim.chain = append(sim.chain, genesis)

	return sim
}

func (sim *BlockchainSimulator) Advance(n uint, verbose bool) {
	if n == 0 {
		return
	}
	// choose tip
	tip := sim.Tip()
	// shallow cop y tip (?)
	newstate := tip
	// select pre-determined random group from tip
	a := tip.SelectedGroupAddress()
	g := sim.grpmap[a]
	// get new group signature
	/* TODO: transition to new logging package
	if verbose {
		spew.Dump(tip.Rand())
		spew.Dump(tip.Rand().Bytes())
		spew.Dump(tip.Rand().String())
	}
	*/
	sig := g.Sign(tip.Rand().Bytes())
	if Double_check {
		if !bls.VerifySig(tip.GroupPubkey(a), tip.Rand().Bytes(), sig) {
			fmt.Println("Error: group signature not valid.")
		}
	}

	// sign new state by group
	newstate.SetSignature(sig)
	
	// append new state	
	sim.chain = append(sim.chain, newstate)
	// recurse
	sim.Advance(n-1, verbose) 
	return
}

func (sim *BlockchainSimulator) Log() {
	seed := sim.seed.Bytes()
	fmt.Printf("BlkCh: (n)%d (k)%d (seed)%x\n", sim.groupSize, sim.threshold, seed[:8])
/*
	fmt.Println("  groups: ", len(sim.group))
	fmt.Println("  processes: ", len(sim.proc))
	fmt.Println("  chain height: ", len(sim.chain))
	sim.chain[len(sim.chain)-1].Log()
	for _, p := range sim.proc {
		p.Log()
	}
	for _, g := range sim.group {
		g.Log()
	}
*/
}

func (sim *BlockchainSimulator) Length() int {
	return len(sim.chain)
}

func (sim *BlockchainSimulator) Tip() state.State {
	return sim.chain[len(sim.chain)-1]
}

func (sim *BlockchainSimulator) Touch() {
}
