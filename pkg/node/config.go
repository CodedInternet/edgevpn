// Copyright © 2021 Ettore Di Giacinto <mudler@mocaccino.org>
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program; if not, see <http://www.gnu.org/licenses/>.

package node

import (
	"context"
	"time"

	"github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/mudler/edgevpn/pkg/blockchain"
	discovery "github.com/mudler/edgevpn/pkg/discovery"
	hub "github.com/mudler/edgevpn/pkg/hub"
	protocol "github.com/mudler/edgevpn/pkg/protocol"
)

// Config is the node configuration
type Config struct {
	// ExchangeKey is a Symmetric key used to seal the messages
	ExchangeKey string

	// RoomName is the OTP token gossip room where all peers are subscribed to
	RoomName string

	// ListenAddresses is the discovery peer initial bootstrap addresses
	ListenAddresses []discovery.AddrList

	// Insecure disables secure p2p e2e encrypted communication
	Insecure bool

	// Handlers are a list of handlers subscribed to messages received by the vpn interface
	Handlers, GenericChannelHandler []Handler

	MaxMessageSize  int
	SealKeyInterval int

	ServiceDiscovery []ServiceDiscovery
	NetworkServices  []NetworkService
	Logger           log.StandardLogger

	SealKeyLength    int
	InterfaceAddress string

	Store blockchain.Store

	// Handle is a handle consumed by HumanInterfaces to handle received messages
	Handle                     func(bool, *hub.Message)
	StreamHandlers             map[protocol.Protocol]StreamHandler
	AdditionalOptions, Options []libp2p.Option

	DiscoveryInterval, LedgerSyncronizationTime, LedgerAnnounceTime time.Duration
	DiscoveryBootstrapPeers                                         discovery.AddrList

	Whitelist, Blacklist []string

	// GenericHub enables generic hub
	GenericHub bool

	Sealer    Sealer
	PeerGater Gater
}

type Gater interface {
	Gate(*Node, peer.ID) bool
	Enable()
	Disable()
	Enabled() bool
}

type Sealer interface {
	Seal(string, string) (string, error)
	Unseal(string, string) (string, error)
}

// NetworkService is a service running over the network. It takes a context, a node and a ledger
type NetworkService func(context.Context, Config, *Node, *blockchain.Ledger) error

type StreamHandler func(*Node, *blockchain.Ledger) func(stream network.Stream)

type Handler func(*blockchain.Ledger, *hub.Message, chan *hub.Message) error

type ServiceDiscovery interface {
	Run(log.StandardLogger, context.Context, host.Host) error
	Option(context.Context) func(c *libp2p.Config) error
}

type Option func(cfg *Config) error

// Apply applies the given options to the config, returning the first error
// encountered (if any).
func (cfg *Config) Apply(opts ...Option) error {
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(cfg); err != nil {
			return err
		}
	}
	return nil
}
