package networkcontroller

import (
	"time"

	pb "github.com/linkernetworks/network-controller/messages"
	"github.com/linkernetworks/vortex/src/entity"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// DEFAULT_CONTROLLER_PORT set the default port as 50051
const DEFAULT_CONTROLLER_PORT = "50051"

// NetworkController is the structure for Network Controller
type NetworkController struct {
	ClientCtl pb.NetworkControlClient
	Context   context.Context
}

// New will Set up a connection to the Network Controller server
func New(serverAddress string) (*NetworkController, error) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	ctx, _ := context.WithTimeout(context.Background(), 45*time.Second)

	return &NetworkController{
		ClientCtl: pb.NewNetworkControlClient(conn),
		Context:   ctx,
	}, nil
}

// CreateOVSNetwork will Create OVS Network by Network Controller
func (nc *NetworkController) CreateOVSNetwork(datapathType string, bridgeName string, phyIfaces []entity.PhyInterface, vlanTags []int32) error {
	if _, err := nc.ClientCtl.CreateBridge(
		nc.Context,
		&pb.CreateBridgeRequest{
			BridgeName:   bridgeName,
			DatapathType: datapathType,
		}); err != nil {
		return err
	}

	for _, phyIface := range phyIfaces {
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
	}
	return nil
}

// CreateOVSDPDKNetwork will Create OVS+DPDK Network by Network Controller
func (nc *NetworkController) CreateOVSDPDKNetwork(bridgeName string, phyIfaces []entity.PhyInterface, vlanTags []int32) error {
	if _, err := nc.ClientCtl.CreateBridge(
		nc.Context,
		&pb.CreateBridgeRequest{
			BridgeName:   bridgeName,
			DatapathType: "netdev",
		}); err != nil {
		return err
	}

	for _, phyIface := range phyIfaces {
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
	}
	return nil
}

// DeleteOVSNetwork will delete OVS network controller
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

// DumpOVSPorts will dump ports information of the target ovs
func (nc *NetworkController) DumpOVSPorts(bridgeName string) ([]*pb.PortInfo, error) {
	data, err := nc.ClientCtl.DumpPorts(
		nc.Context,
		&pb.DumpPortsRequest{
			BridgeName: bridgeName,
		})
	if err != nil {
		return nil, err
	}

	return data.Ports, nil
}
