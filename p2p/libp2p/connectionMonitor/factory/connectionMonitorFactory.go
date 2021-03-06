package factory

import (
	"fmt"

	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-go/p2p"
	"github.com/ElrondNetwork/elrond-go/p2p/libp2p/connectionMonitor"
)

// ArgsConnectionMonitorFactory represents the argument for the connection monitor factory
type ArgsConnectionMonitorFactory struct {
	Reconnecter                p2p.Reconnecter
	Sharder                    p2p.CommonSharder
	ThresholdMinConnectedPeers int
	TargetCount                int
}

// NewConnectionMonitor creates a new ConnectionMonitor instance
func NewConnectionMonitor(arg ArgsConnectionMonitorFactory) (ConnectionMonitor, error) {
	if check.IfNil(arg.Reconnecter) {
		return nil, p2p.ErrNilReconnecter
	}

	switch kadSharder := arg.Sharder.(type) {
	case connectionMonitor.Sharder:
		return connectionMonitor.NewLibp2pConnectionMonitorSimple(arg.Reconnecter, arg.ThresholdMinConnectedPeers, kadSharder)
	default:
		return nil, fmt.Errorf("%w for connection monitor: invalid type %T", p2p.ErrInvalidValue, kadSharder)
	}
}
