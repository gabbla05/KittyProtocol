# UWAGA UWAGA TUTAJ TE WYNIKI SĄ JUŻ DLA NOWEGO PLIKU TESTOWEGO I URUCHAMIA SIE KOMENDĄ:
`go test ./hub -bench=BenchmarkHubRouting -run='^$'` (przynajmniej na linuxie)
# Michał - Arch Linux
| Date | PC Name | Max Cores | Used Cores | Packets | Duration | Latency/pkt | Throughput |
|--- |--- |--- |--- |--- |--- |--- |--- |
| 2026-05-24 22:54:54 | archlinux | 16 | 1 | 10000 | 126ms | 12592.10 ns | 79414.89 msg/s |
| 2026-05-24 22:54:54 | archlinux | 16 | 1 | 10000 | 144ms | 14381.22 ns | 69535.11 msg/s |
| 2026-05-24 22:54:54 | archlinux | 16 | 1 | 10000 | 142ms | 14179.01 ns | 70526.77 msg/s |
| 2026-05-24 22:54:54 | archlinux | 16 | 1 | 10000 | 135ms | 13464.64 ns | 74268.62 msg/s |
| 2026-05-24 22:54:54 | archlinux | 16 | 1 | 10000 | 135ms | 13508.76 ns | 74026.04 msg/s |
| 2026-05-24 22:54:54 | archlinux | 16 | 1 | 10000 | 140ms | 14047.66 ns | 71186.25 msg/s |
| 2026-05-24 22:54:55 | archlinux | 16 | 1 | 10000 | 133ms | 13261.21 ns | 75407.93 msg/s |
| 2026-05-24 22:54:55 | archlinux | 16 | 1 | 10000 | 133ms | 13311.79 ns | 75121.38 msg/s |
| 2026-05-24 22:54:55 | archlinux | 16 | 1 | 10000 | 134ms | 13375.59 ns | 74763.08 msg/s |
| 2026-05-24 22:54:55 | archlinux | 16 | 1 | 10000 | 126ms | 12588.32 ns | 79438.69 msg/s |
| 2026-05-24 22:54:55 | archlinux | 16 | 1 | 10000 | 141ms | 14131.19 ns | 70765.43 msg/s |
| 2026-05-24 22:54:55 | archlinux | 16 | 2 | 10000 | 68ms | 6767.00 ns | 147776.04 msg/s |
| 2026-05-24 22:54:55 | archlinux | 16 | 2 | 10000 | 81ms | 8085.50 ns | 123678.20 msg/s |
| 2026-05-24 22:54:55 | archlinux | 16 | 2 | 10000 | 85ms | 8493.89 ns | 117731.68 msg/s |
| 2026-05-24 22:54:55 | archlinux | 16 | 2 | 10000 | 75ms | 7508.89 ns | 133175.56 msg/s |
| 2026-05-24 22:54:56 | archlinux | 16 | 2 | 10000 | 76ms | 7593.05 ns | 131699.41 msg/s |
| 2026-05-24 22:54:56 | archlinux | 16 | 2 | 10000 | 80ms | 7953.86 ns | 125725.12 msg/s |
| 2026-05-24 22:54:56 | archlinux | 16 | 2 | 10000 | 78ms | 7840.67 ns | 127540.18 msg/s |
| 2026-05-24 22:54:56 | archlinux | 16 | 2 | 10000 | 77ms | 7667.40 ns | 130422.25 msg/s |
| 2026-05-24 22:54:56 | archlinux | 16 | 2 | 10000 | 80ms | 8013.04 ns | 124796.61 msg/s |
| 2026-05-24 22:54:56 | archlinux | 16 | 4 | 10000 | 48ms | 4808.41 ns | 207969.06 msg/s |
| 2026-05-24 22:54:56 | archlinux | 16 | 4 | 10000 | 51ms | 5085.67 ns | 196631.05 msg/s |
| 2026-05-24 22:54:56 | archlinux | 16 | 4 | 10000 | 48ms | 4785.55 ns | 208962.28 msg/s |
| 2026-05-24 22:54:56 | archlinux | 16 | 4 | 10000 | 49ms | 4857.56 ns | 205864.88 msg/s |
| 2026-05-24 22:54:56 | archlinux | 16 | 4 | 10000 | 48ms | 4776.44 ns | 209360.88 msg/s |
| 2026-05-24 22:54:56 | archlinux | 16 | 4 | 10000 | 54ms | 5359.38 ns | 186588.91 msg/s |
| 2026-05-24 22:54:56 | archlinux | 16 | 4 | 10000 | 49ms | 4885.94 ns | 204669.02 msg/s |
| 2026-05-24 22:54:56 | archlinux | 16 | 4 | 10000 | 50ms | 4987.42 ns | 200504.66 msg/s |
| 2026-05-24 22:54:56 | archlinux | 16 | 8 | 10000 | 32ms | 3247.62 ns | 307917.57 msg/s |
| 2026-05-24 22:54:56 | archlinux | 16 | 8 | 10000 | 33ms | 3307.67 ns | 302327.32 msg/s |
| 2026-05-24 22:54:56 | archlinux | 16 | 8 | 10000 | 35ms | 3491.79 ns | 286386.27 msg/s |
| 2026-05-24 22:54:56 | archlinux | 16 | 8 | 10000 | 35ms | 3513.19 ns | 284641.87 msg/s |
| 2026-05-24 22:54:56 | archlinux | 16 | 8 | 10000 | 37ms | 3662.47 ns | 273039.98 msg/s |
| 2026-05-24 22:54:57 | archlinux | 16 | 8 | 10000 | 45ms | 4542.15 ns | 220160.20 msg/s |
| 2026-05-24 22:54:57 | archlinux | 16 | 8 | 10000 | 46ms | 4612.47 ns | 216803.52 msg/s |
| 2026-05-24 22:54:57 | archlinux | 16 | 16 | 10000 | 42ms | 4212.55 ns | 237385.88 msg/s |
| 2026-05-24 22:54:57 | archlinux | 16 | 16 | 10000 | 41ms | 4082.31 ns | 244959.19 msg/s |
| 2026-05-24 22:54:57 | archlinux | 16 | 16 | 10000 | 37ms | 3706.06 ns | 269828.60 msg/s |
| 2026-05-24 22:54:57 | archlinux | 16 | 16 | 10000 | 35ms | 3461.30 ns | 288908.62 msg/s |
| 2026-05-24 22:54:57 | archlinux | 16 | 16 | 10000 | 36ms | 3576.35 ns | 279614.42 msg/s |
| 2026-05-24 22:54:57 | archlinux | 16 | 16 | 10000 | 34ms | 3381.36 ns | 295738.57 msg/s |
| 2026-05-24 22:54:57 | archlinux | 16 | 16 | 10000 | 37ms | 3676.25 ns | 272016.52 msg/s |

## Interpretacja:
- 1 → 2 rdzenie
    
    Przepustowość prawie ×2.
    Idealne skalowanie.

- 2 → 4 rdzenie
    
    Przepustowość prawie ×1.5.
    Bardzo dobre skalowanie.

- 4 → 8 rdzeni

    Przepustowość prawie ×1.4.
    Wciąż dobre skalowanie.

- 8 → 16 rdzeni

    Przepustowość trochę spada.
Dlaczego?

Bo:

- QUIC listener ma jedną goroutine acceptującą,
- routing ma locki,
- JSON ma overhead,
- scheduler Go ma koszty synchronizacji,
- 16 goroutines zaczyna walczyć o zasoby.

To jest normalne — każdy system ma punkt nasycenia. Żaden serwer na świecie nie skaluje się liniowo do 16 rdzeni.

Dla mnie oznacza to, że Hub na mojej maszynie:

    - przy 1 rdzeniu obsługuje ~75k wiadomości na sekundę,

    - przy 8 rdzeniach obsługuje ~300k wiadomości na sekundę.


# KittyProtocol - Hub Routing Performance History (DEPRECATED)

## GOŁEK:
### go test ./hub -bench=BenchmarkHubRouting -v -run=^$

| Date | PC Name | Max Cores | Used Cores | Packets | Duration | Latency/pkt | Throughput |
|--- |--- |--- |--- |--- |--- |--- |--- |
| 2026-05-17 20:55:53 | Dell_Latitude | 12 | 1 | 1 | 1ms | 1032800.00 ns | 968.24 msg/s |
| 2026-05-17 20:56:03 | Dell_Latitude | 12 | 2 | 1 | 1ms | 1381500.00 ns | 723.85 msg/s |
| 2026-05-17 20:56:04 | Dell_Latitude | 12 | 2 | 100 | 55ms | 549923.00 ns | 1818.44 msg/s |
| 2026-05-17 20:56:14 | Dell_Latitude | 12 | 4 | 1 | 1ms | 1004600.00 ns | 995.42 msg/s |
| 2026-05-17 20:56:14 | Dell_Latitude | 12 | 4 | 100 | 54ms | 535490.00 ns | 1867.45 msg/s |
| 2026-05-17 20:56:24 | Dell_Latitude | 12 | 8 | 1 | 1ms | 864700.00 ns | 1156.47 msg/s |
| 2026-05-17 20:56:24 | Dell_Latitude | 12 | 8 | 100 | 54ms | 535014.00 ns | 1869.11 msg/s |
| 2026-05-17 20:56:26 | Dell_Latitude | 12 | 8 | 2233 | 1.211s | 542451.14 ns | 1843.48 msg/s |
| 2026-05-17 20:58:14 | Dell_Latitude | 12 | 1 | 1 | 1ms | 999700.00 ns | 1000.30 msg/s |
| 2026-05-17 20:58:14 | Dell_Latitude | 12 | 1 | 100 | 79ms | 788081.00 ns | 1268.91 msg/s |
| 2026-05-17 20:58:24 | Dell_Latitude | 12 | 2 | 1 | 1ms | 505500.00 ns | 1978.24 msg/s |
| 2026-05-17 20:58:35 | Dell_Latitude | 12 | 4 | 1 | 1ms | 522300.00 ns | 1914.61 msg/s |
| 2026-05-17 20:58:35 | Dell_Latitude | 12 | 4 | 100 | 55ms | 548279.00 ns | 1823.89 msg/s |
| 2026-05-17 20:58:45 | Dell_Latitude | 12 | 8 | 1 | 1ms | 1287400.00 ns | 776.76 msg/s |
| 2026-05-17 20:58:45 | Dell_Latitude | 12 | 8 | 100 | 57ms | 568281.00 ns | 1759.69 msg/s |

## MICHAL:

| Date | PC Name | Max Cores | Used Cores | Packets | Duration | Latency/pkt | Throughput |
|--- |--- |--- |--- |--- |--- |--- |--- |
| 2026-05-21 14:40:50 | archlinux | 16 | 1 | 1 | 0s | 362630.00 ns | 2757.63 msg/s |
| 2026-05-21 14:41:00 | archlinux | 16 | 1 | 100 | 10s | 100003900.02 ns | 10.00 msg/s |
| 2026-05-21 14:41:01 | archlinux | 16 | 2 | 1 | 0s | 392806.00 ns | 2545.79 msg/s |
| 2026-05-21 14:41:11 | archlinux | 16 | 2 | 100 | 10.001s | 100007316.62 ns | 10.00 msg/s |
| 2026-05-21 14:41:11 | archlinux | 16 | 4 | 1 | 0s | 465973.00 ns | 2146.05 msg/s |
| 2026-05-21 14:41:11 | archlinux | 16 | 4 | 100 | 13ms | 127375.59 ns | 7850.80 msg/s |
| 2026-05-21 14:41:21 | archlinux | 16 | 4 | 9136 | 10.002s | 1094738.06 ns | 913.46 msg/s |
| 2026-05-21 14:41:22 | archlinux | 16 | 8 | 1 | 0s | 440666.00 ns | 2269.29 msg/s |
| 2026-05-21 14:41:22 | archlinux | 16 | 8 | 100 | 14ms | 144367.52 ns | 6926.77 msg/s |
| 2026-05-21 14:41:32 | archlinux | 16 | 8 | 8061 | 10.002s | 1240806.52 ns | 805.93 msg/s |
| 2026-05-21 14:41:32 | archlinux | 16 | 16 | 1 | 0s | 324467.00 ns | 3081.98 msg/s |
| 2026-05-21 14:41:42 | archlinux | 16 | 16 | 100 | 10.003s | 100026446.05 ns | 10.00 msg/s |

