package database_test

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/tendermint/tendermint/version"

	stakingtypes "github.com/forbole/bdjuno/x/staking/types"

	"github.com/cosmos/cosmos-sdk/simapp"
	jconfig "github.com/desmos-labs/juno/config"
	"github.com/stretchr/testify/suite"

	"github.com/forbole/bdjuno/database"

	_ "github.com/proullon/ramsql/driver"
)

func TestDatabaseTestSuite(t *testing.T) {
	suite.Run(t, new(DbTestSuite))
}

type DbTestSuite struct {
	suite.Suite

	database *database.BigDipperDb
}

func (suite *DbTestSuite) SetupTest() {
	// Create the codec
	codec := simapp.MakeCodec()

	// Build the database
	config := &jconfig.Config{
		DatabaseConfig: jconfig.DatabaseConfig{
			Type: "psql",
			Config: &jconfig.PostgreSQLConfig{
				Name:     "bdjuno",
				Host:     "localhost",
				Port:     5433,
				User:     "bdjuno",
				Password: "password",
			},
		},
	}

	db, err := database.Builder(config, codec)
	suite.Require().NoError(err)

	bigDipperDb, ok := (db).(*database.BigDipperDb)
	suite.Require().True(ok)

	// Delete the public schema
	_, err = bigDipperDb.Sql.Exec(`DROP SCHEMA public CASCADE;`)
	suite.Require().NoError(err)

	// Re-create the schema
	_, err = bigDipperDb.Sql.Exec(`CREATE SCHEMA public;`)
	suite.Require().NoError(err)

	dirPath := "../schema"
	dir, err := ioutil.ReadDir(dirPath)
	for _, fileInfo := range dir {
		file, err := ioutil.ReadFile(filepath.Join(dirPath, fileInfo.Name()))
		suite.Require().NoError(err)

		commentsRegExp := regexp.MustCompile(`/\*.*\*/`)
		requests := strings.Split(string(file), ";")
		for _, request := range requests {
			_, err := bigDipperDb.Sql.Exec(commentsRegExp.ReplaceAllString(request, ""))
			suite.Require().NoError(err)
		}
	}

	suite.database = bigDipperDb
}

// getBlock builds, stores and returns a block for the provided height
func (suite *DbTestSuite) getBlock(height int64) *tmctypes.ResultBlock {
	validator := suite.getValidator(
		"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
		"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
		"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
	)

	addr, err := sdk.ConsAddressFromBech32(validator.GetConsAddr())
	suite.Require().NoError(err)

	block := &tmctypes.ResultBlock{
		BlockID: tmtypes.BlockID{},
		Block: &tmtypes.Block{
			Header: tmtypes.Header{
				Version:            version.Consensus{},
				ChainID:            "",
				Height:             height,
				Time:               time.Now(),
				LastBlockID:        tmtypes.BlockID{},
				LastCommitHash:     nil,
				DataHash:           nil,
				ValidatorsHash:     []byte("hash"),
				NextValidatorsHash: nil,
				ConsensusHash:      nil,
				AppHash:            nil,
				LastResultsHash:    nil,
				EvidenceHash:       nil,
				ProposerAddress:    tmtypes.Address(addr.Bytes()),
			},
			Data:     tmtypes.Data{},
			Evidence: tmtypes.EvidenceData{},
			LastCommit: &tmtypes.Commit{
				Height:     height - 1,
				Round:      1,
				BlockID:    tmtypes.BlockID{},
				Signatures: nil,
			},
		},
	}

	err = suite.database.SaveBlock(block, 10000, 100)
	suite.Require().NoError(err)
	return block
}

// getValidator stores inside the database a validator having the given
// consensus address, validator address and validator public key
func (suite *DbTestSuite) getValidator(consAddr, valAddr, pubkey string) stakingtypes.Validator {
	selfDelegation := suite.getAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")

	pubKey, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, pubkey)
	suite.Require().NoError(err)

	maxRate := sdk.NewDec(10)
	maxChangeRate := sdk.NewDec(20)

	validator := stakingtypes.NewValidator(consAddr, valAddr, pubKey, selfDelegation.String(), &maxChangeRate, &maxRate)
	err = suite.database.SaveValidatorData(validator)
	suite.Require().NoError(err)

	return validator
}

// getAccount saves inside the database an account having the given address
func (suite *DbTestSuite) getAccount(addr string) sdk.AccAddress {
	delegator, err := sdk.AccAddressFromBech32(addr)
	suite.Require().NoError(err)

	_, err = suite.database.Sql.Exec(`INSERT INTO account (address) VALUES ($1) ON CONFLICT DO NOTHING`, delegator.String())
	suite.Require().NoError(err)

	return delegator
}
