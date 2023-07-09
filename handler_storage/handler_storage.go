package handler_storage

import (
	"context"
	"fmt"

	"github.com/JackalLabs/jackalgo/handler_wallet"
	"github.com/JackalLabs/jackalgo/utils"
	"github.com/cosmos/cosmos-sdk/types"
	storagetypes "github.com/jackalLabs/canine-chain/v3/x/storage/types"
)

func NewStorageHandler(w *handler_wallet.WalletHandler) *StorageHandler {

	s := StorageHandler{
		walletHandler: w,
	}

	return &s
}

func (s *StorageHandler) BuyStorage(forAddress string, duration int64, bytes int64) (*types.TxResponse, error) {
	if duration <= 0 {
		return nil, fmt.Errorf("cannot use less than 0 months of duration")
	}
	monthsAsHours := duration * 720

	_, err := types.AccAddressFromBech32(forAddress)
	if err != nil {
		return nil, err
	}

	buyMsg := storagetypes.MsgBuyStorage{
		Creator:      s.walletHandler.GetAddress(),
		PaymentDenom: "ujkl",
		ForAddress:   forAddress,
		Duration:     fmt.Sprintf("%dh", monthsAsHours),
		Bytes:        fmt.Sprintf("%d", utils.NumTo3xTB(bytes)),
	}
	res, err := s.walletHandler.SendTx(&buyMsg)

	return res, err
}

func (s *StorageHandler) UpgradeStorage(forAddress string, duration int64, bytes int64) (*types.TxResponse, error) {
	if duration <= 0 {
		return nil, fmt.Errorf("cannot use less than 0 months of duration")
	}
	monthsAsHours := duration * 720

	_, err := types.AccAddressFromBech32(forAddress)
	if err != nil {
		return nil, err
	}

	buyMsg := storagetypes.MsgUpgradeStorage{
		Creator:      s.walletHandler.GetAddress(),
		PaymentDenom: "ujkl",
		ForAddress:   forAddress,
		Duration:     fmt.Sprintf("%dh", monthsAsHours),
		Bytes:        fmt.Sprintf("%d", utils.NumTo3xTB(bytes)),
	}
	res, err := s.walletHandler.SendTx(&buyMsg)

	return res, err
}

func (s *StorageHandler) QueryGetPayData(address string) (*storagetypes.QueryPayDataResponse, error) {
	_, err := types.AccAddressFromBech32(address)
	if err != nil {
		return nil, err
	}

	req := storagetypes.QueryPayDataRequest{
		Address: address,
	}

	queryClient := storagetypes.NewQueryClient(s.walletHandler.GetClientCtx())

	res, err := queryClient.GetPayData(context.Background(), &req)

	return res, err
}

func (s *StorageHandler) QueryJackalPrice(bytes int64, duration int64) (*storagetypes.QueryPriceCheckResponse, error) {
	tbs := utils.NumTo3xTB(bytes)
	if duration <= 0 {
		return nil, fmt.Errorf("cannot use less than 0 months of duration")
	}
	monthsAsHours := duration * 720

	req := storagetypes.QueryPriceCheckRequest{
		Bytes:    tbs,
		Duration: fmt.Sprintf("%dh", monthsAsHours),
	}

	queryClient := storagetypes.NewQueryClient(s.walletHandler.GetClientCtx())

	res, err := queryClient.PriceCheck(context.Background(), &req)

	return res, err
}
