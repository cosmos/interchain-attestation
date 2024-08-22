package relayer

import (
	"context"
	"fmt"
	clienttypes "github.com/cosmos/ibc-go/v9/modules/core/02-client/types"
	"gitlab.com/tozd/go/errors"
	"time"
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

func WaitUntilCondition(timeoutAfter, pollingInterval time.Duration, fn func() (bool, error)) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutAfter)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("failed waiting for condition after %f seconds", timeoutAfter.Seconds())
		case <-time.After(pollingInterval):
			reachedCondition, err := fn()
			if err != nil {
				return fmt.Errorf("error occurred while waiting for condition: %s", err)
			}

			if reachedCondition {
				return nil
			}
		}
	}
}