package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	LdapUrl     string `json:"LDAP_URL"`
	LdapBaseDn  string `json:"LDAP_BASE_DN"`
	LdapAdminDn string `json:"LDAP_ADMIN_DN"`
	Users       []struct {
		Name  string   `json:"Name"`
		Value []string `json:"Value"`
	} `json:"Users"`
	Groups []struct {
		Name  string   `json:"Name"`
		Value []string `json:"Value"`
	} `json:"Groups"`
	OrgUnits []struct {
		Name  string   `json:"Name"`
		Value []string `json:"Value"`
	} `json:"OrgUnits"`
}

const fileName = ".goldapconfig.json"

func getConfigPath() (string, error) {
	HomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Error getting config path: %v", err)
	}
	return HomeDir + "/" + fileName, nil
}

func write(c Config) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}
	jsonData, err := json.Marshal(c)
	if err != nil {
		return err
	}
	if err := os.WriteFile(configPath, jsonData, 0666); err != nil {
		return err
	}
	return nil
}

func Read() Config {
	var config Config
	configPath, err := getConfigPath()
	if err != nil {
		return config
	}
	contents, err := os.ReadFile(configPath)
	if err != nil {
		return config
	}
	if err := json.Unmarshal(contents, &config); err != nil {
		return Config{}
	}
	return config
}

func (c Config) SetUrl(ldapUrl string) error {
	c.LdapUrl = ldapUrl
	if err := write(c); err != nil {
		return fmt.Errorf("Error writing config; %v", err)
	}
	return nil
}

func (c Config) SetBaseDn(baseDn string) error {
	c.LdapBaseDn = baseDn
	if err := write(c); err != nil {
		return fmt.Errorf("Error writing config; %v", err)
	}
	return nil
}

func (c Config) SetlAdminDn(adminDn string) error {
	c.LdapAdminDn = adminDn
	if err := write(c); err != nil {
		return fmt.Errorf("Error writing config; %v", err)
	}
	return nil
}
