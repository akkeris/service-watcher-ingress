
package utils

import (
	"net/http"
	"os"
	vault "github.com/akkeris/vault-client"
	"strings"
	"fmt"
	"k8s.io/client-go/rest"
)

var ExternalCert string
var ExternalCa string
var ExternalKey string
var InternalCert string
var InternalCa string
var InternalKey string
var Kubetoken string
var HTTPClient *http.Client
var Regionapilocation string
var Regionapiusername string
var Regionapipassword string
var Ingresstlssecret string
var DefaultDomain string
var InsideDomain string
var Kubernetesapiurl string
var Blacklist map[string]bool
var Client rest.Interface


func SetSecrets() {
	externalsecret := vault.GetSecret(os.Getenv("EXTERNAL_CERT_VAULT_PATH"))
	ExternalCert = strings.Replace(vault.GetFieldFromVaultSecret(externalsecret, "cert"), "\\n", "\n", -1)
	ExternalCert = strings.Replace(ExternalCert, "\\r", "\r", -1)
	ExternalCa = strings.Replace(vault.GetFieldFromVaultSecret(externalsecret, "ca"), "\\n", "\n", -1)
	ExternalCa = strings.Replace(ExternalCa, "\\r", "\r", -1)
	ExternalKey = strings.Replace(vault.GetFieldFromVaultSecret(externalsecret, "key"), "\\n", "\n", -1)

    internalsecret := vault.GetSecret(os.Getenv("INTERNAL_CERT_VAULT_PATH"))
    InternalCert = strings.Replace(vault.GetFieldFromVaultSecret(internalsecret, "cert"), "\\n", "\n", -1)
    InternalCert = strings.Replace(InternalCert, "\\r", "\r", -1)
    InternalCa = strings.Replace(vault.GetFieldFromVaultSecret(internalsecret, "ca"), "\\n", "\n", -1)
    InternalCa = strings.Replace(InternalCa, "\\r", "\r", -1)
    InternalKey = strings.Replace(vault.GetFieldFromVaultSecret(internalsecret, "key"), "\\n", "\n", -1)

	Kubetoken = vault.GetField(os.Getenv("KUBERNETES_TOKEN_VAULT_PATH"), "token")
    regionapisecret := os.Getenv("REGIONAPI_SECRET")
    Regionapiusername = vault.GetField(regionapisecret, "username")
    Regionapipassword = vault.GetField(regionapisecret, "password")
	Regionapilocation = vault.GetField(regionapisecret, "location")
	
	Ingresstlssecret = os.Getenv("INGRESS_TLS_SECRET_NAME")
	DefaultDomain = os.Getenv("DEFAULT_DOMAIN")
	InsideDomain = os.Getenv("INSIDE_DOMAIN")
	Kubernetesapiurl = os.Getenv("KUBERNETES_API_SERVER")
	initBlacklist()
	HTTPClient = &http.Client{}
}

func initBlacklist() {
	Blacklist = make(map[string]bool)
	blackliststring := os.Getenv("NAMESPACE_BLACKLIST")
	blacklistslice := strings.Split(blackliststring, ",")
	for _, element := range blacklistslice {
			Blacklist[element] = true
	}
	keys := make([]string, 0, len(Blacklist))
	for k := range Blacklist {
			keys = append(keys, k)
	}

	fmt.Printf("Setting blacklist to %v\n", strings.Join(keys, ","))

}
