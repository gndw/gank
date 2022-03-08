package model

type Health struct {
	Services []ServiceHealth `json:"services"`
}

type ServiceHealth struct {
	ServiceName string `json:"service_name,omitempty"`
	IsHealthy   bool   `json:"is_healthy,omitempty"`
	Status      string `json:"status,omitempty"`
}
