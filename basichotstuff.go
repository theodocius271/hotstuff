package hotstuff

import (
	"bytes"
	"encoding/json"
	"strconv"

	"github.com/hyperledger/fabric/orderer/consensus"
	"github.com/hyperledger/fabric/protos/common"
	"github.com/niclabs/tcrsa"
	"github.com/theodocius271/hotstuff/config"
	"github.com/theodocius271/hotstuff/crypto"
	"github.com/theodocius271/hotstuff/logging"
	pb "github.com/theodocius271/hotstuff/proto"
	"google.golang.org/protobuf/proto"
)

var logger = logging.GetLogger()

type BasicHotStuff struct {
	HotStuffImpl
	// sometimes, a new view msg will be processed before decide msg.
	// it may cause a bug. The parameter 'decided' is used to avoid it to happen temporarily.
	// TODO: find a better way to fix the bug.
	decided bool
	support consensus.ConsenterSupport
}

// Begin RecvMsg; TODO: Might consider mv 'go RecvMsg' to Chain.Start()
func NewBasicHotStuff(id uint32, support consensus.ConsenterSupport) *BasicHotStuff {
	msgEntrance := make(chan *pb.Msg)
	bhs := &BasicHotStuff{}
	bhs.MsgEntrance = msgEntrance
	bhs.ID = id
	bhs.View = NewView(1, 1)
	logger.Debugf("[HOTSTUFF] Init block storage, replica id: %d", id)
	bhs.BlockStorage = NewBlockStorageImpl(strconv.Itoa(int(id)))
	logger.Debugf("[HOTSTUFF] Generate genesis block")
	genesisBlock := GenerateGenesisBlock()
	err := bhs.BlockStorage.Put(genesisBlock)
	if err != nil {
		logger.Fatal("generate genesis block failed")
	}
	bhs.PrepareQC = &pb.QuorumCert{
		BlockHash: genesisBlock.Hash,
		Type:      pb.MsgType_PREPARE_VOTE,
		ViewNum:   0,
		Signature: nil,
	}
	bhs.PreCommitQC = &pb.QuorumCert{
		BlockHash: genesisBlock.Hash,
		Type:      pb.MsgType_PRECOMMIT_VOTE,
		ViewNum:   0,
		Signature: nil,
	}
	bhs.CommitQC = &pb.QuorumCert{
		BlockHash: genesisBlock.Hash,
		Type:      pb.MsgType_COMMIT_VOTE,
		ViewNum:   0,
		Signature: nil,
	}
	logger.Debugf("[HOTSTUFF] Init command set, replica id: %d", id)
	bhs.ReqSet = NewCmdSet()

	// read config
	bhs.Config = config.HotStuffConfig{}
	bhs.Config.ReadConfig()

	// init timer and stop it
	bhs.TimeChan = NewTimer(bhs.Config.Timeout)
	bhs.TimeChan.Init()

	bhs.BatchTimeChan = NewTimer(bhs.Config.BatchTimeout)
	bhs.BatchTimeChan.Init()

	bhs.CurExec = &CurProposal{
		Node:          nil,
		DocumentHash:  nil,
		PrepareVote:   make([]*tcrsa.SigShare, 0),
		PreCommitVote: make([]*tcrsa.SigShare, 0),
		CommitVote:    make([]*tcrsa.SigShare, 0),
		HighQC:        make([]*pb.QuorumCert, 0),
	}
	privateKey, err := crypto.ReadThresholdPrivateKeyFromFile(bhs.GetSelfInfo().PrivateKey)
	if err != nil {
		logger.Fatal(err)
	}
	bhs.Config.PrivateKey = privateKey
	bhs.decided = false

	bhs.support = support

	go bhs.receiveMsg()

	return bhs
}

/*
func (bhs *BasicHotStuff) RegisterChain(support consensus.ConsenterSupport) {
	if _, ok := bhs.supports[support.ChainID()]; ok {
		return
	}
	logger.Printf("[Node] Register the chain(%s)", support.ChainID())
	bhs.supports[support.ChainID()] = support
}
*/

// receiveMsg receive msg from msg channel
func (bhs *BasicHotStuff) receiveMsg() {
	for {
		select {
		case msg, ok := <-bhs.MsgEntrance:
			if ok && msg.ChannalID == bhs.support.Sequence() {
				go bhs.handleMsg(msg)
			}
		case <-bhs.TimeChan.Timeout():
			logger.Warn("Time out, goto new view")
			// set the duration of the timeout to 2 times
			bhs.TimeChan = NewTimer(bhs.Config.Timeout * 2)
			bhs.TimeChan.Init()
			bhs.ReqSet.UnMark(bhs.CurExec.Node.Commands...)
			bhs.BlockStorage.Put(bhs.CreateLeaf(bhs.CurExec.Node.ParentHash, nil, nil, true))
			bhs.View.ViewNum++
			bhs.View.Primary = bhs.GetLeader()
			// check if self is the next leader
			if bhs.GetLeader() != bhs.ID {
				// if not, send next view mag to the next leader
				newViewMsg := bhs.Msg(pb.MsgType_NEWVIEW, nil, bhs.PrepareQC)
				bhs.Unicast(bhs.GetNetworkInfo()[bhs.GetLeader()], newViewMsg)
				// clear curExec
				bhs.CurExec = NewCurProposal()
			} else {
				bhs.decided = true
			}
		case <-bhs.BatchTimeChan.Timeout():
			bhs.BatchTimeChan.Init()
			bhs.batchEvent(bhs.ReqSet.GetFirst(int(bhs.Config.BatchSize)), true)
			// bhs.batchEvent(bhs.CmdSet.GetFirst(int(bhs.support.SharedConfig().BatchSize().MaxMessageCount)))
		}
	}
}

// handleMsg handle different msg with different way
func (bhs *BasicHotStuff) handleMsg(msg *pb.Msg) {
	switch msg.Payload.(type) {
	case *pb.Msg_NewView:
		logger.Debug("[HOTSTUFF NEWVIEW] Got new view msg")
		// process highqc and node
		bhs.CurExec.HighQC = append(bhs.CurExec.HighQC, msg.GetNewView().PrepareQC)
		if bhs.decided {
			if len(bhs.CurExec.HighQC) >= 2*bhs.Config.F {
				bhs.View.ViewChanging = true
				bhs.HighQC = bhs.PrepareQC
				for _, qc := range bhs.CurExec.HighQC {
					if qc.ViewNum > bhs.HighQC.ViewNum {
						bhs.HighQC = qc
					}
				}
				// TODO sync blocks if fall behind
				bhs.CurExec = NewCurProposal()
				bhs.View.ViewChanging = false
				bhs.BatchTimeChan.SoftStartTimer()
				bhs.decided = false
			}
		}
		break
	case *pb.Msg_Prepare:
		logger.Debug("[HOTSTUFF PREPARE] Got prepare msg")
		if !bhs.MatchingMsg(msg, pb.MsgType_PREPARE) {
			logger.Warn("[HOTSTUFF PREPARE] msg does not match")
			return
		}
		prepare := msg.GetPrepare()
		if !bytes.Equal(prepare.CurProposal.ParentHash, prepare.HighQC.BlockHash) ||
			!bhs.SafeNode(prepare.CurProposal, prepare.HighQC) {
			logger.Warn("[HOTSTUFF PREPARE] node is not correct")
			return
		}
		// create prepare vote msg
		marshal, _ := proto.Marshal(msg)
		bhs.CurExec.DocumentHash, _ = crypto.CreateDocumentHash(marshal, bhs.Config.PublicKey)
		bhs.CurExec.Node = prepare.CurProposal
		partSig, _ := crypto.TSign(bhs.CurExec.DocumentHash, bhs.Config.PrivateKey, bhs.Config.PublicKey)
		partSigBytes, _ := json.Marshal(partSig)
		prepareVoteMsg := bhs.VoteMsg(pb.MsgType_PREPARE_VOTE, bhs.CurExec.Node, nil, partSigBytes)
		// send msg to leader
		bhs.Unicast(bhs.GetNetworkInfo()[bhs.GetLeader()], prepareVoteMsg)
		bhs.TimeChan.SoftStartTimer()
		break
	case *pb.Msg_PrepareVote:
		logger.Debugf("[HOTSTUFF PREPARE-VOTE] Got prepare vote msg")
		if !bhs.MatchingMsg(msg, pb.MsgType_PREPARE_VOTE) {
			logger.Warn("[HOTSTUFF PREPARE-VOTE] Msg not match")
			return
		}
		// verify
		prepareVote := msg.GetPrepareVote()
		partSig := new(tcrsa.SigShare)
		_ = json.Unmarshal(prepareVote.PartialSig, partSig)
		if err := crypto.VerifyPartSig(partSig, bhs.CurExec.DocumentHash, bhs.Config.PublicKey); err != nil {
			logger.Warn("[HOTSTUFF PREPARE-VOTE] Partial signature is not correct")
			return
		}
		// put it into preparevote slice
		bhs.CurExec.PrepareVote = append(bhs.CurExec.PrepareVote, partSig)
		if len(bhs.CurExec.PrepareVote) == bhs.Config.F*2+1 {
			// create full signature
			signature, _ := crypto.CreateFullSignature(bhs.CurExec.DocumentHash, bhs.CurExec.PrepareVote, bhs.Config.PublicKey)
			qc := bhs.QC(pb.MsgType_PREPARE_VOTE, signature, prepareVote.BlockHash)
			bhs.PrepareQC = qc
			preCommitMsg := bhs.Msg(pb.MsgType_PRECOMMIT, bhs.CurExec.Node, qc)
			// broadcast msg
			bhs.Broadcast(preCommitMsg)
			bhs.TimeChan.SoftStartTimer()
		}
		break
	case *pb.Msg_PreCommit:
		logger.Debug("[HOTSTUFF PRECOMMIT] Got precommit msg")
		if !bhs.MatchingQC(msg.GetPreCommit().PrepareQC, pb.MsgType_PREPARE_VOTE) {
			logger.Warn("[HOTSTUFF PRECOMMIT] QC not match")
			return
		}
		bhs.PrepareQC = msg.GetPreCommit().PrepareQC
		partSig, _ := crypto.TSign(bhs.CurExec.DocumentHash, bhs.Config.PrivateKey, bhs.Config.PublicKey)
		partSigBytes, _ := json.Marshal(partSig)
		preCommitVote := bhs.VoteMsg(pb.MsgType_PRECOMMIT_VOTE, bhs.CurExec.Node, nil, partSigBytes)
		bhs.Unicast(bhs.GetNetworkInfo()[bhs.GetLeader()], preCommitVote)
		bhs.TimeChan.SoftStartTimer()
		break
	case *pb.Msg_PreCommitVote:
		logger.Debug("[HOTSTUFF PRECOMMIT-VOTE] Got precommit vote msg")
		if !bhs.MatchingMsg(msg, pb.MsgType_PRECOMMIT_VOTE) {
			logger.Warn("[HOTSTUFF PRECOMMIT-VOTE] Msg not match")
			return
		}
		// verify
		preCommitVote := msg.GetPreCommitVote()
		partSig := new(tcrsa.SigShare)
		_ = json.Unmarshal(preCommitVote.PartialSig, partSig)
		if err := crypto.VerifyPartSig(partSig, bhs.CurExec.DocumentHash, bhs.Config.PublicKey); err != nil {
			logger.Warn("[HOTSTUFF PRECOMMIT-VOTE] Partial signature is not correct")
			return
		}
		bhs.CurExec.PreCommitVote = append(bhs.CurExec.PreCommitVote, partSig)
		if len(bhs.CurExec.PreCommitVote) == 2*bhs.Config.F+1 {
			signature, _ := crypto.CreateFullSignature(bhs.CurExec.DocumentHash, bhs.CurExec.PreCommitVote, bhs.Config.PublicKey)
			preCommitQC := bhs.QC(pb.MsgType_PRECOMMIT_VOTE, signature, bhs.CurExec.Node.Hash)
			// vote self
			bhs.PreCommitQC = preCommitQC
			commitMsg := bhs.Msg(pb.MsgType_COMMIT, bhs.CurExec.Node, preCommitQC)
			bhs.Broadcast(commitMsg)
			bhs.TimeChan.SoftStartTimer()
		}
		break
	case *pb.Msg_Commit:
		logger.Debug("[HOTSTUFF COMMIT] Got commit msg")
		commit := msg.GetCommit()
		if !bhs.MatchingQC(commit.PreCommitQC, pb.MsgType_PRECOMMIT_VOTE) {
			logger.Warn("[HOTSTUFF COMMIT] QC not match")
			return
		}
		bhs.PreCommitQC = commit.PreCommitQC
		partSig, _ := crypto.TSign(bhs.CurExec.DocumentHash, bhs.Config.PrivateKey, bhs.Config.PublicKey)
		partSigBytes, _ := json.Marshal(partSig)
		commitVoteMsg := bhs.VoteMsg(pb.MsgType_COMMIT_VOTE, bhs.CurExec.Node, nil, partSigBytes)
		bhs.Unicast(bhs.GetNetworkInfo()[bhs.GetLeader()], commitVoteMsg)
		bhs.TimeChan.SoftStartTimer()
		break
	case *pb.Msg_CommitVote:
		logger.Debug("[HOTSTUFF COMMIT-VOTE] Got commit vote msg")
		if !bhs.MatchingMsg(msg, pb.MsgType_COMMIT_VOTE) {
			logger.Warn("[HOTSTUFF COMMIT-VOTE] Msg not match")
			return
		}
		commitVoteMsg := msg.GetCommitVote()
		partSig := new(tcrsa.SigShare)
		_ = json.Unmarshal(commitVoteMsg.PartialSig, partSig)
		if err := crypto.VerifyPartSig(partSig, bhs.CurExec.DocumentHash, bhs.Config.PublicKey); err != nil {
			logger.Warn("[HOTSTUFF COMMIT-VOTE] Partial signature is not correct")
			return
		}
		bhs.CurExec.CommitVote = append(bhs.CurExec.CommitVote, partSig)
		if len(bhs.CurExec.CommitVote) == 2*bhs.Config.F+1 {
			signature, _ := crypto.CreateFullSignature(bhs.CurExec.DocumentHash, bhs.CurExec.CommitVote, bhs.Config.PublicKey)
			commitQC := bhs.QC(pb.MsgType_COMMIT_VOTE, signature, bhs.CurExec.Node.Hash)
			// vote self
			bhs.CommitQC = commitQC
			decideMsg := bhs.Msg(pb.MsgType_DECIDE, bhs.CurExec.Node, commitQC)
			bhs.Broadcast(decideMsg)
			bhs.TimeChan.Stop()
			bhs.processProposal()
		}
		break
	case *pb.Msg_Decide:
		logger.Debug("[HOTSTUFF DECIDE] Got decide msg")
		decideMsg := msg.GetDecide()
		if !bhs.MatchingQC(decideMsg.CommitQC, pb.MsgType_COMMIT_VOTE) {
			logger.Warn("[HOTSTUFF DECIDE] QC not match")
			return
		}
		bhs.CommitQC = decideMsg.CommitQC
		bhs.TimeChan.Stop()
		bhs.processProposal()
		break
	case *pb.Msg_Request:
		request := msg.GetRequest()
		logger.Debugf("[HOTSTUFF] Got request msg, content:%s", request.String())
		// put the cmd into the cmdset
		bhs.ReqSet.Add(request.Transaction)
		// send request to the leader, if the replica is not the leader
		if bhs.ID != bhs.GetLeader() {
			bhs.Unicast(bhs.GetNetworkInfo()[bhs.GetLeader()], msg)
			return
		}
		if bhs.CurExec.Node != nil || bhs.View.ViewChanging {
			return
		}

		if !request.IsNormal {
			logger.Debugf("processing config transaction")
			bhs.batchEvent(bhs.ReqSet.GetFirst(1), false)
			break
		}

		// start batch timer
		bhs.BatchTimeChan.SoftStartTimer()
		// if the length of unprocessed cmd equals to batch size, stop timer and call handleMsg to send prepare msg
		logger.Debugf("cmd set size: %d", len(bhs.ReqSet.GetFirst(int(bhs.Config.BatchSize))))
		cmds := bhs.ReqSet.GetFirst(int(bhs.Config.BatchSize))
		if len(cmds) == int(bhs.Config.BatchSize) {
			// stop timer
			bhs.BatchTimeChan.Stop()
			// create prepare msg
			bhs.batchEvent(cmds, true)
		}
		break
	default:
		logger.Warn("Unsupported msg type, drop it.")
		break
	}
}

func (bhs *BasicHotStuff) processProposal() {

	// process proposal, deprecated
	// go bhs.ProcessProposal(bhs.CurExec.Node.Commands, bhs.support)
	go bhs.writeLedger(bhs.CurExec.Node)
	// store block
	// used to have goroutine
	bhs.CurExec.Node.Committed = true
	bhs.BlockStorage.Put(bhs.CurExec.Node)

	// add view number
	bhs.View.ViewNum++
	bhs.View.Primary = bhs.GetLeader()
	// check if self is the next leader
	if bhs.View.Primary != bhs.ID {
		// if not, send next view mag to the next leader
		newViewMsg := bhs.Msg(pb.MsgType_NEWVIEW, nil, bhs.PrepareQC)
		bhs.Unicast(bhs.GetNetworkInfo()[bhs.GetLeader()], newViewMsg)
		// clear curExec
		bhs.CurExec = NewCurProposal()
	} else {
		bhs.decided = true
	}

}

func (bhs *BasicHotStuff) writeLedger(block *pb.Block) {
	if block == nil {
		logger.Errorf("block is nil")
		return
	}

	batch := make([]*common.Envelope, 0)
	transactions := block.Commands

	for _, tx := range transactions {
		if tx.isNormal {
			// Normal Transaction
			batch = append(batch, tx.Envelope)
		} else {
			// Config Message
			if batch != nil {
				fabricNormalBlock := bhs.support.CreateNextBlock(batch)
				bhs.support.WriteBlock(fabricNormalBlock, nil)
				batch = make([]*common.Envelope, 0)
			}
			fabricConfigBlock := bhs.support.CreateNextBlock([]*common.Envelope{tx.Envelope})
			bhs.support.WriteConfigBlock(fabricConfigBlock, nil)
		}
	}

	if batch != nil {
		fabricNormalBlock := bhs.support.CreateNextBlock(batch)
		bhs.support.WriteBlock(fabricNormalBlock, nil)
	}

}

func (bhs *BasicHotStuff) batchEvent(cmds []*pb.Transaction, isNormal bool) {
	if len(cmds) == 0 {
		bhs.BatchTimeChan.SoftStartTimer()
		return
	}
	// create prepare msg
	node := bhs.CreateLeaf(bhs.BlockStorage.GetLastBlockHash(), cmds, nil, isNormal)
	bhs.CurExec.Node = node
	bhs.ReqSet.MarkProposed(cmds...)
	if bhs.HighQC == nil {
		bhs.HighQC = bhs.PrepareQC
	}
	prepareMsg := bhs.Msg(pb.MsgType_PREPARE, node, bhs.HighQC)
	// vote self
	marshal, _ := proto.Marshal(prepareMsg)
	bhs.CurExec.DocumentHash, _ = crypto.CreateDocumentHash(marshal, bhs.Config.PublicKey)
	partSig, _ := crypto.TSign(bhs.CurExec.DocumentHash, bhs.Config.PrivateKey, bhs.Config.PublicKey)
	bhs.CurExec.PrepareVote = append(bhs.CurExec.PrepareVote, partSig)
	// broadcast prepare msg
	bhs.Broadcast(prepareMsg)
	bhs.TimeChan.SoftStartTimer()
}
