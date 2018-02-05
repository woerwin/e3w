package routers

import (
	"fmt"
	"strings"
	"time"

	"github.com/Guazi-inc/e3ch"
	"github.com/Guazi-inc/e3w/conf"
	"github.com/coreos/etcd/clientv3"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

const (
	ETCD_USERNAME_HEADER = "X-Etcd-Username"
	ETCD_PASSWORD_HEADER = "X-Etcd-Password"
	ETCD_ENV_HEADER      = "X-Etcd-Env"
	webAuthUserKey       = "/web_auth_users"
)

type e3chHandler func(*gin.Context, *client.EtcdHRCHYClient) (interface{}, error)

type groupHandler func(e3chHandler) respHandler

func withE3chGroup(e3chClt *client.EtcdHRCHYClient, config *conf.NewConfig) groupHandler {
	return func(h e3chHandler) respHandler {
		return func(c *gin.Context) (interface{}, error) {
			//clt := e3chClt
			//if config.Auth {
			//	username := c.Request.Header.Get(ETCD_USERNAME_HEADER)
			//	password := c.Request.Header.Get(ETCD_PASSWORD_HEADER)
			//	clt, err := e3ch.CloneE3chClient(username, password, e3chClt)
			//	if err != nil {
			//		return nil, err
			//	}
			//	defer clt.EtcdClient().Close()
			//}
			clt, err := NewE3chClient(c, e3chClt, config)
			if err != nil {
				return nil, err
			}
			defer clt.EtcdClient().Close()
			return h(c, clt)
		}
	}
}

type etcdHandler func(*gin.Context, *clientv3.Client) (interface{}, error)

func etcdWrapper(h etcdHandler) e3chHandler {
	return func(c *gin.Context, e3chClt *client.EtcdHRCHYClient) (interface{}, error) {
		return h(c, e3chClt.EtcdClient())
	}
}

func InitRouters(g *gin.Engine, config *conf.NewConfig, e3chClt *client.EtcdHRCHYClient) {
	g.GET("/", func(c *gin.Context) {
		c.File("./static/dist/index.html")
	})
	g.Static("/public", "./static/dist")

	e3chGroup := withE3chGroup(e3chClt, config)

	// key/value actions
	g.GET("/kv/*key", resp(e3chGroup(getKeyHandler)))
	g.POST("/kv/*key", resp(e3chGroup(postKeyHandler)))
	g.PUT("/kv/*key", resp(e3chGroup(putKeyHandler)))
	g.DELETE("/kv/*key", resp(e3chGroup(delKeyHandler)))

	// members actions
	g.GET("/members", resp(e3chGroup(etcdWrapper(getMembersHandler))))

	// roles actions
	g.GET("/roles", resp(e3chGroup(etcdWrapper(getRolesHandler))))
	g.POST("/role", resp(e3chGroup(etcdWrapper(createRoleHandler))))
	g.GET("/role/:name", resp(e3chGroup(getRolePermsHandler)))
	g.DELETE("/role/:name", resp(e3chGroup(etcdWrapper(deleteRoleHandler))))
	g.POST("/role/:name/permission", resp(e3chGroup(createRolePermHandler)))
	g.DELETE("/role/:name/permission", resp(e3chGroup(deleteRolePermHandler)))

	// users actions
	g.GET("/users", resp(e3chGroup(etcdWrapper(getUsersHandler))))
	g.POST("/user", resp(e3chGroup(etcdWrapper(createUserHandler))))
	g.GET("/user/:name", resp(e3chGroup(etcdWrapper(getUserRolesHandler))))
	g.DELETE("/user/:name", resp(e3chGroup(etcdWrapper(deleteUserHandler))))
	g.PUT("/user/:name/password", resp(e3chGroup(etcdWrapper(setUserPasswordHandler))))
	g.PUT("/user/:name/role/:role", resp(e3chGroup(etcdWrapper(grantUserRoleHandler))))
	g.DELETE("/user/:name/role/:role", resp(e3chGroup(etcdWrapper(revokeUserRoleHandler))))

	g.GET("/staff/me", resp(Me))
	g.POST("/staff/logout", resp(LogOut))
	g.GET("/envs", resp(GetEnvs))
}

func NewE3chClient(c *gin.Context, e3chClt *client.EtcdHRCHYClient, config *conf.NewConfig) (*client.EtcdHRCHYClient, error) {

	env := c.Request.Header.Get(ETCD_ENV_HEADER)
	cfg, ok := config.EtcdMap[env]
	if !ok {
		return nil, errors.Errorf("config not found, ENV: %s", env)
	}
	username := c.Request.Header.Get(ETCD_USERNAME_HEADER)
	e3Cfg := clientv3.Config{
		Endpoints:   cfg.EndPoints,
		DialTimeout: 3 * time.Second,
	}
	if cfg.Auth {
		//e3Cfg.Username = strconv.Itoa(userID)
		//e3Cfg.Password = DEFAULT_PASSWORD
		e3Cfg.Username = username
		e3Cfg.Password = c.Request.Header.Get(ETCD_PASSWORD_HEADER)
	}
	clt, err := clientv3.New(e3Cfg)
	if err != nil {
		return nil, err
	}
	if cfg.WebAuth {
		enable, err := isWebAuthUser(clt, username)
		if err != nil {
			return nil, err
		}
		if !enable {
			return nil, errors.New("permission denied")
		}
	}
	return e3chClt.Clone(clt), nil
}

func isWebAuthUser(clt *clientv3.Client, username string) (bool, error) {
	if username == "" {
		return false, errors.New("empty username")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := clt.Get(ctx, webAuthUserKey)
	cancel()
	if err != nil {
		return false, err
	}
	if len(resp.Kvs) == 0 {
		return false, errors.New("empty web auth users")
	}
	fmt.Printf("username: %s, users: %+v\n", username, strings.Split(string(resp.Kvs[0].Value), ","))
	for _, user := range strings.Split(string(resp.Kvs[0].Value), ",") {
		if username == user {
			return true, nil
		}
	}
	return false, nil
}
