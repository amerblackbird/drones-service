package services

import (
	"context"
	"drones/internal/core/domain"
	"drones/internal/ports"
)

type DronesService struct {
	repo           ports.DronesRepository
	cacheService   ports.CacheService
	eventPublisher ports.EventPublisher
	logger         ports.Logger
}

func NewDronesService(
	repo ports.DronesRepository,
	cacheService ports.CacheService,
	eventPublisher ports.EventPublisher,
	logger ports.Logger,
) ports.DronesService {
	return &DronesService{repo: repo, cacheService: cacheService, eventPublisher: eventPublisher, logger: logger}
}

func (s *DronesService) CreateDrone(ctx context.Context, drone *domain.CreateDroneRequest) (*domain.Drone, error) {
	newDrone, err := s.repo.CreateDrone(ctx, drone)
	if err != nil {
		s.logger.Error("Failed to create drone", "error", err)
		return nil, err
	}
	// Cache the newly created drone
	err = s.cacheService.Set(ctx, newDrone.ID, *newDrone, 0)
	if err != nil {
		s.logger.Error("Failed to cache newly created drone", "droneID", newDrone.ID, "error", err)
	}
	return newDrone, nil
}

func (s *DronesService) GetDroneByID(ctx context.Context, droneID string) (*domain.Drone, error) {
	// Check cache first
	var cachedDrone domain.Drone
	cacheKey := "drones:" + droneID
	err := s.cacheService.Get(ctx, cacheKey, &cachedDrone)
	if err == nil && cachedDrone.ID != "" {
		return &cachedDrone, nil
	}
	return s.repo.GetDroneByID(ctx, droneID)
}

func (s *DronesService) UpdateDrone(ctx context.Context, droneID string, update *domain.UpdateDroneRequest) (*domain.Drone, error) {
	_, err := s.GetDroneByID(ctx, droneID)
	if err != nil {
		s.logger.Error("Drone not found for update", "droneID", droneID, "error", err)
		return nil, err
	}
	updatedDrone, err := s.repo.UpdateDrone(ctx, droneID, update)
	if err != nil {
		s.logger.Error("Failed to update drone", "droneID", droneID, "error", err)
		return nil, err
	}
	// Update cache
	cacheKey := "drones:" + droneID
	err = s.cacheService.Set(ctx, cacheKey, *updatedDrone, 0)
	if err != nil {
		s.logger.Error("Failed to update cache for drone", "droneID", droneID, "error", err)
	}
	return updatedDrone, nil
}

func (s *DronesService) ListDrones(ctx context.Context, options domain.PaginationOption[domain.DroneFilter]) (*domain.Pagination[domain.DroneDTO], error) {
	return s.repo.ListDrones(ctx, options)
}

func (s *DronesService) GetDroneByFilter(ctx context.Context, options domain.DroneFilter) (*domain.Drone, error) {
	return s.repo.GetDroneByFilter(ctx, options)
}

func (s *DronesService) NearbyDrones(ctx context.Context, lat, lon, radiusKm float64) ([]*domain.Drone, error) {
	return s.repo.NearbyDrones(ctx, lat, lon, radiusKm)
}

func (s *DronesService) UpdateDroneStatus(ctx context.Context, userID, droneID string, status domain.DroneStatus) (*domain.Drone, error) {
	drone, err := s.GetDroneByID(ctx, droneID)
	if err != nil {
		s.logger.Error("Drone not found for action broken", "droneID", droneID, "error", err)
		return nil, err
	}

	if !status.IsTransitionAllowed(drone.Status) {
		s.logger.Error("Invalid status transition", "droneID", droneID, "from", drone.Status, "to", status)
		return nil, drone.Status.TransitionErr()
	}

	updatedDrone, err := s.repo.UpdateStatusBroken(ctx, userID, droneID, status)
	if err != nil {
		s.logger.Error("Failed to mark drone as broken", "droneID", droneID, "error", err)
		return nil, err
	}
	// Update cache
	cacheKey := "drones:" + droneID
	err = s.cacheService.Set(ctx, cacheKey, *updatedDrone, 0)
	if err != nil {
		s.logger.Error("Failed to update cache for broken drone", "droneID", droneID, "error", err)
	}
	return updatedDrone, nil
}

func (s *DronesService) ProcessHeartbeat(ctx context.Context, droneID string, userId string, req domain.HeartbeatRequest) (*domain.Drone, error) {
	// Validate drone exists
	drone, err := s.GetDroneByID(ctx, droneID)
	if err != nil {
		s.logger.Error("Drone not found for heartbeat", "droneID", droneID, "error", err)
		return nil, err
	}
	// Process heartbeat
	updatedDrone, err := s.repo.ProcessHeartbeat(ctx, drone.ID, userId, req)
	if err != nil {
		s.logger.Error("Failed to process heartbeat", "droneID", droneID, "error", err)
		return nil, err
	}
	// Update cache with the updated drone
	cacheKey := "drones:" + droneID
	err = s.cacheService.Set(ctx, cacheKey, *updatedDrone, 0)
	if err != nil {
		s.logger.Error("Failed to update cache for drone heartbeat", "droneID", droneID, "error", err)
	}
	return updatedDrone, nil
}
