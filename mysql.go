package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

type Database interface {
	Link() string
	Boot() error
	Close() error
}

type MysqlOptions map[string]*MysqlOption

type MysqlOption struct {
	Debug       bool   `json:"debug" toml:"debug" yaml:"debug"`
	Host        string `json:"host" toml:"host" yaml:"host"`
	Port        int    `json:"port" toml:"port" yaml:"port"`
	User        string `json:"user" toml:"user" yaml:"user"`
	Password    string `json:"password" toml:"password" yaml:"password"`
	Database    string `json:"database" toml:"database" yaml:"database"`
	Charset     string `json:"charset" toml:"charset" yaml:"charset"`
	MaxIdleConn int    `json:"max_idle_conn" toml:"max_idle_conn" yaml:"max_idle_conn"`
	MaxOpenConn int    `json:"max_open_conn" toml:"max_open_conn" yaml:"max_open_conn"`
	MaxLifeTime int    `json:"max_life_time" toml:"max_life_time" yaml:"max_life_time"`
}

func (m *MysqlOption) Link() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		m.User, m.Password, m.Host, m.Port, m.Database, m.Charset)
}

func (m *MysqlOption) MakeSource() (*gorm.DB, error) {
	return NewMysql(m)
}

type Mysql struct {
	entity *gorm.DB
	opt    *MysqlOption
}

func NewMysql(opt *MysqlOption) (db *gorm.DB, err error) {
	mysql := &Mysql{opt: opt}
	if err = mysql.Boot(); err != nil {
		return
	}

	return mysql.entity, nil
}

func (m *Mysql) Boot() (err error) {
	m.entity, err = gorm.Open(mysql.Open(m.opt.Link()), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		return
	}

	if m.opt.Debug {
		m.entity.Debug()
	}

	m.pool()
	return
}

func (m *Mysql) pool() {
	pool, err := m.entity.DB()
	if err != nil {
		panic(fmt.Sprintf("mysql get pool failed: %+v\n", err))
	}

	pool.SetMaxIdleConns(m.opt.MaxIdleConn)
	pool.SetMaxOpenConns(m.opt.MaxOpenConn)
	pool.SetConnMaxLifetime(time.Duration(m.opt.MaxLifeTime) * time.Second)
}

func (m *Mysql) Close() {
}
