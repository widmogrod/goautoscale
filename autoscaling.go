package goautoscale

import "math"

// percent is a value between 0 and 100
type percent = float64

// Range in which a value should hold
type Range struct {
	Min percent
	Max percent
}

// Contains verifies if x is in range, when is then returns true, otherwise false
func (ir Range) Contains(x percent) bool {
	if ir.Min > x || ir.Max < x {
		return false
	}

	return true
}

// Context for scaling operation
type Context struct {
	// CPUNoopRange defines boundaries in which CPU should not trigger auto-scaling
	CPUNoopRange Range
	// MaintainsCPUAvg level of CPU utilisation that should be maintained when decision about scaling up or down in made
	MaintainsCPUAvg percent
	// CPUUtilisation represents current CPU utilisation.
	CPUUtilisation percent
	// Instances represents current number of instances of a service.
	Instances int
}

// CPUScale calculates how many instances should be added or removed to maintain given CPU utilization
// A positive number means to scale up, negative scale down, and zero doesn't do anything.
func CPUScale(in Context) int {
	if in.CPUNoopRange.Contains(in.CPUUtilisation) {
		return 0
	}

	return ScaleInstances(float64(in.Instances), in.CPUUtilisation, in.MaintainsCPUAvg)
}

// ScaleInstances calculates how many instances should be added or removed
// to maintain given percentage of utilization of resource with respect to current utilization.
// Utilisation is abstract, and it can be applied to average CPU utilisation, average queue size,...
func ScaleInstances(instances, utilisation, maintain percent) int {
	candidate := instances * utilisation / maintain
	candidate = math.Ceil(candidate - instances)
	return int(candidate)
}

type Recommendation struct {
	ScaleUp   uint
	ScaleDown uint
}

func ToRecommendation(candidate int) Recommendation {
	result := Recommendation{}
	if candidate < 0 {
		result.ScaleDown = uint(-candidate)
	} else if candidate > 0 {
		result.ScaleUp = uint(candidate)
	}

	return result
}
