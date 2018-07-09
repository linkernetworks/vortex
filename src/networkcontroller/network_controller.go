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

func (nc *NetworkController) CreateOVSNetwork(datapathType string, bridgeName string, phyIface entity.PhyInterface, vlanTags []int32) error {
	if _, err := nc.ClientCtl.CreateBridge(
		nc.Context,
		&pb.CreateBridgeRequest{
			BridgeName:   bridgeName,
			DatapathType: datapathType,
		}); err != nil {
		return err
	}

	_, err := nc.ClientCtl.AddPort(
		nc.Context,
		&pb.AddPortRequest{
			BridgeName: bridgeName,
			IfaceName:  phyIface.Name,
		})
	if err != nil {
		return err
	}

	if len(vlanTags) > 0 {
		_, err := nc.ClientCtl.SetPort(
			nc.Context,
			&pb.SetPortRequest{
				IfaceName: phyIface.Name,
				Options: &pb.PortOptions{
					VLANMode: "trunk",
					Trunk:    vlanTags,
				},
			})
		if err != nil {
			return err
		}
	}
	return nil
}

func (nc *NetworkController) CreateOVSDPDKNetwork(bridgeName string, phyIface entity.PhyInterface, vlanTags []int32) error {
	if _, err := nc.ClientCtl.CreateBridge(
		nc.Context,
		&pb.CreateBridgeRequest{
			BridgeName:   bridgeName,
			DatapathType: "netdev",
		}); err != nil {
		return err
	}

	if phyIface.IsDPDKPort == true {
		_, err := nc.ClientCtl.AddDPDKPort(
			nc.Context,
			&pb.AddPortRequest{
				BridgeName:  bridgeName,
				IfaceName:   phyIface.Name,
				DpdkDevargs: phyIface.PCIID,
			})
		if err != nil {
			return err
		}
	} else {
		_, err := nc.ClientCtl.AddPort(
			nc.Context,
			&pb.AddPortRequest{
				BridgeName: bridgeName,
				IfaceName:  phyIface.Name,
			})
		if err != nil {
			return err
		}
	}

	if len(vlanTags) > 0 {
		_, err := nc.ClientCtl.SetPort(
			nc.Context,
			&pb.SetPortRequest{
				IfaceName: phyIface.Name,
				Options: &pb.PortOptions{
					VLANMode: "trunk",
					Trunk:    vlanTags,
				},
			})
		if err != nil {
			return err
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
