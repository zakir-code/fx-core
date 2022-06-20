package v020

import (
	"fmt"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"

	bsctypes "github.com/functionx/fx-core/x/bsc/types"
	polygontypes "github.com/functionx/fx-core/x/polygon/types"
	trontypes "github.com/functionx/fx-core/x/tron/types"

	abci "github.com/tendermint/tendermint/abci/types"

	evmtypes "github.com/evmos/ethermint/x/evm/types"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	ethermint "github.com/evmos/ethermint/types"

	migratetypes "github.com/functionx/fx-core/x/migrate/types"

	erc20types "github.com/functionx/fx-core/x/erc20/types"

	"github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibcconnectiontypes "github.com/cosmos/ibc-go/v3/modules/core/03-connection/types"
	ibckeeper "github.com/cosmos/ibc-go/v3/modules/core/keeper"

	fxtypes "github.com/functionx/fx-core/types"
	erc20keeper "github.com/functionx/fx-core/x/erc20/keeper"
)

// CreateUpgradeHandler creates an SDK upgrade handler for v2
func CreateUpgradeHandler(
	kvStoreKeyMap map[string]*sdk.KVStoreKey,
	mm *module.Manager, configurator module.Configurator,
	bankKeeper bankKeeper.Keeper, accountKeeper authkeeper.AccountKeeper,
	paramsKeeper paramskeeper.Keeper, ibcKeeper *ibckeeper.Keeper,
	erc20Keeper erc20keeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		cacheCtx, commit := ctx.CacheContext()

		// 1. clear testnet module data
		clearTestnetModule(ctx, kvStoreKeyMap)

		// 2. update FX metadata
		updateFXMetadata(cacheCtx, bankKeeper, kvStoreKeyMap)

		// 3. update block params (max_gas:3000000000)
		updateBlockParams(cacheCtx, paramsKeeper)

		// 4. migrate base account to eth account
		migrateAccountToEth(cacheCtx, accountKeeper)

		// set max expected block time parameter. Replace the default with your expected value
		// https://github.com/cosmos/ibc-go/blob/release/v1.0.x/docs/ibc/proto-docs.md#params-2
		ibcKeeper.ConnectionKeeper.SetParams(cacheCtx, ibcconnectiontypes.DefaultParams())

		// cosmos-sdk 0.42.x from version must be empty
		if len(fromVM) != 0 {
			panic("invalid from version map")
		}

		for n, m := range mm.Modules {
			//NOTE: fromVM empty
			if initGenesis[n] {
				continue
			}
			if v, ok := runMigrates[n]; ok {
				fromVM[n] = v
				continue
			}
			fromVM[n] = m.ConsensusVersion()
		}

		if mm.OrderMigrations == nil {
			mm.OrderMigrations = migrationsOrder(mm.ModuleNames())
		}
		cacheCtx.Logger().Info("start to run module v2 migrations...")
		toVersion, err := mm.RunMigrations(cacheCtx, configurator, fromVM)
		if err != nil {
			return nil, fmt.Errorf("run migrations error %s", err.Error())
		}

		// clear metadata except FX
		clearTestnetDenom(ctx, kvStoreKeyMap)

		// register coin
		for _, metadata := range fxtypes.GetMetadata() {
			cacheCtx.Logger().Info("add metadata", "coin", metadata.String())
			pair, err := erc20Keeper.RegisterCoin(cacheCtx, metadata)
			if err != nil {
				return nil, fmt.Errorf("register %s error %s", metadata.Base, err.Error())
			}
			cacheCtx.EventManager().EmitEvent(sdk.NewEvent(
				erc20types.EventTypeRegisterCoin,
				sdk.NewAttribute(erc20types.AttributeKeyDenom, pair.Denom),
				sdk.NewAttribute(erc20types.AttributeKeyTokenAddress, pair.Erc20Address),
			))
		}

		//commit upgrade
		commit()

		return toVersion, nil
	}
}

func updateFXMetadata(ctx sdk.Context, bankKeeper bankKeeper.Keeper, keys map[string]*sdk.KVStoreKey) {
	md := fxtypes.GetFXMetaData(fxtypes.DefaultDenom)
	if err := md.Validate(); err != nil {
		panic(fmt.Sprintf("invalid %s metadata", fxtypes.DefaultDenom))
	}
	key, ok := keys[banktypes.StoreKey]
	if !ok {
		panic("bank key store not found")
	}
	logger := ctx.Logger()
	logger.Info("update FX metadata", "metadata", md.String())
	//delete fx
	fxDenom := strings.ToLower(fxtypes.DefaultDenom)
	denomMetaDataStore := prefix.NewStore(ctx.KVStore(key), banktypes.DenomMetadataKey(fxDenom))
	denomMetaDataStore.Delete([]byte(fxDenom))
	//set FX
	bankKeeper.SetDenomMetaData(ctx, md)
}

func migrateAccountToEth(ctx sdk.Context, ak authkeeper.AccountKeeper) {
	logger := ctx.Logger()
	logger.Info("migrate account to eth", "network", fxtypes.Network())
	// migrate base account to eth account
	ak.IterateAccounts(ctx, func(account authtypes.AccountI) (stop bool) {
		if _, ok := account.(ethermint.EthAccountI); ok {
			return false
		}
		baseAccount, ok := account.(*authtypes.BaseAccount)
		if !ok {
			logger.Info("ignore account", "address", account.GetAddress(), "type", fmt.Sprintf("%T", account))
			return false
		}
		ethAccount := &ethermint.EthAccount{
			BaseAccount: baseAccount,
			CodeHash:    common.BytesToHash(emptyCodeHash).String(),
		}
		ak.SetAccount(ctx, ethAccount)
		logger.Info("migrate account to eth", "address", account.GetAddress())
		return false
	})
}

func updateBlockParams(ctx sdk.Context, pk paramskeeper.Keeper) {
	logger := ctx.Logger()
	logger.Info("update block params", "network", fxtypes.Network())
	baseappSubspace, found := pk.GetSubspace(baseapp.Paramspace)
	if !found {
		panic(fmt.Sprintf("unknown subspace: %s", baseapp.Paramspace))
	}
	var bp abci.BlockParams
	baseappSubspace.Get(ctx, baseapp.ParamStoreKeyBlockParams, &bp)
	logger.Info("update block params", "before update", bp.String())
	bp.MaxGas = blockParamsMaxGas
	baseappSubspace.Set(ctx, baseapp.ParamStoreKeyBlockParams, bp)
	logger.Info("update block params", "after update", bp.String())
}

func migrationsOrder(modules []string) []string {
	modules = module.DefaultMigrationsOrder(modules)
	orders := make([]string, 0, len(modules))
	for _, name := range modules {
		if name == bsctypes.ModuleName || name == polygontypes.ModuleName || name == trontypes.ModuleName ||
			name == feemarkettypes.ModuleName || name == evmtypes.ModuleName ||
			name == erc20types.ModuleName || name == migratetypes.ModuleName {
			continue
		}
		orders = append(orders, name)
	}
	orders = append(orders, []string{
		bsctypes.ModuleName, polygontypes.ModuleName, trontypes.ModuleName,
		feemarkettypes.ModuleName, evmtypes.ModuleName,
		erc20types.ModuleName, migratetypes.ModuleName,
	}...)
	return orders
}

func clearTestnetDenom(ctx sdk.Context, keys map[string]*types.KVStoreKey) {
	if fxtypes.NetworkTestnet() != fxtypes.Network() {
		return
	}
	key, ok := keys[banktypes.StoreKey]
	if !ok {
		panic("bank key store not found")
	}
	logger := ctx.Logger()
	logger.Info("clear testnet metadata", "network", fxtypes.Network())
	for _, md := range fxtypes.GetMetadata() {
		//remove denom except FX
		if md.Base == fxtypes.DefaultDenom {
			continue
		}
		logger.Info("clear testnet metadata", "metadata", md.String())
		denomMetaDataStore := prefix.NewStore(ctx.KVStore(key), banktypes.DenomMetadataKey(md.Base))
		denomMetaDataStore.Delete([]byte(md.Base))
	}
}

func clearTestnetModule(ctx sdk.Context, keys map[string]*types.KVStoreKey) {
	logger := ctx.Logger()
	if fxtypes.NetworkTestnet() != fxtypes.Network() {
		return
	}
	logger.Info("clear kv store", "network", fxtypes.Network())
	cleanModules := []string{feemarkettypes.StoreKey, evmtypes.StoreKey, erc20types.StoreKey, migratetypes.StoreKey}
	multiStore := ctx.MultiStore()
	for _, storeName := range cleanModules {
		logger.Info("clear kv store", "storesName", storeName)
		startTime := time.Now()
		storeKey, ok := keys[storeName]
		if !ok {
			panic(fmt.Sprintf("%s store not found", storeName))
		}
		kvStore := multiStore.GetKVStore(storeKey)
		if err := deleteKVStore(kvStore); err != nil {
			panic(fmt.Sprintf("failed to delete store %s: %s", storeName, err.Error()))
		}
		logger.Info("clear kv store done", "storesName", storeName, "consumeMs", time.Now().UnixNano()-startTime.UnixNano())
	}
}

func deleteKVStore(kv types.KVStore) error {
	// Note that we cannot write while iterating, so load all keys here, delete below
	var keys [][]byte
	itr := kv.Iterator(nil, nil)
	defer itr.Close()

	for itr.Valid() {
		keys = append(keys, itr.Key())
		itr.Next()
	}

	for _, k := range keys {
		kv.Delete(k)
	}
	return nil
}