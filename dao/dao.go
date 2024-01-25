package dao

import "github.com/pkg/errors"

// CloseDB 关闭数据库连接
func CloseDB() error {
	err := DB.Close()
	if err != nil {
		return errors.Wrap(err, "Database operation failed")
	}
	return nil
}

// CloseRedis 关闭 Redis 连接
func CloseRedis() error {
	err := Client.Close()
	if err != nil {
		return errors.Wrap(err, "Redis operation failed")
	}
	return nil
}

// Close 主关闭函数
func Close() error {
	// 尝试关闭数据库连接
	if err := CloseDB(); err != nil {
		logger.Error(err)
		return err
	}

	// 尝试关闭 Redis 连接
	if err := CloseRedis(); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}
