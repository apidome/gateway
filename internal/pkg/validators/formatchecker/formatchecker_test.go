package formatchecker

import "testing"

type data struct {
	data  string
	valid bool
}

func TestIsValidDateTime(t *testing.T) {
	slideData := make([]data, 0)
	t1 := data{
		data:  "1985-04-12T23:20:50.52Z",
		valid: true,
	}
	t2 := data{
		data:  "1996-12-19T16:39:57-08:00",
		valid: true,
	}
	t3 := data{
		data:  "06/19/1963 08:30:06 PST",
		valid: false,
	}
	slideData = append(slideData, t1, t2, t3)

	for _, d := range slideData {
		if valid, _ := IsValidDateTime(d.data); valid != d.valid {
			var valid string
			if !d.valid {
				valid = " not"
			}
			t.Errorf("validate %s against date-time expected to%s be valid", d.data, valid)
		}
	}
}
