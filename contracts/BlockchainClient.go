package contracts

import (
	"fmt"
	"math"
	"math/big"

	. "chainlink-sdet-project/contracts/ethereum"
	. "chainlink-sdet-project/config"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	FUNC_RETURNED = "Function returned: "
)

type BlockchainClient struct {
	fAggr    *FluxAggregator
	decimals uint8
	rounds   map[uint32]float64
}

func NewBlockchainClient(address string) (*BlockchainClient, error) {
	bC := &BlockchainClient{}
	eClient, err := ethclient.Dial(Conf.Wss)
	if err != nil {
		return bC, fmt.Errorf("Error instantiating the ETH client at WebSocket address '%s'. %s\n'%s'\n\n",
			Conf.Wss, FUNC_RETURNED, err.Error())
	}
	bC.fAggr, err = NewFluxAggregator(common.HexToAddress(address), eClient)
	bC.decimals, err = bC.fAggr.Decimals(&bind.CallOpts{})
	if err != nil {
		return bC, fmt.Errorf("Error returning contract's decimals. %s\n'%s'\n\n", FUNC_RETURNED, err.Error())
	}
	return bC, nil
}

func (bC *BlockchainClient) FilterSubmissionReceived() (*FluxAggregatorSubmissionReceivedIterator, error) {
	iter, err := bC.fAggr.FilterSubmissionReceived(&bind.FilterOpts{}, []*big.Int{}, bC.getRoundIds(), []common.Address{})
	if err != nil {
		return nil, fmt.Errorf("Error instantiating filter for 'SubmissionReceived' event. %s\n'%s'\n\n", FUNC_RETURNED, err.Error())
	}
	return iter, nil
}

func (bC *BlockchainClient) FetchLastRounds(noOfRounds int) error {
	lRData, err := bC.fAggr.LatestRoundData(&bind.CallOpts{})
	if err != nil {
		return fmt.Errorf("Error getting latest round data. %s\n'%s'\n\n", FUNC_RETURNED, err.Error())
	}
	roundId := uint32(lRData.RoundId.Int64())
	bC.rounds = map[uint32]float64{
		roundId: bC.GetPrice(lRData.Answer),
	}
	for i := 0; i < noOfRounds-1; i++ {
		roundId -= 1
		rData, err := bC.fAggr.GetRoundData(&bind.CallOpts{}, big.NewInt(int64(roundId)))
		if err != nil {
			return fmt.Errorf("Error getting data for round %d. %s\n'%s'\n\n", roundId, FUNC_RETURNED, err.Error())
		}
		bC.rounds[roundId] = bC.GetPrice(rData.Answer)
	}
	return nil
}

func (bC *BlockchainClient) GetPrice(price *big.Int) float64 {
	return float64(price.Int64()) / math.Pow(10, float64(bC.decimals))
}

func (bC *BlockchainClient) GetAnswerDev(toCompare float64, roundId uint32) float64 {
	return math.Abs(((toCompare - bC.rounds[roundId]) * 100) / bC.rounds[roundId])
}

func (bC *BlockchainClient) getRoundIds() []uint32 {
	roundIds := make([]uint32, len(bC.rounds))
	i := 0
	for k := range bC.rounds {
		roundIds[i] = uint32(k)
		i++
	}
	return roundIds
}
