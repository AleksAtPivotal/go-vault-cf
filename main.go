package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cloudfoundry-community/go-cfenv"
	vault "github.com/mch1307/vaultlib"
)

func main() {
	log.Printf("Started Application...\n")
	http.HandleFunc("/", handle)
	log.Printf("Started Listening on port: %s \n", os.Getenv("PORT"))
	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil)

}

func handle(w http.ResponseWriter, r *http.Request) {
	appEnv, err := cfenv.Current()
	if err != nil {
		fmt.Printf("Error getting CF Service info : %s\n", err)
		return
	}

	fmt.Fprintf(w, "vcap_services is %s\n", os.Getenv("VCAP_SERVICES"))

	// VCAP_Services as go-cfenv Services
	cfenvvault, err := appEnv.Services.WithName("hashicorp-vault")
	if err != nil {
		fmt.Printf("Error getting Vault Service from VCAP : %s\n", err)
		return
	}

	var vaultSB VaultService
	vaultSB.name = cfenvvault.Name
	vaultSB.label = cfenvvault.Label

	for k, v := range cfenvvault.Credentials {
		if k == "address" {
			vaultSB.credentials.address = fmt.Sprintf("%s", v)
		}
		if k == "auth" {
			// We need to parse auth to token
			vaultSB = vaultSB.getTokenfromCF(v)
		}
		if k == "backends" {
			vaultSB = vaultSB.getBackendsfromCF(v)
		}
		if k == "backends_shared" {
			vaultSB = vaultSB.getBackendsSharedfromCF(v)
		}
	}

	log.Printf("Parsed VaultSB value %s: \n", vaultSB)
	log.Printf("Calling Vault API ...\n")

	vcConf := vault.NewConfig()
	vcConf.Address = vaultSB.credentials.address
	vcConf.InsecureSSL = true
	vcConf.Token = vaultSB.credentials.auth.token

	vaultCli, err := vault.NewClient(vcConf)
	if err != nil {
		log.Fatal(err)
	}

	var secretPath string
	secretPath = os.Getenv("VAULT_SEARCH_PATH")
	if secretPath == "" {
		secretPath = vaultSB.credentials.backends.generic
	}

	// Get the Vault secret data
	secret, err := vaultCli.GetSecret(secretPath)
	if err != nil {
		log.Println(err)
	}

	for k, v := range secret.KV {
		log.Printf("secret %v: %v\n", k, v)
	}

}
