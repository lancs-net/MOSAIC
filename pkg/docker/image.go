package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"slices"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/pkg/archive"
)

type ImageBuildOpts struct {
	Tags       []string
	Dockerfile string
	BuildArgs  map[string]*string
	Context    string
}

// ImagePull pulls the specified image from the Docker registry.
// It returns an error if the image pull operation fails.
func (d *DockerClient) ImagePull(image string) error {
	resp, err := d.client.ImagePull(context.Background(), image, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer resp.Close()
	// io.Copy(os.Stdout, resp)
	io.Copy(io.Discard, resp)
	return nil
}

// ImageBuild builds a Docker image using the provided options.
// It takes an ImageBuildOpts struct as input and returns an error if any.
func (d *DockerClient) ImageBuild(opts ImageBuildOpts) error {
	exist, err := d.ImageExists(opts.Tags[0])
	if exist && err == nil {
		fmt.Println("Image already exists for", opts.Tags[0], "Skipping Build...")
		return nil
	}
	if err == nil && !exist {
		ctx := context.Background()
		fmt.Println("Initiating Build for", opts.Tags)
		buildContext, err := archive.TarWithOptions(opts.Context, &archive.TarOptions{})
		if err != nil {
			fmt.Println("Image Build Context Creation Failed")
			return err
		}
		response, err := d.client.ImageBuild(ctx, buildContext, types.ImageBuildOptions{
			Tags:        opts.Tags,
			Remove:      true,
			ForceRemove: true,
			Dockerfile:  opts.Dockerfile,
			BuildArgs:   opts.BuildArgs,
			Labels: map[string]string{
				"auto5gc": "true",
			},
		})
		if err != nil {
			return err
		}
		defer response.Body.Close()
		// io.Copy(io.Discard, response.Body)
		io.Copy(os.Stdout, response.Body)
		return nil
	} else {
		return err
	}

}

// ImageRemove removes a Docker image.
func (d *DockerClient) ImageRemove(image string) error {
	_, err := d.client.ImageRemove(context.Background(), image, types.ImageRemoveOptions{})
	return err
}

func (d *DockerClient) PruneImages() error {
	filter := filters.NewArgs()
	// filter.Add("label", "auto5gc=true")
	filter.Add("dangling", "true")
	_, err := d.client.ImagesPrune(context.Background(), filter)
	return err
}

// ImageExists checks if the specified image exists in the Docker client.
// It returns an error if the image does not exist or if there was an error
// inspecting the image.
func (d *DockerClient) ImageExists(image string) (bool, error) {
	img_id, err := d.ImageNameToID(image)
	if err != nil {
		fmt.Println("Error while finding corresponding image ID...")
		fmt.Println(err)
		return true, err
	}
	if img_id == "" {
		return false, nil
	}
	return true, nil
}

func (d *DockerClient) ImageNameToID(name string) (string, error) {
	images, err := d.client.ImageList(context.Background(), types.ImageListOptions{
		All: true,
	})
	if err != nil {
		return "", err
	}
	for _, i := range images {
		if slices.Contains(i.RepoTags, name) {
			return i.ID, nil
		}
	}
	return "", nil
}

func (d *DockerClient) GetImageHistory(name string) ([]string, error) {
	history, err := d.client.ImageHistory(context.Background(), name)
	var historyList []string
	if err != nil {
		fmt.Println(err)
		return historyList, err
	}
	for i, h := range history {
		if h.ID != "<missing>" && i != 0 {
			historyList = append(historyList, h.ID)
		}
	}
	slices.Reverse(historyList)
	return historyList, nil
}
