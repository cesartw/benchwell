package cmd

import (
	"errors"
	"fmt"
	"os"

	"bitbucket.org/goreorto/sqlhero/config"
	"bitbucket.org/goreorto/sqlhero/ctrl"
	"bitbucket.org/goreorto/sqlhero/gtk"
	"bitbucket.org/goreorto/sqlhero/sqlengine"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var version = "nightly"

const appID = "com.iodone.sqlhero"

var rootCmd = &cobra.Command{
	Use:   "sqlhero",
	Short: "SQLHero: Database",
	Long:  `Visit http://sqlhero.com for more details`,
	RunE: func(cmd *cobra.Command, args []string) error {
		eng := sqlengine.New(config.Env)
		defer eng.Dispose()

		// Create a new application.
		app, err := gtk.New(appID)
		if err != nil {
			return err
		}

		ctl, err := ctrl.MainCtrl{}.Init(ctrl.Options{
			Engine:  eng,
			Config:  config.Env,
			Factory: app,
		})
		if err != nil {
			return err
		}

		// Connect function to application startup event, this is not required.
		app.Connect("startup", func() {
			config.Env.Log.Debug("application startup")
		})

		// Connect function to application activate event
		app.Connect("activate", ctl.OnActivate)

		// Connect function to application shutdown event, this is not required.
		app.Connect("shutdown", func() {
			config.Env.Log.Debug("application shutdown")
		})

		// Launch the application
		if app.Run(nil) != 0 {
			return errors.New("exit with error")
		}

		return nil
	},
}

var cfgFile string
var xdgHome string

func init() {
	cobra.OnInitialize(initConfig)
	viper.SetConfigType("json")
	xdgHome, _ = os.UserConfigDir()

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "",
		fmt.Sprintf("config file (default is %s/config.json)", xdgHome+"/sqlhero/config.json"))
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringP("logfile", "f", "log.txt", "log out file")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("logfile", rootCmd.PersistentFlags().Lookup("logfile"))
}

func initConfig() {
	if cfgFile == "" {
		xdgHome, _ := os.UserConfigDir()
		cfgFile = xdgHome + "/sqlhero/config.json"
	}

	viper.SetConfigFile(cfgFile)

	err := viper.ReadInConfig()
	if err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println("Cannot load config file: ", err.Error())
	}

	if viper.GetBool("verbose") {
		config.Env.Log.SetLevel(logrus.DebugLevel)
	}

	if viper.GetString("logfile") != "" {
		f, err := os.OpenFile(viper.GetString("logfile"), os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModeAppend|os.ModePerm)
		if err != nil {
			panic(err)
		}
		config.Env.Log.SetOutput(f)

	}

	viper.Unmarshal(config.Env)
	config.Env.Version = version
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
