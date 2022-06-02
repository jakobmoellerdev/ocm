// Copyright 2020 Copyright (c) 2020 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"github.com/open-component-model/ocm/pkg/common"
	"github.com/open-component-model/ocm/pkg/contexts/config"
	cfg "github.com/open-component-model/ocm/pkg/contexts/config/cpi"
	"github.com/open-component-model/ocm/pkg/contexts/oci/cpi"
	"github.com/open-component-model/ocm/pkg/runtime"
)

const (
	ConfigType   = "oci.config" + common.TypeGroupSuffix
	ConfigTypeV1 = ConfigType + runtime.VersionSeparator + "v1"
)

func init() {
	cfg.RegisterConfigType(ConfigType, cfg.NewConfigType(ConfigType, &ConfigSpec{}, usage))
	cfg.RegisterConfigType(ConfigTypeV1, cfg.NewConfigType(ConfigTypeV1, &ConfigSpec{}, usage))
}

// ConfigSpec describes a memory based config interface.
type ConfigSpec struct {
	runtime.ObjectVersionedType `json:",inline"`
	Aliases                     map[string]*cpi.GenericRepositorySpec `json:"aliases,omitempty"`
}

// NewConfigSpec creates a new memory ConfigSpec
func NewConfigSpec() *ConfigSpec {
	return &ConfigSpec{
		ObjectVersionedType: runtime.NewVersionedObjectType(ConfigType),
	}
}

func (a *ConfigSpec) GetType() string {
	return ConfigType
}

func (a *ConfigSpec) SetAlias(name string, spec cpi.RepositorySpec) error {
	g, err := cpi.ToGenericRepositorySpec(spec)
	if err != nil {
		return err
	}
	if a.Aliases == nil {
		a.Aliases = map[string]*cpi.GenericRepositorySpec{}
	}
	a.Aliases[name] = g
	return nil
}

func (a *ConfigSpec) ApplyTo(ctx config.Context, target interface{}) error {
	t, ok := target.(cpi.Context)
	if !ok {
		return config.ErrNoContext(ConfigType)
	}
	for n, s := range a.Aliases {
		t.SetAlias(n, s)
	}
	return nil
}

const usage = `
The config type <code>` + ConfigType + `</code> can be used to define
OCI registry aliases:

<pre>
    type: ` + ConfigType + `
    aliases:
       &lt;name>: &lt;OCI registry specification>
       ...
</pre>
`
