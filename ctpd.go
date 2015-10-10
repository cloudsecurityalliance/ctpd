//    Copyright 2015 Cloud Security Alliance EMEA (cloudsecurityalliance.org)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/cloudsecurityalliance/ctpd/server"
	"github.com/cloudsecurityalliance/ctpd/server/ctp"
)

const CTPD_VERSION = 0.1

var configFileFlag string
var versionFlag bool
var helpFlag bool

func init() {
	flag.StringVar(&configFileFlag, "config", "/path/to/file", "Specify an alternative configuration file to use.")
	flag.BoolVar(&versionFlag, "version", false, "Print version information.")
	flag.BoolVar(&helpFlag, "help", false, "Print help.")
}

func main() {
	var ok bool
	var conf ctp.Configuration

	flag.Parse()

	if versionFlag {
		fmt.Println("ctpd version %f.", CTPD_VERSION)
		fmt.Println(" Copyright 2015 Cloud Security Alliance EMEA (cloudsecurityalliance.org).")
		fmt.Println(" ctpd is licensed under the Apache License, Version 2.0.")
		fmt.Println(" see http://www.apache.org/licenses/LICENSE-2.0")
		fmt.Println("")
		return
	}

	if helpFlag {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags]\n", os.Args[0])
		flag.PrintDefaults()
		return
	}

	if configFileFlag == "/path/to/file" {
		conf, ok = ctp.SearchAndLoadConfigurationFile()
	} else {
		conf, ok = ctp.LoadConfigurationFromFile(configFileFlag)
	}

	if !ok {
        log.Printf("No configuration file was loaded, using defaults.")
        conf = ctp.ConfigurationDefaults
	}

	if conf["client"] != "" {
		http.Handle("/", http.FileServer(http.Dir(conf["client"])))
	}

	http.Handle(conf["basepath"], server.NewCtpApiHandlerMux(conf))
	if conf["tls_use"] != "" && conf["tls_use"] != "no" {
		if conf["tls_use"] != "yes" {
			log.Fatal("Configuration: tls_use must be either 'yes' or 'no'")
		}
		if conf["tls_key_file"] == "" || conf["tls_cert_file"] == "" {
			log.Fatal("Missing tls_key_file or tls_cert_file in configuration.")
		}
		log.Printf("Starting ctpd with TLS enabled at %s", conf["listen"])
		log.Fatal(http.ListenAndServeTLS(conf["listen"], conf["tls_cert_file"], conf["tls_key_file"], nil))
	} else {
		log.Printf("Starting ctpd at %s", conf["listen"])
		log.Fatal(http.ListenAndServe(conf["listen"], nil))
	}
}