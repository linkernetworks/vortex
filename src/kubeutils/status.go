package kubeutils

import (
	"fmt"
	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"

	"gopkg.in/mgo.v2/bson"
)

func GetNonCompletedPods(sp *serviceprovider.Container, query bson.M) ([]entity.Pod, error) {
	session := sp.Mongo.NewSession()
	defer session.Close()
	pods := []entity.Pod{}

	err := session.FindAll(entity.PodCollectionName, query, &pods)
	if err != nil {
		return nil, fmt.Errorf("load the database %v fail:%v", query, err)
	}

	ret := []entity.Pod{}
	for _, pod := range pods {
		//Check the pod's status, report error if at least one pod is running.
		//We use the "" to fetch all pods for all namespaces
		currentPod, err := sp.KubeCtl.GetPod(pod.Name, "")
		if err != nil {
			continue
		}

		if !sp.KubeCtl.IsPodCompleted(currentPod) {
			ret = append(ret, pod)
		}
	}

	return ret, nil
}
