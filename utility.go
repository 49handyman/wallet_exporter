package main

import (
	"bufio"
	"os"
	"os/user"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// Prefer commandline flags, fall back to zcash.conf if available
// TODO add env vars to the lookups
func reconcileConfigs() {
	if *rpcPort == "" {
		if zcashConfValues.rpcPort != "" {
			*rpcPort = zcashConfValues.rpcPort
		} else if zcashConfValues.testNet {
			*rpcPort = "18232"
		} else {
			*rpcPort = "8232"
		}
	}
	if *rpcUser == "" {
		if zcashConfValues.rpcUser != "" {
			*rpcUser = zcashConfValues.rpcUser
		} else {
			log.Fatalln("RPC user missing, a value is required")
		}
	}
	if *rpcPassword == "" {
		if zcashConfValues.rpcPassword != "" {
			*rpcPassword = zcashConfValues.rpcPassword
		} else {
			log.Fatalln("RPC password missing, a value is required")
		}
	}
}

// Read from a zcash.conf provided from command line flag
// or from the default location (user's home directory)
// To avoid reading a default file, use --zcash.conf.path="ignore"
func readZcashConf() (zcashConf, error) {
	var zcashConf zcashConf
	if *zcashConfPath == "" {
		user, err := user.Current()
		if err != nil {
			return zcashConf, err
		}
		if _, err := os.Stat(user.HomeDir + "/.zcash/zcash.conf"); os.IsNotExist(err) {
			log.Infoln("No zcash.conf flag provided, and no default found. Skipping")
			return zcashConf, nil
		}
		*zcashConfPath = user.HomeDir + "/.zcash/zcash.conf"
	}
	log.Infoln("Attempting to read zcash.conf at:", *zcashConfPath)
	confFile, err := os.Open(*zcashConfPath)
	if err != nil {
		log.Infoln("Error reading zcash.conf file, set --zcash.conf.path=\"\" to ignore the file")
		return zcashConf, err
	}
	defer confFile.Close()

	scanner := bufio.NewScanner(confFile)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), "=")
		switch strings.TrimSpace(s[0]) {
		case "testnet":
			if strings.TrimSpace(s[1]) == "1" {
				zcashConf.testNet = true
			}
		case "rpcuser":
			zcashConf.rpcUser = strings.TrimSpace(s[1])
		case "rpcpassword":
			zcashConf.rpcPassword = strings.TrimSpace(s[1])
		case "rpcport":
			zcashConf.rpcPort = strings.TrimSpace(s[1])
		}
	}
	return zcashConf, nil
}
