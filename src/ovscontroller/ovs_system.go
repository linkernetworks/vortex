package ovscontroller

import (
	"net"

	"github.com/linkernetworks/network-controller/utils"
	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/networkcontroller"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	"gopkg.in/mgo.v2/bson"
)

func DumpPorts(sp *serviceprovider.Container, nodeName string, bridgeName string) ([]entity.OVSPortInfo, error) {
	nodeIP, err := sp.KubeCtl.GetNodeInternalIP(nodeName)
	if err != nil {
		return nil, err
	}

	nodeAddr := net.JoinHostPort(nodeIP, networkcontroller.DEFAULT_CONTROLLER_PORT)
	nc, err := networkcontroller.New(nodeAddr)

	//Get the information of OVS ports
	retPorts, err := nc.DumpOVSPorts(bridgeName)
	if err != nil {
		return nil, err
	}

	//We need to find a mapping for veth to podName
	//1. lookup the mongodb to find all deployments which bridge name is equal to bridgeName
	//2. lookup all current pods which is belogs to above deployments
	//3. use the pod's UID and the interfae of each entity.DeploymentNetowrk to geneate the vetxXXXXXX
	//4. use the vtxXXXXXXXXX as the key to combine the podName/interfaace and the OVSPorts
	session := sp.Mongo.NewSession()
	defer session.Close()

	//1. lookup the mongodb to find all deployments which bridge name is equal to bridgeName
	//In order to the following use, use the map here and the key is the deploynent name and the value is the deployment object.
	//We need to mapping the deployment.Networks with Pod's UID and the connection is the label of the Pod is vortex=deployment.name.
	deployments := []entity.Deployment{}
	session.FindAll(entity.DeploymentCollectionName, bson.M{"networks.bridgeName": bridgeName}, &deployments)
	deployMap := map[string]*entity.Deployment{}

	for _, v := range deployments {
		deployMap[v.Name] = &v
	}

	//2. lookup all current pods which is belogs to above deployments
	pods, err := sp.KubeCtl.GetPods("")
	if err != nil {
		return nil, err
	}

	//Create a local structure to combine the deploymentNetwork and PodName.
	//We want to know the interface name in that Pod, so we need to keep the deploymentNetwork object.
	type portData struct {
		podName string
		entity.DeploymentNetwork
	}
	interfaces := map[string]*portData{}
	for _, v := range pods {
		name, ok := v.Labels["vortex"]
		if !ok {
			continue
		}
		deploy, ok := deployMap[name]
		if !ok {
			continue
		}

		//3. use the pod's UID and the interfae of each entity.DeploymentNetowrk to geneate the vtxXXXXXXXXX
		//for each deploymentNetwork, we get the vethname via veth+sha256(podUID + interfaceName in container)[0:8]
		uid := v.ObjectMeta.UID
		for _, k := range deploy.Networks {
			vethName := utils.GenerateVethName(string(uid), k.IfName)
			//Use the veth name as the key and the PodName/InterfaceName in the value, we will add those inforamtion
			//to OVSPortInfo later.
			interfaces[vethName] = &portData{
				v.Name,
				k,
			}
		}
	}

	//4. use the vethxxx as the key to combine the podName/interfaace and the OVSPorts
	ports := []entity.OVSPortInfo{}
	for _, v := range retPorts {
		port := entity.OVSPortInfo{
			PortID:     v.ID,
			Name:       v.Name,
			MacAddress: v.MacAddr,
			Received: entity.OVSPortStats{
				Packets: v.Received.Packets,
				Bytes:   v.Received.Byte,
				Dropped: v.Received.Dropped,
				Errors:  v.Received.Errors,
			},
			Transmitted: entity.OVSPortStats{
				Packets: v.Transmitted.Packets,
				Bytes:   v.Transmitted.Byte,
				Dropped: v.Transmitted.Dropped,
				Errors:  v.Transmitted.Errors,
			},
		}

		if i, ok := interfaces[port.Name]; ok {
			port.InterfaceName = i.IfName
			port.PodName = i.podName
		}
		ports = append(ports, port)

	}

	return ports, nil
}
