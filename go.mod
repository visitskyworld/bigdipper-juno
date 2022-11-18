module github.com/forbole/bdjuno/v2

go 1.16

require (
	github.com/cosmos/cosmos-sdk v0.46.5
	github.com/forbole/juno/v2 v2.0.0-20220117075314-1d0a50fab7d4
	github.com/go-co-op/gocron v1.11.0
	github.com/gogo/protobuf v1.3.3
	github.com/jmoiron/sqlx v1.2.1-0.20200324155115-ee514944af4b
	github.com/lib/pq v1.10.6
	github.com/pelletier/go-toml v1.9.5
	github.com/proullon/ramsql v0.0.0-20181213202341-817cee58a244
	github.com/rs/zerolog v1.27.0
	github.com/spf13/cobra v1.6.0
	github.com/stretchr/testify v1.8.0
	github.com/tendermint/tendermint v0.34.23
	google.golang.org/grpc v1.50.1
	gopkg.in/yaml.v3 v3.0.1
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

replace github.com/tendermint/tendermint => github.com/forbole/tendermint v0.34.13-0.20210820072129-a2a4af55563d
