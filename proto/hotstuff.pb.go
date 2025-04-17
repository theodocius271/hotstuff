// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v3.12.4
// source: orderer/consensus/hotstuff/proto/hotstuff.proto

package proto

import (
	common "github.com/hyperledger/fabric/protos/common"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Msg struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Payload:
	//
	//	*Msg_Prepare
	//	*Msg_PrepareVote
	//	*Msg_PreCommit
	//	*Msg_PreCommitVote
	//	*Msg_Commit
	//	*Msg_CommitVote
	//	*Msg_Decide
	//	*Msg_NewView
	//	*Msg_Request
	//	*Msg_Reply
	Payload       isMsg_Payload `protobuf_oneof:"Payload"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Msg) Reset() {
	*x = Msg{}
	mi := &file_orderer_consensus_hotstuff_proto_hotstuff_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Msg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Msg) ProtoMessage() {}

func (x *Msg) ProtoReflect() protoreflect.Message {
	mi := &file_orderer_consensus_hotstuff_proto_hotstuff_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Msg.ProtoReflect.Descriptor instead.
func (*Msg) Descriptor() ([]byte, []int) {
	return file_orderer_consensus_hotstuff_proto_hotstuff_proto_rawDescGZIP(), []int{0}
}

func (x *Msg) GetPayload() isMsg_Payload {
	if x != nil {
		return x.Payload
	}
	return nil
}

func (x *Msg) GetPrepare() *Prepare {
	if x != nil {
		if x, ok := x.Payload.(*Msg_Prepare); ok {
			return x.Prepare
		}
	}
	return nil
}

func (x *Msg) GetPrepareVote() *PrepareVote {
	if x != nil {
		if x, ok := x.Payload.(*Msg_PrepareVote); ok {
			return x.PrepareVote
		}
	}
	return nil
}

func (x *Msg) GetPreCommit() *PreCommit {
	if x != nil {
		if x, ok := x.Payload.(*Msg_PreCommit); ok {
			return x.PreCommit
		}
	}
	return nil
}

func (x *Msg) GetPreCommitVote() *PreCommitVote {
	if x != nil {
		if x, ok := x.Payload.(*Msg_PreCommitVote); ok {
			return x.PreCommitVote
		}
	}
	return nil
}

func (x *Msg) GetCommit() *Commit {
	if x != nil {
		if x, ok := x.Payload.(*Msg_Commit); ok {
			return x.Commit
		}
	}
	return nil
}

func (x *Msg) GetCommitVote() *CommitVote {
	if x != nil {
		if x, ok := x.Payload.(*Msg_CommitVote); ok {
			return x.CommitVote
		}
	}
	return nil
}

func (x *Msg) GetDecide() *Decide {
	if x != nil {
		if x, ok := x.Payload.(*Msg_Decide); ok {
			return x.Decide
		}
	}
	return nil
}

func (x *Msg) GetNewView() *NewView {
	if x != nil {
		if x, ok := x.Payload.(*Msg_NewView); ok {
			return x.NewView
		}
	}
	return nil
}

func (x *Msg) GetRequest() *Request {
	if x != nil {
		if x, ok := x.Payload.(*Msg_Request); ok {
			return x.Request
		}
	}
	return nil
}

func (x *Msg) GetReply() *Reply {
	if x != nil {
		if x, ok := x.Payload.(*Msg_Reply); ok {
			return x.Reply
		}
	}
	return nil
}

type isMsg_Payload interface {
	isMsg_Payload()
}

type Msg_Prepare struct {
	Prepare *Prepare `protobuf:"bytes,1,opt,name=prepare,proto3,oneof"`
}

type Msg_PrepareVote struct {
	PrepareVote *PrepareVote `protobuf:"bytes,2,opt,name=prepareVote,proto3,oneof"`
}

type Msg_PreCommit struct {
	PreCommit *PreCommit `protobuf:"bytes,3,opt,name=preCommit,proto3,oneof"`
}

type Msg_PreCommitVote struct {
	PreCommitVote *PreCommitVote `protobuf:"bytes,4,opt,name=preCommitVote,proto3,oneof"`
}

type Msg_Commit struct {
	Commit *Commit `protobuf:"bytes,5,opt,name=commit,proto3,oneof"`
}

type Msg_CommitVote struct {
	CommitVote *CommitVote `protobuf:"bytes,6,opt,name=commitVote,proto3,oneof"`
}

type Msg_Decide struct {
	Decide *Decide `protobuf:"bytes,7,opt,name=decide,proto3,oneof"`
}

type Msg_NewView struct {
	NewView *NewView `protobuf:"bytes,8,opt,name=newView,proto3,oneof"`
}

type Msg_Request struct {
	Request *Request `protobuf:"bytes,9,opt,name=request,proto3,oneof"`
}

type Msg_Reply struct {
	Reply *Reply `protobuf:"bytes,10,opt,name=reply,proto3,oneof"`
}

func (*Msg_Prepare) isMsg_Payload() {}

func (*Msg_PrepareVote) isMsg_Payload() {}

func (*Msg_PreCommit) isMsg_Payload() {}

func (*Msg_PreCommitVote) isMsg_Payload() {}

func (*Msg_Commit) isMsg_Payload() {}

func (*Msg_CommitVote) isMsg_Payload() {}

func (*Msg_Decide) isMsg_Payload() {}

func (*Msg_NewView) isMsg_Payload() {}

func (*Msg_Request) isMsg_Payload() {}

func (*Msg_Reply) isMsg_Payload() {}

type Empty struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Empty) Reset() {
	*x = Empty{}
	mi := &file_orderer_consensus_hotstuff_proto_hotstuff_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_orderer_consensus_hotstuff_proto_hotstuff_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_orderer_consensus_hotstuff_proto_hotstuff_proto_rawDescGZIP(), []int{1}
}

type Prepare struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	CurProposal   *Block                 `protobuf:"bytes,1,opt,name=curProposal,proto3" json:"curProposal,omitempty"`
	HighQC        *QuorumCert            `protobuf:"bytes,2,opt,name=highQC,proto3" json:"highQC,omitempty"`
	ViewNum       uint64                 `protobuf:"varint,3,opt,name=viewNum,proto3" json:"viewNum,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Prepare) Reset() {
	*x = Prepare{}
	mi := &file_orderer_consensus_hotstuff_proto_hotstuff_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Prepare) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Prepare) ProtoMessage() {}

func (x *Prepare) ProtoReflect() protoreflect.Message {
	mi := &file_orderer_consensus_hotstuff_proto_hotstuff_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Prepare.ProtoReflect.Descriptor instead.
func (*Prepare) Descriptor() ([]byte, []int) {
	return file_orderer_consensus_hotstuff_proto_hotstuff_proto_rawDescGZIP(), []int{2}
}

func (x *Prepare) GetCurProposal() *Block {
	if x != nil {
		return x.CurProposal
	}
	return nil
}

func (x *Prepare) GetHighQC() *QuorumCert {
	if x != nil {
		return x.HighQC
	}
	return nil
}

func (x *Prepare) GetViewNum() uint64 {
	if x != nil {
		return x.ViewNum
	}
	return 0
}

type PrepareVote struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	BlockHash     []byte                 `protobuf:"bytes,1,opt,name=blockHash,proto3" json:"blockHash,omitempty"`
	Qc            *QuorumCert            `protobuf:"bytes,2,opt,name=qc,proto3" json:"qc,omitempty"`
	PartialSig    []byte                 `protobuf:"bytes,3,opt,name=partialSig,proto3" json:"partialSig,omitempty"`
	ViewNum       uint64                 `protobuf:"varint,4,opt,name=viewNum,proto3" json:"viewNum,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PrepareVote) Reset() {
	*x = PrepareVote{}
	mi := &file_orderer_consensus_hotstuff_proto_hotstuff_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PrepareVote) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PrepareVote) ProtoMessage() {}

func (x *PrepareVote) ProtoReflect() protoreflect.Message {
	mi := &file_orderer_consensus_hotstuff_proto_hotstuff_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PrepareVote.ProtoReflect.Descriptor instead.
func (*PrepareVote) Descriptor() ([]byte, []int) {
	return file_orderer_consensus_hotstuff_proto_hotstuff_proto_rawDescGZIP(), []int{3}
}

func (x *PrepareVote) GetBlockHash() []byte {
	if x != nil {
		return x.BlockHash
	}
	return nil
}

func (x *PrepareVote) GetQc() *QuorumCert {
	if x != nil {
		return x.Qc
	}
	return nil
}

func (x *PrepareVote) GetPartialSig() []byte {
	if x != nil {
		return x.PartialSig
	}
	return nil
}

func (x *PrepareVote) GetViewNum() uint64 {
	if x != nil {
		return x.ViewNum
	}
	return 0
}

type PreCommit struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	PrepareQC     *QuorumCert            `protobuf:"bytes,1,opt,name=prepareQC,proto3" json:"prepareQC,omitempty"`
	ViewNum       uint64                 `protobuf:"varint,2,opt,name=viewNum,proto3" json:"viewNum,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PreCommit) Reset() {
	*x = PreCommit{}
	mi := &file_orderer_consensus_hotstuff_proto_hotstuff_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PreCommit) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PreCommit) ProtoMessage() {}

func (x *PreCommit) ProtoReflect() protoreflect.Message {
	mi := &file_orderer_consensus_hotstuff_proto_hotstuff_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PreCommit.ProtoReflect.Descriptor instead.
func (*PreCommit) Descriptor() ([]byte, []int) {
	return file_orderer_consensus_hotstuff_proto_hotstuff_proto_rawDescGZIP(), []int{4}
}

func (x *PreCommit) GetPrepareQC() *QuorumCert {
	if x != nil {
		return x.PrepareQC
	}
	return nil
}

func (x *PreCommit) GetViewNum() uint64 {
	if x != nil {
		return x.ViewNum
	}
	return 0
}

type PreCommitVote struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	BlockHash     []byte                 `protobuf:"bytes,1,opt,name=blockHash,proto3" json:"blockHash,omitempty"`
	Qc            *QuorumCert            `protobuf:"bytes,2,opt,name=qc,proto3" json:"qc,omitempty"`
	PartialSig    []byte                 `protobuf:"bytes,3,opt,name=partialSig,proto3" json:"partialSig,omitempty"`
	ViewNum       uint64                 `protobuf:"varint,4,opt,name=viewNum,proto3" json:"viewNum,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PreCommitVote) Reset() {
	*x = PreCommitVote{}
	mi := &file_orderer_consensus_hotstuff_proto_hotstuff_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PreCommitVote) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PreCommitVote) ProtoMessage() {}

func (x *PreCommitVote) ProtoReflect() protoreflect.Message {
	mi := &file_orderer_consensus_hotstuff_proto_hotstuff_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PreCommitVote.ProtoReflect.Descriptor instead.
func (*PreCommitVote) Descriptor() ([]byte, []int) {
	return file_orderer_consensus_hotstuff_proto_hotstuff_proto_rawDescGZIP(), []int{5}
}

func (x *PreCommitVote) GetBlockHash() []byte {
	if x != nil {
		return x.BlockHash
	}
	return nil
}

func (x *PreCommitVote) GetQc() *QuorumCert {
	if x != nil {
		return x.Qc
	}
	return nil
}

func (x *PreCommitVote) GetPartialSig() []byte {
	if x != nil {
		return x.PartialSig
	}
	return nil
}

func (x *PreCommitVote) GetViewNum() uint64 {
	if x != nil {
		return x.ViewNum
	}
	return 0
}

type Commit struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	PreCommitQC   *QuorumCert            `protobuf:"bytes,1,opt,name=preCommitQC,proto3" json:"preCommitQC,omitempty"`
	ViewNum       uint64                 `protobuf:"varint,2,opt,name=viewNum,proto3" json:"viewNum,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Commit) Reset() {
	*x = Commit{}
	mi := &file_orderer_consensus_hotstuff_proto_hotstuff_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Commit) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Commit) ProtoMessage() {}

func (x *Commit) ProtoReflect() protoreflect.Message {
	mi := &file_orderer_consensus_hotstuff_proto_hotstuff_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Commit.ProtoReflect.Descriptor instead.
func (*Commit) Descriptor() ([]byte, []int) {
	return file_orderer_consensus_hotstuff_proto_hotstuff_proto_rawDescGZIP(), []int{6}
}

func (x *Commit) GetPreCommitQC() *QuorumCert {
	if x != nil {
		return x.PreCommitQC
	}
	return nil
}

func (x *Commit) GetViewNum() uint64 {
	if x != nil {
		return x.ViewNum
	}
	return 0
}

type CommitVote struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	BlockHash     []byte                 `protobuf:"bytes,1,opt,name=blockHash,proto3" json:"blockHash,omitempty"`
	Qc            *QuorumCert            `protobuf:"bytes,2,opt,name=qc,proto3" json:"qc,omitempty"`
	PartialSig    []byte                 `protobuf:"bytes,3,opt,name=partialSig,proto3" json:"partialSig,omitempty"`
	ViewNum       uint64                 `protobuf:"varint,4,opt,name=viewNum,proto3" json:"viewNum,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CommitVote) Reset() {
	*x = CommitVote{}
	mi := &file_orderer_consensus_hotstuff_proto_hotstuff_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CommitVote) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommitVote) ProtoMessage() {}

func (x *CommitVote) ProtoReflect() protoreflect.Message {
	mi := &file_orderer_consensus_hotstuff_proto_hotstuff_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommitVote.ProtoReflect.Descriptor instead.
func (*CommitVote) Descriptor() ([]byte, []int) {
	return file_orderer_consensus_hotstuff_proto_hotstuff_proto_rawDescGZIP(), []int{7}
}

func (x *CommitVote) GetBlockHash() []byte {
	if x != nil {
		return x.BlockHash
	}
	return nil
}

func (x *CommitVote) GetQc() *QuorumCert {
	if x != nil {
		return x.Qc
	}
	return nil
}

func (x *CommitVote) GetPartialSig() []byte {
	if x != nil {
		return x.PartialSig
	}
	return nil
}

func (x *CommitVote) GetViewNum() uint64 {
	if x != nil {
		return x.ViewNum
	}
	return 0
}

type Decide struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	CommitQC      *QuorumCert            `protobuf:"bytes,1,opt,name=commitQC,proto3" json:"commitQC,omitempty"`
	ViewNum       uint64                 `protobuf:"varint,2,opt,name=viewNum,proto3" json:"viewNum,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Decide) Reset() {
	*x = Decide{}
	mi := &file_orderer_consensus_hotstuff_proto_hotstuff_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Decide) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Decide) ProtoMessage() {}

func (x *Decide) ProtoReflect() protoreflect.Message {
	mi := &file_orderer_consensus_hotstuff_proto_hotstuff_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Decide.ProtoReflect.Descriptor instead.
func (*Decide) Descriptor() ([]byte, []int) {
	return file_orderer_consensus_hotstuff_proto_hotstuff_proto_rawDescGZIP(), []int{8}
}

func (x *Decide) GetCommitQC() *QuorumCert {
	if x != nil {
		return x.CommitQC
	}
	return nil
}

func (x *Decide) GetViewNum() uint64 {
	if x != nil {
		return x.ViewNum
	}
	return 0
}

type NewView struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	PrepareQC     *QuorumCert            `protobuf:"bytes,1,opt,name=prepareQC,proto3" json:"prepareQC,omitempty"`
	ViewNum       uint64                 `protobuf:"varint,2,opt,name=viewNum,proto3" json:"viewNum,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *NewView) Reset() {
	*x = NewView{}
	mi := &file_orderer_consensus_hotstuff_proto_hotstuff_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *NewView) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NewView) ProtoMessage() {}

func (x *NewView) ProtoReflect() protoreflect.Message {
	mi := &file_orderer_consensus_hotstuff_proto_hotstuff_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NewView.ProtoReflect.Descriptor instead.
func (*NewView) Descriptor() ([]byte, []int) {
	return file_orderer_consensus_hotstuff_proto_hotstuff_proto_rawDescGZIP(), []int{9}
}

func (x *NewView) GetPrepareQC() *QuorumCert {
	if x != nil {
		return x.PrepareQC
	}
	return nil
}

func (x *NewView) GetViewNum() uint64 {
	if x != nil {
		return x.ViewNum
	}
	return 0
}

type Request struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Envelope      *common.Envelope       `protobuf:"bytes,1,opt,name=envelope,proto3" json:"envelope,omitempty"`
	ChannalId     string                 `protobuf:"bytes,2,opt,name=channal_id,json=channalId,proto3" json:"channal_id,omitempty"`
	ConfigSeq     uint64                 `protobuf:"varint,3,opt,name=config_seq,json=configSeq,proto3" json:"config_seq,omitempty"`
	IsNormal      bool                   `protobuf:"varint,4,opt,name=is_normal,json=isNormal,proto3" json:"is_normal,omitempty"`
	TimeStamp     uint64                 `protobuf:"varint,5,opt,name=time_stamp,json=timeStamp,proto3" json:"time_stamp,omitempty"`
	ClientAddress string                 `protobuf:"bytes,6,opt,name=client_address,json=clientAddress,proto3" json:"client_address,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Request) Reset() {
	*x = Request{}
	mi := &file_orderer_consensus_hotstuff_proto_hotstuff_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file_orderer_consensus_hotstuff_proto_hotstuff_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Request.ProtoReflect.Descriptor instead.
func (*Request) Descriptor() ([]byte, []int) {
	return file_orderer_consensus_hotstuff_proto_hotstuff_proto_rawDescGZIP(), []int{10}
}

func (x *Request) GetEnvelope() *common.Envelope {
	if x != nil {
		return x.Envelope
	}
	return nil
}

func (x *Request) GetChannalId() string {
	if x != nil {
		return x.ChannalId
	}
	return ""
}

func (x *Request) GetConfigSeq() uint64 {
	if x != nil {
		return x.ConfigSeq
	}
	return 0
}

func (x *Request) GetIsNormal() bool {
	if x != nil {
		return x.IsNormal
	}
	return false
}

func (x *Request) GetTimeStamp() uint64 {
	if x != nil {
		return x.TimeStamp
	}
	return 0
}

func (x *Request) GetClientAddress() string {
	if x != nil {
		return x.ClientAddress
	}
	return ""
}

type Reply struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Result        string                 `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"`
	Command       string                 `protobuf:"bytes,2,opt,name=command,proto3" json:"command,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Reply) Reset() {
	*x = Reply{}
	mi := &file_orderer_consensus_hotstuff_proto_hotstuff_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Reply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Reply) ProtoMessage() {}

func (x *Reply) ProtoReflect() protoreflect.Message {
	mi := &file_orderer_consensus_hotstuff_proto_hotstuff_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Reply.ProtoReflect.Descriptor instead.
func (*Reply) Descriptor() ([]byte, []int) {
	return file_orderer_consensus_hotstuff_proto_hotstuff_proto_rawDescGZIP(), []int{11}
}

func (x *Reply) GetResult() string {
	if x != nil {
		return x.Result
	}
	return ""
}

func (x *Reply) GetCommand() string {
	if x != nil {
		return x.Command
	}
	return ""
}

var File_orderer_consensus_hotstuff_proto_hotstuff_proto protoreflect.FileDescriptor

const file_orderer_consensus_hotstuff_proto_hotstuff_proto_rawDesc = "" +
	"\n" +
	"/orderer/consensus/hotstuff/proto/hotstuff.proto\x12\x05proto\x1a*orderer/consensus/hotstuff/proto/cmm.proto\x1a\x1aprotos/common/common.proto\"\xe9\x03\n" +
	"\x03Msg\x12*\n" +
	"\aprepare\x18\x01 \x01(\v2\x0e.proto.PrepareH\x00R\aprepare\x126\n" +
	"\vprepareVote\x18\x02 \x01(\v2\x12.proto.PrepareVoteH\x00R\vprepareVote\x120\n" +
	"\tpreCommit\x18\x03 \x01(\v2\x10.proto.PreCommitH\x00R\tpreCommit\x12<\n" +
	"\rpreCommitVote\x18\x04 \x01(\v2\x14.proto.PreCommitVoteH\x00R\rpreCommitVote\x12'\n" +
	"\x06commit\x18\x05 \x01(\v2\r.proto.CommitH\x00R\x06commit\x123\n" +
	"\n" +
	"commitVote\x18\x06 \x01(\v2\x11.proto.CommitVoteH\x00R\n" +
	"commitVote\x12'\n" +
	"\x06decide\x18\a \x01(\v2\r.proto.DecideH\x00R\x06decide\x12*\n" +
	"\anewView\x18\b \x01(\v2\x0e.proto.NewViewH\x00R\anewView\x12*\n" +
	"\arequest\x18\t \x01(\v2\x0e.proto.RequestH\x00R\arequest\x12$\n" +
	"\x05reply\x18\n" +
	" \x01(\v2\f.proto.ReplyH\x00R\x05replyB\t\n" +
	"\aPayload\"\a\n" +
	"\x05Empty\"~\n" +
	"\aPrepare\x12.\n" +
	"\vcurProposal\x18\x01 \x01(\v2\f.proto.BlockR\vcurProposal\x12)\n" +
	"\x06highQC\x18\x02 \x01(\v2\x11.proto.QuorumCertR\x06highQC\x12\x18\n" +
	"\aviewNum\x18\x03 \x01(\x04R\aviewNum\"\x88\x01\n" +
	"\vPrepareVote\x12\x1c\n" +
	"\tblockHash\x18\x01 \x01(\fR\tblockHash\x12!\n" +
	"\x02qc\x18\x02 \x01(\v2\x11.proto.QuorumCertR\x02qc\x12\x1e\n" +
	"\n" +
	"partialSig\x18\x03 \x01(\fR\n" +
	"partialSig\x12\x18\n" +
	"\aviewNum\x18\x04 \x01(\x04R\aviewNum\"V\n" +
	"\tPreCommit\x12/\n" +
	"\tprepareQC\x18\x01 \x01(\v2\x11.proto.QuorumCertR\tprepareQC\x12\x18\n" +
	"\aviewNum\x18\x02 \x01(\x04R\aviewNum\"\x8a\x01\n" +
	"\rPreCommitVote\x12\x1c\n" +
	"\tblockHash\x18\x01 \x01(\fR\tblockHash\x12!\n" +
	"\x02qc\x18\x02 \x01(\v2\x11.proto.QuorumCertR\x02qc\x12\x1e\n" +
	"\n" +
	"partialSig\x18\x03 \x01(\fR\n" +
	"partialSig\x12\x18\n" +
	"\aviewNum\x18\x04 \x01(\x04R\aviewNum\"W\n" +
	"\x06Commit\x123\n" +
	"\vpreCommitQC\x18\x01 \x01(\v2\x11.proto.QuorumCertR\vpreCommitQC\x12\x18\n" +
	"\aviewNum\x18\x02 \x01(\x04R\aviewNum\"\x87\x01\n" +
	"\n" +
	"CommitVote\x12\x1c\n" +
	"\tblockHash\x18\x01 \x01(\fR\tblockHash\x12!\n" +
	"\x02qc\x18\x02 \x01(\v2\x11.proto.QuorumCertR\x02qc\x12\x1e\n" +
	"\n" +
	"partialSig\x18\x03 \x01(\fR\n" +
	"partialSig\x12\x18\n" +
	"\aviewNum\x18\x04 \x01(\x04R\aviewNum\"Q\n" +
	"\x06Decide\x12-\n" +
	"\bcommitQC\x18\x01 \x01(\v2\x11.proto.QuorumCertR\bcommitQC\x12\x18\n" +
	"\aviewNum\x18\x02 \x01(\x04R\aviewNum\"T\n" +
	"\aNewView\x12/\n" +
	"\tprepareQC\x18\x01 \x01(\v2\x11.proto.QuorumCertR\tprepareQC\x12\x18\n" +
	"\aviewNum\x18\x02 \x01(\x04R\aviewNum\"\xd8\x01\n" +
	"\aRequest\x12,\n" +
	"\benvelope\x18\x01 \x01(\v2\x10.common.EnvelopeR\benvelope\x12\x1d\n" +
	"\n" +
	"channal_id\x18\x02 \x01(\tR\tchannalId\x12\x1d\n" +
	"\n" +
	"config_seq\x18\x03 \x01(\x04R\tconfigSeq\x12\x1b\n" +
	"\tis_normal\x18\x04 \x01(\bR\bisNormal\x12\x1d\n" +
	"\n" +
	"time_stamp\x18\x05 \x01(\x04R\ttimeStamp\x12%\n" +
	"\x0eclient_address\x18\x06 \x01(\tR\rclientAddress\"9\n" +
	"\x05Reply\x12\x16\n" +
	"\x06result\x18\x01 \x01(\tR\x06result\x12\x18\n" +
	"\acommand\x18\x02 \x01(\tR\acommand2\x8c\x01\n" +
	"\x0fHotStuffService\x12%\n" +
	"\aSendMsg\x12\n" +
	".proto.Msg\x1a\f.proto.Empty\"\x00\x12)\n" +
	"\vSendRequest\x12\n" +
	".proto.Msg\x1a\f.proto.Empty\"\x00\x12'\n" +
	"\tSendReply\x12\n" +
	".proto.Msg\x1a\f.proto.Empty\"\x00B@Z>github.com/hyperledger/fabric/orderer/consensus/hotstuff/protob\x06proto3"

var (
	file_orderer_consensus_hotstuff_proto_hotstuff_proto_rawDescOnce sync.Once
	file_orderer_consensus_hotstuff_proto_hotstuff_proto_rawDescData []byte
)

func file_orderer_consensus_hotstuff_proto_hotstuff_proto_rawDescGZIP() []byte {
	file_orderer_consensus_hotstuff_proto_hotstuff_proto_rawDescOnce.Do(func() {
		file_orderer_consensus_hotstuff_proto_hotstuff_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_orderer_consensus_hotstuff_proto_hotstuff_proto_rawDesc), len(file_orderer_consensus_hotstuff_proto_hotstuff_proto_rawDesc)))
	})
	return file_orderer_consensus_hotstuff_proto_hotstuff_proto_rawDescData
}

var file_orderer_consensus_hotstuff_proto_hotstuff_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_orderer_consensus_hotstuff_proto_hotstuff_proto_goTypes = []any{
	(*Msg)(nil),             // 0: proto.Msg
	(*Empty)(nil),           // 1: proto.Empty
	(*Prepare)(nil),         // 2: proto.Prepare
	(*PrepareVote)(nil),     // 3: proto.PrepareVote
	(*PreCommit)(nil),       // 4: proto.PreCommit
	(*PreCommitVote)(nil),   // 5: proto.PreCommitVote
	(*Commit)(nil),          // 6: proto.Commit
	(*CommitVote)(nil),      // 7: proto.CommitVote
	(*Decide)(nil),          // 8: proto.Decide
	(*NewView)(nil),         // 9: proto.NewView
	(*Request)(nil),         // 10: proto.Request
	(*Reply)(nil),           // 11: proto.Reply
	(*Block)(nil),           // 12: proto.Block
	(*QuorumCert)(nil),      // 13: proto.QuorumCert
	(*common.Envelope)(nil), // 14: common.Envelope
}
var file_orderer_consensus_hotstuff_proto_hotstuff_proto_depIdxs = []int32{
	2,  // 0: proto.Msg.prepare:type_name -> proto.Prepare
	3,  // 1: proto.Msg.prepareVote:type_name -> proto.PrepareVote
	4,  // 2: proto.Msg.preCommit:type_name -> proto.PreCommit
	5,  // 3: proto.Msg.preCommitVote:type_name -> proto.PreCommitVote
	6,  // 4: proto.Msg.commit:type_name -> proto.Commit
	7,  // 5: proto.Msg.commitVote:type_name -> proto.CommitVote
	8,  // 6: proto.Msg.decide:type_name -> proto.Decide
	9,  // 7: proto.Msg.newView:type_name -> proto.NewView
	10, // 8: proto.Msg.request:type_name -> proto.Request
	11, // 9: proto.Msg.reply:type_name -> proto.Reply
	12, // 10: proto.Prepare.curProposal:type_name -> proto.Block
	13, // 11: proto.Prepare.highQC:type_name -> proto.QuorumCert
	13, // 12: proto.PrepareVote.qc:type_name -> proto.QuorumCert
	13, // 13: proto.PreCommit.prepareQC:type_name -> proto.QuorumCert
	13, // 14: proto.PreCommitVote.qc:type_name -> proto.QuorumCert
	13, // 15: proto.Commit.preCommitQC:type_name -> proto.QuorumCert
	13, // 16: proto.CommitVote.qc:type_name -> proto.QuorumCert
	13, // 17: proto.Decide.commitQC:type_name -> proto.QuorumCert
	13, // 18: proto.NewView.prepareQC:type_name -> proto.QuorumCert
	14, // 19: proto.Request.envelope:type_name -> common.Envelope
	0,  // 20: proto.HotStuffService.SendMsg:input_type -> proto.Msg
	0,  // 21: proto.HotStuffService.SendRequest:input_type -> proto.Msg
	0,  // 22: proto.HotStuffService.SendReply:input_type -> proto.Msg
	1,  // 23: proto.HotStuffService.SendMsg:output_type -> proto.Empty
	1,  // 24: proto.HotStuffService.SendRequest:output_type -> proto.Empty
	1,  // 25: proto.HotStuffService.SendReply:output_type -> proto.Empty
	23, // [23:26] is the sub-list for method output_type
	20, // [20:23] is the sub-list for method input_type
	20, // [20:20] is the sub-list for extension type_name
	20, // [20:20] is the sub-list for extension extendee
	0,  // [0:20] is the sub-list for field type_name
}

func init() { file_orderer_consensus_hotstuff_proto_hotstuff_proto_init() }
func file_orderer_consensus_hotstuff_proto_hotstuff_proto_init() {
	if File_orderer_consensus_hotstuff_proto_hotstuff_proto != nil {
		return
	}
	file_orderer_consensus_hotstuff_proto_cmm_proto_init()
	file_orderer_consensus_hotstuff_proto_hotstuff_proto_msgTypes[0].OneofWrappers = []any{
		(*Msg_Prepare)(nil),
		(*Msg_PrepareVote)(nil),
		(*Msg_PreCommit)(nil),
		(*Msg_PreCommitVote)(nil),
		(*Msg_Commit)(nil),
		(*Msg_CommitVote)(nil),
		(*Msg_Decide)(nil),
		(*Msg_NewView)(nil),
		(*Msg_Request)(nil),
		(*Msg_Reply)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_orderer_consensus_hotstuff_proto_hotstuff_proto_rawDesc), len(file_orderer_consensus_hotstuff_proto_hotstuff_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_orderer_consensus_hotstuff_proto_hotstuff_proto_goTypes,
		DependencyIndexes: file_orderer_consensus_hotstuff_proto_hotstuff_proto_depIdxs,
		MessageInfos:      file_orderer_consensus_hotstuff_proto_hotstuff_proto_msgTypes,
	}.Build()
	File_orderer_consensus_hotstuff_proto_hotstuff_proto = out.File
	file_orderer_consensus_hotstuff_proto_hotstuff_proto_goTypes = nil
	file_orderer_consensus_hotstuff_proto_hotstuff_proto_depIdxs = nil
}
