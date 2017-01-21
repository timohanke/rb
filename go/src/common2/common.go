package common2 

import (
	"sort"
	"github.com/ethereum/go-ethereum/common"
)

func SortAddresses(addr []common.Address) []common.Address {
	// sort addresses by converting to hex string and back
	hex := make([]string, len(addr))
        for i := 0; i<len(addr); i++ {
		hex[i] = addr[i].Hex()
	}
	sort.Strings(hex)
	sorted := make([]common.Address, len(addr))
        for i := 0; i<len(addr); i++ {
		sorted[i] = common.HexToAddress(hex[i])
	}
	return sorted
}
