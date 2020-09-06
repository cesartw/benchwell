package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"bitbucket.org/goreorto/benchwell/config"
	"bitbucket.org/goreorto/benchwell/ctrl"
	"bitbucket.org/goreorto/benchwell/gtk"
	"bitbucket.org/goreorto/benchwell/sqlengine"
)

var verbose bool
var logfile string

var rootCmd = &cobra.Command{
	Use:   "benchwell",
	Short: "Benchwell",
	Long:  `Visit https://benchwell.io for more details`,
	RunE: func(cmd *cobra.Command, args []string) error {
		config.Init()

		if verbose {
			config.SetLevel(logrus.DebugLevel)
		}

		if logfile != "" {
			f, err := os.OpenFile(logfile, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModeAppend|os.ModePerm)
			if err != nil {
				panic(err)
			}
			config.SetOutput(f)
		}

		config.Debug("application startup")

		eng := sqlengine.New()
		defer eng.Dispose()

		ctr := ctrl.AppCtrl{}.Init(eng)

		var err error
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

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&logfile, "logfile", "f", "", "log out file")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
