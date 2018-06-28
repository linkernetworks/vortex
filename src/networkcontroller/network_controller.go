package networkcontroller

import (
	"time"

	pb "github.com/linkernetworks/network-controller/messages"
	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/kubernetes"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type NetworkController struct {
	KubeCtl   *kubernetes.KubeCtl
	ClientCtl pb.NetworkControlClient
	Network   entity.Network
	Context   context.Context
}

func New(kubeCtl *kubernetes.KubeCtl, network entity.Network) (*NetworkController, error) {
	nodeIP, err := kubeCtl.GetNodeExternalIP(network.NodeName)
	if err != nil {
		return nil, err
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(nodeIP+":50051", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second)

	return &NetworkController{
		KubeCtl:   kubeCtl,
		ClientCtl: pb.NewNetworkControlClient(conn),
		Network:   network,
		Context:   ctx,
	}, nil
}

func (nc *NetworkController) CreateNetwork() error {
	/*
		_, err := nc.ClientCtl.CreateBridge(
			nc.Context,
			&pb.CreateBridgeRequest{
				BridgeName: nc.Network.BridgeName,
			})
		if err != nil {
			return err
		}

		for _, port := range nc.Network.PhysicalPorts {
			_, err := nc.ClientCtl.AddPort(
				nc.Context,
				&pb.AddPortRequest{
					BridgeName: nc.Network.BridgeName,
					IfaceName:  port.Name,
				})
			if err != nil {
				return err
			}
		}

	*/
	return nil
}

func (nc *NetworkController) DeleteNetwork() error {
	_, err := nc.ClientCtl.DeleteBridge(
		nc.Context,
		&pb.DeleteBridgeRequest{
			BridgeName: nc.Network.BridgeName,
		})
	if err != nil {
		return err
	}

	return nil
}
