package wallet_handler

import (
	"context"

	"github.com/JackalLabs/jackalgo/utils"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	ecies "github.com/ecies/go/v2"
	"github.com/spf13/pflag"
)

var (
	Bech32Prefix = "jkl"
	// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address
	Bech32PrefixAccAddr = Bech32Prefix
	// Bech32PrefixAccPub defines the Bech32 prefix of an account's public key
	Bech32PrefixAccPub = Bech32Prefix + sdk.PrefixPublic
	// Bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address
	Bech32PrefixValAddr = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator
	// Bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key
	Bech32PrefixValPub = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator + sdk.PrefixPublic
	// Bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address
	Bech32PrefixConsAddr = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus
	// Bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key
	Bech32PrefixConsPub = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus + sdk.PrefixPublic
)

type WalletHandler struct {
	clientCtx client.Context
	address   string
	flags     *pflag.FlagSet
	key       *cryptotypes.PrivKey
	eciesKey  *ecies.PrivateKey
}

func NewWalletHandler(seedPhrase string) (*WalletHandler, error) {
	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	cfg.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	cfg.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)
	// cfg.SetAddressVerifier(wasmtypes.VerifyAddressLen())
	cfg.Seal()

	var pKey *cryptotypes.PrivKey = nil
	address := ""
	if len(seedPhrase) > 0 {
		pKey = cryptotypes.GenPrivKeyFromSecret([]byte(seedPhrase))
		var err error
		address, err = bech32.ConvertAndEncode(Bech32PrefixAccAddr, pKey.PubKey().Address().Bytes())
		if err != nil {
			return nil, err
		}
	}

	k := &client.Context{}

	srvCtx := utils.NewDefaultContext()
	ctx := context.Background()
	ctx = context.WithValue(ctx, client.ClientContextKey, k)
	ctx = context.WithValue(ctx, utils.JackalGoContextKey, srvCtx)

	var clientCtx client.Context

	if v := ctx.Value(client.ClientContextKey); v != nil {
		clientCtxPtr := v.(*client.Context)
		clientCtx = *clientCtxPtr
	}

	clientCtx = clientCtx.WithChainID("jackal-1")

	flagSet := pflag.NewFlagSet("jackalgo-flags", pflag.PanicOnError)

	AddTxFlagsToCmd(flagSet)

	flagSet.String(flags.FlagFrom, address, "Name or address of private key with which to sign")

	newpkey, err := pKey.Sign([]byte("Initiate Jackal Session"))
	if err != nil {
		return nil, err
	}

	eciesKey := ecies.NewPrivateKeyFromBytes(newpkey[:32])

	w := WalletHandler{
		clientCtx: clientCtx,
		flags:     flagSet,
		key:       pKey,
		address:   address,
		eciesKey:  eciesKey,
	}

	return &w, nil
}
