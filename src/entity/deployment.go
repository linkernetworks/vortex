package entity

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	// DeploymentCollectionName is a const string
	DeploymentCollectionName string = "deployments"
	// DeploymentHostNetwork is the network type for the Deployment
	// host means the deployment use the hostNetwork (share the network with the host machine)
	DeploymentHostNetwork = "host"
	// DeploymentClusterNetwork is cluster which means use the cluster Network, maybe the flannel network
	DeploymentClusterNetwork = "cluster"
	// DeploymentCustomNetwork is custom which means the custom netwokr we created before, it support the OVS and DPDK network for additional network interface card
	DeploymentCustomNetwork = "custom"
)

// DeploymentRoute is the structure for add IP routing table
type DeploymentRoute struct {
	DstCIDR string `bson:"dstCIDR" json:"dstCIDR" validate:"required,cidrv4"`
	Gateway string `bson:"gateway" json:"gateway" validate:"omitempty,ipv4"`
}

// DeploymentNetwork is the structure for deployment network info
type DeploymentNetwork struct {
	Name   string `bson:"name" json:"name" validate:"required"`
	IfName string `bson:"ifName" json:"ifName" validate:"required"`
	// can not validate nil
	VlanTag   *int32            `bson:"vlanTag" json:"vlanTag" validate:"-"`
	IPAddress string            `bson:"ipAddress" json:"ipAddress" validate:"required,ipv4"`
	Netmask   string            `bson:"netmask" json:"netmask" validate:"required,ipv4"`
	Routes    []DeploymentRoute `bson:"routes,omitempty" json:"routes" validate:"required,dive,required"`

	// It's from the entity.Network entity
	BridgeName string `bson:"bridgeName" json:"bridgeName" validate:"-"`
}

// DeploymentVolume is the structure for deployment volume info
type DeploymentVolume struct {
	Name      string `bson:"name" json:"name" validate:"required"`
	MountPath string `bson:"mountPath" json:"mountPath" validate:"required"`
}

// Deployment is the structure for deployment info
type Deployment struct {
	ID           bson.ObjectId       `bson:"_id,omitempty" json:"id" validate:"-"`
	Name         string              `bson:"name" json:"name" validate:"required,k8sname"`
	Namespace    string              `bson:"namespace" json:"namespace" validate:"required"`
	Labels       map[string]string   `bson:"labels,omitempty" json:"labels" validate:"required,dive,keys,printascii,endkeys,required,printascii"`
	EnvVars      map[string]string   `bson:"envVars,omitempty" json:"envVars" validate:"required,dive,keys,printascii,endkeys,required,printascii"`
	Containers   []Container         `bson:"containers" json:"containers" validate:"required,dive,required"`
	Volumes      []DeploymentVolume  `bson:"volumes,omitempty" json:"volumes" validate:"required,dive,required"`
	Networks     []DeploymentNetwork `bson:"networks,omitempty" json:"networks" validate:"required,dive,required"`
	Capability   bool                `bson:"capability" json:"capability" validate:"-"`
	NetworkType  string              `bson:"networkType" json:"networkType" validate:"required,eq=host|eq=cluster|eq=custom"`
	NodeAffinity []string            `bson:"nodeAffinity" json:"nodeAffinity" validate:"required"`
	CreatedAt    *time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty" validate:"-"`

	Replicas int32 `bson:"replicas",json:"replicas" validate:"required"`
}

// GetCollection - get model mongo collection name.
func (m Deployment) GetCollection() string {
	return DeploymentCollectionName
}
