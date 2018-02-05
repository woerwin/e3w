package conf

import "testing"

func TestNewInit(t *testing.T) {
	cfg, err := NewInit("../conf/conf.yml")
	if err != nil {
		t.Error(err)
	}
	t.Log(cfg)
	t.Log(cfg.Etcd)
}
