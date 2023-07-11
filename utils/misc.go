package utils

import (
	"context"
	"fmt"

	"github.com/JackalLabs/jackalgo/handlers/wallet_handler"
	"github.com/JackalLabs/jackalgo/utils/crypt"
	filetreetypes "github.com/jackalLabs/canine-chain/v3/x/filetree/types"
)

func NumTo3xTB(base int64) int64 {
	base *= 1000 /** KB */
	base *= 1000 /** MB */
	base *= 1000 /** GB */
	base *= 1000 /** TB */
	base *= 3    /** Redundancy */
	return base
}

func GetFileTreeData(rawPath string, owner string, wallet *wallet_handler.WalletHandler) (*filetreetypes.QueryFileResponse, error) {
	hexAddress := crypt.MerkleMeBro(rawPath)
	hexedOwner := crypt.HashAndHex(fmt.Sprintf(`o%s%s`, hexAddress, crypt.HashAndHex(owner)))

	queryClient := filetreetypes.NewQueryClient(wallet.GetClientCtx())

	req := filetreetypes.QueryFileRequest{
		Address:      hexAddress,
		OwnerAddress: hexedOwner,
	}

	return queryClient.Files(context.Background(), &req)
}
