package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/russellchadwick/configurationservice"
)

func main() {
	key := flag.String("key", "", "Configuration key")
	flag.Parse()

	if *key == "" {
		flag.Usage()
		log.Error("Required input not present")
		return
	}

	client := configurationservice.Client{}
	val, err := client.GetConfiguration(*key)
	if err != nil {
		log.WithField("error", err).Panic("unable to get configuration")
	}

	log.WithField("value", *val).Info("got configuration")
}
