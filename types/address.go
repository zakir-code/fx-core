package types

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/ethereum/go-ethereum/common"
	tronaddress "github.com/fbsobreira/gotron-sdk/pkg/address"
)

func ParseAddress(addr string) (accAddr sdk.AccAddress, isEvmAddr bool, err error) {
	_, bytes, decodeErr := bech32.DecodeAndConvert(addr)
	if decodeErr == nil {
		return bytes, false, nil
	}
	ethAddrError := ValidateEthereumAddress(addr)
	if ethAddrError == nil {
		return common.HexToAddress(addr).Bytes(), true, nil
	}
	return nil, false, errors.Join(decodeErr, ethAddrError)
}

func AddressToStr(bt []byte, module string) string {
	if module == "tron" {
		return tronaddress.Address(bt).String()
	} else {
		return common.BytesToAddress(bt).String()
	}
}
