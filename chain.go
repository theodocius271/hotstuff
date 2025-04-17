package hotstuff

import (
	"fmt"
	"time"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/hyperledger/fabric/orderer/consensus"
	"github.com/hyperledger/fabric/protos/common"
	"github.com/pkg/errors"
	pb "github.com/theodocius271/hotstuff/proto"
)

type Chain struct {
	bhs      *BasicHotStuff
	support  consensus.ConsenterSupport
	exitChan chan struct{}
	logger   *flogging.FabricLogger
}

func NewChain(support consensus.ConsenterSupport) *Chain {
	logger.Info("NewChain - ", support.ChainID())

}

func (c *Chain) Order(env *common.Envelope, configSeq uint64) error {
	// verify envelope
	if env == nil {
		return errors.New("envelope is nil")
	}
	select {
	case <-c.exitChan:
		logger.Info("[CHAIN error exit config]")
		return fmt.Errorf("Exiting")
	default:
	}
	channelID := c.support.ChainID()

	// create msg{request}
	request := &pb.Request{
		Envelope:      env,
		ChannalId:     channelID,
		ConfigSeq:     configSeq,
		IsNormal:      true,
		TimeStamp:     uint64(time.Now().UnixNano()),
		ClientAddress: c.extractClientInfo(env),
	}
	msg := &pb.Msg{
		Payload: &pb.Msg_Request{
			Request: request,
		},
	}

	// redireat to prime
	c.logger.Debugf("Ordering normal transaction, sending request to primary node")
	return c.bhs.Unicast(c.bhs.GetNetworkInfo()[c.bhs.GetLeader()], msg)
}

func (c *Chain) Configure(env *common.Envelope, configSeq uint64) error {
	if env == nil {
		return errors.New("envelope is nil")
	}
	select {
	case <-c.exitChan:
		logger.Info("[CHAIN error exit config]")
		return fmt.Errorf("Exiting")
	default:
	}
	channelID := c.support.ChainID()

	request := &pb.Request{
		Envelope:      env,
		ChannalId:     channelID,
		ConfigSeq:     configSeq,
		IsNormal:      false,
		TimeStamp:     uint64(time.Now().UnixNano()),
		ClientAddress: c.extractClientInfo(env),
	}
	msg := &pb.Msg{
		Payload: &pb.Msg_Request{
			Request: request,
		},
	}

	c.logger.Debugf("Ordering config transaction with config sequence %d, sending request to primary node", configSeq)
	return c.bhs.Unicast(c.bhs.GetNetworkInfo()[c.bhs.GetLeader()], msg)
}

func (c *Chain) extractClientInfo(env *common.Envelope) string {
	// TODO find addr
	return "0"
}

func (ch *Chain) WaitReady() error {
	logger.Info("wait ready")
	return nil
}

func (ch *Chain) Halt() {
	logger.Info("halt")
	select {
	case <-ch.exitChan:
	default:
		close(ch.exitChan)
	}
}
