// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: core/types/v1/attestation.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	github_com_cosmos_gogoproto_types "github.com/cosmos/gogoproto/types"
	types "github.com/cosmos/ibc-go/v9/modules/core/02-client/types"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	io "io"
	math "math"
	math_bits "math/bits"
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

type Attestation struct {
	ValidatorAddress []byte  `protobuf:"bytes,1,opt,name=validator_address,json=validatorAddress,proto3" json:"validator_address,omitempty"`
	AttestedData     IBCData `protobuf:"bytes,2,opt,name=attested_data,json=attestedData,proto3" json:"attested_data"`
}

func (m *Attestation) Reset()         { *m = Attestation{} }
func (m *Attestation) String() string { return proto.CompactTextString(m) }
func (*Attestation) ProtoMessage()    {}
func (*Attestation) Descriptor() ([]byte, []int) {
	return fileDescriptor_25eb7c0454d2e150, []int{0}
}
func (m *Attestation) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Attestation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Attestation.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Attestation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Attestation.Merge(m, src)
}
func (m *Attestation) XXX_Size() int {
	return m.Size()
}
func (m *Attestation) XXX_DiscardUnknown() {
	xxx_messageInfo_Attestation.DiscardUnknown(m)
}

var xxx_messageInfo_Attestation proto.InternalMessageInfo

func (m *Attestation) GetValidatorAddress() []byte {
	if m != nil {
		return m.ValidatorAddress
	}
	return nil
}

func (m *Attestation) GetAttestedData() IBCData {
	if m != nil {
		return m.AttestedData
	}
	return IBCData{}
}

type IBCData struct {
	ChainId           string       `protobuf:"bytes,1,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	ClientId          string       `protobuf:"bytes,2,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	ClientToUpdate    string       `protobuf:"bytes,3,opt,name=client_to_update,json=clientToUpdate,proto3" json:"client_to_update,omitempty"`
	Height            types.Height `protobuf:"bytes,4,opt,name=height,proto3" json:"height"`
	Timestamp         time.Time    `protobuf:"bytes,5,opt,name=timestamp,proto3,stdtime" json:"timestamp"`
	PacketCommitments [][]byte     `protobuf:"bytes,6,rep,name=packet_commitments,json=packetCommitments,proto3" json:"packet_commitments,omitempty"`
}

func (m *IBCData) Reset()         { *m = IBCData{} }
func (m *IBCData) String() string { return proto.CompactTextString(m) }
func (*IBCData) ProtoMessage()    {}
func (*IBCData) Descriptor() ([]byte, []int) {
	return fileDescriptor_25eb7c0454d2e150, []int{1}
}
func (m *IBCData) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *IBCData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_IBCData.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *IBCData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IBCData.Merge(m, src)
}
func (m *IBCData) XXX_Size() int {
	return m.Size()
}
func (m *IBCData) XXX_DiscardUnknown() {
	xxx_messageInfo_IBCData.DiscardUnknown(m)
}

var xxx_messageInfo_IBCData proto.InternalMessageInfo

func (m *IBCData) GetChainId() string {
	if m != nil {
		return m.ChainId
	}
	return ""
}

func (m *IBCData) GetClientId() string {
	if m != nil {
		return m.ClientId
	}
	return ""
}

func (m *IBCData) GetClientToUpdate() string {
	if m != nil {
		return m.ClientToUpdate
	}
	return ""
}

func (m *IBCData) GetHeight() types.Height {
	if m != nil {
		return m.Height
	}
	return types.Height{}
}

func (m *IBCData) GetTimestamp() time.Time {
	if m != nil {
		return m.Timestamp
	}
	return time.Time{}
}

func (m *IBCData) GetPacketCommitments() [][]byte {
	if m != nil {
		return m.PacketCommitments
	}
	return nil
}

func init() {
	proto.RegisterType((*Attestation)(nil), "core.types.v1.Attestation")
	proto.RegisterType((*IBCData)(nil), "core.types.v1.IBCData")
}

func init() { proto.RegisterFile("core/types/v1/attestation.proto", fileDescriptor_25eb7c0454d2e150) }

var fileDescriptor_25eb7c0454d2e150 = []byte{
	// 420 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x92, 0x41, 0x6f, 0xd3, 0x30,
	0x1c, 0xc5, 0x9b, 0x6e, 0x74, 0xad, 0xd7, 0xa1, 0xcd, 0x42, 0x28, 0x14, 0x29, 0xa9, 0x76, 0x8a,
	0x84, 0x66, 0xab, 0xec, 0xc2, 0xb5, 0x19, 0x07, 0x7a, 0xe0, 0x12, 0x8d, 0x0b, 0x97, 0xc8, 0xb1,
	0x4d, 0x62, 0xd1, 0xc4, 0x51, 0xf2, 0x6f, 0x24, 0x0e, 0x7c, 0x87, 0x1d, 0xf9, 0x48, 0x3b, 0xee,
	0xc8, 0x09, 0x50, 0xfb, 0x45, 0x50, 0xec, 0xa4, 0x81, 0x5b, 0xfc, 0xde, 0xef, 0xef, 0xff, 0xcb,
	0x93, 0x91, 0xcf, 0x75, 0x25, 0x29, 0x7c, 0x2b, 0x65, 0x4d, 0x9b, 0x15, 0x65, 0x00, 0xb2, 0x06,
	0x06, 0x4a, 0x17, 0xa4, 0xac, 0x34, 0x68, 0x7c, 0xd1, 0x02, 0xc4, 0x00, 0xa4, 0x59, 0x2d, 0x5e,
	0xa4, 0x3a, 0xd5, 0xc6, 0xa1, 0xed, 0x97, 0x85, 0x16, 0xbe, 0x4a, 0x38, 0x35, 0x37, 0xf1, 0xad,
	0x92, 0x05, 0xb4, 0x57, 0xd9, 0xaf, 0x1e, 0x48, 0xb5, 0x4e, 0xb7, 0x92, 0x9a, 0x53, 0xb2, 0xfb,
	0x42, 0x41, 0xe5, 0xed, 0xa2, 0xbc, 0xb4, 0xc0, 0xf5, 0x77, 0x74, 0xbe, 0x1e, 0x76, 0xe3, 0x37,
	0xe8, 0xaa, 0x61, 0x5b, 0x25, 0x18, 0xe8, 0x2a, 0x66, 0x42, 0x54, 0xb2, 0xae, 0x5d, 0x67, 0xe9,
	0x04, 0xf3, 0xe8, 0xf2, 0x68, 0xac, 0xad, 0x8e, 0xd7, 0xe8, 0xc2, 0xe6, 0x96, 0x22, 0x16, 0x0c,
	0x98, 0x3b, 0x5e, 0x3a, 0xc1, 0xf9, 0xdb, 0x97, 0xe4, 0xbf, 0xe8, 0x64, 0x13, 0xde, 0xbd, 0x67,
	0xc0, 0xc2, 0xd3, 0xc7, 0x5f, 0xfe, 0x28, 0x9a, 0xf7, 0x23, 0xad, 0x76, 0xfd, 0x63, 0x8c, 0xce,
	0x3a, 0x1f, 0xbf, 0x42, 0x53, 0x9e, 0x31, 0x55, 0xc4, 0x4a, 0x98, 0x95, 0xb3, 0xe8, 0xcc, 0x9c,
	0x37, 0x02, 0xbf, 0x46, 0x33, 0xfb, 0x5b, 0xad, 0x37, 0x36, 0xde, 0xd4, 0x0a, 0x1b, 0x81, 0x03,
	0x74, 0xd9, 0x99, 0xa0, 0xe3, 0x5d, 0x29, 0x18, 0x48, 0xf7, 0xc4, 0x30, 0xcf, 0xad, 0x7e, 0xaf,
	0x3f, 0x19, 0x15, 0xbf, 0x43, 0x93, 0x4c, 0xaa, 0x34, 0x03, 0xf7, 0xd4, 0x24, 0x5d, 0x10, 0x95,
	0x70, 0x9b, 0xb6, 0x6b, 0xad, 0x59, 0x91, 0x0f, 0x86, 0xe8, 0xd2, 0x76, 0x3c, 0x0e, 0xd1, 0xec,
	0xd8, 0x9c, 0xfb, 0xac, 0x1b, 0xb6, 0xdd, 0x92, 0xbe, 0x5b, 0x72, 0xdf, 0x13, 0xe1, 0xb4, 0x1d,
	0x7e, 0xf8, 0xed, 0x3b, 0xd1, 0x30, 0x86, 0x6f, 0x10, 0x2e, 0x19, 0xff, 0x2a, 0x21, 0xe6, 0x3a,
	0xcf, 0x15, 0xe4, 0xb2, 0x80, 0xda, 0x9d, 0x2c, 0x4f, 0x82, 0x79, 0x74, 0x65, 0x9d, 0xbb, 0xc1,
	0x08, 0x3f, 0x3e, 0xee, 0x3d, 0xe7, 0x69, 0xef, 0x39, 0x7f, 0xf6, 0x9e, 0xf3, 0x70, 0xf0, 0x46,
	0x4f, 0x07, 0x6f, 0xf4, 0xf3, 0xe0, 0x8d, 0x3e, 0xdf, 0xa6, 0x0a, 0xb2, 0x5d, 0x42, 0xb8, 0xce,
	0x29, 0xd7, 0x75, 0xae, 0x6b, 0xaa, 0x0a, 0x90, 0x95, 0x69, 0xeb, 0xe6, 0x9f, 0xa7, 0x44, 0x87,
	0x47, 0x96, 0x4c, 0x4c, 0xcc, 0xdb, 0xbf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x58, 0x09, 0x22, 0x47,
	0x79, 0x02, 0x00, 0x00,
}

func (m *Attestation) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Attestation) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Attestation) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.AttestedData.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintAttestation(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.ValidatorAddress) > 0 {
		i -= len(m.ValidatorAddress)
		copy(dAtA[i:], m.ValidatorAddress)
		i = encodeVarintAttestation(dAtA, i, uint64(len(m.ValidatorAddress)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *IBCData) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *IBCData) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *IBCData) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.PacketCommitments) > 0 {
		for iNdEx := len(m.PacketCommitments) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.PacketCommitments[iNdEx])
			copy(dAtA[i:], m.PacketCommitments[iNdEx])
			i = encodeVarintAttestation(dAtA, i, uint64(len(m.PacketCommitments[iNdEx])))
			i--
			dAtA[i] = 0x32
		}
	}
	n2, err2 := github_com_cosmos_gogoproto_types.StdTimeMarshalTo(m.Timestamp, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdTime(m.Timestamp):])
	if err2 != nil {
		return 0, err2
	}
	i -= n2
	i = encodeVarintAttestation(dAtA, i, uint64(n2))
	i--
	dAtA[i] = 0x2a
	{
		size, err := m.Height.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintAttestation(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	if len(m.ClientToUpdate) > 0 {
		i -= len(m.ClientToUpdate)
		copy(dAtA[i:], m.ClientToUpdate)
		i = encodeVarintAttestation(dAtA, i, uint64(len(m.ClientToUpdate)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.ClientId) > 0 {
		i -= len(m.ClientId)
		copy(dAtA[i:], m.ClientId)
		i = encodeVarintAttestation(dAtA, i, uint64(len(m.ClientId)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.ChainId) > 0 {
		i -= len(m.ChainId)
		copy(dAtA[i:], m.ChainId)
		i = encodeVarintAttestation(dAtA, i, uint64(len(m.ChainId)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintAttestation(dAtA []byte, offset int, v uint64) int {
	offset -= sovAttestation(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Attestation) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.ValidatorAddress)
	if l > 0 {
		n += 1 + l + sovAttestation(uint64(l))
	}
	l = m.AttestedData.Size()
	n += 1 + l + sovAttestation(uint64(l))
	return n
}

func (m *IBCData) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.ChainId)
	if l > 0 {
		n += 1 + l + sovAttestation(uint64(l))
	}
	l = len(m.ClientId)
	if l > 0 {
		n += 1 + l + sovAttestation(uint64(l))
	}
	l = len(m.ClientToUpdate)
	if l > 0 {
		n += 1 + l + sovAttestation(uint64(l))
	}
	l = m.Height.Size()
	n += 1 + l + sovAttestation(uint64(l))
	l = github_com_cosmos_gogoproto_types.SizeOfStdTime(m.Timestamp)
	n += 1 + l + sovAttestation(uint64(l))
	if len(m.PacketCommitments) > 0 {
		for _, b := range m.PacketCommitments {
			l = len(b)
			n += 1 + l + sovAttestation(uint64(l))
		}
	}
	return n
}

func sovAttestation(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozAttestation(x uint64) (n int) {
	return sovAttestation(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Attestation) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAttestation
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Attestation: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Attestation: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValidatorAddress", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAttestation
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthAttestation
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthAttestation
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ValidatorAddress = append(m.ValidatorAddress[:0], dAtA[iNdEx:postIndex]...)
			if m.ValidatorAddress == nil {
				m.ValidatorAddress = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AttestedData", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAttestation
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthAttestation
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthAttestation
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.AttestedData.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipAttestation(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthAttestation
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *IBCData) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAttestation
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: IBCData: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: IBCData: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChainId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAttestation
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthAttestation
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAttestation
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ChainId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClientId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAttestation
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthAttestation
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAttestation
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ClientId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClientToUpdate", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAttestation
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthAttestation
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAttestation
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ClientToUpdate = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Height", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAttestation
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthAttestation
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthAttestation
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Height.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAttestation
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthAttestation
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthAttestation
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_cosmos_gogoproto_types.StdTimeUnmarshal(&m.Timestamp, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PacketCommitments", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAttestation
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthAttestation
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthAttestation
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PacketCommitments = append(m.PacketCommitments, make([]byte, postIndex-iNdEx))
			copy(m.PacketCommitments[len(m.PacketCommitments)-1], dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipAttestation(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthAttestation
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipAttestation(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowAttestation
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowAttestation
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowAttestation
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthAttestation
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupAttestation
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthAttestation
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthAttestation        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowAttestation          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupAttestation = fmt.Errorf("proto: unexpected end of group")
)
