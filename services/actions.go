package services

import (
	utils "akkeris-service-watcher-ingress/utils"
	structs "akkeris-service-watcher-ingress/structs"
	corev1 "k8s.io/api/core/v1"
	"encoding/json"
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/base64"
	"bytes"
)

func DeleteIngress(obj interface{}) {
	servicename := obj.(*corev1.Service).ObjectMeta.Name
	namespace := obj.(*corev1.Service).ObjectMeta.Namespace
	appname := servicename + "-" + namespace
	if namespace == "default" {
		appname = servicename
	}
	req, err := http.NewRequest("DELETE", "https://"+utils.Kubernetesapiurl+"/apis/extensions/v1beta1/namespaces/"+namespace+"/ingresses/"+appname, nil)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Authorization", "Bearer "+utils.Kubetoken)

	resp, doerr := utils.HTTPClient.Do(req)
	if doerr != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
    fmt.Println("delete ingress response: "+resp.Status)
}
func InstallIngress(obj interface{}) {
	servicename := obj.(*corev1.Service).ObjectMeta.Name
	namespace := obj.(*corev1.Service).ObjectMeta.Namespace
	appname := servicename + "-" + namespace
	if namespace == "default" {
		appname = servicename
	}
	needssecret, err := needsTLSSecret(namespace, utils.Ingresstlssecret)
	if err != nil {
		fmt.Println(err)
	}
	if needssecret {
		createTLSSecret(namespace, utils.Ingresstlssecret)
	}
	var url string
	internal := isInternal(namespace)
	if internal {
		 url = appname+"."+utils.InsideDomain
	}
	if !internal {
		url= appname+"."+utils.DefaultDomain
	}
	var ingress structs.SimpleAkkerisIngress
	ingress.Kind = "Ingress"
	ingress.Metadata.Annotations.KubernetesIoIngressClass = "nginx"
	ingress.Metadata.Name = appname
	ingress.Metadata.Namespace = namespace
	var backend structs.Backend
	backend.ServiceName = servicename
	backend.ServicePort = 80
	var path structs.Path
	path.Backend = backend
	var paths []structs.Path
	paths = append(paths, path)
	var rule structs.Rule
	rule.Host = url
	rule.HTTP.Paths = paths
	ingress.Spec.Rules = append(ingress.Spec.Rules, rule)
	var tls structs.TLS
	tls.SecretName =utils.Ingresstlssecret
	tls.Hosts = append(tls.Hosts, url)
	ingress.Spec.TLS = append(ingress.Spec.TLS, tls)

	p, err := json.Marshal(ingress)
	if err != nil {
		fmt.Println(err)
	}

	req, err := http.NewRequest("POST", "https://"+utils.Kubernetesapiurl+"/apis/extensions/v1beta1/namespaces/"+namespace+"/ingresses", bytes.NewBuffer(p))
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Authorization", "Bearer "+utils.Kubetoken)

	resp, doerr := utils.HTTPClient.Do(req)
	if doerr != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
        fmt.Println("install ingress response: "+resp.Status)
}

func isInternal(space string) bool {

	req, err := http.NewRequest("GET", utils.Regionapilocation+"/v1/space/"+space, nil)
	req.SetBasicAuth(utils.Regionapiusername, utils.Regionapipassword)
	if err != nil {
		fmt.Println(err)
	}
	resp, err := utils.HTTPClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	bb, err := ioutil.ReadAll(resp.Body)
	var spaceobject structs.Spacespec
	uerr := json.Unmarshal(bb, &spaceobject)
	if uerr != nil {
		fmt.Println(uerr)
	}
	return spaceobject.Internal
}

func needsTLSSecret(space string, secretname string) (b bool, e error) {
	var toreturn bool
	toreturn = true
	req, err := http.NewRequest("GET", "https://"+utils.Kubernetesapiurl+"/api/v1/namespaces/"+space+"/secrets", nil)
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Authorization", "Bearer "+utils.Kubetoken)

	resp, doerr := utils.HTTPClient.Do(req)
	if doerr != nil {
		fmt.Println(err)
		return toreturn, err
	}
	defer resp.Body.Close()
	bodybytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return toreturn, err
	}
	var secretlist structs.SecretsList
	err = json.Unmarshal(bodybytes, &secretlist)
	if err != nil {
		fmt.Println(err)
		return toreturn, err
	}

	for _, element := range secretlist.Items {
		if element.Metadata.Name == secretname && element.Type == "kubernetes.io/tls" {
			toreturn = false
		}
	}
        if toreturn {
           fmt.Println("namespace "+space+" needs secret")
        }
        if !toreturn {
           fmt.Println("namespace "+space+" already has secret")
        }
	return toreturn, nil

}

func createTLSSecret(space string, secretname string) (e error) {
	fmt.Println("adding secret")
	var tlssecret structs.TlsSecret
	tlssecret.Kind = "Secret"
	tlssecret.Metadata.Name = secretname
	tlssecret.Metadata.Namespace = space
        if isInternal(space) {
          tlssecret.Data.TlsKey = utils.InternalKey
	  tlssecret.Data.TlsCrt = utils.InternalCert + utils.InternalCa
        }
        if ! isInternal(space){
          tlssecret.Data.TlsKey = utils.ExternalKey
          tlssecret.Data.TlsCrt = utils.ExternalCert + utils.ExternalCa
        }
	tlssecret.Type = "kubernetes.io/tls"

	tlssecret.Data.TlsCrt = base64.StdEncoding.EncodeToString([]byte(tlssecret.Data.TlsCrt))
	tlssecret.Data.TlsKey = base64.StdEncoding.EncodeToString([]byte(tlssecret.Data.TlsKey))

	p, err := json.Marshal(tlssecret)
	if err != nil {
		fmt.Println(err)
		return err
	}

	req, err := http.NewRequest("POST", "https://"+utils.Kubernetesapiurl+"/api/v1/namespaces/"+space+"/secrets", bytes.NewBuffer(p))
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Authorization", "Bearer "+utils.Kubetoken)

	resp, doerr := utils.HTTPClient.Do(req)
	if doerr != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()
        fmt.Println("create secret response: "+resp.Status)
	return nil
}

