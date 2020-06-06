package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"bitbucket.org/goreorto/sqlaid/assets"
	"bitbucket.org/goreorto/sqlaid/config"
	"bitbucket.org/goreorto/sqlaid/ctrl"
	"bitbucket.org/goreorto/sqlaid/gtk"
	"bitbucket.org/goreorto/sqlaid/sqlengine"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var version = "nightly"

var rootCmd = &cobra.Command{
	Use:   "sqlaid",
	Short: "SQLaid: Database",
	Long:  `Visit https://sqlaid.com for more details`,
	RunE: func(cmd *cobra.Command, args []string) error {
		config.Env.Log.Debug("application startup")
		err := assets.Load()
		if err != nil {
			panic(err)
		}

		eng := sqlengine.New(config.Env)
		defer eng.Dispose()

		ctr := &ctrl.AppCtrl{}
		ctr.Engine = eng
		ctr.Config = config.Env

		// Create a new application.
		ctr.App, err = gtk.Application{}.Init(ctr)
		if err != nil {
			return err
		}

		/*
			systray.RunWithAppWindow(
				"SQLaid",
				400, 400,
				func() { // ready
					systray.SetIcon(TrayIcon)
					systray.SetTitle("SQLaid")
					systray.SetTooltip("Ultimate database GUI")

					mQuit := systray.AddMenuItem("Quit", "I'm out")
					mShow := systray.AddMenuItem("Show All", "Display windows")
					go func() {
						<-mQuit.ClickedCh
						fmt.Println("Requesting quit")
						systray.Quit()
						fmt.Println("Finished quitting")
					}()
					go func() {
						for {
							<-mShow.ClickedCh
							ctl.ShowAll()
						}
					}()

					// Launch the application
					go func() {
						app.Run(nil)
					}()
				},
				func() { //quit
				},
			)*/

		ctr.App.Run(nil)

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
		fmt.Sprintf("config file (default is %s/config.json)", xdgHome+"/sqlaid/config.json"))
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringP("logfile", "f", "", "log out file")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("logfile", rootCmd.PersistentFlags().Lookup("logfile"))
}

func initConfig() {
	if cfgFile == "" {
		xdgHome, _ := os.UserConfigDir()
		cfgFile = xdgHome + "/sqlaid/config.json"
	}

	viper.SetConfigFile(cfgFile)

	err := viper.ReadInConfig()
	if err == nil {
		fmt.Println("Using config file: ", viper.ConfigFileUsed())
	} else {
		err = ioutil.WriteFile(cfgFile, []byte(assets.DEFAULT_CONFIG), os.FileMode(0666))
		if err != nil {
			fmt.Println("Cannot write config file: ", err.Error())
		}
		err = viper.ReadInConfig()
	}

	if err != nil {
		viper.ReadConfig(bytes.NewBuffer([]byte(assets.DEFAULT_CONFIG)))
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

	config.RegisterStyle()
	config.InitKeyChain(config.Env.EncryptMode)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
