# Význam "Dílků" na indikátorech topení a přepočet na Koruny (Kč)

## 1. Co ve skutečnosti znamená "Dílek" (HCA Unit)?
Indikátor topných nákladů (ITN) **neměří přímo fyzikální jednotku energie (kWh ani GJ)**.
Číslo na displeji a v telegramu (`inst_hca_0`) je **bezrozměrná poměrová odečtová jednotka**.

V drtivé většině českých bytových domů se používá tzv. **jednotná (nestupňovaná) stupnice**:
- Všechny indikátory v domě počítají jednotným algoritmem impulzy.
- Samotný surový dílek na velkém radiátoru v obýváku představuje více odevzdaného tepla než stejný dílek na malém radiátoru na chodbě.

---

## 2. Přepočtové koeficienty otopného tělesa ($K$)
Aby bylo vyúčtování spravedlivé, rozúčtovací firma (Techem, Ista, Qundis, Inmes) při roční uzávěrce násobí surové dílky $U$ součinem koeficientů:

$$U_{\text{korig}} = U \times K_Q \times K_C \times K_P$$

1. **Koeficient výkonu radiátoru ($K_Q$):**
   - Poměr jmenovitého výkonu otopného tělesa $Q_N$ (při teplotním spádu 75/65/20 °C) vůči základnímu výkonu 1000 W.
   - Příklad:
     - Standardní panelový radiátor Radik 22-0600-1000 (1700 W) ➔ $K_Q = 1.70$
     - Velký obývákový radiátor Radik 22-0600-1600 (2700 W) ➔ $K_Q = 2.70$
     - Koupelnový žebřík 500 W ➔ $K_Q = 0.50$

2. **Koeficient přestupu tepla ($K_C$):**
   - Zohledňuje tepelný odpor mezi topnou vodou, ocelovým plechem radiátoru a hliníkovou zadní deskou snímače.
   - Pohybuje se typicky v rozmezí **$0.85 \text{ až } 1.15$**.

3. **Koeficient polohy místnosti ($K_P$ dle vyhlášky č. 269/2015 Sb.):**
   - Kompenzuje tepelné ztráty místností s nevýhodnou polohou (ochlazované stěny, pod střechou, nad suterénem):
     - Místnost pod plochou střechou / nárožní: koeficient **$0.70 \text{ až } 0.85$** (sleva 15–30 %)
     - Běžná místnost uprostřed dispozice: koeficient **$1.00$**
     - Vnitřní chráněná místnost (ze všech stran vytápěná): koeficient **$1.05 \text{ až } 1.10$**

---

## 3. Finanční přepočet na Koruny (Kč)
Vyúčtování nákladů na teplo v ČR se dle vyhlášky č. 269/2015 Sb. dělí na dvě složky:
- **Základní složka (30 % až 50 %):** Platí se fixně podle podlahové plochy bytu ($m^2$), bez ohledu na to, zda radiátory hřejí.
- **Spotřební složka (50 % až 70 %):** Rozpočítává se poměrově podle korigovaných dílků $U_{\text{korig}}$.

### Reálné ceny pro Prahu (CZT / Pražská teplárenská):
- Běžná cena dálkového tepla: **1 050 až 1 350 Kč / GJ vč. DPH** ($\approx 3.80\text{ Kč / kWh}$).
- Typický byt (65 $m^2$) spotřebuje za topnou sezónu cca **20–25 GJ** tepla.
- Celkový roční náklad bytu na teplo: cca **22 000 až 30 000 Kč**.
  - Z toho základní složka (plocha): cca 10 000 až 13 000 Kč.
  - Z toho spotřební složka (radiátory): cca 12 000 až 17 000 Kč.
- Běžný počet odečtených dílků na byt za sezónu: **3 000 až 6 000 dílků** celkem (při součtu všech radiátorů v bytě).

### Pravidlo pro přepočet:
$$\mathbf{1\text{ dílek spotřební složky} \approx 2.50\text{ až } 3.50\text{ Kč}}$$
*(Výchozí standardní hodnota: **2.80 Kč za 1 dílek**).*

### Energetický ekvivalent:
$$\mathbf{1\text{ dílek} \approx 0.022\text{ GJ} \approx 6.1\text{ kWh tepla}}$$

---

## 4. Příklad výpočtu na živých datech
Máme-li radiátor s aktuálním stavem **74.0 dílků** (ID 18466389):
- Odhadovaná spotřeba energie: $74 \times 6.1\text{ kWh} \approx \mathbf{451\text{ kWh}}$ ($1.63\text{ GJ}$).
- Odhadovaná cena spotřební složky: $74 \times 2.80\text{ Kč} \approx \mathbf{207.20\text{ Kč}}$.

Radiátor se stavem **534.0 dílků** (ID 18466388):
- Odhadovaná spotřeba energie: $534 \times 6.1\text{ kWh} \approx \mathbf{3 257\text{ kWh}}$ ($11.7\text{ GJ}$).
- Odhadovaná cena spotřební složky: $534 \times 2.80\text{ Kč} \approx \mathbf{1 495.20\text{ Kč}}$.
