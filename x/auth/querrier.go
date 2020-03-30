package auth

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	sdkAuth "github.com/cosmos/cosmos-sdk/x/auth"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	QueryMultiSigAcc = "get_multisig_acc"
	QueryPendingTx   = "get_multisig_pending_tx"
)

func NewQuerier(cdc *codec.Codec, accountKeeper sdkAuth.AccountKeeper) sdkTypes.Querier {
	return func(ctx sdkTypes.Context, path []string, req abci.RequestQuery) ([]byte, sdkTypes.Error) {
		switch path[0] {
		case QueryMultiSigAcc:
			return queryMultiSigAcc(cdc, ctx, path[1:], req, accountKeeper)
		case QueryPendingTx:
			return queryMultiSigPendingTx(cdc, ctx, path[1:], req, accountKeeper)
		default:
			return nil, sdkTypes.ErrUnknownRequest("unknown mxw/Auth query endpoint")
		}
	}
}

func queryMultiSigAcc(cdc *codec.Codec, ctx sdkTypes.Context, path []string, _ abci.RequestQuery, accountKeeper sdkAuth.AccountKeeper) ([]byte, sdkTypes.Error) {

	if len(path) != 1 {
		return nil, sdkTypes.ErrUnknownRequest(fmt.Sprintf("Invalid path %s", strings.Join(path, "/")))
	}

	groupAddr, err := sdkTypes.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, sdkTypes.ErrUnknownAddress(fmt.Sprintf("Invalid group address %s", path[0]))
	}

	groupAcc := accountKeeper.GetAccount(ctx, groupAddr)

	respData := cdc.MustMarshalJSON(groupAcc)

	return respData, nil
}

func queryMultiSigPendingTx(cdc *codec.Codec, ctx sdkTypes.Context, path []string, _ abci.RequestQuery, accountKeeper sdkAuth.AccountKeeper) ([]byte, sdkTypes.Error) {

	if len(path) != 2 {
		return nil, sdkTypes.ErrUnknownRequest(fmt.Sprintf("Invalid path %s", strings.Join(path, "/")))
	}

	groupAddr, err := sdkTypes.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, sdkTypes.ErrUnknownAddress(fmt.Sprintf("Invalid group address %s", path[0]))
	}

	txID, err := strconv.ParseUint(path[1], 0, 64)
	if err != nil {
		return nil, sdkTypes.ErrUnknownRequest(fmt.Sprintf("Invalid txID %s", path[1]))
	}

	groupAcc := accountKeeper.GetAccount(ctx, groupAddr)
	if !groupAcc.IsMultiSig() {
		return nil, nil
	}
	multiSig := groupAcc.GetMultiSig()
	pendingTx := multiSig.GetPendingTx(txID)
	if pendingTx == nil {
		return nil,nil
	}
	tx := pendingTx.GetTx()

	stdTx, ok := tx.(sdkAuth.StdTx)
	if !ok {
		return nil, sdkTypes.ErrInternal("Tx must be StdTx.")
	}

	respData := cdc.MustMarshalJSON(stdTx)

	return respData, nil
}
