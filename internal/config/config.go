package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const fileName = ".goldapconfig.json"

type LdapAttribute struct {
	Name  string   `json:"Name"`
	Value []string `json:"Value"`
}

type Config struct {
	LdapUrl     string          `json:"LDAP_URL"`
	LdapBaseDn  string          `json:"LDAP_BASE_DN"`
	LdapAdminDn string          `json:"LDAP_ADMIN_DN"`
	Users       []LdapAttribute `json:"Users"`
	Groups      []LdapAttribute `json:"Groups"`
	OrgUnits    []LdapAttribute `json:"OrgUnits"`
}

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

func (c *Config) GetLdapAttributes(name string) ([]LdapAttribute, bool) {
	switch name {
	case "Users":
		return c.Users, true
	case "Groups":
		return c.Groups, true
	case "OrgUnits":
		return c.OrgUnits, true
	default:
		return []LdapAttribute{}, false
	}
}
