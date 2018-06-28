package networkcontroller

import (
	"time"

	pb "github.com/linkernetworks/network-controller/messages"
	"github.com/linkernetworks/vortex/src/entity"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type NetworkController struct {
	ClientCtl pb.NetworkControlClient
	Context   context.Context
}

func New(serverAddress string) (*NetworkController, error) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second)

	return &NetworkController{
		ClientCtl: pb.NewNetworkControlClient(conn),
		Context:   ctx,
	}, nil
}

func (nc *NetworkController) CreateOVSNetwork(bridgeName string, ports []entity.PhysicalPort) error {
	if _, err := nc.ClientCtl.CreateBridge(
		nc.Context,
		&pb.CreateBridgeRequest{
			BridgeName: bridgeName,
		}); err != nil {
		return err
	}

	for _, port := range ports {
		_, err := nc.ClientCtl.AddPort(
			nc.Context,
			&pb.AddPortRequest{
				BridgeName: bridgeName,
				IfaceName:  port.Name,
			})
		if err != nil {
			return err
		}
	}
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
