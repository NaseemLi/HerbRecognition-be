package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server       ServerConfig       `mapstructure:"server"`
	Database     DatabaseConfig     `mapstructure:"database"`
	ModelService ModelServiceConfig `mapstructure:"model_service"`
	JWT          JWTConfig          `mapstructure:"jwt"`
	CORS         CORSConfig         `mapstructure:"cors"`
	Admin        AdminConfig        `mapstructure:"admin"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	DBName          string `mapstructure:"dbname"`
	Charset         string `mapstructure:"charset"`
	ParseTime       bool   `mapstructure:"parseTime"`
	Loc             string `mapstructure:"loc"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

// 拼接 DSN 连接字符串
func (c *DatabaseConfig) BuildDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
		c.Charset,
		c.ParseTime,
		c.Loc,
	)
}

type ModelServiceConfig struct {
	URL           string `mapstructure:"url"`
	ONNXModelPath string `mapstructure:"onnx_model_path"`
	ClassesPath   string `mapstructure:"classes_path"`
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

// CORSConfig CORS 配置
type CORSConfig struct {
	AllowOrigins []string `mapstructure:"allow_origins"`
}

// AdminConfig 默认管理员配置
type AdminConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

var Conf *Config

func Init() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	envBindings := map[string]string{
		"server.port":                   "SERVER_PORT",
		"server.mode":                   "SERVER_MODE",
		"database.host":                 "DB_HOST",
		"database.port":                 "DB_PORT",
		"database.user":                 "DB_USER",
		"database.password":             "DB_PASSWORD",
		"database.dbname":               "DB_NAME",
		"model_service.onnx_model_path": "ONNX_MODEL_PATH",
		"model_service.classes_path":    "CLASSES_PATH",
		"admin.username":                "ADMIN_USERNAME",
		"admin.password":                "ADMIN_PASSWORD",
	}

	for key, envName := range envBindings {
		if err := viper.BindEnv(key, envName); err != nil {
			return err
		}
	}

	if configName := strings.TrimSpace(os.Getenv("CONFIG_NAME")); configName != "" {
		viper.SetConfigName(configName)
	}

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	Conf = &Config{}
	if err := viper.Unmarshal(Conf); err != nil {
		return err
	}
	return nil
}
