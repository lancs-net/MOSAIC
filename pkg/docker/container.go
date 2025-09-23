package docker

import (
	"context"
	"slices"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
)

type ContainerCreateOpts struct {
	Name             string
	Image            string
	Binds            []string
	Cmd              []string
	Env              []string
	ExposedPorts     nat.PortSet
	CapAdd           []string
	Net              string
	NetworkAlias     string
	NetworkIPAddress string
	PortBindings     nat.PortMap
	Device           string
}

var (
	containerStopTimeout int = 5
)

func ContainerList() ([]string, error) {
	containers, err := Client.ContainerList(context.Background(), container.ListOptions{
		All: true,
	})
	if err != nil {
		return nil, err
	}
	var names []string
	for _, c := range containers {
		names = append(names, c.Names[0])
	}
	return names, nil
}

// ContainerCreate creates a new container in Docker with the specified options.
func (d *DockerClient) ContainerCreate(opts ContainerCreateOpts) (string, error) {
	hostConfig := container.HostConfig{
		Binds:        opts.Binds,
		CapAdd:       opts.CapAdd,
		PortBindings: opts.PortBindings,
	}
	if opts.Device != "" {
		hostConfig.Devices = []container.DeviceMapping{{
			PathOnHost:        opts.Device,
			PathInContainer:   opts.Device,
			CgroupPermissions: "rwm",
		}}
	}

	// Resources: container.Resources{
	// 	Devices: []container.DeviceMapping{{PathInContainer: opts.Device}}},

	resp, err := d.client.ContainerCreate(context.Background(),
		&container.Config{
			Image:        opts.Image,
			Cmd:          opts.Cmd,
			Env:          opts.Env,
			ExposedPorts: opts.ExposedPorts,
			Labels:       map[string]string{"auto5gc": "true"},
		},
		&hostConfig,
		&network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				opts.Net: {
					Aliases:   []string{opts.NetworkAlias},
					IPAddress: opts.NetworkIPAddress,
				},
			},
		},
		nil,
		opts.Name,
	)
	if err != nil {
		return "", err
	}
	// io.Copy(os.Stdout, response.Body)
	return resp.ID, nil
}

// ContainerStart starts a Docker container with the specified ID.
func (d *DockerClient) ContainerStart(id string) error {
	return d.client.ContainerStart(context.Background(), id, container.StartOptions{})
}

// ContainerStop stops a Docker container with the specified ID.
func (d *DockerClient) ContainerStop(id string) error {
	return d.client.ContainerStop(context.Background(), id, container.StopOptions{
		Timeout: &containerStopTimeout,
	})
}

// ContainerRemove removes a Docker container by its ID.
func (d *DockerClient) ContainerRemove(id string) error {
	return d.client.ContainerRemove(context.Background(), id, container.RemoveOptions{
		Force: true,
	})
}

func (d *DockerClient) PruneContainers() error {
	filter := filters.NewArgs()
	filter.Add("label", "auto5gc=true")
	_, err := d.client.ContainersPrune(context.Background(), filter)
	return err
}

// ContainerNameToID returns the ID of a container given its name.
// It searches for the container by name and returns the first matching ID.
// If no matching container is found, it returns an empty string.
func (d *DockerClient) ContainerNameToID(name string) (string, error) {
	containers, err := d.client.ContainerList(context.Background(), container.ListOptions{
		All:     true,
		Filters: filters.NewArgs(filters.KeyValuePair{Key: "name", Value: name}),
	})
	if err != nil {
		return "", err
	}
	for _, c := range containers {
		for _, n := range c.Names {
			if n == "/"+name {
				return c.ID, nil
			}
		}
	}
	return "", nil
}

func (d *DockerClient) ContainerStatus(name string) (string, error) {
	containers, err := d.client.ContainerList(context.Background(), container.ListOptions{
		All:     true,
		Filters: filters.NewArgs(filters.KeyValuePair{Key: "name", Value: name}),
	})
	if err != nil {
		return "", err
	}
	for _, c := range containers {
		if slices.Contains(c.Names, "/"+name) {
			return c.Status, nil
		}
	}
	return "", nil
}

// func (d *DockerClient) GetConPid(id string) (types.ContainerJSON, error) {
// 	return d.client.ContainerInspect(context.Background(), id)
// }

//logstash!!!
