package cmd

import (
	"fmt"
	"os"

	"bitbucket.org/goreorto/sqlaid/config"
	"bitbucket.org/goreorto/sqlaid/ctrl"
	"bitbucket.org/goreorto/sqlaid/gtk"
	"bitbucket.org/goreorto/sqlaid/sqlengine"
	"github.com/getlantern/systray"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var version = "nightly"

const appID = "com.sqlaid"

var rootCmd = &cobra.Command{
	Use:   "sqlaid",
	Short: "SQLaid: Database",
	Long:  `Visit https://sqlaid.com for more details`,
	RunE: func(cmd *cobra.Command, args []string) error {
		eng := sqlengine.New(config.Env)
		defer eng.Dispose()

		// Create a new application.
		app, err := gtk.New(appID)
		if err != nil {
			return err
		}

		ctl, err := ctrl.AppCtrl{}.Init(ctrl.Options{
			Engine: eng,
			Config: config.Env,
			App:    app,
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

		// tray
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
		)

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
