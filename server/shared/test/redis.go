package test

import (
	"context"
	"fmt"

	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/go-redis/redis/v8"
	"testing"
)

var port string

// RunWithRedisInDocker runs the tests with
// a redis instance in a docker container.
func RunWithRedisInDocker(m *testing.M) int {
	c, err := client.NewClientWithOpts(client.WithVersion("1.41"))
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	resp, err := c.ContainerCreate(ctx,
		&container.Config{
			Image: consts.RedisImage,
			ExposedPorts: nat.PortSet{
				consts.RedisContainerPort: {},
			},
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				consts.RedisContainerPort: []nat.PortBinding{
					{
						HostIP:   consts.RedisContainerIP,
						HostPort: consts.RedisPort,
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
	hostPort := inspRes.NetworkSettings.Ports[consts.RedisContainerPort][0]
	port = hostPort.HostPort
	return m.Run()
}

// NewRedisClient creates a client connected to the redis instance in docker.
func NewRedisClient(c context.Context, db int) (*redis.Client, error) {
	if port == "" {
		return nil, fmt.Errorf("redis port not set.Please run RunWithRedisInDocker in TestMain")
	}
	return redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", consts.RedisContainerIP, port),
		DB:   db,
	}), nil
}
