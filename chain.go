package hotstuff

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/orderer/consensus"
	"github.com/hyperledger/fabric/protos/common"
	"github.com/pkg/errors"
	pb "github.com/theodocius271/hotstuff/proto"
)

type Chain struct {
	bhs      *BasicHotStuff
	support  consensus.ConsenterSupport
	exitChan chan struct{}
}

func NewChain(support consensus.ConsenterSupport) *Chain {
	logger.Info("NewChain - ", support.ChainID())
	ch := &Chain{}
	currentNodeIDStr := os.Getenv("ORDERER_HOTSTUFF_NODEID")
	if currentNodeIDStr == "" {
		logger.Infof("ORDERER_HOTSTUFF_NODEID environment variable not set")
		return nil
	}
	currentNodeID, err := strconv.ParseUint(currentNodeIDStr, 10, 32)
	if err != nil {
		logger.Infof("invalid ORDERER_HOTSTUFF_NODEID: %v", err)
		return nil
	}
	ch.bhs = NewBasicHotStuff(uint32(currentNodeID), nil, support)
	ch.support = support
	ch.exitChan = make(chan struct{})

	return ch
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
	logger.Debugf("Ordering normal transaction, sending request to primary node")
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

	logger.Debugf("Ordering config transaction with config sequence %d, sending request to primary node", configSeq)
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
