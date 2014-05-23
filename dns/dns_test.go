package dns

import (
	"github.com/reds/aws/config"
	"os"
	"path/filepath"
	"testing"
)

func TestGetip(t *testing.T) {
	acct, err := config.AccountConfig(filepath.Join(os.Getenv("HOME"), ".ssh", "amazon.js"),
		"xamazonpt@redmond5.com")
	if err != nil {
		t.Fatal(err)
	}
	rt53, err := NewFromAccount(acct, "flashfb.com")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(rt53.GetIpForDomain("flashfb.com"))
}
