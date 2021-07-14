package cmd

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
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
		Use:   "gac",
		Short: "Find dark (stale) Route53 records",
		// Uncomment the following line if your bare application
		// has an action associated with it:
		//      Run: func(cmd *cobra.Command, args []string) { },
	}
)

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatalf("ERROR: %v", err.Error())
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.google-admin.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".google-admin") // name of config file (without extension)
	viper.AddConfigPath("$HOME")         // adding home directory as first search path
	viper.SetEnvPrefix("google_admin")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

}

// func init() {
// 	parseEnv()
// }

func initSession() {
	s := session.Must(session.NewSession())
	r := route53.New(s)

	if len(ZONES) == 0 {
		ids, err := getPublicZoneIds(r)
		if err != nil {
			log.Fatal(err)
		}
		populateDB(r, ids)
	} else {
		populateDB(r, ZONES)
	}

	scan()
	reportCSV()
	// reportJSON()
}

func populateDB(r *route53.Route53, ids []string) {
	for _, v := range ids {
		getResourceRecords(r, v)
	}
}
