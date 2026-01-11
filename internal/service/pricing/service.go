package pricing

import (
	"context"

	"compConfigurator/internal/repo"
)

type Service struct {
	components *repo.AssemblyComponentsRepo
	details    *repo.AssemblyDetailsRepo
	offers     *repo.OffersRepo
}

func New(components *repo.AssemblyComponentsRepo, details *repo.AssemblyDetailsRepo, offers *repo.OffersRepo) *Service {
	return &Service{components: components, details: details, offers: offers}
}

func (s *Service) GetForOwner(ctx context.Context, userID, assemblyID int64) (PricingResult, error) {
	comp, err := s.components.GetForOwner(ctx, assemblyID, userID)
	if err != nil {
		return PricingResult{}, err
	}

	ram, err := s.details.ListRamKits(ctx, assemblyID)
	if err != nil {
		return PricingResult{}, err
	}

	drives, err := s.details.ListDrives(ctx, assemblyID)
	if err != nil {
		return PricingResult{}, err
	}

	var lines []PricingLine
	var total int64 = 0
	allAvailable := true

	addLine := func(componentType string, componentID int64, title string) error {
		offer, err := s.offers.Best(ctx, componentType, componentID, true /*onlyAvailable*/)
		if err != nil {
			return err
		}
		if offer == nil {
			allAvailable = false
		} else {
			total += offer.PriceCents
			if !offer.Available {
				allAvailable = false
			}
		}
		lines = append(lines, PricingLine{
			ComponentType: componentType,
			ComponentID:   componentID,
			Title:         title,
			Offer:         offer,
		})
		return nil
	}

	// базовые компоненты (если выбраны)
	if comp.CPUId != nil && comp.CPUName != nil {
		if err := addLine("CPU", *comp.CPUId, *comp.CPUName); err != nil {
			return PricingResult{}, err
		}
	} else {
		allAvailable = false
	}

	if comp.MotherboardId != nil && comp.MotherboardName != nil {
		if err := addLine("MOTHERBOARD", *comp.MotherboardId, *comp.MotherboardName); err != nil {
			return PricingResult{}, err
		}
	} else {
		allAvailable = false
	}

	if comp.PSUId != nil && comp.PSUName != nil {
		if err := addLine("PSU", *comp.PSUId, *comp.PSUName); err != nil {
			return PricingResult{}, err
		}
	} else {
		allAvailable = false
	}

	if comp.CaseId != nil && comp.CaseName != nil {
		if err := addLine("CASE", *comp.CaseId, *comp.CaseName); err != nil {
			return PricingResult{}, err
		}
	} else {
		allAvailable = false
	}

	// опциональные
	if comp.GPUId != nil && comp.GPUName != nil {
		if err := addLine("GPU", *comp.GPUId, *comp.GPUName); err != nil {
			return PricingResult{}, err
		}
	}
	if comp.CoolerId != nil && comp.CoolerName != nil {
		if err := addLine("COOLER", *comp.CoolerId, *comp.CoolerName); err != nil {
			return PricingResult{}, err
		}
	}

	// RAM kits
	for _, k := range ram {
		if err := addLine("RAM", k.RamKitID, k.Name); err != nil {
			return PricingResult{}, err
		}
	}

	// Drives
	for _, d := range drives {
		if err := addLine("DRIVE", d.DriveID, d.Name); err != nil {
			return PricingResult{}, err
		}
	}

	return PricingResult{
		AssemblyID:      assemblyID,
		Lines:           lines,
		TotalPriceCents: total,
		AllAvailable:    allAvailable,
	}, nil
}
