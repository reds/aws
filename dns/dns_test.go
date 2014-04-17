package dns

import (
	"github.com/reds/config"
	"testing"
)

func TestGetip(t *testing.T) {
	acct, err := config.AccountConfig(filepath.Join(os.Getenv("HOME"), ".ssh", "amazon.js"),
		"xamazonpt@redmond5.com")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(GetIpForDomain("flashfb.com"))
}
