// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: fx/gov/v1/params.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"
	_ "google.golang.org/protobuf/types/known/durationpb"
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

// Params defines the fx x/gov module params
type Params struct {
	MinInitialDeposit   types.Coin     `protobuf:"bytes,1,opt,name=min_initial_deposit,json=minInitialDeposit,proto3" json:"min_initial_deposit"`
	EgfDepositThreshold types.Coin     `protobuf:"bytes,2,opt,name=egf_deposit_threshold,json=egfDepositThreshold,proto3" json:"egf_deposit_threshold"`
	ClaimRatio          string         `protobuf:"bytes,3,opt,name=claim_ratio,json=claimRatio,proto3" json:"claim_ratio,omitempty"`
	Erc20Quorum         string         `protobuf:"bytes,4,opt,name=erc20_quorum,json=erc20Quorum,proto3" json:"erc20_quorum,omitempty"`
	EvmQuorum           string         `protobuf:"bytes,5,opt,name=evm_quorum,json=evmQuorum,proto3" json:"evm_quorum,omitempty"`
	EgfVotingPeriod     *time.Duration `protobuf:"bytes,6,opt,name=egf_voting_period,json=egfVotingPeriod,proto3,stdduration" json:"egf_voting_period,omitempty"`
	EvmVotingPeriod     *time.Duration `protobuf:"bytes,7,opt,name=evm_voting_period,json=evmVotingPeriod,proto3,stdduration" json:"evm_voting_period,omitempty"`
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_a8e5d06ed1291671, []int{0}
}
func (m *Params) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params.Merge(m, src)
}
func (m *Params) XXX_Size() int {
	return m.Size()
}
func (m *Params) XXX_DiscardUnknown() {
	xxx_messageInfo_Params.DiscardUnknown(m)
}

var xxx_messageInfo_Params proto.InternalMessageInfo

func (m *Params) GetMinInitialDeposit() types.Coin {
	if m != nil {
		return m.MinInitialDeposit
	}
	return types.Coin{}
}

func (m *Params) GetEgfDepositThreshold() types.Coin {
	if m != nil {
		return m.EgfDepositThreshold
	}
	return types.Coin{}
}

func (m *Params) GetClaimRatio() string {
	if m != nil {
		return m.ClaimRatio
	}
	return ""
}

func (m *Params) GetErc20Quorum() string {
	if m != nil {
		return m.Erc20Quorum
	}
	return ""
}

func (m *Params) GetEvmQuorum() string {
	if m != nil {
		return m.EvmQuorum
	}
	return ""
}

func (m *Params) GetEgfVotingPeriod() *time.Duration {
	if m != nil {
		return m.EgfVotingPeriod
	}
	return nil
}

func (m *Params) GetEvmVotingPeriod() *time.Duration {
	if m != nil {
		return m.EvmVotingPeriod
	}
	return nil
}

func init() {
	proto.RegisterType((*Params)(nil), "fx.gov.v1.Params")
}

func init() { proto.RegisterFile("fx/gov/v1/params.proto", fileDescriptor_a8e5d06ed1291671) }

var fileDescriptor_a8e5d06ed1291671 = []byte{
	// 424 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0x3f, 0x8f, 0xd3, 0x30,
	0x18, 0xc6, 0x1b, 0xae, 0x14, 0xd5, 0x45, 0xa0, 0xcb, 0x01, 0xca, 0xdd, 0x90, 0x9e, 0x98, 0xba,
	0xd4, 0x26, 0xc7, 0x37, 0x08, 0x5d, 0x10, 0x03, 0x47, 0x40, 0x0c, 0x2c, 0x51, 0x92, 0xda, 0xae,
	0xa5, 0x3a, 0x6f, 0x70, 0x1c, 0xab, 0x7c, 0x0b, 0x46, 0x3e, 0x08, 0x1f, 0xe2, 0x24, 0x96, 0x13,
	0x13, 0x13, 0xa0, 0xf6, 0x8b, 0xa0, 0xd8, 0x8e, 0x84, 0x10, 0x43, 0xb7, 0xbc, 0x8f, 0x9f, 0xe7,
	0xf7, 0xfe, 0x51, 0xd0, 0x13, 0xb6, 0x23, 0x1c, 0x0c, 0x31, 0x09, 0x69, 0x0a, 0x55, 0xc8, 0x16,
	0x37, 0x0a, 0x34, 0x84, 0x53, 0xb6, 0xc3, 0x1c, 0x0c, 0x36, 0xc9, 0x45, 0x5c, 0x41, 0x2b, 0xa1,
	0x25, 0x65, 0xd1, 0x52, 0x62, 0x92, 0x92, 0xea, 0x22, 0x21, 0x15, 0x88, 0xda, 0x59, 0x2f, 0xce,
	0xdd, 0x7b, 0x6e, 0x2b, 0xe2, 0x0a, 0xff, 0xf4, 0x88, 0x03, 0x07, 0xa7, 0xf7, 0x5f, 0x5e, 0x8d,
	0x39, 0x00, 0xdf, 0x52, 0x62, 0xab, 0xb2, 0x63, 0x64, 0xdd, 0xa9, 0x42, 0x0b, 0xf0, 0xc0, 0xa7,
	0xdf, 0x4e, 0xd0, 0xe4, 0xda, 0x0e, 0x13, 0xbe, 0x46, 0x67, 0x52, 0xd4, 0xb9, 0xa8, 0x85, 0x16,
	0xc5, 0x36, 0x5f, 0xd3, 0x06, 0x5a, 0xa1, 0xa3, 0xe0, 0x32, 0x58, 0xcc, 0xae, 0xce, 0xb1, 0x6f,
	0xd6, 0x4f, 0x86, 0xfd, 0x64, 0xf8, 0x05, 0x88, 0x3a, 0x1d, 0xdf, 0xfc, 0x9c, 0x8f, 0xb2, 0x53,
	0x29, 0xea, 0x97, 0x2e, 0xba, 0x72, 0xc9, 0xf0, 0x2d, 0x7a, 0x4c, 0x39, 0x1b, 0x40, 0xb9, 0xde,
	0x28, 0xda, 0x6e, 0x60, 0xbb, 0x8e, 0xee, 0x1c, 0x87, 0x3c, 0xa3, 0x9c, 0x79, 0xd6, 0xbb, 0x21,
	0x1b, 0xce, 0xd1, 0xac, 0xda, 0x16, 0x42, 0xe6, 0x76, 0x8d, 0xe8, 0xe4, 0x32, 0x58, 0x4c, 0x33,
	0x64, 0xa5, 0xac, 0x57, 0xc2, 0x04, 0xdd, 0xa7, 0xaa, 0xba, 0x7a, 0x96, 0x7f, 0xec, 0x40, 0x75,
	0x32, 0x1a, 0xf7, 0x8e, 0xf4, 0xc1, 0xf7, 0xaf, 0x4b, 0xe4, 0xfb, 0xad, 0x68, 0x95, 0xcd, 0xac,
	0xe7, 0x8d, 0xb5, 0x84, 0x4b, 0x84, 0xa8, 0x91, 0x43, 0xe0, 0xee, 0x7f, 0x03, 0x53, 0x6a, 0xa4,
	0xb7, 0xbf, 0x42, 0xa7, 0xfd, 0x5e, 0x06, 0xb4, 0xa8, 0x79, 0xde, 0x50, 0x25, 0x60, 0x1d, 0x4d,
	0xfc, 0x4e, 0xee, 0xde, 0x78, 0xb8, 0x37, 0x5e, 0xf9, 0x7b, 0xa7, 0xe3, 0x2f, 0xbf, 0xe6, 0x41,
	0xf6, 0x90, 0x72, 0xf6, 0xde, 0x06, 0xaf, 0x6d, 0xce, 0xc2, 0x8c, 0xfc, 0x07, 0x76, 0xef, 0x58,
	0x98, 0x91, 0x7f, 0xc3, 0xd2, 0xf4, 0x66, 0x1f, 0x07, 0xb7, 0xfb, 0x38, 0xf8, 0xbd, 0x8f, 0x83,
	0xcf, 0x87, 0x78, 0x74, 0x7b, 0x88, 0x47, 0x3f, 0x0e, 0xf1, 0xe8, 0xc3, 0x82, 0x0b, 0xbd, 0xe9,
	0x4a, 0x5c, 0x81, 0x24, 0xac, 0xab, 0xab, 0x9e, 0xb2, 0x23, 0x6c, 0xb7, 0xac, 0x40, 0x51, 0xe2,
	0xfe, 0x4b, 0xfd, 0xa9, 0xa1, 0x6d, 0x39, 0xb1, 0xdd, 0x9e, 0xff, 0x09, 0x00, 0x00, 0xff, 0xff,
	0xec, 0xc5, 0x1d, 0x1d, 0xae, 0x02, 0x00, 0x00,
}

func (m *Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Params) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.EvmVotingPeriod != nil {
		n1, err1 := github_com_gogo_protobuf_types.StdDurationMarshalTo(*m.EvmVotingPeriod, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdDuration(*m.EvmVotingPeriod):])
		if err1 != nil {
			return 0, err1
		}
		i -= n1
		i = encodeVarintParams(dAtA, i, uint64(n1))
		i--
		dAtA[i] = 0x3a
	}
	if m.EgfVotingPeriod != nil {
		n2, err2 := github_com_gogo_protobuf_types.StdDurationMarshalTo(*m.EgfVotingPeriod, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdDuration(*m.EgfVotingPeriod):])
		if err2 != nil {
			return 0, err2
		}
		i -= n2
		i = encodeVarintParams(dAtA, i, uint64(n2))
		i--
		dAtA[i] = 0x32
	}
	if len(m.EvmQuorum) > 0 {
		i -= len(m.EvmQuorum)
		copy(dAtA[i:], m.EvmQuorum)
		i = encodeVarintParams(dAtA, i, uint64(len(m.EvmQuorum)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.Erc20Quorum) > 0 {
		i -= len(m.Erc20Quorum)
		copy(dAtA[i:], m.Erc20Quorum)
		i = encodeVarintParams(dAtA, i, uint64(len(m.Erc20Quorum)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.ClaimRatio) > 0 {
		i -= len(m.ClaimRatio)
		copy(dAtA[i:], m.ClaimRatio)
		i = encodeVarintParams(dAtA, i, uint64(len(m.ClaimRatio)))
		i--
		dAtA[i] = 0x1a
	}
	{
		size, err := m.EgfDepositThreshold.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	{
		size, err := m.MinInitialDeposit.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintParams(dAtA []byte, offset int, v uint64) int {
	offset -= sovParams(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.MinInitialDeposit.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.EgfDepositThreshold.Size()
	n += 1 + l + sovParams(uint64(l))
	l = len(m.ClaimRatio)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	l = len(m.Erc20Quorum)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	l = len(m.EvmQuorum)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	if m.EgfVotingPeriod != nil {
		l = github_com_gogo_protobuf_types.SizeOfStdDuration(*m.EgfVotingPeriod)
		n += 1 + l + sovParams(uint64(l))
	}
	if m.EvmVotingPeriod != nil {
		l = github_com_gogo_protobuf_types.SizeOfStdDuration(*m.EvmVotingPeriod)
		n += 1 + l + sovParams(uint64(l))
	}
	return n
}

func sovParams(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozParams(x uint64) (n int) {
	return sovParams(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowParams
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
			return fmt.Errorf("proto: Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinInitialDeposit", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.MinInitialDeposit.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EgfDepositThreshold", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.EgfDepositThreshold.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClaimRatio", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ClaimRatio = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Erc20Quorum", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Erc20Quorum = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EvmQuorum", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.EvmQuorum = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EgfVotingPeriod", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.EgfVotingPeriod == nil {
				m.EgfVotingPeriod = new(time.Duration)
			}
			if err := github_com_gogo_protobuf_types.StdDurationUnmarshal(m.EgfVotingPeriod, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EvmVotingPeriod", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.EvmVotingPeriod == nil {
				m.EvmVotingPeriod = new(time.Duration)
			}
			if err := github_com_gogo_protobuf_types.StdDurationUnmarshal(m.EvmVotingPeriod, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipParams(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthParams
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
func skipParams(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowParams
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
					return 0, ErrIntOverflowParams
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
					return 0, ErrIntOverflowParams
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
				return 0, ErrInvalidLengthParams
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupParams
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthParams
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthParams        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowParams          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupParams = fmt.Errorf("proto: unexpected end of group")
)