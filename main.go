package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Marco98/paperlessmerge/pkg/paperless"
	"github.com/adrg/xdg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

func main() {
	cobra.OnInitialize(initConfig)
	var rootCmd = &cobra.Command{
		Use:   "paperlessmerge docid1 docid2 ...",
		Short: "PaperlessMerge can easily merge documents in paperless-ng(x)",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			i, err := paperless.New(
				viper.GetString("server"),
				viper.GetString("username"),
				viper.GetString("password"),
				viper.GetBool("ignoretls"),
			)
			if err != nil {
				return err
			}
			ids, err := strToIntArr(args)
			if err != nil {
				return err
			}
			if err := i.MergePDF(ids); err != nil {
				return err
			}
			fmt.Println("new document was uploaded")
			if del, err := cmd.Flags().GetBool("delete"); err == nil && del {
				if err := i.DeleteDocuments(ids); err != nil {
					return err
				}
				fmt.Println("old documents were deleted")
			}
			return nil
		},
	}
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")
	rootCmd.PersistentFlags().StringP("server", "s", "", "paperless base url")
	rootCmd.PersistentFlags().StringP("username", "u", "", "paperless auth username")
	rootCmd.PersistentFlags().StringP("password", "p", "", "paperless auth password")
	rootCmd.PersistentFlags().BoolP("ignoretls", "k", false, "do not validate tls certificate")
	rootCmd.Flags().BoolP("delete", "d", false, "deletes merged documents")
	if err := viper.BindPFlag("server", rootCmd.PersistentFlags().Lookup("server")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("username", rootCmd.PersistentFlags().Lookup("username")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("ignoretls", rootCmd.PersistentFlags().Lookup("ignoretls")); err != nil {
		panic(err)
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(filepath.Join(xdg.ConfigHome, "paperlessmerge"))
		viper.SetConfigType("toml")
		viper.SetConfigName("config.toml")
	}
	viper.SetEnvPrefix("plm")
	viper.AutomaticEnv()
	viper.SetDefault("ignoretls", false)

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func strToIntArr(arr []string) ([]int, error) {
	ids := make([]int, len(arr))
	var err error
	for n, v := range arr {
		ids[n], err = strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
	}
	return ids, nil
}
