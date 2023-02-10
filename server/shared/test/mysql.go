package test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var mysqlDSN string

// RunWithMySQLInDocker runs the tests with
// a MySQL instance in a docker container.
func RunWithMySQLInDocker(m *testing.M) int {
	c, err := client.NewClientWithOpts()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	resp, err := c.ContainerCreate(ctx, &container.Config{
		Image: consts.MySQLImage,
		ExposedPorts: nat.PortSet{
			consts.MySQLContainerPort: {},
		},
	}, &container.HostConfig{
		PortBindings: nat.PortMap{
			consts.MySQLContainerPort: []nat.PortBinding{
				{
					HostIP:   consts.MySQLContainerIP,
					HostPort: consts.MySQLPort,
				},
			},
		},
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}
	containerID := resp.ID
	defer func() {
		err := c.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
			Force: true,
		})
		if err != nil {
			panic(err)
		}
	}()

	err = c.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
	if err != nil {
		panic(err)
	}

	inspRes, err := c.ContainerInspect(ctx, containerID)
	if err != nil {
		panic(err)
	}
	hostPort := inspRes.NetworkSettings.Ports[consts.MySQLContainerPort][0]
	port, _ := strconv.Atoi(hostPort.HostPort)
	mysqlDSN = fmt.Sprintf(consts.MySqlDSN, consts.MySQLAdmin, "", hostPort.HostIP, port, consts.TikGok)

	return m.Run()
}

func NewTestDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(mysqlDSN), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}
	return db
}

func SetupDatabase() {
}
