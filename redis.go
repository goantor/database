package database

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type RedisOptions map[string]*RedisOption

type RedisOption struct {
	Host     string `json:"host" toml:"host" yaml:"host"`
	Port     int    `json:"port" toml:"port" yaml:"port"`
	Password string `json:"password" toml:"password" yaml:"password"`
	Database int    `json:"database" toml:"database" yaml:"database"`
}

func (r *RedisOption) Link() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

func (r *RedisOption) MakeSource() (*redis.Client, error) {
	return NewRedis(r)
}

type Redis struct {
	entity *redis.Client
	opt    *RedisOption
}

func NewRedis(opt *RedisOption) (*redis.Client, error) {
	rd := &Redis{opt: opt}
	if err := rd.Boot(); err != nil {
		return nil, err
	}

	return rd.entity, nil
}

func (r *Redis) Boot() (err error) {
	r.entity = redis.NewClient(&redis.Options{
		Addr:     r.opt.Link(),
		Password: r.opt.Password,
		DB:       r.opt.Database,
	})

	_, err = r.entity.Ping(context.TODO()).Result()
	return
}

func (r *Redis) Close() {
	_ = r.entity.Close()
}
