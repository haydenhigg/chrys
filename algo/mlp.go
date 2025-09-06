/**
 * This is not production-ready. It is still in progress, and it may not get
 * finished at all. USE AT YOUR OWN RISK!
 */

package algo

import (
	// "fmt"
	"math"
	"math/rand"
)

// Activation function
func sigmoid(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}

func sigmoidDerivative(s float64) float64 {
	return s * (1 - s)
}

// MLP structure
type MLP struct {
	inputSize, hiddenSize, outputSize int
	weightsInputHidden                [][]float64
	weightsHiddenOutput               [][]float64
	learningRate                      float64
}

// Random float between -1 and 1
func randWeight() float64 {
	return rand.Float64()*2 - 1
}

// Create a new MLP with random weights
func NewMLP(inputSize, hiddenSize, outputSize int, learningRate float64) *MLP {
	wIH := make([][]float64, inputSize)
	for i := range wIH {
		wIH[i] = make([]float64, hiddenSize)
		for j := range wIH[i] {
			wIH[i][j] = randWeight()
		}
	}

	wHO := make([][]float64, hiddenSize)
	for i := range wHO {
		wHO[i] = make([]float64, outputSize)
		for j := range wHO[i] {
			wHO[i][j] = randWeight()
		}
	}

	return &MLP{
		inputSize:           inputSize,
		hiddenSize:          hiddenSize,
		outputSize:          outputSize,
		weightsInputHidden:  wIH,
		weightsHiddenOutput: wHO,
		learningRate:        learningRate,
	}
}

// Forward pass
func (mlp *MLP) forward(input []float64) ([]float64, []float64) {
	// Hidden layer
	hidden := make([]float64, mlp.hiddenSize)
	for i := 0; i < mlp.hiddenSize; i++ {
		sum := 0.
		for j := 0; j < mlp.inputSize; j++ {
			sum += input[j] * mlp.weightsInputHidden[j][i]
		}
		hidden[i] = sigmoid(sum)
	}

	// Output layer
	output := make([]float64, mlp.outputSize)
	for i := 0; i < mlp.outputSize; i++ {
		sum := 0.
		for j := 0; j < mlp.hiddenSize; j++ {
			sum += hidden[j] * mlp.weightsHiddenOutput[j][i]
		}
		output[i] = sigmoid(sum)
	}

	return hidden, output
}

// Train with one data point using backpropagation
func (mlp *MLP) Train(input []float64, target []float64) {
	hidden, output := mlp.forward(input)

	// Calculate output errors
	outputErrors := make([]float64, mlp.outputSize)
	for i := range outputErrors {
		outputErrors[i] = (target[i] - output[i]) * sigmoidDerivative(output[i])
	}

	// Calculate hidden errors
	hiddenErrors := make([]float64, mlp.hiddenSize)
	for i := range hiddenErrors {
		sum := 0.
		for j := 0; j < mlp.outputSize; j++ {
			sum += outputErrors[j] * mlp.weightsHiddenOutput[i][j]
		}
		hiddenErrors[i] = sum * sigmoidDerivative(hidden[i])
	}

	// Update weights hidden->output
	for i := 0; i < mlp.hiddenSize; i++ {
		for j := 0; j < mlp.outputSize; j++ {
			mlp.weightsHiddenOutput[i][j] += mlp.learningRate * outputErrors[j] * hidden[i]
		}
	}

	// Update weights input->hidden
	for i := 0; i < mlp.inputSize; i++ {
		for j := 0; j < mlp.hiddenSize; j++ {
			mlp.weightsInputHidden[i][j] += mlp.learningRate * hiddenErrors[j] * input[i]
		}
	}
}

// Predict (inference only)
func (mlp *MLP) Predict(input []float64) []float64 {
	_, output := mlp.forward(input)
	return output
}

// MSE loss
func MSE(predicted, target []float64) float64 {
	sum := 0.
	for i := range predicted {
		diff := predicted[i] - target[i]
		sum += diff * diff
	}
	return sum / float64(len(predicted))
}

// // Main function
// func main() {
// 	mlp := NewMLP(2, 4, 1, 0.9)

// 	// XOR training data
// 	inputs := [][]float64{
// 		{0, 0},
// 		{0, 1},
// 		{1, 0},
// 		{1, 1},
// 	}
// 	targets := [][]float64{
// 		{0},
// 		{1},
// 		{1},
// 		{0},
// 	}

// // Training loop
// epochs := 10000
// for epoch := 1; epoch <= epochs; epoch++ {
// 	totalLoss := 0.
// 	for i := range inputs {
// 		mlp.Train(inputs[i], targets[i])
// 		pred := mlp.Predict(inputs[i])
// 		totalLoss += mse(pred, targets[i])
// 	}
// 	if epoch%1000 == 0 {
// 		fmt.Printf("Epoch %d: Loss = %.4f\n", epoch, totalLoss)
// 	}
// }

// 	// Test
// 	fmt.Println("\nTrained predictions:")
// 	for i, input := range inputs {
// 		output := mlp.Predict(input)
// 		fmt.Printf("Input: %v -> Output: %.4f (Target: %.0f)\n", input, output[0], targets[i][0])
// 	}
// }

//
// 🧠 What’s Included
// Single hidden layer MLP
// Manual forward and backward pass
// Gradient descent training loop
// Sigmoid activation
//
// 🛠 What’s Missing (can be added if needed)
// Bias terms (can improve convergence)
// Multiple hidden layers
// Other activation functions (e.g., ReLU)
// Batch training
// Model saving/loading
//
