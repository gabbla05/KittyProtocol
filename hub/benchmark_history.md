# KittyProtocol - Hub Routing Performance History

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
