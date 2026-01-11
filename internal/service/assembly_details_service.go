package service

import (
	"context"

	"compConfigurator/internal/repo"
)

type AssemblyDetails struct {
	Assembly      repo.Assembly            `json:"assembly"`
	Compatibility repo.CompatibilityResult `json:"compatibility"`
	RamKits       []repo.RamKit            `json:"ram_kits"`
	Drives        []repo.Drive             `json:"drives"`
}

type AssemblyDetailsService struct {
	assemblies *repo.AssembliesRepo
	details    *repo.AssemblyDetailsRepo
}

func NewAssemblyDetailsService(a *repo.AssembliesRepo, d *repo.AssemblyDetailsRepo) *AssemblyDetailsService {
	return &AssemblyDetailsService{assemblies: a, details: d}
}

func (s *AssemblyDetailsService) Get(ctx context.Context, userID, assemblyID int64) (AssemblyDetails, error) {
	a, err := s.assemblies.GetByID(ctx, assemblyID, userID)
	if err != nil {
		return AssemblyDetails{}, err
	}

	comp, err := s.assemblies.GetCompatibility(ctx, assemblyID)
	if err != nil {
		return AssemblyDetails{}, err
	}

	ram, err := s.details.ListRamKits(ctx, assemblyID)
	if err != nil {
		return AssemblyDetails{}, err
	}

	drives, err := s.details.ListDrives(ctx, assemblyID)
	if err != nil {
		return AssemblyDetails{}, err
	}

	return AssemblyDetails{
		Assembly:      a,
		Compatibility: comp,
		RamKits:       ram,
		Drives:        drives,
	}, nil
}

func (s *AssemblyDetailsService) GetPublic(ctx context.Context, assemblyID int64) (AssemblyDetails, error) {
	a, err := s.assemblies.GetPublicByID(ctx, assemblyID)
	if err != nil {
		return AssemblyDetails{}, err
	}

	comp, err := s.assemblies.GetCompatibility(ctx, assemblyID)
	if err != nil {
		return AssemblyDetails{}, err
	}

	ram, err := s.details.ListRamKits(ctx, assemblyID)
	if err != nil {
		return AssemblyDetails{}, err
	}

	drives, err := s.details.ListDrives(ctx, assemblyID)
	if err != nil {
		return AssemblyDetails{}, err
	}

	return AssemblyDetails{
		Assembly:      a,
		Compatibility: comp,
		RamKits:       ram,
		Drives:        drives,
	}, nil
}
