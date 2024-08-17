package types

func DefaultParams() Params {
	return Params{}
}

// Validate performs basic validation of params
func (p Params) Validate() error {
	return nil
}
