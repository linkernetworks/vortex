package entity

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	// PodCollectionName's const
	PodCollectionName string = "pods"
	// The network type for the Pod
	// host means the pod use the hostNetwork (share the network with the host machine)
	PodHostNetwork = "host"
	// cluster means use the cluster Network, maybe the flannel network
	PodClusterNetwork = "cluster"
	// custom means the custom netwokr we created before, it support the OVS and DPDK network for additional network interface card
	PodCustomNetwork = "custom"
)

// Container is the structure for init Container info
type Container struct {
	Name    string   `bson:"name" json:"name" validate:"required,k8sname"`
	Image   string   `bson:"image" json:"image" validate:"required"`
	Command []string `bson:"command" json:"command" validate:"required,dive,required"`
}

// PodRoute is the structure for add IP routing table
type PodRoute struct {
	DstCIDR string `bson:"dstCIDR" json:"dstCIDR" validate:"required,cidr"`
	Gateway string `bson:"gateway" json:"gateway" validate:"required,ip"`
}

// PodNetwork is the structure for pod network info
type PodNetwork struct {
	Name       string     `bson:"name" json:"name"`
	IfName     string     `bson:"ifName" json:"ifName"`
	VlanTag    *int32     `bson:"vlanTag" json:"vlanTag"`
	IPAddress  string     `bson:"ipAddress" json:"ipAddress"`
	Netmask    string     `bson:"netmask" json:"netmask"`
	Routes     []PodRoute `bson:"routes,omitempty" json:"routes"`
	BridgeName string     `bson:"bridgeName" json:"bridgeName"` //its from the entity.Network entity
}

// PodVolume is the structure for pod volume info
type PodVolume struct {
	Name      string `bson:"name" json:"name" validate:"required"`
	MountPath string `bson:"mountPath" json:"mountPath" validate:"required"`
}

// Pod is the structure for pod info
type Pod struct {
	ID            bson.ObjectId     `bson:"_id,omitempty" json:"id" validate:"-"`
	Name          string            `bson:"name" json:"name" validate:"required,k8sname"`
	Namespace     string            `bson:"namespace" json:"namespace" validate:"required"`
	Labels        map[string]string `bson:"labels,omitempty" json:"labels" validate:"required,dive,keys,alphanum,endkeys,required,alphanum"`
	Containers    []Container       `bson:"containers" json:"containers" validate:"required,dive,required"`
	CreatedAt     *time.Time        `bson:"createdAt,omitempty" json:"createdAt,omitempty" validate:"-"`
	Volumes       []PodVolume       `bson:"volumes,omitempty" json:"volumes" validate:"required,dive,required"`
	Networks      []PodNetwork      `bson:"networks,omitempty" json:"networks" validate:"required,dive,required"`
	RestartPolicy string            `bson:"restartPolicy" json:"restartPolicy" validate:"required,eq=Always|eq=OnFailure|eq=Never`
	Capability    bool              `bson:"capability" json:"capability" validate:"-"`
	HostNetwork   bool              `bson:"hostNetwork" json:"hostNetwork" validate:"-"`
	NetworkType   string            `bson:"networkType" json:"networkType" validate:"required,eq=Host|Cluster|Custom`
}

// GetCollection - get model mongo collection name.
func (m Pod) GetCollection() string {
	return PodCollectionName
}
