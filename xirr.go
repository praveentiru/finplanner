package finplanner

import (
	"fmt"
	"math"
	"time"
)

// Notes: Formula for converting daily return to annual return
// AR = (DR + 1)^365 - 1

type Investment struct {
	Date       time.Time `json:"investment_date,omitempty"`
	Investment float64   `json:"investment,omitempty"`
}

// Data needed to confrom to optimize.Problem of gonum/optimize package
type xirrData struct {
	Cashflow        []float64
	DaysInPortfolio []float64
}

// Calculates if error between portfolio value and calculated portfolio value with current
// xirr guess
func (xd *xirrData) Func(r float64) float64 {

	res := 0.0
	for i, cash := range xd.Cashflow {
		res += cash * math.Pow(1+r, xd.DaysInPortfolio[i])
	}
	return res
}

// Calculates derivative of function F
func (xd *xirrData) Grad(r float64) float64 {

	g := 0.0
	for i, cash := range xd.Cashflow {
		if xd.DaysInPortfolio[i] >= 0.9 {
			g += cash * xd.DaysInPortfolio[i] * math.Pow(1+r, xd.DaysInPortfolio[i]-1)
		}
	}

	return g
}

// XIRR calculates the XIRR return for cash flows provided as input
func XIRR(invs []Investment, currValue Investment) (float64, error) {

	xd, err := buildXIRRData(invs, currValue)
	if err != nil {
		return 0.0, err
	}

	res, err := runNewtonRaphson(0.003, &xd, 1500, 0.00001)
	if err != nil {
		return 0.0, err
	}
	// Convert daily return to annual return
	ar := math.Pow((res+1), 365) - 1

	return ar, nil
}

func buildXIRRData(invs []Investment, currValue Investment) (xirrData, error) {
	// There should be atleast one value in the array
	if len(invs) < 1 {
		return xirrData{}, fmt.Errorf("there should be atleast one investment")
	}

	// There should be atleast one +ve value in the investments
	v := false
	for _, inv := range invs {
		if inv.Investment > 0.0 {
			v = true
			break
		}
	}
	if !v {
		return xirrData{}, fmt.Errorf("there should be at least one +ve investment")
	}

	cash := make([]float64, 0, len(invs)+1)
	days := make([]float64, 0, len(invs)+1)
	cd := currValue.Date
	for _, inv := range invs {
		cash = append(cash, inv.Investment)
		d := cd.Sub(inv.Date)
		days = append(days, d.Hours()/24)
	}

	cash = append(cash, -currValue.Investment)
	days = append(days, 0)

	return xirrData{
		Cashflow:        cash,
		DaysInPortfolio: days,
	}, nil
}

func runNewtonRaphson(guess float64, xd *xirrData, maxIter int, errLimit float64) (float64, error) {

	x := guess
	f := 30000.0
	for i := 0; i < maxIter; i++ {
		f = xd.Func(x)
		if math.Abs(f) <= errLimit {
			break
		}

		x = x - xd.Func(x)/xd.Grad(x)
	}

	if math.Abs(f) > errLimit {
		return 0.0, fmt.Errorf("could not find XIRR after max iterations")
	}

	return x, nil
}
