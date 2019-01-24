package main

import (
	"fmt"
	tg "github.com/galeone/tfgo"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

func main() {
	model := tg.LoadModel("test_models/export", []string{"tag"}, nil)

	fakeInput, _ := tf.NewTensor([1][28][28][1]float32{})
	results := model.Exec([]tf.Output{
		model.Op("LeNetDropout/softmax_linear/Identity", 0),
	}, map[tf.Output]*tf.Tensor{
		model.Op("input_", 0): fakeInput,
	})

	predictions := results[0].Value().([][]float32)
	fmt.Println(predictions)
}
