# TSHARK RESULTS
### Na maszynie wirtualnej, pierwotnie nazwanej klientem, został uruchomiony hub oraz jednocześnie narzędzie tshark do monitorowania przepływu

```bash
Alice@KittyKlient1:~$ sudo tshark -r /tmp/kitty-cross-network.pcapng
Running as user "root" and group "root". This could be dangerous.
    1 0.000000000 185.188.117.206 → 10.0.0.4     QUIC 1322 Initial, DCID=e984ac90b40f028dae89671a5c, PKN: 0, CRYPTO
    2 0.000280475 185.188.117.206 → 10.0.0.4     QUIC 1322 Initial, DCID=e984ac90b40f028dae89671a5c, PKN: 1, PADDING, CRYPTO
    3 0.002548018     10.0.0.4 → 185.188.117.206 QUIC 79 Initial, SCID=41855db5, PKN: 0, ACK
    4 0.005679382     10.0.0.4 → 185.188.117.206 QUIC 1322 Initial, SCID=41855db5, PKN: 1, ACK, PADDING, CRYPTO
    5 0.005690221     10.0.0.4 → 185.188.117.206 QUIC 1297 Handshake, SCID=41855db5
    6 0.005696693     10.0.0.4 → 185.188.117.206 QUIC 133 Protected Payload (KP0)
    7 0.107653002 185.188.117.206 → 10.0.0.4     QUIC 1322 Initial, DCID=41855db5, PKN: 2, ACK, PADDING
    8 0.107653448 185.188.117.206 → 10.0.0.4     QUIC 147 Protected Payload (KP0)
    9 0.107653508 185.188.117.206 → 10.0.0.4     QUIC 119 Protected Payload (KP0)
   10 0.108184514     10.0.0.4 → 185.188.117.206 QUIC 360 Protected Payload (KP0)
   11 0.108773623     10.0.0.4 → 185.188.117.206 QUIC 130 Protected Payload (KP0)
   12 0.352565443     10.0.0.4 → 185.188.117.206 QUIC 352 Protected Payload (KP0)
   13 0.352581736     10.0.0.4 → 185.188.117.206 QUIC 130 Protected Payload (KP0)
   14 0.414993649 185.188.117.206 → 10.0.0.4     QUIC 72 Protected Payload (KP0)
   15 8.191673371 185.188.117.206 → 10.0.0.4     QUIC 1408 Protected Payload (KP0)
   16 8.191983341     10.0.0.4 → 185.188.117.206 QUIC 1408 Protected Payload (KP0)
   17 8.192062083 185.188.117.206 → 10.0.0.4     QUIC 137 Protected Payload (KP0)
   18 8.192152866     10.0.0.4 → 185.188.117.206 QUIC 69 Protected Payload (KP0)
   19 8.279908601     10.0.0.4 → 185.188.117.206 QUIC 127 Protected Payload (KP0)
   20 8.407218196 185.188.117.206 → 10.0.0.4     QUIC 70 Protected Payload (KP0)
   21 13.414776015 185.188.117.206 → 10.0.0.4     QUIC 1322 Initial, DCID=e4e5aa614f56a23f7e6ee72c27, PKN: 0, CRYPTO
   22 13.415140558 185.188.117.206 → 10.0.0.4     QUIC 1322 Initial, DCID=e4e5aa614f56a23f7e6ee72c27, PKN: 1, PADDING, CRYPTO
   23 13.415250708     10.0.0.4 → 185.188.117.206 QUIC 79 Initial, SCID=6ee3493c, PKN: 0, ACK
   24 13.418768033     10.0.0.4 → 185.188.117.206 QUIC 1322 Initial, SCID=6ee3493c, PKN: 1, ACK, PADDING, CRYPTO
   25 13.418778799     10.0.0.4 → 185.188.117.206 QUIC 1299 Handshake, SCID=6ee3493c
   26 13.418784375     10.0.0.4 → 185.188.117.206 QUIC 133 Protected Payload (KP0)
   27 13.535158662 185.188.117.206 → 10.0.0.4     QUIC 1322 Initial, DCID=e4e5aa614f56a23f7e6ee72c27, PKN: 2, CRYPTO
   28 13.535381610     10.0.0.4 → 185.188.117.206 QUIC 79 Initial, SCID=6ee3493c, PKN: 2, ACK
   29 13.535624188 185.188.117.206 → 10.0.0.4     QUIC 1322 Initial, DCID=e4e5aa614f56a23f7e6ee72c27, PKN: 3, PADDING, CRYPTO
   30 13.535692926     10.0.0.4 → 185.188.117.206 QUIC 79 Initial, SCID=6ee3493c, PKN: 3, ACK
   31 13.619169940     10.0.0.4 → 185.188.117.206 QUIC 1322 Initial, SCID=6ee3493c, PKN: 4, PADDING, CRYPTO
   32 13.619199384     10.0.0.4 → 185.188.117.206 QUIC 1322 Initial, SCID=6ee3493c, PKN: 5, PADDING, CRYPTO
   33 13.624772941 185.188.117.206 → 10.0.0.4     QUIC 1322 Handshake, DCID=7fca4721
   34 13.625099828     10.0.0.4 → 185.188.117.206 QUIC 357 Protected Payload (KP0)
   35 13.625340596 185.188.117.206 → 10.0.0.4     QUIC 119 Protected Payload (KP0)
   36 13.625865803     10.0.0.4 → 185.188.117.206 QUIC 138 Protected Payload (KP0)
   37 13.836688386 185.188.117.206 → 10.0.0.4     QUIC 70 Protected Payload (KP0)
   38 17.305395138 185.188.117.206 → 10.0.0.4     QUIC 1408 Protected Payload (KP0)
   39 17.305632710     10.0.0.4 → 185.188.117.206 QUIC 1408 Protected Payload (KP0)
   40 17.305789103 185.188.117.206 → 10.0.0.4     QUIC 137 Protected Payload (KP0)
   41 17.305870695     10.0.0.4 → 185.188.117.206 QUIC 69 Protected Payload (KP0)
   42 17.395736665     10.0.0.4 → 185.188.117.206 QUIC 127 Protected Payload (KP0)
   43 17.438600428 185.188.117.206 → 10.0.0.4     QUIC 71 Protected Payload (KP0)
   44 17.510341867 185.188.117.206 → 10.0.0.4     QUIC 71 Protected Payload (KP0)
   45 20.275481862 185.188.117.206 → 10.0.0.4     QUIC 1451 Protected Payload (KP0)
   46 20.275482319 185.188.117.206 → 10.0.0.4     QUIC 234 Protected Payload (KP0)
   47 20.275758321     10.0.0.4 → 185.188.117.206 QUIC 1451 Protected Payload (KP0)
   48 20.275791134     10.0.0.4 → 185.188.117.206 QUIC 69 Protected Payload (KP0)
   49 20.276270529     10.0.0.4 → 185.188.117.206 QUIC 134 Protected Payload (KP0)
   50 20.276298238     10.0.0.4 → 185.188.117.206 QUIC 1451 Protected Payload (KP0)
   51 20.276304205     10.0.0.4 → 185.188.117.206 QUIC 245 Protected Payload (KP0)
   52 20.486618008 185.188.117.206 → 10.0.0.4     QUIC 1451 Protected Payload (KP0)
   53 20.486618075 185.188.117.206 → 10.0.0.4     QUIC 70 Protected Payload (KP0)
   54 20.486656312 185.188.117.206 → 10.0.0.4     QUIC 70 Protected Payload (KP0)
   55 20.512143664     10.0.0.4 → 185.188.117.206 QUIC 70 Protected Payload (KP0)
   56 23.040230475 185.188.117.206 → 10.0.0.4     QUIC 1472 Protected Payload (KP0)
   57 23.040230892 185.188.117.206 → 10.0.0.4     QUIC 232 Protected Payload (KP0)
   58 23.040590628     10.0.0.4 → 185.188.117.206 QUIC 1472 Protected Payload (KP0)
   59 23.040628372     10.0.0.4 → 185.188.117.206 QUIC 69 Protected Payload (KP0)
   60 23.040810201     10.0.0.4 → 185.188.117.206 QUIC 134 Protected Payload (KP0)
   61 23.040858940     10.0.0.4 → 185.188.117.206 QUIC 1472 Protected Payload (KP0)
   62 23.040888711     10.0.0.4 → 185.188.117.206 QUIC 245 Protected Payload (KP0)
   63 23.250299142 185.188.117.206 → 10.0.0.4     QUIC 70 Protected Payload (KP0)
   64 23.250838481 185.188.117.206 → 10.0.0.4     QUIC 1472 Protected Payload (KP0)
   65 23.250838709 185.188.117.206 → 10.0.0.4     QUIC 70 Protected Payload (KP0)
   66 23.276306720     10.0.0.4 → 185.188.117.206 QUIC 70 Protected Payload (KP0)
   67 38.250713195     10.0.0.4 → 185.188.117.206 QUIC 1483 Protected Payload (KP0)
   68 38.250746994     10.0.0.4 → 185.188.117.206 QUIC 63 Protected Payload (KP0)
   69 38.250943212     10.0.0.4 → 185.188.117.206 QUIC 1483 Protected Payload (KP0)
   70 38.250950302     10.0.0.4 → 185.188.117.206 QUIC 63 Protected Payload (KP0)
   71 38.302849524 185.188.117.206 → 10.0.0.4     QUIC 1483 Protected Payload (KP0)
   72 38.303269717 185.188.117.206 → 10.0.0.4     QUIC 70 Protected Payload (KP0)
   73 38.303452303 185.188.117.206 → 10.0.0.4     QUIC 1483 Protected Payload (KP0)
   74 38.303452419 185.188.117.206 → 10.0.0.4     QUIC 70 Protected Payload (KP0)
   75 38.328898071     10.0.0.4 → 185.188.117.206 QUIC 70 Protected Payload (KP0)
   76 38.328933800     10.0.0.4 → 185.188.117.206 QUIC 70 Protected Payload (KP0)
   77 39.529994174 185.188.117.206 → 10.0.0.4     QUIC 107 Protected Payload (KP0)
   78 39.555440273     10.0.0.4 → 185.188.117.206 QUIC 70 Protected Payload (KP0)
   79 43.008034827 185.188.117.206 → 10.0.0.4     QUIC 268 Protected Payload (KP0)
   80 43.008332818     10.0.0.4 → 185.188.117.206 QUIC 142 Protected Payload (KP0)
   81 43.008367941     10.0.0.4 → 185.188.117.206 QUIC 281 Protected Payload (KP0)
   82 43.141931871 185.188.117.206 → 10.0.0.4     QUIC 71 Protected Payload (KP0)
   83 43.142435213 185.188.117.206 → 10.0.0.4     QUIC 71 Protected Payload (KP0)
   84 46.186665986 185.188.117.206 → 10.0.0.4     QUIC 242 Protected Payload (KP0)
   85 46.187015018     10.0.0.4 → 185.188.117.206 QUIC 142 Protected Payload (KP0)
   86 46.187040099     10.0.0.4 → 185.188.117.206 QUIC 253 Protected Payload (KP0)
   87 46.291250884 185.188.117.206 → 10.0.0.4     QUIC 70 Protected Payload (KP0)
   88 46.316196250 185.188.117.206 → 10.0.0.4     QUIC 71 Protected Payload (KP0)
   89 48.540446712 185.188.117.206 → 10.0.0.4     QUIC 107 Protected Payload (KP0)
   90 48.566314036     10.0.0.4 → 185.188.117.206 QUIC 70 Protected Payload (KP0)
   91 50.483673988 185.188.117.206 → 10.0.0.4     QUIC 238 Protected Payload (KP0)
   92 50.484406248     10.0.0.4 → 185.188.117.206 QUIC 143 Protected Payload (KP0)
   93 50.484426282     10.0.0.4 → 185.188.117.206 QUIC 249 Protected Payload (KP0)
   94 50.616813836 185.188.117.206 → 10.0.0.4     QUIC 71 Protected Payload (KP0)
   95 50.616973890 185.188.117.206 → 10.0.0.4     QUIC 71 Protected Payload (KP0)
   96 52.843930671 185.188.117.206 → 10.0.0.4     QUIC 232 Protected Payload (KP0)
   97 52.844296883     10.0.0.4 → 185.188.117.206 QUIC 142 Protected Payload (KP0)
   98 52.844321055     10.0.0.4 → 185.188.117.206 QUIC 245 Protected Payload (KP0)
   99 53.073789728 185.188.117.206 → 10.0.0.4     QUIC 71 Protected Payload (KP0)
  100 53.075572633 185.188.117.206 → 10.0.0.4     QUIC 71 Protected Payload (KP0)
  101 57.037185466 185.188.117.206 → 10.0.0.4     QUIC 232 Protected Payload (KP0)
  102 57.037538005     10.0.0.4 → 185.188.117.206 QUIC 142 Protected Payload (KP0)
  103 57.037564545     10.0.0.4 → 185.188.117.206 QUIC 245 Protected Payload (KP0)
  104 57.170795816 185.188.117.206 → 10.0.0.4     QUIC 71 Protected Payload (KP0)
  105 57.171136536 185.188.117.206 → 10.0.0.4     QUIC 71 Protected Payload (KP0)
  106 63.488155349 185.188.117.206 → 10.0.0.4     QUIC 128 Protected Payload (KP0)
  107 63.488550996     10.0.0.4 → 185.188.117.206 QUIC 150 Protected Payload (KP0)
  108 63.623772262 185.188.117.206 → 10.0.0.4     QUIC 71 Protected Payload (KP0)
  109 68.608775236 185.188.117.206 → 10.0.0.4     QUIC 130 Protected Payload (KP0)
  110 68.609505998     10.0.0.4 → 185.188.117.206 QUIC 153 Protected Payload (KP0)
  111 68.844570433 185.188.117.206 → 10.0.0.4     QUIC 71 Protected Payload (KP0)
  112 69.529980785 185.188.117.206 → 10.0.0.4     QUIC 107 Protected Payload (KP0)
  113 69.555554349     10.0.0.4 → 185.188.117.206 QUIC 70 Protected Payload (KP0)
  114 72.192585089 185.188.117.206 → 10.0.0.4     QUIC 106 Protected Payload (KP0)
  115 72.192585425 185.188.117.206 → 10.0.0.4     QUIC 69 Protected Payload (KP0)
  116 72.192791046     10.0.0.4 → 185.188.117.206 QUIC 69 Protected Payload (KP0)
  117 72.192887121     10.0.0.4 → 185.188.117.206 QUIC 76 Protected Payload (KP0)
  118 78.131833288 185.188.117.206 → 10.0.0.4     QUIC 128 Protected Payload (KP0)
  119 78.132126751     10.0.0.4 → 185.188.117.206 QUIC 151 Protected Payload (KP0)
  120 78.266043891 185.188.117.206 → 10.0.0.4     QUIC 71 Protected Payload (KP0)
  121 81.721855033 185.188.117.206 → 10.0.0.4     QUIC 106 Protected Payload (KP0)
  122 81.721855484 185.188.117.206 → 10.0.0.4     QUIC 69 Protected Payload (KP0)
  123 81.722057237     10.0.0.4 → 185.188.117.206 QUIC 69 Protected Payload (KP0)
  124 81.722128151     10.0.0.4 → 185.188.117.206 QUIC 76 Protected Payload (KP0)
Alice@KittyKlient1:~$ 
```

# Interpretacja

Linia:

```bash
1 0.000000000 185.188.117.206 → 10.0.0.4     QUIC 1322 Initial, DCID=e984ac90b40f028dae89671a5c, PKN: 0, CRYPTO
```

To oznacza:

- 185.188.117.206 → Twój laptop (publiczne IP operatora)

- 10.0.0.4 → Azure VM (wewnętrzny adres)

- QUIC Initial → pierwszy pakiet QUIC

- CRYPTO → część handshake TLS 1.3

- PKN: 0 → packet number 0

To jest HELLO i początek sesji QUIC.

Dalej 

```bash
2 185.188.117.206 → 10.0.0.4 QUIC Initial, PADDING, CRYPTO
```
To dalszy ciąg handshake

### Odpowiedź Azure:

```bash
3 10.0.0.4 → 185.188.117.206 QUIC Initial, ACK
```
Azure potwierdza odbiór.

---

Potem:

```bash
5 10.0.0.4 → 185.188.117.206 QUIC Handshake
```
To jest TLS handshake — wymiana kluczy.

Oraz

```bash
6 10.0.0.4 → 185.188.117.206 QUIC Protected Payload (KP0)
```

To oznacza:

***Od tego momentu CAŁA komunikacja jest zaszyfrowana.***

I dlatego nie widać JSON‑a: QUIC + TLS 1.3 szyfrują wszystko.

# Gdzie są ramki HELLO / AUTH / DATA / GET_STATUS?

**W środku pakietów „Protected Payload (KP0)”.**

QUIC szyfruje:

- typ ramki,

- msg_id,

- payload,

- MAC,

- ogólnie większość, jak nie wszystko.

Dlatego w tshark widać tylko:

```bash
QUIC Protected Payload (KP0)
```

To jest dowód, że:

- TLS 1.3 działa,

- QUIC działa,

- E2EE działa,

- nikt nie może podsłuchać Twoich wiadomości.

# Co oznaczają duże pakiety 1408 / 1472 bajty?

To są:

- zaszyfrowane DATA,

- zaszyfrowane STATUS_RES,

- zaszyfrowane PING,

- zaszyfrowane ACK.

QUIC pakuje dane w maksymalny rozmiar MTU (~1500 bajtów).

To jest dowód, że routing działa.