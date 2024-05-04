package types

const (
	// ModuleName defines the module name
	ModuleName = "pessimist"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_pessimist"

	ClientType = "69-pessimist"

	ValidatorObjectiveKeyPrefix = "validator_objective/"
)

var (
	ParamsKey = []byte("p_pessimist")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

func ValidatorObjectiveKey(clientID string) []byte {
	var key []byte
	indexBytes := []byte(clientID)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)
	return key
}
