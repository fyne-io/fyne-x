package r2

import (
	"math"
	"testing"
)

func TestVec2(t *testing.T) {
	// Test the Angle() function
	v1 := V2(1, 0)
	if v1.Angle() != 0 {
		t.Errorf("Angle of {1,0} failed. Expected 0, got %f", v1.Angle())
	}
	v1 = V2(1, 0.5)
	tolerance := 0.000001
	if math.Abs(v1.Angle()-0.463647) > tolerance {
		t.Errorf("Angle of {1,0.5} failed. Expected 0.463647, got %f", v1.Angle())
	}
	v1 = V2(-1, 0.5)
	if math.Abs(v1.Angle()-2.677945) > tolerance {
		t.Errorf("Angle of {-1,0.5} failed. Expected 2.677945, got %f", v1.Angle())
	}
	v1 = V2(-1, -0.5)
	if math.Abs(v1.Angle()+2.677945) > tolerance {
		t.Errorf("Angle of {1,0.5} failed. Expected -2.677945, got %f", v1.Angle())
	}
	v1 = V2(1, -0.5)
	if math.Abs(v1.Angle()+0.463647) > tolerance {
		t.Errorf("Angle of {1,0.5} failed. Expected -0.463647, got %f", v1.Angle())
	}
}

func TestAngleSum(t *testing.T) {
	tolerance := 0.000001
	if math.Abs(AddAngles(0.0, 0.0)) > tolerance {
		t.Error("AngleSum 0 failed")
	}
	if math.Abs(AddAngles(0.0, math.Pi)-math.Pi) > tolerance {
		t.Error("AngleSum 0, Pi failed")
	}
	if math.Abs(AddAngles(math.Pi, math.Pi)) > tolerance {
		t.Error("AngleSum Pi, Pi failed")
	}
	if math.Abs(AddAngles(-math.Pi, -math.Pi)) > tolerance {
		t.Error("AngleSum -Pi, -Pi failed")
	}
	if math.Abs(AddAngles(0.0, math.Pi/3)-math.Pi/3) > tolerance {
		t.Error("AngleSum 0, Pi/3 failed")
	}
	if math.Abs(AddAngles(math.Pi/3, math.Pi/3)-2*math.Pi/3) > tolerance {
		t.Error("AngleSum Pi/3, Pi/3 failed")
	}
	// The following sum could return either +Pi or -Pi
	if math.Abs(AddAngles(math.Pi/3, 2*math.Pi/3))-math.Pi > tolerance {
		t.Error("AngleSum Pi/3, 2*Pi/3 failed")
	}
	if math.Abs(AddAngles(2*math.Pi/3, 2*math.Pi/3)+2*math.Pi/3) > tolerance {
		t.Error("AngleSum 2*Pi/3, 2*Pi/3 failed")
	}
	if math.Abs(AddAngles(-2*math.Pi/3, -2*math.Pi/3)-2*math.Pi/3) > tolerance {
		t.Error("AngleSum -2*Pi/3, -2*Pi/3 failed")
	}
}
