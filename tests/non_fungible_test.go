package tests

import (
	"testing"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	nonFungible "github.com/maxonrow/maxonrow-go/x/token/nonfungible"
)

type NonFungibleTokenInfo struct {
	Action                       string
	ApplicationFee               string
	FeeCollector                 string
	Name                         string
	Symbol                       string
	Owner                        string
	NewOwner                     string
	TokenMetadata                string
	ItemID                       string
	Properties                   string
	Metadata                     string
	Approved                     bool
	Frozen                       bool
	Burnable                     bool
	Modifiable                   bool
	Public                       bool
	Provider                     string
	ProviderNonce                string
	Issuer                       string
	FeeSettingName               string
	VerifyTransferTokenOwnership string
	TransferLimit                string
	MintLimit                    string
	EndorserList                 []string
}

func makeNonFungibleTokenTxs() []*testCase {

	tcs := []*testCase{

		// 1.0 create non fungible - without ItemID
		{"nonFungibleToken", false, false, "Create non fungible token - Happy Path.  commit", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT", "acc-40", "", "token metadata", "test test test", "hi bye", "hi bye", false, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "Create non fungible token - Happy Path.", nil},
		{"nonFungibleToken", true, true, "Create non fungible token - Token already exists (TNFT). commit", "acc-29", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT", "acc-29", "", "token metadata", "test test test", "hi bye", "hi bye", false, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "Create non fungible token - Token already exists (TFT).", nil},
		{"nonFungibleToken", true, true, "Create non fungible token - Insufficient fee amount. commit", "acc-29", "0cin", 0, NonFungibleTokenInfo{"create", "0", "nft-mostafa", "TestNonFungibleToken", "TNFT", "acc-29", "", "token metadata", "test test test", "hi bye", "hi bye", false, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "Create non fungible token - Insufficient fee amount.", nil},
		{"nonFungibleToken", true, true, "Create non fungible token - Very long metadata!", "acc-29", "0cin", 0, NonFungibleTokenInfo{"create", "0", "nft-mostafa", "TestNonFungibleToken", "TNFT", "acc-29", "", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaAAAAAA---aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaAAAAAA", "test test test", "hi bye", "hi bye", false, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "Create non fungible token - Very long metadata!", nil},
		{"nonFungibleToken", true, true, "Create non fungible token - Fee collector invalid", "acc-29", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "nft-yk", "TestNonFungibleToken-191", "TNFT-191", "acc-29", "", "token metadata", "test test test", "hi bye", "hi bye", false, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "goh-Create non fungible token - Fee collector invalid", nil},                                                                         // ok
		{"nonFungibleToken", true, true, "Create non fungible token - Invalid fee amount", "acc-29", "100000000cin", 0, NonFungibleTokenInfo{"create", "abcXXX", "nft-mostafa", "TestNonFungibleToken-191", "TNFT-191", "acc-29", "", "token metadata", "test test test", "hi bye", "hi bye", false, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "goh-Create non fungible token - Invalid fee amount", nil},                                                                            // ok
		{"nonFungibleToken", true, true, "Create non fungible token - Insufficient balance to pay for application fee", "acc-29", "100000000cin", 0, NonFungibleTokenInfo{"create", "77999999999999900000000", "nft-mostafa", "TestNonFungibleToken-191", "TNFT-191", "acc-29", "", "token metadata", "test test test", "hi bye", "hi bye", false, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "goh-Create non fungible token - Insufficient balance to pay for application fee", nil}, //ok

		// 1.1 Update-Token-Metadata : Failed
		{"nonFungibleToken", true, true, "Update NFT Metadata non fungible token - Token yet to approved, can not edit-metadata.", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"update-nft-metadata", "0", "nft-mostafa", "TestNonFungibleToken", "TNFT", "acc-40", "", "token metadata", "updated here by goh 0111", "properties", "metadata", false, false, false, false, false, "nft-jeansoon", "0", "nft-carlo", "default", "", "2", "2", []string{"nft-jeansoon", "nft-carlo"}}, "", nil}, //ok

		// 2. approve non fungible - without ItemID
		{"nonFungibleToken", false, false, "Approve non fungible token(TNFT) : TransferLimit(2) Mintlimit(2) Endorser(nft-jeanson,nft-carlo) - Happy path", "nft-mostafa", "0cin", 0, NonFungibleTokenInfo{"approve", "", "", "TestNonFungibleToken", "TNFT", "", "", "", "", "properties", "metadata", false, false, true, false, false, "nft-jeansoon", "0", "nft-carlo", "default", "", "2", "2", []string{"nft-jeansoon", "nft-carlo"}}, "", nil},
		{"nonFungibleToken", true, true, "Approve non fungible token(TNFT) : The Signer Not authorised to approve", "nft-yk", "100000000cin", 0, NonFungibleTokenInfo{"approve", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT", "", "", "", "", "properties", "metadata", false, false, true, false, false, "nft-jeansoon", "0", "nft-carlo", "default", "", "2", "2", []string{"nft-jeansoon", "nft-carlo"}}, "", nil},                    //ok
		{"nonFungibleToken", true, true, "Approve non fungible token(TNFT-191) : The Token symbol does not exist", "nft-mostafa", "100000000cin", 0, NonFungibleTokenInfo{"approve", "10000000", "nft-mostafa", "TestNonFungibleToken-191", "TNFT-191", "", "", "", "", "properties", "metadata", false, false, true, false, false, "nft-jeansoon", "0", "nft-carlo", "default", "", "2", "2", []string{"nft-jeansoon", "nft-carlo"}}, "", nil},        //ok
		{"nonFungibleToken", true, true, "Approve non fungible token(TNFT) : Unauthorized signature - nft-yk", "nft-mostafa", "100000000cin", 0, NonFungibleTokenInfo{"approve", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT", "", "", "", "", "properties", "metadata", false, false, true, false, false, "nft-jeansoon", "0", "nft-yk", "default", "", "2", "2", []string{"nft-jeansoon", "nft-carlo"}}, "", nil},                       //ok
		{"nonFungibleToken", true, true, "Approve non fungible token(TNFT) : Fee setting is not valid - fee-setting-191", "nft-mostafa", "100000000cin", 0, NonFungibleTokenInfo{"approve", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT", "", "", "", "", "properties", "metadata", false, false, true, false, false, "nft-jeansoon", "0", "nft-carlo", "fee-setting-191", "", "2", "2", []string{"nft-jeansoon", "nft-carlo"}}, "", nil}, //ok
		{"nonFungibleToken", true, true, "Approve non fungible token(TNFT) : Token already approved - TNFT", "nft-mostafa", "0cin", 0, NonFungibleTokenInfo{"approve", "", "", "TestNonFungibleToken", "TNFT", "", "", "", "", "properties", "metadata", true, false, true, false, false, "nft-jeansoon", "0", "nft-carlo", "default", "", "2", "2", []string{"nft-jeansoon", "nft-carlo"}}, "", nil},                                                  // ok

		// 2.1 Update-Token-Metadata - without ItemID (Logic : after approved only can apply this process)
		{"nonFungibleToken", false, false, "Update NFT Metadata non fungible token - Happy Path.  commit", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"update-nft-metadata", "0", "nft-mostafa", "TestNonFungibleToken", "TNFT", "acc-40", "", "token metadata updated here by goh 001", "", "", "", false, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},
		{"nonFungibleToken", true, true, "Update NFT Metadata non fungible token - Invalid metadata field length due to very long token metadata!", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"update-nft-metadata", "0", "nft-mostafa", "TestNonFungibleToken", "TNFT", "acc-40", "", "NFT Metadata updated : aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaAAAAAA------aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaBBBBBB----aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaCCCCCC------aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaAAAAAA-DDDDDD------aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaAAAAAA-PPPPPP------aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaAAAAAA-QQQQQQ------aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaAAAAAA-RRRRRR------aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaAAAAAA-SSSSSS", "", "", "", false, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},
		{"nonFungibleToken", true, true, "Update NFT Metadata non fungible token - Invalid account owner.", "maintainer-1", "100000000cin", 0, NonFungibleTokenInfo{"update-nft-metadata", "0", "nft-mostafa", "TestNonFungibleToken", "TNFT", "maintainer-1", "", "token metadata updated here by goh 002", "", "", "", false, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},    //ok
		{"nonFungibleToken", true, true, "Update NFT Metadata non fungible token - Invalid token-symbol.", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"update-nft-metadata", "0", "nft-mostafa", "TestNonFungibleToken", "TNFT-UPDATE-METADATA", "acc-40", "", "token metadata updated here by goh 003", "", "", "", false, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil}, //ok
		{"nonFungibleToken", true, true, "Update NFT Metadata non fungible token - symbol is empty.", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"update-nft-metadata", "0", "nft-mostafa", "TestNonFungibleToken", "", "acc-40", "", "token metadata updated here by goh 004", "", "", "", false, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},                          //ok
		{"nonFungibleToken", true, true, "Update NFT Metadata non fungible token - Owner must passed KYC.", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"update-nft-metadata", "0", "nft-mostafa", "TestNonFungibleToken", "TNFT", "acc-57", "", "token metadata updated here by goh 005", "", "", "", false, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},                //ok

		// 3.0 mint non fungible - with ItemID
		{"nonFungibleToken", false, false, "Mint non fungible token - Happy path", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"mint", "", "", "", "TNFT", "acc-40", "nft-mostafa", "", "112233", "properties", "metadata", true, false, true, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},                     //ok
		{"nonFungibleToken", false, false, "Mint non fungible token - (mint for burn) Happy path", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"mint", "", "", "", "TNFT", "acc-40", "nft-yk", "", "223344", "properties", "metadata", true, false, true, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},          //ok
		{"nonFungibleToken", false, false, "Mint non fungible token - (mint for endorsement) Happy path", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"mint", "", "", "", "TNFT", "acc-40", "nft-nago", "", "334455", "properties", "metadata", true, false, true, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil}, //ok
		{"nonFungibleToken", true, true, "Mint non fungible token - Invalid Token Symbol", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"mint", "", "", "", "TNFT-191", "acc-40", "nft-mostafa", "", "112233", "properties", "metadata", true, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},        //ok
		{"nonFungibleToken", false, true, "Mint non fungible token - Token item id is in used.", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"mint", "", "", "", "TNFT", "acc-40", "nft-mostafa", "", "112233", "properties", "metadata", true, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},      //ok
		{"nonFungibleToken", false, true, "Mint non fungible token - Invalid token minter.", "nft-yk", "100000000cin", 0, NonFungibleTokenInfo{"mint", "", "", "", "TNFT", "nft-yk", "nft-mostafa", "", "112233", "properties", "metadata", true, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},          //ok

		// 3.1 mint item - Modifiable, Public(TRUE)
		{"nonFungibleToken", false, false, "Create (by Public==TRUE) non fungible token(TNFT-public-01) - Happy Path. commit", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT-public-01", "acc-40", "", "token metadata", "test test test", "hi bye", "hi bye", false, false, false, false, true, "", "", "", "", "", "", "", []string{""}}, "Create (by Public==TRUE) non fungible token(TNFT-public-01) - Happy Path. commit", nil},
		{"nonFungibleToken", false, false, "Approve (by Public==TRUE) non fungible token(TNFT-public-01) - Happy path", "nft-mostafa", "0cin", 0, NonFungibleTokenInfo{"approve", "", "", "TestNonFungibleToken", "TNFT-public-01", "", "", "token metadata", "", "properties", "metadata", false, false, false, false, true, "nft-jeansoon", "0", "nft-carlo", "default", "", "2", "2", []string{"nft-jeansoon", "nft-carlo"}}, "Approve (by Public==TRUE) non fungible token(TNFT-public-01) - Happy path", nil},
		{"nonFungibleToken", false, false, "Mint (by Public==TRUE) non fungible token(TNFT-public-01) - Happy path", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"mint", "", "", "", "TNFT-public-01", "acc-40", "acc-40", "token metadata", "223344", "properties", "metadata", true, false, false, false, true, "", "", "", "", "", "", "", []string{""}}, "Mint (by Public==TRUE) non fungible token(TNFT-public-01) - Happy path", nil},
		{"nonFungibleToken", false, true, "Mint (by Public==TRUE) non fungible token(TNFT-public-01) - Error, Public token can only be minted to itself. ", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"mint", "", "", "", "TNFT-public-01", "acc-40", "nft-yk", "token metadata", "223344", "properties", "metadata", true, false, false, false, true, "", "", "", "", "", "", "", []string{""}}, "Mint (by Public==TRUE) non fungible token(TNFT-public-01) - Error, Public token can only be minted to itself.", nil}, //ok

		// 3.2 Update-Item-Metadata - with ItemID
		{"nonFungibleToken", false, false, "Update Item Metadata non fungible token - Happy Path.  commit", "nft-mostafa", "100000000cin", 0, NonFungibleTokenInfo{"update-item-metadata", "", "", "", "TNFT", "nft-mostafa", "", "token metadata", "112233", "properties", "update Item metadata 9991", true, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},
		{"nonFungibleToken", true, true, "Update Item Metadata non fungible token - Invalid account owner.", "maintainer-1", "100000000cin", 0, NonFungibleTokenInfo{"update-item-metadata", "", "", "", "TNFT", "maintainer-1", "", "token metadata", "112233", "properties", "update Item metadata 9991", true, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},    //ok
		{"nonFungibleToken", false, true, "Update Item Metadata non fungible token - Invalid Item Id.", "nft-mostafa", "100000000cin", 0, NonFungibleTokenInfo{"update-item-metadata", "", "", "", "TNFT", "nft-mostafa", "", "token metadata", "771177", "properties", "update Item metadata 8881", true, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},           //ok
		{"nonFungibleToken", false, true, "Update Item Metadata non fungible token - Item owner not match.", "yk", "100000000cin", 0, NonFungibleTokenInfo{"update-item-metadata", "", "", "", "TNFT", "yk", "", "token metadata", "112233", "properties", "update Item metadata 7771", true, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},                        //ok
		{"nonFungibleToken", true, true, "Update Item Metadata non fungible token - Item Id is empty.", "nft-mostafa", "100000000cin", 0, NonFungibleTokenInfo{"update-item-metadata", "", "", "", "TNFT", "nft-mostafa", "", "token metadata", "", "properties", "update Item metadata 8881", true, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},                 //ok, base on 'ValidateBasic()'
		{"nonFungibleToken", true, true, "Update Item Metadata non fungible token - Owner must passed KYC.", "nft-mostafa", "100000000cin", 0, NonFungibleTokenInfo{"update-item-metadata", "", "", "", "TNFT", "acc-57", "", "token metadata", "112233", "properties", "update Item metadata 9991", true, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},           //ok
		{"nonFungibleToken", false, true, "Update Item Metadata non fungible token - Invalid token symbol.", "nft-mostafa", "100000000cin", 0, NonFungibleTokenInfo{"update-item-metadata", "", "", "", "TNFT-9911", "nft-mostafa", "", "token metadata", "771177", "properties", "update Item metadata 8881", true, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil}, //ok
		{"nonFungibleToken", true, true, "Update Item Metadata non fungible token - Invalid metadata field length due to very long Item metadata!", "nft-mostafa", "100000000cin", 0, NonFungibleTokenInfo{"update-item-metadata", "", "", "", "TNFT", "nft-mostafa", "", "token metadata", "112233", "properties", "Item metadata updated : aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaAAAAAA------aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaBBBBBB----aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaCCCCCC------aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaAAAAAA-DDDDDD------aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaAAAAAA-PPPPPP------aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaAAAAAA-QQQQQQ------aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaAAAAAA-RRRRRR------aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaAAAAAA-SSSSSS", true, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},

		//====================== with ItemID :
		// 4. make endorsement - with ItemID
		{"nonFungibleToken", false, false, "endorse a nonfungible item - Happy path", "nft-carlo", "100000000cin", 0, NonFungibleTokenInfo{"endorsement", "", "", "", "TNFT", "nft-carlo", "", "token metadata", "334455", "", "", true, false, false, false, false, "", "", "", "", "", "", "", []string{"nft-jeansoon", "nft-carlo"}}, "", nil},
		{"nonFungibleToken", true, true, "endorse a nonfungible item - Invalid endorser", "nft-yk", "100000000cin", 0, NonFungibleTokenInfo{"endorsement", "", "", "", "TNFT", "nft-yk", "", "token metadata", "778899", "properties", "metadata", true, false, false, false, false, "", "", "", "", "", "", "", []string{"nft-jeansoon", "nft-carlo"}}, "", nil},                // ok
		{"nonFungibleToken", false, true, "endorse a nonfungible item - Invalid Token Symbol", "nft-carlo", "100000000cin", 0, NonFungibleTokenInfo{"endorsement", "", "", "", "TNFT-111", "nft-carlo", "", "token metadata", "334455", "properties", "metadata", true, false, false, false, false, "", "", "", "", "", "", "", []string{"nft-jeansoon", "nft-carlo"}}, "", nil}, //ok
		{"nonFungibleToken", false, true, "endorse a nonfungible item - Invalid Item-ID", "nft-carlo", "100000000cin", 0, NonFungibleTokenInfo{"endorsement", "", "", "", "TNFT", "nft-carlo", "", "token-metadata", "999111", "properties", "metadata", true, false, false, false, false, "", "", "", "", "", "", "", []string{"nft-jeansoon", "nft-carlo"}}, "", nil},          // ok

		// 5. transfer non fungible item - with ItemID
		{"nonFungibleToken", false, false, "Transfer non fungible token item - Happy path", "nft-mostafa", "100000000cin", 0, NonFungibleTokenInfo{"transfer", "", "", "", "TNFT", "nft-mostafa", "nft-yk", "token metadata", "112233", "properties", "metadata", true, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},
		{"nonFungibleToken", true, true, "Transfer non fungible token item - Invalid Token Symbol", "nft-mostafa", "100000000cin", 0, NonFungibleTokenInfo{"transfer", "", "", "", "TNFT-111", "nft-mostafa", "nft-yk", "token metadata", "112233", "properties", "metadata", true, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},         // ok
		{"nonFungibleToken", true, true, "Transfer non fungible token item - Invalid Account to transfer from", "nft-carlo", "100000000cin", 0, NonFungibleTokenInfo{"transfer", "", "", "", "TNFT-111", "nft-carlo", "nft-yk", "token metadata", "112233", "properties", "metadata", true, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil}, // ok
		{"nonFungibleToken", false, true, "Transfer non fungible token item - Invalid Item-ID", "nft-mostafa", "100000000cin", 0, NonFungibleTokenInfo{"transfer", "", "", "", "TNFT", "nft-mostafa", "nft-yk", "token metadata", "999111", "properties", "metadata", true, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},                 // ok

		// 6. burn non fungible item - with ItemID
		{"nonFungibleToken", false, true, "Burn non fungible token item - Invalid token owner", "nft-carlo", "100000000cin", 0, NonFungibleTokenInfo{"burn", "", "", "", "TNFT", "nft-carlo", "", "token metadata", "223344", "", "", false, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil}, //ok
		{"nonFungibleToken", false, false, "Burn non fungible token item -  Happy path", "nft-yk", "100000000cin", 0, NonFungibleTokenInfo{"burn", "", "", "", "TNFT", "nft-yk", "", "token metadata", "223344", "", "", true, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},
		{"nonFungibleToken", true, true, "Burn non fungible token item - Invalid Token Symbol", "nft-yk", "100000000cin", 0, NonFungibleTokenInfo{"burn", "", "", "", "TNFT-111", "nft-yk", "", "token metadata", "223344", "", "", false, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},                         //ok
		{"nonFungibleToken", false, true, "Burn non fungible token item - Invalid Item-ID", "nft-yk", "100000000cin", 0, NonFungibleTokenInfo{"burn", "", "", "", "TNFT", "nft-yk", "", "token metadata", "999111", "", "", false, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},                                 //ok
		{"nonFungibleToken", true, true, "Burn non fungible token item - Invalid account to burn from due to yet pass KYC", "nft-yk", "100000000cin", 0, NonFungibleTokenInfo{"burn", "", "", "", "TNFT", "acc-19", "", "token metadata", "223344", "", "", false, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil}, //ok

		// Special cases-01 =============================== start goh : base on 'TNFT-191'
		// create non fungible :
		{"nonFungibleToken", false, false, "Create non fungible token - Happy Path.  commit", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "nft-mostafa", "TestNonFungibleToken-191", "TNFT-191", "acc-40", "", "token metadata", "", "", "", false, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "Create non fungible token - Happy Path.", nil}, // ok
		// mint non fungible :
		{"nonFungibleToken", true, true, "Mint non fungible token item - Invalid token as yet to approved", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"mint", "", "", "", "TNFT-191", "acc-40", "nft-mostafa", "token metadata", "112233", "properties", "metadata", false, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil}, //ok
		// burn non fungible :
		{"nonFungibleToken", true, true, "Burn non fungible token item - Invalid token", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"burn", "", "", "", "'TNFT-191'", "acc-40", "nft-mostafa", "token metadata", "112233", "properties", "metadata", false, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil}, //ok
		// Special cases-01 =============================== end goh : base on 'TNFT-191'

		//====================== without ItemID :
		// 7. transfer ownership - without ItemID
		{"nonFungibleToken", false, false, "Transfer non fungible token ownership - Happy path", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"transfer-ownership", "", "", "", "TNFT", "acc-40", "nft-yk", "token metadata", "", "", "", false, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},
		{"nonFungibleToken", false, false, "Create non fungible token - TNFT-T1 Happy Path.  commit", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT-T1", "acc-40", "", "token metadata", "test test test", "hi bye", "hi bye", false, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},        //ok
		{"nonFungibleToken", true, true, "Transfer non fungible token ownership - Invalid token as yet to approved", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"transfer-ownership", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT-T1", "acc-40", "nft-yk", "token metadata", "", "", "", false, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil}, //ok
		{"nonFungibleToken", false, false, "Create non fungible token - TNFT-T2 Happy Path.  commit", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT-T2", "acc-40", "", "token metadata", "test test test", "hi bye", "hi bye", false, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},        //ok
		{"nonFungibleToken", true, true, "Transfer non fungible token ownership - Invalid token owner", "nft-yk", "100000000cin", 0, NonFungibleTokenInfo{"transfer-ownership", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT-T2", "nft-yk", "acc-40", "token metadata", "", "", "", false, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},              //ok

		// 8. verify transfer token ownership - without ItemID
		{"nonFungibleToken", false, false, "Approve non fungible token transfer ownership - Happy path for TNFT", "nft-mostafa", "0cin", 0, NonFungibleTokenInfo{"verify-transfer-tokenOwnership", "", "", "", "TNFT", "nft-mostafa", "nft-yk", "token metadata", "", "", "", true, false, false, false, false, "nft-jeansoon", "0", "nft-carlo", "", "APPROVE_TRANFER_TOKEN_OWNERSHIP", "", "", []string{""}}, "", nil},

		// 9. accept token ownership - without ItemID
		{"nonFungibleToken", false, false, "Accept non fungible token ownership - Happy path. commit", "nft-yk", "100000000cin", 0, NonFungibleTokenInfo{"accept-ownership", "", "", "", "TNFT", "", "nft-yk", "token metadata", "", "", "", false, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},
		{"nonFungibleToken", false, false, "Create non fungible token - TNFT-Q1 Happy Path.  commit", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT-Q1", "acc-40", "", "token metadata", "test test test", "hi bye", "hi bye", false, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},                                 //ok
		{"nonFungibleToken", true, true, "Accept non fungible token ownership - Invalid token as yet to approved", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"transfer-ownership", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT-Q1", "acc-40", "nft-yk", "token metadata", "", "", "", false, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},                            //ok
		{"nonFungibleToken", false, false, "Create non fungible token - TNFT-Q2 Happy Path.  commit", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT-Q2", "acc-40", "", "token metadata", "test test test", "hi bye", "hi bye", false, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},                                 //ok
		{"nonFungibleToken", true, true, "Accept non fungible token ownership - Invalid token new-owner", "nft-yk", "100000000cin", 0, NonFungibleTokenInfo{"transfer-ownership", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT-Q2", "nft-yk", "acc-40", "token metadata", "", "", "", false, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},                                     //ok
		{"nonFungibleToken", false, false, "Create non fungible token - TNFT-Q3 Happy Path.  commit", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT-Q3", "acc-40", "", "token, metadata", "test test test", "hi bye", "hi bye", false, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},                                //ok
		{"nonFungibleToken", true, true, "Accept non fungible token ownership - Invalid token due to IsTokenOwnershipTransferrable == FALSE", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"transfer-ownership", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT-Q3", "acc-40", "nft-yk", "token metadata", "", "", "", false, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil}, //ok

		// 10. freeze non fungible token - without ItemID
		{"nonFungibleToken", false, false, "Freeze non fungible token - Happy path. commit", "nft-mostafa", "0cin", 0, NonFungibleTokenInfo{"freeze", "", "", "", "TNFT", "", "", "token metadata", "", "", "", true, false, false, false, false, "nft-jeansoon", "0", "nft-carlo", "", "", "", "", []string{""}}, "", nil},
		{"nonFungibleToken", true, true, "Transfer non fungible token ownership - Invalid token action (due to Token not approved) ", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"transfer-ownership", "", "", "", "TNFT", "acc-40", "nft-yk", "token metadata", "", "", "", false, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "", nil},
		{"nonFungibleToken", false, false, "Create non fungible token - Happy Path.  commit", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "nft-mostafa", "TestNonFungibleToken-192", "TNFT-192", "acc-40", "", "token metadata", "", "", "", false, false, false, false, false, "", "", "", "", "", "", "", []string{""}}, "Create non fungible token - Happy Path.", nil},     // ok
		{"nonFungibleToken", true, true, "Freeze non fungible token - Not authorised to approve due to Invalid Fee collector", "acc-29", "100000000cin", 0, NonFungibleTokenInfo{"freeze", "10000000", "nft-yk", "TestNonFungibleToken-192", "TNFT-192", "", "", "token metadata", "", "", "", true, false, false, false, false, "nft-jeansoon", "0", "nft-carlo", "", "", "", "", []string{""}}, "", nil}, //ok
		{"nonFungibleToken", true, true, "Freeze non fungible token - Invalid Token symbol - TNFT-111", "mostafa", "100000000cin", 0, NonFungibleTokenInfo{"freeze", "10000000", "nft-mostafa", "TestNonFungibleToken-111", "TNFT-111", "", "", "token metadata", "", "", "", true, false, false, false, false, "nft-jeansoon", "0", "nft-carlo", "", "", "", "", []string{""}}, "", nil},                  //ok

		// 11. unfreeze non fungible token - without ItemID
		{"nonFungibleToken", false, false, "Unfreeze non fungible token - Happy path", "nft-mostafa", "0cin", 0, NonFungibleTokenInfo{"unfreeze", "", "", "", "TNFT", "", "", "token metadata", "", "", "", true, false, false, false, false, "nft-jeansoon", "0", "nft-carlo", "", "", "", "", []string{""}}, "", nil},
		{"nonFungibleToken", true, true, "Unfreeze non fungible token - Not authorised to approve due to Invalid Fee collector", "acc-29", "100000000cin", 0, NonFungibleTokenInfo{"unfreeze", "10000000", "nft-yk", "TestNonFungibleToken-192", "TNFT-192", "", "", "token metadata", "", "", "", true, false, false, false, false, "nft-jeansoon", "0", "nft-carlo", "", "", "", "", []string{""}}, "", nil}, // ok
		{"nonFungibleToken", true, true, "Unfreeze non fungible token - Invalid Token symbol - TNFT-111", "nft-mostafa", "100000000cin", 0, NonFungibleTokenInfo{"unfreeze", "10000000", "nft-mostafa", "TestNonFungibleToken-111", "TNFT-111", "", "", "token metadata", "", "", "", true, false, false, false, false, "nft-jeansoon", "0", "nft-carlo", "", "", "", "", []string{""}}, "", nil},              // ok

		//====================== without ItemID :
		// freeze and THEN unfreeze
		{"nonFungibleToken", false, false, "CREATE non fungible token [TNFT-B2] - Happy Path.  commit", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT-B2", "acc-40", "", "token metadata", "", "", "", false, false, false, false, false, "", "", "", "", "", "", "", []string{"nft-jeansoon", "nft-carlo"}}, "", nil},
		{"nonFungibleToken", false, false, "APPROVE non fungible token [TNFT-B2] - Happy path.  commit", "nft-mostafa", "100000000cin", 0, NonFungibleTokenInfo{"approve", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT-B2", "acc-40", "", "token metadata", "", "properties", "metadata", false, false, false, false, false, "nft-jeansoon", "0", "nft-carlo", "default", "", "2", "2", []string{"nft-jeansoon", "nft-carlo"}}, "", nil}, //ok
		{"nonFungibleToken", false, false, "MINT non fungible item [TNFT-B2] - Happy path", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"mint", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT-B2", "acc-40", "nft-mostafa", "token metadata", "001177", "properties", "metadata", true, false, false, false, false, "nft-jeansoon", "0", "nft-carlo", "", "", "", "", []string{"nft-jeansoon", "nft-carlo"}}, "", nil},               //ok
		{"nonFungibleToken", false, false, "FREEZE non fungible item [TNFT-B2] - Happy path.", "nft-mostafa", "100000000cin", 0, NonFungibleTokenInfo{"freeze-item", "10000000", "", "", "TNFT-B2", "nft-mostafa", "", "token metadata", "001177", "properties", "metadata", true, false, false, false, false, "nft-jeansoon", "0", "nft-carlo", "", "", "", "", []string{"nft-jeansoon", "nft-carlo"}}, "", nil},                                     //ok
		{"nonFungibleToken", false, false, "UNFREEZE non fungible item [TNFT-B2] - Happy path.", "nft-mostafa", "100000000cin", 0, NonFungibleTokenInfo{"unfreeze-item", "10000000", "nft-mostafa", "", "TNFT-B2", "nft-mostafa", "", "token metadata", "001177", "properties", "metadata", true, false, false, false, false, "nft-jeansoon", "0", "nft-carlo", "", "", "", "", []string{"nft-jeansoon", "nft-carlo"}}, "", nil},                      //ok

		// 12. freeze non fungible item - with ItemID
		{"nonFungibleToken", false, false, "CREATE non fungible token [TNFT-D2] - Happy Path.  commit", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT-D2", "acc-40", "", "token metadata", "", "", "", false, false, false, false, false, "", "", "", "", "", "", "", []string{"nft-jeansoon", "nft-carlo"}}, "", nil},                                                         //ok
		{"nonFungibleToken", false, false, "APPROVE non fungible token [TNFT-D2] - Happy path.  commit", "nft-mostafa", "100000000cin", 0, NonFungibleTokenInfo{"approve", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT-D2", "acc-40", "", "token metadata", "", "properties", "metadata", false, false, false, false, false, "nft-jeansoon", "0", "nft-carlo", "default", "", "2", "2", []string{"nft-jeansoon", "nft-carlo"}}, "", nil}, //ok
		{"nonFungibleToken", false, false, "MINT non fungible item [TNFT-D2] - Happy path", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"mint", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT-D2", "acc-40", "nft-mostafa", "token metadata", "880099", "properties", "metadata", true, true, false, false, false, "nft-jeansoon", "0", "nft-carlo", "", "", "", "", []string{"nft-jeansoon", "nft-carlo"}}, "", nil},                //ok
		{"nonFungibleToken", false, false, "Freeze non fungible item [TNFT-D2] - Happy path.", "nft-mostafa", "100000000cin", 0, NonFungibleTokenInfo{"freeze-item", "10000000", "", "", "TNFT-D2", "nft-mostafa", "", "token metadata", "880099", "", "", true, true, false, false, false, "nft-jeansoon", "0", "nft-carlo", "", "", "", "", []string{""}}, "", nil},
		{"nonFungibleToken", true, true, "Freeze non fungible item [TNFT-D2] - Invalid signer.", "nft-jeansoon", "100000000cin", 0, NonFungibleTokenInfo{"freeze-item", "10000000", "", "", "TNFT-D2", "nft-jeansoon", "", "token metadata", "880099", "", "", true, false, false, false, false, "nft-jeansoon", "1", "nft-carlo", "", "", "", "", []string{""}}, "", nil},                //ok
		{"nonFungibleToken", false, true, "Freeze non fungible item [TNFT-D2] - No such non fungible token.", "nft-yk", "100000000cin", 0, NonFungibleTokenInfo{"freeze-item", "10000000", "", "", "TNFT-9988", "nft-yk", "", "token metadata", "880099", "", "", true, false, false, false, false, "nft-jeansoon", "1", "nft-carlo", "", "", "", "", []string{""}}, "", nil},             //ok
		{"nonFungibleToken", false, true, "Freeze non fungible item [TNFT-D2] - No such item to freeze.", "nft-yk", "100000000cin", 0, NonFungibleTokenInfo{"freeze-item", "10000000", "", "", "TNFT-D2", "nft-yk", "", "token metadata", "991111", "", "", true, false, false, false, false, "nft-jeansoon", "1", "nft-carlo", "", "", "", "", []string{""}}, "", nil},                   //ok
		{"nonFungibleToken", true, true, "Freeze non fungible item [TNFT-D2] - Not authorised to freeze non token item.", "nft-yk", "100000000cin", 0, NonFungibleTokenInfo{"freeze-item", "10000000", "", "", "TNFT-D2", "nft-yk", "", "token metadata", "880099", "", "", true, false, false, false, false, "nft-jeansoon", "1", "nft-carlo", "", "", "", "", []string{""}}, "", nil},   //ok
		{"nonFungibleToken", true, true, "Freeze non fungible item [TNFT-D2] - Non Fungible item already frozen.", "nft-mostafa", "100000000cin", 0, NonFungibleTokenInfo{"freeze-item", "10000000", "", "", "TNFT-D2", "nft-mostafa", "", "token metadata", "880099", "", "", true, true, false, false, false, "nft-jeansoon", "1", "nft-carlo", "", "", "", "", []string{""}}, "", nil}, //ok
		{"nonFungibleToken", false, true, "Unfreeze non fungible item [TNFT-D2] - Invalid nonce.", "nft-mostafa", "100000000cin", 0, NonFungibleTokenInfo{"unfreeze-item", "10000000", "", "", "TNFT-D2", "nft-mostafa", "", "token metadata", "880099", "", "", true, true, false, false, false, "nft-jeansoon", "2", "nft-carlo", "", "", "", "", []string{""}}, "", nil},               //ok
		//{"nonFungibleToken", false, false, "Unfreeze non fungible item [TNFT-D2] - Happy path.", "nft-mostafa", "100000000cin", 0, NonFungibleTokenInfo{"unfreeze-item", "10000000", "", "", "TNFT-D2", "nft-mostafa", "", "token metadata", "880099", "", "", true, true, false, false, false, "nft-jeansoon", "1", "nft-carlo", "", "", "", "", []string{""}}, "", nil},                 //ok

		// 13. unfreeze non fungible item - with ItemID
		{"nonFungibleToken", false, false, "CREATE non fungible token [TNFT-E2] - Happy Path.  commit", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"create", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT-E2", "acc-40", "", "token metadata", "", "", "", false, false, false, false, false, "", "", "", "", "", "", "", []string{"nft-jeansoon", "nft-carlo"}}, "", nil},
		{"nonFungibleToken", false, false, "APPROVE non fungible token [TNFT-E2] - Happy path.  commit", "nft-mostafa", "100000000cin", 0, NonFungibleTokenInfo{"approve", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT-E2", "acc-40", "", "token metadata", "", "properties", "metadata", false, false, false, false, false, "nft-jeansoon", "1", "nft-carlo", "default", "", "2", "2", []string{"nft-jeansoon", "nft-carlo"}}, "", nil},
		{"nonFungibleToken", false, false, "MINT non fungible item [TNFT-E2] - Happy path", "acc-40", "100000000cin", 0, NonFungibleTokenInfo{"mint", "10000000", "nft-mostafa", "TestNonFungibleToken", "TNFT-E2", "acc-40", "nft-mostafa", "token metadata", "770099", "properties", "metadata", true, false, false, false, false, "nft-jeansoon", "0", "nft-carlo", "", "", "", "", []string{"nft-jeansoon", "nft-carlo"}}, "", nil}, //ok
		{"nonFungibleToken", false, true, "Unfreeze non fungible item [TNFT-E2] - Invalid signer.", "nft-jeansoon", "100000000cin", 0, NonFungibleTokenInfo{"unfreeze-item", "10000000", "", "", "TNFT-E2", "nft-jeansoon", "", "token metadata", "770099", "", "", true, false, false, false, false, "nft-jeansoon", "1", "nft-carlo", "", "", "", "", []string{""}}, "", nil},                                                         //ok
		{"nonFungibleToken", false, true, "Unfreeze non fungible item [TNFT-E2] - No such non fungible token.", "nft-yk", "100000000cin", 0, NonFungibleTokenInfo{"unfreeze-item", "10000000", "", "", "TNFT-9988", "nft-yk", "", "token metadata", "770099", "", "", true, false, false, false, false, "nft-jeansoon", "1", "nft-carlo", "", "", "", "", []string{""}}, "", nil},                                                       //ok
		{"nonFungibleToken", false, true, "Unfreeze non fungible item [TNFT-E2] - No such  non fungible item to unfreeze.", "nft-yk", "100000000cin", 0, NonFungibleTokenInfo{"unfreeze-item", "10000000", "", "", "TNFT-E2", "nft-yk", "", "token metadata", "991111", "", "", true, false, false, false, false, "nft-jeansoon", "1", "nft-carlo", "", "", "", "", []string{""}}, "", nil},                                             //ok
		{"nonFungibleToken", false, true, "Unfreeze non fungible item [TNFT-E2] - Not authorised to unfreeze token account.", "nft-yk", "100000000cin", 0, NonFungibleTokenInfo{"unfreeze-item", "10000000", "", "", "TNFT-E2", "nft-yk", "", "token metadata", "770099", "", "", true, false, false, false, false, "nft-jeansoon", "1", "nft-carlo", "", "", "", "", []string{""}}, "", nil},                                           //ok
		{"nonFungibleToken", false, true, "Unfreeze non fungible item [TNFT-E2] - Non fungible item not frozen.", "nft-mostafa", "100000000cin", 0, NonFungibleTokenInfo{"unfreeze-item", "10000000", "", "", "TNFT-E2", "nft-mostafa", "", "token metadata", "770099", "", "", true, true, false, false, false, "nft-jeansoon", "1", "nft-carlo", "", "", "", "", []string{""}}, "", nil},                                              //ok

	}

	return tcs
}

func makeCreateNonFungibleTokenMsg(t *testing.T, name string, symbol string, metadata string, owner string, applicationFee string, tokenFeeCollector string) sdkTypes.Msg {

	// create new token
	ownerAddr := tKeys[owner].addr
	fee := nonFungible.Fee{
		To:    tKeys[tokenFeeCollector].addr,
		Value: applicationFee,
	}
	msgCreateNonFungibleToken := nonFungible.NewMsgCreateNonFungibleToken(symbol, ownerAddr, name, "", metadata, fee)

	return msgCreateNonFungibleToken
}

func makeApproveNonFungibleTokenMsg(t *testing.T, signer string, provider string, providerNonce string, issuer string, symbol string, status string, feeSettingName string, mintLimit, transferLimit string, endorserList []string, burnable bool, modifiable bool, public bool) sdkTypes.Msg {

	providerAddr := tKeys[provider].addr

	var tokenFee = []nonFungible.TokenFee{
		{
			Action:  "transfer",
			FeeName: feeSettingName,
		},
		{
			Action:  "mint",
			FeeName: feeSettingName,
		},
		{
			Action:  "burn",
			FeeName: feeSettingName,
		},
		{
			Action:  "transferOwnership",
			FeeName: feeSettingName,
		},
		{
			Action:  "acceptOwnership",
			FeeName: feeSettingName,
		},
	}

	mintL := sdkTypes.NewUintFromString(mintLimit)
	transferL := sdkTypes.NewUintFromString(transferLimit)

	var endorsers []sdkTypes.AccAddress

	for _, v := range endorserList {
		addr := tKeys[v].addr
		endorsers = append(endorsers, addr)
	}

	tokenDoc := nonFungible.NewToken(providerAddr, providerNonce, status, symbol, transferL, mintL, tokenFee, endorsers, burnable, true, modifiable, public)

	// provider sign the token
	tokenBz, err := tCdc.MarshalJSON(tokenDoc)
	require.NoError(t, err)
	signedTokenBz, err := tKeys[provider].priv.Sign(sdkTypes.MustSortJSON(tokenBz))
	require.NoError(t, err)

	tokenPayload := nonFungible.NewPayload(*tokenDoc, tKeys[provider].pub, signedTokenBz)

	// issuer sign the tokenPayload
	tokenPayloadBz, err := tCdc.MarshalJSON(tokenPayload)
	require.NoError(t, err)
	signedTokenPayloadBz, err := tKeys[issuer].priv.Sign(sdkTypes.MustSortJSON(tokenPayloadBz))
	require.NoError(t, err)

	var signatures []nonFungible.Signature
	signature := nonFungible.NewSignature(tKeys[issuer].pub, signedTokenPayloadBz)
	signatures = append(signatures, signature)

	msgSetFungibleTokenStatus := nonFungible.NewMsgSetNonFungibleTokenStatus(tKeys[signer].addr, *tokenPayload, signatures)

	return msgSetFungibleTokenStatus
}

//module of transfer
func makeTransferNonFungibleTokenMsg(t *testing.T, owner string, newOwner string, symbol string, itemID string) sdkTypes.Msg {

	ownerAddr := tKeys[owner].addr
	newOwnerAddr := tKeys[newOwner].addr

	msgTransferPayload := nonFungible.NewMsgTransferNonFungibleToken(symbol, ownerAddr, newOwnerAddr, itemID)
	return msgTransferPayload

}

//module of mint
func makeMintNonFungibleTokenMsg(t *testing.T, owner string, newOwner string, symbol string, itemID string, properties, metadata string) sdkTypes.Msg {

	ownerAddr := tKeys[owner].addr
	newOwnerAddr := tKeys[newOwner].addr

	msgMintPayload := nonFungible.NewMsgMintNonFungibleToken(ownerAddr, symbol, newOwnerAddr, itemID, properties, metadata)
	return msgMintPayload

}

//module of burn
func makeBurnNonFungibleTokenMsg(t *testing.T, owner string, symbol string, itemID string) sdkTypes.Msg {

	ownerAddr := tKeys[owner].addr

	msgBurnNonFungible := nonFungible.NewMsgBurnNonFungibleToken(symbol, ownerAddr, itemID)
	return msgBurnNonFungible

}

//moduel of transferOwnership
func makeTransferNonFungibleTokenOwnershipMsg(t *testing.T, owner string, newOwner string, symbol string) sdkTypes.Msg {

	ownerAddr := tKeys[owner].addr
	newOwnerAddr := tKeys[newOwner].addr

	msgTransferOwnershipPayload := nonFungible.NewMsgTransferNonFungibleTokenOwnership(symbol, ownerAddr, newOwnerAddr)
	return msgTransferOwnershipPayload

}

//module of acceptOwnership
func makeAcceptNonFungibleTokenOwnershipMsg(t *testing.T, newOwner string, symbol string) sdkTypes.Msg {

	fromAddr := tKeys[newOwner].addr

	msgAcceptOwnershipPayload := nonFungible.NewMsgAcceptNonFungibleTokenOwnership(symbol, fromAddr)
	return msgAcceptOwnershipPayload

}

func makeFreezeNonFungibleTokenMsg(t *testing.T, signer string, provider string, providerNonce string, issuer string, symbol string, burnable bool, modifiable bool, public bool) sdkTypes.Msg {

	status := "FREEZE"
	providerAddr := tKeys[provider].addr

	tokenDoc := nonFungible.NewToken(providerAddr, providerNonce, status, symbol, sdkTypes.NewUint(0), sdkTypes.NewUint(0), nil, nil, burnable, true, modifiable, public) // status : FREEZE / UNFREEZE

	// provider sign the token
	tokenBz, err := tCdc.MarshalJSON(tokenDoc)
	require.NoError(t, err)
	signedTokenBz, err := tKeys[provider].priv.Sign(sdkTypes.MustSortJSON(tokenBz))
	require.NoError(t, err)

	tokenPayload := nonFungible.NewPayload(*tokenDoc, tKeys[provider].pub, signedTokenBz)

	// issuer sign the tokenPayload
	tokenPayloadBz, err := tCdc.MarshalJSON(tokenPayload)
	require.NoError(t, err)
	signedTokenPayloadBz, err := tKeys[issuer].priv.Sign(sdkTypes.MustSortJSON(tokenPayloadBz))
	require.NoError(t, err)

	var signatures []nonFungible.Signature
	signature := nonFungible.NewSignature(tKeys[issuer].pub, signedTokenPayloadBz)
	signatures = append(signatures, signature)

	msgSetNonFungibleTokenStatus := nonFungible.NewMsgSetNonFungibleTokenStatus(tKeys[signer].addr, *tokenPayload, signatures)

	return msgSetNonFungibleTokenStatus
}

func makeUnfreezeNonFungibleTokenMsg(t *testing.T, signer string, provider string, providerNonce string, issuer string, symbol string, burnable bool, modifiable bool, public bool) sdkTypes.Msg {

	status := "UNFREEZE"
	providerAddr := tKeys[provider].addr

	tokenDoc := nonFungible.NewToken(providerAddr, providerNonce, status, symbol, sdkTypes.NewUint(0), sdkTypes.NewUint(0), nil, nil, burnable, true, modifiable, public) // status : FREEZE / UNFREEZE

	// provider sign the token
	tokenBz, err := tCdc.MarshalJSON(tokenDoc)
	require.NoError(t, err)
	signedTokenBz, err := tKeys[provider].priv.Sign(sdkTypes.MustSortJSON(tokenBz))
	require.NoError(t, err)

	tokenPayload := nonFungible.NewPayload(*tokenDoc, tKeys[provider].pub, signedTokenBz)

	// issuer sign the tokenPayload
	tokenPayloadBz, err := tCdc.MarshalJSON(tokenPayload)
	require.NoError(t, err)
	signedTokenPayloadBz, err := tKeys[issuer].priv.Sign(sdkTypes.MustSortJSON(tokenPayloadBz))
	require.NoError(t, err)

	var signatures []nonFungible.Signature
	signature := nonFungible.NewSignature(tKeys[issuer].pub, signedTokenPayloadBz)
	signatures = append(signatures, signature)

	msgSetNonFungibleTokenStatus := nonFungible.NewMsgSetNonFungibleTokenStatus(tKeys[signer].addr, *tokenPayload, signatures)

	return msgSetNonFungibleTokenStatus
}

func makeVerifyTransferNonFungibleTokenOwnership(t *testing.T, signer, provider, providerNonce, issuer, symbol, action string, burnable bool, modifiable bool, public bool) sdkTypes.Msg {

	providerAddr := tKeys[provider].addr

	// burnable and tokenfees is not in used for verifying transfer token status, we just set it to false and leave it empty.
	verifyTransferTokenOwnershipDoc := nonFungible.NewToken(providerAddr, providerNonce, action, symbol, sdkTypes.NewUint(0), sdkTypes.NewUint(0), []nonFungible.TokenFee{}, nil, burnable, true, modifiable, public)

	// provider sign
	verifyTransferTokenOwnershipDocBz, err := tCdc.MarshalJSON(verifyTransferTokenOwnershipDoc)
	require.NoError(t, err)
	signedVerifyTransferTokenOwnershipDoc, err := tKeys[provider].priv.Sign(sdkTypes.MustSortJSON(verifyTransferTokenOwnershipDocBz))
	require.NoError(t, err)

	verifyTransferTokenOwnershipPayload := nonFungible.NewPayload(*verifyTransferTokenOwnershipDoc, tKeys[provider].pub, signedVerifyTransferTokenOwnershipDoc)

	// issuer sign
	verifyTransferPayloadBz, err := tCdc.MarshalJSON(verifyTransferTokenOwnershipPayload)
	require.NoError(t, err)
	signedVerifyTransferPayloadBz, err := tKeys[issuer].priv.Sign(sdkTypes.MustSortJSON(verifyTransferPayloadBz))
	require.NoError(t, err)

	var signatures []nonFungible.Signature
	signature := nonFungible.NewSignature(tKeys[issuer].pub, signedVerifyTransferPayloadBz)
	signatures = append(signatures, signature)

	msgVerifyTransferNonFungibleTokenOwnership := nonFungible.NewMsgSetNonFungibleTokenStatus(tKeys[signer].addr, *verifyTransferTokenOwnershipPayload, signatures)

	return msgVerifyTransferNonFungibleTokenOwnership
}

func makeEndorsement(t *testing.T, signer, to, symbol string, itemID string) sdkTypes.Msg {

	signerAddr := tKeys[signer].addr

	return nonFungible.NewMsgEndorsement(symbol, signerAddr, itemID)
}

// Freeze Item
func makeFreezeNonFungibleItemMsg(t *testing.T, signer string, provider string, providerNonce string, issuer string, symbol string, itemID string) sdkTypes.Msg {

	status := "FREEZE_ITEM"
	providerAddr := tKeys[provider].addr

	itemDetails := nonFungible.NewItemDetails(providerAddr, providerNonce, status, symbol, itemID) // status : FREEZE / UNFREEZE

	// provider sign the item
	itemBz, err := tCdc.MarshalJSON(itemDetails)
	require.NoError(t, err)
	signedItemBz, err := tKeys[provider].priv.Sign(sdkTypes.MustSortJSON(itemBz))
	require.NoError(t, err)

	itemPayload := nonFungible.NewItemPayload(*itemDetails, tKeys[provider].pub, signedItemBz)

	// issuer sign the itemPayload
	itemPayloadBz, err := tCdc.MarshalJSON(itemPayload)
	require.NoError(t, err)
	signedItemPayloadBz, err := tKeys[issuer].priv.Sign(sdkTypes.MustSortJSON(itemPayloadBz))
	require.NoError(t, err)

	var signatures []nonFungible.Signature
	signature := nonFungible.NewSignature(tKeys[issuer].pub, signedItemPayloadBz)
	signatures = append(signatures, signature)

	msgSetNonFungibleTokenStatus := nonFungible.NewMsgSetNonFungibleItemStatus(tKeys[signer].addr, *itemPayload, signatures)
	return msgSetNonFungibleTokenStatus
}

// UnFreeze Item
func makeUnfreezeNonFungibleItemMsg(t *testing.T, signer string, provider string, providerNonce string, issuer string, symbol string, itemID string) sdkTypes.Msg {

	status := "UNFREEZE_ITEM"
	providerAddr := tKeys[provider].addr

	itemDetails := nonFungible.NewItemDetails(providerAddr, providerNonce, status, symbol, itemID) // status : FREEZE / UNFREEZE

	// provider sign the item
	itemBz, err := tCdc.MarshalJSON(itemDetails)
	require.NoError(t, err)
	signedItemBz, err := tKeys[provider].priv.Sign(sdkTypes.MustSortJSON(itemBz))
	require.NoError(t, err)

	itemPayload := nonFungible.NewItemPayload(*itemDetails, tKeys[provider].pub, signedItemBz)

	// issuer sign the itemPayload
	itemPayloadBz, err := tCdc.MarshalJSON(itemPayload)
	require.NoError(t, err)
	signedItemPayloadBz, err := tKeys[issuer].priv.Sign(sdkTypes.MustSortJSON(itemPayloadBz))
	require.NoError(t, err)

	var signatures []nonFungible.Signature
	signature := nonFungible.NewSignature(tKeys[issuer].pub, signedItemPayloadBz)
	signatures = append(signatures, signature)

	msgSetNonFungibleTokenStatus := nonFungible.NewMsgSetNonFungibleItemStatus(tKeys[signer].addr, *itemPayload, signatures)
	return msgSetNonFungibleTokenStatus
}

// UpdateItemMetadata
func makeUpdateItemMetadataMsg(t *testing.T, symbol string, owner string, itemID string, metadata string) sdkTypes.Msg {

	ownerAddr := tKeys[owner].addr

	msgUpdateItemMetadataPayload := nonFungible.NewMsgUpdateItemMetadata(symbol, ownerAddr, itemID, metadata)
	return msgUpdateItemMetadataPayload

}

// UpdateNFTMetadata
func makeUpdateNFTMetadataMsg(t *testing.T, symbol string, owner string, metadata string) sdkTypes.Msg {

	ownerAddr := tKeys[owner].addr

	msgUpdateNFTMetadataPayload := nonFungible.NewMsgUpdateNFTMetadata(symbol, ownerAddr, metadata)
	return msgUpdateNFTMetadataPayload

}
