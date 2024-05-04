package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"time"

	vault "github.com/hashicorp/vault/api"
	"gopkg.in/yaml.v2"
)

// CLI arguments
var (
	configPath = flag.String("config", "./vault-unseal.yaml", "The path to the configuration file")
)

type UnsealConfig struct {
	Nodes         []string      `yaml:"nodes"`
	Tokens        []string      `yaml:"tokens"`
	CheckInterval time.Duration `yaml:"checkInterval"`
	TlsSkipVerify bool          `yaml:"tlsSkipVerify"`
}

func checkVaultReady(client *vault.Client) bool {
	_, err := client.Sys().Health()
	if err != nil {
		return false
	}
	return true
}

func getVaultSealStatus(client *vault.Client) bool {
	status, err := client.Sys().SealStatus()
	if err != nil {
		fmt.Errorf("checking seal status: %w", err)
	}

	return status.Sealed
}

func newVaultClient(addr string, tlsSkipVerify bool) (vaultClient *vault.Client) {
	var err error
	vaultConfig := vault.DefaultConfig()
	vaultConfig.Address = addr

	if err = vaultConfig.ConfigureTLS(&vault.TLSConfig{Insecure: tlsSkipVerify}); err != nil {
		fmt.Errorf("error initializing tls config")
	}

	if vaultClient, err = vault.NewClient(vaultConfig); err != nil {
		fmt.Errorf("error creating vault client: %v", err)
	}

	return vaultClient
}

func main() {
	var unsealConfig UnsealConfig
	flag.Parse()

	yamlFile, err := ioutil.ReadFile(*configPath)
	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(yamlFile, &unsealConfig); err != nil {
		panic(err)
	}

	for {
		for _, node := range unsealConfig.Nodes {
			client := newVaultClient(node, unsealConfig.TlsSkipVerify)

			if !checkVaultReady(client) {
				fmt.Printf("Node %s is not ready\n", node)
				continue
			}

			if getVaultSealStatus(client) {
				fmt.Println("Vault is seal, start unseal")
				for _, token := range unsealConfig.Tokens {
					client.Sys().Unseal(token)
				}
			} else {
				fmt.Println("Vault is unseal")
			}
		}

		time.Sleep(unsealConfig.CheckInterval * time.Second)
	}
}
