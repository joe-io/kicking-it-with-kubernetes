package main

import (
	"bufio"
	"bytes"
	"errors"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
)

type ClassifyResult struct {
	Url    string        `json:"url"`
	Labels []LabelResult `json:"labels"`
}

type LabelResult struct {
	Label       string  `json:"label"`
	Probability float32 `json:"probability"`
}

var (
	graphModel   *tf.Graph
	sessionModel *tf.Session
	labels       []string
)

func classifyImage(url string) (*ClassifyResult, error) {
	// Fetch the image at the URL
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.New("Could not download the image.")
	}
	// Ensure that we close the body, when this function exits
	defer resp.Body.Close()

	// Read the image data into a buffer
	var imageBuffer bytes.Buffer
	io.Copy(&imageBuffer, resp.Body)

	// Make tensor
	tensor, err := makeTensorFromImage(&imageBuffer, url)
	if err != nil {
		return nil, errors.New("Invalid image")
	}

	// Run inference
	output, err := sessionModel.Run(
		map[tf.Output]*tf.Tensor{
			graphModel.Operation("input").Output(0): tensor,
		},
		[]tf.Output{
			graphModel.Operation("output").Output(0),
		},
		nil)
	if err != nil {
		return nil, errors.New("Could not run inference")
	}

	// Return best labels
	return &ClassifyResult{
		Url:    url,
		Labels: findBestLabels(output[0].Value().([][]float32)[0]),
	}, nil
}

func loadModel() error {
	// Load inception model
	model, err := ioutil.ReadFile("./model/tensorflow_inception_graph.pb")
	if err != nil {
		return err
	}
	graphModel = tf.NewGraph()
	if err := graphModel.Import(model, ""); err != nil {
		return err
	}

	sessionModel, err = tf.NewSession(graphModel, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Load labels
	labelsFile, err := os.Open("./model/imagenet_comp_graph_label_strings.txt")
	if err != nil {
		return err
	}
	defer labelsFile.Close()
	scanner := bufio.NewScanner(labelsFile)
	// Labels are separated by newlines
	for scanner.Scan() {
		labels = append(labels, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

type ByProbability []LabelResult

func (a ByProbability) Len() int           { return len(a) }
func (a ByProbability) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByProbability) Less(i, j int) bool { return a[i].Probability > a[j].Probability }

func findBestLabels(probabilities []float32) []LabelResult {
	// Make a list of label/probability pairs
	var resultLabels []LabelResult
	for i, p := range probabilities {
		if i >= len(labels) {
			break
		}
		resultLabels = append(resultLabels, LabelResult{Label: labels[i], Probability: p})
	}
	// Sort by probability
	sort.Sort(ByProbability(resultLabels))
	// Return top 5 labels
	return resultLabels[:5]
}

// Based on the great work done by @tinrab here: https://github.com/tinrab/go-tensorflow-image-recognition
