// Expected EnvVar
/*
version: '2'
services:
  orderer1.example.com:
    container_name: orderer1.example.com
    image: hyperledger/fabric-orderer:1.4.4
    environment:
      - ORDERER_GENERAL_LOGLEVEL=INFO
      - ORDERER_GENERAL_LISTENADDRESS=0.0.0.0
      - ORDERER_GENERAL_LISTENPORT=7050
      - ORDERER_GENERAL_GENESISMETHOD=file
      - ORDERER_GENERAL_GENESISFILE=/var/hyperledger/orderer/orderer.genesis.block
      - ORDERER_GENERAL_LOCALMSPID=OrdererMSP
      - ORDERER_GENERAL_LOCALMSPDIR=/var/hyperledger/orderer/msp
      # HotStuff 特定的配置
      - ORDERER_GENERAL_GENESISPROFILE=SampleHotStuffConsortium
      - ORDERER_GENERAL_LEDGERTYPE=file
      # 改为 hotstuff
      - ORDERER_GENERAL_CONSENSUSTYPE=hotstuff
      # HotStuff 环境变量
      - ORDERER_HOTSTUFF_NODEID=1
      - ORDERER_HOTSTUFF_CLUSTER=4:7050,1:7050,2:7050,3:7050
      - ORDERER_HOTSTUFF_TIMEOUT=3000
      - ORDERER_HOTSTUFF_BATCHTIMEOUT=6000
      - ORDERER_HOTSTUFF_PUBLICKEYPATH=base64key1,base64key2,base64key3,base64key4
      - ORDERER_HOTSTUFF_PRIVATEKEYPATH=base64key1,base64key2,base64key3,base64key4

*/

package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/niclabs/tcrsa"
	"github.com/theodocius271/hotstuff/crypto"
	"github.com/theodocius271/hotstuff/logging"
)

var logger = logging.GetLogger()

type HotStuffConfig struct {
	BatchTimeout time.Duration
	Timeout      time.Duration
	PublicKey    *tcrsa.KeyMeta
	PrivateKey   *tcrsa.KeyShare
	Cluster      []*ReplicaInfo
	N            int
	F            int
}

type ReplicaInfo struct {
	ID         uint32
	Address    string
	PrivateKey string
}

func NewHotStuffConfig() *HotStuffConfig {
	config := &HotStuffConfig{}
	config.ReadConfig()
	return config
}

func (config *HotStuffConfig) ReadConfig() {
	logger.Debug("[HOTSTUFF] Read config")

	timeoutStr := os.Getenv("ORDERER_HOTSTUFF_TIMEOUT")
	if timeoutStr == "" {
		logger.Infof("ORDERER_HOTSTUFF_TIMEOUT environment variable not set")
		return
	}
	timeout, err := strconv.Atoi(timeoutStr)
	if err == nil {
		config.Timeout = time.Duration(timeout) * time.Millisecond
	} else {
		config.Timeout = 2000 * time.Millisecond
	}

	batchTimeoutStr := os.Getenv("ORDERER_HOTSTUFF_TIMEOUT")
	if batchTimeoutStr == "" {
		logger.Infof("ORDERER_HOTSTUFF_TIMEOUT environment variable not set")
		return
	}
	batchTimeout, err := strconv.Atoi(batchTimeoutStr)
	if err == nil {
		config.BatchTimeout = time.Duration(batchTimeout) * time.Millisecond
	} else {
		config.BatchTimeout = 1000 * time.Millisecond
	}

	publicKeyPath := os.Getenv("ORDERER_HOTSTUFF_PUBLICKEYPATH")
	if publicKeyPath == "" {
		logger.Infof("ORDERER_HOTSTUFF_PUBLICKEYPATH environment variable not set")
		return
	}
	publicKey, err := crypto.ReadThresholdPublicKeyFromFile(publicKeyPath)
	if err != nil {
		logger.Fatal(err)
	}
	config.PublicKey = publicKey

	clusterStr := os.Getenv("ORDERER_HOTSTUFF_CLUSTER")
	if clusterStr == "" {
		logger.Infof("ORDERER_HOTSTUFF_CLUSTER environment variable not set")
		return
	}
	replicaPairs := strings.Split(clusterStr, ",")
	replicas := make([]*ReplicaInfo, 0, len(replicaPairs))
	for _, pair := range replicaPairs {
		parts := strings.Split(pair, ":")
		if len(parts) != 2 {
			logger.Infof("invalid replica format: %s, expected format 'id:address'", pair)
			return
		}
		id, err := strconv.ParseUint(parts[0], 10, 32)
		if err != nil {
			logger.Infof("invalid replica ID: %s: %v", parts[0], err)
			return
		}
		replica := &ReplicaInfo{
			ID:      uint32(id),
			Address: parts[1],
		}
		replicas = append(replicas, replica)
	}

	config.Cluster = replicas
	config.N = len(replicas)
	config.F = (config.N - 1) / 3

}
