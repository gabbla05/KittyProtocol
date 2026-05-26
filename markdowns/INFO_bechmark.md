# Co dokładnie porównujesz między maszynami?

## W Twoim formacie:

- Latency/pkt → „ile ns kosztuje obsługa jednej wiadomości na tej maszynie / tym OS”

- Throughput → „ile wiadomości/s ta maszyna realnie przepchnie przez Hub+QUIC+JSON”


## To są idealne metryki porównawcze:

- między Linux vs Windows vs WSL,

- między laptopem A vs laptopem B,

- między starym CPU vs nowym CPU.


## Na co uważać przy porównaniach?

* WSL:

    - ma dodatkową warstwę (kernel + wirtualizacja),

    - wyniki będą zwykle gorsze niż natywny Linux, ale względne różnice (np. 1 vs 2 vs 4 rdzenie) nadal mają sens.

* Windows natywny:

    - inny scheduler,

    - inne timery,

    - inne zachowanie stosu sieciowego,

    - ale dalej: throughput i latency są porównywalne jako „wydajność tej platformy”.

* Różne CPU:

    - tu wyniki są wręcz idealne do porównań — zobaczysz realny zysk z lepszego CPU.