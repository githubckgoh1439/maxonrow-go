package nonfungible

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/maxonrow/maxonrow-go/types"
)

func (k *Keeper) MintNonFungibleToken(ctx sdkTypes.Context, symbol string, from sdkTypes.AccAddress, to sdkTypes.AccAddress, itemID, properties, metadata string) sdkTypes.Result {

	nonFungibleToken := new(Token)

	if exists := k.getTokenData(ctx, symbol, nonFungibleToken); !exists {
		return types.ErrInvalidTokenSymbol(symbol).Result()
	}

	// get minter account.
	// if token is public that means minter can be anyone.
	minterAccount := k.accountKeeper.GetAccount(ctx, from)
	if minterAccount == nil {
		return types.ErrInvalidTokenAccount().Result()
	}

	if !nonFungibleToken.Flags.HasFlag(PubFlag) {
		if !nonFungibleToken.Owner.Equals(from) {
			return types.ErrInvalidTokenMinter().Result()
		}
	} else {
		if !from.Equals(to) {
			return sdkTypes.ErrInternal("Public token can only be minted to oneself.").Result()
		}
	}

	if !nonFungibleToken.Flags.HasFlag(MintFlag) {
		return types.ErrInvalidTokenAction().Result()
	}

	if !nonFungibleToken.Flags.HasFlag(ApprovedFlag) {
		return types.ErrTokenInvalid().Result()
	}

	if nonFungibleToken.Flags.HasFlag(FrozenFlag) {
		return types.ErrTokenFrozen().Result()
	}

	if !k.IsItemIDUnique(ctx, nonFungibleToken.Symbol, itemID) {
		return types.ErrTokenItemIDInUsed().Result()
	}

	amt := sdkTypes.NewUint(1)
	nonFungibleToken.TotalSupply = nonFungibleToken.TotalSupply.Add(amt)

	// check mint limit, if token mint limit !=0
	if !nonFungibleToken.MintLimit.IsZero() {

		mintLimitKey := getMintItemLimitKey(nonFungibleToken.Symbol, to)

		store := ctx.KVStore(k.key)
		limit := store.Get(mintLimitKey)
		if limit != nil {
			if sdkTypes.NewUintFromString(string(limit)).GTE(nonFungibleToken.MintLimit) {
				return sdkTypes.ErrInternal("Holding limit existed.").Result()
			}

			k.increaseMintItemLimit(ctx, symbol, to)
		}
	}

	k.storeToken(ctx, symbol, nonFungibleToken)

	item := k.createNonFungibleItem(ctx, nonFungibleToken.Symbol, to, itemID, properties, metadata)

	eventParam := []string{symbol, string(item.ID), from.String(), to.String()}
	eventSignature := "MintedNonFungibleItem(string,string,string,string)"

	accountSequence := minterAccount.GetSequence()
	resultLog := types.NewResultLog(accountSequence, ctx.TxBytes())

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, from.String(), eventParam),
		Log:    resultLog.String(),
	}
}

//* TransferNonFungibleToken
func (k *Keeper) TransferNonFungibleToken(ctx sdkTypes.Context, symbol string, from, to sdkTypes.AccAddress, itemID string) sdkTypes.Result {
	var token = new(Token)
	if exists := k.getTokenData(ctx, symbol, token); !exists {
		return types.ErrTokenInvalid().Result()
	}

	if !token.Flags.HasFlag(TransferableFlag) {
		return types.ErrInvalidTokenAction().Result()
	}

	fromAccount := k.accountKeeper.GetAccount(ctx, from)
	if fromAccount == nil {
		return types.ErrInvalidTokenAccount().Result()
	}

	if token.Flags.HasFlag(FrozenFlag) {
		return types.ErrTokenFrozen().Result()
	}

	itemKey := getNonFungibleItemKey(symbol, []byte(itemID))
	ownerKey := getNonFungibleItemOwnerKey(symbol, []byte(itemID))

	store := ctx.KVStore(k.key)

	ownerValue := store.Get(ownerKey)
	if ownerValue == nil {
		return types.ErrInvalidTokenOwner().Result()
	}

	itemValue := store.Get(itemKey)
	if itemValue == nil {
		return types.ErrInvalidTokenOwner().Result()
	}

	var item = new(Item)
	k.cdc.MustUnmarshalBinaryLengthPrefixed(itemValue, item)

	if item.TransferLimit.GTE(token.TransferLimit) {

		// TO-DO: own error message.
		return sdkTypes.ErrInternal("Item has existed transfer limit.").Result()
	}

	// delete old owner
	store.Delete(ownerKey)

	// set to new owner
	store.Set(ownerKey, to.Bytes())

	// increase the transfer limit and set
	item.TransferLimit = item.TransferLimit.Add(sdkTypes.NewUint(1))
	itemData := k.cdc.MustMarshalBinaryLengthPrefixed(item)
	store.Set(itemKey, itemData)

	eventParam := []string{symbol, string(itemID), from.String(), to.String()}
	eventSignature := "TransferredNonFungibleItem(string,string,string,string)"

	accountSequence := fromAccount.GetSequence()
	resultLog := types.NewResultLog(accountSequence, ctx.TxBytes())

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, from.String(), eventParam),
		Log:    resultLog.String(),
	}

}

// BurnFungibleToken
func (k *Keeper) BurnNonFungibleToken(ctx sdkTypes.Context, symbol string, from sdkTypes.AccAddress, itemID string) sdkTypes.Result {
	var token = new(Token)
	if exists := k.getTokenData(ctx, symbol, token); !exists {
		return types.ErrInvalidTokenSymbol(symbol).Result()
	}

	if !token.Flags.HasFlag(BurnFlag) {
		return types.ErrInvalidTokenAction().Result()
	}

	fromAccount := k.accountKeeper.GetAccount(ctx, from)
	if fromAccount == nil {
		return sdkTypes.ErrInvalidSequence("Invalid account to burn from.").Result()
	}

	if !token.Flags.HasFlag(ApprovedFlag) {
		return types.ErrTokenInvalid().Result()
	}

	if token.Flags.HasFlag(FrozenFlag) {
		return types.ErrTokenFrozen().Result()
	}

	item := k.getNonFungibleItem(ctx, symbol, itemID)
	if item == nil {
		return types.ErrInvalidTokenOwner().Result()
	}

	itemOwner := k.getNonFungibleItemOwner(ctx, symbol, itemID)
	if !itemOwner.Equals(from) {
		return types.ErrInvalidTokenOwner().Result()
	}

	ownerKey := getNonFungibleItemOwnerKey(symbol, []byte(itemID))
	itemKey := getNonFungibleItemKey(symbol, []byte(itemID))

	store := ctx.KVStore(k.key)

	store.Delete(itemKey)
	store.Delete(ownerKey)

	eventParam := []string{symbol, string(item.ID), from.String()}
	eventSignature := "BurnedNonFungibleItem(string,string,string)"

	accountSequence := fromAccount.GetSequence()
	resultLog := types.NewResultLog(accountSequence, ctx.TxBytes())

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, from.String(), eventParam),
		Log:    resultLog.String(),
	}

}

// TO-DO: proper implementation to cater nonfungibletoken transfer ownership
func (k *Keeper) TransferTokenOwnership(ctx sdkTypes.Context, symbol string, from sdkTypes.AccAddress, to sdkTypes.AccAddress) sdkTypes.Result {

	var token = new(Token)

	err := k.mustGetTokenData(ctx, symbol, token)
	if err != nil {
		return err.Result()
	}

	return (k.transferNonFungibleTokenOwnership(ctx, from, to, token))

}

func (k *Keeper) transferNonFungibleTokenOwnership(ctx sdkTypes.Context, from sdkTypes.AccAddress, to sdkTypes.AccAddress, token *Token) sdkTypes.Result {

	ownerWalletAccount := k.accountKeeper.GetAccount(ctx, token.Owner)
	if ownerWalletAccount == nil {
		return types.ErrInvalidTokenOwner().Result()
	}

	if ownerWalletAccount != nil && !token.Owner.Equals(from) {
		return types.ErrInvalidTokenOwner().Result()
	}

	if !token.IsApproved() {
		// TODO: Please define an error code
		return sdkTypes.ErrUnknownRequest("Token is not approved.").Result()
	}

	if token.IsFrozen() {
		return types.ErrTokenFrozen().Result()
	}

	// set token newowner to new owner, pending for accepting by new owner
	token.NewOwner = to
	token.Flags.AddFlag(TransferTokenOwnershipFlag)

	k.storeToken(ctx, token.Symbol, token)

	eventParam := []string{token.Symbol, from.String(), to.String()}
	eventSignature := "TransferredNonFungibleTokenOwnership(string,string,string)"

	accountSequence := ownerWalletAccount.GetSequence()
	resultLog := types.NewResultLog(accountSequence, ctx.TxBytes())

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, from.String(), eventParam),
		Log:    resultLog.String(),
	}
}

// TO-DO: proper implementation to cater nonfungibletoken accept ownership
func (k *Keeper) AcceptTokenOwnership(ctx sdkTypes.Context, symbol string, from sdkTypes.AccAddress) sdkTypes.Result {

	var token = new(Token)

	err := k.mustGetTokenData(ctx, symbol, token)
	if err != nil {
		return err.Result()
	}

	return (k.acceptNonFungibleTokenOwnership(ctx, from, token))
}

func (k *Keeper) acceptNonFungibleTokenOwnership(ctx sdkTypes.Context, from sdkTypes.AccAddress, token *Token) sdkTypes.Result {

	if !token.Flags.HasFlag(AcceptTokenOwnershipFlag) && !token.Flags.HasFlag(ApproveTransferTokenOwnershipFlag) && !token.Flags.HasFlag(TransferTokenOwnershipFlag) {
		return types.ErrInvalidTokenAction().Result()
	}

	// validation of exisisting owner account
	ownerWalletAccount := k.accountKeeper.GetAccount(ctx, token.Owner)
	if ownerWalletAccount == nil {
		return sdkTypes.ErrInvalidSequence("Invalid token owner.").Result()
	}

	// validation of new owner account
	newOwnerWalletAccount := k.accountKeeper.GetAccount(ctx, token.Owner)
	if newOwnerWalletAccount == nil {
		return sdkTypes.ErrInvalidSequence("Invalid token owner.").Result()
	}

	if newOwnerWalletAccount != nil && token.NewOwner.String() != from.String() {
		return types.ErrInvalidTokenNewOwner().Result()
	}

	if !token.Flags.HasFlag(ApprovedFlag) {
		return sdkTypes.ErrUnknownRequest("Token is not approved.").Result()
	}

	if token.Flags.HasFlag(FrozenFlag) {
		return types.ErrTokenFrozen().Result()
	}

	//TO-DO: if there is need to set token.NewOwner to empty
	// accepting token ownership, remove newowner move the newowner into owner.
	var emptyAccAddr sdkTypes.AccAddress
	token.Owner = from
	token.NewOwner = emptyAccAddr

	token.Flags.RemoveFlag(ApproveTransferTokenOwnershipFlag)
	token.Flags.RemoveFlag(AcceptTokenOwnershipFlag)
	token.Flags.RemoveFlag(TransferTokenOwnershipFlag)
	k.storeToken(ctx, token.Symbol, token)

	eventParam := []string{token.Symbol, from.String()}
	eventSignature := "AcceptedNonFungibleTokenOwnership(string,string)"

	accountSequence := newOwnerWalletAccount.GetSequence()
	resultLog := types.NewResultLog(accountSequence, ctx.TxBytes())

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, from.String(), eventParam),
		Log:    resultLog.String(),
	}

}

func (k *Keeper) MakeEndorsement(ctx sdkTypes.Context, symbol string, from sdkTypes.AccAddress, itemID string) sdkTypes.Result {

	// validation of exisisting owner account
	ownerWalletAccount := k.accountKeeper.GetAccount(ctx, from)
	if ownerWalletAccount == nil {
		return sdkTypes.ErrInvalidSequence("Invalid endorser.").Result()
	}

	item := k.getNonFungibleItem(ctx, symbol, itemID)
	if item == nil {
		return types.ErrTokenInvalid().Result()
	}

	eventParam := []string{symbol, string(itemID), from.String()}
	eventSignature := "EndorsedNonFungibleItem(string,string,string)"

	accountSequence := ownerWalletAccount.GetSequence()
	resultLog := types.NewResultLog(accountSequence, ctx.TxBytes())

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, from.String(), eventParam),
		Log:    resultLog.String(),
	}

}

func (k *Keeper) UpdateItemMetadata(ctx sdkTypes.Context, symbol string, from sdkTypes.AccAddress, itemID string, metadata string) sdkTypes.Result {

	// validation of exisisting owner account
	ownerWalletAccount := k.accountKeeper.GetAccount(ctx, from)
	if ownerWalletAccount == nil {
		return sdkTypes.ErrInvalidSequence("Invalid item owner.").Result()
	}

	item := k.getNonFungibleItem(ctx, symbol, itemID)
	if item == nil {
		return types.ErrTokenInvalid().Result()
	}

	itemOwnerAddr := k.getNonFungibleItemOwner(ctx, symbol, itemID)

	if !ownerWalletAccount.GetAddress().Equals(itemOwnerAddr) {
		return sdkTypes.ErrUnknownAddress("Invalid item owner.").Result()
	}

	item.Metadata = metadata
	k.storeNonFungibleItem(ctx, symbol, from, item)

	accountSequence := ownerWalletAccount.GetSequence()
	resultLog := types.NewResultLog(accountSequence, ctx.TxBytes())

	eventParam := []string{symbol, string(itemID), from.String()}
	eventSignature := "UpdatedNonFungibleItemMetadata(string,string,string)"

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, from.String(), eventParam),
		Log:    resultLog.String(),
	}
}

func (k *Keeper) UpdateNFTMetadata(ctx sdkTypes.Context, symbol string, from sdkTypes.AccAddress, metadata string) sdkTypes.Result {

	// validation of exisisting owner account
	ownerWalletAccount := k.accountKeeper.GetAccount(ctx, from)
	if ownerWalletAccount == nil {
		return types.ErrInvalidTokenOwner().Result()
	}

	var token = new(Token)

	err := k.mustGetTokenData(ctx, symbol, token)
	if err != nil {
		return err.Result()
	}

	if !token.Owner.Equals(from) {
		return types.ErrInvalidTokenOwner().Result()
	}

	if !token.Flags.HasFlag(ApprovedFlag) {
		return types.ErrTokenInvalid().Result()
	}

	token.Metadata = metadata
	k.storeToken(ctx, symbol, token)

	accountSequence := ownerWalletAccount.GetSequence()
	resultLog := types.NewResultLog(accountSequence, ctx.TxBytes())

	eventParam := []string{symbol, from.String()}
	eventSignature := "UpdatedNonFungibleTokenMetadata(string,string)"

	return sdkTypes.Result{
		Events: types.MakeMxwEvents(eventSignature, from.String(), eventParam),
		Log:    resultLog.String(),
	}
}

func (k *Keeper) IsTokenOwnershipAcceptable(ctx sdkTypes.Context, symbol string) bool {

	var token = new(Token)

	err := k.mustGetTokenData(ctx, symbol, token)
	if err != nil {
		return false
	}

	if token.Flags.HasFlag(TransferTokenOwnershipFlag) &&
		token.Flags.HasFlag(AcceptTokenOwnershipFlag) &&
		token.Flags.HasFlag(ApproveTransferTokenOwnershipFlag) {
		return true
	}

	return false
}

func (k *Keeper) IsTokenOwnershipTransferrable(ctx sdkTypes.Context, symbol string) bool {

	var token = new(Token)

	err := k.mustGetTokenData(ctx, symbol, token)
	if err != nil {
		return false
	}

	if token.Flags.HasFlag(TransferTokenOwnershipFlag) || token.Flags.HasFlag(AcceptTokenOwnershipFlag) || token.Flags.HasFlag(ApproveTransferTokenOwnershipFlag) {
		return false
	}

	return true
}

func (k *Keeper) IsTokenOwner(ctx sdkTypes.Context, symbol string, owner sdkTypes.AccAddress) bool {
	var token = new(Token)

	err := k.mustGetTokenData(ctx, symbol, token)
	if err != nil {
		return false
	}

	if token.Owner.Equals(owner) {
		return true
	}

	return false
}

func (k *Keeper) IsTokenNewOwner(ctx sdkTypes.Context, symbol string, newOwner sdkTypes.AccAddress) bool {

	var token = new(Token)

	err := k.mustGetTokenData(ctx, symbol, token)
	if err != nil {
		return false
	}

	if token.NewOwner.Equals(newOwner) {
		return true
	}

	return false
}
