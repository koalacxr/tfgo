/*
Copyright 2017 Paolo Galeone. All right reserved.
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

package tfgo

import (
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"github.com/tensorflow/tensorflow/tensorflow/go/op"
)

// Batchify creates a batch of tensors, concatenating them along the first dimension
func Batchify(scope *op.Scope, tensors []tf.Output) tf.Output {
	s := scope.SubScope("batchify")
	var tensors4d []tf.Output
	for _, tensor := range tensors {
		tensors4d = append(tensors4d, op.ExpandDims(s.SubScope("ExpandDims"), tensor, op.Const(s.SubScope("axis"), []int32{0})))
	}
	return op.ConcatV2(s.SubScope("ConcatV2"), tensors4d, op.Const(s.SubScope("axis"), int32(0)))
}

// Cast casts value to the specified dtype
func Cast(scope *op.Scope, value tf.Output, dtype tf.DataType) tf.Output {
	if value.DataType() == dtype {
		return value
	}
	return op.Cast(scope.SubScope("Cast"), value, dtype)
}

// NewRoot creates a new *op.Scope, empty
func NewRoot() *op.Scope {
	return op.NewScope()
}

// Const creates a constant value within the specified scope
func Const(scope *op.Scope, value interface{}) tf.Output {
	return op.Const(scope.SubScope("Const"), value)
}

// IsClose defines the isclose operation between a and b.
// Returns a conditional node that is true when a is close to b.
// relTol is the relative tolerance
// absTol is the absolute tolerance
func IsClose(scope *op.Scope, a, b tf.Output, relTol, absTol tf.Output) tf.Output {
	s := scope.SubScope("IsClose")
	return op.LessEqual(s.SubScope("LessEqual"),
		op.Abs(s.SubScope("Abs"),
			op.Sub(s.SubScope("Sub"), a, b)),
		op.Maximum(s.SubScope("Maximum"),
			op.Mul(s.SubScope("Mul"), relTol,
				op.Maximum(s.SubScope("Maximum"), op.Abs(s.SubScope("Abs"), a), op.Abs(s.SubScope("Abs"), b))), absTol))
}

// Exec creates the computation graph from the scope, then executes
// the operations required to compute each element of tensors.
// Node in the graph can be overwritten with feedDict.
// The session options can be specified using the session parameter.
// Returns the evaluated tensors. Panics on error.
func Exec(scope *op.Scope, tensors []tf.Output, feedDict map[tf.Output]*tf.Tensor, options *tf.SessionOptions) []*tf.Tensor {
	graph, err := scope.Finalize()
	if err != nil {
		panic(err.Error())
	}

	if sess, err := tf.NewSession(graph, options); err == nil {
		if results, err := sess.Run(feedDict, tensors, nil); err == nil {
			return results
		}
		panic(err)
	} else {
		panic(err)
	}
}
