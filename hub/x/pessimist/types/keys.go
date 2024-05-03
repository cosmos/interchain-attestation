package types

const (
	// ModuleName defines the module name
	ModuleName = "pessimist"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_pessimist"
)

var (
	ParamsKey = []byte("p_pessimist")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
