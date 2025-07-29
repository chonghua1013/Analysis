package database

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

// MySQLPool 封装MySQL连接池
type MySQLPool struct {
	db     *sql.DB
	config *MySQLConfig
	logger *zap.Logger
	mu     sync.Mutex
}

// MySQLConfig MySQL配置结构
type MySQLConfig struct {
	Username        string        `yaml:"username" json:"username"`
	Password        string        `yaml:"password" json:"password"`
	Host            string        `yaml:"host" json:"host"`
	Port            string        `yaml:"port" json:"port"`
	Database        string        `yaml:"database" json:"database"`
	MaxOpenConns    int           `yaml:"max_open_conns" json:"max_open_conns"`         // 最大打开连接数
	MaxIdleConns    int           `yaml:"max_idle_conns" json:"max_idle_conns"`         // 最大空闲连接数
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" json:"conn_max_lifetime"`   // 连接最大存活时间
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time" json:"conn_max_idle_time"` // 连接最大空闲时间
}

// DefaultConfig 返回默认配置
func DefaultConfig() *MySQLConfig {
	return &MySQLConfig{
		MaxOpenConns:    50,
		MaxIdleConns:    20,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 1 * time.Minute,
	}
}

// NewMySQLPool 创建新的MySQL连接池
func NewMySQLPool(config *MySQLConfig, logger *zap.Logger) (*MySQLPool, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=5s",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open mysql connection: %w", err)
	}

	// 设置连接池参数
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)
	db.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	// 测试连接
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping mysql: %w", err)
	}

	logger.Info("MySQL connection pool initialized",
		zap.String("host", config.Host),
		zap.String("database", config.Database),
	)

	return &MySQLPool{
		db:     db,
		config: config,
		logger: logger,
	}, nil
}

// GetDB 获取数据库连接
func (p *MySQLPool) GetDB() *sql.DB {
	return p.db
}

// Close 关闭连接池
func (p *MySQLPool) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.db != nil {
		err := p.db.Close()
		if err != nil {
			return fmt.Errorf("failed to close mysql connection: %w", err)
		}
		p.db = nil
		p.logger.Info("MySQL connection pool closed")
	}
	return nil
}

// Stats 返回连接池状态
func (p *MySQLPool) Stats() sql.DBStats {
	return p.db.Stats()
}

// WithTransaction 执行事务操作
func (p *MySQLPool) WithTransaction(fn func(tx *sql.Tx) error) error {
	tx, err := p.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p) // 重新抛出panic
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction failed: %v, rollback also failed: %w", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
