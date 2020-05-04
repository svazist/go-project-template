package cmd

import (
	"fmt"
	"github.com/apex/log"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/svazist/go-project-template/logger"
	"os"
	"strings"
)

var Version string

var (
	cfgFile, logLevel  string
	cfgFileDefaultName = "config.yaml"
	logCtx             log.Interface
	envPrefix          = "bs"
	initMetrics        = false
	IsDebug            = false
)

var RootCmd = &cobra.Command{
	Use:   "batch-server",
	Short: "Batch JSON requests",
	Long:  `Batch JSON requests`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	viper.SetEnvPrefix(envPrefix)
	cobra.OnInitialize(initConfig, initLogger)

	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", cfgFileDefaultName, "config file (default is $HOME/.php-batch_server.yaml)")
	RootCmd.PersistentFlags().StringVar(&logLevel, "log.level", "Info", "Only log messages with the given severity or above. Valid levels: [debug, info, warn, error, fatal]")

	envs := map[string]string{
		"LOG_LEVEL": "log.level",
	}

	mapEnvVars(envs, RootCmd)

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	viper.SetConfigType("yaml")

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".php-fpm_exporter" (without extension).
		viper.AddConfigPath(".")
		viper.AddConfigPath("/etc/batch-server")
		viper.AddConfigPath(home)
		viper.SetConfigName(cfgFileDefaultName)
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv() // read in environment variables that match

}

// initLogger configures the log level
func initLogger() {

	//fmt.Println("before init logger")

	logger.InitLogger()
	//server.SetLogger(log)

	//if value := os.Getenv("BATCH_SERVER_LOG_LEVEL"); value != "" {
	//	logLevel = value
	//}
	//
	//lvl, err := logrus.ParseLevel(logLevel)
	//if err != nil {
	//	lvl = logrus.InfoLevel
	//	log.Fatalf("Could not set log level to '%v'.", logLevel)
	//}
	//
	//log.SetLevel(lvl)
}

func mapEnvVars(envs map[string]string, cmd *cobra.Command) {

	for env, flagName := range envs {
		flag := cmd.Flags().Lookup(flagName)
		if flag == nil && cmd.HasPersistentFlags() {
			flag = cmd.PersistentFlags().Lookup(flagName)
		}

		if flag == nil {
			log.Errorf("Params %s not found", flagName)
			return
		}

		flag.Usage = fmt.Sprintf("%v [env %v]", flag.Usage, env)
		if value := os.Getenv(env); value != "" {

			if err := flag.Value.Set(value); err != nil {
				log.Errorf("set value error : %s", err)
			} else {
				viper.Set(flagName, value)
			}
		}
	}
}
