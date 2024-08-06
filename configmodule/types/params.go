package types

func DefaultParams() Params {
	return Params{
		MinimumPower: 0,
	}
}

// Validate performs basic validation of params
func (p Params) Validate() error {
	return nil
}
