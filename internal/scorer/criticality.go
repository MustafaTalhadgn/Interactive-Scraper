package scorer

// Criticality represents threat criticality levels
type Criticality string

const (
	CriticalityLow      Criticality = "low"
	CriticalityMedium   Criticality = "medium"
	CriticalityHigh     Criticality = "high"
	CriticalityCritical Criticality = "critical"
)

type CriticalityThresholds struct {
	Low      int 
	Medium   int 
	High     int 
	Critical int
}


func DefaultThresholds() *CriticalityThresholds {
	return &CriticalityThresholds{
		Low:      25,
		Medium:   50,
		High:     75,
		Critical: 100,
	}
}

func ClassifyCriticality(score int, thresholds *CriticalityThresholds) Criticality {
	if thresholds == nil {
		thresholds = DefaultThresholds()
	}
	
	switch {
	case score >= 76:
		return CriticalityCritical
	case score >= 51:
		return CriticalityHigh
	case score >= 26:
		return CriticalityMedium
	default:
		return CriticalityLow
	}
}


func GetCriticalityColor(criticality Criticality) string {
	switch criticality {
	case CriticalityCritical:
		return "#DC2626" 
	case CriticalityHigh:
		return "#F59E0B" 
	case CriticalityMedium:
		return "#FBBF24" 
	case CriticalityLow:
		return "#10B981" 
	default:
		return "#6B7280" 
	}
}

func GetCriticalityPriority(criticality Criticality) int {
	switch criticality {
	case CriticalityCritical:
		return 4
	case CriticalityHigh:
		return 3
	case CriticalityMedium:
		return 2
	case CriticalityLow:
		return 1
	default:
		return 0
	}
}