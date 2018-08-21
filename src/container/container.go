// Copyright 2017 The Kubernetes Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package container

import (
	"io"
	"io/ioutil"

	restful "github.com/emicklei/go-restful"
	"github.com/linkernetworks/vortex/src/serviceprovider"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
)

// maximum number of lines loaded from the apiserver
var lineReadLimit int64 = 5000

// maximum number of bytes loaded from the apiserver
var byteReadLimit int64 = 500000

// GetLogDetails returns logs for particular pod and container. When container is null, logs for the first one
// are returned. Previous indicates to read archived logs created by log rotation or container crash
func GetLogDetails(sp *serviceprovider.Container, namespace, podID string, container string,
	logSelector *Selection, usePreviousLogs bool) (*LogDetails, error) {

	logOptions := mapToLogOptions(container, logSelector, usePreviousLogs)
	rawLogs, err := readRawLogs(sp, namespace, podID, logOptions)
	if err != nil {
		return nil, err
	}
	details := ConstructLogDetails(podID, rawLogs, container, logSelector)
	return details, nil
}

// Maps the log selection to the corresponding api object
// Read limits are set to avoid out of memory issues
func mapToLogOptions(container string, logSelector *Selection, previous bool) *v1.PodLogOptions {
	logOptions := &v1.PodLogOptions{
		Container:  container,
		Follow:     false,
		Previous:   previous,
		Timestamps: true,
	}

	if logSelector.LogFilePosition == Beginning {
		logOptions.LimitBytes = &byteReadLimit
	} else {
		logOptions.TailLines = &lineReadLimit
	}

	return logOptions
}

// Construct a request for getting the logs for a pod and retrieves the logs.
func readRawLogs(sp *serviceprovider.Container, namespace, podID string, logOptions *v1.PodLogOptions) (
	string, error) {
	readCloser, err := openStream(sp, namespace, podID, logOptions)
	if err != nil {
		return err.Error(), nil
	}

	defer readCloser.Close()

	result, err := ioutil.ReadAll(readCloser)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

// GetLogFile returns a stream to the log file which can be piped directly to the response. This avoids out of memory
// issues. Previous indicates to read archived logs created by log rotation or container crash
func GetLogFile(sp *serviceprovider.Container, namespace, podID string, container string, usePreviousLogs bool) (io.ReadCloser, error) {
	logOptions := &v1.PodLogOptions{
		Container:  container,
		Follow:     false,
		Previous:   usePreviousLogs,
		Timestamps: false,
	}
	logStream, err := openStream(sp, namespace, podID, logOptions)
	return logStream, err
}

func openStream(sp *serviceprovider.Container, namespace, podID string, logOptions *v1.PodLogOptions) (io.ReadCloser, error) {
	return sp.KubeCtl.Clientset.CoreV1().RESTClient().Get().
		Namespace(namespace).
		Name(podID).
		Resource("pods").
		SubResource("log").
		VersionedParams(logOptions, scheme.ParameterCodec).Stream()
}

// ConstructLogDetails creates a new log details structure for given parameters.
func ConstructLogDetails(podID string, rawLogs string, container string, logSelector *Selection) *LogDetails {
	parsedLines := ToLogLines(rawLogs)
	logLines, fromDate, toDate, logSelection, lastPage := parsedLines.SelectLogs(logSelector)

	readLimitReached := isReadLimitReached(int64(len(rawLogs)), int64(len(parsedLines)), logSelector.LogFilePosition)
	truncated := readLimitReached && lastPage

	info := LogInfo{
		PodName:       podID,
		ContainerName: container,
		FromDate:      fromDate,
		ToDate:        toDate,
		Truncated:     truncated,
	}
	return &LogDetails{
		Info:      info,
		Selection: logSelection,
		LogLines:  logLines,
	}
}

// Checks if the amount of log file returned from the apiserver is equal to the read limits
func isReadLimitReached(bytesLoaded int64, linesLoaded int64, logFilePosition string) bool {
	return (logFilePosition == Beginning && bytesLoaded >= byteReadLimit) ||
		(logFilePosition == End && linesLoaded >= lineReadLimit)
}

func HandleDownload(resp *restful.Response, result io.ReadCloser) error {
	resp.AddHeader(restful.HEADER_ContentType, "text/plain")
	defer result.Close()
	_, err := io.Copy(resp, result)
	if err != nil {
		return err
	}
	return nil
}
