// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v3.12.4
// source: orderer/consensus/hotstuff/proto/common.proto

package deprecated

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

type MsgType int32

const (
	MsgType_PREPARE        MsgType = 0
	MsgType_PREPARE_VOTE   MsgType = 1
	MsgType_PRECOMMIT      MsgType = 2
	MsgType_PRECOMMIT_VOTE MsgType = 3
	MsgType_COMMIT         MsgType = 4
	MsgType_COMMIT_VOTE    MsgType = 5
	MsgType_NEWVIEW        MsgType = 6
	MsgType_DECIDE         MsgType = 7
)

// Enum value maps for MsgType.
var (
	MsgType_name = map[int32]string{
		0: "PREPARE",
		1: "PREPARE_VOTE",
		2: "PRECOMMIT",
		3: "PRECOMMIT_VOTE",
		4: "COMMIT",
		5: "COMMIT_VOTE",
		6: "NEWVIEW",
		7: "DECIDE",
	}
	MsgType_value = map[string]int32{
		"PREPARE":        0,
		"PREPARE_VOTE":   1,
		"PRECOMMIT":      2,
		"PRECOMMIT_VOTE": 3,
		"COMMIT":         4,
		"COMMIT_VOTE":    5,
		"NEWVIEW":        6,
		"DECIDE":         7,
	}
)

func (x MsgType) Enum() *MsgType {
	p := new(MsgType)
	*p = x
	return p
}

func (x MsgType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MsgType) Descriptor() protoreflect.EnumDescriptor {
	return file_orderer_consensus_hotstuff_proto_common_proto_enumTypes[0].Descriptor()
}

func (MsgType) Type() protoreflect.EnumType {
	return &file_orderer_consensus_hotstuff_proto_common_proto_enumTypes[0]
}

func (x MsgType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MsgType.Descriptor instead.
func (MsgType) EnumDescriptor() ([]byte, []int) {
	return file_orderer_consensus_hotstuff_proto_common_proto_rawDescGZIP(), []int{0}
}

type Block struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ParentHash    []byte                 `protobuf:"bytes,1,opt,name=ParentHash,proto3" json:"ParentHash,omitempty"`
	Hash          []byte                 `protobuf:"bytes,2,opt,name=Hash,proto3" json:"Hash,omitempty"`
	Height        uint64                 `protobuf:"varint,3,opt,name=height,proto3" json:"height,omitempty"`
	Commands      []*common.Envelope     `protobuf:"bytes,4,rep,name=commands,proto3" json:"commands,omitempty"`
	Justify       *QuorumCert            `protobuf:"bytes,5,opt,name=Justify,proto3" json:"Justify,omitempty"`
	Committed     bool                   `protobuf:"varint,6,opt,name=committed,proto3" json:"committed,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Block) Reset() {
	*x = Block{}
	mi := &file_orderer_consensus_hotstuff_proto_common_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Block) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Block) ProtoMessage() {}

func (x *Block) ProtoReflect() protoreflect.Message {
	mi := &file_orderer_consensus_hotstuff_proto_common_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Block.ProtoReflect.Descriptor instead.
func (*Block) Descriptor() ([]byte, []int) {
	return file_orderer_consensus_hotstuff_proto_common_proto_rawDescGZIP(), []int{0}
}

func (x *Block) GetParentHash() []byte {
	if x != nil {
		return x.ParentHash
	}
	return nil
}

func (x *Block) GetHash() []byte {
	if x != nil {
		return x.Hash
	}
	return nil
}

func (x *Block) GetHeight() uint64 {
	if x != nil {
		return x.Height
	}
	return 0
}

func (x *Block) GetCommands() []*common.Envelope {
	if x != nil {
		return x.Commands
	}
	return nil
}

func (x *Block) GetJustify() *QuorumCert {
	if x != nil {
		return x.Justify
	}
	return nil
}

func (x *Block) GetCommitted() bool {
	if x != nil {
		return x.Committed
	}
	return false
}

type QuorumCert struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	BlockHash     []byte                 `protobuf:"bytes,1,opt,name=BlockHash,proto3" json:"BlockHash,omitempty"`
	Type          MsgType                `protobuf:"varint,2,opt,name=type,proto3,enum=proto.MsgType" json:"type,omitempty"`
	ViewNum       uint64                 `protobuf:"varint,3,opt,name=viewNum,proto3" json:"viewNum,omitempty"`
	Signature     []byte                 `protobuf:"bytes,4,opt,name=signature,proto3" json:"signature,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *QuorumCert) Reset() {
	*x = QuorumCert{}
	mi := &file_orderer_consensus_hotstuff_proto_common_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QuorumCert) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QuorumCert) ProtoMessage() {}

func (x *QuorumCert) ProtoReflect() protoreflect.Message {
	mi := &file_orderer_consensus_hotstuff_proto_common_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QuorumCert.ProtoReflect.Descriptor instead.
func (*QuorumCert) Descriptor() ([]byte, []int) {
	return file_orderer_consensus_hotstuff_proto_common_proto_rawDescGZIP(), []int{1}
}

func (x *QuorumCert) GetBlockHash() []byte {
	if x != nil {
		return x.BlockHash
	}
	return nil
}

func (x *QuorumCert) GetType() MsgType {
	if x != nil {
		return x.Type
	}
	return MsgType_PREPARE
}

func (x *QuorumCert) GetViewNum() uint64 {
	if x != nil {
		return x.ViewNum
	}
	return 0
}

func (x *QuorumCert) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

var File_orderer_consensus_hotstuff_proto_common_proto protoreflect.FileDescriptor

const file_orderer_consensus_hotstuff_proto_common_proto_rawDesc = "" +
	"\n" +
	"-orderer/consensus/hotstuff/proto/common.proto\x12\x05proto\x1a\x1aprotos/common/common.proto\"\xcc\x01\n" +
	"\x05Block\x12\x1e\n" +
	"\n" +
	"ParentHash\x18\x01 \x01(\fR\n" +
	"ParentHash\x12\x12\n" +
	"\x04Hash\x18\x02 \x01(\fR\x04Hash\x12\x16\n" +
	"\x06height\x18\x03 \x01(\x04R\x06height\x12,\n" +
	"\bcommands\x18\x04 \x03(\v2\x10.common.EnvelopeR\bcommands\x12+\n" +
	"\aJustify\x18\x05 \x01(\v2\x11.proto.QuorumCertR\aJustify\x12\x1c\n" +
	"\tcommitted\x18\x06 \x01(\bR\tcommitted\"\x86\x01\n" +
	"\n" +
	"QuorumCert\x12\x1c\n" +
	"\tBlockHash\x18\x01 \x01(\fR\tBlockHash\x12\"\n" +
	"\x04type\x18\x02 \x01(\x0e2\x0e.proto.MsgTypeR\x04type\x12\x18\n" +
	"\aviewNum\x18\x03 \x01(\x04R\aviewNum\x12\x1c\n" +
	"\tsignature\x18\x04 \x01(\fR\tsignature*\x81\x01\n" +
	"\aMsgType\x12\v\n" +
	"\aPREPARE\x10\x00\x12\x10\n" +
	"\fPREPARE_VOTE\x10\x01\x12\r\n" +
	"\tPRECOMMIT\x10\x02\x12\x12\n" +
	"\x0ePRECOMMIT_VOTE\x10\x03\x12\n" +
	"\n" +
	"\x06COMMIT\x10\x04\x12\x0f\n" +
	"\vCOMMIT_VOTE\x10\x05\x12\v\n" +
	"\aNEWVIEW\x10\x06\x12\n" +
	"\n" +
	"\x06DECIDE\x10\aB)Z'github.com/theodocius271/hotstuff/protob\x06proto3"

var (
	file_orderer_consensus_hotstuff_proto_common_proto_rawDescOnce sync.Once
	file_orderer_consensus_hotstuff_proto_common_proto_rawDescData []byte
)

func file_orderer_consensus_hotstuff_proto_common_proto_rawDescGZIP() []byte {
	file_orderer_consensus_hotstuff_proto_common_proto_rawDescOnce.Do(func() {
		file_orderer_consensus_hotstuff_proto_common_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_orderer_consensus_hotstuff_proto_common_proto_rawDesc), len(file_orderer_consensus_hotstuff_proto_common_proto_rawDesc)))
	})
	return file_orderer_consensus_hotstuff_proto_common_proto_rawDescData
}

var file_orderer_consensus_hotstuff_proto_common_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_orderer_consensus_hotstuff_proto_common_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_orderer_consensus_hotstuff_proto_common_proto_goTypes = []any{
	(MsgType)(0),            // 0: proto.MsgType
	(*Block)(nil),           // 1: proto.Block
	(*QuorumCert)(nil),      // 2: proto.QuorumCert
	(*common.Envelope)(nil), // 3: common.Envelope
}
var file_orderer_consensus_hotstuff_proto_common_proto_depIdxs = []int32{
	3, // 0: proto.Block.commands:type_name -> common.Envelope
	2, // 1: proto.Block.Justify:type_name -> proto.QuorumCert
	0, // 2: proto.QuorumCert.type:type_name -> proto.MsgType
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_orderer_consensus_hotstuff_proto_common_proto_init() }
func file_orderer_consensus_hotstuff_proto_common_proto_init() {
	if File_orderer_consensus_hotstuff_proto_common_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_orderer_consensus_hotstuff_proto_common_proto_rawDesc), len(file_orderer_consensus_hotstuff_proto_common_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_orderer_consensus_hotstuff_proto_common_proto_goTypes,
		DependencyIndexes: file_orderer_consensus_hotstuff_proto_common_proto_depIdxs,
		EnumInfos:         file_orderer_consensus_hotstuff_proto_common_proto_enumTypes,
		MessageInfos:      file_orderer_consensus_hotstuff_proto_common_proto_msgTypes,
	}.Build()
	File_orderer_consensus_hotstuff_proto_common_proto = out.File
	file_orderer_consensus_hotstuff_proto_common_proto_goTypes = nil
	file_orderer_consensus_hotstuff_proto_common_proto_depIdxs = nil
}
