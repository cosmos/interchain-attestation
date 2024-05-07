package pessimisticinterchaintest

import (
	testifysuite "github.com/stretchr/testify/suite"
	"testing"
)

func TestE2ETestSuite(t *testing.T) {
	testifysuite.Run(t, new(E2ETestSuite))
}

func (s *E2ETestSuite) TestTheKitchenSink() {
	s.NotNil(s.ic)
}

