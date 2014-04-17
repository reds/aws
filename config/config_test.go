package config

import (
	"github.com/reds/aws/dns"
	"os"
	"path/filepath"
	"testing"
)

func TestConfig(t *testing.T) {
	cfg, err := LoadConfig(filepath.Join(os.Getenv("HOME"), ".ssh", "amazon.js"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(cfg)
	acct, err := cfg.AccountConfig("xamazonpt@redmond5.com")
	if err != nil {
		t.Fatal(err)
	}
	rt53, err := acct.Route53ZoneIdConfig("flashfb.com")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(dns.GetIpForDomain("flashfb.com", acct.AccessKeyId, acct.SecretAccessKey, rt53.Id))
}
