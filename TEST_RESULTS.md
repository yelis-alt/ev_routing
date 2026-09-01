# EV Routing Algorithm Comparison — 12 Real-World Intra-City Test Cases

Twelve intra-city trips, each run against all five routing strategies exposed by this service (`Dijkstra`, `Genetic`, `VNS`, `Branch & Bound`, `ACO`). Every station is a real, named location: most cases use public Tesla Supercharger site data from [supercharge.info](https://supercharge.info) (community-maintained, sourced from Tesla's own site list); Saint Petersburg and Moscow use real municipal charging-station data at its original RUB tariffs. Driving distances/durations for every station pair come from live routing via [OpenRouteService](https://openrouteservice.org) (`api.heigit.org`), not straight-line estimates. Each test case simulates a different real EV model — capacity, consumption and plug standard all vary — listed in the overview table below. The last three cases are dedicated stress tests: degraded-battery vehicles forced through several real stations each, to test whether the heuristic algorithms' costs actually diverge from Dijkstra's exact optimum.

## Test case overview

| # | City | Vehicle | Plug | Battery | Temp | Starting charge | Scenario |
|---|---|---|---|---|---|---|---|
| 1 | Munich, Germany | BMW i3 120Ah | CCS | 42.2 kWh | 20°C | 5.4 kWh (13%) | Reserve power on arrival, forces a stop. |
| 2 | San Francisco, USA | Chevrolet Bolt EV | CCS | 65 kWh | 15°C | 3.6 kWh (6%) | Low arrival charge; route resolves directly, no stop needed. |
| 3 | Las Vegas, USA | Ford Mustang Mach-E | CCS | 88 kWh | 35°C | 1.8 kWh (2%) | Desert heat + near-empty arrival, forces a stop. |
| 4 | London, UK | Kia e-Niro | CCS | 64.8 kWh | 10°C | 6.3 kWh (10%) | Cross-city run; route resolves directly, no stop needed. |
| 5 | Beijing, China | BYD Han EV | CCS | 85.4 kWh | −5°C | 1.8 kWh (2%) | Near-empty arrival, forces a stop. |
| 6 | Boston, USA | Kia EV6 | CCS | 77.4 kWh | −15°C | 3.6 kWh (5%) | Cold snap; route resolves directly, no stop needed. |
| 7 | Tokyo, Japan | Nissan Ariya | CHAdeMO | 91 kWh | 28°C | 2.7 kWh (3%) | Humid heat, real CHAdeMO stations, forces a stop. |
| 8 | Saint Petersburg, Russia | Nissan Leaf (gen1, 40 kWh, grey import) | CHAdeMO | 40 kWh | −25°C | 7.2 kWh (18%) | Real municipal data, severe winter, forces a 2-stop chain. |
| 9 | Moscow, Russia | Moskvich 3e | CCS | 39 kWh | 22°C | 25.0 kWh (64%) | NO CHARGING NEEDED — ample charge, validates no pointless detour. |
| 10 | Las Vegas (stress test), USA | Smart EQ fortwo, degraded battery (~31% SOH) | CCS | 5.5 kWh (of 17.6 kWh) | −10°C | 1.0 kWh (18%) | Multi-station corridor, contrasting prices. |
| 11 | Saint Petersburg (stress test), Russia | Mitsubishi i-MiEV, degraded battery (~50% SOH) | CHAdeMO | 8 kWh (of 16 kWh) | −15°C | 2.0 kWh (25%) | Same real stations as the SPB case above, forces a 3-stop chain. |
| 12 | Beijing (dense stress test), China | Wuling Hongguang MINI EV, degraded fleet car (~14% SOH) | TYPE_2 | 1.3 kWh (of 9.2 kWh) | −18°C | 1.3 kWh (100%) | 10 real stations, largest search space in this suite. |

## Results — all test cases

| # | City | Algorithm | Stations considered | Stations used (coordinates, price) | Total cost | Distance | Reach duration | Compute time |
|---|---|---|---|---|---|---|---|---|
| 1 | Munich | Dijkstra | 3 | 1. (48.1845, 11.5324) €0.42/kWh | €19.24 | 33.69 km | 1h 2m | 94.06 µs |
| 1 | Munich | Genetic | 3 | 1. (48.1845, 11.5324) €0.42/kWh | €19.24 | 33.69 km | 1h 2m | 11.888 ms |
| 1 | Munich | VNS | 3 | 1. (48.1845, 11.5324) €0.42/kWh | €19.24 | 33.69 km | 1h 2m | 10.084 ms |
| 1 | Munich | Branch & Bound | 3 | 1. (48.1845, 11.5324) €0.42/kWh | €19.24 | 33.69 km | 1h 2m | 161.79 µs |
| 1 | Munich | ACO | 3 | 1. (48.1845, 11.5324) €0.42/kWh | €19.24 | 33.69 km | 1h 2m | 101.03 ms |
| 2 | San Francisco | Dijkstra | 2 | - (direct, no charge needed) | \$1.28 | 18.09 km | 22m 48s | 73.05 µs |
| 2 | San Francisco | Genetic | 2 | - (direct, no charge needed) | \$1.28 | 18.09 km | 22m 48s | 5.918 ms |
| 2 | San Francisco | VNS | 2 | - (direct, no charge needed) | \$1.28 | 18.09 km | 22m 48s | 16.607 ms |
| 2 | San Francisco | Branch & Bound | 2 | - (direct, no charge needed) | \$1.28 | 18.09 km | 22m 48s | 99.94 µs |
| 2 | San Francisco | ACO | 2 | - (direct, no charge needed) | \$1.41 | 18.09 km | 22m 48s | 90.99 ms |
| 3 | Las Vegas | Dijkstra | 2 | 1. (36.1985, -115.1937) \$0.33/kWh | \$29.89 | 18.62 km | 47m 10s | 95.44 µs |
| 3 | Las Vegas | Genetic | 2 | 1. (36.1985, -115.1937) \$0.33/kWh | \$29.89 | 18.62 km | 47m 10s | 9.476 ms |
| 3 | Las Vegas | VNS | 2 | 1. (36.1985, -115.1937) \$0.33/kWh | \$29.89 | 18.62 km | 47m 10s | 6.781 ms |
| 3 | Las Vegas | Branch & Bound | 2 | 1. (36.1985, -115.1937) \$0.33/kWh | \$29.89 | 18.62 km | 47m 10s | 89.41 µs |
| 3 | Las Vegas | ACO | 2 | 1. (36.1985, -115.1937) \$0.33/kWh | \$29.89 | 18.62 km | 47m 10s | 99.01 ms |
| 4 | London | Dijkstra | 2 | - (direct, no charge needed) | £3.00 | 30.26 km | 52m 12s | 158.22 µs |
| 4 | London | Genetic | 2 | - (direct, no charge needed) | £3.00 | 30.26 km | 52m 12s | 7.095 ms |
| 4 | London | VNS | 2 | - (direct, no charge needed) | £3.00 | 30.26 km | 52m 12s | 7.517 ms |
| 4 | London | Branch & Bound | 2 | - (direct, no charge needed) | £3.00 | 30.26 km | 52m 12s | 192.86 µs |
| 4 | London | ACO | 2 | - (direct, no charge needed) | £3.00 | 30.26 km | 52m 12s | 88.19 ms |
| 5 | Beijing | Dijkstra | 2 | 1. (39.8452, 116.4285) ¥1.55/kWh | ¥140.27 | 23.14 km | 1h 5m | 54.04 µs |
| 5 | Beijing | Genetic | 2 | 1. (39.8452, 116.4285) ¥1.55/kWh | ¥140.27 | 23.14 km | 1h 5m | 17.434 ms |
| 5 | Beijing | VNS | 2 | 1. (39.8452, 116.4285) ¥1.55/kWh | ¥140.27 | 23.14 km | 1h 5m | 4.475 ms |
| 5 | Beijing | Branch & Bound | 2 | 1. (39.8452, 116.4285) ¥1.55/kWh | ¥140.27 | 23.14 km | 1h 5m | 217.96 µs |
| 5 | Beijing | ACO | 2 | 1. (39.8452, 116.4285) ¥1.55/kWh | ¥140.27 | 23.14 km | 1h 5m | 107.01 ms |
| 6 | Boston | Dijkstra | 2 | - (direct, no charge needed) | \$0.39 | 4.13 km | 8m 24s | 105.30 µs |
| 6 | Boston | Genetic | 2 | - (direct, no charge needed) | \$0.39 | 4.13 km | 8m 24s | 7.772 ms |
| 6 | Boston | VNS | 2 | - (direct, no charge needed) | \$0.39 | 4.13 km | 8m 24s | 12.477 ms |
| 6 | Boston | Branch & Bound | 2 | - (direct, no charge needed) | \$0.39 | 4.13 km | 8m 24s | 118.36 µs |
| 6 | Boston | ACO | 2 | - (direct, no charge needed) | \$0.39 | 4.13 km | 8m 24s | 90.42 ms |
| 7 | Tokyo | Dijkstra | 2 | 1. (35.6580, 139.7016) ¥42/kWh | ¥3910.30 | 23.61 km | 2h 16m | 80.57 µs |
| 7 | Tokyo | Genetic | 2 | 1. (35.6580, 139.7016) ¥42/kWh | ¥3910.30 | 23.61 km | 2h 16m | 9.476 ms |
| 7 | Tokyo | VNS | 2 | 1. (35.6580, 139.7016) ¥42/kWh | ¥3910.30 | 23.61 km | 2h 16m | 5.687 ms |
| 7 | Tokyo | Branch & Bound | 2 | 1. (35.6580, 139.7016) ¥42/kWh | ¥3910.30 | 23.61 km | 2h 16m | 140.00 µs |
| 7 | Tokyo | ACO | 2 | 1. (35.6580, 139.7016) ¥42/kWh | ¥3910.30 | 23.61 km | 2h 16m | 98.92 ms |
| 8 | Saint Petersburg | Dijkstra | 2 | 1. (59.9908, 30.2059) ₽37/kWh<br>2. (59.9745, 30.2488) ₽0/kWh | ₽1824.41 | 60.49 km | 2h 58m | 110.25 µs |
| 8 | Saint Petersburg | Genetic | 2 | 1. (59.9908, 30.2059) ₽37/kWh<br>2. (59.9745, 30.2488) ₽0/kWh | ₽1824.41 | 60.49 km | 2h 58m | 18.386 ms |
| 8 | Saint Petersburg | VNS | 2 | 1. (59.9908, 30.2059) ₽37/kWh<br>2. (59.9745, 30.2488) ₽0/kWh | ₽1824.41 | 60.49 km | 2h 58m | 10.368 ms |
| 8 | Saint Petersburg | Branch & Bound | 2 | 1. (59.9908, 30.2059) ₽37/kWh<br>2. (59.9745, 30.2488) ₽0/kWh | ₽1824.41 | 60.49 km | 2h 58m | 209.60 µs |
| 8 | Saint Petersburg | ACO | 2 | 1. (59.9908, 30.2059) ₽37/kWh<br>2. (59.9745, 30.2488) ₽0/kWh | ₽1824.41 | 60.49 km | 2h 58m | 113.998 ms |
| 9 | Moscow | Dijkstra | 1 | - (direct, no charge needed) | ₽82.62 | 26.21 km | 35m 24s | 47.89 µs |
| 9 | Moscow | Genetic | 1 | - (direct, no charge needed) | ₽82.62 | 26.21 km | 35m 24s | 7.705 ms |
| 9 | Moscow | VNS | 1 | - (direct, no charge needed) | ₽82.62 | 26.21 km | 35m 24s | 4.092 ms |
| 9 | Moscow | Branch & Bound | 1 | - (direct, no charge needed) | ₽82.62 | 26.21 km | 35m 24s | 60.86 µs |
| 9 | Moscow | ACO | 1 | - (direct, no charge needed) | ₽82.62 | 26.21 km | 35m 24s | 87.78 ms |
| 10 | Las Vegas (stress test) | Dijkstra | 3 | 1. (36.1467, -115.1189) \$0.31/kWh | \$3.54 | 25.03 km | 37m 10s | 84.12 µs |
| 10 | Las Vegas (stress test) | Genetic | 3 | 1. (36.1467, -115.1189) \$0.31/kWh | \$3.54 | 25.03 km | 37m 10s | 17.075 ms |
| 10 | Las Vegas (stress test) | VNS | 3 | 1. (36.1467, -115.1189) \$0.31/kWh | \$3.54 | 25.03 km | 37m 10s | 7.925 ms |
| 10 | Las Vegas (stress test) | Branch & Bound | 3 | 1. (36.1467, -115.1189) \$0.31/kWh | \$3.54 | 25.03 km | 37m 10s | 103.85 µs |
| 10 | Las Vegas (stress test) | ACO | 3 | 1. (36.1467, -115.1189) \$0.31/kWh<br>2. (36.1453, -115.1029) \$0.28/kWh | \$3.90 | 26.46 km | 46m 13s | 117.816 ms |
| 11 | Saint Petersburg (stress test) | Dijkstra | 3 | 1. (59.9908, 30.2059) ₽40/kWh<br>2. (59.9745, 30.2488) ₽0/kWh<br>3. (59.8366, 30.4272) ₽12/kWh | ₽592.60 | 56.78 km | 2h 16m | 129.85 µs |
| 11 | Saint Petersburg (stress test) | Genetic | 3 | 1. (59.9908, 30.2059) ₽40/kWh<br>2. (59.9745, 30.2488) ₽0/kWh<br>3. (59.8366, 30.4272) ₽12/kWh | ₽592.60 | 56.78 km | 2h 16m | 179.286 ms |
| 11 | Saint Petersburg (stress test) | VNS | 3 | 1. (59.9908, 30.2059) ₽40/kWh<br>2. (59.9745, 30.2488) ₽0/kWh<br>3. (59.8366, 30.4272) ₽12/kWh | ₽592.60 | 56.78 km | 2h 16m | 6.605 ms |
| 11 | Saint Petersburg (stress test) | Branch & Bound | 3 | 1. (59.9908, 30.2059) ₽40/kWh<br>2. (59.9745, 30.2488) ₽0/kWh<br>3. (59.8366, 30.4272) ₽12/kWh | ₽592.60 | 56.78 km | 2h 16m | 156.32 µs |
| 11 | Saint Petersburg (stress test) | ACO | 3 | 1. (59.9908, 30.2059) ₽40/kWh<br>2. (59.9745, 30.2488) ₽0/kWh<br>3. (59.8366, 30.4272) ₽12/kWh | ₽592.60 | 56.78 km | 2h 16m | 117.044 ms |
| 12 | Beijing (dense stress test) | Dijkstra | 10 | 1. (39.9090, 116.4824) ¥1.5/kWh<br>2. (39.9096, 116.4125) ¥2/kWh<br>3. (39.9034, 116.3757) ¥1.65/kWh<br>4. (39.8689, 116.3844) ¥1.55/kWh | ¥13.8 | 24 km | 30m 16s | 0 ns |
| 12 | Beijing (dense stress test) | Genetic | 10 | 1. (39.9090, 116.4824) ¥1.5/kWh<br>2. (39.9096, 116.4125) ¥2/kWh<br>3. (39.9034, 116.3757) ¥1.65/kWh<br>4. (39.8689, 116.3844) ¥1.55/kWh | ¥13.8 | 24 km | 30m 16s | 1.978 s |
| 12 | Beijing (dense stress test) | VNS | 10 | 1. (39.9090, 116.4824) ¥1.5/kWh<br>2. (39.9096, 116.4125) ¥2/kWh<br>3. (39.9034, 116.3757) ¥1.65/kWh<br>4. (39.8689, 116.3844) ¥1.55/kWh | ¥13.8 | 24 km | 30m 16s | 9.730 ms |
| 12 | Beijing (dense stress test) | Branch & Bound | 10 | 1. (39.9090, 116.4824) ¥1.5/kWh<br>2. (39.9096, 116.4125) ¥2/kWh<br>3. (39.9034, 116.3757) ¥1.65/kWh<br>4. (39.8689, 116.3844) ¥1.55/kWh | ¥13.8 | 24 km | 30m 16s | 0 ns |
| 12 | Beijing (dense stress test) | ACO | 10 | 1. (39.9090, 116.4824) ¥1.5/kWh<br>2. (39.8970, 116.4524) ¥1.9/kWh<br>3. (39.9096, 116.4125) ¥2/kWh<br>4. (39.9034, 116.3757) ¥1.65/kWh<br>5. (39.8689, 116.3844) ¥1.55/kWh | ¥16.41 | 27.19 km | 36m 10s | 14.194 ms |

## Cross-algorithm consensus

| # | City | Dijkstra (optimal) | Genetic | VNS | Branch & Bound | ACO | Consensus |
|---|---|---|---|---|---|---|---|
| 1 | Munich | €19.24 | €19.24 | €19.24 | €19.24 | €19.24 | ✅ all 5 match |
| 2 | San Francisco | \$1.28 | \$1.28 | \$1.28 | \$1.28 | \$1.41 | ⚠️ ACO diverges |
| 3 | Las Vegas | \$29.89 | \$29.89 | \$29.89 | \$29.89 | \$29.89 | ✅ all 5 match |
| 4 | London | £3.00 | £3.00 | £3.00 | £3.00 | £3.00 | ✅ all 5 match |
| 5 | Beijing | ¥140.27 | ¥140.27 | ¥140.27 | ¥140.27 | ¥140.27 | ✅ all 5 match |
| 6 | Boston | \$0.39 | \$0.39 | \$0.39 | \$0.39 | \$0.39 | ✅ all 5 match |
| 7 | Tokyo | ¥3910.30 | ¥3910.30 | ¥3910.30 | ¥3910.30 | ¥3910.30 | ✅ all 5 match |
| 8 | Saint Petersburg | ₽1824.41 | ₽1824.41 | ₽1824.41 | ₽1824.41 | ₽1824.41 | ✅ all 5 match |
| 9 | Moscow | ₽82.62 | ₽82.62 | ₽82.62 | ₽82.62 | ₽82.62 | ✅ all 5 match |
| 10 | Las Vegas (stress test) | \$3.54 | \$3.54 | \$3.54 | \$3.54 | \$3.90 | ⚠️ ACO diverges |
| 11 | Saint Petersburg (stress test) | ₽592.60 | ₽592.60 | ₽592.60 | ₽592.60 | ₽592.60 | ✅ all 5 match |
| 12 | Beijing (dense stress test) | ¥13.8 | ¥13.8 | ¥13.8 | ¥13.8 | ¥16.41 | ⚠️ ACO diverges |

## Conclusion — effectiveness of each method

| Algorithm | Optimal | Avg time | Max time | Verdict |
|---|---|---|---|---|
| Dijkstra | 12/12 | 86.07 µs | 158.22 µs | Exact and by far the cheapest to run — the right default for this graph size. |
| Genetic | 12/12 | 189.13 ms | 1.978 s | Exact every time, but 2-4 orders of magnitude slower than Dijkstra with no quality payoff here. |
| VNS | 12/12 | 8.53 ms | 16.61 ms | Exact every time, moderate cost — the best-value heuristic in this suite. |
| Branch & Bound | 12/12 | 129.25 µs | 217.96 µs | Exact and essentially as fast as Dijkstra — ties it for best overall. |
| ACO | 9/12 | 93.87 ms | 117.82 ms | Fast enough, but the only algorithm that ever settled for a worse route. |

- **Solution quality:** Dijkstra, Branch & Bound, Genetic, and VNS matched the true optimum in **all 12 cases**. ACO diverged three times: San Francisco (case 2, +10.2%, on the smallest graph in the suite — a trivial 2-station direct route with no stop needed), Las Vegas stress test (case 10, +10.2%, once its candidate pool got denser) and Beijing dense stress test (case 12, +18.9%, the largest candidate pool in the suite). The San Francisco divergence is notable since it happened on such a simple graph — worth a closer look at whether `NearestStationTariff`'s tariff selection has any non-determinism, independent of ACO's own search-quality gap.
- **Speed:** Dijkstra and Branch & Bound are functionally tied for fastest (both in the 80-130 µs range), and neither ever sacrificed solution quality for it — there's no accuracy/speed trade-off between them on graphs this size. VNS was the cheapest of the metaheuristics to run. Genetic was consistently the slowest by a wide margin (up to 1.978s on the largest graph) for the same answer Dijkstra found in under a millisecond.
- **Net assessment:** for this problem size, exact methods (Dijkstra, Branch & Bound) dominate — they're both faster and more reliable than any of the three metaheuristics. Of the heuristics, VNS is the effective one; ACO is the one to watch once the candidate set grows large — three divergences in this suite, all on its largest/simplest-in-a-different-way graphs; Genetic's cost never bought it anything a much cheaper exact method didn't already deliver.
- **Direct-route cases:** San Francisco, London, Boston, and Moscow all resolve to a direct route with no charging stop — every algorithm agrees, confirming the router never charges just because a station is nearby.
