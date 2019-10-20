# cni-plugin-logger-cloudwatch
A CNI plugin publishing its requests to CloudWatch Logs

## CNI

CNI(Container Network Interface) serves as an interface between container runtime and network implementaion. 

## Demo 

1. Create CNI Configuration

    ```
    {
    "cniVersion": "0.3.1",
    "name": "mynet",
    "plugins": [
        {
            "type": "bridge",
            "isGateway": true,
            "isDefaultGateway": true,
            "ipMasq": true,
            "bridge": "br0",
            "ipam": {
                "type": "host-local",
                "subnet": "10.11.0.0/24",
                "gateway": "10.11.0.1",
                "routes": [
                    { "dst": "0.0.0.0/0" }
                ],
            "dataDir": "/run/ipam-out-net"
            },
            "dns": {
            "nameservers": [ "8.8.8.8" ]
            }
        },
        {
        "type": "portmap",
        "capabilities": {"portMappings": true},
        "snat": false
        },
        {
        "type":"cni-plugin-logger-cloudwatch",
        "debug":true, 
        "debugDir": "/var/vcap/data/cni-configs/net-debug",
        "logGroupName": "cni-plugin-logger-cloudwatch"
        }
    ]
    }
    ```

1. Create Network Namespace
    `sudo ip netns add cnitest1`


1. Use the [cnitool](https://github.com/containernetworking/cni/tree/master/cnitool) to test the 

    `sudo CNI_PATH=/home/ec2-user/cni/plugins NETCONFPATH=/home/ec2-user/cni/net.d /cnitool add mynet /var/run/netns/cnitest1`

