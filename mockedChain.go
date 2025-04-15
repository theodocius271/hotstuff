package hotstuff

import (
	"time"

	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric/common/flogging"
	"github.com/hyperledger/fabric/protoutil"
)

func (c *HotStuffChain) main() {
	logger := flogging.MustGetLogger("orderer.consensus.hotstuff")

	// 批处理超时定时器
	var batchTimer <-chan time.Time

	// 视图超时定时器
	viewTimer := c.setupViewTimer()

	for {
		select {
		case <-c.haltChan:
			// 优雅退出...
			return

		case msg := <-c.msgChan:
			// 处理各种共识消息...
			c.handleConsensusMsg(msg)

		case env := <-c.normalChan:
			// 处理客户端普通交易请求
			if !c.isLeader() {
				// 非主节点不处理客户端请求
				logger.Debugf("Not leader, dropping client request")
				continue
			}

			// 使用BlockCutter处理交易
			batches, pending := c.support.BlockCutter().Ordered(env)

			// 处理已准备好的批次
			for _, batch := range batches {
				block := c.createNextBlock(batch)
				c.broadcastProposal(block)
			}

			// 如果有待处理批次且定时器未启动，则启动批处理定时器
			if pending && batchTimer == nil {
				batchTimeout, _ := time.ParseDuration(c.support.BatchTimeout().String())
				batchTimer = time.After(batchTimeout)
			}

		case <-batchTimer:
			// 批处理超时，切割当前批次
			batch := c.support.BlockCutter().Cut()
			if len(batch) > 0 {
				block := c.createNextBlock(batch)
				c.broadcastProposal(block)
			}
			batchTimer = nil

		case <-viewTimer:
			// 处理视图超时...
			c.handleViewTimeout()
			viewTimer = c.setupViewTimer()
		}
	}
}

func (c *HotStuffChain) createNextBlock(batch []*common.Envelope) *common.Block {
	// 使用ConsenterSupport创建区块
	block := c.support.CreateNextBlock(batch)

	// 获取区块元数据
	blockMetadata := c.support.BlockMetadata()

	// 设置最新配置区块编号
	lastConfigBlockNum := c.support.LastConfigBlockNum()

	// 填充元数据
	c.support.BlockMetadataModifier().SetLastConfigBlockNum(block, lastConfigBlockNum)

	// 记录区块信息
	c.logger.Infof("Created block [%d], transactions: %d", block.Header.Number, len(batch))

	return block
}

func (c *HotStuffChain) broadcastProposal(block *common.Block) {
	// 创建HotStuff提议消息
	proposal := &protos.HotStuffMessage{
		Type:      protos.HotStuffMessage_PROPOSAL,
		View:      c.currentView,
		SeqNum:    block.Header.Number,
		BlockHash: protoutil.BlockHeaderHash(block.Header),
		Block:     block,
		Sender:    c.selfID,
	}

	// 向所有节点广播提议
	c.broadcastToAllNodes(proposal)

	// 本地处理提议
	c.processProposal(proposal)
}

func (c *HotStuffChain) broadcastToAllNodes(msg *protos.HotStuffMessage) {
	for _, node := range c.nodes {
		if node.ID != c.selfID { // 不发送给自己
			c.transport.Send(msg, node.ID)
		}
	}
}

func (c *HotStuffChain) processProposal(proposal *protos.HotStuffMessage) {
	// 验证提议
	if !c.validateProposal(proposal) {
		c.logger.Warningf("Invalid proposal received: %v", proposal)
		return
	}

	// 更新HotStuff状态
	c.blockInPhase = proposal.Block
	c.currentPhase = PhasePrePrepare

	// 向所有节点发送prepare消息
	prepare := &protos.HotStuffMessage{
		Type:      protos.HotStuffMessage_PREPARE,
		View:      c.currentView,
		SeqNum:    proposal.SeqNum,
		BlockHash: proposal.BlockHash,
		Sender:    c.selfID,
	}
	c.broadcastToAllNodes(prepare)

	// 更新状态并进入下一阶段
	c.enterPreparePhase()
}

func (c *HotStuffChain) processDecide(msg *protos.HotStuffMessage) {
	// 验证决定消息
	if !c.validateDecide(msg) {
		return
	}

	// 获取对应的区块
	block := c.blockPool[msg.BlockHash]
	if block == nil {
		c.logger.Warningf("Block not found for decide message: %v", msg)
		return
	}

	// 使用ConsenterSupport写入区块
	c.support.WriteBlock(block, nil)

	// 清理状态
	c.cleanupForNextRound()

	// 如果是主节点并且有足够的交易，立即开始下一个区块
	if c.isLeader() && c.support.BlockCutter().PendingBatchSize() > 0 {
		batch := c.support.BlockCutter().Cut()
		if len(batch) > 0 {
			nextBlock := c.createNextBlock(batch)
			c.broadcastProposal(nextBlock)
		}
	}
}

func (c *HotStuffChain) processConfigTx(env *common.Envelope) {
	// 使用ConsenterSupport处理配置交易
	configSeq := c.support.Sequence()
	err := c.support.ProcessConfigMsg(env)
	if err != nil {
		c.logger.Warningf("Error processing config transaction: %v", err)
		return
	}

	// 创建一个只包含配置交易的区块
	batch := []*common.Envelope{env}
	block := c.createNextBlock(batch)

	// 作为主节点发起提议
	if c.isLeader() {
		c.broadcastProposal(block)
	}
}

func (c *HotStuffChain) Start() {
	c.logger.Info("Starting HotStuff chain")

	// 初始化各种通道和状态
	c.initializeState()

	// 启动主循环
	go c.main()

	// 启动网络监听
	go c.listenForNetworkMessages()

	c.logger.Info("HotStuff chain started")
}
