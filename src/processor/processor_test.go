package processor

import (
	"testing"

	"gotest.tools/assert"
)

type Assertion struct {
	outcome []string
	total   int64
	stops   []int
	err     error
}

const (
	defaultBet = 5

	Atkins       = "Atkins"
	Steak        = "Steak"
	Ham          = "Ham"
	BuffaloWings = "BuffaloWings"
	Sausage      = "Sausage"
	Eggs         = "Eggs"
	Butter       = "Butter"
	Cheese       = "Cheese"
	Bacon        = "Bacon"
	Mayonnaise   = "Mayonnaise"
	Scale        = "Scale"
)

var (
	lines = [20][]int{
		{2, 2, 2, 2, 2},
		{1, 1, 1, 1, 1},
		{3, 3, 3, 3, 3},
		{1, 2, 3, 2, 1},
		{3, 2, 1, 2, 3},
		{2, 1, 1, 1, 2},
		{2, 3, 3, 3, 2},
		{1, 1, 2, 3, 3},
		{3, 3, 2, 1, 1},
		{2, 1, 2, 3, 2},
		{2, 3, 2, 1, 2},
		{1, 2, 2, 2, 1},
		{3, 2, 2, 2, 3},
		{1, 2, 1, 2, 1},
		{3, 2, 3, 2, 3},
		{2, 2, 1, 2, 2},
		{2, 2, 3, 2, 2},
		{1, 1, 3, 1, 1},
		{3, 3, 1, 3, 3},
		{1, 3, 3, 3, 1},
	}

	prices = map[string][4]int64{
		Atkins:       {5, 50, 500, 5000},
		Steak:        {3, 40, 200, 1000},
		Ham:          {2, 30, 150, 500},
		BuffaloWings: {2, 25, 100, 300},
		Sausage:      {0, 20, 75, 200},
		Eggs:         {0, 20, 75, 200},
		Butter:       {0, 15, 50, 100},
		Cheese:       {0, 15, 50, 100},
		Bacon:        {0, 10, 25, 50},
		Mayonnaise:   {0, 10, 25, 50},
	}
)

func TestCheck(t *testing.T) {
	assertions := []Assertion{
		{
			outcome: []string{
				Sausage, BuffaloWings, Steak,
				BuffaloWings, Scale, Mayonnaise,
				Atkins, BuffaloWings, Bacon,
				Cheese, BuffaloWings, Bacon,
				Bacon, Scale, Steak,
			},
			total: 250,
			stops: []int{6, 10},
			err:   nil,
		},
		{
			outcome: []string{
				BuffaloWings, Steak, Butter,
				Ham, Atkins, Butter,
				Mayonnaise, BuffaloWings, Ham,
				Atkins, Scale, Butter,
				Steak, Mayonnaise, Sausage,
			},
			total: 190,
			stops: []int{1, 4, 12, 14, 16, 17},
			err:   nil,
		},
		{
			outcome: []string{
				Eggs, Cheese, Mayonnaise,
				Ham, Atkins, Butter,
				Sausage, Bacon, Steak,
				Bacon, Cheese, Sausage,
				Steak, Ham, Cheese,
			},
			total: 0,
			stops: []int{},
			err:   nil,
		},
		{
			outcome: []string{
				Mayonnaise, Ham, Cheese,
				Sausage, Mayonnaise, Ham,
				Butter, Mayonnaise, Cheese,
				Ham, Sausage, Steak,
				Sausage, Ham, Butter,
			},
			total: 70,
			stops: []int{7, 11, 12},
			err:   nil,
		},
		{
			outcome: []string{
				Sausage, BuffaloWings, Steak,
				Eggs, BuffaloWings, Mayonnaise,
				Eggs, Bacon, Mayonnaise,
				Cheese, Atkins, Scale,
				Atkins, Butter, BuffaloWings,
			},
			total: 30,
			stops: []int{1, 16, 17},
			err:   nil,
		},
	}

	for _, a := range assertions {
		res, err := check(a.outcome, defaultBet, lines, Scale, Atkins, 3, prices)

		assert.NilError(t, err)
		assert.Equal(t, res.Total, a.total)
		assert.DeepEqual(t, res.Stops, a.stops)
	}
}
