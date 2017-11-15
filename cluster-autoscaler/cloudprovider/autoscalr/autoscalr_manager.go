/*
Copyright 2017 AutoScalr

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

package autoscalr

import (
	"io"
)

// AutoScalrManager is handles communication and data caching.
type AutoScalrManager struct {
	random   string
}

func createAutoScalrManagerInternal(configReader io.Reader) (*AutoScalrManager, error) {
	manager := &AutoScalrManager{
		random: "Test-jay",
	}
	return manager, nil
}

func CreateAutoScalrManager(configReader io.Reader) (*AutoScalrManager, error) {
	return createAutoScalrManagerInternal(configReader)
}
