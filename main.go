package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/cloudfoundry-community/go-cfenv"
	vault "github.com/mch1307/vaultlib"
)

func main() {
	http.HandleFunc("/", handle)
	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil)
}

func handle(w http.ResponseWriter, r *http.Request) {
	appEnv, _ := cfenv.Current()
	fmt.Println("Services:", appEnv.Services)
	fmt.Fprintf(w, "vcap_services is %s\n", os.Getenv("VCAP_SERVICES"))
	vcConf := vault.NewConfig()
	cfvault, err := appEnv.Services.WithName("hashicorp-vault")
	if err != nil {
		fmt.Println("Error getting Vault Service from VCAP : %s", err)
	}

	vcConf.Address = cfvault.Credentials

}
