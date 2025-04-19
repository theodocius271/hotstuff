package hotstuff

import (
	"errors"
	"testing"

	"github.com/hyperledger/fabric/common/channelconfig"
	"github.com/hyperledger/fabric/orderer/common/blockcutter"
	"github.com/hyperledger/fabric/orderer/common/msgprocessor"
	"github.com/hyperledger/fabric/protos/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	pb "github.com/theodocius271/hotstuff/proto"
)

// 禁用日志输出，以防止测试输出中有太多日志
func init() {
	// flogging.SetModuleLevel("hotstuff", "ERROR")
}

// 创建一个ConsenterSupport的模拟实现，基于Fabric 1.4.4
type mockConsenterSupport struct {
	mock.Mock
}

func (mcs *mockConsenterSupport) Sequence() uint64 {
	args := mcs.Called()
	return args.Get(0).(uint64)
}

// 修正 ProcessNormalMsg 方法的签名以匹配 Fabric 1.4.4
func (mcs *mockConsenterSupport) ProcessNormalMsg(env *common.Envelope) (uint64, error) {
	args := mcs.Called(env)
	return args.Get(0).(uint64), args.Error(1)
}

func (mcs *mockConsenterSupport) ProcessConfigMsg(env *common.Envelope) (*common.Envelope, uint64, error) {
	args := mcs.Called(env)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).(*common.Envelope), args.Get(1).(uint64), args.Error(2)
}

func (mcs *mockConsenterSupport) CreateNextBlock(messages []*common.Envelope) *common.Block {
	args := mcs.Called(messages)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*common.Block)
}

func (mcs *mockConsenterSupport) WriteBlock(block *common.Block, encodedMetadataValue []byte) {
	mcs.Called(block, encodedMetadataValue)
}

func (mcs *mockConsenterSupport) WriteConfigBlock(block *common.Block, encodedMetadataValue []byte) {
	mcs.Called(block, encodedMetadataValue)
}

// Fabric 1.4.4 中的其他ConsenterSupport接口方法
func (mcs *mockConsenterSupport) Height() uint64 {
	args := mcs.Called()
	return args.Get(0).(uint64)
}

func (mcs *mockConsenterSupport) ChainID() string {
	return "test-channel"
}

func (mcs *mockConsenterSupport) Sign(message []byte) ([]byte, error) {
	args := mcs.Called(message)
	return args.Get(0).([]byte), args.Error(1)
}

func (mcs *mockConsenterSupport) NewSignatureHeader() (*common.SignatureHeader, error) {
	args := mcs.Called()
	return args.Get(0).(*common.SignatureHeader), args.Error(1)
}

func (mcs *mockConsenterSupport) ClassifyMsg(chdr *common.ChannelHeader) msgprocessor.Classification {
	args := mcs.Called(chdr)
	return args.Get(0).(msgprocessor.Classification)
}

func (mcs *mockConsenterSupport) BlockCutter() blockcutter.Receiver {
	args := mcs.Called()
	return args.Get(0).(blockcutter.Receiver)
}

func (mcs *mockConsenterSupport) SharedConfig() channelconfig.Orderer {
	args := mcs.Called()
	return args.Get(0).(channelconfig.Orderer)
}

// 添加Fabric 1.4.4中ConsenterSupport接口所需的Append方法
func (mcs *mockConsenterSupport) Append(block *common.Block) error {
	args := mcs.Called(block)
	return args.Error(0)
}

// 添加Fabric 1.4.4中ConsenterSupport接口所需的Block方法
func (mcs *mockConsenterSupport) Block(number uint64) *common.Block {
	args := mcs.Called(number)
	return args.Get(0).(*common.Block)
}

func (mcs *mockConsenterSupport) VerifyBlockSignature(sd []*common.SignedData, config *common.ConfigEnvelope) error {
	args := mcs.Called(sd, config)
	return args.Error(0)
}

// 添加Fabric 1.4.4中ConsenterSupport接口所需的ChannelConfig方法
func (mcs *mockConsenterSupport) ChannelConfig() channelconfig.Channel {
	args := mcs.Called()
	return args.Get(0).(channelconfig.Channel)
}

// 添加Fabric 1.4.4中ConsenterSupport接口所需的ProcessConfigUpdateMsg方法
func (mcs *mockConsenterSupport) ProcessConfigUpdateMsg(env *common.Envelope) (*common.Envelope, uint64, error) {
	args := mcs.Called(env)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).(*common.Envelope), args.Get(1).(uint64), args.Error(2)
}

// 创建测试用的交易
func createTestTransaction(configSeq uint64) *pb.Transaction {
	return &pb.Transaction{
		Envelope:  &common.Envelope{Payload: []byte("test-payload")},
		ChannalId: "test-channel",
		ConfigSeq: configSeq,
	}
}

// 创建一个基本的HotStuff实例用于测试
func setupBasicHotStuff() (*BasicHotStuff, *mockConsenterSupport) {
	mockSupport := new(mockConsenterSupport)
	bhs := &BasicHotStuff{
		support: mockSupport,
	}
	return bhs, mockSupport
}

func TestWriteLedger_NilBlock(t *testing.T) {
	bhs, _ := setupBasicHotStuff()

	// 测试空区块
	bhs.writeLedger(nil)
	// 由于函数只是记录错误并返回，我们无法直接断言其行为
	// 这里主要测试函数不会panic
}

func TestWriteLedger_NormalBlock_Success(t *testing.T) {
	bhs, mockSupport := setupBasicHotStuff()

	// 设置当前序列号
	currentSeq := uint64(10)
	mockSupport.On("Sequence").Return(currentSeq)

	// 创建两个具有不同configSeq的交易
	tx1 := createTestTransaction(currentSeq - 1) // 需要处理
	tx2 := createTestTransaction(currentSeq)     // 不需要处理

	// 设置ProcessNormalMsg的预期行为 - 成功处理
	mockSupport.On("ProcessNormalMsg", tx1.Envelope).Return(tx1.Envelope, nil)

	// 预期的创建区块行为
	expectedEnvelopes := []*common.Envelope{tx1.Envelope, tx2.Envelope}
	fabricBlock := &common.Block{Header: &common.BlockHeader{Number: 1}}
	mockSupport.On("CreateNextBlock", expectedEnvelopes).Return(fabricBlock)

	// 预期的写入区块行为
	mockSupport.On("WriteBlock", fabricBlock, mock.Anything).Return()

	// 创建包含两个交易的区块
	block := &pb.Block{
		IsNormal: true,
		Commands: []*pb.Transaction{tx1, tx2},
	}

	// 执行测试
	bhs.writeLedger(block)

	// 验证所有预期调用都发生了
	mockSupport.AssertExpectations(t)
}

func TestWriteLedger_NormalBlock_ProcessFails(t *testing.T) {
	bhs, mockSupport := setupBasicHotStuff()

	// 设置当前序列号
	currentSeq := uint64(10)
	mockSupport.On("Sequence").Return(currentSeq)

	// 创建两个交易，一个处理失败
	tx1 := createTestTransaction(currentSeq - 1) // 需要处理，但处理失败
	tx2 := createTestTransaction(currentSeq)     // 不需要处理

	// 设置tx1处理失败 - 注意这里返回的是 uint64 和 error
	mockSupport.On("ProcessNormalMsg", tx1.Envelope).Return(uint64(0), errors.New("processing failed"))

	// 由于tx1处理失败，只有tx2会被包含
	expectedEnvelopes := []*common.Envelope{tx2.Envelope}
	fabricBlock := &common.Block{Header: &common.BlockHeader{Number: 1}}
	mockSupport.On("CreateNextBlock", expectedEnvelopes).Return(fabricBlock)

	// 预期的写入区块行为
	mockSupport.On("WriteBlock", fabricBlock, mock.Anything).Return()

	// 创建包含两个交易的区块
	block := &pb.Block{
		IsNormal: true,
		Commands: []*pb.Transaction{tx1, tx2},
	}

	// 执行测试
	bhs.writeLedger(block)

	// 验证所有预期调用都发生了
	mockSupport.AssertExpectations(t)
}

func TestWriteLedger_NormalBlock_NoValidTx(t *testing.T) {
	bhs, mockSupport := setupBasicHotStuff()

	// 设置当前序列号
	currentSeq := uint64(10)
	mockSupport.On("Sequence").Return(currentSeq)

	// 创建一个将处理失败的交易
	tx := createTestTransaction(currentSeq - 1)

	// 设置tx处理失败
	mockSupport.On("ProcessNormalMsg", tx.Envelope).Return(uint64(0), errors.New("processing failed"))

	// 创建区块
	block := &pb.Block{
		IsNormal: true,
		Commands: []*pb.Transaction{tx},
	}

	// 执行测试 - 由于没有有效交易，CreateNextBlock不会被调用
	bhs.writeLedger(block)

	// 验证调用
	mockSupport.AssertExpectations(t)
	mockSupport.AssertNotCalled(t, "CreateNextBlock", mock.Anything)
	mockSupport.AssertNotCalled(t, "WriteBlock", mock.Anything, mock.Anything)
}

func TestWriteLedger_NormalBlock_EmptyTransaction(t *testing.T) {
	bhs, mockSupport := setupBasicHotStuff()

	// 设置当前序列号
	mockSupport.On("Sequence").Return(uint64(10))

	// 创建空交易和有效交易的混合区块
	nilTx := (*pb.Transaction)(nil)
	emptyTx := &pb.Transaction{Envelope: nil}
	validTx := createTestTransaction(10)

	// 创建区块
	block := &pb.Block{
		IsNormal: true,
		Commands: []*pb.Transaction{nilTx, emptyTx, validTx},
	}

	// 预期的创建和写入区块行为
	fabricBlock := &common.Block{Header: &common.BlockHeader{Number: 1}}
	mockSupport.On("CreateNextBlock", []*common.Envelope{validTx.Envelope}).Return(fabricBlock)
	mockSupport.On("WriteBlock", fabricBlock, mock.Anything).Return()

	// 执行测试
	bhs.writeLedger(block)

	// 验证调用
	mockSupport.AssertExpectations(t)
}

func TestWriteLedger_NormalBlock_CreateBlockFails(t *testing.T) {
	bhs, mockSupport := setupBasicHotStuff()

	// 设置当前序列号
	mockSupport.On("Sequence").Return(uint64(10))

	// 创建交易
	tx := createTestTransaction(10)

	// 创建区块
	block := &pb.Block{
		IsNormal: true,
		Commands: []*pb.Transaction{tx},
	}

	// 设置CreateNextBlock失败
	mockSupport.On("CreateNextBlock", []*common.Envelope{tx.Envelope}).Return(nil)

	// 执行测试
	bhs.writeLedger(block)

	// 验证调用
	mockSupport.AssertExpectations(t)
	mockSupport.AssertNotCalled(t, "WriteBlock", mock.Anything, mock.Anything)
}

func TestWriteLedger_ConfigBlock_Success(t *testing.T) {
	bhs, mockSupport := setupBasicHotStuff()

	// 设置当前序列号
	currentSeq := uint64(10)
	mockSupport.On("Sequence").Return(currentSeq)

	// 创建配置交易
	tx := createTestTransaction(currentSeq - 1) // 需要处理

	// 设置ProcessConfigMsg的预期行为 - 成功处理
	processedEnv := &common.Envelope{Payload: []byte("processed-config")}
	mockSupport.On("ProcessConfigMsg", tx.Envelope).Return(processedEnv, uint64(currentSeq), nil)

	// 预期的创建区块行为
	fabricBlock := &common.Block{Header: &common.BlockHeader{Number: 1}}
	mockSupport.On("CreateNextBlock", []*common.Envelope{processedEnv}).Return(fabricBlock)

	// 预期的写入配置区块行为
	mockSupport.On("WriteConfigBlock", fabricBlock, mock.Anything).Return()

	// 创建包含配置交易的区块
	block := &pb.Block{
		IsNormal: false,
		Commands: []*pb.Transaction{tx},
	}

	// 执行测试
	bhs.writeLedger(block)

	// 验证所有预期调用都发生了
	mockSupport.AssertExpectations(t)
}

func TestWriteLedger_ConfigBlock_NewerConfigSeq(t *testing.T) {
	bhs, mockSupport := setupBasicHotStuff()

	// 设置当前序列号
	currentSeq := uint64(10)
	mockSupport.On("Sequence").Return(currentSeq)

	// 创建配置交易，ConfigSeq比当前更新
	tx := createTestTransaction(currentSeq + 1)

	// 预期的创建区块行为 - 使用原始信封
	fabricBlock := &common.Block{Header: &common.BlockHeader{Number: 1}}
	mockSupport.On("CreateNextBlock", []*common.Envelope{tx.Envelope}).Return(fabricBlock)

	// 预期的写入配置区块行为
	mockSupport.On("WriteConfigBlock", fabricBlock, mock.Anything).Return()

	// 创建包含配置交易的区块
	block := &pb.Block{
		IsNormal: false,
		Commands: []*pb.Transaction{tx},
	}

	// 执行测试
	bhs.writeLedger(block)

	// 验证所有预期调用都发生了
	mockSupport.AssertExpectations(t)
	mockSupport.AssertNotCalled(t, "ProcessConfigMsg", mock.Anything)
}

func TestWriteLedger_ConfigBlock_ProcessConfigFails(t *testing.T) {
	bhs, mockSupport := setupBasicHotStuff()

	// 设置当前序列号
	currentSeq := uint64(10)
	mockSupport.On("Sequence").Return(currentSeq)

	// 创建配置交易
	tx := createTestTransaction(currentSeq - 1)

	// 设置ProcessConfigMsg失败
	mockSupport.On("ProcessConfigMsg", tx.Envelope).Return(nil, uint64(0), errors.New("config processing failed"))

	// 创建包含配置交易的区块
	block := &pb.Block{
		IsNormal: false,
		Commands: []*pb.Transaction{tx},
	}

	// 执行测试
	bhs.writeLedger(block)

	// 验证调用
	mockSupport.AssertExpectations(t)
	mockSupport.AssertNotCalled(t, "CreateNextBlock", mock.Anything)
	mockSupport.AssertNotCalled(t, "WriteConfigBlock", mock.Anything, mock.Anything)
}

func TestWriteLedger_ConfigBlock_NoConfigTx(t *testing.T) {
	bhs, mockSupport := setupBasicHotStuff()

	// 设置当前序列号
	mockSupport.On("Sequence").Return(uint64(10))

	// 创建空的配置区块
	block := &pb.Block{
		IsNormal: false,
		Commands: []*pb.Transaction{},
	}

	// 执行测试
	bhs.writeLedger(block)

	// 验证调用
	mockSupport.AssertNotCalled(t, "ProcessConfigMsg", mock.Anything)
	mockSupport.AssertNotCalled(t, "CreateNextBlock", mock.Anything)
	mockSupport.AssertNotCalled(t, "WriteConfigBlock", mock.Anything, mock.Anything)
}

func TestWriteLedger_ConfigBlock_EmptyConfigTx(t *testing.T) {
	bhs, mockSupport := setupBasicHotStuff()

	// 设置当前序列号
	mockSupport.On("Sequence").Return(uint64(10))

	// 创建空交易的配置区块
	nilTx := (*pb.Transaction)(nil)
	block := &pb.Block{
		IsNormal: false,
		Commands: []*pb.Transaction{nilTx},
	}

	// 执行测试
	bhs.writeLedger(block)

	// 验证调用
	mockSupport.AssertNotCalled(t, "ProcessConfigMsg", mock.Anything)
	mockSupport.AssertNotCalled(t, "CreateNextBlock", mock.Anything)
	mockSupport.AssertNotCalled(t, "WriteConfigBlock", mock.Anything, mock.Anything)
}

func TestWriteLedger_ConfigBlock_CreateBlockFails(t *testing.T) {
	bhs, mockSupport := setupBasicHotStuff()

	// 设置当前序列号
	currentSeq := uint64(10)
	mockSupport.On("Sequence").Return(currentSeq)

	// 创建配置交易
	tx := createTestTransaction(currentSeq)

	// 设置CreateNextBlock失败
	mockSupport.On("CreateNextBlock", []*common.Envelope{tx.Envelope}).Return(nil)

	// 创建包含配置交易的区块
	block := &pb.Block{
		IsNormal: false,
		Commands: []*pb.Transaction{tx},
	}

	// 执行测试
	bhs.writeLedger(block)

	// 验证调用
	mockSupport.AssertExpectations(t)
	mockSupport.AssertNotCalled(t, "WriteConfigBlock", mock.Anything, mock.Anything)
}

func TestWriteLedger_EdgeCases(t *testing.T) {
	t.Run("Normal block with nil commands", func(t *testing.T) {
		bhs, mockSupport := setupBasicHotStuff()
		mockSupport.On("Sequence").Return(uint64(10))

		block := &pb.Block{
			IsNormal: true,
			Commands: nil,
		}

		bhs.writeLedger(block)
		mockSupport.AssertNotCalled(t, "CreateNextBlock", mock.Anything)
	})

	t.Run("Config block with nil and empty envelope", func(t *testing.T) {
		bhs, mockSupport := setupBasicHotStuff()
		mockSupport.On("Sequence").Return(uint64(10))

		// 包含nil交易的区块
		block := &pb.Block{
			IsNormal: false,
			Commands: []*pb.Transaction{
				{Envelope: nil},
			},
		}

		bhs.writeLedger(block)
		mockSupport.AssertNotCalled(t, "ProcessConfigMsg", mock.Anything)
		mockSupport.AssertNotCalled(t, "CreateNextBlock", mock.Anything)
	})

	t.Run("Normal block with mixed valid and invalid transactions", func(t *testing.T) {
		bhs, mockSupport := setupBasicHotStuff()
		mockSupport.On("Sequence").Return(uint64(10))

		// 有效交易
		validTx := createTestTransaction(10)

		// 创建一个无效交易
		invalidTx := createTestTransaction(9)
		mockSupport.On("ProcessNormalMsg", invalidTx.Envelope).Return(uint64(0), errors.New("invalid tx"))

		// 创建区块
		block := &pb.Block{
			IsNormal: true,
			Commands: []*pb.Transaction{invalidTx, nil, validTx},
		}

		// 设置预期行为
		fabricBlock := &common.Block{Header: &common.BlockHeader{Number: 1}}
		mockSupport.On("CreateNextBlock", []*common.Envelope{validTx.Envelope}).Return(fabricBlock)
		mockSupport.On("WriteBlock", fabricBlock, mock.Anything).Return()

		// 执行测试
		bhs.writeLedger(block)

		// 验证调用
		mockSupport.AssertExpectations(t)
	})
}

// 测试Assertions的辅助函数
func TestAssertionHelpers(t *testing.T) {
	// 创建模拟对象
	mockSupport := new(mockConsenterSupport)

	// 测试预期的调用
	mockSupport.On("Sequence").Return(uint64(10))

	// 执行调用
	seq := mockSupport.Sequence()

	// 验证返回值
	assert.Equal(t, uint64(10), seq)

	// 验证调用发生
	mockSupport.AssertCalled(t, "Sequence")

	// 验证未发生的调用
	mockSupport.AssertNotCalled(t, "ProcessNormalMsg", mock.Anything)
}
