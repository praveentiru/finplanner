package finplanner

import (
	"math"
	"reflect"
	"testing"
	"time"
)

func Test_buildXIRRData(t *testing.T) {
	type args struct {
		invs      []Investment
		currValue Investment
	}
	tests := []struct {
		name    string
		args    args
		want    xirrData
		wantErr bool
	}{
		{
			name: "All green",
			args: args{
				invs: []Investment{
					{
						Investment: 10000,
						Date:       time.Date(2015, time.April, 15, 0, 0, 0, 0, time.UTC),
					},
				},
				currValue: Investment{
					Investment: 21589.25,
					Date:       time.Date(2025, time.April, 15, 0, 0, 0, 0, time.UTC),
				},
			},
			want: xirrData{
				Cashflow:        []float64{10000.0, -21589.25},
				DaysInPortfolio: []float64{3653, 0},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildXIRRData(tt.args.invs, tt.args.currValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildXIRRData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildXIRRData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestXIRR(t *testing.T) {
	type args struct {
		invs      []Investment
		currValue Investment
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "All green",
			args: args{
				invs: []Investment{
					{
						Investment: 10000,
						Date:       time.Date(2015, time.April, 15, 0, 0, 0, 0, time.UTC),
					},
				},
				currValue: Investment{
					Investment: 21589.25,
					Date:       time.Date(2025, time.April, 15, 0, 0, 0, 0, time.UTC),
				},
			},
			want:    0.0799317,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := XIRR(tt.args.invs, tt.args.currValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("XIRR() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if math.Abs(got-tt.want) > 0.000001 {
				t.Errorf("XIRR() = %v, want %v", got, tt.want)
			}
		})
	}
}
