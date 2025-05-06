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
	server   *Server
	support  consensus.ConsenterSupport
	exitChan chan struct{}
	started  bool
}

func NewChain(support consensus.ConsenterSupport) *Chain {
	logger.Info("NewChain - ", support.ChainID())
	ch := &Chain{
		support:  support,
		started:  true,
		exitChan: make(chan struct{}),
	}
	currentNodeIDStr := os.Getenv("ORDERER_HOTSTUFF_NODEID")
	if currentNodeIDStr == "" {
		logger.Errorf("ORDERER_HOTSTUFF_NODEID environment variable not set")
		return nil
	}
	currentNodeID, err := strconv.ParseUint(currentNodeIDStr, 10, 32)
	if err != nil {
		logger.Errorf("invalid ORDERER_HOTSTUFF_NODEID: %v", err)
		return nil
	}
	ch.server = NewServer(uint32(currentNodeID), support)
	err = ch.server.RunAsync()
	if err != nil {
		logger.Errorf("Cannot Start Server")
		return nil
	}

	return ch
}

func (c *Chain) Order(env *common.Envelope, configSeq uint64) error {
	// verify envelope
	if env == nil {
		return errors.New("envelope is nil")
	}
	select {
	case <-c.exitChan:
		logger.Info("[CHAIN error exit normal]")
		return fmt.Errorf("Exiting")
	default:
	}

	seq := c.support.Sequence()
	if configSeq < seq {
		logger.Warnf("Normal message was validated against %d, although current config seq has advanced (%d)", configSeq, seq)
		if _, err := c.support.ProcessNormalMsg(env); err != nil {
			return errors.Errorf("bad normal message: %s", err)
		}
	}

	channelID := c.support.ChainID()

	transaction := &pb.Transaction{
		Envelope:  env,
		ChannalId: channelID,
		ConfigSeq: configSeq,
	}
	// create msg{request}
	request := &pb.Request{
		Transaction: transaction,
		IsNormal:    true,
		TimeStamp:   uint64(time.Now().UnixNano()),
	}
	msg := &pb.Msg{
		Payload: &pb.Msg_Request{
			Request: request,
		},
		ChannalID: channelID,
	}

	// redireat to self
	logger.Debugf("Ordering normal transaction, pakaging as Req")
	// return c.bhs.Unicast(c.bhs.GetNetworkInfo()[c.bhs.GetLeader()], msg)
	impl := c.server.GetHotStuffImpl()
	// return impl.Unicast(impl.GetSelfInfo().Address, msg)
	impl.GetMsgEntrance() <- msg
	return nil
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

	seq := c.support.Sequence()
	var err error
	if configSeq < seq {
		logger.Warnf("Normal message was validated against %d, although current config seq has advanced (%d)", configSeq, seq)
		if env, _, err = c.support.ProcessConfigMsg(env); err != nil {
			return errors.Errorf("bad normal message: %s", err)
		}
	}

	transaction := &pb.Transaction{
		Envelope:  env,
		ChannalId: channelID,
		ConfigSeq: configSeq,
	}
	// create msg{request}
	request := &pb.Request{
		Transaction: transaction,
		IsNormal:    false,
		TimeStamp:   uint64(time.Now().UnixNano()),
	}
	msg := &pb.Msg{
		Payload: &pb.Msg_Request{
			Request: request,
		},
		ChannalID: channelID,
	}

	// redireat to self
	logger.Debugf("Ordering Config transaction, pakaging as Req")
	// return c.bhs.Unicast(c.bhs.GetNetworkInfo()[c.bhs.GetLeader()], msg)
	impl := c.server.GetHotStuffImpl()
	// return impl.Unicast(impl.GetSelfInfo().Address, msg)
	impl.GetMsgEntrance() <- msg
	return nil
}

func (ch *Chain) WaitReady() error {
	logger.Info("wait ready")
	return nil
}

func (ch *Chain) Halt() {
	logger.Info("halt")
	if !ch.started {
		logger.Info("Chain not running")
		return
	}
	logger.Info("Halting HotStuff consensus chain")
	ch.server.Shutdown()
	close(ch.exitChan)
	ch.started = false

}

func (c *Chain) Start() {
	logger.Info("Start() called, but chain was already started during creation")
}

func (ch *Chain) Errored() <-chan struct{} {
	logger.Errorf("We Are Doomed!")
	return ch.exitChan
}
