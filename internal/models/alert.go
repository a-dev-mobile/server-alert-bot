package models

// Alert defines the structure for an alert
type Alert struct {
	Labels struct {
		Alertname string `json:"alertname"`
		Instance  string `json:"instance"`
		Job       string `json:"job"`
		Severity  string `json:"severity"`
	} `json:"labels"`
	Annotations struct {
		Description string `json:"description"`
		Summary     string `json:"summary"`
	} `json:"annotations"`
	State    string `json:"state"`
	ActiveAt string `json:"activeAt"`
	Value    string `json:"value"`
}

// AlertData holds the response structure for alert data
type AlertData struct {
	Status string `json:"status"`
	Data   struct {
		Alerts []Alert `json:"alerts"`
	} `json:"data"`
}
