package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	_ "modernc.org/sqlite"
)

type MeterInfo struct {
	ID           int     `json:"id"`
	Manufacturer string  `json:"manufacturer"`
	Model        string  `json:"model"`
	RoomLabel    string  `json:"room_label"`
	FirstSeen    string  `json:"first_seen"`
	LastSeen     string  `json:"last_seen"`
	LastHCA      float64 `json:"last_hca"`
	PrevHCA      float64 `json:"prev_hca"`
	BillingDate  string  `json:"billing_date"`
	TotalPackets int     `json:"total_packets"`
	CostCZK      float64 `json:"cost_czk"`
	EnergyKWh    float64 `json:"energy_kwh"`
}

type ReadingRecord struct {
	ID          int     `json:"id"`
	Timestamp   string  `json:"timestamp"`
	HCAVal      float64 `json:"hca_val"`
	HCAPrev     float64 `json:"hca_prev"`
	BillingDate string  `json:"billing_date"`
}

type SummaryStats struct {
	TotalMeters int     `json:"total_meters"`
	TotalUnits  float64 `json:"total_units"`
	TotalCZK    float64 `json:"total_czk"`
	TotalKWh    float64 `json:"total_kwh"`
	TotalGJ     float64 `json:"total_gj"`
	LastUpdated string  `json:"last_updated"`
}

var (
	dataLock    sync.RWMutex
	meters      map[int]*MeterInfo
	readings    []ReadingRecord
	DB_PATH     = ""
	JSONL_PATH  = ""
	PORT        = 8082
	PRICE_PER_UNIT = 2.80 // 2.80 Kč za dílek spotřební složky
	KWH_PER_UNIT   = 6.10 // 6.1 kWh tepla na 1 dílek
)

func init() {
	meters = make(map[int]*MeterInfo)
	readings = make([]ReadingRecord, 0)
}

func loadFromDatabase(dbFile string) error {
	dataLock.Lock()
	defer dataLock.Unlock()

	db, err := sql.Open("sqlite3", dbFile+"?mode=ro")
	if err != nil {
		// Try without query param
		db, err = sql.Open("sqlite", dbFile)
		if err != nil {
			return err
		}
	}
	defer db.Close()

	// 1. Read meters
	rows, err := db.Query("SELECT id, manufacturer, model, room_label, first_seen, last_seen, last_hca, prev_hca, billing_date, total_packets FROM meters ORDER BY last_seen DESC")
	if err == nil {
		defer rows.Close()
		meters = make(map[int]*MeterInfo)
		for rows.Next() {
			var m MeterInfo
			var rLabel sql.NullString
			if err := rows.Scan(&m.ID, &m.Manufacturer, &m.Model, &rLabel, &m.FirstSeen, &m.LastSeen, &m.LastHCA, &m.PrevHCA, &m.BillingDate, &m.TotalPackets); err == nil {
				if rLabel.Valid && rLabel.String != "" {
					m.RoomLabel = rLabel.String
				} else {
					m.RoomLabel = fmt.Sprintf("Radiátor #%d", m.ID)
				}
				m.CostCZK = m.LastHCA * PRICE_PER_UNIT
				m.EnergyKWh = m.LastHCA * KWH_PER_UNIT
				meters[m.ID] = &m
			}
		}
	}

	// 2. Read readings
	rRows, err := db.Query("SELECT id, timestamp, hca_val, hca_prev, billing_date FROM readings ORDER BY timestamp ASC")
	if err == nil {
		defer rRows.Close()
		readings = make([]ReadingRecord, 0)
		for rRows.Next() {
			var r ReadingRecord
			if err := rRows.Scan(&r.ID, &r.Timestamp, &r.HCAVal, &r.HCAPrev, &r.BillingDate); err == nil {
				readings = append(readings, r)
			}
		}
	}

	return nil
}

// Fallback loader if jsonl used directly
func loadFromJSONL(jsonlFile string) error {
	dataLock.Lock()
	defer dataLock.Unlock()

	data, err := os.ReadFile(jsonlFile)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	readings = make([]ReadingRecord, 0)
	meters = make(map[int]*MeterInfo)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var raw map[string]interface{}
		if err := json.Unmarshal([]byte(line), &raw); err != nil {
			continue
		}

		idFloat, ok := raw["id"].(float64)
		if !ok {
			continue
		}
		id := int(idFloat)
		timeStr, _ := raw["time"].(string)
		mfg, _ := raw["M"].(string)

		hcaStr, _ := raw["inst_hca_0"].(string)
		hcaVal, _ := strconv.ParseFloat(strings.TrimSpace(hcaStr), 64)

		prevStr, _ := raw["inst_hca_17"].(string)
		prevVal, _ := strconv.ParseFloat(strings.TrimSpace(prevStr), 64)

		bDate, _ := raw["inst_date_17"].(string)

		rec := ReadingRecord{
			ID:          id,
			Timestamp:   timeStr,
			HCAVal:      hcaVal,
			HCAPrev:     prevVal,
			BillingDate: bDate,
		}
		readings = append(readings, rec)

		if m, exists := meters[id]; exists {
			m.LastSeen = timeStr
			m.LastHCA = hcaVal
			m.PrevHCA = prevVal
			m.BillingDate = bDate
			m.TotalPackets++
			m.CostCZK = m.LastHCA * PRICE_PER_UNIT
			m.EnergyKWh = m.LastHCA * KWH_PER_UNIT
		} else {
			meters[id] = &MeterInfo{
				ID:           id,
				Manufacturer: mfg,
				Model:        "Heat Cost Allocator",
				RoomLabel:    fmt.Sprintf("Radiátor #%d", id),
				FirstSeen:    timeStr,
				LastSeen:     timeStr,
				LastHCA:      hcaVal,
				PrevHCA:      prevVal,
				BillingDate:  bDate,
				TotalPackets: 1,
				CostCZK:      hcaVal * PRICE_PER_UNIT,
				EnergyKWh:    hcaVal * KWH_PER_UNIT,
			}
		}
	}
	return nil
}

func reloadData() {
	if DB_PATH != "" && fileExists(DB_PATH) {
		if err := loadFromDatabase(DB_PATH); err == nil {
			return
		}
	}
	if JSONL_PATH != "" && fileExists(JSONL_PATH) {
		loadFromJSONL(JSONL_PATH)
	}
}

func fileExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}

func getSummaryStats() SummaryStats {
	dataLock.RLock()
	defer dataLock.RUnlock()

	totalUnits := 0.0
	lastTs := ""
	for _, m := range meters {
		totalUnits += m.LastHCA
		if m.LastSeen > lastTs {
			lastTs = m.LastSeen
		}
	}

	totalCZK := totalUnits * PRICE_PER_UNIT
	totalKWh := totalUnits * KWH_PER_UNIT
	totalGJ := totalKWh * 0.0036

	return SummaryStats{
		TotalMeters: len(meters),
		TotalUnits:  totalUnits,
		TotalCZK:    totalCZK,
		TotalKWh:    totalKWh,
		TotalGJ:     totalGJ,
		LastUpdated: lastTs,
	}
}

func httpserver(w http.ResponseWriter, request *http.Request) {
	reloadData()

	dataLock.RLock()
	stats := getSummaryStats()

	// Extract sorted meter IDs
	var meterIDs []int
	for id := range meters {
		meterIDs = append(meterIDs, id)
	}
	sort.Ints(meterIDs)

	// Collect unique timestamps for X Axis
	timeMap := make(map[string]bool)
	for _, r := range readings {
		if r.Timestamp != "" {
			timeMap[r.Timestamp] = true
		}
	}
	var times []string
	for t := range timeMap {
		times = append(times, t)
	}
	sort.Strings(times)

	// Build Series per meter
	seriesData := make(map[int][]opts.LineData)
	for _, id := range meterIDs {
		seriesData[id] = make([]opts.LineData, 0)
	}

	// Latest known val per meter
	lastVal := make(map[int]float64)
	for _, t := range times {
		// find if any reading happened at time t
		for _, r := range readings {
			if r.Timestamp == t {
				lastVal[r.ID] = r.HCAVal
			}
		}
		for _, id := range meterIDs {
			val := lastVal[id]
			seriesData[id] = append(seriesData[id], opts.LineData{Value: val})
		}
	}
	dataLock.RUnlock()

	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Theme:  types.ThemeWesteros,
			Width:  "100%",
			Height: "60vh",
		}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Odečet tepla SOB23 // Radiátory & Indikátory",
			Subtitle: fmt.Sprintf("868.95 MHz Wireless M-Bus Mode C/T • Qundis QDS • %d radiátorů celkem", stats.TotalMeters),
			TitleStyle: &opts.TextStyle{
				Color: "#f8fafc",
			},
			SubtitleStyle: &opts.TextStyle{
				Color: "#94a3b8",
			},
		}),
		charts.WithColorsOpts(opts.Colors{
			"#ff5722", // warm coral
			"#00bcd4", // cyan
			"#10b981", // emerald
			"#f59e0b", // amber
			"#8b5cf6", // violet
			"#ec4899", // pink
			"#3b82f6", // blue
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:    true,
			Trigger: "axis",
			AxisPointer: &opts.AxisPointer{
				Type: "cross",
			},
		}),
		charts.WithLegendOpts(opts.Legend{
			Show: true,
			Top:  "32px",
			TextStyle: &opts.TextStyle{
				Color: "#cbd5e1",
			},
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:  "slider",
			Start: 0,
			End:   100,
		}),
	)

	if len(times) == 0 {
		times = []string{time.Now().Format("2006-01-02 15:04:05")}
	}

	line.SetXAxis(times)
	for _, id := range meterIDs {
		label := meters[id].RoomLabel
		data := seriesData[id]
		if len(data) == 0 {
			data = []opts.LineData{{Value: meters[id].LastHCA}}
		}
		line.AddSeries(label, data, charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
	}
	line.SetSeriesOptions(
		charts.WithLineStyleOpts(opts.LineStyle{Width: 2.5}),
	)

	var chartBuf bytes.Buffer
	line.Render(&chartBuf)
	chartHTML := chartBuf.String()

	// Prepare Meter Table
	dataLock.RLock()
	var meterList []*MeterInfo
	for _, m := range meters {
		meterList = append(meterList, m)
	}
	sort.Slice(meterList, func(i, j int) bool {
		return meterList[i].LastSeen > meterList[j].LastSeen
	})
	dataLock.RUnlock()

	pageData := struct {
		Stats     SummaryStats
		Meters    []*MeterInfo
		ChartHTML template.HTML
	}{
		Stats:     stats,
		Meters:    meterList,
		ChartHTML: template.HTML(chartHTML),
	}

	tmpl := template.Must(template.New("dashboard").Parse(dashboardTemplate))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, pageData)
}

func eventsHandler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	targetFile := DB_PATH
	if targetFile == "" || !fileExists(targetFile) {
		targetFile = JSONL_PATH
	}

	lastModTime := time.Time{}
	if fi, err := os.Stat(targetFile); err == nil {
		lastModTime = fi.ModTime()
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	fmt.Fprintf(w, ": ping\n\n")
	flusher.Flush()

	notify := r.Context().Done()
	for {
		select {
		case <-notify:
			return
		case <-ticker.C:
			if fi, err := os.Stat(targetFile); err == nil {
				if fi.ModTime().After(lastModTime) {
					lastModTime = fi.ModTime()
					fmt.Fprintf(w, "data: reload\n\n")
					flusher.Flush()
				}
			}
		}
	}
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	reloadData()
	dataLock.RLock()
	defer dataLock.RUnlock()

	out := struct {
		Stats    SummaryStats           `json:"stats"`
		Meters   map[int]*MeterInfo     `json:"meters"`
		Readings []ReadingRecord        `json:"readings"`
	}{
		Stats:    getSummaryStats(),
		Meters:   meters,
		Readings: readings,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

func lastTimeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	stats := getSummaryStats()
	fmt.Fprintf(w, "%s", stats.LastUpdated)
}

func main() {
	flag.StringVar(&DB_PATH, "db", "/home/kof/src/observator/data/heat_meters.db", "Path to SQLite database.")
	flag.StringVar(&JSONL_PATH, "jsonl", "/home/kof/src/observator/data/heat_meters.jsonl", "Path to JSONL log file.")
	flag.IntVar(&PORT, "p", 8082, "HTTP port to listen on.")
	flag.Float64Var(&PRICE_PER_UNIT, "price", 2.80, "Price per heat unit in CZK.")
	flag.Parse()

	// If absolute path doesn't exist, try local
	if !fileExists(DB_PATH) && fileExists("heat_meters.db") {
		DB_PATH = "heat_meters.db"
	}
	if !fileExists(JSONL_PATH) && fileExists("heat_meters.jsonl") {
		JSONL_PATH = "heat_meters.jsonl"
	}

	reloadData()

	r := mux.NewRouter()
	r.HandleFunc("/", httpserver).Methods("GET")
	r.HandleFunc("/events", eventsHandler).Methods("GET")
	r.HandleFunc("/heat.json", jsonHandler).Methods("GET")
	r.HandleFunc("/last_time", lastTimeHandler).Methods("GET")

	c := cors.Default().Handler(r)
	log.Printf("[*] Observatoř Heat Dashboard SOB23 listening on http://0.0.0.0:%d\n", PORT)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), c); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

const dashboardTemplate = `
<!DOCTYPE html>
<html lang="cs">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Odečet tepla SOB23 // Radiátory & Indikátory</title>
    <style>
        :root {
            --bg-main: #0b0f19;
            --bg-card: #151c2c;
            --border-card: #222d45;
            --text-main: #f1f5f9;
            --text-muted: #94a3b8;
            --accent-fire: #ff5722;
            --accent-cyan: #00bcd4;
            --accent-green: #10b981;
            --accent-amber: #f59e0b;
        }
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body {
            background-color: var(--bg-main);
            color: var(--text-main);
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
            padding: 18px 24px;
            min-height: 100vh;
        }
        .header-bar {
            display: flex;
            justify-content: space-between;
            align-items: center;
            border-bottom: 1px solid var(--border-card);
            padding-bottom: 14px;
            margin-bottom: 20px;
        }
        .header-title h1 {
            font-size: 22px;
            font-weight: 700;
            letter-spacing: -0.5px;
            color: var(--text-main);
        }
        .header-title p {
            font-size: 13px;
            color: var(--text-muted);
            margin-top: 4px;
        }
        .badge-live {
            background: rgba(16, 185, 129, 0.15);
            color: var(--accent-green);
            border: 1px solid rgba(16, 185, 129, 0.3);
            padding: 4px 10px;
            border-radius: 9999px;
            font-size: 12px;
            font-weight: 600;
            display: inline-flex;
            align-items: center;
            gap: 6px;
        }
        .badge-live::before {
            content: "";
            width: 8px;
            height: 8px;
            background-color: var(--accent-green);
            border-radius: 50%;
            display: inline-block;
            box-shadow: 0 0 8px var(--accent-green);
        }
        .stats-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 16px;
            margin-bottom: 22px;
        }
        .stat-card {
            background: var(--bg-card);
            border: 1px solid var(--border-card);
            border-radius: 10px;
            padding: 16px 20px;
        }
        .stat-label {
            font-size: 12px;
            text-transform: uppercase;
            letter-spacing: 0.5px;
            color: var(--text-muted);
            margin-bottom: 6px;
        }
        .stat-value {
            font-size: 26px;
            font-weight: 700;
            color: var(--text-main);
        }
        .stat-sub {
            font-size: 12px;
            color: var(--text-muted);
            margin-top: 4px;
        }
        .chart-container {
            background: var(--bg-card);
            border: 1px solid var(--border-card);
            border-radius: 12px;
            padding: 16px;
            margin-bottom: 24px;
        }
        .table-container {
            background: var(--bg-card);
            border: 1px solid var(--border-card);
            border-radius: 12px;
            padding: 18px 20px;
            overflow-x: auto;
        }
        .table-title {
            font-size: 16px;
            font-weight: 600;
            margin-bottom: 14px;
            color: var(--text-main);
        }
        table {
            width: 100%;
            border-collapse: collapse;
            font-size: 13px;
        }
        th {
            text-align: left;
            padding: 10px 14px;
            color: var(--text-muted);
            border-bottom: 1px solid var(--border-card);
            font-weight: 600;
        }
        td {
            padding: 12px 14px;
            border-bottom: 1px solid rgba(255, 255, 255, 0.04);
            color: var(--text-main);
        }
        tr:hover td {
            background-color: rgba(255, 255, 255, 0.02);
        }
        .tag-mfg {
            background: rgba(0, 188, 212, 0.15);
            color: var(--accent-cyan);
            padding: 2px 8px;
            border-radius: 4px;
            font-size: 11px;
            font-family: monospace;
        }
        .price-badge {
            color: var(--accent-amber);
            font-weight: 600;
        }
    </style>
    <script>
        const evtSource = new EventSource('/events');
        evtSource.onmessage = function(event) {
            if (event.data === 'reload') {
                window.location.reload();
            }
        };
    </script>
</head>
<body>
    <div class="header-bar">
        <div class="header-title">
            <h1>Odečet tepla SOB23 // Radiátory & Indikátory</h1>
            <p>Wireless M-Bus 868.95 MHz Mode C&T • QUNDIS QDS • ČSN EN 834</p>
        </div>
        <div>
            <span class="badge-live">ŽIVÝ PŘÍJEM AKTIVNÍ</span>
        </div>
    </div>

    <div class="stats-grid">
        <div class="stat-card">
            <div class="stat-label">Naměřená spotřeba</div>
            <div class="stat-value" style="color: var(--accent-fire);">{{printf "%.1f" .Stats.TotalUnits}} <span style="font-size: 16px;">dílků</span></div>
            <div class="stat-sub">Součet všech zachycených radiátorů</div>
        </div>
        <div class="stat-card">
            <div class="stat-label">Odhad nákladů (Spotřební)</div>
            <div class="stat-value" style="color: var(--accent-amber);">{{printf "%.2f" .Stats.TotalCZK}} <span style="font-size: 16px;">Kč</span></div>
            <div class="stat-sub">Sazba 2.80 Kč / dílek (Pražská tepl.)</div>
        </div>
        <div class="stat-card">
            <div class="stat-label">Tepelná energie</div>
            <div class="stat-value" style="color: var(--accent-cyan);">{{printf "%.1f" .Stats.TotalKWh}} <span style="font-size: 16px;">kWh</span></div>
            <div class="stat-sub">{{printf "%.2f" .Stats.TotalGJ}} GJ odevzdaného tepla</div>
        </div>
        <div class="stat-card">
            <div class="stat-label">Aktivní tělesa</div>
            <div class="stat-value" style="color: var(--accent-green);">{{.Stats.TotalMeters}}</div>
            <div class="stat-sub">Poslední paket: {{.Stats.LastUpdated}}</div>
        </div>
    </div>

    <div class="chart-container">
        {{.ChartHTML}}
    </div>

    <div class="table-container">
        <div class="table-title">Přehled otopných těles a rozdělovačů nákladů</div>
        <table>
            <thead>
                <tr>
                    <th>ID Měřiče</th>
                    <th>Výrobce</th>
                    <th>Místnost / Popis</th>
                    <th>Stav (Dílky)</th>
                    <th>Minulé období</th>
                    <th>Odhad Kč</th>
                    <th>Odhad kWh</th>
                    <th>Poslední signál</th>
                    <th>Paketů</th>
                </tr>
            </thead>
            <tbody>
                {{range .Meters}}
                <tr>
                    <td><strong>#{{.ID}}</strong></td>
                    <td><span class="tag-mfg">{{.Manufacturer}}</span></td>
                    <td>{{.RoomLabel}}</td>
                    <td><strong style="color: var(--accent-fire);">{{printf "%.1f" .LastHCA}}</strong></td>
                    <td style="color: var(--text-muted);">{{printf "%.1f" .PrevHCA}}</td>
                    <td class="price-badge">{{printf "%.1f Kč" .CostCZK}}</td>
                    <td>{{printf "%.1f kWh" .EnergyKWh}}</td>
                    <td>{{.LastSeen}}</td>
                    <td>{{.TotalPackets}}</td>
                </tr>
                {{end}}
            </tbody>
        </table>
    </div>
</body>
</html>
`
