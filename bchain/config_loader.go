//go:build integration

package bchain

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	buildcfg "github.com/trezor/blockbook/build/tools"
)

// BlockchainCfg contains fields read from blockbook's blockchaincfg.json after being rendered from templates.
type BlockchainCfg struct {
	// more fields can be added later as needed
	RpcUrl string `json:"rpc_url"`
}

// LoadBlockchainCfg returns the resolved blockchaincfg.json (env overrides are honored in tests)
func LoadBlockchainCfg(t *testing.T, coinAlias string) BlockchainCfg {
	t.Helper()

	configsDir, err := repoConfigsDir()
	if err != nil {
		t.Fatalf("integration config path error: %v", err)
	}
	templatesDir, err := repoTemplatesDir(configsDir)
	if err != nil {
		t.Fatalf("integration templates path error: %v", err)
	}

	config, err := buildcfg.LoadConfig(configsDir, coinAlias)
	if err != nil {
		t.Fatalf("load config for %s: %v", coinAlias, err)
	}

	outputDir, err := os.MkdirTemp("", "integration_blockchaincfg")
	if err != nil {
		t.Fatalf("integration temp dir error: %v", err)
	}
	t.Cleanup(func() {
		_ = os.RemoveAll(outputDir)
	})

	// Render templates so tests read the same generated blockchaincfg.json as packaging.
	if err := buildcfg.GeneratePackageDefinitions(config, templatesDir, outputDir); err != nil {
		t.Fatalf("generate package definitions for %s: %v", coinAlias, err)
	}

	blockchainCfg, err := readBlockchainCfg(filepath.Join(outputDir, "blockbook", "blockchaincfg.json"))
	if err != nil {
		t.Fatalf("read blockchain config for %s: %v", coinAlias, err)
	}
	if blockchainCfg.RpcUrl == "" {
		t.Fatalf("empty rpc_url for %s", coinAlias)
	}
	return blockchainCfg
}

// readBlockchainCfg loads the rendered blockchain config for test assertions.
func readBlockchainCfg(path string) (BlockchainCfg, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return BlockchainCfg{}, err
	}
	var cfg BlockchainCfg
	if err := json.Unmarshal(b, &cfg); err != nil {
		return BlockchainCfg{}, err
	}
	return cfg, nil
}

// repoTemplatesDir locates build/templates relative to the repo root.
func repoTemplatesDir(configsDir string) (string, error) {
	repoRoot := filepath.Dir(configsDir)
	templatesDir := filepath.Join(repoRoot, "build", "templates")
	if _, err := os.Stat(templatesDir); err == nil {
		return templatesDir, nil
	} else if os.IsNotExist(err) {
		return "", fmt.Errorf("build/templates not found near %s", configsDir)
	} else {
		return "", err
	}
}

// repoConfigsDir finds configs/coins from the caller path so tests can run from any subdir.
func repoConfigsDir() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("unable to resolve caller path")
	}
	dir := filepath.Dir(file)
	// Walk up so tests can run from any subdir while still locating configs.
	for i := 0; i < 3; i++ {
		configsDir := filepath.Join(dir, "configs")
		if _, err := os.Stat(filepath.Join(configsDir, "coins")); err == nil {
			return configsDir, nil
		} else if !os.IsNotExist(err) {
			return "", err
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", errors.New("configs/coins not found from caller path")
}
