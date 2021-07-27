package cmd

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	JSONOut bool
	name    string
	scanAll bool
	timeout int
	ports   = []string{}
	zones   = []string{}

	RRtypes = map[string]bool{ // which record types to scan
		"CNAME": true,
	}

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
	RootCmd.PersistentFlags().BoolVarP(&JSONOut, "json", "j", false, "json output")
	RootCmd.PersistentFlags().StringVarP(&name, "name", "n", "report", "report file name")
	RootCmd.PersistentFlags().BoolVarP(&scanAll, "all", "a", false, "scan A/AAAA records (in addition to aliases)")
	RootCmd.PersistentFlags().StringSliceVarP(&ports, "ports", "p", ports, "TCP ports to scan")
	RootCmd.PersistentFlags().IntVarP(&timeout, "timeout", "t", 10, "port scan timeout")
	RootCmd.PersistentFlags().StringSliceVarP(&zones, "zones", "z", zones, "zone ids to scan")
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

	// scanAll = viper.Get("all").(bool)
	if scanAll {
		RRtypes["A"] = true
		RRtypes["AAAA"] = true
	}

}

func scanner(cmd *cobra.Command, args []string) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}

	ctx := context.Background()
	r53Client := route53.NewFromConfig(cfg)

	if len(zones) == 0 {
		zones = viper.GetStringSlice("zones")
		if len(zones) == 0 {
			zones, err = getPublicZoneIDs(ctx, r53Client)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	populateDB(ctx, r53Client, zones)

	scan()

	if JSONOut {
		reportJSON()
	} else {
		reportCSV()
	}
}
