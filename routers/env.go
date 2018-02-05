package routers

import (
	"github.com/Guazi-inc/e3w/conf"
	"github.com/gin-gonic/gin"
)

func GetEnvs(c *gin.Context) (interface{}, error) {
	var envs []string
	for _, ey := range conf.MainConfig.Etcd {
		envs = append(envs, ey.Env)
	}
	return envs, nil
}
