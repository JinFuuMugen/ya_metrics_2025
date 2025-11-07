package handler

import (
	"html/template"
	"net/http"

	"github.com/JinFuuMugen/ya_metrics_2025/internal/logger"
	"github.com/JinFuuMugen/ya_metrics_2025/internal/storage"
)

const pageTmpl = `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>Metrics</title>
	<style>
		body { font-family: sans-serif; }
		table { border-collapse: collapse; margin-bottom: 2rem; }
		th, td { border: 1px solid #999; padding: .4rem .8rem; }
		th { background: #eee; }
	</style>
</head>
<body>
	<h1>Counters</h1>
	<table>
		<tr><th>Name</th><th>Type</th><th>Value</th></tr>
		{{range .Counters}}
		<tr>
			<td>{{.Name}}</td><td>{{.Type}}</td><td>{{.Value}}</td>
		</tr>
		{{end}}
	</table>

	<h1>Gauges</h1>
	<table>
		<tr><th>Name</th><th>Type</th><th>Value</th></tr>
		{{range .Gauges}}
		<tr>
			<td>{{.Name}}</td><td>{{.Type}}</td><td>{{printf "%.3f" .Value}}</td>
		</tr>
		{{end}}
	</table>
</body>
</html>
`

type PageData struct {
	Counters []storage.Counter
	Gauges   []storage.Gauge
}

var tmpl = template.Must(template.New("page").Parse(pageTmpl))

func InfoPageHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-Type", "text/html")

	counters := storage.GetCounters()
	gauges := storage.GetGauges()

	data := PageData{Counters: counters, Gauges: gauges}
	if err := tmpl.Execute(w, data); err != nil {
		logger.Errorf("unable to parse info page: %w", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

}
