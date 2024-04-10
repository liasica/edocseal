// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v4.25.3
// source: contract.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// 合同签署主体
type ContractSignSubject int32

const (
	ContractSignSubject_Enterprise ContractSignSubject = 0 // 企业
	ContractSignSubject_Personal   ContractSignSubject = 1 // 个人
)

// Enum value maps for ContractSignSubject.
var (
	ContractSignSubject_name = map[int32]string{
		0: "Enterprise",
		1: "Personal",
	}
	ContractSignSubject_value = map[string]int32{
		"Enterprise": 0,
		"Personal":   1,
	}
)

func (x ContractSignSubject) Enum() *ContractSignSubject {
	p := new(ContractSignSubject)
	*p = x
	return p
}

func (x ContractSignSubject) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ContractSignSubject) Descriptor() protoreflect.EnumDescriptor {
	return file_contract_proto_enumTypes[0].Descriptor()
}

func (ContractSignSubject) Type() protoreflect.EnumType {
	return &file_contract_proto_enumTypes[0]
}

func (x ContractSignSubject) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ContractSignSubject.Descriptor instead.
func (ContractSignSubject) EnumDescriptor() ([]byte, []int) {
	return file_contract_proto_rawDescGZIP(), []int{0}
}

// 创建合同字段
type ContractFromField struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Value:
	//
	//	*ContractFromField_Text
	//	*ContractFromField_Checkbox
	Value isContractFromField_Value `protobuf_oneof:"value"`
}

func (x *ContractFromField) Reset() {
	*x = ContractFromField{}
	if protoimpl.UnsafeEnabled {
		mi := &file_contract_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ContractFromField) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContractFromField) ProtoMessage() {}

func (x *ContractFromField) ProtoReflect() protoreflect.Message {
	mi := &file_contract_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContractFromField.ProtoReflect.Descriptor instead.
func (*ContractFromField) Descriptor() ([]byte, []int) {
	return file_contract_proto_rawDescGZIP(), []int{0}
}

func (m *ContractFromField) GetValue() isContractFromField_Value {
	if m != nil {
		return m.Value
	}
	return nil
}

func (x *ContractFromField) GetText() string {
	if x, ok := x.GetValue().(*ContractFromField_Text); ok {
		return x.Text
	}
	return ""
}

func (x *ContractFromField) GetCheckbox() bool {
	if x, ok := x.GetValue().(*ContractFromField_Checkbox); ok {
		return x.Checkbox
	}
	return false
}

type isContractFromField_Value interface {
	isContractFromField_Value()
}

type ContractFromField_Text struct {
	Text string `protobuf:"bytes,1,opt,name=text,proto3,oneof"` // 文本
}

type ContractFromField_Checkbox struct {
	Checkbox bool `protobuf:"varint,2,opt,name=checkbox,proto3,oneof"` // 勾选
}

func (*ContractFromField_Text) isContractFromField_Value() {}

func (*ContractFromField_Checkbox) isContractFromField_Value() {}

// 创建合同请求
type ContractCreateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TemplateId string                        `protobuf:"bytes,1,opt,name=template_id,json=templateId,proto3" json:"template_id,omitempty"`                                                               // 模板编号
	Values     map[string]*ContractFromField `protobuf:"bytes,2,rep,name=values,proto3" json:"values,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"` // 合同模板字段，key为字段名，value为字段值
}

func (x *ContractCreateRequest) Reset() {
	*x = ContractCreateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_contract_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ContractCreateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContractCreateRequest) ProtoMessage() {}

func (x *ContractCreateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_contract_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContractCreateRequest.ProtoReflect.Descriptor instead.
func (*ContractCreateRequest) Descriptor() ([]byte, []int) {
	return file_contract_proto_rawDescGZIP(), []int{1}
}

func (x *ContractCreateRequest) GetTemplateId() string {
	if x != nil {
		return x.TemplateId
	}
	return ""
}

func (x *ContractCreateRequest) GetValues() map[string]*ContractFromField {
	if x != nil {
		return x.Values
	}
	return nil
}

// 创建合同响应
type ContractCreateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url   string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`                  // 创建合同合同文件地址
	DocId string `protobuf:"bytes,2,opt,name=doc_id,json=docId,proto3" json:"doc_id,omitempty"` // 待签约文档编号
}

func (x *ContractCreateResponse) Reset() {
	*x = ContractCreateResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_contract_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ContractCreateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContractCreateResponse) ProtoMessage() {}

func (x *ContractCreateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_contract_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContractCreateResponse.ProtoReflect.Descriptor instead.
func (*ContractCreateResponse) Descriptor() ([]byte, []int) {
	return file_contract_proto_rawDescGZIP(), []int{2}
}

func (x *ContractCreateResponse) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *ContractCreateResponse) GetDocId() string {
	if x != nil {
		return x.DocId
	}
	return ""
}

// 合同签署实体
type ContractSignEntity struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Subject  ContractSignSubject `protobuf:"varint,1,opt,name=subject,proto3,enum=pb.ContractSignSubject" json:"subject,omitempty"` // 签署主体，必填
	Field    string              `protobuf:"bytes,2,opt,name=field,proto3" json:"field,omitempty"`                                  // 签章字段，必填
	Image    *string             `protobuf:"bytes,3,opt,name=image,proto3,oneof" json:"image,omitempty"`                            // 签章图片，BASE64编码，个人必填
	Name     *string             `protobuf:"bytes,4,opt,name=name,proto3,oneof" json:"name,omitempty"`                              // 签署人，个人必填
	Province *string             `protobuf:"bytes,5,opt,name=province,proto3,oneof" json:"province,omitempty"`                      // 省份，个人必填
	City     *string             `protobuf:"bytes,6,opt,name=city,proto3,oneof" json:"city,omitempty"`                              // 城市，个人必填
	Address  *string             `protobuf:"bytes,7,opt,name=address,proto3,oneof" json:"address,omitempty"`                        // 地址，个人必填
	Phone    *string             `protobuf:"bytes,8,opt,name=phone,proto3,oneof" json:"phone,omitempty"`                            // 手机号，个人必填
	Idcard   *string             `protobuf:"bytes,9,opt,name=idcard,proto3,oneof" json:"idcard,omitempty"`                          // 身份证号，个人必填
}

func (x *ContractSignEntity) Reset() {
	*x = ContractSignEntity{}
	if protoimpl.UnsafeEnabled {
		mi := &file_contract_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ContractSignEntity) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContractSignEntity) ProtoMessage() {}

func (x *ContractSignEntity) ProtoReflect() protoreflect.Message {
	mi := &file_contract_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContractSignEntity.ProtoReflect.Descriptor instead.
func (*ContractSignEntity) Descriptor() ([]byte, []int) {
	return file_contract_proto_rawDescGZIP(), []int{3}
}

func (x *ContractSignEntity) GetSubject() ContractSignSubject {
	if x != nil {
		return x.Subject
	}
	return ContractSignSubject_Enterprise
}

func (x *ContractSignEntity) GetField() string {
	if x != nil {
		return x.Field
	}
	return ""
}

func (x *ContractSignEntity) GetImage() string {
	if x != nil && x.Image != nil {
		return *x.Image
	}
	return ""
}

func (x *ContractSignEntity) GetName() string {
	if x != nil && x.Name != nil {
		return *x.Name
	}
	return ""
}

func (x *ContractSignEntity) GetProvince() string {
	if x != nil && x.Province != nil {
		return *x.Province
	}
	return ""
}

func (x *ContractSignEntity) GetCity() string {
	if x != nil && x.City != nil {
		return *x.City
	}
	return ""
}

func (x *ContractSignEntity) GetAddress() string {
	if x != nil && x.Address != nil {
		return *x.Address
	}
	return ""
}

func (x *ContractSignEntity) GetPhone() string {
	if x != nil && x.Phone != nil {
		return *x.Phone
	}
	return ""
}

func (x *ContractSignEntity) GetIdcard() string {
	if x != nil && x.Idcard != nil {
		return *x.Idcard
	}
	return ""
}

// 合同签署请求
// 个签属性和团签属性只能二选一（需要进行限制判定）
type ContractSignRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DocId      string              `protobuf:"bytes,1,opt,name=doc_id,json=docId,proto3" json:"doc_id,omitempty"` // 待签约文档编号
	Enterprise *ContractSignEntity `protobuf:"bytes,2,opt,name=enterprise,proto3" json:"enterprise,omitempty"`    // 企业签署人
	Personal   *ContractSignEntity `protobuf:"bytes,3,opt,name=personal,proto3" json:"personal,omitempty"`        // 个人签署人
}

func (x *ContractSignRequest) Reset() {
	*x = ContractSignRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_contract_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ContractSignRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContractSignRequest) ProtoMessage() {}

func (x *ContractSignRequest) ProtoReflect() protoreflect.Message {
	mi := &file_contract_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContractSignRequest.ProtoReflect.Descriptor instead.
func (*ContractSignRequest) Descriptor() ([]byte, []int) {
	return file_contract_proto_rawDescGZIP(), []int{4}
}

func (x *ContractSignRequest) GetDocId() string {
	if x != nil {
		return x.DocId
	}
	return ""
}

func (x *ContractSignRequest) GetEnterprise() *ContractSignEntity {
	if x != nil {
		return x.Enterprise
	}
	return nil
}

func (x *ContractSignRequest) GetPersonal() *ContractSignEntity {
	if x != nil {
		return x.Personal
	}
	return nil
}

// 合同签署响应
type ContractSignResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status  string `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`   // 签署状态 SUCCESS:成功，FAIL:失败
	File    string `protobuf:"bytes,2,opt,name=file,proto3" json:"file,omitempty"`       // 已签署合同文件地址
	Message string `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"` // 其他消息，例如错误消息，成功时为空
}

func (x *ContractSignResponse) Reset() {
	*x = ContractSignResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_contract_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ContractSignResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContractSignResponse) ProtoMessage() {}

func (x *ContractSignResponse) ProtoReflect() protoreflect.Message {
	mi := &file_contract_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContractSignResponse.ProtoReflect.Descriptor instead.
func (*ContractSignResponse) Descriptor() ([]byte, []int) {
	return file_contract_proto_rawDescGZIP(), []int{5}
}

func (x *ContractSignResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *ContractSignResponse) GetFile() string {
	if x != nil {
		return x.File
	}
	return ""
}

func (x *ContractSignResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_contract_proto protoreflect.FileDescriptor

var file_contract_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x02, 0x70, 0x62, 0x22, 0x50, 0x0a, 0x11, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74,
	0x46, 0x72, 0x6f, 0x6d, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x14, 0x0a, 0x04, 0x74, 0x65, 0x78,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x12,
	0x1c, 0x0a, 0x08, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x62, 0x6f, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x08, 0x48, 0x00, 0x52, 0x08, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x62, 0x6f, 0x78, 0x42, 0x07, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0xc9, 0x01, 0x0a, 0x15, 0x43, 0x6f, 0x6e, 0x74, 0x72,
	0x61, 0x63, 0x74, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x49,
	0x64, 0x12, 0x3d, 0x0a, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x25, 0x2e, 0x70, 0x62, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73,
	0x1a, 0x50, 0x0a, 0x0b, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x2b, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x15, 0x2e, 0x70, 0x62, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x46, 0x72,
	0x6f, 0x6d, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02,
	0x38, 0x01, 0x22, 0x41, 0x0a, 0x16, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03,
	0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x15,
	0x0a, 0x06, 0x64, 0x6f, 0x63, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x64, 0x6f, 0x63, 0x49, 0x64, 0x22, 0xec, 0x02, 0x0a, 0x12, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61,
	0x63, 0x74, 0x53, 0x69, 0x67, 0x6e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x12, 0x31, 0x0a, 0x07,
	0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x17, 0x2e,
	0x70, 0x62, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x53, 0x69, 0x67, 0x6e, 0x53,
	0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x52, 0x07, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x12,
	0x14, 0x0a, 0x05, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x66, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x19, 0x0a, 0x05, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x05, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x88, 0x01, 0x01,
	0x12, 0x17, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x48, 0x01,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x88, 0x01, 0x01, 0x12, 0x1f, 0x0a, 0x08, 0x70, 0x72, 0x6f,
	0x76, 0x69, 0x6e, 0x63, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x48, 0x02, 0x52, 0x08, 0x70,
	0x72, 0x6f, 0x76, 0x69, 0x6e, 0x63, 0x65, 0x88, 0x01, 0x01, 0x12, 0x17, 0x0a, 0x04, 0x63, 0x69,
	0x74, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x48, 0x03, 0x52, 0x04, 0x63, 0x69, 0x74, 0x79,
	0x88, 0x01, 0x01, 0x12, 0x1d, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x09, 0x48, 0x04, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x88,
	0x01, 0x01, 0x12, 0x19, 0x0a, 0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28,
	0x09, 0x48, 0x05, 0x52, 0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x88, 0x01, 0x01, 0x12, 0x1b, 0x0a,
	0x06, 0x69, 0x64, 0x63, 0x61, 0x72, 0x64, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x48, 0x06, 0x52,
	0x06, 0x69, 0x64, 0x63, 0x61, 0x72, 0x64, 0x88, 0x01, 0x01, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x69,
	0x6d, 0x61, 0x67, 0x65, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x42, 0x0b, 0x0a,
	0x09, 0x5f, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x6e, 0x63, 0x65, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x63,
	0x69, 0x74, 0x79, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x42,
	0x08, 0x0a, 0x06, 0x5f, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x69, 0x64,
	0x63, 0x61, 0x72, 0x64, 0x22, 0x98, 0x01, 0x0a, 0x13, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63,
	0x74, 0x53, 0x69, 0x67, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x15, 0x0a, 0x06,
	0x64, 0x6f, 0x63, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x64, 0x6f,
	0x63, 0x49, 0x64, 0x12, 0x36, 0x0a, 0x0a, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x70, 0x72, 0x69, 0x73,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x70, 0x62, 0x2e, 0x43, 0x6f, 0x6e,
	0x74, 0x72, 0x61, 0x63, 0x74, 0x53, 0x69, 0x67, 0x6e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x52,
	0x0a, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x70, 0x72, 0x69, 0x73, 0x65, 0x12, 0x32, 0x0a, 0x08, 0x70,
	0x65, 0x72, 0x73, 0x6f, 0x6e, 0x61, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e,
	0x70, 0x62, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x53, 0x69, 0x67, 0x6e, 0x45,
	0x6e, 0x74, 0x69, 0x74, 0x79, 0x52, 0x08, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x61, 0x6c, 0x22,
	0x5c, 0x0a, 0x14, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x53, 0x69, 0x67, 0x6e, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12,
	0x12, 0x0a, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x66,
	0x69, 0x6c, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2a, 0x33, 0x0a,
	0x13, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x53, 0x69, 0x67, 0x6e, 0x53, 0x75, 0x62,
	0x6a, 0x65, 0x63, 0x74, 0x12, 0x0e, 0x0a, 0x0a, 0x45, 0x6e, 0x74, 0x65, 0x72, 0x70, 0x72, 0x69,
	0x73, 0x65, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x50, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x61, 0x6c,
	0x10, 0x01, 0x32, 0x8a, 0x01, 0x0a, 0x08, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x12,
	0x41, 0x0a, 0x06, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x19, 0x2e, 0x70, 0x62, 0x2e, 0x43,
	0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x70, 0x62, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61,
	0x63, 0x74, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x12, 0x3b, 0x0a, 0x04, 0x73, 0x69, 0x67, 0x6e, 0x12, 0x17, 0x2e, 0x70, 0x62, 0x2e,
	0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x53, 0x69, 0x67, 0x6e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x70, 0x62, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63,
	0x74, 0x53, 0x69, 0x67, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42,
	0x23, 0x5a, 0x21, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6c, 0x69,
	0x61, 0x73, 0x69, 0x63, 0x61, 0x2f, 0x65, 0x64, 0x6f, 0x63, 0x73, 0x65, 0x61, 0x6c, 0x2f, 0x70,
	0x62, 0x3b, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_contract_proto_rawDescOnce sync.Once
	file_contract_proto_rawDescData = file_contract_proto_rawDesc
)

func file_contract_proto_rawDescGZIP() []byte {
	file_contract_proto_rawDescOnce.Do(func() {
		file_contract_proto_rawDescData = protoimpl.X.CompressGZIP(file_contract_proto_rawDescData)
	})
	return file_contract_proto_rawDescData
}

var file_contract_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_contract_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_contract_proto_goTypes = []interface{}{
	(ContractSignSubject)(0),       // 0: pb.ContractSignSubject
	(*ContractFromField)(nil),      // 1: pb.ContractFromField
	(*ContractCreateRequest)(nil),  // 2: pb.ContractCreateRequest
	(*ContractCreateResponse)(nil), // 3: pb.ContractCreateResponse
	(*ContractSignEntity)(nil),     // 4: pb.ContractSignEntity
	(*ContractSignRequest)(nil),    // 5: pb.ContractSignRequest
	(*ContractSignResponse)(nil),   // 6: pb.ContractSignResponse
	nil,                            // 7: pb.ContractCreateRequest.ValuesEntry
}
var file_contract_proto_depIdxs = []int32{
	7, // 0: pb.ContractCreateRequest.values:type_name -> pb.ContractCreateRequest.ValuesEntry
	0, // 1: pb.ContractSignEntity.subject:type_name -> pb.ContractSignSubject
	4, // 2: pb.ContractSignRequest.enterprise:type_name -> pb.ContractSignEntity
	4, // 3: pb.ContractSignRequest.personal:type_name -> pb.ContractSignEntity
	1, // 4: pb.ContractCreateRequest.ValuesEntry.value:type_name -> pb.ContractFromField
	2, // 5: pb.Contract.create:input_type -> pb.ContractCreateRequest
	5, // 6: pb.Contract.sign:input_type -> pb.ContractSignRequest
	3, // 7: pb.Contract.create:output_type -> pb.ContractCreateResponse
	6, // 8: pb.Contract.sign:output_type -> pb.ContractSignResponse
	7, // [7:9] is the sub-list for method output_type
	5, // [5:7] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_contract_proto_init() }
func file_contract_proto_init() {
	if File_contract_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_contract_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ContractFromField); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_contract_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ContractCreateRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_contract_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ContractCreateResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_contract_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ContractSignEntity); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_contract_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ContractSignRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_contract_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ContractSignResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_contract_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*ContractFromField_Text)(nil),
		(*ContractFromField_Checkbox)(nil),
	}
	file_contract_proto_msgTypes[3].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_contract_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_contract_proto_goTypes,
		DependencyIndexes: file_contract_proto_depIdxs,
		EnumInfos:         file_contract_proto_enumTypes,
		MessageInfos:      file_contract_proto_msgTypes,
	}.Build()
	File_contract_proto = out.File
	file_contract_proto_rawDesc = nil
	file_contract_proto_goTypes = nil
	file_contract_proto_depIdxs = nil
}
