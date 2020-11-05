package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name:  "OK",
		Price: 1.0,
		SKU:   "abc-abc-def",
	}
	err := p.Validate()
	if err != nil {
		t.Fatal(err)
	}
}
