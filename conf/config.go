package conf

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/ini.v1"
	"gopkg.in/yaml.v2"
)

var MainConfig *NewConfig

type EtcdConfig struct {
	Env       string   `yaml:"env"`
	Auth      bool     `yaml:"auth"`
	RootKey   string   `yaml:"root_key"`
	EndPoints []string `yaml:"addr"`
	DirValue  string
	WebAuth   bool `yaml:"web_auth"`
}

type NewConfig struct {
	App struct {
		Port string `yaml:"port"`
	} `yaml:"app"`
	Etcd    []*EtcdConfig `yaml:"etcd"`
	EtcdMap map[string]*EtcdConfig
}

type Config struct {
	Port          string
	Auth          bool
	EtcdRootKey   string
	DirValue      string
	EtcdEndPoints []string
	EtcdUsername  string
	EtcdPassword  string
	CertFile      string
	KeyFile       string
	CAFile        string
}

func Init(filepath string) (*Config, error) {
	cfg, err := ini.Load(filepath)
	if err != nil {
		return nil, err
	}

	c := &Config{}

	appSec := cfg.Section("app")
	c.Port = appSec.Key("port").Value()
	c.Auth = appSec.Key("auth").MustBool()

	etcdSec := cfg.Section("etcd")
	c.EtcdRootKey = etcdSec.Key("root_key").Value()
	c.DirValue = etcdSec.Key("dir_value").Value()
	c.EtcdEndPoints = etcdSec.Key("addr").Strings(",")
	c.EtcdUsername = etcdSec.Key("username").Value()
	c.EtcdPassword = etcdSec.Key("password").Value()
	c.CertFile = etcdSec.Key("cert_file").Value()
	c.KeyFile = etcdSec.Key("key_file").Value()
	c.CAFile = etcdSec.Key("ca_file").Value()

	return c, nil
}

func NewInit(path string) (*NewConfig, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg NewConfig
	err = yaml.Unmarshal(buf, &cfg)
	if err != nil {
		return nil, err
	}
	if len(cfg.Etcd) == 0 {
		return nil, errors.New("missing etcd config")
	}
	cfg.EtcdMap = map[string]*EtcdConfig{}
	for _, ey := range cfg.Etcd {
		cfg.EtcdMap[ey.Env] = ey
	}
	MainConfig = &cfg
	return MainConfig, nil
}
