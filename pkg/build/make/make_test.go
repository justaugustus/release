/*
Copyright 2020 The Kubernetes Authors.

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

package make_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"k8s.io/release/pkg/build/make"
	"k8s.io/release/pkg/build/make/makefakes"
)

var err = errors.New("error")

func TestMakeCross(t *testing.T) {
	for _, tc := range []struct {
		prepare     func(*makefakes.FakeImpl)
		shouldError bool
	}{
		{ // success
			prepare:     func(*makefakes.FakeImpl) {},
			shouldError: false,
		},
		{ // OpenRepo fails
			prepare: func(mock *makefakes.FakeImpl) {
				mock.OpenRepoReturns(nil, err)
			},
			shouldError: true,
		},
		{ // Checkout fails
			prepare: func(mock *makefakes.FakeImpl) {
				mock.CheckoutReturns(err)
			},
			shouldError: true,
		},
		{ // Command fails
			prepare: func(mock *makefakes.FakeImpl) {
				mock.CommandReturns(err)
			},
			shouldError: true,
		},
		{ // Rename fails
			prepare: func(mock *makefakes.FakeImpl) {
				mock.RenameReturns(err)
			},
			shouldError: true,
		},
		{ // Command fails on second call
			prepare: func(mock *makefakes.FakeImpl) {
				mock.CommandReturnsOnCall(1, err)
			},
			shouldError: true,
		},
	} {
		sut := make.New()
		mock := &makefakes.FakeImpl{}
		tc.prepare(mock)
		sut.SetImpl(mock)
		err := sut.MakeCross("v1.20.0")
		if tc.shouldError {
			require.NotNil(t, err)
		} else {
			require.Nil(t, err)
		}
	}
}
