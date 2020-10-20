package mcping

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ResolverTests struct {
	suite.Suite
}

func (suite *ResolverTests) TestNXDomain() {
	hostport, err := ResolveMinecraftHostPort(context.Background(), nil, "doesnotexist.example.com")
	suite.Assert().NoError(err)
	suite.Assert().Empty(hostport)
}

func TestResolver(t *testing.T) {
	suite.Run(t, &ResolverTests{})
}
