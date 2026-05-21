🟦 Gołson (Dell Latitude, 12 wątków)

    throughput ~1800 msg/s przy 4–8 cores

    latency ~0.53–0.57 ms/pkt

    stabilny do ~2000 pakietów

🟩 Michał (Ryzen 5800H, 16 wątków)

    throughput ~7860 msg/s przy 4 cores
    → ~4x szybciej niż Gołson

    latency ~0.127 ms/pkt
    → ~4x niższa latencja

    FAIL przy dużym b.N, bo test jest agresywniejszy

📌 Wniosek:

Hub Michała działa szybciej niż u Gołsona, ale test jest tak agresywny, że nawet szybki Hub nie wyrabia przy dużym b.N.