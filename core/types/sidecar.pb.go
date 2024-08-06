// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: core/sidecar/v1/sidecar.proto

package types

import (
	context "context"
	fmt "fmt"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	proto "github.com/cosmos/gogoproto/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type AttestationRequest struct {
	ChainId string `protobuf:"bytes,1,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
}

func (m *AttestationRequest) Reset()         { *m = AttestationRequest{} }
func (m *AttestationRequest) String() string { return proto.CompactTextString(m) }
func (*AttestationRequest) ProtoMessage()    {}
func (*AttestationRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8ce634b51eec8241, []int{0}
}
func (m *AttestationRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AttestationRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AttestationRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AttestationRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AttestationRequest.Merge(m, src)
}
func (m *AttestationRequest) XXX_Size() int {
	return m.Size()
}
func (m *AttestationRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AttestationRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AttestationRequest proto.InternalMessageInfo

func (m *AttestationRequest) GetChainId() string {
	if m != nil {
		return m.ChainId
	}
	return ""
}

type AttestationResponse struct {
	Attestation *Attestation `protobuf:"bytes,1,opt,name=attestation,proto3" json:"attestation,omitempty"`
}

func (m *AttestationResponse) Reset()         { *m = AttestationResponse{} }
func (m *AttestationResponse) String() string { return proto.CompactTextString(m) }
func (*AttestationResponse) ProtoMessage()    {}
func (*AttestationResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8ce634b51eec8241, []int{1}
}
func (m *AttestationResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AttestationResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AttestationResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AttestationResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AttestationResponse.Merge(m, src)
}
func (m *AttestationResponse) XXX_Size() int {
	return m.Size()
}
func (m *AttestationResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AttestationResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AttestationResponse proto.InternalMessageInfo

func (m *AttestationResponse) GetAttestation() *Attestation {
	if m != nil {
		return m.Attestation
	}
	return nil
}

func init() {
	proto.RegisterType((*AttestationRequest)(nil), "core.sidecar.v1.AttestationRequest")
	proto.RegisterType((*AttestationResponse)(nil), "core.sidecar.v1.AttestationResponse")
}

func init() { proto.RegisterFile("core/sidecar/v1/sidecar.proto", fileDescriptor_8ce634b51eec8241) }

var fileDescriptor_8ce634b51eec8241 = []byte{
	// 273 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x4d, 0xce, 0x2f, 0x4a,
	0xd5, 0x2f, 0xce, 0x4c, 0x49, 0x4d, 0x4e, 0x2c, 0xd2, 0x2f, 0x33, 0x84, 0x31, 0xf5, 0x0a, 0x8a,
	0xf2, 0x4b, 0xf2, 0x85, 0xf8, 0x41, 0xd2, 0x7a, 0x30, 0xb1, 0x32, 0x43, 0x29, 0x79, 0xb0, 0xfa,
	0x92, 0xca, 0x82, 0xd4, 0x62, 0x90, 0xea, 0xc4, 0x92, 0x92, 0xd4, 0xe2, 0x92, 0xc4, 0x92, 0xcc,
	0xfc, 0x3c, 0x88, 0x0e, 0x25, 0x7d, 0x2e, 0x21, 0x47, 0x84, 0x60, 0x50, 0x6a, 0x61, 0x69, 0x6a,
	0x71, 0x89, 0x90, 0x24, 0x17, 0x47, 0x72, 0x46, 0x62, 0x66, 0x5e, 0x7c, 0x66, 0x8a, 0x04, 0xa3,
	0x02, 0xa3, 0x06, 0x67, 0x10, 0x3b, 0x98, 0xef, 0x99, 0xa2, 0x14, 0xcc, 0x25, 0x8c, 0xa2, 0xa1,
	0xb8, 0x20, 0x3f, 0xaf, 0x38, 0x55, 0xc8, 0x86, 0x8b, 0x1b, 0xc9, 0x70, 0xb0, 0x26, 0x6e, 0x23,
	0x29, 0x3d, 0xb0, 0x7b, 0xc0, 0xd6, 0xeb, 0x95, 0x19, 0xea, 0x21, 0x6b, 0x44, 0x56, 0x6e, 0x94,
	0xc1, 0xc5, 0x1e, 0x0c, 0x71, 0xb4, 0x50, 0x2c, 0x17, 0x9f, 0x7b, 0x6a, 0x09, 0x92, 0x4a, 0x21,
	0x65, 0x3d, 0x34, 0x5f, 0xe9, 0x61, 0xba, 0x58, 0x4a, 0x05, 0xbf, 0x22, 0x88, 0x2b, 0x95, 0x18,
	0x9c, 0x42, 0x4f, 0x3c, 0x92, 0x63, 0xbc, 0xf0, 0x48, 0x8e, 0xf1, 0xc1, 0x23, 0x39, 0xc6, 0x09,
	0x8f, 0xe5, 0x18, 0x2e, 0x3c, 0x96, 0x63, 0xb8, 0xf1, 0x58, 0x8e, 0x21, 0xca, 0x3a, 0x3d, 0xb3,
	0x24, 0xa3, 0x34, 0x49, 0x2f, 0x39, 0x3f, 0x57, 0x3f, 0x3d, 0x2b, 0xb5, 0x28, 0xb7, 0x34, 0x2f,
	0x25, 0x3d, 0xb1, 0x28, 0x31, 0x29, 0x51, 0xbf, 0x20, 0xb5, 0xb8, 0x38, 0x33, 0x37, 0xb3, 0xb8,
	0x24, 0x33, 0x59, 0xb7, 0x2c, 0x31, 0x27, 0x33, 0x05, 0x6c, 0xaa, 0x3e, 0x22, 0x6c, 0x93, 0xd8,
	0xc0, 0xa1, 0x69, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0xc2, 0x85, 0x4a, 0xda, 0xa0, 0x01, 0x00,
	0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// SidecarClient is the client API for Sidecar service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SidecarClient interface {
	GetAttestation(ctx context.Context, in *AttestationRequest, opts ...grpc.CallOption) (*AttestationResponse, error)
}

type sidecarClient struct {
	cc grpc1.ClientConn
}

func NewSidecarClient(cc grpc1.ClientConn) SidecarClient {
	return &sidecarClient{cc}
}

func (c *sidecarClient) GetAttestation(ctx context.Context, in *AttestationRequest, opts ...grpc.CallOption) (*AttestationResponse, error) {
	out := new(AttestationResponse)
	err := c.cc.Invoke(ctx, "/core.sidecar.v1.Sidecar/GetAttestation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SidecarServer is the server API for Sidecar service.
type SidecarServer interface {
	GetAttestation(context.Context, *AttestationRequest) (*AttestationResponse, error)
}

// UnimplementedSidecarServer can be embedded to have forward compatible implementations.
type UnimplementedSidecarServer struct {
}

func (*UnimplementedSidecarServer) GetAttestation(ctx context.Context, req *AttestationRequest) (*AttestationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAttestation not implemented")
}

func RegisterSidecarServer(s grpc1.Server, srv SidecarServer) {
	s.RegisterService(&_Sidecar_serviceDesc, srv)
}

func _Sidecar_GetAttestation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AttestationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SidecarServer).GetAttestation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.sidecar.v1.Sidecar/GetAttestation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SidecarServer).GetAttestation(ctx, req.(*AttestationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Sidecar_serviceDesc = grpc.ServiceDesc{
	ServiceName: "core.sidecar.v1.Sidecar",
	HandlerType: (*SidecarServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAttestation",
			Handler:    _Sidecar_GetAttestation_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "core/sidecar/v1/sidecar.proto",
}

func (m *AttestationRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AttestationRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AttestationRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ChainId) > 0 {
		i -= len(m.ChainId)
		copy(dAtA[i:], m.ChainId)
		i = encodeVarintSidecar(dAtA, i, uint64(len(m.ChainId)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *AttestationResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AttestationResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AttestationResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Attestation != nil {
		{
			size, err := m.Attestation.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSidecar(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintSidecar(dAtA []byte, offset int, v uint64) int {
	offset -= sovSidecar(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *AttestationRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.ChainId)
	if l > 0 {
		n += 1 + l + sovSidecar(uint64(l))
	}
	return n
}

func (m *AttestationResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Attestation != nil {
		l = m.Attestation.Size()
		n += 1 + l + sovSidecar(uint64(l))
	}
	return n
}

func sovSidecar(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozSidecar(x uint64) (n int) {
	return sovSidecar(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *AttestationRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSidecar
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
			return fmt.Errorf("proto: AttestationRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AttestationRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChainId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSidecar
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
				return ErrInvalidLengthSidecar
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSidecar
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ChainId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSidecar(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSidecar
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
func (m *AttestationResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSidecar
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
			return fmt.Errorf("proto: AttestationResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AttestationResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Attestation", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSidecar
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
				return ErrInvalidLengthSidecar
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSidecar
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Attestation == nil {
				m.Attestation = &Attestation{}
			}
			if err := m.Attestation.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSidecar(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSidecar
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
func skipSidecar(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowSidecar
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
					return 0, ErrIntOverflowSidecar
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
					return 0, ErrIntOverflowSidecar
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
				return 0, ErrInvalidLengthSidecar
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupSidecar
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthSidecar
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthSidecar        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowSidecar          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupSidecar = fmt.Errorf("proto: unexpected end of group")
)