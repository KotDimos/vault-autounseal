package main

import (
	"log"

	vaultapi "github.com/hashicorp/vault/api"
)

func checkVaultReady(client *vaultapi.Client) bool {
	_, err := client.Sys().Health()

	return err == nil
}

func getVaultSealStatus(client *vaultapi.Client) bool {
	status, err := client.Sys().SealStatus()
	if err != nil {
		log.Fatalf("Error checking seal status: %v", err)
	}

	return status.Sealed
}

func newVaultClient(addr string, tlsSkipVerify bool) (vaultClient *vaultapi.Client) {
	var err error

	vaultConfig := vaultapi.DefaultConfig()
	vaultConfig.Address = addr

	if err = vaultConfig.ConfigureTLS(&vaultapi.TLSConfig{Insecure: tlsSkipVerify}); err != nil {
		log.Fatalf("Error initializing tls config - %v", err)
	}

	if vaultClient, err = vaultapi.NewClient(vaultConfig); err != nil {
		log.Fatalf("Error creating vault client - %v", err)
	}

	return vaultClient
}
