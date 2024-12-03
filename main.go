package main

import (
	"flag"
	"log"
	"os"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
	"gopkg.in/yaml.v3"
)

// CLI arguments
var (
	configPath = flag.String("config", "./vault-unseal.yaml", "The path to the configuration file")
)

var (
	l = log.New(os.Stdout, "", log.Ldate|log.Ltime)
)

type UnsealConfig struct {
	Nodes           []string `yaml:"nodes"`
	UnsealTokens    []string `yaml:"unsealTokens"`
	CheckInterval   int      `yaml:"checkInterval"`
	TlsSkipVerify   bool     `yaml:"tlsSkipVerify"`
	PrintUnsealLogs bool     `yaml:"printUnsealLogs"`
}

func checkVaultReady(client *vaultapi.Client) bool {
	_, err := client.Sys().Health()
	return err == nil
}

func getVaultSealStatus(client *vaultapi.Client) bool {
	status, err := client.Sys().SealStatus()
	if err != nil {
		l.Fatalf("Error checking seal status: %v", err)
	}

	return status.Sealed
}

func newVaultClient(addr string, tlsSkipVerify bool) (vaultClient *vaultapi.Client) {
	var err error
	vaultConfig := vaultapi.DefaultConfig()
	vaultConfig.Address = addr

	if err = vaultConfig.ConfigureTLS(&vaultapi.TLSConfig{Insecure: tlsSkipVerify}); err != nil {
		l.Fatalf("Error initializing tls config - %v", err)
	}

	if vaultClient, err = vaultapi.NewClient(vaultConfig); err != nil {
		l.Fatalf("Error creating vault client - %v", err)
	}

	return vaultClient
}

func main() {
	var unsealConfig UnsealConfig
	flag.Parse()

	yamlConig, err := os.ReadFile(*configPath)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlConig, &unsealConfig)
	if err != nil {
		panic(err)
	}

	if len(unsealConfig.UnsealTokens) == 0 {
		l.Fatalf("Error tokens not founds")
		os.Exit(1)
	}

	if len(unsealConfig.Nodes) == 0 {
		l.Fatalf("Error nodes not founds")
		os.Exit(1)
	}

	for {
		for _, node := range unsealConfig.Nodes {
			client := newVaultClient(node, unsealConfig.TlsSkipVerify)

			if !checkVaultReady(client) {
				l.Printf("Node ''%s' is not ready\n", node)
				continue
			}

			if getVaultSealStatus(client) {
				l.Printf("Node '%s' is seal, start unseal\n", node)
				for _, token := range unsealConfig.UnsealTokens {
					_, err := client.Sys().Unseal(token)
					if err != nil {
						l.Fatalf("Error unseal vault - %v", err)
					}
				}
			} else if unsealConfig.PrintUnsealLogs {
				l.Printf("Node '%s' is unseal\n", node)
			}
		}

		time.Sleep(time.Duration(unsealConfig.CheckInterval) * time.Second)
	}
}
