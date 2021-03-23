package hype

import (
	"github.com/CityOfZion/neo-storm/examples/token/nep5"

	"github.com/CityOfZion/neo-storm/interop/runtime"
	"github.com/CityOfZion/neo-storm/interop/storage"
	"github.com/CityOfZion/neo-storm/interop/util"
)

const (
	decimals   = 4
	multiplier = decimals * 10
)

var owner = util.FromAddress("AUxNUdS3Vxwv2giJT8VA89dyo7paM3jqAG")

// CreateToken initializes the Token Interface for the Smart Contract to operate with
func CreateToken() nep5.Token {
	return nep5.Token{
		Name:           "Hypecoin",
		Symbol:         "HYPE",
		Decimals:       decimals,
		Owner:          owner,
		TotalSupply:    1350 * multiplier,
		CirculationKey: "TokenCirculation",
	}
}

// Main function = contract entry
func Main(operation string, args []interface{}) interface{} {
	token := CreateToken()

	if operation == "name" {
		return token.Name
	}
	if operation == "symbol" {
		return token.Symbol
	}
	if operation == "decimals" {
		return token.Decimals
	}

	// The following operations need ctx
	ctx := storage.GetContext()

	if operation == "totalSupply" {
		return token.GetSupply(ctx)
	}
	if operation == "balanceOf" && CheckArgs(args, 1) {
		hodler := args[0].([]byte)
		return token.BalanceOf(ctx, hodler)
	}
	if operation == "transfer" && CheckArgs(args, 3) {
		from := args[0].([]byte)
		to := args[1].([]byte)
		amount := args[2].(int)
		return token.Transfer(ctx, from, to, amount)
	}
	if operation == "deploy" {
		return Deploy(ctx, token)
	}
	if operation == "circulation" {
		return GetCirculation(ctx, token)
	}

	return true
}

// CheckArgs checks args array against a length indicator
func CheckArgs(args []interface{}, length int) bool {
	if len(args) == length {
		return true
	}

	return false
}

// Deploy NEP-5 Token to blockchain
func Deploy(ctx storage.Context, t nep5.Token) bool {
	if !runtime.CheckWitness(t.Owner) {
		return false
	}

	if !storage.Get(ctx, "initialized").(bool) {
		storage.Put(ctx, "initialized", 1)
		storage.Put(ctx, t.Owner, t.TotalSupply)
		return AddToCirculation(ctx, t)
	}

	return false

}

// AddToCirculation add NEP-5 into circulation
func AddToCirculation(ctx storage.Context, t nep5.Token) bool {
	currentSupply := storage.Get(ctx, t.CirculationKey).(int)

	currentSupply += t.TotalSupply
	storage.Put(ctx, t.CirculationKey, currentSupply)
	return true
}

// GetCirculation get current circulation amount
func GetCirculation(ctx storage.Context, t nep5.Token) interface{} {
	return storage.Get(ctx, t.CirculationKey)
}
