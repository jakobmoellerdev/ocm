// Copyright 2022 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package genericocireg

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/gardener/ocm/pkg/common/accessio"
	"github.com/gardener/ocm/pkg/common/accessobj"
	"github.com/gardener/ocm/pkg/errors"
	"github.com/gardener/ocm/pkg/oci"
	"github.com/gardener/ocm/pkg/oci/artdesc"
	"github.com/gardener/ocm/pkg/oci/repositories/ctf/format"
	"github.com/gardener/ocm/pkg/ocm/compdesc"
	"github.com/gardener/ocm/pkg/ocm/cpi"
	"github.com/gardener/ocm/pkg/ocm/repositories/ctf/comparch"
	"github.com/gardener/ocm/pkg/ocm/repositories/genericocireg/componentmapping"
	"github.com/opencontainers/go-digest"
	ociv1 "github.com/opencontainers/image-spec/specs-go/v1"
)

func NewState(mode accessobj.AccessMode, name, version string, access oci.ManifestAccess) (accessobj.State, error) {
	return accessobj.NewState(mode, NewStateAccess(access), NewStateHandler(name, version))
}

// StateAccess handles the component descriptor persistence in an OCI Manifest
type StateAccess struct {
	access oci.ManifestAccess
}

var _ accessobj.StateAccess = (*StateAccess)(nil)

func NewStateAccess(access oci.ManifestAccess) accessobj.StateAccess {
	return &StateAccess{
		access: access,
	}
}

func (s *StateAccess) Get() (accessio.BlobAccess, error) {
	mediaType := s.access.GetDescriptor().Config.MediaType
	switch mediaType {
	case componentmapping.ComponentDescriptorConfigMimeType, componentmapping.ComponentDescriptorLegacyConfigMimeType:
		return s.get()
	case "":
		return nil, errors.ErrNotFound(cpi.KIND_COMPONENTVERSION)
	default:
		return nil, errors.Newf("artefact is no component: %s", mediaType)
	}
}

func (s *StateAccess) get() (accessio.BlobAccess, error) {
	var config ComponentDescriptorConfig

	data, err := accessio.BlobData(s.access.GetConfigBlob())

	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	if config.ComponentDescriptorLayer == nil || config.ComponentDescriptorLayer.Digest == "" {
		return nil, errors.ErrInvalid("component descriptor config")
	}
	switch config.ComponentDescriptorLayer.MediaType {
	case componentmapping.ComponentDescriptorJSONMimeType, componentmapping.ComponentDescriptorYAMLMimeType:
		return s.access.GetBlob(config.ComponentDescriptorLayer.Digest)
	case componentmapping.ComponentDescriptorTarMimeType, componentmapping.LegacyComponentDescriptorTarMimeType:
		d, err := s.access.GetBlob(config.ComponentDescriptorLayer.Digest)
		if err != nil {
			return nil, err
		}
		r, err := d.Reader()
		if err != nil {
			return nil, err
		}
		defer r.Close()
		data, err := s.readComponentDescriptorFromTar(r)
		if err != nil {
			return nil, err
		}
		return accessio.BlobAccessForData(componentmapping.ComponentDescriptorYAMLMimeType, data), nil
	default:
		return nil, errors.ErrInvalid("config mediatype", config.ComponentDescriptorLayer.MediaType)
	}
}

// readComponentDescriptorFromTar reads the component descriptor from a tar.
// The component is expected to be inside the tar at "/component-descriptor.yaml"
func (s *StateAccess) readComponentDescriptorFromTar(r io.Reader) ([]byte, error) {
	tr := tar.NewReader(r)
	for {
		header, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				return nil, errors.New("no component descriptor found in tar")
			}
			return nil, fmt.Errorf("unable to read tar: %w", err)
		}

		if strings.TrimLeft(header.Name, "/") != comparch.ComponentDescriptorFileName {
			continue
		}

		var data bytes.Buffer
		if _, err := io.Copy(&data, tr); err != nil {
			return nil, fmt.Errorf("erro while reading component descriptor file from tar: %w", err)
		}
		return data.Bytes(), err
	}
}

func (s StateAccess) Digest() digest.Digest {
	blob, err := s.access.GetConfigBlob()
	if err != nil {
		return ""
	}
	return blob.Digest()
}

func (s *StateAccess) Put(data []byte) error {
	desc := s.access.GetDescriptor()
	mediaType := desc.Config.MediaType
	if mediaType == "" {
		mediaType = componentmapping.ComponentDescriptorConfigMimeType
		desc.Config.MediaType = mediaType
	}

	arch, err := s.writeComponentDescriptorFromTar(data)
	if err != nil {
		return err
	}
	config := ComponentDescriptorConfig{
		ComponentDescriptorLayer: artdesc.DefaultBlobDescriptor(arch),
	}

	configdata, err := json.Marshal(&config)
	if err != nil {
		return err
	}

	err = s.access.AddBlob(arch)
	if err != nil {
		return err
	}
	configblob := accessio.BlobAccessForData(componentmapping.ComponentDescriptorConfigMimeType, configdata)
	err = s.access.AddBlob(configblob)
	if err != nil {
		return err
	}
	desc.Config = *artdesc.DefaultBlobDescriptor(configblob)
	if len(desc.Layers) < 2 {
		desc.Layers = []ociv1.Descriptor{*artdesc.DefaultBlobDescriptor(arch)}
	} else {
		desc.Layers[0] = *artdesc.DefaultBlobDescriptor(arch)
	}
	return nil
}

// readComponentDescriptorFromTar reads the component descriptor from a tar.
// The component is expected to be inside the tar at "/component-descriptor.yaml"
func (s *StateAccess) writeComponentDescriptorFromTar(data []byte) (cpi.BlobAccess, error) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	err := tw.WriteHeader(&tar.Header{
		Typeflag: tar.TypeReg,
		Name:     comparch.ComponentDescriptorFileName,
		Size:     int64(len(data)),
		ModTime:  format.ModTime,
	})
	if err != nil {
		return nil, errors.Newf("unable to add component descriptor header: %s", err)
	}
	if _, err := io.Copy(tw, bytes.NewBuffer(data)); err != nil {
		return nil, errors.Newf("unable to write component-descriptor to tar: %s", err)
	}
	if err := tw.Close(); err != nil {
		return nil, errors.Newf("unable to close tar writer: %s", err)
	}
	return accessio.BlobAccessForData(componentmapping.ComponentDescriptorTarMimeType, buf.Bytes()), nil
}

// ComponentDescriptorConfig is a Component-Descriptor OCI configuration that is used to store the reference to the
// (pseudo-)layer used to store the Component-Descriptor in.
type ComponentDescriptorConfig struct {
	ComponentDescriptorLayer *ociv1.Descriptor `json:"componentDescriptorLayer,omitempty"`
}

////////////////////////////////////////////////////////////////////////////////

// StateHandler handles the encoding of a component descriptor
type StateHandler struct {
	name    string
	version string
}

var _ accessobj.StateHandler = (*StateHandler)(nil)

func NewStateHandler(name, version string) accessobj.StateHandler {
	return &StateHandler{
		name:    name,
		version: version,
	}
}

func (i StateHandler) Initial() interface{} {
	return compdesc.New(i.name, i.version)
}

// Encode always provides a yaml representation
func (i StateHandler) Encode(d interface{}) ([]byte, error) {
	desc := d.(*compdesc.ComponentDescriptor)
	desc.Name = i.name
	desc.Version = i.version
	return compdesc.Encode(desc)
}

// Decode always accepts a yaml representation, and therefore json, also
func (i StateHandler) Decode(data []byte) (interface{}, error) {
	return compdesc.Decode(data)
}

func (i StateHandler) Equivalent(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}
