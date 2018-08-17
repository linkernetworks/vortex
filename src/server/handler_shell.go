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

package server

import (
	"net/http"
	"os"
	"path/filepath"

	response "github.com/linkernetworks/vortex/src/net/http"
	"github.com/linkernetworks/vortex/src/web"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

// TerminalResponse is sent by handleExecShell. The Id is a random session id that binds the original REST request and the SockJS connection.
// Any clientapi in possession of this Id can hijack the terminal session.
type TerminalResponse struct {
	Id string `json:"id"`
}

// Handles execute shell API call
// func handleExecShell(request *restful.Request, response *restful.Response) {
func handleExecShell(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response
	sessionId, err := genTerminalSessionId()
	if err != nil {
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}

	kubeconfigPath := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)

	terminalSessions[sessionId] = TerminalSession{
		id:       sessionId,
		bound:    make(chan error),
		sizeChan: make(chan remotecommand.TerminalSize),
	}
	go WaitForTerminal(sp.KubeCtl.Clientset, cfg, req, sessionId)
	resp.WriteHeaderAndEntity(http.StatusOK, TerminalResponse{Id: sessionId})
}
