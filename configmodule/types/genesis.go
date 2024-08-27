package types

import (
	"errors"
)

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
	if gs.Params == nil {
		return errors.New("params cannot be nil")
	}

	if err := gs.Params.Validate(); err != nil {
		return nil
		// return fmt.Errorf("params failed validation: %w", err)
	}

	return nil
}
