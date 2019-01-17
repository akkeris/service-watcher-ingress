# service-watcher-ingress


Environment Variables

CLUSTER - identifier for clustser (example "maru")

DEFAULT_DOMAIN - domain name for external apps that will be appended to the application name (example "maruapp.octanner.io")

EXTERNAL_CERT_VAULT_PATH - location in vault of the secret containing the fields "cert," "ca," and "key"  of SSL Certificate to be use by external applications (example secret/ops/certs/star.maruappi.octanner.io) 

INGRESS_TLS_SECRET_NAME - the name to be used for the kubernetes secret for the certificat material

INTERNAL_CERT_VAULT_PATH - see above

INSIDE_DOMAIN - see above 

KUBERNETES_TOKEN_VAULT_PATH - location in vault of the secret containing the "token" field for the kubernetes bearer token

KUBERNETES_API_SERVER - URL of kubernetes api servers.  (no https://)

NAMESPACE_BLACKLIST - comma separated list of spaces to ignore

REGIONAPI_SECRET - location in vault of the secret containing the "location," "username," and "password" fields used to access the region api

VAULT_ADDR - URL of vault

VAULT_TOKEN - Token to be used for vault
