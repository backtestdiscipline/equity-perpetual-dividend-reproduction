# Reproducing five dividend adjustments in Binance equity perpetuals

Five Binance equity-perpetual special funding rates in captured exchange records numerically match the formula reported by Seo et al. in [“Basis-Only Funding Is Not Dividend Fair”](https://arxiv.org/abs/2610.22839), version 1. For each event the absolute difference is below the paper's 0.000005 funding-rate threshold. This is a narrow retrospective accounting reproduction; it does not independently establish the exchange's historical calculation method.

The paper's results and the same API events were inspected before our method was frozen. This is therefore an exploratory reproduction of known results, not a prospective confirmation. The study uses the paper's five dividend amounts and ex-dates as inputs; later documentary checks are described below and are not part of the sealed run.

## Method

For each of the paper's five Table 1 events, the study reads a captured Binance USDⓈ-M funding-history response and selects the unique `Special` entry for the named contract within 60 seconds of 00:00 UTC on the listed ex-date. All five events occurred during US daylight saving time: 00:00 UTC is 20:00 Eastern on the previous calendar day. This convention is specific to these events, not a general rule for all dates.

Let D be the dividend per share from the paper, M the `markPrice` field in the selected API entry, and r its `fundingRate`. Using exact rational arithmetic, the evaluator computes the absolute residual |r − (−D/(M−D))|. Each event matches if that residual is strictly less than 0.000005; all five are required for the overall match. The selection rule, inputs and tolerance are in the frozen [specification](records/equity-perpetual-dividend-reproduction/spec.json). A missing or ambiguous event is incomplete, rather than a zero payment. Displayed residuals below are rounded to 12 decimal places; the threshold comparison uses the unrounded rational values.

## Results

| Contract | Ex-date | Dividend ($) | Special rate | API mark ($) | Absolute residual |
|---|---|---:|---:|---:|---:|
| AAPLUSDT | 2026-05-11 | 0.27 | -0.00092193 | 293.13371232 | 0.000000000538 |
| MSFTUSDT | 2026-05-21 | 0.91 | -0.00217887 | 418.55739734 | 0.000000001473 |
| NVDAUSDT | 2026-06-04 | 0.25 | -0.00116812 | 214.27000000 | 0.000000004871 |
| GOOGLUSDT | 2026-06-08 | 0.22 | -0.00060111 | 366.21000000 | 0.000000000680 |
| METAUSDT | 2026-06-15 | 0.525 | -0.00091155 | 576.46461396 | 0.000000003898 |

The table is rendered from the [sealed result](.research-platform/runs/20261111T125325Z-0a66a12c1a04/artifacts/result.json) by `python3 scripts/result-table.py`. The run outcome was COMPLETE and its comparison status was MATCH. The standard funding record is separate from the special record and is not included in this formula comparison. The captures are venue-provided market records; they are not individual account cash-flow statements.

## Interpretation and limitations

Binance's [current methodology](https://www.binance.com/en/support/faq/detail/7ced719b5e9a4859a1864c2fe657309f), viewed on 2026-11-11 and shown as updated on 2026-10-01, specifies −D/M with M the real-time mark at execution. The historical API records match −D/(API mark−D) more closely. We have not established the historical terms or whether those mark definitions are equivalent, so this does not demonstrate a rule change or identify Binance's actual historical calculation. The sealed results also report the residual from −D/M as a diagnostic; it was not an acceptance test.

Post-seal documentary checks support the five dividend amounts and provide independent ex-date evidence, with qualifications for Apple and Alphabet. Those sources were not included in the frozen lineage or its semantic replay; see [FOLLOWUP.md](FOLLOWUP.md) for source-by-source evidence and access limitations. The paper's dividend inputs were not independently verified by the sealed calculation itself.

Platform integrity verification passed, meaning the run's recorded artifacts and bindings checked out against the retained bytes. Study semantic verification also passed, but it reuses the evaluator; it is deterministic replay, not an independent implementation or a proof of the financial interpretation. The finding does not reproduce the paper's mathematical proof, establish cross-venue fairness, verify actual account payments, or demonstrate an executable or profitable trade.

The [handoff](handoffs/20261111T125325Z-0a66a12c1a04.json) and [run artifacts](.research-platform/runs/20261111T125325Z-0a66a12c1a04/) preserve the complete result. Source files are excluded from distribution; [README.md](README.md) explains the resulting limits on exact independent verification.
