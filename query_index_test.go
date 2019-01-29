package tormenta_test

import (
	"testing"

	"github.com/jpincas/tormenta"
	"github.com/jpincas/tormenta/testtypes"
)

// Helper for making groups of depatments
func getDept(i int) int {
	if i <= 10 {
		return 1
	} else if i <= 20 {
		return 2
	} else {
		return 3
	}
}

// Test aggregation on an index
func Test_Aggregation(t *testing.T) {
	var fullStructs []tormenta.Record

	for i := 1; i <= 30; i++ {
		fullStruct := &testtypes.FullStruct{
			FloatField: float64(i),
			IntField:   i,
		}

		fullStructs = append(fullStructs, fullStruct)
	}

	tormenta.RandomiseRecords(fullStructs)

	db, _ := tormenta.OpenTest("data/tests", tormenta.DefaultOptions)
	defer db.Close()
	db.Save(fullStructs...)

	results := []testtypes.FullStruct{}
	var intSum int32
	var floatSum float64
	expected := 465

	// Int32

	_, err := db.Find(&results).Range("intfield", 1, 30).QuickSum(&intSum)
	if err != nil {
		t.Error("Testing int32 agreggation.  Got error")
	}

	expectedIntSum := int32(expected)
	if intSum != expectedIntSum {
		t.Errorf("Testing int32 agreggation. Expteced %v, got %v", expectedIntSum, intSum)
	}

	// Float64

	_, err = db.Find(&results).Range("floatfield", 1.00, 30.00).QuickSum(&floatSum)
	if err != nil {
		t.Error("Testing float64 agreggation.  Got error")
	}

	expectedFloatSum := float64(expected)
	if floatSum != expectedFloatSum {
		t.Errorf("Testing float64 agreggation. Expteced %v, got %v", expectedFloatSum, floatSum)
	}
}
