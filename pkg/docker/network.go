package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
)

type NetworkCreateOpts struct {
	Name   string
	Subnet string
}

// NetworkCreate creates a new network in Docker with the specified options.
// It returns an error if the network creation fails.
func (d *DockerClient) NetworkCreate(opts NetworkCreateOpts) (string, error) {
	resp, err := d.client.NetworkCreate(context.Background(), opts.Name, types.NetworkCreate{
		Driver: "bridge",
		IPAM: &network.IPAM{
			Config: []network.IPAMConfig{
				{
					Subnet: opts.Subnet,
				},
			},
		},
		Options: map[string]string{
			"com.docker.network.bridge.name": opts.Name + "-f5gc",
		},
		Labels: map[string]string{"mosaic": "true"},
	})
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}

// NetworkRemove removes a network by its name.
func (d *DockerClient) NetworkRemove(name string) error {
	return d.client.NetworkRemove(context.Background(), name)
}

func (d *DockerClient) PruneNetworks() error {
	filter := filters.NewArgs()
	filter.Add("label", "mosaic=true")
	_, err := d.client.NetworksPrune(context.Background(), filter)
	return err
}

// NetworkExists checks if a network with the given name exists.
// It returns true if the network exists, false otherwise.
// An error is returned if there was a problem listing the networks.
func (d *DockerClient) NetworkExists(name string) (bool, error) {
	networks, err := d.client.NetworkList(context.Background(), types.NetworkListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{Key: "name", Value: name}),
	})
	if err != nil {
		return false, err
	}
	for _, n := range networks {
		if n.Name == name {
			return true, nil
		}
	}
	return false, nil
}

// NetworkNameToID retrieves the ID of a Docker network based on its name.
// It takes a network name as input and returns the corresponding network ID.
// If the network is not found, it returns an empty string.
// If an error occurs during the network retrieval, it returns the error.
func (d *DockerClient) NetworkNameToID(name string) (string, error) {
	networks, err := d.client.NetworkList(context.Background(), types.NetworkListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{Key: "name", Value: name}),
	})
	if err != nil {
		return "", err
	}
	if len(networks) == 0 {
		return "", nil
	}
	return networks[0].ID, nil
}
