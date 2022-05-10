package evm_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	fxtypes "github.com/functionx/fx-core/types"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/rpc/client/http"

	"github.com/functionx/fx-core/app"
	_ "github.com/functionx/fx-core/app"
)

func TestQueryBalance(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	client, err := ethclient.Dial("http://0.0.0.0:8545")
	require.NoError(t, err)

	addressBytes, err := sdk.AccAddressFromBech32("fx17ykqect7ee5e9r4l2end78d8gmp6mauzj87cwz")
	require.NoError(t, err)

	address := common.BytesToAddress(addressBytes)
	println(address.Hex())
	balanceRes, err := client.BalanceAt(context.Background(), address, nil)
	require.NoError(t, err)
	println(balanceRes.String())
}

func TestQueryTransaction(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	client, err := ethclient.Dial("http://0.0.0.0:8545")
	require.NoError(t, err)

	transactionReceipt, err := client.TransactionReceipt(context.Background(), common.HexToHash("0x74a90ed91f42baa375804c22e2fa17087a6060bbca4ffb8f1e0fc1446883a0f7"))
	require.NoError(t, err)
	t.Logf("transactionReceipt:%+#v", transactionReceipt)
}

func TestQueryTransactionRaw(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	client, err := ethclient.Dial("http://0.0.0.0:8545")
	require.NoError(t, err)

	tx, _, err := client.TransactionByHash(context.Background(), common.HexToHash("0x815fced350c7d84ab36e1aa2ff392d55c4b810876f8b6bbf01312b5148f8f543"))
	require.NoError(t, err)

	bz, err := tx.MarshalBinary()
	require.NoError(t, err)

	t.Logf("raw hex %x", bz)
}

func TestQueryFxTxByEvmHash(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	client, err := ethclient.Dial("http://0.0.0.0:8545")
	require.NoError(t, err)

	transactionReceipt, err := client.TransactionReceipt(context.Background(), common.HexToHash("0x74a90ed91f42baa375804c22e2fa17087a6060bbca4ffb8f1e0fc1446883a0f7"))
	require.NoError(t, err)
	t.Logf("transactionReceipt:%+#v", transactionReceipt)

	fxClient, err := http.New("http://0.0.0.0:26657", "/websocket")
	require.NoError(t, err)
	evmHashBlockNumber := transactionReceipt.BlockNumber.Int64()
	blockData, err := fxClient.Block(context.Background(), &evmHashBlockNumber)
	require.NoError(t, err)
	require.True(t, uint(len(blockData.Block.Txs)) > transactionReceipt.TransactionIndex)
	fxTx := blockData.Block.Txs[transactionReceipt.TransactionIndex]
	encodingConfig := app.MakeEncodingConfig()
	tx, err := encodingConfig.TxConfig.TxDecoder()(fxTx)
	require.NoError(t, err)
	txJsonStr, err := encodingConfig.TxConfig.TxJSONEncoder()(tx)
	require.NoError(t, err)
	//marshalIndent, err := json.MarshalIndent(string(txJsonStr), "", "  ")
	//require.NoError(t, err)
	t.Logf("\nTxHash:%x\nData:\n%v", fxTx.Hash(), string(txJsonStr))

}
func TestMnemonicToFxPrivate(t *testing.T) {
	privKey, err := mnemonicToFxPrivKey("december slow blue fury silly bread friend unknown render resource dry buyer brand final abstract gallery slow since hood shadow neglect travel convince foil")
	require.NoError(t, err)
	t.Logf("%x", privKey.Key)
}

func TestEthPrivateKeyToAddress(t *testing.T) {
	//privateKey, err := crypto.GenerateKey()
	//require.NoError(t, err)
	//fromECDSA := crypto.FromECDSA(privateKey)
	//t.Logf("fromEc:%x", fromECDSA)

	// 1ce31354ff0a3f057c9b70ebbbdacb68ace4bf9c008ac722f2b996328ab3ca08
	hexPrivKey := "86b87f127b6e0901f7f00aa77b6c82624847f2628a901bf1833b2d48883b73d3"
	recoverPrivKey, err := crypto.HexToECDSA(hexPrivKey)
	require.NoError(t, err)
	address := crypto.PubkeyToAddress(recoverPrivKey.PublicKey)
	t.Logf("Eth address:%v, FxAddress:%v", address.Hex(), sdk.AccAddress(address.Bytes()).String())
}

func TestEthAddressToFxAddress(t *testing.T) {
	ethAddress := common.HexToAddress("0xf12C0Ce17eCE69928ebf5666Df1Da746c3adf782")
	t.Logf("%o", ethAddress.Bytes())
	t.Logf("EthAddress:%v, FxAddress:%v", ethAddress.Hex(), sdk.AccAddress(ethAddress.Bytes()).String())
}

func TestFxAddressToEthAddress(t *testing.T) {
	fxAddress, err := sdk.AccAddressFromBech32("fx10kg059hhxc2pevxssszunvgc70jpmxsjal4xf6")
	require.NoError(t, err)
	ethAddress := common.BytesToAddress(fxAddress)
	t.Logf("EthAddress:%v, FxAddress:%v", ethAddress.Hex(), sdk.AccAddress(ethAddress.Bytes()).String())
}

func TestTraverseBlockERC20(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	fxClient, err := http.New("http://127.0.0.1:26657", "/websocket")
	require.NoError(t, err)

	ctx := context.Background()
	info, err := fxClient.Status(ctx)
	if err != nil {
		t.Fatal(err)
	}
	for i := int64(1); i < info.SyncInfo.LatestBlockHeight; i++ {
		block, err := fxClient.BlockResults(ctx, &i)
		if err != nil {
			t.Fatal(err)
		}
		for _, v := range block.EndBlockEvents {
			for _, vv := range v.Attributes {
				if strings.EqualFold("fip20_symbol", string(vv.Key)) {
					fmt.Println(i, "fip20 symbol:", string(vv.Value))
				}
				if strings.EqualFold("fip20_token", string(vv.Key)) {
					fmt.Println(i, "fip20 address:", string(vv.Value))
				}
			}
		}
	}
}

func mnemonicToFxPrivKey(mnemonic string) (*secp256k1.PrivKey, error) {
	algo := hd.Secp256k1
	bytes, err := algo.Derive()(mnemonic, "", "m/44'/118'/0'/0/0")
	if err != nil {
		return nil, err
	}
	privKey := algo.Generate()(bytes)
	priv, ok := privKey.(*secp256k1.PrivKey)
	if !ok {
		return nil, fmt.Errorf("not secp256k1.PrivKey")
	}
	return priv, nil
}

func TestFIP20Code(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	codeAddr := "0x5f123738067a8BAA3E9bb8Cd7e4A8827474a2F53"
	codeBase64 := "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"

	code, codeNew := replayCodeAddress(t, codeBase64, codeAddr, fxtypes.FIP20LogicAddress)
	t.Log("code", hex.EncodeToString(code))
	t.Log("new-code", hex.EncodeToString(codeNew))
}

func TestWFXCode(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	codeAddr := "0x5882566ad042E16F14633a0D77b705E9a912e94d"
	codeBase64 := "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"

	code, codeNew := replayCodeAddress(t, codeBase64, codeAddr, fxtypes.WFXLogicAddress)
	t.Log("code", hex.EncodeToString(code))
	t.Log("new-code", hex.EncodeToString(codeNew))
}

func replayCodeAddress(t *testing.T, codeBase64, addr, addrNew string) (code, codeNew []byte) {
	bz, err := base64.StdEncoding.DecodeString(codeBase64)
	require.NoError(t, err)

	addr1 := common.HexToAddress(addr)
	addr2 := common.HexToAddress(addrNew)

	bzZero := bytes.ReplaceAll(bz, addr1.Bytes(), common.HexToAddress(fxtypes.EmptyEvmAddress).Bytes())
	bzNew := bytes.ReplaceAll(bz, addr1.Bytes(), addr2.Bytes())

	return bzZero, bzNew
}
