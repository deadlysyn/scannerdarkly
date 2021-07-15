package cmd

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
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
	scanArecords bool
	cfgFile      string
	outputFormat string
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

	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.yml", "config file")
	RootCmd.PersistentFlags().StringVarP(&outputFormat, "format", "f", "csv", "output format")
	RootCmd.PersistentFlags().BoolVarP(&scanArecords, "alias-only", "a", false, "scan A/AAAA records")
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

	// for _, v := range viper.GetStringSlice("zones") {
	// 	result, _ := getParam(r53Client, fmt.Sprintf("%s/%s", viper.GetString("ssm.prefix"), v))
	// 	creds[v] = aws.ToString(result.Parameter.Value)
	// }
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
