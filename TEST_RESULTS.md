# EV Routing Algorithm Comparison — 12 Real-World Intra-City Test Cases

Twelve intra-city trips, each run against all five routing strategies exposed by this service (`Dijkstra`, `Genetic`, `VNS`, `Branch & Bound`, `ACO`). Every station is a real, named location: nine cases use public Tesla Supercharger site data from [supercharge.info](https://supercharge.info) (community-maintained, sourced from Tesla's own site list); Saint Petersburg and Moscow use the real municipal charging-station data the user supplied from their own `electrocar.station` database, at its original RUB tariffs. Driving distances/durations for every station pair come from live routing via [OpenRouteService](https://openrouteservice.org) (`api.heigit.org`), not straight-line estimates. Each test case simulates a different real EV model — capacity, consumption and plug standard all vary — listed in the overview table below. The last three cases are dedicated stress tests: degraded-battery vehicles forced through 6-10 real stations each, to test whether the heuristic algorithms' costs actually diverge from Dijkstra's exact optimum.

## Test case overview

| # | City | Vehicle | Plug | Battery | Temp | Starting charge | Scenario |
|---|---|---|---|---|---|---|---|
| 1 | Munich, Germany | BMW i3 120Ah | CCS | 42.2 kWh | 20°C | 5.4 kWh (13%) | Reserve power on arrival, forces a stop. |
| 2 | San Francisco, USA | Chevrolet Bolt EV | CCS | 65 kWh | 15°C | 3.6 kWh (6%) | Low arrival charge, 5 stations to pick from. |
| 3 | Las Vegas, USA | Ford Mustang Mach-E | CCS | 88 kWh | 35°C | 1.8 kWh (2%) | Desert heat + near-empty arrival, 6 stations. |
| 4 | London, UK | Kia e-Niro | CCS | 64.8 kWh | 10°C | 6.3 kWh (10%) | Long cross-city run, 5 stations en route. |
| 5 | Beijing, China | BYD Han EV | CCS | 85.4 kWh | −5°C | 1.8 kWh (2%) | Dense downtown cluster, near-empty arrival. |
| 6 | Boston, USA | Kia EV6 | CCS | 77.4 kWh | −15°C | 3.6 kWh (5%) | Cold snap, 4 stations cross-town. |
| 7 | Tokyo, Japan | Nissan Ariya | CHAdeMO | 91 kWh | 28°C | 2.7 kWh (3%) | Humid heat, 6 real CHAdeMO stations. |
| 8 | Saint Petersburg, Russia | Nissan Leaf (gen1, 40 kWh, grey import) | CHAdeMO | 40 kWh | −25°C | 7.2 kWh (18%) | Real municipal data, severe winter, 1 free station. |
| 9 | Moscow, Russia | Moskvich 3e | CCS | 39 kWh | 22°C | 25.0 kWh (64%) | NO CHARGING NEEDED — ample charge, validates no pointless detour. |
| 10 | Las Vegas (stress test), USA | Smart EQ fortwo, degraded battery (~31% SOH) | CCS | 5.5 kWh (of 17.6 kWh) | −10°C | 1.0 kWh (18%) | Multi-stop chain, 8 stations along the Fremont St.→Henderson corridor, contrasting prices. |
| 11 | Saint Petersburg (stress test), Russia | Mitsubishi i-MiEV, degraded battery (~50% SOH) | CHAdeMO | 8 kWh (of 16 kWh) | −15°C | 2.0 kWh (25%) | Same real 6 stations as the SPB case above, longer forced chain. |
| 12 | Beijing (dense stress test), China | Wuling Hongguang MINI EV, degraded fleet car (~14% SOH) | TYPE_2 | 1.3 kWh (of 9.2 kWh) | −18°C | 1.3 kWh (100%) | 10 real stations, largest search space in this suite. |

## Results — all test cases

| # | City | Algorithm | Stations considered | Stations used (coordinates, price) | Total cost | Distance | Reach duration | Compute time |
|---|---|---|---|---|---|---|---|---|
| 1 | Munich | Dijkstra | 3 | 1. (48.1845, 11.5324) €0.42/kWh | €19.04 | 33.69 km | 1h 34m | 0 ns |
| 1 | Munich | Genetic | 3 | 1. (48.1845, 11.5324) €0.42/kWh | €19.04 | 33.69 km | 1h 34m | 22.150 ms |
| 1 | Munich | VNS | 3 | 1. (48.1845, 11.5324) €0.42/kWh | €19.04 | 33.69 km | 1h 34m | 11.259 ms |
| 1 | Munich | Branch & Bound | 3 | 1. (48.1845, 11.5324) €0.42/kWh | €19.04 | 33.69 km | 1h 34m | 0 ns |
| 1 | Munich | ACO | 3 | 1. (48.1845, 11.5324) €0.42/kWh | €19.04 | 33.69 km | 1h 34m | 19.940 ms |
| 2 | San Francisco | Dijkstra | 5 | 1. (37.7387, -122.4062) \$0.42/kWh | \$28.46 | 17.44 km | 37m 10s | 0 ns |
| 2 | San Francisco | Genetic | 5 | 1. (37.7387, -122.4062) \$0.42/kWh | \$28.46 | 17.44 km | 37m 10s | 36.169 ms |
| 2 | San Francisco | VNS | 5 | 1. (37.7387, -122.4062) \$0.42/kWh | \$28.46 | 17.44 km | 37m 10s | 5.710 ms |
| 2 | San Francisco | Branch & Bound | 5 | 1. (37.7387, -122.4062) \$0.42/kWh | \$28.46 | 17.44 km | 37m 10s | 0 ns |
| 2 | San Francisco | ACO | 5 | 1. (37.7387, -122.4062) \$0.42/kWh | \$28.46 | 17.44 km | 37m 10s | 7.916 ms |
| 3 | Las Vegas | Dijkstra | 6 | 1. (36.1985, -115.1937) \$0.33/kWh | \$30.72 | 24.29 km | 48m 37s | 507.90 µs |
| 3 | Las Vegas | Genetic | 6 | 1. (36.1985, -115.1937) \$0.33/kWh | \$30.72 | 24.29 km | 48m 37s | 30.916 ms |
| 3 | Las Vegas | VNS | 6 | 1. (36.1985, -115.1937) \$0.33/kWh | \$30.72 | 24.29 km | 48m 37s | 28.103 ms |
| 3 | Las Vegas | Branch & Bound | 6 | 1. (36.1985, -115.1937) \$0.33/kWh | \$30.72 | 24.29 km | 48m 37s | 506.50 µs |
| 3 | Las Vegas | ACO | 6 | 1. (36.1985, -115.1937) \$0.33/kWh | \$30.72 | 24.29 km | 48m 37s | 31.386 ms |
| 4 | London | Dijkstra | 5 | 1. (51.5063, -0.0728) £0.65/kWh | £46.42 | 47.28 km | 2h 0m | 0 ns |
| 4 | London | Genetic | 5 | 1. (51.5063, -0.0728) £0.65/kWh | £46.42 | 47.28 km | 2h 0m | 33.624 ms |
| 4 | London | VNS | 5 | 1. (51.5063, -0.0728) £0.65/kWh | £46.42 | 47.28 km | 2h 0m | 16.269 ms |
| 4 | London | Branch & Bound | 5 | 1. (51.5063, -0.0728) £0.65/kWh | £46.42 | 47.28 km | 2h 0m | 0 ns |
| 4 | London | ACO | 5 | 1. (51.5063, -0.0728) £0.65/kWh | £46.42 | 47.28 km | 2h 0m | 26.231 ms |
| 5 | Beijing | Dijkstra | 6 | 1. (39.8611, 116.4660) ¥1.6/kWh | ¥145.07 | 21.08 km | 42m 39s | 0 ns |
| 5 | Beijing | Genetic | 6 | 1. (39.8611, 116.4660) ¥1.6/kWh | ¥145.07 | 21.08 km | 42m 39s | 32.358 ms |
| 5 | Beijing | VNS | 6 | 1. (39.8611, 116.4660) ¥1.6/kWh | ¥145.07 | 21.08 km | 42m 39s | 19.506 ms |
| 5 | Beijing | Branch & Bound | 6 | 1. (39.8611, 116.4660) ¥1.6/kWh | ¥145.07 | 21.08 km | 42m 39s | 0 ns |
| 5 | Beijing | ACO | 6 | 1. (39.8611, 116.4660) ¥1.6/kWh | ¥145.07 | 21.08 km | 42m 39s | 31.281 ms |
| 6 | Boston | Dijkstra | 4 | 1. (42.3472, -71.0810) \$0.33/kWh | \$26.73 | 9.52 km | 1h 21m | 0 ns |
| 6 | Boston | Genetic | 4 | 1. (42.3472, -71.0810) \$0.33/kWh | \$26.73 | 9.52 km | 1h 21m | 71.159 ms |
| 6 | Boston | VNS | 4 | 1. (42.3472, -71.0810) \$0.33/kWh | \$26.73 | 9.52 km | 1h 21m | 2.572 ms |
| 6 | Boston | Branch & Bound | 4 | 1. (42.3472, -71.0810) \$0.33/kWh | \$26.73 | 9.52 km | 1h 21m | 0 ns |
| 6 | Boston | ACO | 4 | 1. (42.3472, -71.0810) \$0.33/kWh | \$26.73 | 9.52 km | 1h 21m | 7.777 ms |
| 7 | Tokyo | Dijkstra | 6 | 1. (35.6491, 139.6996) ¥44/kWh | ¥4142.46 | 16.02 km | 45m 4s | 0 ns |
| 7 | Tokyo | Genetic | 6 | 1. (35.6491, 139.6996) ¥44/kWh | ¥4142.46 | 16.02 km | 45m 4s | 10.792 ms |
| 7 | Tokyo | VNS | 6 | 1. (35.6491, 139.6996) ¥44/kWh | ¥4142.46 | 16.02 km | 45m 4s | 7.201 ms |
| 7 | Tokyo | Branch & Bound | 6 | 1. (35.6491, 139.6996) ¥44/kWh | ¥4142.46 | 16.02 km | 45m 4s | 0 ns |
| 7 | Tokyo | ACO | 6 | 1. (35.6491, 139.6996) ¥44/kWh | ¥4142.46 | 16.02 km | 45m 4s | 8.946 ms |
| 8 | Saint Petersburg | Dijkstra | 6 | 1. (59.9908, 30.2059) ₽37/kWh<br>2. (59.9745, 30.2488) ₽0/kWh | ₽1846.86 | 69.17 km | 3h 28m | 0 ns |
| 8 | Saint Petersburg | Genetic | 6 | 1. (59.9908, 30.2059) ₽37/kWh<br>2. (59.9745, 30.2488) ₽0/kWh | ₽1846.86 | 69.17 km | 3h 28m | 39.938 ms |
| 8 | Saint Petersburg | VNS | 6 | 1. (59.9908, 30.2059) ₽37/kWh<br>2. (59.9745, 30.2488) ₽0/kWh | ₽1846.86 | 69.17 km | 3h 28m | 7.439 ms |
| 8 | Saint Petersburg | Branch & Bound | 6 | 1. (59.9908, 30.2059) ₽37/kWh<br>2. (59.9745, 30.2488) ₽0/kWh | ₽1846.86 | 69.17 km | 3h 28m | 0 ns |
| 8 | Saint Petersburg | ACO | 6 | 1. (59.9908, 30.2059) ₽37/kWh<br>2. (59.9745, 30.2488) ₽0/kWh | ₽1846.86 | 69.17 km | 3h 28m | 8.223 ms |
| 9 | Moscow | Dijkstra | 1 | - | ₽133.62 | 28.11 km | 34m 11s | 0 ns |
| 9 | Moscow | Genetic | 1 | - | ₽133.62 | 28.11 km | 34m 11s | 7.431 ms |
| 9 | Moscow | VNS | 1 | - | ₽133.62 | 28.11 km | 34m 11s | 1.785 ms |
| 9 | Moscow | Branch & Bound | 1 | - | ₽133.62 | 28.11 km | 34m 11s | 0 ns |
| 9 | Moscow | ACO | 1 | - | ₽133.62 | 28.11 km | 34m 11s | 6.275 ms |
| 10 | Las Vegas (stress test) | Dijkstra | 8 | 1. (36.1659, -115.1387) \$0.4/kWh<br>2. (36.1467, -115.1189) \$0.31/kWh | \$4.26 | 26.03 km | 36m 55s | 175.37 µs |
| 10 | Las Vegas (stress test) | Genetic | 8 | 1. (36.1659, -115.1387) \$0.4/kWh<br>2. (36.1467, -115.1189) \$0.31/kWh | \$4.26 | 26.03 km | 36m 55s | 19.924 ms |
| 10 | Las Vegas (stress test) | VNS | 8 | 1. (36.1659, -115.1387) \$0.4/kWh<br>2. (36.1467, -115.1189) \$0.31/kWh | \$4.26 | 26.03 km | 36m 55s | 19.952 ms |
| 10 | Las Vegas (stress test) | Branch & Bound | 8 | 1. (36.1659, -115.1387) \$0.4/kWh<br>2. (36.1467, -115.1189) \$0.31/kWh | \$4.26 | 26.03 km | 36m 55s | 115.74 µs |
| 10 | Las Vegas (stress test) | ACO | 8 | 1. (36.1659, -115.1387) \$0.4/kWh<br>2. (36.1453, -115.1029) \$0.28/kWh<br>3. (36.1467, -115.1189) \$0.31/kWh | \$4.96 | 29.63 km | 45m 6s | 30.752 ms |
| 11 | Saint Petersburg (stress test) | Dijkstra | 6 | 1. (59.9908, 30.2059) ₽40/kWh<br>2. (59.9745, 30.2488) ₽0/kWh<br>3. (59.8366, 30.4272) ₽12/kWh | ₽599.86 | 56.56 km | 2h 37m | 0 ns |
| 11 | Saint Petersburg (stress test) | Genetic | 6 | 1. (59.9908, 30.2059) ₽40/kWh<br>2. (59.9745, 30.2488) ₽0/kWh<br>3. (59.8366, 30.4272) ₽12/kWh | ₽599.86 | 56.56 km | 2h 37m | 511.316 ms |
| 11 | Saint Petersburg (stress test) | VNS | 6 | 1. (59.9908, 30.2059) ₽40/kWh<br>2. (59.9745, 30.2488) ₽0/kWh<br>3. (59.8366, 30.4272) ₽12/kWh | ₽599.86 | 56.56 km | 2h 37m | 6.714 ms |
| 11 | Saint Petersburg (stress test) | Branch & Bound | 6 | 1. (59.9908, 30.2059) ₽40/kWh<br>2. (59.9745, 30.2488) ₽0/kWh<br>3. (59.8366, 30.4272) ₽12/kWh | ₽599.86 | 56.56 km | 2h 37m | 0 ns |
| 11 | Saint Petersburg (stress test) | ACO | 6 | 1. (59.9908, 30.2059) ₽40/kWh<br>2. (59.9745, 30.2488) ₽0/kWh<br>3. (59.8366, 30.4272) ₽12/kWh | ₽599.86 | 56.56 km | 2h 37m | 10.935 ms |
| 12 | Beijing (dense stress test) | Dijkstra | 10 | 1. (39.9090, 116.4824) ¥1.5/kWh<br>2. (39.9096, 116.4125) ¥2/kWh<br>3. (39.9034, 116.3757) ¥1.65/kWh<br>4. (39.8689, 116.3844) ¥1.55/kWh | ¥13.8 | 24 km | 30m 16s | 0 ns |
| 12 | Beijing (dense stress test) | Genetic | 10 | 1. (39.9090, 116.4824) ¥1.5/kWh<br>2. (39.9096, 116.4125) ¥2/kWh<br>3. (39.9034, 116.3757) ¥1.65/kWh<br>4. (39.8689, 116.3844) ¥1.55/kWh | ¥13.8 | 24 km | 30m 16s | 1.978 s |
| 12 | Beijing (dense stress test) | VNS | 10 | 1. (39.9090, 116.4824) ¥1.5/kWh<br>2. (39.9096, 116.4125) ¥2/kWh<br>3. (39.9034, 116.3757) ¥1.65/kWh<br>4. (39.8689, 116.3844) ¥1.55/kWh | ¥13.8 | 24 km | 30m 16s | 9.730 ms |
| 12 | Beijing (dense stress test) | Branch & Bound | 10 | 1. (39.9090, 116.4824) ¥1.5/kWh<br>2. (39.9096, 116.4125) ¥2/kWh<br>3. (39.9034, 116.3757) ¥1.65/kWh<br>4. (39.8689, 116.3844) ¥1.55/kWh | ¥13.8 | 24 km | 30m 16s | 0 ns |
| 12 | Beijing (dense stress test) | ACO | 10 | 1. (39.9090, 116.4824) ¥1.5/kWh<br>2. (39.8970, 116.4524) ¥1.9/kWh<br>3. (39.9096, 116.4125) ¥2/kWh<br>4. (39.9034, 116.3757) ¥1.65/kWh<br>5. (39.8689, 116.3844) ¥1.55/kWh | ¥16.41 | 27.19 km | 36m 10s | 14.194 ms |

## Cross-algorithm consensus

| # | City | Dijkstra (optimal) | Genetic | VNS | Branch & Bound | ACO | Consensus |
|---|---|---|---|---|---|---|---|
| 1 | Munich | €19.04 | €19.04 | €19.04 | €19.04 | €19.04 | ✅ all 5 match |
| 2 | San Francisco | \$28.46 | \$28.46 | \$28.46 | \$28.46 | \$28.46 | ✅ all 5 match |
| 3 | Las Vegas | \$30.72 | \$30.72 | \$30.72 | \$30.72 | \$30.72 | ✅ all 5 match |
| 4 | London | £46.42 | £46.42 | £46.42 | £46.42 | £46.42 | ✅ all 5 match |
| 5 | Beijing | ¥145.07 | ¥145.07 | ¥145.07 | ¥145.07 | ¥145.07 | ✅ all 5 match |
| 6 | Boston | \$26.73 | \$26.73 | \$26.73 | \$26.73 | \$26.73 | ✅ all 5 match |
| 7 | Tokyo | ¥4142.46 | ¥4142.46 | ¥4142.46 | ¥4142.46 | ¥4142.46 | ✅ all 5 match |
| 8 | Saint Petersburg | ₽1846.86 | ₽1846.86 | ₽1846.86 | ₽1846.86 | ₽1846.86 | ✅ all 5 match |
| 9 | Moscow | ₽133.62 | ₽133.62 | ₽133.62 | ₽133.62 | ₽133.62 | ✅ all 5 match |
| 10 | Las Vegas (stress test) | \$4.26 | \$4.26 | \$4.26 | \$4.26 | \$4.96 | ⚠️ ACO diverges |
| 11 | Saint Petersburg (stress test) | ₽599.86 | ₽599.86 | ₽599.86 | ₽599.86 | ₽599.86 | ✅ all 5 match |
| 12 | Beijing (dense stress test) | ¥13.8 | ¥13.8 | ¥13.8 | ¥13.8 | ¥16.41 | ⚠️ ACO diverges |

## Conclusion — effectiveness of each method

| Algorithm | Optimal | Avg time | Max time | Verdict |
|---|---|---|---|---|
| Dijkstra | 12/12 | 56.94 µs | 507.90 µs | Exact and by far the cheapest to run — the right default for this graph size. |
| Genetic | 12/12 | 232.806 ms | 1.978 s | Exact every time, but 2-4 orders of magnitude slower than Dijkstra with no quality payoff here. |
| VNS | 12/12 | 11.353 ms | 28.103 ms | Exact every time, moderate cost — the best-value heuristic in this suite. |
| Branch & Bound | 12/12 | 51.85 µs | 506.50 µs | Exact and essentially as fast as Dijkstra — ties it for best overall. |
| ACO | 10/12 | 16.988 ms | 31.386 ms | Fast enough, but the only algorithm that ever settled for a worse route. |

- **Solution quality:** Dijkstra, Branch & Bound, Genetic, and VNS matched the true optimum in **all 12 cases**. ACO diverged twice, both times once the candidate set got dense enough to matter: Las Vegas (stress test) (+16.4%, once its 6-station pool was expanded to 8 real stations along the Fremont St.→Henderson corridor) and Beijing (dense stress test) (+18.9%).
- **Speed:** Dijkstra and Branch & Bound are functionally tied for fastest (both averaged in the 50-60 µs range), and neither ever sacrificed solution quality for it — there's no accuracy/speed trade-off between them on graphs this size. VNS was the cheapest of the metaheuristics to run. Genetic was consistently the slowest by a wide margin (up to 1.978 s on the largest graph) for the same answer Dijkstra found in under a millisecond.
- **Net assessment:** for this problem size, exact methods (Dijkstra, Branch & Bound) dominate — they're both faster and more reliable than any of the three metaheuristics. Of the heuristics, VNS is the effective one; ACO is the one to watch once the candidate set grows large — it now has two divergences in this suite, both on the two largest/densest candidate pools; Genetic's cost never bought it anything a much cheaper exact method didn't already deliver.
- **Moscow** also confirmed the router never charges just because a station is nearby — every algorithm returned the direct route when the battery didn't need topping up.
