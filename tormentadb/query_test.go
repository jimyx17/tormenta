package tormentadb_test

import (
	"testing"
	"time"

	"github.com/jpincas/gouuidv6"
	"github.com/jpincas/tormenta/demo"
	tormenta "github.com/jpincas/tormenta/tormentadb"
)

// Basic Queries

func Test_BasicQuery(t *testing.T) {
	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()

	// 1 order
	order1 := demo.Order{}
	db.Save(&order1)

	var orders []demo.Order
	n, err := db.Find(&orders).Run()

	if err != nil {
		t.Error("Testing basic querying - got error")
	}

	if len(orders) != 1 || n != 1 {
		t.Errorf("Testing querying with 1 entity saved. Expecting 1 entity - got %v/%v", len(orders), n)
	}

	orders = []demo.Order{}
	c, err := db.Find(&orders).Count()
	if c != 1 {
		t.Errorf("Testing count 1 entity saved. Expecting 1 - got %v", c)
	}

	// 2 orders
	order2 := demo.Order{}
	db.Save(&order2)

	orders = []demo.Order{}
	if n, _ := db.Find(&orders).Run(); n != 2 {
		t.Errorf("Testing querying with 2 entity saved. Expecting 2 entities - got %v", n)
	}

	if c, _ := db.Find(&orders).Count(); c != 2 {
		t.Errorf("Testing count 2 entities saved. Expecting 2 - got %v", c)
	}
	if order1.ID == order2.ID {
		t.Errorf("Testing querying with 2 entities saved. 2 entities saved both have same ID")
	}
	if orders[0].ID == orders[1].ID {
		t.Errorf("Testing querying with 2 entities saved. 2 results returned. Both have same ID")
	}

	// Limit
	orders = []demo.Order{}
	if n, _ := db.Find(&orders).Limit(1).Run(); n != 1 {
		t.Errorf("Testing querying with 2 entities saved + limit. Wrong number of results received")
	}

	// Reverse - simple, only tests number received
	orders = []demo.Order{}
	if n, _ := db.Find(&orders).Reverse().Run(); n != 2 {
		t.Errorf("Testing querying with 2 entities saved + reverse. Expected %v, got %v", 2, n)
	}

	// Reverse + Limit - simple, only tests number received
	orders = []demo.Order{}
	if n, _ := db.Find(&orders).Reverse().Limit(1).Run(); n != 1 {
		t.Errorf("Testing querying with 2 entities saved + reverse + limit. Expected %v, got %v", 1, n)
	}

}

func Test_BasicQuery_First(t *testing.T) {
	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()

	order1 := demo.Order{}
	order2 := demo.Order{}
	db.Save(&order1, &order2)

	var order demo.Order
	n, err := db.First(&order).Run()

	if err != nil {
		t.Error("Testing first - got error")
	}

	if n != 1 {
		t.Errorf("Testing first. Expecting 1 entity - got %v", n)
	}

	if order.ID.IsNil() {
		t.Errorf("Testing first. Nil ID retrieved")
	}

	if order.ID != order1.ID {
		t.Errorf("Testing first. Order IDs are not equal - wrong order retrieved")
	}

	// Test nothing found (impossible range)
	n, _ = db.First(&order).From(time.Now()).To(time.Now()).Run()
	if n != 0 {
		t.Errorf("Testing first when nothing should be found.  Got n = %v", n)
	}
}

func Test_BasicQuery_DateRange(t *testing.T) {
	// Create a list of orders over a date range
	var orders []tormenta.Tormentable
	dates := []time.Time{
		// Now
		time.Now(),

		// Over the last week
		time.Now().Add(-1 * 24 * time.Hour),
		time.Now().Add(-2 * 24 * time.Hour),
		time.Now().Add(-3 * 24 * time.Hour),
		time.Now().Add(-4 * 24 * time.Hour),
		time.Now().Add(-5 * 24 * time.Hour),
		time.Now().Add(-6 * 24 * time.Hour),
		time.Now().Add(-7 * 24 * time.Hour),

		// Specific years
		time.Date(2009, time.January, 1, 1, 0, 0, 0, time.UTC),
		time.Date(2010, time.January, 1, 1, 0, 0, 0, time.UTC),
		time.Date(2011, time.January, 1, 1, 0, 0, 0, time.UTC),
		time.Date(2012, time.January, 1, 1, 0, 0, 0, time.UTC),
		time.Date(2013, time.January, 1, 1, 0, 0, 0, time.UTC),
	}

	for _, date := range dates {
		orders = append(orders, &demo.Order{
			Model: tormenta.Model{
				ID: gouuidv6.NewFromTime(date),
			},
		})
	}

	// Save the orders
	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()
	db.Save(orders...)

	// Also another entity, to make sure there is no crosstalk
	db.Save(&demo.Product{
		Code:          "001",
		Name:          "Computer",
		Price:         999.99,
		StartingStock: 50,
		Description:   demo.DefaultDescription})

	// Quick check that all orders have saved correctly
	var results []demo.Order
	n, _ := db.Find(&results).Run()

	if len(results) != len(orders) || n != len(orders) {
		t.Errorf("Testing range query. Haven't even got to ranges yet. Just basic query expected %v - got %v/%v", len(orders), len(results), n)
		t.FailNow()
	}

	// Range test cases
	testCases := []struct {
		testName  string
		from, to  time.Time
		expected  int
		includeTo bool
		limit     int
		reverse   bool
		offset    int
	}{
		{"from right now - no orders expected, no 'to'", time.Now(), time.Time{}, 0, false, 0, false, 0},
		{"from beginning of time - all orders should be included, no 'to'", time.Time{}, time.Time{}, len(orders), false, 0, false, 0},
		{"from beginning of time - offset 1", time.Time{}, time.Time{}, len(orders) - 1, false, 0, false, 1},
		{"from beginning of time - offset 2", time.Time{}, time.Time{}, len(orders) - 2, false, 0, false, 2},
		{"from 2014, no 'to'", time.Date(2014, time.January, 1, 1, 0, 0, 0, time.UTC), time.Time{}, 8, false, 0, false, 0},
		{"from 1 hour ago, no 'to'", time.Now().Add(-1 * time.Hour), time.Time{}, 1, false, 0, false, 0},
		{"from beginning of time to now - expect all", time.Time{}, time.Now(), len(orders), true, 0, false, 0},
		{"from beginning of time to 2014 - expect 5", time.Time{}, time.Date(2014, time.January, 1, 1, 0, 0, 0, time.UTC), 5, true, 0, false, 0},
		{"from beginning of time to an hour ago - expect all but 1", time.Time{}, time.Now().Add(-1 * time.Hour), len(orders) - 1, true, 0, false, 0},
		{"from beginning of time - limit 1", time.Time{}, time.Time{}, 1, false, 1, false, 0},
		{"from beginning of time - limit 10", time.Time{}, time.Time{}, 10, false, 10, false, 0},
		{"from beginning of time - limit 10 - offset 2 (shouldnt affect number of results)", time.Time{}, time.Time{}, 10, false, 10, false, 2},
		{"from beginning of time - limit more than there are", time.Time{}, time.Time{}, len(orders), false, 0, false, 0},
		{"reversed - from beginning of time", time.Time{}, time.Time{}, 0, false, 0, true, 0},
		{"reverse - from now - no to", time.Now(), time.Time{}, len(orders), false, 0, true, 0},
		{"reverse - from now to 2014 - expect 8", time.Now(), time.Date(2014, time.January, 1, 1, 0, 0, 0, time.UTC), 8, true, 0, true, 0},
		{"reverse - from now to 2014 - limit 5 - expect 5", time.Now(), time.Date(2014, time.January, 1, 1, 0, 0, 0, time.UTC), 5, true, 5, true, 0},
	}

	for _, testCase := range testCases {
		rangequeryResults := []demo.Order{}
		query := db.Find(&rangequeryResults).From(testCase.from)

		if testCase.includeTo {
			query = query.To(testCase.to)
		}

		if testCase.limit > 0 {
			query = query.Limit(testCase.limit)
		}

		if testCase.reverse {
			query = query.Reverse()
		}

		if testCase.offset > 0 {
			query = query.Offset(testCase.offset)
		}

		n, _ := query.Run()
		c, _ := query.Count()

		// Count should always equal number of results
		if c != n {
			t.Errorf("Testing %s. Number of results does not equal count. Count: %v, Results: %v", testCase.testName, c, n)
		}

		// Test number of records retrieved
		if n != testCase.expected {
			t.Errorf("Testing %s (number orders retrieved). Expected %v - got %v", testCase.testName, testCase.expected, n)
		}

		// Test Count
		if c != testCase.expected {
			t.Errorf("Testing %s (count). Expected %v - got %v", testCase.testName, testCase.expected, c)
		}

	}

}
