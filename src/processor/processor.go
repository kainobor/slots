package processor

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/kainobor/slots/config"
)

// Processor processing outcomes
type Processor struct {
	config       *config.ProcessorConfig
	random       *rand.Rand
	outcomeLines int
}

// SpinResult is struct for storing result of current spin
type SpinResult struct {
	Type   string `json:"type"`
	Total  int64  `json:"total"`
	Stops  []int  `json:"stops"`
	IsFree bool
}

// New processor with config "c"
func New(c *config.ProcessorConfig) *Processor {
	p := new(Processor)
	p.config = c

	rs := rand.NewSource(time.Now().Unix())
	p.random = rand.New(rs)

	// Process maximum number of line for outcome from config
	for _, line := range c.Lines {
		for _, k := range line {
			if k > p.outcomeLines {
				p.outcomeLines = k
			}
		}
	}

	return p
}

// GenerateOutcome return random outcome of display symbols
func (p *Processor) GenerateOutcome() (outcome []string) {
	for reelNum := 0; reelNum < len(p.config.Lines[0]); reelNum++ {
		cellNum := p.random.Intn(len(p.config.Reels))
		for i := 0; i < p.outcomeLines; i++ {
			currentCellNum := cellNum + i
			if currentCellNum > len(p.config.Reels)-1 {
				currentCellNum = currentCellNum - len(p.config.Reels)
			}
			outcome = append(outcome, p.config.Reels[currentCellNum][reelNum])
		}
	}

	return
}

// CheckWins returns results of current outcome and bet
func (p *Processor) CheckWins(outcome []string, bet int64) (*SpinResult, error) {
	return check(outcome, bet, p.config.Lines, p.config.EmptyCell, p.config.MutableCell, p.outcomeLines, p.config.Prices)
}

func check(outcome []string, bet int64, lines [20][]int, emptyCell string, mutCell string, outcomeLines int, prices map[string][4]int64) (*SpinResult, error) {
	var posInOutcome int
	res := new(SpinResult)
	res.Stops = make([]int, 0)
	for lineNum, positions := range lines {
		cellName := ""
		amount := 0

		for index, cellNum := range positions {
			posInOutcome = index*outcomeLines + cellNum - 1
			if cellName == "" {
				amount = 1
				cellName = outcome[posInOutcome]
			} else if cellName == outcome[posInOutcome] || outcome[posInOutcome] == mutCell {
				amount++
			} else {
				break
			}
		}

		if amount < 2 {
			continue
		}

		if cellName == emptyCell {
			if amount >= 3 {
				res.IsFree = true
			}
			continue
		}

		prices, ok := prices[cellName]
		if !ok {
			return nil, fmt.Errorf("wrong cell name: %s", cellName)
		}

		priceNum := amount - 2
		if prices[priceNum] > 0 {
			res.Total = res.Total + (prices[priceNum] * bet)
			res.Stops = append(res.Stops, lineNum+1)
		}
	}

	return res, nil
}
