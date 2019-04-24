package main

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"

	"github.com/rdxsl/docker-logs-agent/controllers"
	_ "github.com/rdxsl/docker-logs-agent/routers"

	"github.com/astaxie/beego"
)

var Version = "version notset"

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	controllers.Version = Version

	cert, err := tls.LoadX509KeyPair("conf/dev/mtls/certs/server.pem", "conf/dev/mtls/certs/server.key")
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)

	}
	certpool := x509.NewCertPool()
	pem, err := ioutil.ReadFile("conf/dev/mtls/certs/ca.pem")
	if err != nil {
		log.Fatalf("Failed to read client certificate authority: %v", err)
	}
	if !certpool.AppendCertsFromPEM(pem) {
		log.Fatalf("Can't parse client certificate authority")
	}

	config := tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
		ClientCAs:    certpool,
	}
	config.Rand = rand.Reader

	beego.BeeApp.Server.TLSConfig = &config

	beego.Run()
}
