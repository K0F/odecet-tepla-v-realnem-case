# NÁVOD K ODEČTU A INTERPRETACI INDIKÁTORŮ TOPENÍ (MĚŘÁKŮ)
**Příručka uživatele bytu, fyzikální principy, dešifrování displeje a přepočet dílků na koruny (Kč)**

*Vydáno v rámci The Observatoř Telemetry Suite (Praha, listopad 2026)*  
*Normy a legislativa: ČSN EN 834, Zákon č. 67/2013 Sb., Vyhláška č. 269/2015 Sb. a novela č. 274/2023 Sb.*

---

## 1. Co to vlastně visí na radiátoru?

Zařízení upevněné uprostřed vašeho radiátoru **není měřič tepla (kalorimetr)** v pravém slova smyslu, ale **indikátor otopných nákladů (ITN)** neboli poměrový rozdělovač topných nákladů (Heat Cost Allocator – HCA).

Nejčastější výrobci v ČR:
- **QUNDIS** (řada *Q caloric 5.5*, *QDS*, *caloric compact 200/201/202*)
- **Siemens** (řada *WHE30*, *WHE50*, *WHE52*)
- **Techem** (řada *radio 3*, *radio 4*)
- **Ista** (řada *doprimo 3 radio*)

Tyto přístroje neměří proteklou vodu ani gigajouly přímo. Registrují **tepelnou expozici otopného tělesa v čase vůči vzduchu v místnosti**.

---

## 2. Jak číst LCD displej přístroje (Cykly zobrazení)

Všechny moderní indikátory (včetně Qundis) mají malý LCD displej napájený lithiovou baterií (životnost 10+2 roky). Protože mají pouze jeden malý jednořádkový displej, střídají v nekonečné smyčce několik obrazovek (každá svítí 2 až 4 sekundy):

```
       ┌──────────────────────────────┐
       │   [Aktuální stav]            │
       │          3 0 7               │  <-- Hrubý stav spotřeby v aktuální topné sezóně
       └──────────────────────────────┘
                      ▼
       ┌──────────────────────────────┐
       │   [Minulý stav k 31.12.]     │
       │   M      2 7 4               │  <-- Stav k rozhodnému dni (loňské vyúčtování)
       │    (nebo bliká datum)        │
       └──────────────────────────────┘
                      ▼
       ┌──────────────────────────────┐
       │   [Datum rozhodného dne]     │
       │       3 1 . 1 2 .            │  <-- Rozhodné datum předchozího odečtu
       └──────────────────────────────┘
                      ▼
       ┌──────────────────────────────┐
       │   [Test displeje]            │
       │       8 8 8 8 8              │  <-- Všechny segmenty (ověření bezvadnosti)
       └──────────────────────────────┘
```

### Přehled stavů na displeji:

1. **Aktuální stav (Current / Bez písmene):**
   - Čisté číslo bez doplňkových znaků (např. `534` nebo `74`).
   - Udává počet spotřebních dílků napočítaných **od posledního nulování** (typicky od 1. ledna půlnoci do této chvíle).
   - Toto číslo v průběhu topné sezóny neustále narůstá.

2. **Minulý stav k rozhodnému dni (Past / Symbol `M` nebo ikona paměti):**
   - Číslo doplněné písmenem `M` (Memory) nebo symbolem kalendáře/vlajky.
   - **TOTO JE HODNOTA, KTERÁ JDE DO ROČNÍHO VYÚČTOVÁNÍ!**
   - Přístroj ji automaticky zamrazí k 31. 12. ve 23:59:59 a uchovává v paměti celý následující rok.
   - Pokud vám v květnu přijde tištěné vyúčtování za loňský rok, na faktuře v kolonce *Konečný stav měřiče* **musí být přesně toto číslo**.

3. **Chybové stavy a hlášení:**
   - `Err` / `F01` / `Open`: Detekce manipulace. Na zadní straně je mikrospínač. Pokud se indikátor uvolní z montážní lišty navařené na radiátoru, zapíše se tamper alarm a odečet je neplatný.
   - `bAt`: Slabá baterie (nutná výměna servisní firmou).

---

## 3. Kde najít výrobní číslo (ID zařízení)

Každý indikátor má své unikátní celosvětové číslo (obvykle 8 číslic). Najdete ho:
1. **Přímo na čelní straně plastového krytu:** Vytištěno černou barvou nad nebo pod displejem (např. `18466388`).
2. **U čárového / datamatrix kódu:** Sériové číslo je v něm zakódováno pro optickou čtečku technika.
3. **V ročním vyúčtování:** V tabulce otopných těles je u každé místnosti uveden sloupec *„Výrobní číslo indikátoru“*.
4. **V rádiovém vysílání (Wireless M-Bus):** Přístroj toto ID vysílá v každém telegramu v poli `Device ID`.

---

## 4. Fyzikální princip: Co jsou to vlastně „dílky“?

Dílek na indikátoru **není energetická jednotka** (ani kWh, ani GJ). Je to **bezrozměrná poměrná jednotka**.

### Jak indikátor počítá (Dvoučidlové měření dle ČSN EN 834):
Moderní přístroje mají **dvě teplotní čidla**:
1. **$t_Z$ (teplota otopného tělesa):** Čidlo na zadní hliníkové desce pevně přitisknuté k radiátoru.
2. **$t_L$ (teplota místnosti):** Čidlo na přední straně krytu měřící teplotu vzduchu v pokoji.

Přístroj počítá teplotní rozdíl:
$$\Delta t = t_Z - t_L$$

Rychlost přibývání dílků je dána integrálem:
$$\Delta U \sim \int (\Delta t)^n \, dt$$
*(kde exponent $n$ je obvykle cca $1.25$ až $1.35$ dle typu tělesa).*

### Důležité prahové podmínky (Ochrana proti falešnému čítání):
- **Startovací teplota:** Indikátor začíná čítat teprve tehdy, když teplota radiátoru $t_Z$ překročí **$23^\circ\text{C}$** (u některých typů $28^\circ\text{C}$).
- **Minimální spád:** Radiátor musí být teplejší než vzduch v místnosti alespoň o **$4.5\text{ K}$** ($\Delta t \ge 4.5^\circ\text{C}$).
- **Letní režim (Ochrana proti slunci):** Pokud v létě praží na radiátor slunce, zahřejí se obě čidla ($t_Z$ i $t_L$) na vysokou teplotu (např. $32^\circ\text{C}$), ale jejich rozdíl $\Delta t$ je blízký nule. Indikátor proto **v létě nečítá**, i když je radiátor horký od slunce.

---

## 5. Jak se dílky přepočítávají na koruny (Kč)?

Vyúčtování tepla v bytových domech se řídí zákonem č. 67/2013 Sb. a prováděcí vyhláškou **č. 269/2015 Sb.** (novelovanou vyhláškou č. 274/2023 Sb.).

Celkové náklady domu na vytápění se dělí na **dvě složky**:

```
      ┌────────────────────────────────────────────────────────┐
      │          CELKOVÉ NÁKLADY DOMU NA TEPLO (100 %)         │
      └────────────────────────────────────────────────────────┘
                 │                                  │
                 ▼                                  ▼
      ┌─────────────────────────┐        ┌─────────────────────────┐
      │   ZÁKLADNÍ SLOŽKA       │        │   SPOTŘEBNÍ SLOŽKA      │
      │      (30 % až 50 %)     │        │      (50 % až 70 %)     │
      ├─────────────────────────┤        ├─────────────────────────┤
      │ Rozpočítává se přísně   │        │ Rozpočítává se podle    │
      │ podle podlahové plochy  │        │ KORIGOVANÝCH DÍLKŮ      │
      │ bytu (m²).              │        │ z radiátorových měřičů. │
      │                         │        │                         │
      │ Platí každý, i když má  │        │ Kdo netopí, neplatí ji. │
      │ celou zimu vypnuté      │        │ Kdo topí hodně, zaplatí │
      │ radiátory!              │        │ větší podíl domu.       │
      └─────────────────────────┘        └─────────────────────────┘
```

### Výpočet korigovaných dílků:
Dílek na malém radiátoru na záchodě neznamená stejné teplo jako dílek na obřím třídeskovém radiátoru v obýváku. Proto rozúčtovací firma násobí surový náměr z displeje třemi koeficienty:

$$N_{\text{korig}} = N_{\text{náměr}} \times K_Q \times K_C \times K_P$$

1. **$K_Q$ – Koeficient tepelného výkonu radiátoru:**
   - Zohledňuje velikost a typ otopného tělesa (udává se ve wattech).
   - Radiátor 2000 W má zhruba $2\times$ vyšší koeficient než radiátor 1000 W.
   *(Poznámka: U tzv. produktově škálovaných indikátorů zadává montér výkon do paměti přístroje při montáži a dílky na LCD už $K_Q$ obsahují. U neškálovaných se násobí až v počítači rozúčtovatele).*
2. **$K_C$ – Koeficient přestupu tepla:**
   - Zohledňuje tepelný odpor mezi radiátorem a zadní deskou indikátoru (litinový radiátor vs. ocelový plech Radik vs. hliníkový článek). Bývá okolo $0.90$ až $1.15$.
3. **$K_P$ – Polohový koeficient bytu (Korekce na nepříznivou polohu dle vyhlášky):**
   - Spravedlivá ochrana bytů, které mají vyšší tepelné ztráty ne vlastní vinou:
     - Rohový pokoj, severní strana: koeficient např. **0.75 – 0.85** (dílky se zlevňují o 15–25 %).
     - Pod nevytápěnou střechou / nad sklepem: koeficient např. **0.80 – 0.90**.
     - Vnitřní byt chráněný ze všech stran sousedy: koeficient **1.00**.

### Kolik stojí 1 dílek v praxi v Praze?
- V Praze se cena dálkového tepla (Pražská teplárenská) pohybuje typicky okolo **1 100 až 1 350 Kč za 1 GJ** s DPH (cca **4.00 až 4.90 Kč / kWh**).
- V průměrném bytovém domě se spotřební složka dělí celkovým součtem všech korigovaných dílků v celém domě.
- Na běžném panelákovém/cihlovém radiátoru (výkon ~1 400 W, $K_Q \times K_C \approx 1.0$):
  - **1 surový dílek na displeji $\approx$ 2.50 až 3.20 Kč ve spotřební složce**.
  - 1 dílek odpovídá odevzdanému teplu zhruba **5.5 až 6.5 kWh** ($\approx 0.022\text{ GJ}$).

---

## 6. Jak si zkontrolovat vyúčtování (Krok za krokem)

1. **Krok 1 (Silvestr / 31. 12. půlnoc):**
   - Podívejte se na displej radiátoru nebo si uložte snímek z našeho SDR monitoru `heat status`.
   - Zapište si číslo s písmenem `M` nebo hodnotu `hca_prev`.
2. **Krok 2 (Květen – příchod papírového vyúčtování):**
   - Otevřete rozúčtování nákladů na vytápění.
   - Najděte tabulku *Indikátory vytápění*.
   - Zkontrolujte výrobní čísla indikátorů u jednotlivých místností.
   - Zkontrolujte, zda sloupec *„Konečný stav k 31.12.“* přesně sedí s hodnotou `M` z vašeho radiátoru.
3. **Krok 3 (Limity vyhlášky č. 269/2015 Sb.):**
   - Vyhláška chrání nájemníky před extrémy:
     - Náklad na $m^2$ vašeho bytu **nesmí klesnout pod $-20\text{ \%}$** průměru celého domu (i kdybyste vůbec netopili, doplatíte do $-20\text{ \%}$).
     - Náklad na $m^2$ vašeho bytu **nesmí překročit $+100\text{ \%}$** průměru celého domu (ani kdybyste topili naplno s otevřenými okny, strop je dvojnásobek průměru domu).

---

## 7. Rychlá nápověda pro CLI nástroj `heat`

V naší Observatoři máte k dispozici terminálový nástroj:

```bash
# Zobrazit aktuální stav všech zachycených měřičů v domě
heat status

# Živý poslech rádio telegramů (868.95 MHz) v reálném čase
heat live

# Pojmenovat radiátor podle vašeho pokoje (např. kuchyň / pracovna)
heat label 18466388 "Observatoř - Stůl"
heat label 18466389 "Observatoř - Okno"

# Vypsat historii odečtů v JSON formátu
heat json
```

*Dokumentace vygenerována pro projekt `rtl_sdr-heat_wmbus_parser`.*
