package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestSetRoutesRegistersExpectedEndpoints(t *testing.T) {
	router := mux.NewRouter()
	SetRoutes(router)

	tests := []struct {
		method string
		path   string
	}{
		{http.MethodPost, "/api/auth/login"},
		{http.MethodPost, "/api/auth/register"},

		{http.MethodGet, "/api/accounts"},
		{http.MethodPost, "/api/accounts"},
		{http.MethodGet, "/api/transactions"},
		{http.MethodPost, "/api/transactions"},
		{http.MethodGet, "/api/bills"},
		{http.MethodGet, "/api/net-worth/history"},
		{http.MethodGet, "/api/reports/insights"},
		{http.MethodGet, "/api/assets"},
		{http.MethodPost, "/api/assets"},
		{http.MethodPut, "/api/assets/1"},
		{http.MethodDelete, "/api/assets/1"},
		{http.MethodGet, "/api/assets/1/valuations"},
		{http.MethodPost, "/api/assets/1/valuations"},
		{http.MethodGet, "/api/portfolio/holdings"},
		{http.MethodPost, "/api/portfolio/holdings"},
		{http.MethodPut, "/api/portfolio/holdings/1"},
		{http.MethodDelete, "/api/portfolio/holdings/1"},
		{http.MethodGet, "/api/portfolio/holdings/1/valuations"},
		{http.MethodPost, "/api/portfolio/holdings/1/valuations"},
		{http.MethodGet, "/api/transactions/1"},
		{http.MethodPut, "/api/transactions/1"},
		{http.MethodDelete, "/api/transactions/1"},

		{http.MethodGet, "/api/expenses"},
		{http.MethodPost, "/api/expenses"},
		{http.MethodGet, "/api/expenses/1"},
		{http.MethodPut, "/api/expenses/1"},
		{http.MethodDelete, "/api/expenses/1"},
		{http.MethodGet, "/api/expenses/cumulative"},

		{http.MethodGet, "/api/categories"},
		{http.MethodPost, "/api/categories"},
		{http.MethodGet, "/api/categories/1"},
		{http.MethodPut, "/api/categories/1"},
		{http.MethodDelete, "/api/categories/1"},

		{http.MethodGet, "/api/budgets"},
		{http.MethodPost, "/api/budgets"},
		{http.MethodGet, "/api/budgets/1"},
		{http.MethodPut, "/api/budgets/1"},
		{http.MethodDelete, "/api/budgets/1"},

		{http.MethodGet, "/api/goals"},
		{http.MethodPost, "/api/goals"},
		{http.MethodGet, "/api/goals/1"},
		{http.MethodPut, "/api/goals/1"},
		{http.MethodDelete, "/api/goals/1"},
		{http.MethodGet, "/api/planning/goals"},
		{http.MethodPost, "/api/planning/goals"},
		{http.MethodDelete, "/api/planning/goals/1"},
		{http.MethodGet, "/api/planning/goals/1/contributions"},
		{http.MethodPost, "/api/planning/goals/1/contributions"},
		{http.MethodGet, "/api/planning/envelopes"},
		{http.MethodPost, "/api/planning/envelopes"},
		{http.MethodPut, "/api/planning/envelopes/1"},
		{http.MethodDelete, "/api/planning/envelopes/1"},
		{http.MethodPost, "/api/scheduled-flows/1/occurrences"},
		{http.MethodGet, "/api/trends/spending"},

		{http.MethodGet, "/api/investments"},
		{http.MethodPost, "/api/investments"},
		{http.MethodGet, "/api/investments/1"},
		{http.MethodPut, "/api/investments/1"},
		{http.MethodDelete, "/api/investments/1"},
		{http.MethodGet, "/api/investments/cumulative"},

		{http.MethodGet, "/api/recurring/expenses"},
		{http.MethodPost, "/api/recurring/expenses"},
		{http.MethodGet, "/api/recurring/expenses/1"},
		{http.MethodPut, "/api/recurring/expenses/1"},
		{http.MethodDelete, "/api/recurring/expenses/1"},
		{http.MethodGet, "/api/recurring/investments"},
		{http.MethodPost, "/api/recurring/investments"},
		{http.MethodGet, "/api/recurring/investments/1"},
		{http.MethodPut, "/api/recurring/investments/1"},
		{http.MethodDelete, "/api/recurring/investments/1"},
	}

	for _, test := range tests {
		t.Run(test.method+" "+test.path, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.path, nil)
			var match mux.RouteMatch
			if !router.Match(request, &match) {
				t.Errorf("route %s %s is not registered", test.method, test.path)
			}
		})
	}
}

func TestSetRoutesRejectsUnknownAndWrongMethod(t *testing.T) {
	router := mux.NewRouter()
	SetRoutes(router)

	unknown := httptest.NewRequest(http.MethodGet, "/api/unknown", nil)
	var match mux.RouteMatch
	if router.Match(unknown, &match) {
		t.Error("unknown route should not match")
	}

	wrongMethod := httptest.NewRequest(http.MethodPatch, "/api/expenses", nil)
	if router.Match(wrongMethod, &match) {
		t.Error("PATCH /api/expenses should not match")
	}
}
