package main

import (
	"flag"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type UnsealConfig struct {
	Nodes           []string `yaml:"nodes"`
	UnsealTokens    []string `yaml:"unsealTokens"`
	CheckInterval   int      `yaml:"checkInterval"`
	TLSSkipVerify   bool     `yaml:"tlsSkipVerify"`
	PrintUnsealLogs bool     `yaml:"printUnsealLogs"`
}

func main() {
	var unsealConfig UnsealConfig

	configPath := flag.String("config", "./vault-unseal.yaml", "The path to the configuration file")
	flag.Parse()

	yamlConig, err := os.ReadFile(*configPath)
	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(yamlConig, &unsealConfig); err != nil {
		panic(err)
	}

	if len(unsealConfig.UnsealTokens) == 0 {
		log.Fatalf("Error tokens not founds")
	}

	if len(unsealConfig.Nodes) == 0 {
		log.Fatalf("Error nodes not founds")
	}

	for {
		for _, node := range unsealConfig.Nodes {
			client := newVaultClient(node, unsealConfig.TLSSkipVerify)

			if !checkVaultReady(client) {
				log.Printf("Node '%s' is not ready\n", node)

				continue
			}

			if getVaultSealStatus(client) {
				log.Printf("Node '%s' is seal, start unseal\n", node)

				for _, token := range unsealConfig.UnsealTokens {
					if _, err := client.Sys().Unseal(token); err != nil {
						log.Fatalf("Error unseal vault - %v", err)
					}
				}
			} else if unsealConfig.PrintUnsealLogs {
				log.Printf("Node '%s' is unseal\n", node)
			}
		}

		time.Sleep(time.Duration(unsealConfig.CheckInterval) * time.Second)
	}
}
