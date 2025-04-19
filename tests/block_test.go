package tests

/*
import (
	"bytes"
	"crypto/sha256"
	"os"
	"path/filepath"
	"testing"

	"github.com/hyperledger/fabric/protos/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pb "github.com/theodocius271/hotstuff/proto"
)

// 创建测试用的临时数据库目录
func setupTestDB(t *testing.T) string {
	tempDir, err := os.MkdirTemp("", "hotstuff-test")
	require.NoError(t, err)

	// 确保测试结束后清理临时文件
	t.Cleanup(func() {
		os.RemoveAll(tempDir)
	})

	// 创建 dbfile 子目录
	dbDir := filepath.Join(tempDir, "dbfile")
	err = os.MkdirAll(dbDir, 0755)
	require.NoError(t, err)

	// 覆盖默认路径，重新定向到临时目录
	origDir := "/opt/hotstuff/dbfile"
	os.Symlink(dbDir, origDir)
	t.Cleanup(func() {
		os.Remove(origDir)
	})

	return tempDir
}

// 创建测试用的区块
func createTestBlock(height uint64, parentHash []byte, commands []*pb.Transaction, justify *pb.QuorumCert) *pb.Block {
	block := &pb.Block{
		ParentHash: parentHash,
		Height:     height,
		Commands:   commands,
		Justify:    justify,
		IsNormal:   true,
	}
	// 计算并设置区块哈希
	block.Hash = Hash(block)
	return block
}

// 创建测试用的交易
func createTestTransaction() *pb.Transaction {
	return &pb.Transaction{
		Envelope:  &common.Envelope{Payload: []byte("test-payload")},
		ChannalId: "test-channel",
		ConfigSeq: 1,
	}
}

// 创建测试用的证书
func createTestQuorumCert(blockHash []byte) *pb.QuorumCert {
	return &pb.QuorumCert{
		BlockHash: blockHash,
		ViewNum:   1,
		Type:      pb.MsgType_PREPARE,
		Signature: []byte("test-signature"),
	}
}

func TestHash(t *testing.T) {
	// 创建一个简单的区块
	parentHash := []byte("parent-hash")
	tx := createTestTransaction()
	qc := createTestQuorumCert([]byte("some-block-hash"))

	block := createTestBlock(1, parentHash, []*pb.Transaction{tx}, qc)

	// 清除哈希以重新计算
	originalHash := block.Hash
	block.Hash = nil

	// 重新计算哈希
	newHash := Hash(block)

	// 验证哈希一致性
	assert.NotNil(t, newHash)
	assert.Equal(t, originalHash, newHash)

	// 验证如果区块内容改变，哈希也会改变
	block.Height = 2
	block.Hash = nil
	changedHash := Hash(block)
	assert.NotEqual(t, originalHash, changedHash)
}

func TestString(t *testing.T) {
	parentHash := []byte("parent-hash")
	blockHash := []byte("block-hash")

	block := &pb.Block{
		ParentHash: parentHash,
		Hash:       blockHash,
		Height:     42,
	}

	blockStr := String(block)

	// 验证字符串表示包含必要信息
	assert.Contains(t, blockStr, "ParentHash:")
	assert.Contains(t, blockStr, "Hash:")
	assert.Contains(t, blockStr, "Height: 42")
}

func TestNewBlockStorageImpl(t *testing.T) {
	// 设置临时测试目录
	setupTestDB(t)

	// 创建一个新的区块存储
	blockStorage := NewBlockStorageImpl("1")
	defer blockStorage.Close()

	// 验证存储已正确初始化
	assert.NotNil(t, blockStorage.db)
	assert.Nil(t, blockStorage.Tip)
}

func TestBlockStorageImpl_PutAndGet(t *testing.T) {
	// 设置临时测试目录
	setupTestDB(t)

	// 创建一个新的区块存储
	blockStorage := NewBlockStorageImpl("2")
	defer blockStorage.Close()

	// 创建测试区块
	parentHash := sha256.Sum256([]byte("parent"))
	tx := createTestTransaction()
	qc := createTestQuorumCert([]byte("some-block-hash"))

	block := createTestBlock(1, parentHash[:], []*pb.Transaction{tx}, qc)

	// 保存区块
	err := blockStorage.Put(block)
	require.NoError(t, err)

	// 验证 Tip 已更新
	assert.Equal(t, block.Hash, blockStorage.Tip)

	// 获取区块并验证
	retrievedBlock, err := blockStorage.Get(block.Hash)
	require.NoError(t, err)

	assert.Equal(t, block.Height, retrievedBlock.Height)
	assert.Equal(t, block.ParentHash, retrievedBlock.ParentHash)
	assert.True(t, bytes.Equal(block.Hash, retrievedBlock.Hash))
	assert.Equal(t, len(block.Commands), len(retrievedBlock.Commands))
}

func TestBlockStorageImpl_BlockOf(t *testing.T) {
	// 设置临时测试目录
	setupTestDB(t)

	// 创建一个新的区块存储
	blockStorage := NewBlockStorageImpl("3")
	defer blockStorage.Close()

	// 创建测试区块
	parentHash := sha256.Sum256([]byte("parent"))
	tx := createTestTransaction()
	qc := createTestQuorumCert([]byte("some-block-hash"))

	block := createTestBlock(1, parentHash[:], []*pb.Transaction{tx}, qc)

	// 保存区块
	err := blockStorage.Put(block)
	require.NoError(t, err)

	// 创建一个引用该区块的证书
	cert := createTestQuorumCert(block.Hash)

	// 通过证书获取区块
	retrievedBlock, err := blockStorage.BlockOf(cert)
	require.NoError(t, err)

	// 验证获取的是正确的区块
	assert.Equal(t, block.Height, retrievedBlock.Height)
	assert.True(t, bytes.Equal(block.Hash, retrievedBlock.Hash))
}

func TestBlockStorageImpl_ParentOf(t *testing.T) {
	// 设置临时测试目录
	setupTestDB(t)

	// 创建一个新的区块存储
	blockStorage := NewBlockStorageImpl("4")
	defer blockStorage.Close()

	// 创建父区块
	parentHash := sha256.Sum256([]byte("grandparent"))
	parentQC := createTestQuorumCert([]byte("grandparent-block-hash"))
	parentTx := createTestTransaction()

	parentBlock := createTestBlock(1, parentHash[:], []*pb.Transaction{parentTx}, parentQC)

	// 保存父区块
	err := blockStorage.Put(parentBlock)
	require.NoError(t, err)

	// 创建子区块，引用父区块
	childQC := createTestQuorumCert(parentBlock.Hash)
	childTx := createTestTransaction()

	childBlock := createTestBlock(2, parentBlock.Hash, []*pb.Transaction{childTx}, childQC)

	// 保存子区块
	err = blockStorage.Put(childBlock)
	require.NoError(t, err)

	// 获取子区块的父区块
	retrievedParent, err := blockStorage.ParentOf(childBlock)
	require.NoError(t, err)

	// 验证获取的是正确的父区块
	assert.Equal(t, parentBlock.Height, retrievedParent.Height)
	assert.True(t, bytes.Equal(parentBlock.Hash, retrievedParent.Hash))
}

func TestBlockStorageImpl_GetLastBlockHash(t *testing.T) {
	// 设置临时测试目录
	setupTestDB(t)

	// 创建一个新的区块存储
	blockStorage := NewBlockStorageImpl("5")
	defer blockStorage.Close()

	// 初始时，Tip 应为 nil
	assert.Nil(t, blockStorage.GetLastBlockHash())

	// 添加一个区块
	parentHash := sha256.Sum256([]byte("parent"))
	tx := createTestTransaction()
	qc := createTestQuorumCert([]byte("some-block-hash"))

	block1 := createTestBlock(1, parentHash[:], []*pb.Transaction{tx}, qc)

	err := blockStorage.Put(block1)
	require.NoError(t, err)

	// Tip 应该更新为新区块的哈希
	assert.Equal(t, block1.Hash, blockStorage.GetLastBlockHash())

	// 添加另一个区块
	block2 := createTestBlock(2, block1.Hash, []*pb.Transaction{tx}, qc)

	err = blockStorage.Put(block2)
	require.NoError(t, err)

	// Tip 应该更新为最新区块的哈希
	assert.Equal(t, block2.Hash, blockStorage.GetLastBlockHash())
}

func TestBlockStorageImpl_RestoreStatus(t *testing.T) {
	// 设置临时测试目录
	setupTestDB(t)

	// 创建区块存储并添加一些区块
	{
		blockStorage := NewBlockStorageImpl("6")

		parentHash := sha256.Sum256([]byte("parent"))
		tx := createTestTransaction()
		qc := createTestQuorumCert([]byte("some-block-hash"))

		block := createTestBlock(1, parentHash[:], []*pb.Transaction{tx}, qc)

		err := blockStorage.Put(block)
		require.NoError(t, err)

		// 保存最后区块哈希
		lastHash := blockStorage.GetLastBlockHash()

		blockStorage.Close()

		// 再次打开存储并恢复状态
		blockStorage2 := NewBlockStorageImpl("6")
		blockStorage2.RestoreStatus()

		// 验证 Tip 已正确恢复
		assert.Equal(t, lastHash, blockStorage2.GetLastBlockHash())

		blockStorage2.Close()
	}
}

func TestBlockStorageImpl_UpdateState(t *testing.T) {
	// 设置临时测试目录
	setupTestDB(t)

	// 创建一个新的区块存储
	blockStorage := NewBlockStorageImpl("7")
	defer blockStorage.Close()

	// 创建测试区块
	parentHash := sha256.Sum256([]byte("parent"))
	tx := createTestTransaction()
	qc := createTestQuorumCert([]byte("some-block-hash"))

	block := createTestBlock(1, parentHash[:], []*pb.Transaction{tx}, qc)

	// 保存区块
	err := blockStorage.Put(block)
	require.NoError(t, err)

	// 验证初始状态
	retrievedBlock, err := blockStorage.Get(block.Hash)
	require.NoError(t, err)
	assert.False(t, retrievedBlock.Committed)

	// 更新状态为已提交
	block.Committed = true
	err = blockStorage.UpdateState(block)
	require.NoError(t, err)

	// 验证状态已更新
	retrievedBlock, err = blockStorage.Get(block.Hash)
	require.NoError(t, err)
	assert.True(t, retrievedBlock.Committed)

	// 测试空区块
	err = blockStorage.UpdateState(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "block is null")

	// 测试无哈希区块
	invalidBlock := &pb.Block{Height: 2}
	err = blockStorage.UpdateState(invalidBlock)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "block is null")
}

func TestBlockStorageImpl_GetNonExistentBlock(t *testing.T) {
	// 设置临时测试目录
	setupTestDB(t)

	// 创建一个新的区块存储
	blockStorage := NewBlockStorageImpl("8")
	defer blockStorage.Close()

	// 尝试获取不存在的区块
	nonExistentHash := sha256.Sum256([]byte("non-existent"))
	block, err := blockStorage.Get(nonExistentHash[:])

	// 应返回错误
	assert.Error(t, err)
	assert.Nil(t, block)
}

func TestBlockStorageImpl_BlockOfNonExistentCert(t *testing.T) {
	// 设置临时测试目录
	setupTestDB(t)

	// 创建一个新的区块存储
	blockStorage := NewBlockStorageImpl("9")
	defer blockStorage.Close()

	// 创建引用不存在区块的证书
	nonExistentHash := sha256.Sum256([]byte("non-existent"))
	cert := createTestQuorumCert(nonExistentHash[:])

	// 尝试通过证书获取不存在的区块
	block, err := blockStorage.BlockOf(cert)

	// 应返回错误
	assert.Error(t, err)
	assert.Nil(t, block)
}
*/
