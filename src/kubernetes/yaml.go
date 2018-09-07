package kubernetes

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
)

func ParseK8SYAML(data []byte) (runtime.Object, error) {
	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, _, err := decode(data, nil, nil)
	if err != nil {
		return obj, fmt.Errorf("Error while decoding YAML object. Err was: %v", err)
	}
	return obj, nil
}
