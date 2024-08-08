package gorm

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var defaultConfigData = []byte(
	`# GORM configuration file

# Database driver (e.g., postgres, mysql, sqlite, sqlserver)
driver: "postgres"
# SQLite file path
sqlite_file: "test.db"
# Database host
host: "localhost"
# Database port
port: 5432
# Database user
user: "postgres"
# Database password
password: "666666"
# Database name
database: "testdb"
# Postgres connection parameters
postgres_params: "TimeZone=Asia/Shanghai"
# MySQL connection parameters
mysql_params: "charset=utf8mb4&parseTime=True&loc=Local"
# 4 # Log level for GORM logger, 1 - Silent, 2 - Error, 3 - Warn, 4 - Info
log_level: 4
`)

type Config struct {
	Driver         string `mapstructure:"driver"`
	SQLiteFile     string `mapstructure:"sqlite_file"`
	Host           string `mapstructure:"host"`
	Port           int    `mapstructure:"port"`
	User           string `mapstructure:"user"`
	LogLevel       int    `mapstructure:"log_level"` //  1 - Silent, 2 - Error, 3 - Warn, 4 - Info
	Password       string `mapstructure:"password"`
	Database       string `mapstructure:"database"`
	PostgresParams string `mapstructure:"postgres_params"`
	MysqlParams    string `mapstructure:"mysql_params"`
}

func (c *Config) NewLogger() logger.Interface {
	return logger.Default.LogMode(logger.LogLevel(c.LogLevel))
}

func (c *Config) Connect() (db *gorm.DB, close func() error) {
	var dialer gorm.Dialector
	switch c.Driver {
	case "postgres":
		dialer = c.postgresDialer()
	case "mysql":
		dialer = c.mysqlDialer()
	case "sqlite":
		dialer = c.sqliteDialer()
	case "sqlserver":
		dialer = c.sqlserverDialer()
	default:
		panic(fmt.Errorf("unsupported database driver: %s", c.Driver))
	}
	var err error
	db, err = gorm.Open(dialer, &gorm.Config{
		Logger: c.NewLogger(),
	})
	if err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	return db, sqlDB.Close
}

func (c *Config) postgresDialer() gorm.Dialector {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s %s",
		c.Host, c.Port, c.User, c.Password, c.Database, c.PostgresParams,
	)
	return postgres.Open(dsn)
}

func (c *Config) mysqlDialer() gorm.Dialector {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?%s",
		c.User, c.Password, c.Host, c.Port, c.Database, c.MysqlParams,
	)
	return mysql.Open(dsn)
}

func (c *Config) sqliteDialer() gorm.Dialector {
	return sqlite.Open(c.SQLiteFile)
}

func (c *Config) sqlserverDialer() gorm.Dialector {
	dsn := fmt.Sprintf(
		"sqlserver://%s:%s@%s:%d?database=%s",
		c.User, c.Password, c.Host, c.Port, c.Database,
	)
	return sqlserver.Open(dsn)
}
