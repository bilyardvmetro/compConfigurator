package service

import (
	"context"

	"compConfigurator/internal/repo"
)

type AssemblyPartsService struct {
	assemblies *repo.AssembliesRepo
	parts      *repo.AssemblyPartsRepo
}

func NewAssemblyPartsService(a *repo.AssembliesRepo, p *repo.AssemblyPartsRepo) *AssemblyPartsService {
	return &AssemblyPartsService{assemblies: a, parts: p}
}

func (s *AssemblyPartsService) AddRamKit(ctx context.Context, userID, assemblyID, ramKitID int64) (AssemblyView, error) {
	// проверим, что сборка принадлежит юзеру
	if _, err := s.assemblies.GetByID(ctx, assemblyID, userID); err != nil {
		return AssemblyView{}, err
	}
	// чтобы на FK не ловить “внутреннюю” ошибку — проверим существование
	if err := s.parts.EnsureRamKitExists(ctx, ramKitID); err != nil {
		return AssemblyView{}, err
	}

	if err := s.parts.AddRamKit(ctx, assemblyID, ramKitID); err != nil {
		return AssemblyView{}, err
	}

	a, err := s.assemblies.GetByID(ctx, assemblyID, userID)
	if err != nil {
		return AssemblyView{}, err
	}
	comp, err := s.assemblies.GetCompatibility(ctx, assemblyID)
	if err != nil {
		return AssemblyView{}, err
	}
	return AssemblyView{Assembly: a, Compatibility: comp}, nil
}

func (s *AssemblyPartsService) RemoveRamKit(ctx context.Context, userID, assemblyID, ramKitID int64) (AssemblyView, error) {
	if _, err := s.assemblies.GetByID(ctx, assemblyID, userID); err != nil {
		return AssemblyView{}, err
	}

	if err := s.parts.RemoveRamKit(ctx, assemblyID, ramKitID); err != nil {
		return AssemblyView{}, err
	}

	a, err := s.assemblies.GetByID(ctx, assemblyID, userID)
	if err != nil {
		return AssemblyView{}, err
	}
	comp, err := s.assemblies.GetCompatibility(ctx, assemblyID)
	if err != nil {
		return AssemblyView{}, err
	}
	return AssemblyView{Assembly: a, Compatibility: comp}, nil
}

func (s *AssemblyPartsService) AddDrive(ctx context.Context, userID, assemblyID, driveID int64, mountType *string) (AssemblyView, error) {
	if _, err := s.assemblies.GetByID(ctx, assemblyID, userID); err != nil {
		return AssemblyView{}, err
	}
	if err := s.parts.EnsureDriveExists(ctx, driveID); err != nil {
		return AssemblyView{}, err
	}

	if err := s.parts.AddDrive(ctx, assemblyID, driveID, mountType); err != nil {
		return AssemblyView{}, err
	}

	a, err := s.assemblies.GetByID(ctx, assemblyID, userID)
	if err != nil {
		return AssemblyView{}, err
	}
	comp, err := s.assemblies.GetCompatibility(ctx, assemblyID)
	if err != nil {
		return AssemblyView{}, err
	}
	return AssemblyView{Assembly: a, Compatibility: comp}, nil
}

func (s *AssemblyPartsService) RemoveDrive(ctx context.Context, userID, assemblyID, driveID int64) (AssemblyView, error) {
	if _, err := s.assemblies.GetByID(ctx, assemblyID, userID); err != nil {
		return AssemblyView{}, err
	}

	if err := s.parts.RemoveDrive(ctx, assemblyID, driveID); err != nil {
		return AssemblyView{}, err
	}

	a, err := s.assemblies.GetByID(ctx, assemblyID, userID)
	if err != nil {
		return AssemblyView{}, err
	}
	comp, err := s.assemblies.GetCompatibility(ctx, assemblyID)
	if err != nil {
		return AssemblyView{}, err
	}
	return AssemblyView{Assembly: a, Compatibility: comp}, nil
}
