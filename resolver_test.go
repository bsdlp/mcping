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
	hostports, err := ResolveMinecraftHostPort(context.Background(), nil, "doesnotexist.example.com")
	suite.Assert().Error(err)
	suite.Assert().Empty(hostports)
}

func (suite *ResolverTests) TestSingleRecord() {
	hostports, err := ResolveMinecraftHostPort(context.Background(), nil, "single.bsdlp.dev")
	suite.Assert().NoError(err)
	suite.Assert().Len(hostports, 1)
	suite.Assert().Equal("mc.tonkat.su.:25565", hostports[0])
}

func (suite *ResolverTests) TestMultipleRecord() {
	hostports, err := ResolveMinecraftHostPort(context.Background(), nil, "multi.bsdlp.dev")
	suite.Assert().NoError(err)
	suite.Assert().Len(hostports, 2)
	suite.Assert().Equal("mc.tonkat.su.:25565", hostports[0])
	suite.Assert().Equal("mc.tonkat.su.:25566", hostports[1])
}

func TestResolver(t *testing.T) {
	suite.Run(t, &ResolverTests{})
}
