package service

import (
	"context"

	"compConfigurator/internal/repo"
)

type AssemblyService struct {
	assemblies *repo.AssembliesRepo
}

func NewAssemblyService(a *repo.AssembliesRepo) *AssemblyService {
	return &AssemblyService{assemblies: a}
}

type AssemblyView struct {
	Assembly      repo.Assembly            `json:"assembly"`
	Compatibility repo.CompatibilityResult `json:"compatibility"`
}

func (s *AssemblyService) Create(
	ctx context.Context,
	userID int64,
	name string,
	cpuID, motherboardID, psuID, caseID int64,
	gpuID, coolerID *int64,
) (AssemblyView, error) {
	a, err := s.assemblies.Create(ctx, userID, name, cpuID, motherboardID, psuID, caseID, gpuID, coolerID)
	if err != nil {
		return AssemblyView{}, err
	}

	_ = s.assemblies.RecalcTotalPrice(ctx, a.AssemblyID)

	comp, err := s.assemblies.GetCompatibility(ctx, a.AssemblyID)
	if err != nil {
		return AssemblyView{}, err
	}

	// подтянем обновлённую цену
	a, _ = s.assemblies.GetByID(ctx, a.AssemblyID, userID)

	return AssemblyView{Assembly: a, Compatibility: comp}, nil
}

func (s *AssemblyService) Update(
	ctx context.Context,
	userID, assemblyID int64,
	name string,
	isPublic bool,
	cpuID, motherboardID, psuID, caseID int64,
	gpuID, coolerID *int64,
) (AssemblyView, error) {
	a, err := s.assemblies.Update(ctx, userID, assemblyID, name, isPublic, cpuID, motherboardID, psuID, caseID, gpuID, coolerID)
	if err != nil {
		return AssemblyView{}, err
	}

	// триггер уже пересчитал цену; но совместимость в ответе мы хотим показать:
	comp, err := s.assemblies.GetCompatibility(ctx, a.AssemblyID)
	if err != nil {
		return AssemblyView{}, err
	}

	// обновлённая цена уже в таблице:
	a, _ = s.assemblies.GetByID(ctx, a.AssemblyID, userID)

	return AssemblyView{Assembly: a, Compatibility: comp}, nil
}

func (s *AssemblyService) List(ctx context.Context, userID int64, limit, offset int) ([]AssemblyView, error) {
	items, err := s.assemblies.ListByUser(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	out := make([]AssemblyView, 0, len(items))
	for _, a := range items {
		comp, err := s.assemblies.GetCompatibility(ctx, a.AssemblyID)
		if err != nil {
			return nil, err
		}
		out = append(out, AssemblyView{Assembly: a, Compatibility: comp})
	}
	return out, nil
}

func (s *AssemblyService) Get(ctx context.Context, userID, assemblyID int64) (AssemblyView, error) {
	a, err := s.assemblies.GetByID(ctx, assemblyID, userID)
	if err != nil {
		return AssemblyView{}, err
	}
	comp, err := s.assemblies.GetCompatibility(ctx, a.AssemblyID)
	if err != nil {
		return AssemblyView{}, err
	}
	return AssemblyView{Assembly: a, Compatibility: comp}, nil
}
