/*
Copyright 2017 The Nuclio Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package runtime

import (
	"github.com/nuclio/nuclio/pkg/cmdrunner"
	"github.com/nuclio/nuclio/pkg/dockerclient"
	"github.com/nuclio/nuclio/pkg/errors"

	"github.com/nuclio/nuclio-sdk"
)

type Runtime interface {

	// GetDefaultProcessorBaseImageName returns the image name of the default processor base image
	GetDefaultProcessorBaseImageName() string

	// DetectFunctionHandlers returns a list of all the handlers
	// in that directory given a path holding a function (or functions)
	DetectFunctionHandlers(functionPath string) ([]string, error)

	// OnAfterStagingDirCreated prepares anything it may need in that directory
	// towards building a functioning processor,
	OnAfterStagingDirCreated(stagingDir string) error

	// GetProcessorImageObjectPaths returns a map of objects the runtime needs to copy into the processor image
	// the key can be a dir, a file or a url of a file
	// the value is an absolute path into the docker image
	GetProcessorImageObjectPaths() map[string]string

	// GetExtension returns the source extension of the runtime (e.g. .go)
	GetExtension() string

	// GetName returns the name of the runtime, including version if applicable
	GetName() string
}

type Configuration interface {
	GetFunctionPath() string

	GetFunctionDir() string

	GetFunctionName() string

	GetFunctionHandler() string

	GetNuclioSourceDir() string

	GetNuclioSourceURL() string

	GetStagingDir() string

	GetNoBaseImagePull() bool
}

type Factory interface {
	Create(logger nuclio.Logger,
		configuration Configuration) (Runtime, error)
}

type AbstractRuntime struct {
	Logger        nuclio.Logger
	Configuration Configuration
	DockerClient  dockerclient.Client
	CmdRunner     cmdrunner.CmdRunner
}

func NewAbstractRuntime(logger nuclio.Logger, configuration Configuration) (*AbstractRuntime, error) {
	var err error

	newRuntime := &AbstractRuntime{
		Logger:        logger,
		Configuration: configuration,
	}

	// create a docker client
	newRuntime.DockerClient, err = dockerclient.NewShellClient(newRuntime.Logger)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create docker client")
	}

	// set cmdrunner
	newRuntime.CmdRunner, err = cmdrunner.NewShellRunner(newRuntime.Logger)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create command runner")
	}

	return newRuntime, nil
}

func (ar *AbstractRuntime) OnAfterStagingDirCreated(stagingDir string) error {
	return nil
}

func (ar *AbstractRuntime) GetProcessorConfigFileContents() string {
	return ""
}

// return a map of objects the runtime needs to copy into the processor image
// the key can be a dir, a file or a url of a file
// the value is an absolute path into the docker image
func (ar *AbstractRuntime) GetProcessorImageObjectPaths() map[string]string {
	return nil
}
