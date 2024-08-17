package voteextension_test

import (
	"testing"
)

func TestExtendVote(t *testing.T) {
	/*testKey := storetypes.NewKVStoreKey("upgrade")
	ctx := testutil.DefaultContext(testKey, storetypes.NewTransientStoreKey("transient_test"))

	err := os.Setenv(attestationabci.SidecarAddressEnv, "")
	require.NoError(t, err)

	appModule := attestationabci.NewAppModule()
	responseExtendVote, err := appModule.ExtendVote(ctx, &abci.RequestExtendVote{})
	require.NoError(t, err)
	_ = responseExtendVote*/
	// TODO: Create a mock sidecar server and set the address to the environment variable and test with that
}
