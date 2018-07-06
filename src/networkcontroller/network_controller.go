package networkcontroller

import (
	"time"

	pb "github.com/linkernetworks/network-controller/messages"
	"github.com/linkernetworks/vortex/src/entity"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const DEFAULT_CONTROLLER_PORT = "50051"

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
			BridgeName:   bridgeName,
			DatapathType: "system",
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

		if len(port.VlanTags) > 0 {
			_, err := nc.ClientCtl.SetPort(
				nc.Context,
				&pb.SetPortRequest{
					IfaceName: port.Name,
					Options: &pb.PortOptions{
						VLANMode: "trunk",
						Trunk:    port.VlanTags,
					},
				})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (nc *NetworkController) DeleteOVSNetwork(bridgeName string) error {
	_, err := nc.ClientCtl.DeleteBridge(
		nc.Context,
		&pb.DeleteBridgeRequest{
			BridgeName: bridgeName,
		})
	if err != nil {
		return err
	}

	return nil
}
