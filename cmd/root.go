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
	cfgFile      string
	outputFormat string
	scanArecords bool
	scanPorts    []string
	scanTimeout  int
	zoneIDs      []string

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

	RootCmd.PersistentFlags().BoolVarP(&scanArecords, "all", "a", false, "scan A/AAAA records (in addition to CNAMEs)")
	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.yml", "config file")
	RootCmd.PersistentFlags().StringVarP(&outputFormat, "format", "f", "csv", "output format")
	RootCmd.PersistentFlags().StringSliceVarP(&scanPorts, "port", "p", []string{}, "TCP ports to scan")
	RootCmd.PersistentFlags().IntVarP(&scanTimeout, "timeout", "t", 10, "port scan timeout")
	RootCmd.PersistentFlags().StringSliceVarP(&zoneIDs, "zone-id", "z", []string{}, "zone ids to scan")
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
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}

	ctx := context.Background()
	r53Client := route53.NewFromConfig(cfg)

	if len(zoneIDs) == 0 {
		zoneIDs = viper.GetStringSlice("zones")
		if len(zoneIDs) == 0 {
			zoneIDs, err = getPublicZoneIDs(ctx, r53Client)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	populateDB(ctx, r53Client, zoneIDs)
	scan()
	// 	reportCSV()
	// reportJSON()
}
