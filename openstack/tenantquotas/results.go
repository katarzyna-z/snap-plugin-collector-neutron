// +build linux

/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2016 Intel Corporation

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

package tenantquotas

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
)

//Result represents the result of a get operation.
type Result struct {
	gophercloud.Result
}

// Extract will get the Volume object out of the commonResult object.
func (r Result) Extract() (map[string]map[string]int64, error) {
	var resp map[string]map[string]int64
	err := mapstructure.Decode(r.Body, &resp)
	return resp, err
}
