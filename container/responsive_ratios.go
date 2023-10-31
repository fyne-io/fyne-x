package container

// Ratio is a helper type to define a fraction of a container. It's a float32 alias.
type Ratio = float32

// Ful l is the full size of the container.
func FullRatio() Ratio { return 1.0 }

// Half is the half size of the container.
func HalfRatio() Ratio { return 0.5 }

// OneThird is the one third size of the container.
func OneThirdRatio() Ratio { return 1.0 / 3.0 }

// TwoThird is the two thirds size of the container.
func TwoThirdRatio() Ratio { return OneThirdRatio() * 2.0 }

// OneQuarter is the one quarter size of the container.
func OneQuarterRatio() Ratio { return 0.25 }

// ThreeQuarter is the three quarters size of the container.
func ThreeQuarterRatio() Ratio { return 0.75 }

// OneFifth is the one fifth size of the container.
func OneFifthRatio() Ratio { return 0.2 }
