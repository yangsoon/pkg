/*
Copyright 2023 The KubeVela Authors.

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

package util

import (
	"sort"
	"strconv"
	"strings"

	"cuelang.org/go/cue"

	"github.com/kubevela/pkg/util/slices"
)

const orderKey = "step"

// Iterate over all fields of the cue.Value with fn, if fn returns true,
// iteration stops
func Iterate(value cue.Value, fn func(v cue.Value) (stop bool)) (stop bool) {
	var it *cue.Iterator
	// skip definition
	if strings.Contains(value.Path().String(), "#") {
		return false
	}
	switch value.Kind() {
	case cue.ListKind:
		_it, _ := value.List()
		it = &_it
	default:
		it, _ = value.Fields(cue.Optional(true), cue.Hidden(true))
	}
	values := slices.IterToArray[cue.Iterator, cue.Value](it)
	sort.Slice(values, func(i, j int) bool {
		xOrder, yOrder := values[i].Attribute(orderKey), values[j].Attribute(orderKey)
		x, e1 := strconv.ParseInt(xOrder.Contents(), 10, 32)
		y, e2 := strconv.ParseInt(yOrder.Contents(), 10, 32)
		switch {
		case e1 != nil && e2 != nil:
			return i < j
		case e1 != nil:
			return false
		case e2 != nil:
			return true
		default:
			return x < y
		}
	})
	for _, val := range values {
		if Iterate(val, fn) {
			return true
		}
	}
	return fn(value)
}
