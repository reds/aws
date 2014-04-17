package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

type User struct {
	Name            string `json:"name"`
	AccessKeyId     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
}

type Account struct {
	Id              string          `json:"id"`
	AccessKeyId     string          `json:"accessKeyId"`
	SecretAccessKey string          `json:"secretAccessKey"`
	Users           []User          `json:"users,omitempty"`
	Route53ZoneIds  []Route53ZoneId `json:"route53ZoneIds,omitempty"`
}

type Route53ZoneId struct {
	Domain string `json:"domain"`
	Id     string `json:"id"`
}

type Config struct {
	Accounts []Account         `json:"accounts,omitempty"`
	Params   map[string]string `json:"params,omitempty"`
}

func LoadConfig(fn string) (*Config, error) {
	buf, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, errors.New("LoadConfig: " + err.Error())
	}
	var cfg Config
	err = json.Unmarshal(buf, &cfg)
	if err != nil {
		return nil, errors.New("LoadConfig: " + err.Error())
	}
	return &cfg, nil
}

func AccountConfig(fn, an string) (*Account, error) {
	cfg, err := LoadConfig(fn)
	if err != nil {
		return nil, err
	}
	return cfg.AccountConfig(an)
}

func (cfg *Config) AccountConfig(n string) (*Account, error) {
	for _, v := range cfg.Accounts {
		if n == v.Id {
			return &v, nil
		}
	}
	return nil, fmt.Errorf("account not found: " + n)
}

func UserConfig(fn, an, un string) (*User, error) {
	cfg, err := LoadConfig(fn)
	if err != nil {
		return nil, err
	}
	return cfg.UserConfig(an, un)
}

func (cfg *Config) UserConfig(an, un string) (*User, error) {
	for _, v := range cfg.Accounts {
		if an == v.Id {
			return v.UserConfig(un)
		}
	}
	return nil, fmt.Errorf("user not found: " + un + "," + an)
}

func (accts *Account) UserConfig(un string) (*User, error) {
	for _, v := range accts.Users {
		if un == v.Name {
			return &v, nil
		}
	}
	return nil, fmt.Errorf("user not found: " + un)
}

func Route53ZoneIdConfig(fn, acct, domain string) (*Route53ZoneId, error) {
	cfg, err := LoadConfig(fn)
	if err != nil {
		return nil, err
	}
	return cfg.Route53ZoneIdConfig(acct, domain)
}

func (cfg *Config) Route53ZoneIdConfig(an, domain string) (*Route53ZoneId, error) {
	acct, err := cfg.AccountConfig(an)
	if err != nil {
		return nil, err
	}
	return acct.Route53ZoneIdConfig(domain)
}

func (acct *Account) Route53ZoneIdConfig(domain string) (*Route53ZoneId, error) {
	for _, v := range acct.Route53ZoneIds {
		if domain == v.Domain {
			return &v, nil
		}
	}
	return nil, fmt.Errorf("domain not found: %s", domain)
}
