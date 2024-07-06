package types

// NewGenesisState creates a new genesis state.
func NewGenesisState(params Params) *GenesisState {
	return &GenesisState{
		Params: &params,
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(DefaultParams())
}

// Validate performs basic genesis state validation returning an error upon any failure.
func (gs GenesisState) Validate() error {
	return nil
}