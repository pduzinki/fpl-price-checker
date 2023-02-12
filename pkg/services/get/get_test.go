package get

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/pduzinki/fpl-price-checker/pkg/domain"
)

var (
	errReportGetterFailure = fmt.Errorf("report getter failure")

	records = []domain.Record{
		{
			Name:        "Mitoma",
			OldPrice:    "-",
			NewPrice:    "5.0",
			Description: "new",
		},
		{
			Name:        "Kane",
			OldPrice:    "11.9",
			NewPrice:    "12.0",
			Description: "rise",
		},
	}

	report = domain.PriceChangeReport{
		Date:    "2012-12-12",
		Records: records,
	}

	ReportGetterOk = PriceChangeReportGetterMock{
		GetByDateFunc: func(ctx context.Context, date string) (domain.PriceChangeReport, error) {
			return report, nil
		},
	}

	ReportGetterFailure = PriceChangeReportGetterMock{
		GetByDateFunc: func(ctx context.Context, date string) (domain.PriceChangeReport, error) {
			return domain.PriceChangeReport{}, errReportGetterFailure
		},
	}
)

func TestGetLatestReport(t *testing.T) {
	testcases := []struct {
		name    string
		rg      PriceChangeReportGetter
		want    domain.PriceChangeReport
		wantErr error
	}{
		{
			name:    "sunny scenario",
			rg:      &ReportGetterOk,
			want:    report,
			wantErr: nil,
		},
		{
			name:    "ReportGetter failure",
			rg:      &ReportGetterFailure,
			want:    domain.PriceChangeReport{},
			wantErr: errReportGetterFailure,
		},
	}

	for _, test := range testcases {
		test := test

		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			g := NewGetService(test.rg)

			got, gotErr := g.GetLatestReport(ctx)
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("want: %v, got: %v", test.want, got)
			}

			if !errors.Is(gotErr, test.wantErr) {
				t.Errorf("want: %v, got: %v", test.wantErr, gotErr)
			}
		})
	}
}

func TestGetReportByDate(t *testing.T) {
	testcases := []struct {
		name    string
		rg      PriceChangeReportGetter
		date    string
		want    domain.PriceChangeReport
		wantErr error
	}{
		{
			name:    "sunny scenario",
			rg:      &ReportGetterOk,
			date:    "2012-12-12",
			want:    report,
			wantErr: nil,
		},
		{
			name:    "ReportGetter failure",
			rg:      &ReportGetterFailure,
			date:    "2012-12-12",
			want:    domain.PriceChangeReport{},
			wantErr: errReportGetterFailure,
		},
		// TODO fix testcase below
		// {
		// 	name: "Wrong date format",
		// 	rg:   &ReportGetterOk,
		// 	date: "not-even-a-date",
		// 	want: domain.PriceChangeReport{},
		// 	wantErr: &time.ParseError{
		// 		Layout:     domain.DateFormat,
		// 		Value:      "not-even-a-date",
		// 		LayoutElem: "2006",
		// 		ValueElem:  "not-even-a-date",
		// 		Message:    "",
		// 	},
		// },
	}

	for _, test := range testcases {
		test := test

		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			g := NewGetService(test.rg)

			got, gotErr := g.GetReportByDate(ctx, test.date)
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("want: %v, got: %v", test.want, got)
			}

			if !errors.Is(gotErr, test.wantErr) {
				t.Errorf("want: %v, got: %v", test.wantErr, gotErr)
			}
		})
	}
}
