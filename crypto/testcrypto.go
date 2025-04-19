package crypto

import (
	"testing"

	"github.com/niclabs/tcrsa"
)

// GenerateTestKeys 生成测试用的门限签名密钥
func GenerateTestKeys(t *testing.T) (*tcrsa.KeyShare, *tcrsa.KeyMeta) {
	// 为测试生成简单参数的门限密钥
	k := 2 // 最小签名数
	n := 4 // 总节点数

	// 生成密钥
	keyShares, keyMeta, err := tcrsa.NewKey(512, uint16(k), uint16(n), nil)
	if err != nil {
		t.Fatalf("Failed to generate test keys: %v", err)
	}

	// 返回第一个密钥份额和公共密钥元数据
	return keyShares[0], keyMeta
}
