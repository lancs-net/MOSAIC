package docker

import (
	"context"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
)

func (d *DockerClient) VolumeCreate(name string) error {
	_, err := d.client.VolumeCreate(context.Background(), volume.CreateOptions{Name: name})
	if err != nil {
		return err
	}
	return nil
}

func (d *DockerClient) VolumeRemove(volID string) error {
	return d.client.VolumeRemove(context.Background(), volID, false)
}

func (d *DockerClient) PruneVolumes() error {
	filter := filters.NewArgs()
	filter.Add("label", "auto5gc=true")
	_, err := d.client.VolumesPrune(context.Background(), filter)
	return err
}

func (d *DockerClient) VolumeExists(volID string) error {
	_, _, err := d.client.ImageInspectWithRaw(context.Background(), volID)
	if err != nil {
		return err
	}
	return nil
}
