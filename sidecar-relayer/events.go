package main

import (
	"fmt"
	abcitypes "github.com/cometbft/cometbft/abci/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	"slices"
)

// Repurposed from cosmos relayer
func parseClientIDFromEvents(events []abcitypes.Event) (string, error) {
	return parseAttributeFromEvents(events, []string{clienttypes.EventTypeCreateClient}, clienttypes.AttributeKeyClientID)
}

func parseConnectionIDFromEvents(events []abcitypes.Event) (string, error) {
	return parseAttributeFromEvents(events, []string{connectiontypes.EventTypeConnectionOpenInit, connectiontypes.EventTypeConnectionOpenTry}, connectiontypes.AttributeKeyConnectionID)
}

func parseChannelIDFromEvents(events []abcitypes.Event) (string, error) {
	return parseAttributeFromEvents(events, []string{channeltypes.EventTypeChannelOpenInit, channeltypes.EventTypeChannelOpenTry}, channeltypes.AttributeKeyChannelID)
}

func parseAttributeFromEvents(events []abcitypes.Event, eventTypes []string, attributeKey string) (string, error) {
	for _, event := range events {
		if slices.Contains(eventTypes, event.Type) {
			for _, attr := range event.Attributes {
				if attr.Key == attributeKey {
					return attr.Value, nil
				}
			}
		}
	}

	return "", fmt.Errorf("attribute not found")
}
