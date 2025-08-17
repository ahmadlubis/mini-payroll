package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL  string         `yaml:"database_url" mapstructure:"database_url"`
	JWTSecret    string         `yaml:"jwt_secret" mapstructure:"jwt_secret"`
	SeedDatabase bool           `yaml:"seed_database" mapstructure:"seed_database"`
	Environment  string         `yaml:"environment" mapstructure:"environment"`
	LogLevel     string         `yaml:"log_level" mapstructure:"log_level"`
	Server       ServerConfig   `yaml:"server" mapstructure:"server"`
	Database     DatabaseConfig `yaml:"database" mapstructure:"database"`
}

type ServerConfig struct {
	Port         int    `yaml:"port" mapstructure:"port"`
	Host         string `yaml:"host" mapstructure:"host"`
	ReadTimeout  int    `yaml:"read_timeout" mapstructure:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout" mapstructure:"write_timeout"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host" mapstructure:"host"`
	Port     int    `yaml:"port" mapstructure:"port"`
	User     string `yaml:"user" mapstructure:"user"`
	Password string `yaml:"password" mapstructure:"password"`
	DBName   string `yaml:"dbname" mapstructure:"dbname"`
	SSLMode  string `yaml:"sslmode" mapstructure:"sslmode"`
}

// Load loads configuration from YAML file with fallback to environment variables
func Load() *Config {
	config := &Config{}

	// Try to load from YAML first
	LoadConfigYml(config, "config", "configs")

	// Set defaults for any missing values
	setDefaults(config)

	return config
}

// LoadConfigYml loads configuration from YAML file
func LoadConfigYml(configuration interface{}, configName string, configPath string) {
	// Get the project root directory
	projectRoot, err := getProjectRoot()
	if err != nil {
		printRed(fmt.Sprintf("Could not find project root: %v", err))
		return
	}

	// Build the full path from project root
	fileName := fmt.Sprintf("%s.yaml", configName)
	filePath := filepath.Join(projectRoot, configPath, fileName)

	if _, err := os.Stat(filePath); err != nil {
		printRed(fmt.Sprintf("Config file not found: %s (will use defaults/environment variables)", filePath))
		return
	}

	printGreen(fmt.Sprintf("Loading config %s", filePath))

	viper.SetConfigFile(filePath)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	if err := viper.Unmarshal(configuration); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}
}

// setDefaults sets default values for any missing configuration
func setDefaults(config *Config) {
	if config.DatabaseURL == "" {
		config.DatabaseURL = "postgres://user:password@localhost/payslip_db?sslmode=disable"
	}

	if config.JWTSecret == "" {
		config.JWTSecret = "your-secret-key"
	}

	if config.Environment == "" {
		config.Environment = "development"
	}

	if config.LogLevel == "" {
		config.LogLevel = "info"
	}

	// Server defaults
	if config.Server.Port == 0 {
		config.Server.Port = 8080
	}

	if config.Server.Host == "" {
		config.Server.Host = "localhost"
	}

	if config.Server.ReadTimeout == 0 {
		config.Server.ReadTimeout = 30
	}

	if config.Server.WriteTimeout == 0 {
		config.Server.WriteTimeout = 30
	}

	// Database defaults
	if config.Database.Host == "" {
		config.Database.Host = "localhost"
	}

	if config.Database.Port == 0 {
		config.Database.Port = 5432
	}

	if config.Database.User == "" {
		config.Database.User = "user"
	}

	if config.Database.Password == "" {
		config.Database.Password = "password"
	}

	if config.Database.DBName == "" {
		config.Database.DBName = "payslip_db"
	}

	if config.Database.SSLMode == "" {
		config.Database.SSLMode = "disable"
	}
}

// getProjectRoot finds the project root by looking for go.mod
func getProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("could not find go.mod file")
}

// printYellow prints text in yellow color (simple implementation)
// You can replace this with your loghelper.PrintYellowf function
func printYellow(text string) {
	fmt.Printf("\033[33m%s\033[0m\n", text)
}

// printRed prints text in red color
func printRed(text string) {
	fmt.Printf("\033[31m%s\033[0m\n", text)
}

// Additional color functions for completeness
func printGreen(text string) {
	fmt.Printf("\033[32m%s\033[0m\n", text)
}

func printBlue(text string) {
	fmt.Printf("\033[34m%s\033[0m\n", text)
}

func printMagenta(text string) {
	fmt.Printf("\033[35m%s\033[0m\n", text)
}

func printCyan(text string) {
	fmt.Printf("\033[36m%s\033[0m\n", text)
}
