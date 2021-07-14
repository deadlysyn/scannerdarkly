package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type dnsRecord struct {
	Name   string
	Type   string
	Alias  bool
	Values []string
	Active []string
}

var (
	cfgFile string

	// DB maps Route53 Zone IDs to slices of dnsRecords
	DB = make(map[string][]dnsRecord)

	RootCmd = &cobra.Command{
		Use:   "d",
		Short: "Scan for dark (stale) Route53 records",
		Run:   scanner,
	}
)

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatalf("ERROR: %v", err.Error())
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.yml", "config file")
}

func initConfig() {
	viper.SetConfigType("yaml")

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	viper.SetEnvPrefix("pkd")
	viper.AutomaticEnv()

	c := viper.Get("config")
	if c != nil {
		viper.SetConfigFile(c.(string))
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("ERROR: %v", err)
	}
}

func scanner(cmd *cobra.Command, args []string) {
	fmt.Println("hello from scanner")
}

// func initSession() {
// 	s := session.Must(session.NewSession())
// 	r := route53.New(s)

// 	if len(ZONES) == 0 {
// 		ids, err := getPublicZoneIds(r)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		populateDB(r, ids)
// 	} else {
// 		populateDB(r, ZONES)
// 	}

// 	scan()
// 	reportCSV()
// 	// reportJSON()
// }

// func populateDB(r *route53.Route53, ids []string) {
// 	for _, v := range ids {
// 		getResourceRecords(r, v)
// 	}
// }
