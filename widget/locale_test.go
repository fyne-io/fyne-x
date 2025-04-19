package widget

import "testing"

func Test_minusRadixThou(t *testing.T) {
	minus, radix, thou := minusRadixThou("en-US")
	if minus != '-' {
		t.Errorf("minus should be '-' but is %x", minus)
	}
	if radix != '.' {
		t.Errorf("radix should be '.' but is %x", radix)
	}
	if thou != ',' {
		t.Errorf("thou should be ',' but is %x", thou)
	}

	minus, radix, thou = minusRadixThou("de-DE")
	if minus != '-' {
		t.Errorf("minus should be '-' but is %x", minus)
	}
	if radix != ',' {
		t.Errorf("radix should be ',' but is %x", radix)
	}
	if thou != '.' {
		t.Errorf("thou should be '.' but is %x", thou)
	}

	minus, radix, thou = minusRadixThou("de-AT")
	if minus != '-' {
		t.Errorf("minus should be '-' but is %x", minus)
	}
	if radix != ',' {
		t.Errorf("radix should be ',' but is %x", radix)
	}
	if thou != 0xa0 {
		t.Errorf("thou should be '0xa0 but is %x", thou)
	}

	minus, radix, thou = minusRadixThou("de-LI")
	if minus != '-' {
		t.Errorf("minus should be '-' but is %x", minus)
	}
	if radix != '.' {
		t.Errorf("radix should be '.' but is %x", radix)
	}
	if thou != 0x2019 {
		t.Errorf("thou should be 0x2019 but is %x", thou)
	}

	minus, radix, thou = minusRadixThou("fr-FR")
	if minus != '-' {
		t.Errorf("minus should be '-' but is %x", minus)
	}
	if radix != ',' {
		t.Errorf("radix should be ',' but is %x", radix)
	}
	if thou != 0xa0 {
		t.Errorf("thou should be 0xa0 but is %x", thou)
	}

	minus, radix, thou = minusRadixThou("et-EE")
	if minus != 0x2212 {
		t.Errorf("minus should be 0x2212 but is %x", minus)
	}
	if radix != ',' {
		t.Errorf("radix should be ',' but is %x", radix)
	}
	if thou != 0xa0 {
		t.Errorf("thou should be 0xa0 but is %x", thou)
	}

	minus, radix, thou = minusRadixThou("rw-RW")
	if minus != '-' {
		t.Errorf("minus should be '-' but is %x", minus)
	}
	if radix != ',' {
		t.Errorf("radix should be ',' but is %x", radix)
	}
	if thou != '.' {
		t.Errorf("thou should be 0xa0 but is %x", thou)
	}

	minus, radix, thou = minusRadixThou("tr-TR")
	if minus != '-' {
		t.Errorf("minus should be '-' but is %x", minus)
	}
	if radix != ',' {
		t.Errorf("radix should be ',' but is %x", radix)
	}
	if thou != '.' {
		t.Errorf("thou should be 0xa0 but is %x", thou)
	}

	minus, radix, thou = minusRadixThou("tk-TM")
	if minus != '-' {
		t.Errorf("minus should be '-' but is %x", minus)
	}
	if radix != ',' {
		t.Errorf("radix should be ',' but is %x", radix)
	}
	if thou != 0xa0 {
		t.Errorf("thou should be 0xa0 but is %x", thou)
	}

	minus, radix, thou = minusRadixThou("kea-CV")
	if minus != '-' {
		t.Errorf("minus should be '-' but is %x", minus)
	}
	if radix != ',' {
		t.Errorf("radix should be ',' but is %x", radix)
	}
	if thou != 0xa0 {
		t.Errorf("thou should be 0xa0 but is %x", thou)
	}

	minus, radix, thou = minusRadixThou("mas-KE")
	if minus != '-' {
		t.Errorf("minus should be '-' but is %x", minus)
	}
	if radix != '.' {
		t.Errorf("radix should be '.' but is %x", radix)
	}
	if thou != ',' {
		t.Errorf("thou should be ',' but is %x", thou)
	}
}
