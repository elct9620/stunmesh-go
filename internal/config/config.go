package config

import (
	"errors"
	"time"

	"github.com/google/wire"
	"github.com/spf13/viper"
	"github.com/tjjh89017/stunmesh-go/internal/entity"
)

var DefaultSet = wire.NewSet(
	Load,
	NewDeviceConfig,
	wire.Bind(new(entity.PeerAllower), new(*DeviceConfig)),
)

const Name = "config"

var Paths []string = []string{
	"$STUNMESH_CONFIG_DIR",
	"/etc/stunmesh",
	"$HOME/.stunmesh",
	".",
}

var (
	ErrBindEnv         = errors.New("failed to bind env")
	ErrReadConfig      = errors.New("failed to read config")
	ErrUnmarshalConfig = errors.New("failed to unmarshal config")
)

var envs = map[string][]string{
	"cloudflare.api_key":   {"CF_API_KEY", "CLOUDFLARE_API_KEY"},
	"cloudflare.api_email": {"CF_API_EMAIL", "CLOUDFLARE_API_EMAIL"},
	"cloudflare.api_token": {"CF_API_TOKEN", "CLOUDFLARE_API_TOKEN"},
	"cloudflare.zone_name": {"CF_ZONE_NAME", "CLOUDFLARE_ZONE_NAME"},
	"refresh_interval":     {"REFRESH_INTERVAL"},
}

type Logger struct {
	Level string `mapstructure:"level"`
}

type Stun struct {
	Address string `mapstructure:"address"`
}

type Cloudflare struct {
	ApiKey   string `mapstructure:"api_key"`
	ApiEmail string `mapstructure:"api_email"`
	ApiToken string `mapstructure:"api_token"`
	ZoneName string `mapstructure:"zone_name"`
}

type Config struct {
	Interfaces      Interfaces    `mapstructure:"interfaces"`
	RefreshInterval time.Duration `mapstructure:"refresh_interval"`
	Log             Logger        `mapstructure:"log"`
	Stun            Stun          `mapstructure:"stun"`
	Cloudflare      Cloudflare    `mapstructure:"cloudflare"`
}

func Load() (*Config, error) {
	viper.SetConfigName(Name)
	for _, path := range Paths {
		viper.AddConfigPath(path)
	}
	viper.AutomaticEnv()

	viper.SetDefault("refresh_interval", time.Duration(10)*time.Minute)
	viper.SetDefault("stun.addr", "stun.l.google.com:19302")

	for envName, keys := range envs {
		binding := []string{envName}
		binding = append(binding, keys...)

		if err := viper.BindEnv(binding...); err != nil {
			return nil, errors.Join(ErrBindEnv, err)
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, errors.Join(ErrReadConfig, err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, errors.Join(ErrUnmarshalConfig, err)
	}

	return &cfg, nil
}
