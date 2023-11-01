package wallet_handler

import (
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

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
	key       types.PrivKey
	eciesKey  *ecies.PrivateKey
}

func DefaultWalletHandler(seedPhrase string) (*WalletHandler, error) {
	return NewWalletHandler(seedPhrase, "https://rpc.jackalprotocol.com:443", "jackal-1")
}

func createFlags(gas string, address string) *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("jackalgo-flags", pflag.PanicOnError)

	AddTxFlags(flagSet)
	// --gas can accept integers and "auto"
	flagSet.String(flags.FlagGas, gas, fmt.Sprintf("gas limit to set per-transaction; set to %q to calculate sufficient gas automatically (default %d)", flags.GasFlagAuto, flags.DefaultGasLimit))

	flagSet.String(flags.FlagFrom, address, "Name or address of private key with which to sign")

	return flagSet
}

func NewWalletHandler(seedPhrase string, rpc string, chainId string) (*WalletHandler, error) {
	encodingConfig := params.MakeTestEncodingConfig()
	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	simapp.ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	simapp.ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	cfg.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	cfg.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)
	// cfg.SetAddressVerifier(wasmtypes.VerifyAddressLen())
	cfg.Seal()

	hdPath := hd.CreateHDPath(118, 0, 0).String()
	// create master key and derive first key for keyring
	derivedPriv, err := hd.Secp256k1.Derive()(seedPhrase, "", hdPath)
	if err != nil {
		return nil, err
	}

	pKey := hd.Secp256k1.Generate()(derivedPriv)

	// check if the key already exists with the same address and return an error
	// if found
	address := sdk.AccAddress(pKey.PubKey().Address())

	cl, err := client.NewClientFromNode(rpc)
	if err != nil {
		return nil, err
	}

	clientCtx := client.Context{}.
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastBlock).
		WithViper("").
		WithNodeURI(rpc).
		WithClient(cl).
		WithChainID(chainId)

	newpkey, err := pKey.Sign([]byte("Initiate Jackal Session"))
	if err != nil {
		return nil, err
	}

	eciesKey := ecies.NewPrivateKeyFromBytes(newpkey[:32])

	flgs := createFlags("auto", address.String())
	flgs.Float64(flags.FlagGasAdjustment, 1.5, fmt.Sprintf("the gas adjustment, the default is %f", 1.5))

	w := WalletHandler{
		clientCtx: clientCtx,
		flags:     flgs,
		key:       pKey,
		address:   address.String(),
		eciesKey:  eciesKey,
	}

	return &w, nil
}

func (w *WalletHandler) WithGas(gas string) *WalletHandler {
	flgs := createFlags(gas, w.address)

	newWallet := WalletHandler{
		clientCtx: w.clientCtx,
		flags:     flgs,
		key:       w.key,
		address:   w.address,
		eciesKey:  w.eciesKey,
	}
	return &newWallet
}
