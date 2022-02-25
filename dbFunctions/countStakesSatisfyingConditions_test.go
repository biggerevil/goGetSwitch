package dbFunctions

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"goGetSwitch/producerCode"
	"testing"
)

func TestGetPairnamesFromCombination(t *testing.T) {
	// Arrange
	testTable := []struct {
		combinationToTest producerCode.Combination
		expected          []string
	}{
		// Проверяем на пустой комбинации
		{
			combinationToTest: producerCode.Combination{
				Conditions: []producerCode.Condition{{}}},
			expected: []string(nil),
		},
		// Проверяем на комбинации с одной парой
		{
			combinationToTest: producerCode.Combination{
				Conditions: []producerCode.Condition{{"Pairname", "EUR/JPY"}}},
			expected: []string{"EUR/JPY"},
		},
		// Проверяем на комбинации с двумя парами
		{
			combinationToTest: producerCode.Combination{
				Conditions: []producerCode.Condition{{"Pairname", "EUR/JPY"}, {"Pairname", "EUR/USD"}}},
			expected: []string{"EUR/JPY", "EUR/USD"},
		},
		// Проверяем на комбинации с двумя парами и таймфреймом (который не должен быть возвращён)
		{
			combinationToTest: producerCode.Combination{
				Conditions: []producerCode.Condition{{"Timeframe", "900"}, {"Pairname", "EUR/JPY"}, {"Pairname", "EUR/USD"}}},
			expected: []string{"EUR/JPY", "EUR/USD"},
		},
	}

	for _, testCase := range testTable {
		// Act
		result := getPairnamesFromCombination(testCase.combinationToTest)

		//t.Logf("Calling Max(%v), result %d\n", testCase.numbers, result)

		// Assert
		assert.Equal(t, testCase.expected, result,
			fmt.Sprintf("Incorrect result. Expect %v, got %v", testCase.expected, result))
	}
}
