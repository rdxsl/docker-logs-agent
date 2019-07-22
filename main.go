package main

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"

	"github.com/rdxsl/docker-agent-proxy/controllers"
	_ "github.com/rdxsl/docker-agent-proxy/routers"

	"github.com/astaxie/beego"
)

var Version = "version notset"

func setupMTSL() (cert tls.Certificate, certpool *x509.CertPool) {
	cert, err := tls.LoadX509KeyPair(beego.AppConfig.String("HTTPSCertFile"), beego.AppConfig.String("HTTPSKeyFile"))
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}

	certpool = x509.NewCertPool()
	pem, err := ioutil.ReadFile(beego.AppConfig.String("HTTPSCAFile"))
	if err != nil {
		log.Fatalf("Failed to read client certificate authority: %v", err)
	}
	if !certpool.AppendCertsFromPEM(pem) {
		log.Fatalf("Can't parse client certificate authority")
	}
	return
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	controllers.Version = Version

	cert, certpool := setupMTSL()

	config := tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
		ClientCAs:    certpool,
	}
	config.Rand = rand.Reader

	beego.BeeApp.Server.TLSConfig = &config

	beego.Run()
}
