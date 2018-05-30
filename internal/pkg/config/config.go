package config

import (
	"path/filepath"
	"strings"

	pimapi "git.tmaws.io/ProductInventoryManagement/go-pimapi"
	"github.com/mgoodness/namespace-controller/pkg/ldap"
	"github.com/mgoodness/namespace-controller/pkg/tiller"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config holds configuration
type Config struct {
	Debug      bool
	Kubeconfig string
	Ldap       ldap.Config
	Manifests  string
	Namespaces string
	Pim        pimapi.Config
	Tiller     tiller.Config
}

// New creates a new Config object
func New(configFile *string) (config *Config, err error) {
	viper.BindPFlags(pflag.CommandLine)

	viper.BindEnv("debug")
	viper.BindEnv("ldap_enabled")
	viper.BindEnv("ldap_environment")
	viper.BindEnv("ldap_hostname")
	viper.BindEnv("ldap_password")
	viper.BindEnv("ldap_username")

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	configPath, _ := filepath.Abs(*configFile)
	viper.SetConfigName(strings.TrimSuffix(filepath.Base(configPath), filepath.Ext(configPath)))
	viper.AddConfigPath(filepath.Dir(configPath))

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	if err = viper.Unmarshal(&config); err != nil {
		return
	}

	return
}
