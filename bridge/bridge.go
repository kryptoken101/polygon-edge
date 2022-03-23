package bridge

import (
	"github.com/0xPolygon/polygon-edge/bridge/sam"
	"github.com/0xPolygon/polygon-edge/bridge/statesync"
	"github.com/0xPolygon/polygon-edge/bridge/utils"
	"github.com/0xPolygon/polygon-edge/network"
	"github.com/0xPolygon/polygon-edge/types"
	"github.com/hashicorp/go-hclog"
)

type Bridge interface {
	Start() error
	Close() error
	SetValidators([]types.Address, uint64)
	StateSync() statesync.StateSync
}

type bridge struct {
	logger    hclog.Logger
	stateSync statesync.StateSync

	validatorSet utils.ValidatorSet
}

func NewBridge(
	logger hclog.Logger,
	network *network.Server,
	signer sam.Signer,
	dataDirURL string,
	config *Config,
) (Bridge, error) {
	bridgeLogger := logger.Named("bridge")

	valSet := utils.NewValidatorSet(nil, 0)
	stateSync, err := statesync.NewStateSync(
		bridgeLogger,
		network,
		signer,
		valSet,
		dataDirURL,
		config.RootChainURL.String(),
		config.RootChainContract,
		config.Confirmations,
	)

	if err != nil {
		return nil, err
	}

	return &bridge{
		logger:       bridgeLogger,
		stateSync:    stateSync,
		validatorSet: valSet,
	}, nil
}

func (b *bridge) Start() error {
	if err := b.stateSync.Start(); err != nil {
		return err
	}

	return nil
}

func (b *bridge) Close() error {
	if err := b.stateSync.Close(); err != nil {
		return err
	}

	return nil
}

func (b *bridge) SetValidators(validators []types.Address, threshold uint64) {
	b.validatorSet.SetValidators(validators, threshold)
}

func (b *bridge) StateSync() statesync.StateSync {
	return b.stateSync
}
