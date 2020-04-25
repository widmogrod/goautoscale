package goautoscale

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/widmogrod/gonumutil"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
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

func TestAutoScalingVisualize(t *testing.T) {
	pCPU, err := plot.New()
	if err != nil {
		panic(err)
	}
	pInst, err := plot.New()
	if err != nil {
		panic(err)
	}

	pCPU.Title.Text = "Average CPU utilization in percentage"
	pCPU.X.Tick.Marker = gonumutil.NewConstantNumTicker(1)
	pCPU.Y.Tick.Marker = gonumutil.NewConstantNumTicker(5)

	pInst.Title.Text = "Auto Scaling decision to change number of instances with respect to average CPU utilisation "
	pInst.X.Tick.Marker = gonumutil.NewConstantNumTicker(1)
	pInst.Y.Tick.Marker = gonumutil.NewConstantNumTicker(1)

	ctx := Context{
		CPUNoopRange: Range{
			Min: 80,
			Max: 90,
		},
		MaintainsCPUAvg: 85,
		CPUUtilisation:  85,
		Instances:       3,
	}

	// Let's take a look at rate of change
	var utilization, boundaryMax, boundaryMin, instances plotter.XYs
	for i := 0; i < 50; i++ {
		var avgCPU float64 = 80
		if i >= 5 {
			avgCPU = 91
		}
		if i >= 10 {
			avgCPU = 90
		}
		if i >= 15 {
			avgCPU = 89
		}
		if i >= 20 {
			avgCPU = 98
		}
		if i >= 25 {
			avgCPU = 89
		}
		if i >= 30 {
			avgCPU = 81
		}
		if i >= 35 {
			avgCPU = 50
		}
		if i >= 37 {
			avgCPU = 84
		}

		ctx.CPUUtilisation = avgCPU

		scaleInstances := CPUScale(ctx)
		ctx.Instances += scaleInstances

		utilization = append(utilization, plotter.XY{
			X: float64(i),
			Y: ctx.CPUUtilisation,
		})

		instances = append(instances, plotter.XY{
			X: float64(i),
			Y: float64(scaleInstances),
		})

		boundaryMax = append(boundaryMax, plotter.XY{
			X: float64(i),
			Y: ctx.CPUNoopRange.Max,
		})
		boundaryMin = append(boundaryMin, plotter.XY{
			X: float64(i),
			Y: ctx.CPUNoopRange.Min,
		})
	}

	err = plotutil.AddLinePoints(pCPU,
		"CPU utilization", utilization,
		"CPU Max", boundaryMax,
		"CPU Min", boundaryMin,
	)
	if err != nil {
		t.Fatal(err)
	}

	err = plotutil.AddLinePoints(pInst,
		"Change of instances", instances,
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := pCPU.Save(9*vg.Inch, 4*vg.Inch, "autoscaling_by_test_cpu.png"); err != nil {
		t.Fatal(err)
	}

	if err := pInst.Save(9*vg.Inch, 4*vg.Inch, "autoscaling_by_test_inst.png"); err != nil {
		t.Fatal(err)
	}
}

// Positive number means to scale up, negative scale down, and zero don't do anything
func ExampleCPUScale() {
	ctx := Context{
		CPUNoopRange: Range{
			Min: 80,
			Max: 90,
		},
		MaintainsCPUAvg: 85,
		CPUUtilisation:  99,
		Instances:       3,
	}

	fmt.Println(CPUScale(ctx))
	// Output: 1
}
