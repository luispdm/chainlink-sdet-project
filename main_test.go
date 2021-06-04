package main

import (
	"fmt"
	"testing"

	. "chainlink-sdet-project/config"
	. "chainlink-sdet-project/contracts"

	"github.com/stretchr/testify/suite"
)

const (
	BTC        = "0xf570deefff684d964dc3e15e1f9414283e3f7419"
	ETH        = "0x00c7a37b03690fb9f41b5c5af8131735c7275446"
	LINK       = "0x8cde021f0bfa5f82610e8ce46493cf66ac04af53"
	DOGE       = "0x0227903281b0421666f1e9161e8828c7112b8e86"
	MAX_ROUNDS = 100
)

type Test struct {
	Name       string
	Address    string
	NoOfRounds int
	Threshold  float64
}

type ChainlinkTestSuite struct {
	suite.Suite
}

func (suite *ChainlinkTestSuite) SetupSuite() {
	LoadConfig()
}

func (suite *ChainlinkTestSuite) SetupTest() {
	// demo hook - not needed if empty
}

func (suite *ChainlinkTestSuite) TearDownTest() {
	// demo hook - not needed if empty
}

func (suite *ChainlinkTestSuite) TearDownSuite() {
	// demo hook - not needed if empty
}

func TestChainlinkTestSuite(t *testing.T) {
	suite.Run(t, new(ChainlinkTestSuite))
}

func (suite *ChainlinkTestSuite) TestAnswerDeviation() {
	tests := []Test{
		{Name: "BTC", Address: BTC, Threshold: 10, NoOfRounds: 5},
		{Name: "DOGE", Address: DOGE, Threshold: 3, NoOfRounds: 5},
		{Name: "ETH", Address: ETH, Threshold: 4, NoOfRounds: 5},
		{Name: "LINK", Address: LINK, Threshold: 4, NoOfRounds: 5},
	}
	for _, test := range tests {
		test := test // required for parallel execution

		suite.Run(test.Name, func() {
			if Conf.Parallel {
				suite.T().Parallel()
			}
			suite.Require().LessOrEqual(test.NoOfRounds, MAX_ROUNDS,
				fmt.Sprintf("Test '%s': number of rounds exceeded. '%d' allowed, got '%d'",
					test.Name, MAX_ROUNDS, test.NoOfRounds))

			bClient, err := NewBlockchainClient(test.Address)
			suite.Require().NoError(err)
			err = bClient.FetchLastRounds(test.NoOfRounds)
			suite.Require().NoError(err)
			iter, err := bClient.FilterSubmissionReceived()
			suite.Require().NoError(err)

			defer iter.Close()
			for iter.Next() {
				if iter.Error() != nil {
					suite.T().Logf("Iteration error: %s", iter.Error().Error())
					continue
				}
				oraclePrice := bClient.GetPrice(iter.Event.Submission)
				answerDev := bClient.GetAnswerDev(oraclePrice, iter.Event.Round)
				errMessage := fmt.Sprintf("Test '%s': answer deviation of oracle '%v' at round '%v' should be "+
					"less than or equal to '%v%%', got '%v%%'",
					test.Name, iter.Event.Oracle, iter.Event.Round, test.Threshold, answerDev)

				suite.Assert().LessOrEqual(answerDev, test.Threshold, errMessage)
			}
		})
	}
}
