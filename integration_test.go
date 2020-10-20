package mcping_test

import (
	"context"
	"net"
	"strconv"
	"testing"

	"github.com/bsdlp/mcping"
	"github.com/stretchr/testify/suite"
)

type IntegrationTests struct {
	suite.Suite
}

func (suite *IntegrationTests) TestNXDomain() {
	hostports, err := mcping.ResolveMinecraftHostPort(context.Background(), nil, "doesnotexist.example.com")
	suite.Assert().NoError(err)
	suite.Assert().Empty(hostports)
}

func (suite *IntegrationTests) TestNoSRVDomain() {
	hostports, err := mcping.ResolveMinecraftHostPort(context.Background(), nil, "mc.tonkat.su")
	suite.Assert().NoError(err)
	suite.Assert().Equal(mcping.Server{Host: "mc.tonkat.su", Port: 25565}, hostports[0])
}

func (suite *IntegrationTests) TestSingleRecord() {
	hostports, err := mcping.ResolveMinecraftHostPort(context.Background(), nil, "single.bsdlp.dev")
	suite.Assert().NoError(err)
	suite.Assert().Len(hostports, 1)
	suite.Assert().Equal(mcping.Server{Host: "mc.tonkat.su.", Port: 25565}, hostports[0])
}

func (suite *IntegrationTests) TestMultipleRecord() {
	hostports, err := mcping.ResolveMinecraftHostPort(context.Background(), nil, "multi.bsdlp.dev")
	suite.Assert().NoError(err)
	suite.Assert().Len(hostports, 2)
	suite.Assert().Equal(mcping.Server{Host: "mc.tonkat.su.", Port: 25565}, hostports[0])
	suite.Assert().Equal(mcping.Server{Host: "mc.tonkat.su.", Port: 25566}, hostports[1])
}

func (suite *IntegrationTests) TestPing() {
	hostports, err := mcping.ResolveMinecraftHostPort(context.Background(), nil, "mc.sep.gg")
	suite.Require().NoError(err)
	conn, err := net.Dial("tcp", net.JoinHostPort(hostports[0].Host, strconv.FormatUint(uint64(hostports[0].Port), 10)))
	suite.Assert().NoError(err)

	response, err := mcping.Ping(context.Background(), conn, hostports[0])
	suite.Assert().NoError(err)
	suite.Assert().Empty(response)
}

func TestIntegration(t *testing.T) {
	suite.Run(t, &IntegrationTests{})
}
