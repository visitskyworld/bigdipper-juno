package main

import (
	"github.com/desmos-labs/juno/cmd"
	junoparse "github.com/desmos-labs/juno/cmd/parse"

	"github.com/forbole/bdjuno/x/messages"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x"
)

func main() {
	// Setup the config
	config := junoparse.NewConfig("bdjuno").
		WithRegistrar(x.NewModulesRegistrar(messages.AddressesParser)).
		WithDBBuilder(database.Builder)

	// Run the command
	executor := cmd.BuildDefaultExecutor(config)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}
