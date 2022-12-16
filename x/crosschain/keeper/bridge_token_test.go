package keeper_test

import (
	"fmt"

	"github.com/functionx/fx-core/v3/app/helpers"
	"github.com/functionx/fx-core/v3/x/crosschain/types"
)

func (suite *KeeperTestSuite) TestKeeper_BridgeToken() {
	tokenContract := helpers.GenerateAddress().Hex()
	denom, err := suite.Keeper().SetIbcDenomTrace(suite.ctx, tokenContract, "")
	suite.NoError(err)
	suite.Equal(fmt.Sprintf("%s%s", suite.chainName, tokenContract), denom)

	suite.Keeper().AddBridgeToken(suite.ctx, tokenContract, denom)

	bridgeToken := &types.BridgeToken{Token: tokenContract, Denom: denom}
	suite.EqualValues(bridgeToken, suite.Keeper().GetBridgeTokenDenom(suite.ctx, tokenContract))

	suite.EqualValues(bridgeToken, suite.Keeper().GetDenomByBridgeToken(suite.ctx, denom))
}
