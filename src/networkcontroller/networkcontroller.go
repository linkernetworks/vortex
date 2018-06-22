package networkcontroller

import (
	pb "github.com/linkernetworks/network-controller/messages"
	"github.com/linkernetworks/vortex/src/entity"
	k8sCtl "github.com/linkernetworks/vortex/src/kubernetes"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"k8s.io/client-go/kubernetes"
	"time"
)

type NetworkController struct {
	Clientset  kubernetes.Interface
	ClientConn *grpc.ClientConn
	Network    entity.Network
}

func New(clientset kubernetes.Interface, network entity.Network) (*NetworkController, error) {
	node, err := k8sCtl.GetNode(clientset, network.Node)
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

	return &NetworkController{
		Clientset:  clientset,
		ClientConn: conn,
		Network:    network,
	}, nil
}

func (nc *NetworkController) CreateNetwork() error {
	clientControl := pb.NewNetworkControlClient(nc.ClientConn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	for _, port := range nc.Network.PhysicalPorts {
		_, err := clientControl.AddPort(ctx, &pb.AddPortRequest{
			BridgeName: nc.Network.BridgeName,
			IfaceName:  port.Name})
		if err != nil {
			return err
		}
	}

	return nil
}
