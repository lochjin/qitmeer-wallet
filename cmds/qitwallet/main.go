package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/HalalChain/qitmeer-wallet/config"
	"github.com/HalalChain/qitmeer-wallet/version"
	"github.com/HalalChain/qitmeer-wallet/wallet"
	"github.com/HalalChain/qitmeer-wallet/wserver"
)

func main() {
	cmd := newCmd()
	err := cmd.Execute()
	if err != nil {
		log.Errorf("main err: %s", err)
	}
}

func newCmd() (cmd *cobra.Command) {
	var cfg *config.Config

	// root Command
	cmd = &cobra.Command{
		Use:   "qitmeer-wallet",
		Short: `qitmeer-wallet`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			cfg, err = config.LoadConfig(cfg.ConfigFile, false, cfg)
			if err != nil {
				return fmt.Errorf("cmd PersistentPreRunE err: %s", err)
			}
			logLevel, err := log.ParseLevel(cfg.DebugLevel)
			if err != nil {
				return fmt.Errorf("cmd LogLevl err: %s", err)
			}
			log.SetLevel(logLevel)

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			QitmeerMain(cfg)
		},
	}

	// Create Wallet Command
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "create new wallet or recover wallet from seed",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("wallet create")
		},
	}

	// version Command
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "show version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version.Version())
		},
	}

	cmd.AddCommand(versionCmd)
	cmd.AddCommand(createCmd)

	//cmd flags
	gFlags := cmd.PersistentFlags()
	cfg = config.NewDefaultConfig()
	gFlags.StringVarP(&cfg.ConfigFile, "config", "C", cfg.ConfigFile, "Path to configuration file")
	gFlags.StringVarP(&cfg.AppDataDir, "appdata", "A", cfg.AppDataDir, "Application data directory for wallet config, databases and logs")
	gFlags.StringVarP(&cfg.DebugLevel, "debuglevel", "d", cfg.DebugLevel, "Logging level {trace, debug, info, warn, error, critical}")
	gFlags.BoolVar(&cfg.UI, "ui", true, "Start Wallet with RPC and webUI interface")

	return
}

// QitmeerMain wallet main
func QitmeerMain(cfg *config.Config) {
	log.Trace("Qitmeer Main")

	wt, err := wallet.NewWallet(cfg)
	if err != nil {
		log.Errorf("newWallet err: %s", err)
		return
	}
	wt.Start()

	wsvr, err := wserver.NewWalletServer(cfg, wt)
	if err != nil {
		log.Errorf("NewWalletServer err: %s", err)
		return
	}
	wsvr.Start()

	exitCh := make(chan int)
	<-exitCh
}