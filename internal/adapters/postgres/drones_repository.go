package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"math"

	"drones/internal/core/domain"
	"drones/internal/ports"
)

type DronesRepository struct {
	db                  *sql.DB
	logger              ports.Logger
	createStmt          *sql.Stmt
	getByIDStmt         *sql.Stmt
	getByIdentifierStmt *sql.Stmt
	getByUserIDStmt     *sql.Stmt
	updateStmt          *sql.Stmt
	deleteStmt          *sql.Stmt
}

func NewDronesRepository(db *sql.DB, logger ports.Logger) ports.DronesRepository {
	repo := &DronesRepository{
		db:     db,
		logger: logger,
	}
	// Prepare statements
	if err := repo.prepareStatements(); err != nil {
		logger.Error("Failed to prepare drone repository statements", "error", err)
	}

	return repo
}

func (r *DronesRepository) prepareStatements() error {
	var err error

	// Get by ID statement
	r.getByIDStmt, err = r.db.Prepare(`
		SELECT
			id, drone_identifier, user_id, model, serial_number, manufacturer,
			max_weight_kg, max_speed_kmh, max_range_km, battery_capacity_mah,
			status, battery_level_percent, current_lat, current_lon, current_altitude,
			last_location_update_at, total_flight_hours, total_deliveries,
			last_maintenance_at, next_maintenance_due_at,
			created_at, updated_at, active, created_by_id, updated_by_id
		FROM drones
		WHERE id = $1 AND active = TRUE`)
	if err != nil {
		return fmt.Errorf("failed to prepare getByIDStmt: %w", err)
	}

	r.createStmt, err = r.db.Prepare(`
		INSERT INTO drones (
			user_id, model, serial_number, manufacturer, battery_capacity_mah, max_weight_kg, created_by_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING
			id, drone_identifier, user_id, model, serial_number, manufacturer,
			max_weight_kg, max_speed_kmh, max_range_km, battery_capacity_mah,
			status, battery_level_percent, current_lat, current_lon, current_altitude,
			last_location_update_at, total_flight_hours, total_deliveries,
			last_maintenance_at, next_maintenance_due_at,
			created_at, updated_at, active, created_by_id, updated_by_id`)
	if err != nil {
		return fmt.Errorf("failed to prepare createStmt: %w", err)
	}

	// Get by identifier statement
	r.getByIdentifierStmt, err = r.db.Prepare(`
		SELECT
			id, drone_identifier, user_id, model, serial_number, manufacturer,
			max_weight_kg, max_speed_kmh, max_range_km, battery_capacity_mah,
			status, battery_level_percent, current_lat, current_lon, current_altitude,
			last_location_update_at, total_flight_hours, total_deliveries,
			last_maintenance_at, next_maintenance_due_at,
			created_at, updated_at, active, created_by_id, updated_by_id
		FROM drones
		WHERE drone_identifier = $1 AND active = TRUE`)
	if err != nil {
		return fmt.Errorf("failed to prepare getByIdentifierStmt: %w", err)
	}

	// Get by user ID statement
	r.getByUserIDStmt, err = r.db.Prepare(`
		SELECT
			id, drone_identifier, user_id, model, serial_number, manufacturer,
			max_weight_kg, max_speed_kmh, max_range_km, battery_capacity_mah,
			status, battery_level_percent, current_lat, current_lon, current_altitude,
			last_location_update_at, total_flight_hours, total_deliveries,
			last_maintenance_at, next_maintenance_due_at,
			created_at, updated_at, active, created_by_id, updated_by_id
		FROM drones
		WHERE user_id = $1 AND active = TRUE
		ORDER BY created_at DESC`)
	if err != nil {
		return fmt.Errorf("failed to prepare getByUserIDStmt: %w", err)
	}

	// Update statement
	r.updateStmt, err = r.db.Prepare(`
		UPDATE drones SET
			model = COALESCE($2, model),
			serial_number = COALESCE($3, serial_number),
			manufacturer = COALESCE($4, manufacturer),
			max_weight_kg = COALESCE($5, max_weight_kg),
			battery_capacity_mah = COALESCE($6, battery_capacity_mah),
			status = COALESCE($7, status),
			current_lat = COALESCE($8, current_lat),
			current_lon = COALESCE($9, current_lon),
			current_altitude = COALESCE($10, current_altitude),
			last_maintenance_at = COALESCE($11, last_maintenance_at),
			next_maintenance_due_at = COALESCE($12, next_maintenance_due_at),
			updated_by_id = COALESCE($13, updated_by_id),
			updated_at = NOW()
		WHERE id = $1 AND active = TRUE
		RETURNING
			id, drone_identifier, user_id, model, serial_number, manufacturer,
			max_weight_kg, max_speed_kmh, max_range_km, battery_capacity_mah,
			status, battery_level_percent, current_lat, current_lon, current_altitude,
			last_location_update_at, total_flight_hours, total_deliveries,
			last_maintenance_at, next_maintenance_due_at,
			created_at, updated_at, active, created_by_id, updated_by_id`)
	if err != nil {
		return fmt.Errorf("failed to prepare updateStmt: %w", err)
	}

	// Delete statement
	r.deleteStmt, err = r.db.Prepare(`DELETE FROM drones WHERE id = $1`)
	if err != nil {
		return fmt.Errorf("failed to prepare deleteStmt: %w", err)
	}

	return nil
}

func (r *DronesRepository) Close() error {
	if r.getByIDStmt != nil {
		if err := r.getByIDStmt.Close(); err != nil {
			return err
		}
	}
	if r.getByIdentifierStmt != nil {
		if err := r.getByIdentifierStmt.Close(); err != nil {
			return err
		}
	}
	if r.getByUserIDStmt != nil {
		if err := r.getByUserIDStmt.Close(); err != nil {
			return err
		}
	}
	if r.updateStmt != nil {
		if err := r.updateStmt.Close(); err != nil {
			return err
		}
	}
	if r.deleteStmt != nil {
		if err := r.deleteStmt.Close(); err != nil {
			return err
		}
	}
	if r.createStmt != nil {
		if err := r.createStmt.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (r *DronesRepository) GetDB() *sql.DB {
	return r.db
}

// scanDrone scans a row into a Drone struct
func (r *DronesRepository) scanDrone(scanner interface {
	Scan(dest ...interface{}) error
}) (*domain.Drone, error) {
	var drone domain.Drone
	err := scanner.Scan(
		&drone.ID,
		&drone.DroneIdentifier,
		&drone.UserID,
		&drone.Model,
		&drone.SerialNumber,
		&drone.Manufacturer,
		&drone.MaxWeightKg,
		&drone.MaxSpeedKmh,
		&drone.MaxRangeKm,
		&drone.BatteryCapacityMah,
		&drone.Status,
		&drone.BatteryLevelPercent,
		&drone.CurrentLat,
		&drone.CurrentLon,
		&drone.CurrentAltitude,
		&drone.LastLocationUpdateAt,
		&drone.TotalFlightHours,
		&drone.TotalDeliveries,
		&drone.LastMaintenanceAt,
		&drone.NextMaintenanceDueAt,
		&drone.CreatedAt,
		&drone.UpdatedAt,
		&drone.Active,
		&drone.CreatedByID,
		&drone.UpdatedByID,
	)
	if err != nil {
		return nil, err
	}
	return &drone, nil
}

// CreateDrone creates a new drone record
func (r *DronesRepository) CreateDrone(ctx context.Context, drone *domain.CreateDroneRequest) (*domain.Drone, error) {
	var newDrone *domain.Drone
	var err error

	if r.createStmt != nil {
		newDrone, err = r.scanDrone(r.createStmt.QueryRowContext(
			ctx,
			drone.CreatedByID,
			drone.Model,
			drone.SerialNumber,
			drone.Manufacturer,
			drone.BatteryCapacity,
			drone.PayloadCapacity,
			drone.CreatedByID,
		))
	} else {
		query := `
			INSERT INTO drones (
				user_id, model, serial_number, manufacturer, battery_capacity_mah, max_weight_kg, created_by_id
			) VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING
				id, drone_identifier, user_id, model, serial_number, manufacturer,
				max_weight_kg, max_speed_kmh, max_range_km, battery_capacity_mah,
				status, battery_level_percent, current_lat, current_lon, current_altitude,
				last_location_update_at, total_flight_hours, total_deliveries,
				last_maintenance_at, next_maintenance_due_at,
				created_at, updated_at, active, created_by_id, updated_by_id`
		newDrone, err = r.scanDrone(r.db.QueryRowContext(
			ctx,
			query,
			drone.CreatedByID,
			drone.Model,
			drone.SerialNumber,
			drone.Manufacturer,
			drone.BatteryCapacity,
			drone.PayloadCapacity,
			drone.CreatedByID,
		))
	}

	if err != nil {
		r.logger.Error("Failed to create drone", "error", err)
		return nil, err
	}

	return newDrone, nil
}

// UpdateDrone updates an existing drone record
func (r *DronesRepository) UpdateDrone(ctx context.Context, droneID string, req *domain.UpdateDroneRequest) (*domain.Drone, error) {
	var updatedDrone *domain.Drone
	var err error

	if r.updateStmt != nil {
		updatedDrone, err = r.scanDrone(r.updateStmt.QueryRowContext(
			ctx,
			droneID,
			req.Model,
			req.SerialNumber,
			req.Manufacturer,
			req.PayloadCapacity,
			req.BatteryCapacity,
			req.Status,
			req.LastKnownLat,
			req.LastKnownLng,
			req.LastAltitudeM,
			req.LastMaintenanceAt,
			req.NextMaintenanceAt,
			req.UpdatedByID,
		))
	} else {
		query := `
			UPDATE drones SET
				model = COALESCE($2, model),
				serial_number = COALESCE($3, serial_number),
				manufacturer = COALESCE($4, manufacturer),
				max_weight_kg = COALESCE($5, max_weight_kg),
				battery_capacity_mah = COALESCE($6, battery_capacity_mah),
				status = COALESCE($7, status),
				current_lat = COALESCE($8, current_lat),
				current_lon = COALESCE($9, current_lon),
				current_altitude = COALESCE($10, current_altitude),
				last_maintenance_at = COALESCE($11, last_maintenance_at),
				next_maintenance_due_at = COALESCE($12, next_maintenance_due_at),
				updated_by_id = COALESCE($13, updated_by_id),
				updated_at = NOW()
			WHERE id = $1 AND active = TRUE
			RETURNING
				id, drone_identifier, user_id, model, serial_number, manufacturer,
				max_weight_kg, max_speed_kmh, max_range_km, battery_capacity_mah,
				status, battery_level_percent, current_lat, current_lon, current_altitude,
				last_location_update_at, total_flight_hours, total_deliveries,
				last_maintenance_at, next_maintenance_due_at,
				created_at, updated_at, active, created_by_id, updated_by_id`
		updatedDrone, err = r.scanDrone(r.db.QueryRowContext(
			ctx,
			query,
			droneID,
			req.Model,
			req.SerialNumber,
			req.Manufacturer,
			req.PayloadCapacity,
			req.BatteryCapacity,
			req.Status,
			req.LastKnownLat,
			req.LastKnownLng,
			req.LastAltitudeM,
			req.LastMaintenanceAt,
			req.NextMaintenanceAt,
			req.UpdatedByID,
		))
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrDroneNotFound
		}
		r.logger.Error("Failed to update drone", "droneID", droneID, "error", err)
		return nil, err
	}

	return updatedDrone, nil
}

// applyDroneFilters applies drone filters to a query and returns the updated query string and arguments
func (r *DronesRepository) applyDroneFilters(baseQuery string, filter *domain.DroneFilter, startParamCount int) (string, []interface{}, int) {
	query := baseQuery
	args := []interface{}{}
	paramCount := startParamCount

	if filter == nil {
		return query, args, paramCount
	}

	if filter.Status != nil && *filter.Status != "" {
		paramCount++
		query += fmt.Sprintf(" AND status = $%d", paramCount)
		args = append(args, *filter.Status)
	}

	if len(filter.Statuses) > 0 {
		paramCount++
		query += fmt.Sprintf(" AND status = ANY($%d)", paramCount)
		args = append(args, filter.Statuses)
	}

	if filter.Active != nil {
		paramCount++
		query += fmt.Sprintf(" AND active = $%d", paramCount)
		args = append(args, *filter.Active)
	}
	if filter.UserID != nil && *filter.UserID != "" {
		paramCount++
		query += fmt.Sprintf(" AND user_id = $%d", paramCount)
		args = append(args, filter.UserID)
	}

	// Apply location-based filter (with bounding box and Haversine)
	if filter.Lat != nil && filter.Lon != nil && filter.Radius != nil {
		lat := *filter.Lat
		lon := *filter.Lon
		radiusKm := *filter.Radius

		earthRadiusKm := 6371.0
		latDelta := (radiusKm / earthRadiusKm) * (180.0 / math.Pi)
		lonDelta := (radiusKm / (earthRadiusKm * math.Cos(lat*math.Pi/180.0))) * (180.0 / math.Pi)

		minLat := lat - latDelta
		maxLat := lat + latDelta
		minLon := lon - lonDelta
		maxLon := lon + lonDelta

		query += " AND current_lat IS NOT NULL AND current_lon IS NOT NULL"

		paramCount++
		latParam := paramCount
		paramCount++
		lonParam := paramCount
		paramCount++
		minLatParam := paramCount
		paramCount++
		maxLatParam := paramCount
		paramCount++
		minLonParam := paramCount
		paramCount++
		maxLonParam := paramCount
		paramCount++
		radiusParam := paramCount

		query += fmt.Sprintf(" AND current_lat BETWEEN $%d AND $%d", minLatParam, maxLatParam)
		query += fmt.Sprintf(" AND current_lon BETWEEN $%d AND $%d", minLonParam, maxLonParam)
		query += fmt.Sprintf(` AND (
			6371 * acos(
				cos(radians($%d)) * cos(radians(current_lat)) *
				cos(radians(current_lon) - radians($%d)) +
				sin(radians($%d)) * sin(radians(current_lat))
			)
		) <= $%d`, latParam, lonParam, latParam, radiusParam)

		args = append(args, lat, lon, minLat, maxLat, minLon, maxLon, radiusKm)
	}

	return query, args, paramCount
}

// GetDroneByID retrieves a drone by its ID
func (r *DronesRepository) GetDroneByID(ctx context.Context, droneID string) (*domain.Drone, error) {
	var drone *domain.Drone
	var err error

	if r.getByIDStmt != nil {
		drone, err = r.scanDrone(r.getByIDStmt.QueryRowContext(ctx, droneID))
	} else {
		query := `
			SELECT
				id, drone_identifier, user_id, model, serial_number, manufacturer,
				max_weight_kg, max_speed_kmh, max_range_km, battery_capacity_mah,
				status, battery_level_percent, current_lat, current_lon, current_altitude,
				last_location_update_at, total_flight_hours, total_deliveries,
				last_maintenance_at, next_maintenance_due_at,
				created_at, updated_at, active, created_by_id, updated_by_id
			FROM drones
			WHERE id = $1 AND active = TRUE`
		drone, err = r.scanDrone(r.db.QueryRowContext(ctx, query, droneID))
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrDroneNotFound
		}
		return nil, err
	}

	return drone, nil
}

// GetDroneByIdentifier retrieves a drone by its identifier
func (r *DronesRepository) GetDroneByIdentifier(ctx context.Context, identifier string) (*domain.Drone, error) {
	var drone *domain.Drone
	var err error

	if r.getByIdentifierStmt != nil {
		drone, err = r.scanDrone(r.getByIdentifierStmt.QueryRowContext(ctx, identifier))
	} else {
		query := `
			SELECT
				id, drone_identifier, user_id, model, serial_number, manufacturer,
				max_weight_kg, max_speed_kmh, max_range_km, battery_capacity_mah,
				status, battery_level_percent, current_lat, current_lon, current_altitude,
				last_location_update_at, total_flight_hours, total_deliveries,
				last_maintenance_at, next_maintenance_due_at,
				created_at, updated_at, active, created_by_id, updated_by_id
			FROM drones
			WHERE drone_identifier = $1 AND active = TRUE`
		drone, err = r.scanDrone(r.db.QueryRowContext(ctx, query, identifier))
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrDroneNotFound
		}
		return nil, err
	}

	return drone, nil
}

// GetDronesByUserID retrieves all drones for a specific user
func (r *DronesRepository) GetDronesByUserID(ctx context.Context, userID string) ([]*domain.Drone, error) {
	var rows *sql.Rows
	var err error

	if r.getByUserIDStmt != nil {
		rows, err = r.getByUserIDStmt.QueryContext(ctx, userID)
	} else {
		query := `
			SELECT
				id, drone_identifier, user_id, model, serial_number, manufacturer,
				max_weight_kg, max_speed_kmh, max_range_km, battery_capacity_mah,
				status, battery_level_percent, current_lat, current_lon, current_altitude,
				last_location_update_at, total_flight_hours, total_deliveries,
				last_maintenance_at, next_maintenance_due_at,
				created_at, updated_at, active, created_by_id, updated_by_id
			FROM drones
			WHERE user_id = $1 AND active = TRUE
			ORDER BY created_at DESC`
		rows, err = r.db.QueryContext(ctx, query, userID)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var drones []*domain.Drone
	for rows.Next() {
		drone, err := r.scanDrone(rows)
		if err != nil {
			return nil, err
		}
		drones = append(drones, drone)
	}

	return drones, rows.Err()
}

// NearbyDrones retrieves drones near a specific location within a given radius
// Uses Haversine formula for distance calculation
func (r *DronesRepository) NearbyDrones(ctx context.Context, lat, lon, radiusKm float64) ([]*domain.Drone, error) {
	// Calculate bounding box for initial filter (optimization)
	// Earth radius in km
	earthRadiusKm := 6371.0
	latDelta := (radiusKm / earthRadiusKm) * (180.0 / math.Pi)
	lonDelta := (radiusKm / (earthRadiusKm * math.Cos(lat*math.Pi/180.0))) * (180.0 / math.Pi)

	minLat := lat - latDelta
	maxLat := lat + latDelta
	minLon := lon - lonDelta
	maxLon := lon + lonDelta

	// Query with Haversine distance calculation
	query := `
		SELECT
			id, drone_identifier, user_id, model, serial_number, manufacturer,
			max_weight_kg, max_speed_kmh, max_range_km, battery_capacity_mah,
			status, battery_level_percent, current_lat, current_lon, current_altitude,
			last_location_update_at, total_flight_hours, total_deliveries,
			last_maintenance_at, next_maintenance_due_at,
			created_at, updated_at, active, created_by_id, updated_by_id,
			(
				6371 * acos(
					cos(radians($1)) * cos(radians(current_lat)) *
					cos(radians(current_lon) - radians($2)) +
					sin(radians($1)) * sin(radians(current_lat))
				)
			) AS distance
		FROM drones
		WHERE active = TRUE
			AND current_lat IS NOT NULL
			AND current_lon IS NOT NULL
			AND current_lat BETWEEN $3 AND $4
			AND current_lon BETWEEN $5 AND $6
		HAVING distance <= $7
		ORDER BY distance ASC`

	rows, err := r.db.QueryContext(ctx, query, lat, lon, minLat, maxLat, minLon, maxLon, radiusKm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var drones []*domain.Drone
	for rows.Next() {
		var drone domain.Drone
		var distance float64
		err := rows.Scan(
			&drone.ID,
			&drone.DroneIdentifier,
			&drone.UserID,
			&drone.Model,
			&drone.SerialNumber,
			&drone.Manufacturer,
			&drone.MaxWeightKg,
			&drone.MaxSpeedKmh,
			&drone.MaxRangeKm,
			&drone.BatteryCapacityMah,
			&drone.Status,
			&drone.BatteryLevelPercent,
			&drone.CurrentLat,
			&drone.CurrentLon,
			&drone.CurrentAltitude,
			&drone.LastLocationUpdateAt,
			&drone.TotalFlightHours,
			&drone.TotalDeliveries,
			&drone.LastMaintenanceAt,
			&drone.NextMaintenanceDueAt,
			&drone.CreatedAt,
			&drone.UpdatedAt,
			&drone.Active,
			&drone.CreatedByID,
			&drone.UpdatedByID,
			&distance,
		)
		if err != nil {
			return nil, err
		}
		drones = append(drones, &drone)
	}

	return drones, rows.Err()
}

// ListDrones retrieves drones with filtering and pagination
func (r *DronesRepository) ListDrones(ctx context.Context, options domain.PaginationOption[domain.DroneFilter]) (*domain.Pagination[domain.DroneDTO], error) {
	filter := options.Filter
	offset := options.Offset
	limit := options.Limit

	// Build dynamic count query with all available filters
	countQuery, countArgs, _ := r.applyDroneFilters(`SELECT COUNT(*) FROM drones WHERE active = TRUE`, filter, 0)

	// Count total records
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, err
	}

	// Build dynamic query with all available filters
	query, args, paramCount := r.applyDroneFilters(`
		SELECT
			id, drone_identifier, user_id, model, serial_number, manufacturer,
			max_weight_kg, max_speed_kmh, max_range_km, battery_capacity_mah,
			status, battery_level_percent, current_lat, current_lon, current_altitude,
			last_location_update_at, total_flight_hours, total_deliveries,
			last_maintenance_at, next_maintenance_due_at,
			created_at, updated_at, active, created_by_id, updated_by_id
		FROM drones
		WHERE active = TRUE`, filter, 0)

	query += " ORDER BY created_at DESC"

	// Add limit and offset
	paramCount++
	query += fmt.Sprintf(" LIMIT $%d", paramCount)
	args = append(args, limit)

	paramCount++
	query += fmt.Sprintf(" OFFSET $%d", paramCount)
	args = append(args, offset)

	// Get paginated results
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var drones []*domain.DroneDTO
	for rows.Next() {
		drone, err := r.scanDrone(rows)
		if err != nil {
			return nil, err
		}
		drones = append(drones, drone.ToDTO())
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &domain.Pagination[domain.DroneDTO]{
		Data:       drones,
		Total:      total,
		Page:       offset/limit + 1,
		PageSize:   limit,
		TotalPages: (total + limit - 1) / limit,
	}, nil
}

// GetDroneByFilter retrieves a drone by filter criteria
func (r *DronesRepository) GetDroneByFilter(ctx context.Context, filter domain.DroneFilter) (*domain.Drone, error) {
	baseQuery := `
		SELECT
			id, drone_identifier, user_id, model, serial_number, manufacturer,
			max_weight_kg, max_speed_kmh, max_range_km, battery_capacity_mah,
			status, battery_level_percent, current_lat, current_lon, current_altitude,
			last_location_update_at, total_flight_hours, total_deliveries,
			last_maintenance_at, next_maintenance_due_at,
			created_at, updated_at, active, created_by_id, updated_by_id
		FROM drones
		WHERE active = TRUE`

	query, args, _ := r.applyDroneFilters(baseQuery, &filter, 0)
	query += " LIMIT 1"

	drone, err := r.scanDrone(r.db.QueryRowContext(ctx, query, args...))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrDroneNotFound
		}
		return nil, err
	}

	return drone, nil
}

// UpdateStatusBroken updates drone status and if status is broken, updates associated active orders to drone_failed
func (r *DronesRepository) UpdateStatusBroken(ctx context.Context, userID, droneID string, status domain.DroneStatus) (*domain.Drone, error) {
	// Begin transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Error("Failed to begin transaction", "error", err)
		return nil, err
	}
	defer tx.Rollback()

	// Update drone status
	var updatedDrone domain.Drone
	err = tx.QueryRowContext(ctx, `
		UPDATE drones SET
			status = $2,
			updated_by_id = $3,
			updated_at = NOW()
		WHERE id = $1
		RETURNING
			id, drone_identifier, user_id, model, serial_number, manufacturer,
			max_weight_kg, max_speed_kmh, max_range_km, battery_capacity_mah,
			status, battery_level_percent, current_lat, current_lon, current_altitude,
			last_location_update_at, total_flight_hours, total_deliveries,
			last_maintenance_at, next_maintenance_due_at,
			created_at, updated_at, active, created_by_id, updated_by_id`,
		droneID, status, userID).Scan(
		&updatedDrone.ID,
		&updatedDrone.DroneIdentifier,
		&updatedDrone.UserID,
		&updatedDrone.Model,
		&updatedDrone.SerialNumber,
		&updatedDrone.Manufacturer,
		&updatedDrone.MaxWeightKg,
		&updatedDrone.MaxSpeedKmh,
		&updatedDrone.MaxRangeKm,
		&updatedDrone.BatteryCapacityMah,
		&updatedDrone.Status,
		&updatedDrone.BatteryLevelPercent,
		&updatedDrone.CurrentLat,
		&updatedDrone.CurrentLon,
		&updatedDrone.CurrentAltitude,
		&updatedDrone.LastLocationUpdateAt,
		&updatedDrone.TotalFlightHours,
		&updatedDrone.TotalDeliveries,
		&updatedDrone.LastMaintenanceAt,
		&updatedDrone.NextMaintenanceDueAt,
		&updatedDrone.CreatedAt,
		&updatedDrone.UpdatedAt,
		&updatedDrone.Active,
		&updatedDrone.CreatedByID,
		&updatedDrone.UpdatedByID,
	)
	if err != nil {

		if err == sql.ErrNoRows {
			return nil, domain.ErrDroneNotFound
		}
		r.logger.Error("Failed to update drone status", "droneID", droneID, "error", err)
		return nil, err
	}

	// If status is broken, update associated active orders to drone_failed
	if status == "broken" {
		_, err = tx.ExecContext(ctx, `
			UPDATE orders SET
				status = 'handoff',
				drone_id = NULL,
				updated_by_id = $2,
				updated_at = NOW()
			WHERE (drone_id = $1)
			AND status NOT IN ('delivered', 'cancelled', 'handoff')
			AND active = TRUE`,
			droneID, userID)
		if err != nil {
			r.logger.Error("Failed to update orders for broken drone", "droneID", droneID, "error", err)
			return nil, err
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {

		r.logger.Error("Failed to commit transaction", "error", err)
		return nil, err
	}

	return &updatedDrone, nil
}

// Use transaction to ensure data consistency
// ProcessHeartbeat processes a heartbeat from a drone and updates its status and location
// If drone have active order, update order location
func (r *DronesRepository) ProcessHeartbeat(ctx context.Context, droneID string, userId string, req domain.HeartbeatRequest) (*domain.Drone, error) {
	// Begin transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Error("Failed to begin transaction", "error", err)
		return nil, err
	}
	defer tx.Rollback()

	// Update drone location and battery level
	var updatedDrone domain.Drone
	err = tx.QueryRowContext(ctx, `
		UPDATE drones SET
			current_lat = $2,
			current_lon = $3,
			current_altitude = $4,
			battery_level_percent = $5,
			last_location_update_at = NOW(),
			updated_by_id = $6,
			updated_at = NOW()
		WHERE id = $1 AND active = TRUE
		RETURNING
			id, drone_identifier, user_id, model, serial_number, manufacturer,
			max_weight_kg, max_speed_kmh, max_range_km, battery_capacity_mah,
			status, battery_level_percent, current_lat, current_lon, current_altitude,
			last_location_update_at, total_flight_hours, total_deliveries,
			last_maintenance_at, next_maintenance_due_at,
			created_at, updated_at, active, created_by_id, updated_by_id`,
		droneID,
		req.Latitude,
		req.Longitude,
		req.Altitude,
		req.Battery,
		userId,
	).Scan(
		&updatedDrone.ID,
		&updatedDrone.DroneIdentifier,
		&updatedDrone.UserID,
		&updatedDrone.Model,
		&updatedDrone.SerialNumber,
		&updatedDrone.Manufacturer,
		&updatedDrone.MaxWeightKg,
		&updatedDrone.MaxSpeedKmh,
		&updatedDrone.MaxRangeKm,
		&updatedDrone.BatteryCapacityMah,
		&updatedDrone.Status,
		&updatedDrone.BatteryLevelPercent,
		&updatedDrone.CurrentLat,
		&updatedDrone.CurrentLon,
		&updatedDrone.CurrentAltitude,
		&updatedDrone.LastLocationUpdateAt,
		&updatedDrone.TotalFlightHours,
		&updatedDrone.TotalDeliveries,
		&updatedDrone.LastMaintenanceAt,
		&updatedDrone.NextMaintenanceDueAt,
		&updatedDrone.CreatedAt,
		&updatedDrone.UpdatedAt,
		&updatedDrone.Active,
		&updatedDrone.CreatedByID,
		&updatedDrone.UpdatedByID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrDroneNotFound
		}
		r.logger.Error("Failed to update drone in heartbeat", "droneID", droneID, "error", err)
		return nil, err
	}

	// Update location of active order assigned to this drone
	_, err = tx.ExecContext(ctx, `
		UPDATE orders SET
			current_lat = $1,
			current_lon = $2,
			updated_at = NOW()
		WHERE drone_id = $3
		AND status IN ('picked_up', 'in_transit', 'arrived', 'handoff', 'reassigned')
		AND active = TRUE`,
		req.Latitude,
		req.Longitude,
		droneID,
	)
	if err != nil {
		r.logger.Error("Failed to update order location in heartbeat", "droneID", droneID, "error", err)
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		r.logger.Error("Failed to commit heartbeat transaction", "error", err)
		return nil, err
	}

	return &updatedDrone, nil
}
