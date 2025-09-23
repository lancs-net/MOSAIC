package f5gc

import (
	"strconv"

	"github.com/docker/go-connections/nat"
	"github.com/lancs-net/MOSAIC/pkg/common"
	"github.com/lancs-net/MOSAIC/pkg/docker"
)

func dbConf() docker.ContainerCreateOpts {
	var conOpts docker.ContainerCreateOpts

	conOpts.Name = "mongodb"
	conOpts.Image = "mongo"
	conOpts.Binds = []string{"dbdata:/data/db"}
	conOpts.Cmd = []string{"mongod", "--port", "27017"}
	conOpts.ExposedPorts = map[nat.Port]struct{}{"27017/tcp": {}}
	conOpts.Net = "ground"
	conOpts.NetworkAlias = "db"

	return conOpts
}

func webuiConf(path string) (docker.ImageBuildOpts, docker.ContainerCreateOpts) {
	var opts docker.ImageBuildOpts
	var conOpts docker.ContainerCreateOpts

	debugTools := "false"

	conOpts.Name = "webui"
	conOpts.Image = "webui:latest"

	conOpts.Cmd = []string{"./webui", "-c", "/free5gc/config/webuicfg.yaml"}
	conOpts.ExposedPorts = map[nat.Port]struct{}{"2121/tcp": {}, "2122/tcp": {}}
	conOpts.Env = []string{"GIN_MODE=release"}

	conOpts.Binds = []string{
		path + "config/webuicfg.yaml:/free5gc/config/webuicfg.yaml",
	}

	conOpts.PortBindings = map[nat.Port][]nat.PortBinding{
		"5000/tcp": {{
			HostIP:   "0.0.0.0",
			HostPort: "5000",
		}},
	}

	conOpts.Net = "ground"
	conOpts.NetworkAlias = "webui"

	opts.BuildArgs = map[string]*string{"DEBUG_TOOLS": &debugTools}
	opts.Dockerfile = "Dockerfile"
	opts.Context = path + "webui"
	opts.Tags = []string{conOpts.Image}

	// depends on db
	return opts, conOpts
}

// nrf, udr (udr depends on nrf as well), chf (depends on nrf and webui as well)
func commonConf(nfinfo common.NFInfo, givpath string, num int) (docker.ImageBuildOpts, docker.ContainerCreateOpts) {
	var opts docker.ImageBuildOpts
	var conOpts docker.ContainerCreateOpts

	debugTools := "false"
	NfName := nfinfo.Name
	path := givpath + "config/" + NfName + "cfg.yaml"

	conOpts.Name = NfName
	conOpts.Image = NfName + ":latest"
	conOpts.Binds = []string{
		path + ":/free5gc/config/" + NfName + "cfg.yaml",
	}
	conOpts.Cmd = []string{"./" + NfName, "-c", "/free5gc/config/" + NfName + "cfg.yaml"}

	conOpts.Net = nfinfo.Network
	conOpts.NetworkAlias = NfName + strconv.Itoa(num) + ".free5gc.org"

	conOpts.ExposedPorts = map[nat.Port]struct{}{"8000/tcp": {}}
	if NfName == "nrf" || NfName == "udr" || NfName == "chf" {
		conOpts.Env = []string{"GIN_MODE=release", "DB_URI=mongodb://db/free5gc"}
	} else {
		conOpts.Env = []string{"GIN_MODE=release"}
	}

	opts.BuildArgs = map[string]*string{"DEBUG_TOOLS": &debugTools}
	opts.Dockerfile = "Dockerfile"
	opts.Context = path + "nf_" + NfName
	opts.Tags = []string{conOpts.Image}

	// depends on db
	return opts, conOpts
}

func upfConf(nfinfo common.NFInfo, path string, num int) (docker.ImageBuildOpts, docker.ContainerCreateOpts) {
	var opts docker.ImageBuildOpts
	var conOpts docker.ContainerCreateOpts

	debugTools := "false"

	opts.BuildArgs = map[string]*string{"DEBUG_TOOLS": &debugTools}
	opts.Context = path + "nf_upf"
	opts.Dockerfile = "Dockerfile"
	opts.Tags = []string{"upf:latest"}

	conOpts.Name = "upf" + strconv.Itoa(num)
	conOpts.Image = "upf:latest"
	conOpts.Binds = []string{
		path + "config/upfcfg.yaml:/free5gc/config/upfcfg.yaml",
		path + "config/upf-iptables.sh:/free5gc/upf-iptables.sh",
	}
	conOpts.Cmd = []string{"bash", "-c", "/free5gc/upf-iptables.sh && ./upf	-c /free5gc/config/upfcfg.yaml"}
	conOpts.CapAdd = []string{"NET_ADMIN"}
	conOpts.Net = nfinfo.Network
	conOpts.NetworkAlias = "upf" + strconv.Itoa(num) + ".free5gc.org"

	return opts, conOpts
}

func smfConf(nfinfo common.NFInfo, path string, num int) (docker.ImageBuildOpts, docker.ContainerCreateOpts) {
	var opts docker.ImageBuildOpts
	var conOpts docker.ContainerCreateOpts

	debugTools := "false"

	conOpts.Name = "smf" + strconv.Itoa(num)
	conOpts.Image = "smf:latest"
	conOpts.Binds = []string{
		path + "config/smfcfg.yaml:/free5gc/config/smfcfg.yaml",
		path + "config/uerouting.yaml:/free5gc/config/uerouting.yaml",
	}
	conOpts.Cmd = []string{"./smf", "-c", "/free5gc/config/smfcfg.yaml", "-u", "/free5gc/config/uerouting.yaml"}
	conOpts.ExposedPorts = map[nat.Port]struct{}{"8000/tcp": {}}
	conOpts.Env = []string{"GIN_MODE=release"}
	conOpts.Net = nfinfo.Network
	conOpts.NetworkAlias = "smf" + strconv.Itoa(num) + ".free5gc.org"

	opts.BuildArgs = map[string]*string{"DEBUG_TOOLS": &debugTools}
	opts.Dockerfile = "Dockerfile"
	opts.Context = path + "nf_smf"
	opts.Tags = []string{conOpts.Image}

	// depends on nrf and upf
	return opts, conOpts
}

func n3iwfConf(nfinfo common.NFInfo, path string, num int) (docker.ImageBuildOpts, docker.ContainerCreateOpts) {
	var opts docker.ImageBuildOpts
	var conOpts docker.ContainerCreateOpts

	debugTools := "false"

	conOpts.Name = "n3iwf" + strconv.Itoa(num)
	conOpts.Image = "n3iwf:latest"

	conOpts.Cmd = []string{"sh", "-c", "/free5gc/n3iwf-ipsec.sh && ./n3iwf -c /free5gc/config/n3iwfcfg.yaml"}
	conOpts.Env = []string{"GIN_MODE=release"}

	conOpts.Binds = []string{
		path + "config/n3iwfcfg.yaml:/free5gc/config/n3iwfcfg.yaml",
		path + "config/n3iwf-ipsec.sh:/free5gc/n3iwf-ipsec.sh",
	}
	conOpts.CapAdd = []string{"NET_ADMIN"}

	conOpts.Net = nfinfo.Network
	conOpts.NetworkAlias = "n3iwf" + strconv.Itoa(num) + ".free5gc.org"
	conOpts.NetworkIPAddress = "10.100.200.15"

	opts.BuildArgs = map[string]*string{"DEBUG_TOOLS": &debugTools}
	opts.Dockerfile = "Dockerfile"
	opts.Context = path + "nf_n3iwf"
	opts.Tags = []string{conOpts.Image}

	// depends on amf, smf, upf
	return opts, conOpts
}

func n3iwueConf(nfinfo common.NFInfo, path string, num int) (docker.ImageBuildOpts, docker.ContainerCreateOpts) {
	var opts docker.ImageBuildOpts
	var conOpts docker.ContainerCreateOpts

	debugTools := "false"

	conOpts.Name = "n3iwue"

	conOpts.Cmd = []string{"sleep infinity"}

	conOpts.Binds = []string{
		path + "config/n3uecfg.yaml:/n3iwue/config/n3ue.yaml",
	}
	conOpts.CapAdd = []string{"NET_ADMIN"}
	conOpts.Device = "/dev/net/tun"

	conOpts.Net = nfinfo.Network
	conOpts.NetworkAlias = "n3ue" + strconv.Itoa(num) + ".free5gc.org"

	opts.BuildArgs = map[string]*string{"DEBUG_TOOLS": &debugTools}
	opts.Dockerfile = "Dockerfile"
	opts.Context = path + "n3iwue"
	opts.Tags = []string{conOpts.Image}

	// depends on n3iwf
	return opts, conOpts
}

func ueransimConf(nfinfo common.NFInfo, path string, num int) (docker.ImageBuildOpts, docker.ContainerCreateOpts) {
	var opts docker.ImageBuildOpts
	var conOpts docker.ContainerCreateOpts

	debugTools := "false"

	conOpts.Name = "ueransim" + strconv.Itoa(num)
	conOpts.Image = "ueransim:latest"

	conOpts.Cmd = []string{"./nr-gnb", "-c", "/ueransim/config/gnbcfg.yaml"}

	conOpts.Binds = []string{
		path + "config/gnbcfg.yaml:/ueransim/config/gnbcfg.yaml",
		path + "config/uecfg.yaml:/ueransim/config/uecfg.yaml",
	}
	conOpts.CapAdd = []string{"NET_ADMIN"}
	// conOpts.Device = []container.DeviceMapping{{PathInContainer:
	// "/dev/net/tun"}}
	conOpts.Device = "/dev/net/tun"

	conOpts.Net = nfinfo.Network
	conOpts.NetworkAlias = "gnb" + strconv.Itoa(num) + ".free5gc.org"

	opts.BuildArgs = map[string]*string{"DEBUG_TOOLS": &debugTools}
	opts.Dockerfile = "Dockerfile"
	opts.Context = path + "ueransim"
	opts.Tags = []string{conOpts.Image}

	// depends on amf and upf
	return opts, conOpts
}
