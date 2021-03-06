package factory

import (
	"errors"
	"testing"

	"github.com/ElrondNetwork/elrond-go/config"
	"github.com/ElrondNetwork/elrond-go/factory/mock"
	"github.com/ElrondNetwork/elrond-go/p2p/libp2p"
	"github.com/stretchr/testify/require"
)

func TestNewNetworkComponentsFactory_NilStatusHandlerShouldErr(t *testing.T) {
	t.Parallel()

	ncf, err := NewNetworkComponentsFactory(config.P2PConfig{}, config.Config{}, nil, &mock.MarshalizerMock{})
	require.Nil(t, ncf)
	require.Equal(t, ErrNilStatusHandler, err)
}

func TestNewNetworkComponentsFactory_NilMarshalizerShouldErr(t *testing.T) {
	t.Parallel()

	ncf, err := NewNetworkComponentsFactory(config.P2PConfig{}, config.Config{}, &mock.AppStatusHandlerMock{}, nil)
	require.Nil(t, ncf)
	require.True(t, errors.Is(err, ErrNilMarshalizer))
}

func TestNewNetworkComponentsFactory_OkValsShouldWork(t *testing.T) {
	t.Parallel()

	ncf, err := NewNetworkComponentsFactory(config.P2PConfig{}, config.Config{}, &mock.AppStatusHandlerMock{}, &mock.MarshalizerMock{})
	require.NoError(t, err)
	require.NotNil(t, ncf)
}

func TestNetworkComponentsFactory_Create_ShouldErrDueToBadConfig(t *testing.T) {
	//TODO remove skip when external library is concurrent safe
	if testing.Short() {
		t.Skip("this test fails with race detector on because of the github.com/koron/go-ssdp lib")
	}

	ncf, _ := NewNetworkComponentsFactory(config.P2PConfig{}, config.Config{}, &mock.AppStatusHandlerMock{}, &mock.MarshalizerMock{})

	nc, err := ncf.Create()
	require.Error(t, err)
	require.Nil(t, nc)
}

func TestNetworkComponentsFactory_Create_ShouldWork(t *testing.T) {
	//TODO remove skip when external library is concurrent safe
	if testing.Short() {
		t.Skip("this test fails with race detector on because of the github.com/koron/go-ssdp lib")
	}

	p2pConfig := config.P2PConfig{
		Node: config.NodeConfig{
			Port: "0",
			Seed: "seed",
		},
		KadDhtPeerDiscovery: config.KadDhtPeerDiscoveryConfig{
			Enabled:                          false,
			RefreshIntervalInSec:             10,
			RandezVous:                       "erd/kad/1.0.0",
			InitialPeerList:                  []string{"peer0", "peer1"},
			BucketSize:                       10,
			RoutingTableRefreshIntervalInSec: 5,
		},
		Sharding: config.ShardingConfig{
			TargetPeerCount:         10,
			MaxIntraShardValidators: 10,
			MaxCrossShardValidators: 10,
			MaxIntraShardObservers:  10,
			MaxCrossShardObservers:  10,
			Type:                    "NilListSharder",
		},
	}
	ncf, _ := NewNetworkComponentsFactory(
		p2pConfig,
		config.Config{
			Debug: config.DebugConfig{
				Antiflood: config.AntifloodDebugConfig{
					Enabled:                    true,
					CacheSize:                  100,
					IntervalAutoPrintInSeconds: 1,
				},
			},
		},
		&mock.AppStatusHandlerMock{},
		&mock.MarshalizerMock{},
	)

	ncf.SetListenAddress(libp2p.ListenLocalhostAddrWithIp4AndTcp)

	nc, err := ncf.Create()
	require.NoError(t, err)
	require.NotNil(t, nc)
}
