package domain

// Stock represents a single stock entity in the system.
type Stock struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}
