package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/creasty/defaults"
	"gopkg.in/yaml.v3"
)

type UnsealConfig struct {
	Nodes           []string `yaml:"nodes"`
	UnsealTokens    []string `yaml:"unsealTokens"`
	CheckInterval   int      `default:"15"        yaml:"checkInterval,omitempty"`
	TLSSkipVerify   bool     `default:"true"      yaml:"tlsSkipVerify,omitempty"`
	PrintUnsealLogs bool     `default:"false"     yaml:"printUnsealLogs,omitempty"`
}

func parseUnsealConfig() UnsealConfig {
	unsealConfig := &UnsealConfig{}
	if err := defaults.Set(unsealConfig); err != nil {
		panic(err)
	}

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
		log.Fatalf("Error tokens not founds in %s", *configPath)
	}

	if len(unsealConfig.Nodes) == 0 {
		log.Fatalf("Error nodes not founds in %s", *configPath)
	}

	return *unsealConfig
}

func main() {
	unsealConfig := parseUnsealConfig()

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
