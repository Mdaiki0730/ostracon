// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: proto/state/types.proto

package state

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	_ "github.com/golang/protobuf/ptypes/timestamp"
	types "github.com/tendermint/tendermint/abci/types"
	types1 "github.com/tendermint/tendermint/proto/types"
	version "github.com/tendermint/tendermint/proto/version"
	math "math"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// ABCIResponses retains the responses
// of the various ABCI calls during block processing.
// It is persisted to disk for each height before calling Commit.
type ABCIResponses struct {
	DeliverTxs           []*types.ResponseDeliverTx `protobuf:"bytes,1,rep,name=deliver_txs,json=deliverTxs,proto3" json:"deliver_txs,omitempty"`
	EndBlock             *types.ResponseEndBlock    `protobuf:"bytes,2,opt,name=end_block,json=endBlock,proto3" json:"end_block,omitempty"`
	BeginBlock           *types.ResponseBeginBlock  `protobuf:"bytes,3,opt,name=begin_block,json=beginBlock,proto3" json:"begin_block,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                   `json:"-"`
	XXX_unrecognized     []byte                     `json:"-"`
	XXX_sizecache        int32                      `json:"-"`
}

func (m *ABCIResponses) Reset()         { *m = ABCIResponses{} }
func (m *ABCIResponses) String() string { return proto.CompactTextString(m) }
func (*ABCIResponses) ProtoMessage()    {}
func (*ABCIResponses) Descriptor() ([]byte, []int) {
	return fileDescriptor_00e69fef8162ea9b, []int{0}
}
func (m *ABCIResponses) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ABCIResponses.Unmarshal(m, b)
}
func (m *ABCIResponses) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ABCIResponses.Marshal(b, m, deterministic)
}
func (m *ABCIResponses) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ABCIResponses.Merge(m, src)
}
func (m *ABCIResponses) XXX_Size() int {
	return xxx_messageInfo_ABCIResponses.Size(m)
}
func (m *ABCIResponses) XXX_DiscardUnknown() {
	xxx_messageInfo_ABCIResponses.DiscardUnknown(m)
}

var xxx_messageInfo_ABCIResponses proto.InternalMessageInfo

func (m *ABCIResponses) GetDeliverTxs() []*types.ResponseDeliverTx {
	if m != nil {
		return m.DeliverTxs
	}
	return nil
}

func (m *ABCIResponses) GetEndBlock() *types.ResponseEndBlock {
	if m != nil {
		return m.EndBlock
	}
	return nil
}

func (m *ABCIResponses) GetBeginBlock() *types.ResponseBeginBlock {
	if m != nil {
		return m.BeginBlock
	}
	return nil
}

// ValidatorsInfo represents the latest validator set, or the last height it changed
type ValidatorsInfo struct {
	ValidatorSet         *types1.ValidatorSet `protobuf:"bytes,1,opt,name=validator_set,json=validatorSet,proto3" json:"validator_set,omitempty"`
	LastHeightChanged    int64                `protobuf:"varint,2,opt,name=last_height_changed,json=lastHeightChanged,proto3" json:"last_height_changed,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *ValidatorsInfo) Reset()         { *m = ValidatorsInfo{} }
func (m *ValidatorsInfo) String() string { return proto.CompactTextString(m) }
func (*ValidatorsInfo) ProtoMessage()    {}
func (*ValidatorsInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_00e69fef8162ea9b, []int{1}
}
func (m *ValidatorsInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ValidatorsInfo.Unmarshal(m, b)
}
func (m *ValidatorsInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ValidatorsInfo.Marshal(b, m, deterministic)
}
func (m *ValidatorsInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ValidatorsInfo.Merge(m, src)
}
func (m *ValidatorsInfo) XXX_Size() int {
	return xxx_messageInfo_ValidatorsInfo.Size(m)
}
func (m *ValidatorsInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_ValidatorsInfo.DiscardUnknown(m)
}

var xxx_messageInfo_ValidatorsInfo proto.InternalMessageInfo

func (m *ValidatorsInfo) GetValidatorSet() *types1.ValidatorSet {
	if m != nil {
		return m.ValidatorSet
	}
	return nil
}

func (m *ValidatorsInfo) GetLastHeightChanged() int64 {
	if m != nil {
		return m.LastHeightChanged
	}
	return 0
}

// ConsensusParamsInfo represents the latest consensus params, or the last height it changed
type ConsensusParamsInfo struct {
	ConsensusParams      types1.ConsensusParams `protobuf:"bytes,1,opt,name=consensus_params,json=consensusParams,proto3" json:"consensus_params"`
	LastHeightChanged    int64                  `protobuf:"varint,2,opt,name=last_height_changed,json=lastHeightChanged,proto3" json:"last_height_changed,omitempty"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *ConsensusParamsInfo) Reset()         { *m = ConsensusParamsInfo{} }
func (m *ConsensusParamsInfo) String() string { return proto.CompactTextString(m) }
func (*ConsensusParamsInfo) ProtoMessage()    {}
func (*ConsensusParamsInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_00e69fef8162ea9b, []int{2}
}
func (m *ConsensusParamsInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConsensusParamsInfo.Unmarshal(m, b)
}
func (m *ConsensusParamsInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConsensusParamsInfo.Marshal(b, m, deterministic)
}
func (m *ConsensusParamsInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConsensusParamsInfo.Merge(m, src)
}
func (m *ConsensusParamsInfo) XXX_Size() int {
	return xxx_messageInfo_ConsensusParamsInfo.Size(m)
}
func (m *ConsensusParamsInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_ConsensusParamsInfo.DiscardUnknown(m)
}

var xxx_messageInfo_ConsensusParamsInfo proto.InternalMessageInfo

func (m *ConsensusParamsInfo) GetConsensusParams() types1.ConsensusParams {
	if m != nil {
		return m.ConsensusParams
	}
	return types1.ConsensusParams{}
}

func (m *ConsensusParamsInfo) GetLastHeightChanged() int64 {
	if m != nil {
		return m.LastHeightChanged
	}
	return 0
}

type Version struct {
	Consensus            version.Consensus `protobuf:"bytes,1,opt,name=consensus,proto3" json:"consensus"`
	Software             string            `protobuf:"bytes,2,opt,name=software,proto3" json:"software,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Version) Reset()         { *m = Version{} }
func (m *Version) String() string { return proto.CompactTextString(m) }
func (*Version) ProtoMessage()    {}
func (*Version) Descriptor() ([]byte, []int) {
	return fileDescriptor_00e69fef8162ea9b, []int{3}
}
func (m *Version) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Version.Unmarshal(m, b)
}
func (m *Version) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Version.Marshal(b, m, deterministic)
}
func (m *Version) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Version.Merge(m, src)
}
func (m *Version) XXX_Size() int {
	return xxx_messageInfo_Version.Size(m)
}
func (m *Version) XXX_DiscardUnknown() {
	xxx_messageInfo_Version.DiscardUnknown(m)
}

var xxx_messageInfo_Version proto.InternalMessageInfo

func (m *Version) GetConsensus() version.Consensus {
	if m != nil {
		return m.Consensus
	}
	return version.Consensus{}
}

func (m *Version) GetSoftware() string {
	if m != nil {
		return m.Software
	}
	return ""
}

type State struct {
	Version Version `protobuf:"bytes,1,opt,name=version,proto3" json:"version"`
	// immutable
	ChainID string `protobuf:"bytes,2,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	// voter_params
	VoterParams *types1.VoterParams `protobuf:"bytes,3,opt,name=voter_params,json=voterParams,proto3" json:"voter_params,omitempty"`
	// LastBlockHeight=0 at genesis (ie. block(H=0) does not exist)
	LastBlockHeight int64          `protobuf:"varint,4,opt,name=last_block_height,json=lastBlockHeight,proto3" json:"last_block_height,omitempty"`
	LastBlockID     types1.BlockID `protobuf:"bytes,5,opt,name=last_block_id,json=lastBlockId,proto3" json:"last_block_id"`
	LastBlockTime   time.Time      `protobuf:"bytes,6,opt,name=last_block_time,json=lastBlockTime,proto3,stdtime" json:"last_block_time"`
	// vrf hash from proof
	LastProofHash []byte `protobuf:"bytes,7,opt,name=last_proof_hash,json=lastProofHash,proto3" json:"last_proof_hash,omitempty"`
	// LastValidators is used to validate block.LastCommit.
	// Validators are persisted to the database separately every time they change,
	// so we can query for historical validator sets.
	// Note that if s.LastBlockHeight causes a valset change,
	// we set s.LastHeightValidatorsChanged = s.LastBlockHeight + 1 + 1
	// Extra +1 due to nextValSet delay.
	NextValidators              *types1.ValidatorSet `protobuf:"bytes,8,opt,name=next_validators,json=nextValidators,proto3" json:"next_validators,omitempty"`
	Validators                  *types1.ValidatorSet `protobuf:"bytes,9,opt,name=validators,proto3" json:"validators,omitempty"`
	NextVoters                  *types1.VoterSet     `protobuf:"bytes,10,opt,name=next_voters,json=nextVoters,proto3" json:"next_voters,omitempty"`
	Voters                      *types1.VoterSet     `protobuf:"bytes,11,opt,name=voters,proto3" json:"voters,omitempty"`
	LastVoters                  *types1.VoterSet     `protobuf:"bytes,12,opt,name=last_voters,json=lastVoters,proto3" json:"last_voters,omitempty"`
	LastHeightValidatorsChanged int64                `protobuf:"varint,13,opt,name=last_height_validators_changed,json=lastHeightValidatorsChanged,proto3" json:"last_height_validators_changed,omitempty"`
	// Consensus parameters used for validating blocks.
	// Changes returned by EndBlock and updated after Commit.
	ConsensusParams                  types1.ConsensusParams `protobuf:"bytes,14,opt,name=consensus_params,json=consensusParams,proto3" json:"consensus_params"`
	LastHeightConsensusParamsChanged int64                  `protobuf:"varint,15,opt,name=last_height_consensus_params_changed,json=lastHeightConsensusParamsChanged,proto3" json:"last_height_consensus_params_changed,omitempty"`
	// Merkle root of the results from executing prev block
	LastResultsHash []byte `protobuf:"bytes,16,opt,name=last_results_hash,json=lastResultsHash,proto3" json:"last_results_hash,omitempty"`
	// the latest AppHash we've received from calling abci.Commit()
	AppHash              []byte   `protobuf:"bytes,17,opt,name=app_hash,json=appHash,proto3" json:"app_hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *State) Reset()         { *m = State{} }
func (m *State) String() string { return proto.CompactTextString(m) }
func (*State) ProtoMessage()    {}
func (*State) Descriptor() ([]byte, []int) {
	return fileDescriptor_00e69fef8162ea9b, []int{4}
}
func (m *State) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_State.Unmarshal(m, b)
}
func (m *State) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_State.Marshal(b, m, deterministic)
}
func (m *State) XXX_Merge(src proto.Message) {
	xxx_messageInfo_State.Merge(m, src)
}
func (m *State) XXX_Size() int {
	return xxx_messageInfo_State.Size(m)
}
func (m *State) XXX_DiscardUnknown() {
	xxx_messageInfo_State.DiscardUnknown(m)
}

var xxx_messageInfo_State proto.InternalMessageInfo

func (m *State) GetVersion() Version {
	if m != nil {
		return m.Version
	}
	return Version{}
}

func (m *State) GetChainID() string {
	if m != nil {
		return m.ChainID
	}
	return ""
}

func (m *State) GetVoterParams() *types1.VoterParams {
	if m != nil {
		return m.VoterParams
	}
	return nil
}

func (m *State) GetLastBlockHeight() int64 {
	if m != nil {
		return m.LastBlockHeight
	}
	return 0
}

func (m *State) GetLastBlockID() types1.BlockID {
	if m != nil {
		return m.LastBlockID
	}
	return types1.BlockID{}
}

func (m *State) GetLastBlockTime() time.Time {
	if m != nil {
		return m.LastBlockTime
	}
	return time.Time{}
}

func (m *State) GetLastProofHash() []byte {
	if m != nil {
		return m.LastProofHash
	}
	return nil
}

func (m *State) GetNextValidators() *types1.ValidatorSet {
	if m != nil {
		return m.NextValidators
	}
	return nil
}

func (m *State) GetValidators() *types1.ValidatorSet {
	if m != nil {
		return m.Validators
	}
	return nil
}

func (m *State) GetNextVoters() *types1.VoterSet {
	if m != nil {
		return m.NextVoters
	}
	return nil
}

func (m *State) GetVoters() *types1.VoterSet {
	if m != nil {
		return m.Voters
	}
	return nil
}

func (m *State) GetLastVoters() *types1.VoterSet {
	if m != nil {
		return m.LastVoters
	}
	return nil
}

func (m *State) GetLastHeightValidatorsChanged() int64 {
	if m != nil {
		return m.LastHeightValidatorsChanged
	}
	return 0
}

func (m *State) GetConsensusParams() types1.ConsensusParams {
	if m != nil {
		return m.ConsensusParams
	}
	return types1.ConsensusParams{}
}

func (m *State) GetLastHeightConsensusParamsChanged() int64 {
	if m != nil {
		return m.LastHeightConsensusParamsChanged
	}
	return 0
}

func (m *State) GetLastResultsHash() []byte {
	if m != nil {
		return m.LastResultsHash
	}
	return nil
}

func (m *State) GetAppHash() []byte {
	if m != nil {
		return m.AppHash
	}
	return nil
}

func init() {
	proto.RegisterType((*ABCIResponses)(nil), "tendermint.proto.state.ABCIResponses")
	proto.RegisterType((*ValidatorsInfo)(nil), "tendermint.proto.state.ValidatorsInfo")
	proto.RegisterType((*ConsensusParamsInfo)(nil), "tendermint.proto.state.ConsensusParamsInfo")
	proto.RegisterType((*Version)(nil), "tendermint.proto.state.Version")
	proto.RegisterType((*State)(nil), "tendermint.proto.state.State")
}

func init() { proto.RegisterFile("proto/state/types.proto", fileDescriptor_00e69fef8162ea9b) }

var fileDescriptor_00e69fef8162ea9b = []byte{
	// 811 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x56, 0xdd, 0x4e, 0xdb, 0x4a,
	0x10, 0x3e, 0x3e, 0x01, 0x92, 0xac, 0x13, 0x02, 0x8b, 0xc4, 0xf1, 0x09, 0xd2, 0x49, 0x94, 0x83,
	0x20, 0xad, 0x2a, 0x47, 0xa2, 0x37, 0xbd, 0xab, 0x70, 0xd2, 0x16, 0x57, 0xb4, 0x42, 0x06, 0x21,
	0xd4, 0x1b, 0x6b, 0x13, 0x6f, 0x6c, 0xab, 0x89, 0xd7, 0xf2, 0x6e, 0x52, 0x78, 0x86, 0xde, 0xf4,
	0x0d, 0xfa, 0x28, 0xbd, 0xed, 0x53, 0x50, 0x89, 0xeb, 0x3e, 0x44, 0xb5, 0x3f, 0xfe, 0x09, 0x3f,
	0x05, 0xa4, 0x5e, 0x65, 0xbd, 0x33, 0xdf, 0x37, 0xdf, 0xcc, 0xce, 0x8c, 0x02, 0xfe, 0x89, 0x13,
	0xc2, 0x48, 0x8f, 0x32, 0xc4, 0x70, 0x8f, 0x5d, 0xc4, 0x98, 0x9a, 0xe2, 0x06, 0x6e, 0x32, 0x1c,
	0x79, 0x38, 0x99, 0x86, 0x11, 0x93, 0x37, 0xa6, 0xf0, 0x69, 0xee, 0xb0, 0x20, 0x4c, 0x3c, 0x37,
	0x46, 0x09, 0xbb, 0xe8, 0x49, 0xb0, 0x4f, 0x7c, 0x92, 0x9f, 0xa4, 0x77, 0x73, 0x13, 0x0d, 0x47,
	0xa1, 0x64, 0x2c, 0xf2, 0x36, 0x55, 0xc0, 0x9b, 0x86, 0xad, 0xa2, 0x61, 0x8e, 0x26, 0xa1, 0x87,
	0x18, 0x49, 0x94, 0xd1, 0x28, 0x1a, 0x63, 0x94, 0xa0, 0xe9, 0xad, 0x7c, 0x73, 0xc2, 0x70, 0xb2,
	0xc8, 0x37, 0xc7, 0x09, 0x0d, 0x49, 0x94, 0xfe, 0x2a, 0x63, 0xcb, 0x27, 0xc4, 0x9f, 0x60, 0x99,
	0xc0, 0x70, 0x36, 0xee, 0xb1, 0x70, 0x8a, 0x29, 0x43, 0xd3, 0x58, 0x3a, 0x74, 0x7e, 0x6a, 0xa0,
	0xbe, 0x6f, 0xf5, 0x6d, 0x07, 0xd3, 0x98, 0x44, 0x14, 0x53, 0x68, 0x03, 0xdd, 0xc3, 0x93, 0x70,
	0x8e, 0x13, 0x97, 0x9d, 0x53, 0x43, 0x6b, 0x97, 0xba, 0xfa, 0x5e, 0xd7, 0x2c, 0x94, 0x89, 0x67,
	0x6c, 0xca, 0x94, 0x52, 0xd8, 0x40, 0x22, 0x4e, 0xce, 0x1d, 0xe0, 0xa5, 0x47, 0x0a, 0x07, 0xa0,
	0x8a, 0x23, 0xcf, 0x1d, 0x4e, 0xc8, 0xe8, 0xa3, 0xf1, 0x77, 0x5b, 0xeb, 0xea, 0x7b, 0xbb, 0xf7,
	0x10, 0xbd, 0x8a, 0x3c, 0x8b, 0xbb, 0x3b, 0x15, 0xac, 0x4e, 0xf0, 0x2d, 0xd0, 0x87, 0xd8, 0x0f,
	0x23, 0xc5, 0x53, 0x12, 0x3c, 0x4f, 0xee, 0xe1, 0xb1, 0x38, 0x42, 0x32, 0x81, 0x61, 0x76, 0xee,
	0x7c, 0xd6, 0xc0, 0xea, 0x69, 0x5a, 0x73, 0x6a, 0x47, 0x63, 0x02, 0x6d, 0x50, 0xcf, 0x5e, 0xc1,
	0xa5, 0x98, 0x19, 0x9a, 0x08, 0xb0, 0x6d, 0xde, 0x68, 0x0c, 0x19, 0x21, 0x83, 0x1f, 0x63, 0xe6,
	0xd4, 0xe6, 0x85, 0x2f, 0x68, 0x82, 0x8d, 0x09, 0xa2, 0xcc, 0x0d, 0x70, 0xe8, 0x07, 0xcc, 0x1d,
	0x05, 0x28, 0xf2, 0xb1, 0x27, 0x32, 0x2f, 0x39, 0xeb, 0xdc, 0x74, 0x20, 0x2c, 0x7d, 0x69, 0xe8,
	0x7c, 0xd5, 0xc0, 0x46, 0x9f, 0xab, 0x8d, 0xe8, 0x8c, 0x1e, 0x89, 0xd7, 0x16, 0x92, 0xce, 0xc0,
	0xda, 0x28, 0xbd, 0x76, 0x65, 0x17, 0x28, 0x55, 0xbb, 0x77, 0xa9, 0xba, 0x46, 0x63, 0x2d, 0x7d,
	0xbf, 0x6c, 0xfd, 0xe5, 0x34, 0x46, 0x8b, 0xd7, 0x8f, 0x56, 0x18, 0x81, 0xf2, 0xa9, 0x6c, 0x28,
	0xf8, 0x06, 0x54, 0x33, 0x36, 0xa5, 0xe6, 0xff, 0x9b, 0x6a, 0xd2, 0xf6, 0xcb, 0xf4, 0x28, 0x25,
	0x39, 0x16, 0x36, 0x41, 0x85, 0x92, 0x31, 0xfb, 0x84, 0x12, 0x2c, 0x02, 0x57, 0x9d, 0xec, 0xbb,
	0xf3, 0xad, 0x02, 0x96, 0x8f, 0xf9, 0xfc, 0xc1, 0x97, 0xa0, 0xac, 0xb8, 0x54, 0xb0, 0x96, 0x79,
	0xfb, 0xa4, 0x9a, 0x4a, 0xa0, 0x0a, 0x94, 0xa2, 0xe0, 0x0e, 0xa8, 0x8c, 0x02, 0x14, 0x46, 0x6e,
	0x28, 0xf3, 0xab, 0x5a, 0xfa, 0xd5, 0x65, 0xab, 0xdc, 0xe7, 0x77, 0xf6, 0xc0, 0x29, 0x0b, 0xa3,
	0xed, 0xc1, 0xd7, 0xa0, 0x26, 0xc6, 0x29, 0x2d, 0x74, 0xe9, 0xae, 0xd4, 0xd4, 0xf3, 0x73, 0x5f,
	0x59, 0x4d, 0x47, 0x9f, 0xe7, 0x1f, 0xf0, 0x29, 0x10, 0xf5, 0x93, 0x5d, 0xaa, 0x0a, 0x6c, 0x2c,
	0x89, 0xc2, 0x36, 0xb8, 0x41, 0x34, 0xa0, 0xac, 0x2e, 0x3c, 0x03, 0xf5, 0x82, 0x6f, 0xe8, 0x19,
	0xcb, 0x77, 0xa5, 0x28, 0x83, 0x0a, 0xac, 0x3d, 0xb0, 0x36, 0x78, 0x8a, 0x57, 0x97, 0x2d, 0xfd,
	0x30, 0x25, 0xb4, 0x07, 0x8e, 0x9e, 0xb1, 0xdb, 0x1e, 0x3c, 0x04, 0x8d, 0x02, 0x33, 0x9f, 0x76,
	0x63, 0x45, 0x70, 0x37, 0x4d, 0xb9, 0x0a, 0xcc, 0x74, 0x15, 0x98, 0x27, 0xe9, 0x2a, 0xb0, 0x2a,
	0x9c, 0xf6, 0xcb, 0x8f, 0x96, 0xe6, 0xd4, 0x33, 0x2e, 0x6e, 0x85, 0x3b, 0x8a, 0x2d, 0x4e, 0x08,
	0x19, 0xbb, 0x01, 0xa2, 0x81, 0x51, 0x6e, 0x6b, 0xdd, 0x9a, 0xf4, 0x3b, 0xe2, 0xb7, 0x07, 0x88,
	0x06, 0xf0, 0x1d, 0x68, 0x44, 0xf8, 0x9c, 0xb9, 0xd9, 0x34, 0x50, 0xa3, 0xf2, 0x88, 0x29, 0x5a,
	0xe5, 0xe0, 0x7c, 0x2c, 0xe1, 0x00, 0x80, 0x02, 0x53, 0xf5, 0x11, 0x4c, 0x05, 0x1c, 0xdc, 0x07,
	0xba, 0x14, 0xc5, 0x1f, 0x89, 0x1a, 0x40, 0xd0, 0xb4, 0x7f, 0xfb, 0xae, 0x82, 0x42, 0x88, 0x11,
	0x18, 0xf8, 0x02, 0xac, 0x28, 0xb4, 0xfe, 0x40, 0xb4, 0xf2, 0xe7, 0xc1, 0x45, 0xe5, 0x14, 0xbc,
	0xf6, 0xd0, 0xe0, 0x1c, 0xa4, 0x82, 0xf7, 0xc1, 0x7f, 0xc5, 0x59, 0xcd, 0x33, 0xcb, 0xc6, 0xb6,
	0x2e, 0xba, 0x6b, 0x2b, 0x1f, 0xdb, 0xbc, 0x86, 0x6a, 0x80, 0x6f, 0x5d, 0x25, 0xab, 0x7f, 0x64,
	0x95, 0xbc, 0x07, 0xdb, 0x0b, 0xab, 0xe4, 0x5a, 0x94, 0x4c, 0x64, 0x43, 0x88, 0x6c, 0x17, 0x76,
	0xcb, 0x22, 0x51, 0xaa, 0x34, 0x9d, 0x9f, 0x04, 0xd3, 0xd9, 0x84, 0x51, 0xd9, 0x6d, 0x6b, 0xa2,
	0xdb, 0x44, 0x13, 0x3a, 0xf2, 0x5e, 0xf4, 0xdb, 0xbf, 0xa0, 0x82, 0xe2, 0x58, 0xba, 0xac, 0x0b,
	0x97, 0x32, 0x8a, 0x63, 0x6e, 0xb2, 0xcc, 0x0f, 0xcf, 0xfc, 0x90, 0x05, 0xb3, 0xa1, 0x39, 0x22,
	0xd3, 0x5e, 0x9e, 0x62, 0xf1, 0x58, 0xf8, 0x2f, 0x30, 0x5c, 0x11, 0x1f, 0xcf, 0x7f, 0x05, 0x00,
	0x00, 0xff, 0xff, 0xb5, 0xbf, 0x08, 0x4b, 0x21, 0x08, 0x00, 0x00,
}
