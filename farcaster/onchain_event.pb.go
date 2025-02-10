// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v5.29.3
// source: onchain_event.proto

package farcaster

import (
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

type OnChainEventType int32

const (
	OnChainEventType_EVENT_TYPE_NONE            OnChainEventType = 0
	OnChainEventType_EVENT_TYPE_SIGNER          OnChainEventType = 1
	OnChainEventType_EVENT_TYPE_SIGNER_MIGRATED OnChainEventType = 2
	OnChainEventType_EVENT_TYPE_ID_REGISTER     OnChainEventType = 3
	OnChainEventType_EVENT_TYPE_STORAGE_RENT    OnChainEventType = 4
)

// Enum value maps for OnChainEventType.
var (
	OnChainEventType_name = map[int32]string{
		0: "EVENT_TYPE_NONE",
		1: "EVENT_TYPE_SIGNER",
		2: "EVENT_TYPE_SIGNER_MIGRATED",
		3: "EVENT_TYPE_ID_REGISTER",
		4: "EVENT_TYPE_STORAGE_RENT",
	}
	OnChainEventType_value = map[string]int32{
		"EVENT_TYPE_NONE":            0,
		"EVENT_TYPE_SIGNER":          1,
		"EVENT_TYPE_SIGNER_MIGRATED": 2,
		"EVENT_TYPE_ID_REGISTER":     3,
		"EVENT_TYPE_STORAGE_RENT":    4,
	}
)

func (x OnChainEventType) Enum() *OnChainEventType {
	p := new(OnChainEventType)
	*p = x
	return p
}

func (x OnChainEventType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (OnChainEventType) Descriptor() protoreflect.EnumDescriptor {
	return file_onchain_event_proto_enumTypes[0].Descriptor()
}

func (OnChainEventType) Type() protoreflect.EnumType {
	return &file_onchain_event_proto_enumTypes[0]
}

func (x OnChainEventType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use OnChainEventType.Descriptor instead.
func (OnChainEventType) EnumDescriptor() ([]byte, []int) {
	return file_onchain_event_proto_rawDescGZIP(), []int{0}
}

type SignerEventType int32

const (
	SignerEventType_SIGNER_EVENT_TYPE_NONE        SignerEventType = 0
	SignerEventType_SIGNER_EVENT_TYPE_ADD         SignerEventType = 1
	SignerEventType_SIGNER_EVENT_TYPE_REMOVE      SignerEventType = 2
	SignerEventType_SIGNER_EVENT_TYPE_ADMIN_RESET SignerEventType = 3
)

// Enum value maps for SignerEventType.
var (
	SignerEventType_name = map[int32]string{
		0: "SIGNER_EVENT_TYPE_NONE",
		1: "SIGNER_EVENT_TYPE_ADD",
		2: "SIGNER_EVENT_TYPE_REMOVE",
		3: "SIGNER_EVENT_TYPE_ADMIN_RESET",
	}
	SignerEventType_value = map[string]int32{
		"SIGNER_EVENT_TYPE_NONE":        0,
		"SIGNER_EVENT_TYPE_ADD":         1,
		"SIGNER_EVENT_TYPE_REMOVE":      2,
		"SIGNER_EVENT_TYPE_ADMIN_RESET": 3,
	}
)

func (x SignerEventType) Enum() *SignerEventType {
	p := new(SignerEventType)
	*p = x
	return p
}

func (x SignerEventType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SignerEventType) Descriptor() protoreflect.EnumDescriptor {
	return file_onchain_event_proto_enumTypes[1].Descriptor()
}

func (SignerEventType) Type() protoreflect.EnumType {
	return &file_onchain_event_proto_enumTypes[1]
}

func (x SignerEventType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SignerEventType.Descriptor instead.
func (SignerEventType) EnumDescriptor() ([]byte, []int) {
	return file_onchain_event_proto_rawDescGZIP(), []int{1}
}

type IdRegisterEventType int32

const (
	IdRegisterEventType_ID_REGISTER_EVENT_TYPE_NONE            IdRegisterEventType = 0
	IdRegisterEventType_ID_REGISTER_EVENT_TYPE_REGISTER        IdRegisterEventType = 1
	IdRegisterEventType_ID_REGISTER_EVENT_TYPE_TRANSFER        IdRegisterEventType = 2
	IdRegisterEventType_ID_REGISTER_EVENT_TYPE_CHANGE_RECOVERY IdRegisterEventType = 3
)

// Enum value maps for IdRegisterEventType.
var (
	IdRegisterEventType_name = map[int32]string{
		0: "ID_REGISTER_EVENT_TYPE_NONE",
		1: "ID_REGISTER_EVENT_TYPE_REGISTER",
		2: "ID_REGISTER_EVENT_TYPE_TRANSFER",
		3: "ID_REGISTER_EVENT_TYPE_CHANGE_RECOVERY",
	}
	IdRegisterEventType_value = map[string]int32{
		"ID_REGISTER_EVENT_TYPE_NONE":            0,
		"ID_REGISTER_EVENT_TYPE_REGISTER":        1,
		"ID_REGISTER_EVENT_TYPE_TRANSFER":        2,
		"ID_REGISTER_EVENT_TYPE_CHANGE_RECOVERY": 3,
	}
)

func (x IdRegisterEventType) Enum() *IdRegisterEventType {
	p := new(IdRegisterEventType)
	*p = x
	return p
}

func (x IdRegisterEventType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (IdRegisterEventType) Descriptor() protoreflect.EnumDescriptor {
	return file_onchain_event_proto_enumTypes[2].Descriptor()
}

func (IdRegisterEventType) Type() protoreflect.EnumType {
	return &file_onchain_event_proto_enumTypes[2]
}

func (x IdRegisterEventType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use IdRegisterEventType.Descriptor instead.
func (IdRegisterEventType) EnumDescriptor() ([]byte, []int) {
	return file_onchain_event_proto_rawDescGZIP(), []int{2}
}

type OnChainEvent struct {
	state           protoimpl.MessageState `protogen:"open.v1"`
	Type            OnChainEventType       `protobuf:"varint,1,opt,name=type,proto3,enum=OnChainEventType" json:"type,omitempty"`
	ChainId         uint32                 `protobuf:"varint,2,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	BlockNumber     uint32                 `protobuf:"varint,3,opt,name=block_number,json=blockNumber,proto3" json:"block_number,omitempty"`
	BlockHash       []byte                 `protobuf:"bytes,4,opt,name=block_hash,json=blockHash,proto3" json:"block_hash,omitempty"`
	BlockTimestamp  uint64                 `protobuf:"varint,5,opt,name=block_timestamp,json=blockTimestamp,proto3" json:"block_timestamp,omitempty"`
	TransactionHash []byte                 `protobuf:"bytes,6,opt,name=transaction_hash,json=transactionHash,proto3" json:"transaction_hash,omitempty"`
	LogIndex        uint32                 `protobuf:"varint,7,opt,name=log_index,json=logIndex,proto3" json:"log_index,omitempty"`
	Fid             uint64                 `protobuf:"varint,8,opt,name=fid,proto3" json:"fid,omitempty"`
	// Types that are valid to be assigned to Body:
	//
	//	*OnChainEvent_SignerEventBody
	//	*OnChainEvent_SignerMigratedEventBody
	//	*OnChainEvent_IdRegisterEventBody
	//	*OnChainEvent_StorageRentEventBody
	Body          isOnChainEvent_Body `protobuf_oneof:"body"`
	TxIndex       uint32              `protobuf:"varint,13,opt,name=tx_index,json=txIndex,proto3" json:"tx_index,omitempty"`
	Version       uint32              `protobuf:"varint,14,opt,name=version,proto3" json:"version,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *OnChainEvent) Reset() {
	*x = OnChainEvent{}
	mi := &file_onchain_event_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *OnChainEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OnChainEvent) ProtoMessage() {}

func (x *OnChainEvent) ProtoReflect() protoreflect.Message {
	mi := &file_onchain_event_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OnChainEvent.ProtoReflect.Descriptor instead.
func (*OnChainEvent) Descriptor() ([]byte, []int) {
	return file_onchain_event_proto_rawDescGZIP(), []int{0}
}

func (x *OnChainEvent) GetType() OnChainEventType {
	if x != nil {
		return x.Type
	}
	return OnChainEventType_EVENT_TYPE_NONE
}

func (x *OnChainEvent) GetChainId() uint32 {
	if x != nil {
		return x.ChainId
	}
	return 0
}

func (x *OnChainEvent) GetBlockNumber() uint32 {
	if x != nil {
		return x.BlockNumber
	}
	return 0
}

func (x *OnChainEvent) GetBlockHash() []byte {
	if x != nil {
		return x.BlockHash
	}
	return nil
}

func (x *OnChainEvent) GetBlockTimestamp() uint64 {
	if x != nil {
		return x.BlockTimestamp
	}
	return 0
}

func (x *OnChainEvent) GetTransactionHash() []byte {
	if x != nil {
		return x.TransactionHash
	}
	return nil
}

func (x *OnChainEvent) GetLogIndex() uint32 {
	if x != nil {
		return x.LogIndex
	}
	return 0
}

func (x *OnChainEvent) GetFid() uint64 {
	if x != nil {
		return x.Fid
	}
	return 0
}

func (x *OnChainEvent) GetBody() isOnChainEvent_Body {
	if x != nil {
		return x.Body
	}
	return nil
}

func (x *OnChainEvent) GetSignerEventBody() *SignerEventBody {
	if x != nil {
		if x, ok := x.Body.(*OnChainEvent_SignerEventBody); ok {
			return x.SignerEventBody
		}
	}
	return nil
}

func (x *OnChainEvent) GetSignerMigratedEventBody() *SignerMigratedEventBody {
	if x != nil {
		if x, ok := x.Body.(*OnChainEvent_SignerMigratedEventBody); ok {
			return x.SignerMigratedEventBody
		}
	}
	return nil
}

func (x *OnChainEvent) GetIdRegisterEventBody() *IdRegisterEventBody {
	if x != nil {
		if x, ok := x.Body.(*OnChainEvent_IdRegisterEventBody); ok {
			return x.IdRegisterEventBody
		}
	}
	return nil
}

func (x *OnChainEvent) GetStorageRentEventBody() *StorageRentEventBody {
	if x != nil {
		if x, ok := x.Body.(*OnChainEvent_StorageRentEventBody); ok {
			return x.StorageRentEventBody
		}
	}
	return nil
}

func (x *OnChainEvent) GetTxIndex() uint32 {
	if x != nil {
		return x.TxIndex
	}
	return 0
}

func (x *OnChainEvent) GetVersion() uint32 {
	if x != nil {
		return x.Version
	}
	return 0
}

type isOnChainEvent_Body interface {
	isOnChainEvent_Body()
}

type OnChainEvent_SignerEventBody struct {
	SignerEventBody *SignerEventBody `protobuf:"bytes,9,opt,name=signer_event_body,json=signerEventBody,proto3,oneof"`
}

type OnChainEvent_SignerMigratedEventBody struct {
	SignerMigratedEventBody *SignerMigratedEventBody `protobuf:"bytes,10,opt,name=signer_migrated_event_body,json=signerMigratedEventBody,proto3,oneof"`
}

type OnChainEvent_IdRegisterEventBody struct {
	IdRegisterEventBody *IdRegisterEventBody `protobuf:"bytes,11,opt,name=id_register_event_body,json=idRegisterEventBody,proto3,oneof"`
}

type OnChainEvent_StorageRentEventBody struct {
	StorageRentEventBody *StorageRentEventBody `protobuf:"bytes,12,opt,name=storage_rent_event_body,json=storageRentEventBody,proto3,oneof"`
}

func (*OnChainEvent_SignerEventBody) isOnChainEvent_Body() {}

func (*OnChainEvent_SignerMigratedEventBody) isOnChainEvent_Body() {}

func (*OnChainEvent_IdRegisterEventBody) isOnChainEvent_Body() {}

func (*OnChainEvent_StorageRentEventBody) isOnChainEvent_Body() {}

type SignerEventBody struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Key           []byte                 `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	KeyType       uint32                 `protobuf:"varint,2,opt,name=key_type,json=keyType,proto3" json:"key_type,omitempty"`
	EventType     SignerEventType        `protobuf:"varint,3,opt,name=event_type,json=eventType,proto3,enum=SignerEventType" json:"event_type,omitempty"`
	Metadata      []byte                 `protobuf:"bytes,4,opt,name=metadata,proto3" json:"metadata,omitempty"`
	MetadataType  uint32                 `protobuf:"varint,5,opt,name=metadata_type,json=metadataType,proto3" json:"metadata_type,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SignerEventBody) Reset() {
	*x = SignerEventBody{}
	mi := &file_onchain_event_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SignerEventBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignerEventBody) ProtoMessage() {}

func (x *SignerEventBody) ProtoReflect() protoreflect.Message {
	mi := &file_onchain_event_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignerEventBody.ProtoReflect.Descriptor instead.
func (*SignerEventBody) Descriptor() ([]byte, []int) {
	return file_onchain_event_proto_rawDescGZIP(), []int{1}
}

func (x *SignerEventBody) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *SignerEventBody) GetKeyType() uint32 {
	if x != nil {
		return x.KeyType
	}
	return 0
}

func (x *SignerEventBody) GetEventType() SignerEventType {
	if x != nil {
		return x.EventType
	}
	return SignerEventType_SIGNER_EVENT_TYPE_NONE
}

func (x *SignerEventBody) GetMetadata() []byte {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *SignerEventBody) GetMetadataType() uint32 {
	if x != nil {
		return x.MetadataType
	}
	return 0
}

type SignerMigratedEventBody struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	MigratedAt    uint32                 `protobuf:"varint,1,opt,name=migratedAt,proto3" json:"migratedAt,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SignerMigratedEventBody) Reset() {
	*x = SignerMigratedEventBody{}
	mi := &file_onchain_event_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SignerMigratedEventBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignerMigratedEventBody) ProtoMessage() {}

func (x *SignerMigratedEventBody) ProtoReflect() protoreflect.Message {
	mi := &file_onchain_event_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignerMigratedEventBody.ProtoReflect.Descriptor instead.
func (*SignerMigratedEventBody) Descriptor() ([]byte, []int) {
	return file_onchain_event_proto_rawDescGZIP(), []int{2}
}

func (x *SignerMigratedEventBody) GetMigratedAt() uint32 {
	if x != nil {
		return x.MigratedAt
	}
	return 0
}

type IdRegisterEventBody struct {
	state           protoimpl.MessageState `protogen:"open.v1"`
	To              []byte                 `protobuf:"bytes,1,opt,name=to,proto3" json:"to,omitempty"`
	EventType       IdRegisterEventType    `protobuf:"varint,2,opt,name=event_type,json=eventType,proto3,enum=IdRegisterEventType" json:"event_type,omitempty"`
	From            []byte                 `protobuf:"bytes,3,opt,name=from,proto3" json:"from,omitempty"`
	RecoveryAddress []byte                 `protobuf:"bytes,4,opt,name=recovery_address,json=recoveryAddress,proto3" json:"recovery_address,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *IdRegisterEventBody) Reset() {
	*x = IdRegisterEventBody{}
	mi := &file_onchain_event_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *IdRegisterEventBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IdRegisterEventBody) ProtoMessage() {}

func (x *IdRegisterEventBody) ProtoReflect() protoreflect.Message {
	mi := &file_onchain_event_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IdRegisterEventBody.ProtoReflect.Descriptor instead.
func (*IdRegisterEventBody) Descriptor() ([]byte, []int) {
	return file_onchain_event_proto_rawDescGZIP(), []int{3}
}

func (x *IdRegisterEventBody) GetTo() []byte {
	if x != nil {
		return x.To
	}
	return nil
}

func (x *IdRegisterEventBody) GetEventType() IdRegisterEventType {
	if x != nil {
		return x.EventType
	}
	return IdRegisterEventType_ID_REGISTER_EVENT_TYPE_NONE
}

func (x *IdRegisterEventBody) GetFrom() []byte {
	if x != nil {
		return x.From
	}
	return nil
}

func (x *IdRegisterEventBody) GetRecoveryAddress() []byte {
	if x != nil {
		return x.RecoveryAddress
	}
	return nil
}

type StorageRentEventBody struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Payer         []byte                 `protobuf:"bytes,1,opt,name=payer,proto3" json:"payer,omitempty"`
	Units         uint32                 `protobuf:"varint,2,opt,name=units,proto3" json:"units,omitempty"`
	Expiry        uint32                 `protobuf:"varint,3,opt,name=expiry,proto3" json:"expiry,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StorageRentEventBody) Reset() {
	*x = StorageRentEventBody{}
	mi := &file_onchain_event_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StorageRentEventBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StorageRentEventBody) ProtoMessage() {}

func (x *StorageRentEventBody) ProtoReflect() protoreflect.Message {
	mi := &file_onchain_event_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StorageRentEventBody.ProtoReflect.Descriptor instead.
func (*StorageRentEventBody) Descriptor() ([]byte, []int) {
	return file_onchain_event_proto_rawDescGZIP(), []int{4}
}

func (x *StorageRentEventBody) GetPayer() []byte {
	if x != nil {
		return x.Payer
	}
	return nil
}

func (x *StorageRentEventBody) GetUnits() uint32 {
	if x != nil {
		return x.Units
	}
	return 0
}

func (x *StorageRentEventBody) GetExpiry() uint32 {
	if x != nil {
		return x.Expiry
	}
	return 0
}

var File_onchain_event_proto protoreflect.FileDescriptor

var file_onchain_event_proto_rawDesc = string([]byte{
	0x0a, 0x13, 0x6f, 0x6e, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x88, 0x05, 0x0a, 0x0c, 0x4f, 0x6e, 0x43, 0x68, 0x61, 0x69,
	0x6e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x25, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x11, 0x2e, 0x4f, 0x6e, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x19, 0x0a,
	0x08, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x07, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x49, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x62, 0x6c, 0x6f, 0x63,
	0x6b, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0b,
	0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x1d, 0x0a, 0x0a, 0x62,
	0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x09, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68, 0x12, 0x27, 0x0a, 0x0f, 0x62, 0x6c,
	0x6f, 0x63, 0x6b, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x0e, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x12, 0x29, 0x0a, 0x10, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0f, 0x74,
	0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x61, 0x73, 0x68, 0x12, 0x1b,
	0x0a, 0x09, 0x6c, 0x6f, 0x67, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x08, 0x6c, 0x6f, 0x67, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x10, 0x0a, 0x03, 0x66,
	0x69, 0x64, 0x18, 0x08, 0x20, 0x01, 0x28, 0x04, 0x52, 0x03, 0x66, 0x69, 0x64, 0x12, 0x3e, 0x0a,
	0x11, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x72, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x62, 0x6f,
	0x64, 0x79, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x65,
	0x72, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x42, 0x6f, 0x64, 0x79, 0x48, 0x00, 0x52, 0x0f, 0x73, 0x69,
	0x67, 0x6e, 0x65, 0x72, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x57, 0x0a,
	0x1a, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x72, 0x5f, 0x6d, 0x69, 0x67, 0x72, 0x61, 0x74, 0x65, 0x64,
	0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x0a, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x18, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x65, 0x72, 0x4d, 0x69, 0x67, 0x72, 0x61, 0x74,
	0x65, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x42, 0x6f, 0x64, 0x79, 0x48, 0x00, 0x52, 0x17, 0x73,
	0x69, 0x67, 0x6e, 0x65, 0x72, 0x4d, 0x69, 0x67, 0x72, 0x61, 0x74, 0x65, 0x64, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x4b, 0x0a, 0x16, 0x69, 0x64, 0x5f, 0x72, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x62, 0x6f, 0x64, 0x79,
	0x18, 0x0b, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x49, 0x64, 0x52, 0x65, 0x67, 0x69, 0x73,
	0x74, 0x65, 0x72, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x42, 0x6f, 0x64, 0x79, 0x48, 0x00, 0x52, 0x13,
	0x69, 0x64, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x42,
	0x6f, 0x64, 0x79, 0x12, 0x4e, 0x0a, 0x17, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x5f, 0x72,
	0x65, 0x6e, 0x74, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x0c,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x52, 0x65,
	0x6e, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x42, 0x6f, 0x64, 0x79, 0x48, 0x00, 0x52, 0x14, 0x73,
	0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x52, 0x65, 0x6e, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x42,
	0x6f, 0x64, 0x79, 0x12, 0x19, 0x0a, 0x08, 0x74, 0x78, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18,
	0x0d, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x74, 0x78, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x18,
	0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x42, 0x06, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79,
	0x22, 0xb0, 0x01, 0x0a, 0x0f, 0x53, 0x69, 0x67, 0x6e, 0x65, 0x72, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x42, 0x6f, 0x64, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x19, 0x0a, 0x08, 0x6b, 0x65, 0x79, 0x5f, 0x74, 0x79,
	0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x6b, 0x65, 0x79, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x2f, 0x0a, 0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x10, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x65, 0x72, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x23,
	0x0a, 0x0d, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0c, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x54,
	0x79, 0x70, 0x65, 0x22, 0x39, 0x0a, 0x17, 0x53, 0x69, 0x67, 0x6e, 0x65, 0x72, 0x4d, 0x69, 0x67,
	0x72, 0x61, 0x74, 0x65, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x1e,
	0x0a, 0x0a, 0x6d, 0x69, 0x67, 0x72, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x0a, 0x6d, 0x69, 0x67, 0x72, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0x99,
	0x01, 0x0a, 0x13, 0x49, 0x64, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x0e, 0x0a, 0x02, 0x74, 0x6f, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x02, 0x74, 0x6f, 0x12, 0x33, 0x0a, 0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x14, 0x2e, 0x49, 0x64, 0x52,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65,
	0x52, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x66,
	0x72, 0x6f, 0x6d, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x12,
	0x29, 0x0a, 0x10, 0x72, 0x65, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x5f, 0x61, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0f, 0x72, 0x65, 0x63, 0x6f, 0x76,
	0x65, 0x72, 0x79, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x22, 0x5a, 0x0a, 0x14, 0x53, 0x74,
	0x6f, 0x72, 0x61, 0x67, 0x65, 0x52, 0x65, 0x6e, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x42, 0x6f,
	0x64, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x61, 0x79, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x05, 0x70, 0x61, 0x79, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x75, 0x6e, 0x69, 0x74,
	0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x75, 0x6e, 0x69, 0x74, 0x73, 0x12, 0x16,
	0x0a, 0x06, 0x65, 0x78, 0x70, 0x69, 0x72, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06,
	0x65, 0x78, 0x70, 0x69, 0x72, 0x79, 0x2a, 0x97, 0x01, 0x0a, 0x10, 0x4f, 0x6e, 0x43, 0x68, 0x61,
	0x69, 0x6e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x13, 0x0a, 0x0f, 0x45,
	0x56, 0x45, 0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x4e, 0x4f, 0x4e, 0x45, 0x10, 0x00,
	0x12, 0x15, 0x0a, 0x11, 0x45, 0x56, 0x45, 0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x53,
	0x49, 0x47, 0x4e, 0x45, 0x52, 0x10, 0x01, 0x12, 0x1e, 0x0a, 0x1a, 0x45, 0x56, 0x45, 0x4e, 0x54,
	0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x53, 0x49, 0x47, 0x4e, 0x45, 0x52, 0x5f, 0x4d, 0x49, 0x47,
	0x52, 0x41, 0x54, 0x45, 0x44, 0x10, 0x02, 0x12, 0x1a, 0x0a, 0x16, 0x45, 0x56, 0x45, 0x4e, 0x54,
	0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x49, 0x44, 0x5f, 0x52, 0x45, 0x47, 0x49, 0x53, 0x54, 0x45,
	0x52, 0x10, 0x03, 0x12, 0x1b, 0x0a, 0x17, 0x45, 0x56, 0x45, 0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50,
	0x45, 0x5f, 0x53, 0x54, 0x4f, 0x52, 0x41, 0x47, 0x45, 0x5f, 0x52, 0x45, 0x4e, 0x54, 0x10, 0x04,
	0x2a, 0x89, 0x01, 0x0a, 0x0f, 0x53, 0x69, 0x67, 0x6e, 0x65, 0x72, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x1a, 0x0a, 0x16, 0x53, 0x49, 0x47, 0x4e, 0x45, 0x52, 0x5f, 0x45,
	0x56, 0x45, 0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x4e, 0x4f, 0x4e, 0x45, 0x10, 0x00,
	0x12, 0x19, 0x0a, 0x15, 0x53, 0x49, 0x47, 0x4e, 0x45, 0x52, 0x5f, 0x45, 0x56, 0x45, 0x4e, 0x54,
	0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x41, 0x44, 0x44, 0x10, 0x01, 0x12, 0x1c, 0x0a, 0x18, 0x53,
	0x49, 0x47, 0x4e, 0x45, 0x52, 0x5f, 0x45, 0x56, 0x45, 0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45,
	0x5f, 0x52, 0x45, 0x4d, 0x4f, 0x56, 0x45, 0x10, 0x02, 0x12, 0x21, 0x0a, 0x1d, 0x53, 0x49, 0x47,
	0x4e, 0x45, 0x52, 0x5f, 0x45, 0x56, 0x45, 0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x41,
	0x44, 0x4d, 0x49, 0x4e, 0x5f, 0x52, 0x45, 0x53, 0x45, 0x54, 0x10, 0x03, 0x2a, 0xac, 0x01, 0x0a,
	0x13, 0x49, 0x64, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x1f, 0x0a, 0x1b, 0x49, 0x44, 0x5f, 0x52, 0x45, 0x47, 0x49, 0x53,
	0x54, 0x45, 0x52, 0x5f, 0x45, 0x56, 0x45, 0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x4e,
	0x4f, 0x4e, 0x45, 0x10, 0x00, 0x12, 0x23, 0x0a, 0x1f, 0x49, 0x44, 0x5f, 0x52, 0x45, 0x47, 0x49,
	0x53, 0x54, 0x45, 0x52, 0x5f, 0x45, 0x56, 0x45, 0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f,
	0x52, 0x45, 0x47, 0x49, 0x53, 0x54, 0x45, 0x52, 0x10, 0x01, 0x12, 0x23, 0x0a, 0x1f, 0x49, 0x44,
	0x5f, 0x52, 0x45, 0x47, 0x49, 0x53, 0x54, 0x45, 0x52, 0x5f, 0x45, 0x56, 0x45, 0x4e, 0x54, 0x5f,
	0x54, 0x59, 0x50, 0x45, 0x5f, 0x54, 0x52, 0x41, 0x4e, 0x53, 0x46, 0x45, 0x52, 0x10, 0x02, 0x12,
	0x2a, 0x0a, 0x26, 0x49, 0x44, 0x5f, 0x52, 0x45, 0x47, 0x49, 0x53, 0x54, 0x45, 0x52, 0x5f, 0x45,
	0x56, 0x45, 0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x43, 0x48, 0x41, 0x4e, 0x47, 0x45,
	0x5f, 0x52, 0x45, 0x43, 0x4f, 0x56, 0x45, 0x52, 0x59, 0x10, 0x03, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
})

var (
	file_onchain_event_proto_rawDescOnce sync.Once
	file_onchain_event_proto_rawDescData []byte
)

func file_onchain_event_proto_rawDescGZIP() []byte {
	file_onchain_event_proto_rawDescOnce.Do(func() {
		file_onchain_event_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_onchain_event_proto_rawDesc), len(file_onchain_event_proto_rawDesc)))
	})
	return file_onchain_event_proto_rawDescData
}

var file_onchain_event_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_onchain_event_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_onchain_event_proto_goTypes = []any{
	(OnChainEventType)(0),           // 0: OnChainEventType
	(SignerEventType)(0),            // 1: SignerEventType
	(IdRegisterEventType)(0),        // 2: IdRegisterEventType
	(*OnChainEvent)(nil),            // 3: OnChainEvent
	(*SignerEventBody)(nil),         // 4: SignerEventBody
	(*SignerMigratedEventBody)(nil), // 5: SignerMigratedEventBody
	(*IdRegisterEventBody)(nil),     // 6: IdRegisterEventBody
	(*StorageRentEventBody)(nil),    // 7: StorageRentEventBody
}
var file_onchain_event_proto_depIdxs = []int32{
	0, // 0: OnChainEvent.type:type_name -> OnChainEventType
	4, // 1: OnChainEvent.signer_event_body:type_name -> SignerEventBody
	5, // 2: OnChainEvent.signer_migrated_event_body:type_name -> SignerMigratedEventBody
	6, // 3: OnChainEvent.id_register_event_body:type_name -> IdRegisterEventBody
	7, // 4: OnChainEvent.storage_rent_event_body:type_name -> StorageRentEventBody
	1, // 5: SignerEventBody.event_type:type_name -> SignerEventType
	2, // 6: IdRegisterEventBody.event_type:type_name -> IdRegisterEventType
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_onchain_event_proto_init() }
func file_onchain_event_proto_init() {
	if File_onchain_event_proto != nil {
		return
	}
	file_onchain_event_proto_msgTypes[0].OneofWrappers = []any{
		(*OnChainEvent_SignerEventBody)(nil),
		(*OnChainEvent_SignerMigratedEventBody)(nil),
		(*OnChainEvent_IdRegisterEventBody)(nil),
		(*OnChainEvent_StorageRentEventBody)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_onchain_event_proto_rawDesc), len(file_onchain_event_proto_rawDesc)),
			NumEnums:      3,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_onchain_event_proto_goTypes,
		DependencyIndexes: file_onchain_event_proto_depIdxs,
		EnumInfos:         file_onchain_event_proto_enumTypes,
		MessageInfos:      file_onchain_event_proto_msgTypes,
	}.Build()
	File_onchain_event_proto = out.File
	file_onchain_event_proto_goTypes = nil
	file_onchain_event_proto_depIdxs = nil
}
