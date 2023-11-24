package nic

const (
	INET_ADDRESS_TYPE_IPV4 = uint8(0)
	INET_ADDRESS_TYPE_IPV6 = uint8(1)
)

type IINetAddress interface {
	Type() uint8
	Address() []byte
	NetMask() []byte
	AddressString() string
	NetMaskStr() string
	String() string
}
