// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package get

import (
	"context"
	"github.com/pduzinki/fpl-price-checker/internal/domain"
	"sync"
)

// Ensure, that PriceChangeReportGetterMock does implement PriceChangeReportGetter.
// If this is not the case, regenerate this file with moq.
var _ PriceChangeReportGetter = &PriceChangeReportGetterMock{}

// PriceChangeReportGetterMock is a mock implementation of PriceChangeReportGetter.
//
//	func TestSomethingThatUsesPriceChangeReportGetter(t *testing.T) {
//
//		// make and configure a mocked PriceChangeReportGetter
//		mockedPriceChangeReportGetter := &PriceChangeReportGetterMock{
//			GetByDateFunc: func(ctx context.Context, date string) (domain.PriceChangeReport, error) {
//				panic("mock out the GetByDate method")
//			},
//		}
//
//		// use mockedPriceChangeReportGetter in code that requires PriceChangeReportGetter
//		// and then make assertions.
//
//	}
type PriceChangeReportGetterMock struct {
	// GetByDateFunc mocks the GetByDate method.
	GetByDateFunc func(ctx context.Context, date string) (domain.PriceChangeReport, error)

	// calls tracks calls to the methods.
	calls struct {
		// GetByDate holds details about calls to the GetByDate method.
		GetByDate []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Date is the date argument value.
			Date string
		}
	}
	lockGetByDate sync.RWMutex
}

// GetByDate calls GetByDateFunc.
func (mock *PriceChangeReportGetterMock) GetByDate(ctx context.Context, date string) (domain.PriceChangeReport, error) {
	if mock.GetByDateFunc == nil {
		panic("PriceChangeReportGetterMock.GetByDateFunc: method is nil but PriceChangeReportGetter.GetByDate was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		Date string
	}{
		Ctx:  ctx,
		Date: date,
	}
	mock.lockGetByDate.Lock()
	mock.calls.GetByDate = append(mock.calls.GetByDate, callInfo)
	mock.lockGetByDate.Unlock()
	return mock.GetByDateFunc(ctx, date)
}

// GetByDateCalls gets all the calls that were made to GetByDate.
// Check the length with:
//
//	len(mockedPriceChangeReportGetter.GetByDateCalls())
func (mock *PriceChangeReportGetterMock) GetByDateCalls() []struct {
	Ctx  context.Context
	Date string
} {
	var calls []struct {
		Ctx  context.Context
		Date string
	}
	mock.lockGetByDate.RLock()
	calls = mock.calls.GetByDate
	mock.lockGetByDate.RUnlock()
	return calls
}