package hotstuff

import (
	"encoding/base64"
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

// Order 处理普通交易的接收，将其封装为Request并转发给主节点
func (c *Chain) Order(env *common.Envelope, configSeq uint64) error {
	// 1. 验证Envelope的有效性
	if env == nil {
		return errors.New("envelope is nil")
	}
	select {
	case <-c.exitChan:
		logger.Info("[CHAIN error exit config]")
		return fmt.Errorf("Exiting")
	default:
	}

	// 2. 从交易中提取channelID
	channelID := c.support.ChainID()

	// 3. 创建Request对象封装交易
	request := &pb.Request{
		Envelope:      env,
		ChannalId:     channelID,
		ConfigSeq:     configSeq,
		IsNormal:      true, // 普通交易，非配置交易
		TimeStamp:     uint64(time.Now().UnixNano()),
		ClientAddress: c.extractClientInfo(env), // 假设有方法提取客户端信息
	}

	// 4. 创建包含Request的Msg对象
	msg := &pb.Msg{
		Payload: &pb.Msg_Request{
			Request: request,
		},
	}

	// 5. 转发给主节点
	c.logger.Debugf("Ordering normal transaction, sending request to primary node")
	return c.bhs.Unicast(c.bhs.GetNetworkInfo()[c.bhs.GetLeader()], msg)
}

// Configure 处理配置交易的接收，将其封装为Request并转发给主节点
func (c *Chain) Configure(env *common.Envelope, configSeq uint64) error {
	// 1. 验证Envelope的有效性
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

	// 4. 创建Request对象封装配置交易
	request := &pb.Request{
		Envelope:      env,
		ChannalId:     channelID,
		ConfigSeq:     configSeq,
		IsNormal:      false, // 配置交易，非普通交易
		TimeStamp:     uint64(time.Now().UnixNano()),
		ClientAddress: c.extractClientInfo(env), // 假设有方法提取客户端信息
	}

	// 5. 创建包含Request的Msg对象
	msg := &pb.Msg{
		Payload: &pb.Msg_Request{
			Request: request,
		},
	}

	// 6. 转发给主节点
	c.logger.Debugf("Ordering config transaction with config sequence %d, sending request to primary node", configSeq)
	return c.bhs.Unicast(c.bhs.GetNetworkInfo()[c.bhs.GetLeader()], msg)
}

// 辅助方法：从交易中提取客户端信息
func (c *Chain) extractClientInfo(env *common.Envelope) string {
	// 实现从envelope中提取客户端信息的逻辑
	// 这可能包括从签名中提取身份信息，或从应用头部提取客户端字段

	// 简单实现示例
	sighdr, err := protoutil.SignatureHeader(env)
	if err != nil {
		c.logger.Warningf("Failed to extract signature header: %s", err)
		return "unknown-client"
	}

	// 使用创建者字段作为客户端地址标识
	// 在实际实现中，可能需要解析此字段以获取更有意义的标识
	return base64.StdEncoding.EncodeToString(sighdr.Creator)
}
