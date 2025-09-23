package f5gc

import (
	"fmt"
	"strconv"

	"github.com/lancs-net/MOSAIC/pkg/common"
	"github.com/lancs-net/MOSAIC/pkg/docker"
)

var ContainerID = map[string]string{}

func nfBase(nf string, path string) docker.ImageBuildOpts {
	var opts docker.ImageBuildOpts
	opts.BuildArgs = map[string]*string{"F5GC_MODULE": &nf}
	opts.Context = path
	opts.Dockerfile = "Dockerfile.nf"
	tag := "free5gc/" + nf + "-base:latest"
	opts.Tags = []string{tag}
	return opts
}

func Remove(d *docker.DockerClient, function string) error {
	if function == "base" {
		imgID, err := d.ImageNameToID("free5gc/" + function + ":latest")
		fmt.Println("Removing Image: free5gc/" + function + ":latest")
		if err != nil {
			return err
		}
		err = d.ImageRemove(imgID)
		if err != nil {
			return err
		}
		return nil
	}
	imgID, err := d.ImageNameToID("free5gc/" + function + "-base:latest")
	if err != nil {
		return err
	}
	err = d.ImageRemove(imgID)
	if err != nil {
		return err
	}
	return nil
}

func GetConf(function common.NFInfo, path string, num int) (docker.ImageBuildOpts, docker.ContainerCreateOpts) {
	switch function.Name {
	case "upf":
		imgOpts, conOpts := upfConf(function, path, num)
		return imgOpts, conOpts
	case "nrf":
		imgOpts, conOpts := commonConf(function, path, num)
		return imgOpts, conOpts
	case "amf":
		imgOpts, conOpts := commonConf(function, path, num)
		return imgOpts, conOpts
	case "ausf":
		imgOpts, conOpts := commonConf(function, path, num)
		return imgOpts, conOpts
	case "nssf":
		imgOpts, conOpts := commonConf(function, path, num)
		return imgOpts, conOpts
	case "pcf":
		imgOpts, conOpts := commonConf(function, path, num)
		return imgOpts, conOpts
	case "smf":
		imgOpts, conOpts := smfConf(function, path, num)
		return imgOpts, conOpts
	case "udm":
		imgOpts, conOpts := commonConf(function, path, num)
		return imgOpts, conOpts
	case "udr":
		imgOpts, conOpts := commonConf(function, path, num)
		return imgOpts, conOpts
	case "chf":
		imgOpts, conOpts := commonConf(function, path, num)
		return imgOpts, conOpts
	case "n3iwf":
		imgOpts, conOpts := n3iwfConf(function, path, num)
		return imgOpts, conOpts
	case "ueransim":
		imgOpts, conOpts := ueransimConf(function, path, num)
		return imgOpts, conOpts
	}
	return docker.ImageBuildOpts{}, docker.ContainerCreateOpts{}
}

func Stop(d *docker.DockerClient, function string) (bool, error) {
	conID, err := d.ContainerNameToID(function)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	err = d.ContainerStop(conID)
	if err != nil {
		fmt.Println(err)
		return true, err
	}
	return true, nil
}

func StopAll(d *docker.DockerClient, nfinfo []common.NFInfo) (bool, error) {
	funcExists, err := Stop(d, "mongodb")
	if err != nil {
		return funcExists, err
	}
	funcExists, err = Stop(d, "webui")
	if err != nil {
		return funcExists, err
	}

	for i := range nfinfo {
		for j := 1; j <= nfinfo[i].NumInstances; j++ {
			funcExists, err := Stop(d, nfinfo[i].Name+strconv.Itoa(j))
			if err != nil {
				return funcExists, err
			}
		}
		fmt.Println(nfinfo[i].Name + " Stop Unsucessful...")
	}
	return true, nil
}

func Start(d *docker.DockerClient, function string) (bool, error) {
	conID, err := d.ContainerNameToID(function)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	err = d.ContainerStart(conID)
	if err != nil {
		fmt.Println(function + " Container Start Failed...")
		fmt.Println(err)
		return true, err
	}
	ContainerID[function] = conID
	fmt.Println(function + " Container Start Successful...")
	return true, nil
}

func StartAll(d *docker.DockerClient, nfinfo []common.NFInfo) (bool, error) {
	for i := range nfinfo {
		for j := 1; j <= nfinfo[i].NumInstances; j++ {
			funcExists, err := Start(d, nfinfo[i].Name+strconv.Itoa(j))
			if err != nil {
				return funcExists, err
			}
		}
		fmt.Println(nfinfo[i].Name + " Start Unsucessful...")
	}
	return true, nil
}

func RemoveAllCon(d *docker.DockerClient, nfinfo []common.NFInfo) (bool, error) {
	funcExists, err := RemoveCon(d, "mongodb")
	if err != nil {
		return funcExists, err
	}
	funcExists, err = RemoveCon(d, "webui")
	if err != nil {
		return funcExists, err
	}
	for i := range nfinfo {
		for j := 1; j <= nfinfo[i].NumInstances; j++ {
			funcExists, err := RemoveCon(d, nfinfo[i].Name+strconv.Itoa(j))
			if err != nil {
				return funcExists, err
			}
		}
		fmt.Println(nfinfo[i].Name + " Container Remove Unsucessful...")
	}
	return true, nil
}

func RemoveCon(d *docker.DockerClient, function string) (bool, error) {
	fmt.Println("Removing Container: ", function)

	conID, err := d.ContainerNameToID(function)
	if err != nil {
		return false, err
	}
	err = d.ContainerRemove(conID)
	if err != nil {
		return true, err
	}
	// delete(ContainerID, function)
	return true, nil
}
