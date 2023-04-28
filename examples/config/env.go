package config

import (
	"fmt"
	"os"
	"strings"
)

type Example struct {
	Env           Env
	RunSlowStream bool
	RunTrades     bool
	RunPerpTrades bool
}

func Load() (Example, error) {
	env, err := loadEnv()
	if err != nil {
		return Example{}, err
	}

	runSlowStream := true
	rtsV := os.Getenv("RUN_SLOW_STREAM")
	if rtsV == "false" {
		runSlowStream = false
	}

	runTrades := true
	rtV := os.Getenv("RUN_TRADES")
	if rtV == "false" {
		runTrades = false
	}

	runPerpTrades := true
	rtP := os.Getenv("RUN_PERP_TRADES")
	if rtP == "false" {
		runPerpTrades = false
	}

	return Example{
		Env:           env,
		RunTrades:     runTrades,
		RunPerpTrades: runPerpTrades,
		RunSlowStream: runSlowStream,
	}, nil
}

type Env string

const (
	EnvMainnet Env = "mainnet"
	EnvTestnet Env = "testnet"
	EnvLocal   Env = "local"
)

func loadEnv() (Env, error) {
	v, ok := os.LookupEnv("API_ENV")
	if !ok {
		return EnvMainnet, nil
	}

	switch Env(strings.ToLower(v)) {
	case EnvLocal:
		return EnvLocal, nil
	case EnvTestnet:
		return EnvTestnet, nil
	case EnvMainnet:
		return EnvMainnet, nil
	default:
		return EnvMainnet, fmt.Errorf("API_ENV %v not supported", v)
	}
}
