package docker

import "github.com/docker/docker/client"

var (
	LocalClient *DockerClient
	Client      *client.Client
)

// DockerClient represents a client for interacting with the Docker API.
type DockerClient struct {
	client client.APIClient
}

// NewClient creates a new Docker client and returns a pointer to the DockerClient struct.
// It uses the client.NewClientWithOpts function to create the client with default options.
// If an error occurs during client creation, it returns nil and the error.
func NewClient() (*DockerClient, *client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, nil, err
	}
	return &DockerClient{client: cli}, cli, nil
}

// Close closes the Docker client connection.
// It returns an error if there was a problem closing the connection.
func (d *DockerClient) Close() error {
	return d.client.Close()
}

func init() {
	var err error
	LocalClient, Client, err = NewClient()
	if err != nil {
		panic(err)
	}
}
