//go:build integration

package bchain

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	buildcfg "github.com/trezor/blockbook/build/tools"
)

// RPCURLFromConfig renders ipc.rpc_url_template from the coin config for integration tests.
func RPCURLFromConfig(t *testing.T, coinAlias string) string {
	t.Helper()
	configsDir, err := repoConfigsDir()
	if err != nil {
		t.Fatalf("integration config path error: %v", err)
	}
	cfg, err := buildcfg.LoadConfig(configsDir, coinAlias)
	if err != nil {
		t.Fatalf("load config for %s: %v", coinAlias, err)
	}
	templ := cfg.ParseTemplate()
	var out bytes.Buffer
	if err := templ.ExecuteTemplate(&out, "IPC.RPCURLTemplate", cfg); err != nil {
		t.Fatalf("render rpc_url_template for %s: %v", coinAlias, err)
	}
	rpcURL := strings.TrimSpace(out.String())
	if rpcURL == "" {
		t.Fatalf("empty rpc url from config for %s", coinAlias)
	}
	return rpcURL
}

func repoConfigsDir() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("unable to resolve caller path")
	}
	dir := filepath.Dir(file)
	// search the config directory in the parent folders so it is agnostic to the caller location
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
