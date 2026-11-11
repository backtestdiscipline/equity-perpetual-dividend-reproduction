package study

import (
	"strings"
	"testing"
)

func TestCalculate(t *testing.T) {
	data := map[string][]byte{}
	for _, e := range events {
		data["data/"+strings.ToLower(e.Symbol)+".json"] = []byte(`[]`)
	}
	r, err := calculate(data)
	if err != nil || r.Status != "INCOMPLETE" || len(r.Issues) != 5 {
		t.Fatalf("missing: %+v %v", r, err)
	}
	data["data/aapl.json"] = []byte(`[{"symbol":"AAPLUSDT","fundingTime":1778457600005,"fundingRate":"-0.00270712","markPrice":"100","rateType":"Special"}]`)
	r, err = calculate(data)
	if err != nil || len(r.Events) != 1 || !r.Events[0].Match {
		t.Fatalf("event: %+v %v", r, err)
	}
	data["data/aapl.json"] = []byte(`[{"symbol":"AAPLUSDT","fundingTime":1778457600005,"fundingRate":"-0.1","markPrice":"100","rateType":"Special"},{"symbol":"AAPLUSDT","fundingTime":1778457600005,"fundingRate":"-0.1","markPrice":"100","rateType":"Special"}]`)
	r, err = calculate(data)
	if err != nil || len(r.Events) != 0 {
		t.Fatalf("duplicate: %+v %v", r, err)
	}
}
