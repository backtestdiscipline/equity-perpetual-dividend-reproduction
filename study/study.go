// Package study owns event interpretation and semantic verification.
package study

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"path"
	"strings"
	"time"

	"research-platform/seal"
)

type Evaluation struct {
	Outcome seal.TerminalOutcome
	Payload []byte
}
type Verification struct {
	Status string `json:"status"`
	Notes  string `json:"notes"`
}
type event struct{ Symbol, Date, Dividend string }

var events = []event{{"AAPL", "2026-05-11", "0.27"}, {"MSFT", "2026-05-21", "0.91"}, {"NVDA", "2026-06-04", "0.25"}, {"GOOGL", "2026-06-08", "0.22"}, {"META", "2026-06-15", "0.525"}}

type row struct {
	Symbol      string `json:"symbol"`
	FundingTime int64  `json:"fundingTime"`
	FundingRate string `json:"fundingRate"`
	MarkPrice   string `json:"markPrice"`
	RateType    string `json:"rateType"`
}
type resultEvent struct {
	Symbol                 string `json:"symbol"`
	ExDate                 string `json:"ex_date"`
	Dividend               string `json:"dividend"`
	FundingTime            int64  `json:"funding_time"`
	Rate                   string `json:"rate"`
	Mark                   string `json:"mark"`
	PaperResidual          string `json:"paper_residual"`
	CurrentFormulaResidual string `json:"current_formula_residual"`
	Match                  bool   `json:"match"`
}
type result struct {
	Status string        `json:"status"`
	Events []resultEvent `json:"events"`
	Issues []string      `json:"issues"`
}

func decimal(s string) (*big.Rat, bool) {
	if s == "" || strings.ContainsAny(s, "eE/ ") {
		return nil, false
	}
	r, ok := new(big.Rat).SetString(s)
	return r, ok
}
func calculate(data map[string][]byte) (result, error) {
	out := result{Status: "MATCH", Events: []resultEvent{}, Issues: []string{}}
	for _, e := range events {
		b, ok := data["data/"+strings.ToLower(e.Symbol)+".json"]
		if !ok {
			out.Issues = append(out.Issues, e.Symbol+": missing source")
			continue
		}
		var rows []row
		if err := json.Unmarshal(b, &rows); err != nil {
			return out, fmt.Errorf("%s: invalid JSON: %w", e.Symbol, err)
		}
		t, _ := time.Parse("2006-01-02", e.Date)
		center := t.UnixMilli()
		var matches []row
		for _, r := range rows {
			if r.Symbol == e.Symbol+"USDT" && r.RateType == "Special" && r.FundingTime >= center-60000 && r.FundingTime <= center+60000 {
				matches = append(matches, r)
			}
		}
		if len(matches) != 1 {
			out.Issues = append(out.Issues, fmt.Sprintf("%s: expected one special entry, got %d", e.Symbol, len(matches)))
			continue
		}
		r := matches[0]
		d, _ := decimal(e.Dividend)
		m, mok := decimal(r.MarkPrice)
		rate, rok := decimal(r.FundingRate)
		if !mok || !rok || m.Cmp(d) <= 0 {
			out.Issues = append(out.Issues, e.Symbol+": invalid rate or mark")
			continue
		}
		paper := new(big.Rat).Neg(new(big.Rat).Quo(d, new(big.Rat).Sub(m, d)))
		current := new(big.Rat).Neg(new(big.Rat).Quo(d, m))
		diff := new(big.Rat).Abs(new(big.Rat).Sub(rate, paper))
		alt := new(big.Rat).Abs(new(big.Rat).Sub(rate, current))
		tol, _ := decimal("0.000005")
		pass := diff.Cmp(tol) < 0
		out.Events = append(out.Events, resultEvent{e.Symbol, e.Date, e.Dividend, r.FundingTime, r.FundingRate, r.MarkPrice, diff.FloatString(12), alt.FloatString(12), pass})
		if !pass {
			out.Status = "DOES_NOT_MATCH"
		}
	}
	if len(out.Issues) > 0 {
		out.Status = "INCOMPLETE"
	}
	return out, nil
}
func Evaluate(in Inputs) (Evaluation, error) {
	data := map[string][]byte{}
	for _, ref := range in.Sources {
		r, err := in.OpenSource(ref)
		if err != nil {
			return Evaluation{}, err
		}
		b, err := io.ReadAll(io.LimitReader(r, 1000001))
		closeErr := r.Close()
		if err != nil {
			return Evaluation{}, err
		}
		if closeErr != nil {
			return Evaluation{}, closeErr
		}
		if len(b) > 1000000 {
			return Evaluation{}, fmt.Errorf("source too large: %s", ref.Path)
		}
		data[path.Clean(ref.Path)] = b
	}
	out, err := calculate(data)
	if err != nil {
		b, _ := json.Marshal(map[string]string{"error": err.Error()})
		return Evaluation{seal.TechnicalFailure, b}, nil
	}
	b, err := json.Marshal(out)
	if err != nil {
		return Evaluation{}, err
	}
	outcome := seal.Complete
	if out.Status == "INCOMPLETE" {
		outcome = seal.Incomplete
	}
	return Evaluation{outcome, append(b, '\n')}, nil
}
func Verify(in Inputs, payload []byte) (Verification, error) {
	e, err := Evaluate(in)
	if err != nil {
		return Verification{}, err
	}
	if e.Outcome == seal.TechnicalFailure || !bytes.Equal(e.Payload, payload) {
		return Verification{Status: "FAIL", Notes: "Semantic replay differs or input is invalid."}, nil
	}
	return Verification{Status: "PASS", Notes: "Exact deterministic event calculations replay; this is not independent dividend or venue-rule verification."}, nil
}
