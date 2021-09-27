package remote

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/desmos-labs/juno/node/remote"
	bankkeeper "github.com/forbole/bdjuno/modules/bank/source"
	"github.com/forbole/bdjuno/types"
)

var (
	_ bankkeeper.Source = &Source{}
)

type Source struct {
	*remote.Source
	bankClient banktypes.QueryClient
}

// NewSource builds a new Source instance
func NewSource(source *remote.Source, bankClient banktypes.QueryClient) *Source {
	return &Source{
		Source:     source,
		bankClient: bankClient,
	}
}

// GetBalances implements bankkeeper.Source
func (k Source) GetBalances(addresses []string, height int64) ([]types.AccountBalance, error) {
	header := remote.GetHeightRequestHeader(height)

	var balances []types.AccountBalance
	for _, address := range addresses {
		balRes, err := k.bankClient.AllBalances(k.Ctx, &banktypes.QueryAllBalancesRequest{Address: address}, header)
		if err != nil {
			return nil, fmt.Errorf("error while getting all balances: %s", err)
		}

		balances = append(balances, types.NewAccountBalance(
			address,
			balRes.Balances,
			height,
		))
	}

	return balances, nil
}

// GetSupply implements bankkeeper.Source
func (k Source) GetSupply(height int64) (sdk.Coins, error) {
	header := remote.GetHeightRequestHeader(height)
	res, err := k.bankClient.TotalSupply(k.Ctx, &banktypes.QueryTotalSupplyRequest{}, header)
	if err != nil {
		return nil, fmt.Errorf("error while getting total supply: %s", err)
	}

	return res.Supply, nil
}
