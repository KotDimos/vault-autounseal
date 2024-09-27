package main

import (
	"flag"
	"log"
	"os"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
	"gopkg.in/yaml.v2"
)

// CLI arguments
var (
	configPath = flag.String("config", "./vault-unseal.yaml", "The path to the configuration file")
)

var (
	l = log.New(os.Stdout, "", log.Ldate|log.Ltime)
)

type UnsealConfig struct {
	Nodes         []string      `yaml:"nodes"`
	UnsealTokens  []string      `yaml:"unsealTokens"`
	CheckInterval time.Duration `yaml:"checkInterval"`
	TlsSkipVerify bool          `yaml:"tlsSkipVerify"`
}

func checkVaultReady(client *vaultapi.Client) bool {
	_, err := client.Sys().Health()
	return err == nil
}

func getVaultSealStatus(client *vaultapi.Client) bool {
	status, err := client.Sys().SealStatus()
	if err != nil {
		l.Fatalf("checking seal status: %v", err)
	}

	return status.Sealed
}

func newVaultClient(addr string, tlsSkipVerify bool) (vaultClient *vaultapi.Client) {
	var err error
	vaultConfig := vaultapi.DefaultConfig()
	vaultConfig.Address = addr

	if err = vaultConfig.ConfigureTLS(&vaultapi.TLSConfig{Insecure: tlsSkipVerify}); err != nil {
		l.Fatalf("Error initializing tls config")
	}

	if vaultClient, err = vaultapi.NewClient(vaultConfig); err != nil {
		l.Fatalf("Error creating vault client: %v", err)
	}

	return vaultClient
}

func main() {
	var unsealConfig UnsealConfig
	flag.Parse()

	yamlFile, err := os.ReadFile(*configPath)
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
				l.Printf("Node %s is not ready\n", node)
				continue
			}

			if getVaultSealStatus(client) {
				l.Println("Vault is seal, start unseal")
				for _, token := range unsealConfig.UnsealTokens {
					_, err := client.Sys().Unseal(token)
					if err != nil {
						l.Fatalf("Error unseal vault: %v", err)
					}
				}
			} else {
				l.Println("Vault is unseal")
			}
		}

		time.Sleep(unsealConfig.CheckInterval * time.Second)
	}
}
