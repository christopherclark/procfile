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

package procfile

import (
	"fmt"
	"os"
	"sort"

	"github.com/buildpacks/libcnb"
	"github.com/paketoio/libpak"
	"github.com/paketoio/libpak/bard"
)

type Build struct {
	Logger bard.Logger
}

func NewBuild() Build {
	return Build{Logger: bard.NewLogger(os.Stdout)}
}

func (b Build) Build(context libcnb.BuildContext) (libcnb.BuildResult, error) {
	r := libpak.PlanEntryResolver{Plan: context.Plan}

	e, ok, err := r.Resolve("procfile")
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to resolve buildpack plan entry procfile: %w", err)
	} else if !ok {
		return libcnb.BuildResult{}, nil
	}

	b.Logger.Title(context.Buildpack)
	result := libcnb.BuildResult{}

	for k, v := range e.Metadata {
		result.Processes = append(result.Processes, libcnb.Process{Type: k, Command: v.(string)})
	}

	sort.Slice(result.Processes, func(i int, j int) bool {
		return result.Processes[i].Type < result.Processes[j].Type
	})

	return result, nil
}