package tests

import (
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1betal "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v3/app"
	fxtypes "github.com/functionx/fx-core/v3/types"
	crosschaintypes "github.com/functionx/fx-core/v3/x/crosschain/types"
	erc20types "github.com/functionx/fx-core/v3/x/erc20/types"
	ibctransfertypes "github.com/functionx/fx-core/v3/x/ibc/applications/transfer/types"
)

// nolint
func TestAminoEncode(t *testing.T) {
	oneDec := sdk.NewDec(1)
	oneInt := sdk.NewInt(1)

	testcases := []struct {
		name     string
		expected string
		msg      interface{}
	}{
		{
			name:     "upgrade-SoftwareUpgradeProposal",
			expected: `{"type":"cosmos-sdk/SoftwareUpgradeProposal","value":{"description":"foo","plan":{"height":"123","info":"foo","name":"foo","time":"0001-01-01T00:00:00Z"},"title":"v2"}}`,
			msg: upgradetypes.SoftwareUpgradeProposal{
				Title:       "v2",
				Description: "foo",
				Plan: upgradetypes.Plan{
					Name:   "foo",
					Time:   time.Time{},
					Height: 123,
					Info:   "foo",
				},
			},
		},
		{
			name:     "ibc-MsgTransfer",
			expected: `{"type":"fxtransfer/MsgTransfer","value":{"fee":{"amount":"0","denom":"FX"},"receiver":"0x001","sender":"fx1001","source_channel":"channel-0","source_port":"transfer","timeout_height":{},"timeout_timestamp":"1675063442000000000","token":{"amount":"1","denom":"FX"}}}`,
			msg: ibctransfertypes.MsgTransfer{
				SourcePort:       "transfer",
				SourceChannel:    "channel-0",
				Token:            sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1)),
				Sender:           "fx1001",
				Receiver:         "0x001",
				TimeoutHeight:    clienttypes.Height{},
				TimeoutTimestamp: 1675063442000000000,
				Router:           "",
				Fee:              sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.ZeroInt()),
				Memo:             "",
			},
		},
		{
			name:     "erc20-RegisterCoinProposal",
			expected: `{"type":"erc20/RegisterCoinProposal","value":{"description":"foo","metadata":{"base":"test","denom_units":[{"aliases":["ethtest"],"denom":"test"},{"denom":"TEST","exponent":18}],"description":"test","display":"test","name":"test name","symbol":"TEST"},"title":"v2"}}`,
			msg: erc20types.RegisterCoinProposal{
				Title:       "v2",
				Description: "foo",
				Metadata: types.Metadata{
					Description: "test",
					DenomUnits: []*types.DenomUnit{
						{
							Denom:    "test",
							Exponent: 0,
							Aliases: []string{
								"ethtest",
							},
						},
						{
							Denom:    "TEST",
							Exponent: 18,
							Aliases:  []string{},
						},
					},
					Base:    "test",
					Display: "test",
					Name:    "test name",
					Symbol:  "TEST",
				},
			},
		},
		{
			name:     "staking-MsgEditValidator",
			expected: `{"type":"cosmos-sdk/MsgEditValidator","value":{"commission_rate":"1.000000000000000000","description":{"details":"foo","identity":"foo","moniker":"foo","security_contact":"foo","website":"foo"},"min_self_delegation":"1","validator_address":"cosmosvaloper1h6lrm4uusd46tu4slkg620hylv46hhff7a8su6"}}`,
			msg: stakingtypes.MsgEditValidator{
				Description: stakingtypes.Description{
					Moniker:         "foo",
					Identity:        "foo",
					Website:         "foo",
					SecurityContact: "foo",
					Details:         "foo",
				},
				ValidatorAddress:  "cosmosvaloper1h6lrm4uusd46tu4slkg620hylv46hhff7a8su6",
				CommissionRate:    &oneDec,
				MinSelfDelegation: &oneInt,
			},
		},
		{
			name:     "staking-MsgEditValidator",
			expected: `{"type":"cosmos-sdk/MsgEditValidator","value":{"description":{"details":"foo","moniker":"foo","security_contact":"foo","website":"foo"},"validator_address":"cosmosvaloper1h6lrm4uusd46tu4slkg620hylv46hhff7a8su6"}}`,
			msg: stakingtypes.MsgEditValidator{
				Description: stakingtypes.Description{
					Moniker:         "foo",
					Identity:        "",
					Website:         "foo",
					SecurityContact: "foo",
					Details:         "foo",
				},
				ValidatorAddress:  "cosmosvaloper1h6lrm4uusd46tu4slkg620hylv46hhff7a8su6",
				CommissionRate:    nil,
				MinSelfDelegation: nil,
			},
		},
		{
			name:     "gov-TextProposal",
			expected: `{"type":"cosmos-sdk/TextProposal","value":{"description":"foo desc","title":"foo title"}}`,
			msg: govv1betal.TextProposal{
				Title:       "foo title",
				Description: "foo desc",
			},
		},
		{
			name:     "gov-v1-MsgSubmitProposal-crosschain-MsgUpdateParams",
			expected: `{"type":"cosmos-sdk/v1/MsgSubmitProposal","value":{"initial_deposit":null,"messages":[{"type":"crosschain/MsgUpdateParams","value":{"authority":"1","chain_name":"1","params":{"average_block_time":"1","average_external_block_time":"1","delegate_multiple":"1","delegate_threshold":{"amount":"1","denom":"FX"},"external_batch_timeout":"1","gravity_id":"1","ibc_transfer_timeout_height":"1","oracle_set_update_power_change_percent":"1.000000000000000000","signed_window":"1","slash_fraction":"1.000000000000000000"}}}]}}`,
			msg: govv1.MsgSubmitProposal{
				Messages: []*codectypes.Any{
					NewAnyWithValue(&crosschaintypes.MsgUpdateParams{
						ChainName: "1",
						Authority: "1",
						Params: crosschaintypes.Params{
							GravityId:                         "1",
							AverageBlockTime:                  1,
							ExternalBatchTimeout:              1,
							AverageExternalBlockTime:          1,
							SignedWindow:                      1,
							SlashFraction:                     sdk.NewDec(1),
							OracleSetUpdatePowerChangePercent: sdk.NewDec(1),
							IbcTransferTimeoutHeight:          1,
							Oracles:                           nil,
							DelegateThreshold:                 sdk.NewCoin("FX", sdk.NewInt(1)),
							DelegateMultiple:                  1,
						},
					}),
				},
				InitialDeposit: nil,
				Proposer:       "",
				Metadata:       "",
			},
		},
		{
			name:     "gov-v1-MsgSubmitProposal-erc20-MsgUpdateParams",
			expected: `{"type":"cosmos-sdk/v1/MsgSubmitProposal","value":{"initial_deposit":null,"messages":[{"type":"erc20/MsgUpdateParams","value":{"authority":"1","params":{"enable_erc20":true,"enable_evm_hook":true,"ibc_timeout":"1"}}}]}}`,
			msg: govv1.MsgSubmitProposal{
				Messages: []*codectypes.Any{
					NewAnyWithValue(&erc20types.MsgUpdateParams{
						Authority: "1",
						Params: erc20types.Params{
							EnableErc20:   true,
							EnableEVMHook: true,
							IbcTimeout:    1,
						},
					}),
				},
				InitialDeposit: nil,
				Proposer:       "",
				Metadata:       "",
			},
		},
	}

	encode := app.MakeEncodingConfig()
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			aminoJson, err := encode.Amino.MarshalJSON(testcase.msg)
			require.NoError(t, err)
			require.Equal(t, testcase.expected, string(sdk.MustSortJSON(aminoJson)))
		})
	}
}

func NewAnyWithValue(msg proto.Message) *codectypes.Any {
	value, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		panic(err)
	}
	return value
}
