package wallet_handler

import (
	"fmt"

	"github.com/JackalLabs/jackalgo/tx"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/spf13/pflag"
)

// AddTxFlagsToCmd adds common flags to a module tx command.
func AddTxFlagsToCmd(set *pflag.FlagSet) {
	set.Uint64P(flags.FlagAccountNumber, "a", 0, "The account number of the signing account (offline mode only)")
	set.Uint64P(flags.FlagSequence, "s", 0, "The sequence number of the signing account (offline mode only)")
	set.String(flags.FlagNote, "", "Note to add a description to the transaction (previously --memo)")
	set.String(flags.FlagFees, "", "Fees to pay along with transaction; eg: 10uatom")
	set.String(flags.FlagGasPrices, "", "Gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)")
	set.String(flags.FlagNode, "tcp://localhost:26657", "<host>:<port> to tendermint rpc interface for this chain")
	set.Bool(flags.FlagUseLedger, false, "Use a connected Ledger device")
	set.Float64(flags.FlagGasAdjustment, flags.DefaultGasAdjustment, "adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored ")
	set.StringP(flags.FlagBroadcastMode, "b", flags.BroadcastBlock, "Transaction broadcasting mode (sync|async|block)")
	set.Bool(flags.FlagDryRun, false, "ignore the --gas flag and perform a simulation of a transaction, but don't broadcast it (when enabled, the local Keybase is not accessible)")
	set.Bool(flags.FlagGenerateOnly, false, "Build an unsigned transaction and write it to STDOUT (when enabled, the local Keybase is not accessible)")
	set.Bool(flags.FlagOffline, false, "Offline mode (does not allow any online functionality")
	set.BoolP(flags.FlagSkipConfirmation, "y", false, "Skip tx broadcasting prompt confirmation")
	set.String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|kwallet|pass|test|memory)")
	set.String(flags.FlagSignMode, "direct", "Choose sign mode (direct|amino-json), this is an advanced feature")
	set.Uint64(flags.FlagTimeoutHeight, 0, "Set a block timeout height to prevent the tx from being committed past a certain height")
	set.String(flags.FlagFeeAccount, "", "Fee account pays fees for the transaction instead of deducting from the signer")

	// --gas can accept integers and "auto"
	set.String(flags.FlagGas, "", fmt.Sprintf("gas limit to set per-transaction; set to %q to calculate sufficient gas automatically (default %d)", flags.GasFlagAuto, flags.DefaultGasLimit))
}

func (w *WalletHandler) SendTx(msg types.Msg) (*types.TxResponse, error) {

	res, err := tx.SendTx(w.clientCtx, w.flags, msg)

	return res, err
}

func (w *WalletHandler) SendTokens(toAddress string, amount types.Coins) (*types.TxResponse, error) {

	sendMsg := banktypes.MsgSend{
		FromAddress: w.address,
		ToAddress:   toAddress,
		Amount:      amount,
	}

	res, err := w.SendTx(&sendMsg)

	return res, err
}
