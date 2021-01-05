package override_test

import (
	"fmt"

	"github.com/engage-ehs/override"
)

type DeploymentConfig struct {
	Environment string
	Geography   string

	Author string `canset:"no"`
}

func Example() {
	// arguments from the command-line, such as returned from flag.Args
	cl := []string{"environment=prod", "geography=eu"}

	// some settings can be collected from the environment, and we do not want them to be
	// changed by the command-line arguments
	whoami := func() string { return "current user" }

	cfg := DeploymentConfig{
		Environment: "staging", // ship to staging by default
		Author:      whoami(),
	}
	if err := override.Scan(cl, &cfg); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(cfg)
	// Output: {prod eu current user}
}
