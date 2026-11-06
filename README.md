# Odečet tepla v reálném čase

> **Autonomní telemetrie vytápění přes Wireless M-Bus (868.95 MHz), RTL-SDR přijímač, interaktivní Go ECharts dashboard a příručka k indikátorům topných nákladů (Qundis QDS / Q caloric 5.5).**

[![License: CC0-1.0](https://img.shields.io/badge/License-CC0_1.0-lightgrey.svg)](https://creativecommons.org/publicdomain/zero/1.0/)
[![Go Report Card](https://goreportcard.com/badge/github.com/K0F/odecet-tepla-v-realnem-case)](https://goreportcard.com/report/github.com/K0F/odecet-tepla-v-realnem-case)
[![Frequency: 868.95 MHz](https://img.shields.io/badge/RF-868.95%20MHz%20Mode%20C%2FT-blue.svg)]()

---

## 📖 Dokumentace a manuály k měřákům

- **[Návod k odečtu a interpretaci indikátorů topení (měřáků)](doc/NAVOD_ODECET_MERAKU.md)**:
  - Jak číst jednořádkový LCD displej měřiče (aktuální stav, loňský odečet `M` k 31. 12., kontrolní čísla).
  - Kde najít výrobní číslo (ID zařízení).
  - Co jsou to „dílky“ dle ČSN EN 834 (dvoučidlové měření $t_Z - t_L$, startovací teplota, ochrana proti letnímu slunci).
  - Jak se dílky přepočítávají na koruny (Kč) ve vyúčtování (základní složka 30–50 % dle plochy, spotřební složka 50–70 % dle korigovaných dílků $K_Q, K_C, K_P$).
  - Postup kontroly ročního vyúčtování od SVJ/družstva krok za krokem.
- **[Technická specifikace Qundis QDS](doc/QUNDIS_METERS.md)**:
  - Protokol Wireless M-Bus (EN 13757-3/4), OMS Mode C/T, rádiové telegramy a parametry vysílání.
- **[Dílky, fyzika a ceník dálkového tepla v Praze](doc/DILKY_A_CENY.md)**:
  - Podrobný rozbor přepočtu dílků na energii (kWh / GJ) a ceny Pražské teplárenské (~2.80 Kč / hrubý dílek).

---

## 🏛 Architektura systému

```
  [ Radiátory v domě ]
 (Qundis QDS / OMS C&T)
           │
           │  868.95 MHz RF bursty (každých 1–2 minuty)
           ▼
    [ RTL-SDR Dongle ]  (Fitipower FC0013 / Realtek RTL2832U)
           │
           ▼
  ┌──────────────────────────────────────────────────────────┐
  │ sdr-heat-daemon (Python / rtl_433)                       │
  │  - Bezhlavý příjem, parsování OMS Mode C/T telegramů     │
  │  - Zápis do SQLite DB (heat_meters.db)                   │
  │  - Append log (heat_meters.jsonl)                        │
  │  - Zero-wear RAM cache v /dev/shm (atomic JSON)          │
  │  - Real-time OSC vysílání do SuperCollideru (:57120)     │
  └──────────────────────────────────────────────────────────┘
           │
           ▼
  ┌──────────────────────────────────────────────────────────┐
  │ heat_graph (Go Web Server + go-echarts)                  │
  │  - Tmavé téma (Westeros), hladké křivky v reálném čase   │
  │  - Okamžitý přepočet: Dílky -> Kč -> kWh                 │
  │  - Server-Sent Events (/events) pro live auto-reload     │
  │  - Nativní podpora amd64 i arm64 (Raspberry Pi)          │
  └──────────────────────────────────────────────────────────┘
           │
           ▼
  [ Webový prohlížeč / Dashboard ]
```

---

## 🚀 Rychlý start

### 1. Požadavky
- RTL-SDR USB tuner (např. RTL2832U, R820T2, FC0013)
- `rtl_433` s podporou protokolu 104 (`Wireless M-Bus Mode C&T`)
- Python 3 & Go 1.21+

### 2. Spuštění přijímacího démona
```bash
# Ruční spuštění
python3 scripts/sdr_heat_daemon.py

# Nebo instalace systemd uživatelské služby
systemctl --user enable --now sdr-heat.service
```

### 3. Terminálový nástroj `heat`
```bash
# Aktuální stav radiátorů v dosahu
heat status

# Živý poslech radio telegramů
heat live

# Nastavení lidského popisku radiátoru podle místnosti
heat label 18466388 "Observatoř - Stůl"
heat label 18466389 "Observatoř - Okno"
```

### 4. Spuštění Go ECharts webového dashboardu
```bash
# Sestavení binárky
go build -o heat_graph main.go

# Spuštění na portu 8082
./heat_graph -p 8082 -jsonl /path/to/heat_meters.jsonl
```
Dashboard bude dostupný na `http://localhost:8082/`.

---

## 📄 Licence
Public Domain / CC0 1.0 Universal.
Anti-copyright – volné pro libovolné komunitní a nezávislé použití.
