package hotstuff

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/hyperledger/fabric/orderer/consensus"
	"github.com/hyperledger/fabric/protos/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/theodocius271/hotstuff/config"
	"github.com/theodocius271/hotstuff/crypto"
	pb "github.com/theodocius271/hotstuff/proto"
	"google.golang.org/protobuf/proto"
)

// 模拟 ConsenterSupport 接口
type mockConsenterSupport struct {
	mock.Mock
}

func (m *mockConsenterSupport) Chain() string {
	args := m.Called()
	return args.String(0)
}

func (m *mockConsenterSupport) SharedConfig() consensus.SharedConfig {
	args := m.Called()
	return args.Get(0).(consensus.SharedConfig)
}

func (m *mockConsenterSupport) ChannelID() string {
	args := m.Called()
	return args.String(0)
}

func (m *mockConsenterSupport) Sequence() uint64 {
	args := m.Called()
	return args.Get(0).(uint64)
}

func (m *mockConsenterSupport) ProcessNormalMsg(env *common.Envelope) (*common.Envelope, error) {
	args := m.Called(env)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*common.Envelope), args.Error(1)
}

func (m *mockConsenterSupport) ProcessConfigMsg(env *common.Envelope) (*common.Envelope, uint64, error) {
	args := m.Called(env)
	if args.Get(0) == nil {
		return nil, args.Get(1).(uint64), args.Error(2)
	}
	return args.Get(0).(*common.Envelope), args.Get(1).(uint64), args.Error(2)
}

func (m *mockConsenterSupport) CreateNextBlock(messages []*common.Envelope) *common.Block {
	args := m.Called(messages)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*common.Block)
}

func (m *mockConsenterSupport) WriteBlock(block *common.Block, encodedMetadataValue []byte) {
	m.Called(block, encodedMetadataValue)
}

func (m *mockConsenterSupport) WriteConfigBlock(block *common.Block, encodedMetadataValue []byte) {
	m.Called(block, encodedMetadataValue)
}

// 模拟 SharedConfig 接口
type mockSharedConfig struct {
	mock.Mock
}

func (m *mockSharedConfig) BatchSize() consensus.BatchSize {
	args := m.Called()
	return args.Get(0).(consensus.BatchSize)
}

func (m *mockSharedConfig) ConsensusType() string {
	args := m.Called()
	return args.String(0)
}

// 模拟 BatchSize 接口
type mockBatchSize struct {
	maxMessageCount   uint32
	absoluteMaxBytes  uint32
	preferredMaxBytes uint32
}

func (m mockBatchSize) MaxMessageCount() uint32 {
	return m.maxMessageCount
}

func (m mockBatchSize) AbsoluteMaxBytes() uint32 {
	return m.absoluteMaxBytes
}

func (m mockBatchSize) PreferredMaxBytes() uint32 {
	return m.preferredMaxBytes
}

// 创建测试环境的辅助函数
func createTestHotStuffEnv(t *testing.T) (*BasicHotStuff, *mockConsenterSupport) {
	// 创建模拟支持
	mockSupport := new(mockConsenterSupport)
	mockSharedConfig := new(mockSharedConfig)
	mockBatchSize := mockBatchSize{
		maxMessageCount:   10,
		absoluteMaxBytes:  1024 * 1024,
		preferredMaxBytes: 512 * 1024,
	}

	// 设置期望的调用
	mockSharedConfig.On("BatchSize").Return(mockBatchSize)
	mockSharedConfig.On("ConsensusType").Return("hotstuff")
	mockSupport.On("SharedConfig").Return(mockSharedConfig)
	mockSupport.On("ChannelID").Return("test-channel")
	mockSupport.On("Sequence").Return(uint64(0))
	mockSupport.On("CreateNextBlock", mock.Anything).Return(&common.Block{})
	mockSupport.On("WriteBlock", mock.Anything, mock.Anything).Return()
	mockSupport.On("WriteConfigBlock", mock.Anything, mock.Anything).Return()
	mockSupport.On("ProcessNormalMsg", mock.Anything).Return(&common.Envelope{}, nil)
	mockSupport.On("ProcessConfigMsg", mock.Anything).Return(&common.Envelope{}, uint64(0), nil)

	// 在测试之前确保创建了测试配置文件
	testConfig := config.HotStuffConfig{
		F:            1,  // 容忍的故障数量
		N:            4,  // 总节点数量
		BatchSize:    3,  // 批处理大小
		BatchTimeout: 2,  // 批处理超时（秒）
		Timeout:      10, // 视图更改超时
	}

	// 临时写入测试配置
	config.WriteTestConfig(&testConfig, t)

	// 创建 BasicHotStuff 实例并返回
	bhs := NewBasicHotStuff(1, mockSupport)

	// 替换默认密钥为测试密钥
	testPrivateKey, testPublicKey := crypto.GenerateTestKeys(t)
	bhs.Config.PrivateKey = testPrivateKey
	bhs.Config.PublicKey = testPublicKey

	return bhs, mockSupport
}

// 测试 NewBasicHotStuff 函数
func TestNewBasicHotStuff(t *testing.T) {
	bhs, _ := createTestHotStuffEnv(t)

	// 检查初始化是否正确
	assert.NotNil(t, bhs)
	assert.Equal(t, uint32(1), bhs.ID)
	assert.Equal(t, uint64(1), bhs.View.ViewNum)
	assert.Equal(t, uint32(1), bhs.View.Primary)
	assert.NotNil(t, bhs.BlockStorage)
	assert.NotNil(t, bhs.PrepareQC)
	assert.NotNil(t, bhs.PreCommitQC)
	assert.NotNil(t, bhs.CommitQC)
	assert.NotNil(t, bhs.CmdSet)
	assert.False(t, bhs.decided)

	// 检查是否有创世区块
	genesisBlock, err := bhs.BlockStorage.Get([]byte{0})
	assert.NoError(t, err)
	assert.NotNil(t, genesisBlock)
}

// 测试处理请求消息
func TestHandleMsgRequest(t *testing.T) {
	bhs, _ := createTestHotStuffEnv(t)

	// 创建mock交易
	tx := &pb.Transaction{
		Envelope:  &common.Envelope{},
		ConfigSeq: 0,
	}

	// 创建请求消息
	request := &pb.Request{
		Transaction: tx,
		IsNormal:    true,
	}

	// 创建包含请求的消息
	msg := &pb.Msg{
		Type: pb.MsgType_REQUEST,
		From: 2, // 从另一个节点发送
		To:   1, // 发送到当前节点
		Payload: &pb.Msg_Request{
			Request: request,
		},
	}

	// 发送消息前记录命令集大小
	beforeCmdSetSize := bhs.CmdSet.Size()

	// 处理消息
	bhs.handleMsg(msg)

	// 验证命令已添加到命令集
	assert.Equal(t, beforeCmdSetSize+1, bhs.CmdSet.Size())
}

// 测试节点ID和领导者计算
func TestLeaderComputation(t *testing.T) {
	bhs, _ := createTestHotStuffEnv(t)

	// 检查初始视图和领导者
	assert.Equal(t, uint64(1), bhs.View.ViewNum)
	assert.Equal(t, uint32(1), bhs.View.Primary)

	// 模拟视图变更
	bhs.View.ViewNum = 2
	bhs.View.Primary = bhs.GetLeader()
	assert.Equal(t, uint32(2), bhs.View.Primary)

	bhs.View.ViewNum = 3
	bhs.View.Primary = bhs.GetLeader()
	assert.Equal(t, uint32(3), bhs.View.Primary)

	bhs.View.ViewNum = 4
	bhs.View.Primary = bhs.GetLeader()
	assert.Equal(t, uint32(0), bhs.View.Primary)
}

// 测试批处理事件
func TestBatchEvent(t *testing.T) {
	bhs, _ := createTestHotStuffEnv(t)

	// 创建一些测试交易
	txs := make([]*pb.Transaction, 3)
	for i := 0; i < 3; i++ {
		txs[i] = &pb.Transaction{
			Envelope:  &common.Envelope{},
			ConfigSeq: 0,
		}
	}

	// 初始化CurExec为空
	bhs.CurExec = NewCurProposal()

	// 执行批处理
	bhs.batchEvent(txs, true)

	// 验证结果
	assert.NotNil(t, bhs.CurExec.Node)
	assert.Equal(t, 3, len(bhs.CurExec.Node.Commands))
	assert.NotNil(t, bhs.CurExec.DocumentHash)
	assert.Equal(t, 1, len(bhs.CurExec.PrepareVote))
	assert.True(t, bhs.CurExec.Node.IsNormal)
}

// 测试准备阶段消息处理
func TestHandlePrepareMsg(t *testing.T) {
	bhs, _ := createTestHotStuffEnv(t)

	// 创建一个区块
	block := &pb.Block{
		Height:     1,
		ParentHash: []byte{0},
		Hash:       []byte{1},
		Commands:   []*pb.Transaction{},
		Committed:  false,
		IsNormal:   true,
	}

	// 创建一个高QC
	highQC := &pb.QuorumCert{
		BlockHash: []byte{0},
		Type:      pb.MsgType_PREPARE_VOTE,
		ViewNum:   0,
		Signature: nil,
	}

	// 创建一个Prepare消息
	prepareMsg := &pb.Msg{
		Type: pb.MsgType_PREPARE,
		From: 0, // 从领导者发送
		To:   1, // 发送到当前节点
		View: 1,
		Payload: &pb.Msg_Prepare{
			Prepare: &pb.Prepare{
				CurProposal: block,
				HighQC:      highQC,
			},
		},
	}

	// 设置当前执行状态
	bhs.CurExec = NewCurProposal()

	// 处理消息前创建钩子以拦截Unicast调用
	unicastCalled := false
	originalUnicast := bhs.Unicast
	bhs.Unicast = func(addr string, msg *pb.Msg) {
		unicastCalled = true
		assert.Equal(t, pb.MsgType_PREPARE_VOTE, msg.Type)
		originalUnicast(addr, msg)
	}

	// 处理消息
	bhs.handleMsg(prepareMsg)

	// 验证结果
	assert.True(t, unicastCalled)
	assert.NotNil(t, bhs.CurExec.DocumentHash)
	assert.NotNil(t, bhs.CurExec.Node)
}

// 测试预提交消息处理
func TestHandlePreCommitMsg(t *testing.T) {
	bhs, _ := createTestHotStuffEnv(t)

	// 创建一个区块
	block := &pb.Block{
		Height:     1,
		ParentHash: []byte{0},
		Hash:       []byte{1},
		Commands:   []*pb.Transaction{},
		Committed:  false,
		IsNormal:   true,
	}

	// 创建一个PrepareQC
	prepareQC := &pb.QuorumCert{
		BlockHash: []byte{1},
		Type:      pb.MsgType_PREPARE_VOTE,
		ViewNum:   1,
		Signature: []byte{1, 2, 3},
	}

	// 创建一个PreCommit消息
	preCommitMsg := &pb.Msg{
		Type: pb.MsgType_PRECOMMIT,
		From: 0, // 从领导者发送
		To:   1, // 发送到当前节点
		View: 1,
		Payload: &pb.Msg_PreCommit{
			PreCommit: &pb.PreCommit{
				PrepareQC: prepareQC,
			},
		},
	}

	// 设置当前执行状态
	bhs.CurExec = NewCurProposal()
	bhs.CurExec.Node = block

	// 准备文档哈希 (在实际代码中这是由Prepare消息处理产生的)
	testMsg := &pb.Msg{}
	msgBytes, _ := proto.Marshal(testMsg)
	docHash, _ := crypto.CreateDocumentHash(msgBytes, bhs.Config.PublicKey)
	bhs.CurExec.DocumentHash = docHash

	// 处理消息前创建钩子以拦截Unicast调用
	unicastCalled := false
	originalUnicast := bhs.Unicast
	bhs.Unicast = func(addr string, msg *pb.Msg) {
		unicastCalled = true
		assert.Equal(t, pb.MsgType_PRECOMMIT_VOTE, msg.Type)
		originalUnicast(addr, msg)
	}

	// 处理消息
	bhs.handleMsg(preCommitMsg)

	// 验证结果
	assert.True(t, unicastCalled)
	assert.Equal(t, prepareQC, bhs.PrepareQC)
}

// 测试视图变更的新视图消息处理
func TestHandleNewViewMsg(t *testing.T) {
	bhs, _ := createTestHotStuffEnv(t)

	// 创建PrepareQC
	prepareQC := &pb.QuorumCert{
		BlockHash: []byte{1},
		Type:      pb.MsgType_PREPARE_VOTE,
		ViewNum:   2, // 高于当前视图
		Signature: []byte{1, 2, 3},
	}

	// 创建新视图消息
	newViewMsg := &pb.Msg{
		Type: pb.MsgType_NEWVIEW,
		From: 2, // 从另一个节点发送
		To:   1, // 发送到当前节点
		View: 2,
		Payload: &pb.Msg_NewView{
			NewView: &pb.NewView{
				PrepareQC: prepareQC,
			},
		},
	}

	// 设置当前状态
	bhs.CurExec = NewCurProposal()
	bhs.decided = true

	// 处理多个新视图消息，以达到法定数量
	for i := 0; i < int(2*bhs.Config.F); i++ {
		newViewMsg.From = uint32(i + 2) // 不同节点发送
		bhs.handleMsg(newViewMsg)
	}

	// 验证结果 - 应该更新HighQC，并重置decided标志
	assert.Equal(t, prepareQC, bhs.HighQC)
	assert.False(t, bhs.decided)
}

// 测试投票消息 (PrepareVote) 处理
func TestHandlePrepareVoteMsg(t *testing.T) {
	bhs, _ := createTestHotStuffEnv(t)

	// 创建区块
	block := &pb.Block{
		Height:     1,
		ParentHash: []byte{0},
		Hash:       []byte{1},
		Commands:   []*pb.Transaction{},
		Committed:  false,
		IsNormal:   true,
	}

	// 设置当前执行状态
	bhs.CurExec = NewCurProposal()
	bhs.CurExec.Node = block

	// 准备文档哈希
	testMsg := &pb.Msg{}
	msgBytes, _ := proto.Marshal(testMsg)
	docHash, _ := crypto.CreateDocumentHash(msgBytes, bhs.Config.PublicKey)
	bhs.CurExec.DocumentHash = docHash

	// 创建部分签名
	partSig, _ := crypto.TSign(docHash, bhs.Config.PrivateKey, bhs.Config.PublicKey)
	partSigBytes, _ := json.Marshal(partSig)

	// 创建PrepareVote消息
	prepareVoteMsg := &pb.Msg{
		Type: pb.MsgType_PREPARE_VOTE,
		From: 2, // 从另一个节点发送
		To:   1, // 发送到当前节点
		View: 1,
		Payload: &pb.Msg_PrepareVote{
			PrepareVote: &pb.PrepareVote{
				BlockHash:  []byte{1},
				PartialSig: partSigBytes,
			},
		},
	}

	// 处理消息前创建钩子以拦截Broadcast调用
	broadcastCalled := false
	originalBroadcast := bhs.Broadcast
	bhs.Broadcast = func(msg *pb.Msg) {
		broadcastCalled = true
		assert.Equal(t, pb.MsgType_PRECOMMIT, msg.Type)
		originalBroadcast(msg)
	}

	// 设置当前节点为领导者
	bhs.View.Primary = bhs.ID

	// 处理足够多的PrepareVote消息，以形成QC (2f+1)
	for i := 0; i < int(2*bhs.Config.F+1); i++ {
		prepareVoteMsg.From = uint32(i + 2) // 不同节点发送
		bhs.handleMsg(prepareVoteMsg)
	}

	// 验证结果
	assert.True(t, broadcastCalled)
	assert.NotNil(t, bhs.PrepareQC)
	assert.Equal(t, pb.MsgType_PREPARE_VOTE, bhs.PrepareQC.Type)
}

// 测试决定消息处理
func TestHandleDecideMsg(t *testing.T) {
	bhs, _ := createTestHotStuffEnv(t)

	// 创建区块
	block := &pb.Block{
		Height:     1,
		ParentHash: []byte{0},
		Hash:       []byte{1},
		Commands:   []*pb.Transaction{},
		Committed:  false,
		IsNormal:   true,
	}

	// 创建CommitQC
	commitQC := &pb.QuorumCert{
		BlockHash: []byte{1},
		Type:      pb.MsgType_COMMIT_VOTE,
		ViewNum:   1,
		Signature: []byte{1, 2, 3},
	}

	// 创建Decide消息
	decideMsg := &pb.Msg{
		Type: pb.MsgType_DECIDE,
		From: 0, // 从领导者发送
		To:   1, // 发送到当前节点
		View: 1,
		Payload: &pb.Msg_Decide{
			Decide: &pb.Decide{
				CurProposal: block,
				CommitQC:    commitQC,
			},
		},
	}

	// 设置当前执行状态
	bhs.CurExec = NewCurProposal()
	bhs.CurExec.Node = block

	// 存储区块到区块存储
	bhs.BlockStorage.Put(block)

	// 处理消息
	bhs.handleMsg(decideMsg)

	// 验证结果
	assert.Equal(t, commitQC, bhs.CommitQC)

	// 视图应该增加
	assert.Equal(t, uint64(2), bhs.View.ViewNum)
}

// 测试写入账本功能
func TestWriteLedger(t *testing.T) {
	bhs, mockSupport := createTestHotStuffEnv(t)

	// 创建有普通交易的区块
	normalTx := &pb.Transaction{
		Envelope:  &common.Envelope{},
		ConfigSeq: 0,
	}

	normalBlock := &pb.Block{
		Height:     1,
		ParentHash: []byte{0},
		Hash:       []byte{1},
		Commands:   []*pb.Transaction{normalTx},
		Committed:  false,
		IsNormal:   true,
	}

	// 写入普通区块
	bhs.writeLedger(normalBlock)

	// 验证CreateNextBlock和WriteBlock被调用
	mockSupport.AssertCalled(t, "CreateNextBlock", mock.Anything)
	mockSupport.AssertCalled(t, "WriteBlock", mock.Anything, mock.Anything)

	// 创建配置交易区块
	configTx := &pb.Transaction{
		Envelope:  &common.Envelope{},
		ConfigSeq: 0,
	}

	configBlock := &pb.Block{
		Height:     2,
		ParentHash: []byte{1},
		Hash:       []byte{2},
		Commands:   []*pb.Transaction{configTx},
		Committed:  false,
		IsNormal:   false,
	}

	// 写入配置区块
	bhs.writeLedger(configBlock)

	// 验证ProcessConfigMsg和WriteConfigBlock被调用
	mockSupport.AssertCalled(t, "ProcessConfigMsg", mock.Anything)
	mockSupport.AssertCalled(t, "WriteConfigBlock", mock.Anything, mock.Anything)
}

// 测试超时处理
func TestTimeoutHandling(t *testing.T) {
	bhs, _ := createTestHotStuffEnv(t)

	// 设置短超时进行测试
	bhs.Config.Timeout = 100 * time.Millisecond

	// 设置当前状态
	bhs.CurExec = NewCurProposal()
	bhs.CurExec.Node = &pb.Block{
		Height:     1,
		ParentHash: []byte{0},
		Hash:       []byte{1},
		Commands:   []*pb.Transaction{},
		Committed:  false,
		IsNormal:   true,
	}

	// 创建跟踪通道
	timeoutOccurred := make(chan bool, 1)

	// 处理消息前创建钩子以拦截Unicast调用
	originalUnicast := bhs.Unicast
	bhs.Unicast = func(addr string, msg *pb.Msg) {
		if msg.Type == pb.MsgType_NEWVIEW {
			timeoutOccurred <- true
		}
		originalUnicast(addr, msg)
	}

	// 启动超时器
	bhs.TimeChan.SoftStartTimer()

	// 等待超时或超时
	select {
	case <-timeoutOccurred:
		// 验证超时处理正确
		assert.Equal(t, uint64(2), bhs.View.ViewNum) // 视图应增加
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout did not occur")
	}
}
