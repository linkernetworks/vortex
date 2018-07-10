package pod

import (
	"github.com/linkernetworks/vortex/src/entity"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreatePod(sp *serviceprovider.Container, pod *entity.Pod) error {
	var containers []corev1.Container
	for _, container := range pod.Containers {
		c := corev1.Container{}
		c.Name = container.Name
		c.Image = container.Image
		c.Command = container.Command
		containers = append(containers, c)
	}
	p := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: pod.Name,
		},
		Spec: corev1.PodSpec{
			Containers: containers,
		},
	}
	_, err := sp.KubeCtl.CreatePod(&p)
	return err
}

func DeletePod(sp *serviceprovider.Container, podName string) error {
	return sp.KubeCtl.DeletePod(podName)
}
