package wallet_handler

import "github.com/JackalLabs/jackalgo/tx"

func (w *WalletHandler) SendTx(msg types.Msg) (*types.TxResponse, error) {

	res, err := tx.SendTx(w.clientCtx, nil, msg)

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
