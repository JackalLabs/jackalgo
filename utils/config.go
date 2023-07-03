package utils

import (
	"errors"
	"os"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/viper"
	tmcfg "github.com/tendermint/tendermint/config"
	tmlog "github.com/tendermint/tendermint/libs/log"
)

type Context struct {
	Viper  *viper.Viper
	Config *Config
	Logger tmlog.Logger
}

const JackalGoContextKey = sdk.ContextKey("jackalgo.context")

func NewContext(v *viper.Viper, config *Config, logger tmlog.Logger) *Context {
	return &Context{v, config, logger}
}

func DefaultBaseConfig() BaseConfig {
	return BaseConfig{
		LogLevel:  tmcfg.DefaultLogLevel,
		LogFormat: tmcfg.LogFormatPlain,
	}
}

// DefaultConfig returns a default configuration for a Tendermint node
func DefaultConfig() *Config {
	return &Config{
		BaseConfig: DefaultBaseConfig(),
	}
}

type BaseConfig struct {
	// chainID is unexposed and immutable but here for convenience
	//nolint:all
	chainID string

	// The root directory for all data.
	// This should be set in viper so it can unmarshal into this struct
	RootDir string `mapstructure:"home"`

	LogLevel string `mapstructure:"log_level"`

	// Output format: 'plain' (colored text) or 'json'
	LogFormat string `mapstructure:"log_format"`
}

type Config struct {
	BaseConfig `mapstructure:",squash"`
}

func (cfg BaseConfig) ValidateBasic() error {
	switch cfg.LogFormat {
	case tmcfg.LogFormatPlain, tmcfg.LogFormatJSON:
	default:
		return errors.New("unknown log_format (must be 'plain' or 'json')")
	}
	return nil
}

func (cfg *Config) ValidateBasic() error {
	if err := cfg.BaseConfig.ValidateBasic(); err != nil {
		return err
	}

	return nil
}

func (cfg *Config) SetRoot(root string) *Config {
	cfg.BaseConfig.RootDir = root
	return cfg
}

func NewDefaultContext() *Context {
	return NewContext(
		viper.New(),
		DefaultConfig(),
		tmlog.NewTMLogger(tmlog.NewSyncWriter(os.Stdout)),
	)
}
