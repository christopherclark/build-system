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
	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
)

//go:generate mockery -name System -case=underscore

type System interface {
	CachePath() (string, error)
	Detect(context libcnb.DetectContext, result *libcnb.DetectResult) error
	DefaultArguments() []string
	DefaultTarget() string
	Distribution(layersPath string) string
	DistributionLayer(resolver libpak.DependencyResolver, cache libpak.DependencyCache, plan *libcnb.BuildpackPlan) (libcnb.LayerContributor, error)
	Participate(resolver libpak.PlanEntryResolver) (bool, error)
	Wrapper() string
}
