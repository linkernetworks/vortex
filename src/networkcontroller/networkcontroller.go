package networkcontroller

import (
	pb "github.com/linkernetworks/network-controller/messages"
	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/kubernetes"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"time"
)

type NetworkController struct {
	KubeCtl   *kubernetes.KubeCtl
	ClientCtl pb.NetworkControlClient
	Network   entity.Network
	Context   context.Context
}

func New(kubeCtl *kubernetes.KubeCtl, network entity.Network) (*NetworkController, error) {
	node, err := kubeCtl.GetNode(network.NodeName)
	if err != nil {
		return nil, err
	}

	var nodeIP string
	for _, addr := range node.Status.Addresses {
		if addr.Type == "ExternalIP" {
			nodeIP = addr.Address
			break
		}
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(nodeIP, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	return &NetworkController{
		KubeCtl:   kubeCtl,
		ClientCtl: pb.NewNetworkControlClient(conn),
		Network:   network,
		Context:   ctx,
	}, nil
}

func (nc *NetworkController) CreateNetwork() error {
	for _, port := range nc.Network.PhysicalPorts {
		_, err := nc.ClientCtl.AddPort(nc.Context, &pb.AddPortRequest{
			BridgeName: nc.Network.BridgeName,
			IfaceName:  port.Name})
		if err != nil {
			return err
		}
	}

	return nil
}
