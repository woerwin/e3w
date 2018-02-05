package routers

import (
	"testing"

	"fmt"

	"github.com/Guazi-inc/e3w/conf"
	"github.com/Guazi-inc/e3w/e3ch"
	"github.com/coreos/etcd/clientv3"
)

func oneTestScope(f func(client *clientv3.Client)) {
	config, err := conf.Init("../conf/config.dev.ini")
	if err != nil {
		panic(err)
	}

	fmt.Println(config)
	client, err := e3ch.NewE3chClient(config)
	if err != nil {
		panic(err)
	}
	clt, err := e3ch.CloneE3chClient(config.EtcdUsername, config.EtcdPassword, client)
	if err != nil {
		panic(err)
	}
	f(clt.EtcdClient())
}

func TestUserAdd(t *testing.T) {
	r := &createUserRequest{
		Name: "u4",
	}

	oneTestScope(func(client *clientv3.Client) {
		fmt.Println(client.Endpoints())
		//fmt.Println(client.Endpoints())

		resp, err := client.UserAdd(newEtcdCtx(), r.Name, r.Password)
		t.Log(err)
		t.Log(resp)
	})

}
