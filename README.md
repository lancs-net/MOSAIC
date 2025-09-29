# MOSAIC: Piecing Together 5G and LEOs for NTN Integration Experimentation

This repository contains the source code for the paper [MOSAIC: Piecing Together 5G and LEOs for NTN Integration Experimentation](https://dl.acm.org/doi/10.1145/3748749.3749091) published in 3rd ACM Workshop on LEO Networking and Communication (LEO-NET'25) co-located with ACM SIGCOMM 2025.

This repository integrates with the following open-source projects:
- [Free5GC Compose] (https://github.com/free5gc/free5gc-compose)
- [UERANSIM] (https://github.com/aligungr/UERANSIM)

It also supports the following proprietary software:
- [Fraunhofer FOKUS Open5GCore] (https://www.open5gcore.org)

Both Free5GC and Open5GCore can be used as the 5G Core with this repository.
Support for OpenAirInterface (OAI) is currently under development.

## Prerequisites
- Docker and Docker Compose
- GoLang (unless you se the docker deployment for the API)
- GTP5G Kernel Module (for Free5GC)

## Deployment
Topology, Scenario and Networks can be defined in a json file stored in ```pkg/files/*.json``` (Example files present). In the network configuration file, you need to provide the name and a subnet for the networks you want to create.
```
{
    "satellite": "10.0.0.0/24",
    "ground": "10.1.0.0/24"
}
```
In the scenario configuration file, you can define the nodes you want to create, their type (e.g., gNB, UE, 5GC, etc.), and the networks they should be connected to. This file will contain __all__ the network functions you want to deploy.
```
{
    "Path": "/home/tudor/",
    "NetworkFunctions": [
        {
            "name": "amf",
            "number": 1,
            "network": "ground"
        },
        {
            "name": "upf",
            "number": 1,
            "network": "satellite"
        }
        ...
    ]
}
```

MOSAIC provides an API to manage the lifecycle of the experiment from the deployment to the testing. Now that you have the configuration files ready, you can start the API server:
```
cd cmd/mosaic
go build .
go run .
```

The API server will start on port 8000. The following endpoints are available:
+ ```GET /```: Returns a welcome message.
+ ```POST /net```: Configures and deploys the networks defined in the network configuration file. Needs a JSON input (Example provided in /pkg/files/networks.json).
+ ```DELETE /net```: Removes the networks defined in the network configuration file.
+ ```POST /f5gc/base/:path```: Builds the base docker images for the Free5GC network functions. Needs a path parameter to the Free5GC source code. (Example: /home/user/free5gc)
+ ```DELETE /f5gc/base```: Removes the base docker images for the Free5GC network functions.
+ ```POST /f5gc```: Configures and deploys the Free5GC network functions defined in the scenario configuration file. Needs a JSON input (Example provided in /pkg/files/nfscenario.json).
+ ```DELETE /f5gc/:nf```: Removes the all Free5GC network functions or a specifed network function. Needs a parameter to the network function name (Example: all, amf, upf, etc.).
+ ```PUT /f5gc/:step/:nf```: Starts/ Stops a specifed network function. Needs a parameter to the step (start/stop) and the network function name (Example: start, amf).
+ ```GET /status```: Returns the status of all the deployed network functions.


Note:
Dangling images may be created during the build process. It is advised to remove them from time to time to free up disk space.
``` docker rmi $(docker images -f "dangling=true" -q) ```
