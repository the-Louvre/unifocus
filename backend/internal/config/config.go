package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Config 应用配置结构
type Config struct {
	Server     ServerConfig     `yaml:"server"`
	Database   DatabaseConfig   `yaml:"database"`
	Redis      RedisConfig      `yaml:"redis"`
	JWT        JWTConfig        `yaml:"jwt"`
	Crawler    CrawlerConfig    `yaml:"crawler"`
	NLPService NLPServiceConfig `yaml:"nlp_service"`
	Log        LogConfig        `yaml:"log"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         int    `yaml:"port"`
	Mode         string `yaml:"mode"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	User            string `yaml:"user"`
	Password        string `yaml:"password"`
	DBName          string `yaml:"dbname"`
	SSLMode         string `yaml:"sslmode"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
}

// GetDSN 返回PostgreSQL连接字符串
func (d *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.DBName, d.SSLMode,
	)
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	PoolSize int    `yaml:"pool_size"`
}

// GetAddr 返回Redis地址
func (r *RedisConfig) GetAddr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret      string `yaml:"secret"`
	ExpireHours int    `yaml:"expire_hours"`
}

// GetExpireDuration 返回过期时间
func (j *JWTConfig) GetExpireDuration() time.Duration {
	return time.Duration(j.ExpireHours) * time.Hour
}

// CrawlerConfig 爬虫配置
type CrawlerConfig struct {
	WorkerCount    int       `yaml:"worker_count"`
	RequestTimeout int       `yaml:"request_timeout"`
	UserAgents     []string  `yaml:"user_agents"`
	RateLimit      RateLimit `yaml:"rate_limit"`
}

// RateLimit 频率限制配置
type RateLimit struct {
	RequestsPerSecond float64 `yaml:"requests_per_second"`
	Burst             int     `yaml:"burst"`
}

// NLPServiceConfig NLP服务配置
type NLPServiceConfig struct {
	URL     string `yaml:"url"`
	Timeout int    `yaml:"timeout"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `yaml:"level"`
	Output     string `yaml:"output"`
	FilePath   string `yaml:"file_path"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
}

// globalConfig 全局配置实例
// 注意: 此变量在Load()函数中设置，在Get()函数中读取
// 并发安全: Load()通常在应用启动时调用一次，之后只读，因此不需要加锁
// 如果需要在运行时动态更新配置，需要添加sync.RWMutex保护
var globalConfig *Config

// Load 加载配置文件
// 配置加载优先级:
// 1. 如果configPath为空，根据环境变量APP_ENV选择配置文件（默认dev）
// 2. 读取YAML配置文件
// 3. 使用os.ExpandEnv解析环境变量（支持${VAR}格式）
// 4. 解析YAML内容到Config结构体
// 5. 验证配置有效性
// 6. 设置全局配置实例
func Load(configPath string) (*Config, error) {
	// 如果路径为空，根据环境变量选择配置文件
	if configPath == "" {
		env := os.Getenv("APP_ENV")
		if env == "" {
			env = "dev"
		}
		configPath = filepath.Join("configs", fmt.Sprintf("config.%s.yaml", env))
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// 解析环境变量（支持${VAR}格式，如 ${DATABASE_PASSWORD}）
	expandedData := os.ExpandEnv(string(data))

	// 解析YAML
	var cfg Config
	if err := yaml.Unmarshal([]byte(expandedData), &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// 验证配置
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	globalConfig = &cfg
	return &cfg, nil
}

// Get 获取全局配置
func Get() *Config {
	if globalConfig == nil {
		panic("config not loaded, call Load() first")
	}
	return globalConfig
}

// validate 验证配置有效性
func (c *Config) validate() error {
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	if c.Database.Host == "" {
		return fmt.Errorf("database host cannot be empty")
	}

	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT secret cannot be empty")
	}

	return nil
}
