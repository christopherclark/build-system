/*
 * Copyright 2018-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package system

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
)

type Build struct {
	Logger  bard.Logger
	Systems []System
}

func NewBuild() Build {
	return Build{
		Logger:  bard.NewLogger(os.Stdout),
		Systems: []System{Gradle{}, Maven{}},
	}
}

func (b Build) Build(context libcnb.BuildContext) (libcnb.BuildResult, error) {
	pr := libpak.PlanEntryResolver{Plan: context.Plan}

	dr, err := libpak.NewDependencyResolver(context)
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to create dependency resolver: %w", err)
	}

	dc := libpak.NewDependencyCache(context.Buildpack)

	b.Logger.Title(context.Buildpack)
	result := libcnb.BuildResult{}

	for _, s := range b.Systems {
		if ok, err := s.Participate(pr); err != nil {
			return libcnb.BuildResult{}, fmt.Errorf("unable to determine participation: %w", err)
		} else if !ok {
			continue
		}

		var command string
		wrapper := filepath.Join(context.Application.Path, s.Wrapper())
		if _, err := os.Stat(wrapper); os.IsNotExist(err) {
			command = s.Distribution(context.Layers.Path)

			layer, err := s.DistributionLayer(dr, dc, &result.Plan)
			if err != nil {
				return libcnb.BuildResult{}, fmt.Errorf("unable to create distribution layer: %w", err)
			}
			result.Layers = append(result.Layers, layer)
		} else if err != nil {
			return libcnb.BuildResult{}, fmt.Errorf("unable to stat %s: %w", wrapper, err)
		} else {
			command = wrapper
		}

		cache, err := s.CachePath()
		if err != nil {
			return libcnb.BuildResult{}, fmt.Errorf("unable to determine cache location: %w", err)
		}
		result.Layers = append(result.Layers, NewCache(cache))

		run, err := NewApplication(context.Application.Path, command, s.DefaultArguments(), s.DefaultTarget())
		if err != nil {
			return libcnb.BuildResult{}, fmt.Errorf("unable to create run layer: %w", err)
		}

		result.Layers = append(result.Layers, run)
	}

	return result, nil
}
