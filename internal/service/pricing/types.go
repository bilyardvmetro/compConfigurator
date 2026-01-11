package pricing

import "compConfigurator/internal/repo"

type PricingLine struct {
	ComponentType string      `json:"component_type"` // CPU/GPU/...
	ComponentID   int64       `json:"component_id"`
	Title         string      `json:"title"` // имя компонента (для UI)
	Offer         *repo.Offer `json:"offer,omitempty"`
}

type PricingResult struct {
	AssemblyID      int64         `json:"assembly_id"`
	Lines           []PricingLine `json:"lines"`
	TotalPriceCents int64         `json:"total_price_cents"`
	AllAvailable    bool          `json:"all_available"`
}
