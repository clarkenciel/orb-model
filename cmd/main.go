package main

import (
	orb "github.com/clarkenciel/orb-model"
)

func main() {
	ids := []orb.PerformerId{
		"tom",
		"dick",
		"harry",
	}

	script := orb.Script{
		"one",
		"two",
		"three",
	}

	router := orb.PerformerRouter{
		"tom": orb.AddressSet{
			orb.Address{"dick", orb.Left}: true,
			orb.Address{"harry", orb.Right}: true,
		},
		"dick": orb.AddressSet{
			orb.Address{"harry", orb.Left}: true,
			orb.Address{"tom", orb.Right}: true,
		},
		"harry": orb.AddressSet{
			orb.Address{"tom", orb.Left}: true,
			orb.Address{"dick", orb.Right}: true,
		},
	}

	
}
