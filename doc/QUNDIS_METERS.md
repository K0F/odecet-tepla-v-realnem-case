# Dokumentace: Měřiče tepla QUNDIS (Q caloric 5.5 / WHE5)

## 1. Technická specifikace a identifikace
- **Výrobce:** QUNDIS GmbH, Erfurt, Německo
- **Identifikátor v telegramu Wireless M-Bus:** `QDS` (Manufacturer Code)
- **Typ zařízení (Type Code):** `0x08` (Heat Cost Allocator / Poměrový indikátor topných nákladů)
- **Komunikační standard:** 
  - ČSN EN 13757-4 (Wireless M-Bus)
  - OMS (Open Metering System) specifikace, Annex C
- **Frekvence a modulace:**
  - **Frekvence:** 868.950 MHz (pásmo SRD 868 MHz)
  - **Režim:** Mode C1 (C-Mode) nebo Mode T1 (T-Mode)
  - **Rychlost přenosu:** 100 kbps (NRZ u Mode C, 3-of-6 kódování u Mode T)
  - **Vysílací interval:** Periodicky každé 2 až 4 minuty v topné sezóně

---

## 2. Dvousnímačový systém měření (ČSN EN 834)
Indikátory Qundis pracují na dvousnímačovém principu:
1. **Snímač otopného tělesa ($t_R$):** Kontaktní snímač připevněný na zadní desce indikátoru v přímém kontaktu s kovem radiátoru.
2. **Snímač teploty vzduchu ($t_L$):** Snímač měřící teplotu vnitřního vzduchu v místnosti.

### Algoritmus výpočtu přírůstku dílků:
Indikátor začíná počítat přírůstky, jakmile:
$$t_R \ge 28^\circ\text{C} \quad \text{a zároveň} \quad (t_R - t_L) \ge 4.5\text{ K}$$
Integrovaný tepelný tok:
$$\Delta U = \int c \cdot (t_R - t_L)^n \, d\tau$$
kde $n \approx 1.15 \text{ až } 1.30$ je exponent odpovídající přestupové křivce otopného tělesa.

---

## 3. Struktura přijímaného Wireless M-Bus telegramu
Příklad zachyceného paketu dekodérem `rtl_433 -R 104`:
```json
{
  "time": "2026-11-06 01:13:12",
  "model": "Wireless-MBus",
  "mode": "C",
  "M": "QDS",
  "id": 18466419,
  "version": 54,
  "type_string": "Heat Cost Allocator",
  "C": 68,
  "data": "314493441964461836087a130000200b6e0100004b6e000000426cffffcb086e010000c2086c5f3a326cffff046d0f01463b",
  "inst_hca_0": "1.000",
  "inst_hca_17": "1.000",
  "inst_date_17": "26-10-31",
  "inst_timedate_0": "26-11-06T01:15:00"
}
```

### Význam polí:
| Pole | Význam | Popis |
|:---|:---|:---|
| `id` | Sériové číslo měřiče | 8místné unikátní číslo vyražené též na krytu tělesa měřiče |
| `M` | Kód výrobce | `QDS` = Qundis |
| `inst_hca_0` | Aktuální hodnota spotřeby | Počet dílků načtených od začátku aktuálního zúčtovacího období |
| `inst_hca_17`| Hodnota k uzávěrce | Počet dílků načtených za minulé zúčtovací období k datu uzávěrky |
| `inst_date_17`| Datum uzávěrky | Datum konce zúčtovacího období (např. `26-10-31` = 31. října 2026) |
| `inst_timedate_0`| Interní hodiny | Čas vnitřních hodin měřiče (slouží k synchronizaci a ochraně proti manipulaci) |
| `CW` | Stavové slovo (Config/Status) | Bity indikující sejmutí z radiátoru, manipulaci nebo nízký stav baterie |
