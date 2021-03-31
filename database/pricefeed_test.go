package database_test

import (
	"fmt"
	"time"

	pricefeedtypes "github.com/forbole/bdjuno/x/pricefeed/types"

	dbtypes "github.com/forbole/bdjuno/database/types"
)

func (suite *DbTestSuite) insertToken(name string) {
	query := fmt.Sprintf(
		`INSERT INTO token (name, traded_unit) VALUES ('%[2]s', '%[1]s')`, name, name)
	_, err := suite.database.Sql.Query(query)
	suite.Require().NoError(err)

	query = fmt.Sprintf(
		`INSERT INTO token_unit (token_name, denom, exponent) VALUES ('%[1]s', 'u%[1]s', 0), ('%[1]s', 'm%[1]s', 3), ('%[1]s', '%[1]s', 6)`,
		name)
	_, err = suite.database.Sql.Query(query)
	suite.Require().NoError(err)
}

func (suite *DbTestSuite) Test_GetTradedNames() {
	suite.insertToken("desmos")
	suite.insertToken("daric")

	tradedNames, err := suite.database.GetTradedNames()
	suite.Require().NoError(err)

	var expected = []string{"desmos", "daric"}
	suite.Require().Len(tradedNames, len(expected))
	for _, name := range expected {
		suite.Require().Contains(tradedNames, name)
	}
}

func (suite *DbTestSuite) TestBigDipperDb_SaveTokenPrice() {
	suite.insertToken("desmos")
	suite.insertToken("atom")

	tickers := pricefeedtypes.MarketTickers{
		pricefeedtypes.NewMarketTicker(
			"desmos",
			100.01,
			10,
			time.Date(2020, 10, 10, 15, 00, 00, 000, time.UTC),
		),
		pricefeedtypes.NewMarketTicker(
			"atom",
			200.01,
			20,
			time.Date(2020, 10, 10, 15, 00, 00, 000, time.UTC),
		),
	}

	err := suite.database.SaveTokensPrices(tickers)
	suite.Require().NoError(err)

	expected := []dbtypes.TokenPriceRow{
		dbtypes.NewTokenPriceRow("desmos",
			100.01,
			10,
			time.Date(2020, 10, 10, 15, 00, 00, 000, time.UTC)),
		dbtypes.NewTokenPriceRow("atom",
			200.01,
			20,
			time.Date(2020, 10, 10, 15, 00, 00, 000, time.UTC),
		),
	}
	var rows []dbtypes.TokenPriceRow
	err = suite.database.Sqlx.Select(&rows, `SELECT name, price, market_cap, timestamp FROM token_price`)
	suite.Require().NoError(err)
	for i, row := range rows {
		suite.Require().True(expected[i].Equals(row))
	}
}
