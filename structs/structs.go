package structs

type Spacespec struct {
	Name     string `json:"name"}`
	Internal bool   `json:"internal"}`
}

type TlsSecret struct {
	Kind     string `json:"kind"`
	Metadata struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
	} `json:"metadata"`
	Data struct {
		TlsCrt string `json:"tls.crt"`
		TlsKey string `json:"tls.key"`
	} `json:"data"`
	Type string `json:"type"`
}

type SecretsList struct {
	Items []struct {
		Metadata struct {
			Name      string `json:"name"`
			Namespace string `json:"namespace"`
		} `json:"metadata"`
		Type string `json:"type"`
	} `json:"items"`
}

type SimpleAkkerisIngress struct {
	Kind     string `json:"kind"`
	Metadata struct {
		Annotations struct {
			KubernetesIoIngressClass string `json:"kubernetes.io/ingress.class"`
		} `json:"annotations"`
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
	} `json:"metadata"`
	Spec struct {
		Rules []Rule `json:"rules"`
		TLS   []TLS  `json:"tls"`
	} `json:"spec"`
}

type TLS struct {
	Hosts      []string `json:"hosts"`
	SecretName string   `json:"secretName"`
}

type Rule struct {
	Host string `json:"host"`
	HTTP struct {
		Paths []Path `json:"paths"`
	} `json:"http"`
}

type Path struct {
	Backend Backend `json:"backend"`
}

type Backend struct {
	ServiceName string `json:"serviceName"`
	ServicePort int    `json:"servicePort"`
}
