package mocked

import (
	"github.com/hyperledger/fabric/common/channelconfig"
	"github.com/hyperledger/fabric/orderer/common/blockcutter"
	"github.com/hyperledger/fabric/orderer/common/msgprocessor"
	"github.com/hyperledger/fabric/protos/common"
	"github.com/stretchr/testify/mock"
)

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
