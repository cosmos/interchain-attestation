package relayer

import (
	clienttypes "github.com/cosmos/ibc-go/v9/modules/core/02-client/types"
	"gitlab.com/tozd/go/errors"
)

func GetClientType(clientID string) (ClientType, error) {
	clientTypeStr, _, err := clienttypes.ParseClientIdentifier(clientID)
	if err != nil {
		return 0, errors.Errorf("failed to parse client identifier: %w", err)
	}

	return ConvertClientType(clientTypeStr)
}

func ConvertClientType(clientType string) (ClientType, error) {
	switch clientType {
	case "tendermint":
		return TENDERMINT, nil
	case "07-tendermint":
		return TENDERMINT, nil
	case "attestation":
		return ATTESTATION, nil
	case "10-attestation":
		return ATTESTATION, nil
	default:
		return 0, errors.Errorf("invalid client type: %s", clientType)
	}
}