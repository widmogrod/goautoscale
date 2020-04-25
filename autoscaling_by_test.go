package goautoscale

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAutoScalingBy(t *testing.T) {
	useCases := map[string]struct {
		ctx      Context
		expected Recommendation
	}{
		"maintain": {
			ctx: Context{
				CPUNoopRange: Range{
					Min: 80,
					Max: 90,
				},
				MaintainsCPUAvg: 85,
				CPUUtilisation:  85,
				Instances:       3,
			},
			expected: Recommendation{
				ScaleUp:   0,
				ScaleDown: 0,
			},
		},
		"CPUScale up - small": {
			ctx: Context{
				CPUNoopRange: Range{
					Min: 80,
					Max: 90,
				},
				MaintainsCPUAvg: 85,
				CPUUtilisation:  91,
				Instances:       3,
			},
			expected: Recommendation{
				ScaleUp:   1,
				ScaleDown: 0,
			},
		},
		"CPUScale up - big": {
			ctx: Context{
				CPUNoopRange: Range{
					Min: 80,
					Max: 90,
				},
				MaintainsCPUAvg: 85,
				CPUUtilisation:  99,
				Instances:       30,
			},
			expected: Recommendation{
				ScaleUp:   5,
				ScaleDown: 0,
			},
		},
		"CPUScale down - small": {
			ctx: Context{
				CPUNoopRange: Range{
					Min: 80,
					Max: 90,
				},
				MaintainsCPUAvg: 85,
				CPUUtilisation:  67,
				Instances:       4,
			},
			expected: Recommendation{
				ScaleUp:   0,
				ScaleDown: 0,
			},
		},
		"CPUScale down - big": {
			ctx: Context{
				CPUNoopRange: Range{
					Min: 80,
					Max: 90,
				},
				MaintainsCPUAvg: 85,
				CPUUtilisation:  71,
				Instances:       33,
			},
			expected: Recommendation{
				ScaleUp:   0,
				ScaleDown: 5,
			},
		},
	}
	for name, uc := range useCases {
		t.Run(name, func(t *testing.T) {
			result := ToRecommendation(CPUScale(uc.ctx))
			assert.Equal(t, uc.expected, result)
		})
	}
}

// The example shows recommendation when a service needs to be scaled up due to CPU utilization bellow recommended
func ExampleCPUScale_up() {
	ctx := Context{
		// Auto scaling will not trigger when CPU utilisation is in given range
		CPUNoopRange: Range{
			Min: 80,
			Max: 90,
		},
		MaintainsCPUAvg: 85,
		CPUUtilisation:  99,
		Instances:       30,
	}

	fmt.Println(CPUScale(ctx))
	// Output: 5
}

// The example shows recommendation when a service needs to be scaled down due to CPU utilization bellow recommended
func ExampleCPUScale_down() {
	ctx := Context{
		// Auto scaling will not trigger when CPU utilisation is in given range
		CPUNoopRange: Range{
			Min: 80,
			Max: 90,
		},
		MaintainsCPUAvg: 85,
		CPUUtilisation:  50,
		Instances:       30,
	}

	fmt.Println(CPUScale(ctx))
	// Output: -12
}

// The example shows recommendation when a service don't need any change in number of instance
// due to CPI utilization being in range.
func ExampleCPUScale_noop() {
	ctx := Context{
		// Auto scaling will not trigger when CPU utilisation is in given range
		CPUNoopRange: Range{
			Min: 80,
			Max: 90,
		},
		MaintainsCPUAvg: 85,
		CPUUtilisation:  88,
		Instances:       30,
	}

	fmt.Println(CPUScale(ctx))
	// Output: 0
}
