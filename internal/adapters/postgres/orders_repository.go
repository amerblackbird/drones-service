package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"drones/internal/core/domain"
	"drones/internal/ports"
)

type OrdersRepositoryImpl struct {
	db                    *sql.DB
	logger                ports.Logger
	createStmt            *sql.Stmt
	updateStmt            *sql.Stmt
	deleteStmt            *sql.Stmt
	getByOrderNumberStmt  *sql.Stmt
	updateStatusStmt      *sql.Stmt
	getByUserIDStmt       *sql.Stmt
	getOrdersByStatusStmt *sql.Stmt
	countOrdersStmt       *sql.Stmt
}

func NewOrdersRepository(db *sql.DB, logger ports.Logger) ports.OrdersRepository {
	repo := &OrdersRepositoryImpl{
		db:     db,
		logger: logger,
	}
	// Prepare statements
	if err := repo.prepareStatements(); err != nil {
		logger.Error("Orders: Failed to prepare statements", "error", err)
		return repo // Return anyway, fallback to non-prepared queries
	}

	return repo
}

func (r *OrdersRepositoryImpl) prepareStatements() error {
	var err error

	// Create statement with comprehensive fields
	r.createStmt, err = r.db.Prepare(`
		INSERT INTO orders (
			user_id, receiver_name, receiver_phone, delivery_note,
			package_weight_kg, origin_address, origin_lat, origin_lon, destination_address,
			destination_lat, destination_lon, scheduled_at, created_by_id
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
	) RETURNING
		id, order_number, user_id, receiver_name, receiver_phone, delivery_note,
		package_weight_kg, origin_address, origin_lat, origin_lon, destination_address,
		destination_lat, destination_lon, status, scheduled_at, delivered_at, cancelled_at,
		delivered_by_drone_id, drone_id, withdrawn_at, current_lat, current_lon, current_altitude,
		last_location_update_at, estimated_arrival_at, created_at, updated_at, active`)
	if err != nil {
		return err
	}

	// Update statement with comprehensive fields
	r.updateStmt, err = r.db.Prepare(`
		UPDATE orders SET
			receiver_name = COALESCE($2, receiver_name),
			receiver_phone = COALESCE($3, receiver_phone),
			delivery_note = COALESCE($4, delivery_note),
			package_weight_kg = COALESCE($5, package_weight_kg),
			origin_address = COALESCE($6, origin_address),
			origin_lat = COALESCE($7, origin_lat),
			origin_lon = COALESCE($8, origin_lon),
			destination_address = COALESCE($9, destination_address),
			destination_lat = COALESCE($10, destination_lat),
			destination_lon = COALESCE($11, destination_lon),
			status = COALESCE($12, status),
			scheduled_at = COALESCE($13, scheduled_at),
			delivered_at = COALESCE($14, delivered_at),
			cancelled_at = COALESCE($15, cancelled_at),
			delivered_by_drone_id = COALESCE($16, delivered_by_drone_id),
			updated_by_id = COALESCE($17, updated_by_id),
			withdrawn_at = COALESCE($18, withdrawn_at),
			current_lat = COALESCE($19, current_lat),
			current_lon = COALESCE($20, current_lon),
			current_altitude = COALESCE($21, current_altitude),
			last_location_update_at = COALESCE($22, last_location_update_at),
			estimated_arrival_at = COALESCE($23, estimated_arrival_at),
			drone_id = COALESCE($24, drone_id),
			updated_at = NOW()
		WHERE id = $1 AND active = TRUE
	RETURNING
		id, order_number, user_id, receiver_name, receiver_phone, delivery_note,
		package_weight_kg, origin_address, origin_lat, origin_lon, destination_address,
		destination_lat, destination_lon, status, scheduled_at, delivered_at, cancelled_at,
		delivered_by_drone_id, drone_id, withdrawn_at, current_lat, current_lon, current_altitude,
		last_location_update_at, estimated_arrival_at, created_at, updated_at, active`)
	if err != nil {
		return err
	}

	// Permanent delete statement
	r.deleteStmt, err = r.db.Prepare(`DELETE FROM orders WHERE id = $1`)
	if err != nil {
		return err
	}

	// Get by order number
	r.getByOrderNumberStmt, err = r.db.Prepare(`
		SELECT
			id, order_number, user_id, receiver_name, receiver_phone, delivery_note,
			package_weight_kg, origin_address, origin_lat, origin_lon, destination_address,
			destination_lat, destination_lon, status, scheduled_at, delivered_at, cancelled_at,
			delivered_by_drone_id,drone_id , withdrawn_at, current_lat, current_lon, current_altitude,
			last_location_update_at, estimated_arrival_at, created_at, updated_at, active
		FROM orders
		WHERE order_number = $1 AND active = TRUE`)
	if err != nil {
		return err
	}

	// Update status only
	r.updateStatusStmt, err = r.db.Prepare(`
		UPDATE orders SET
			status = $2,
			updated_by_id = $3,
			updated_at = NOW()
	WHERE id = $1 AND active = TRUE
	RETURNING
		id, order_number, user_id, receiver_name, receiver_phone, delivery_note,
		package_weight_kg, origin_address, origin_lat, origin_lon, destination_address,
		destination_lat, destination_lon, status, scheduled_at, delivered_at, cancelled_at,
		delivered_by_drone_id,drone_id, withdrawn_at, current_lat, current_lon, current_altitude,
		last_location_update_at, estimated_arrival_at, created_at, updated_at, active`)
	if err != nil {
		return err
	}

	// Get orders by user ID
	r.getByUserIDStmt, err = r.db.Prepare(`
		SELECT
			id, order_number, user_id, receiver_name, receiver_phone, delivery_note,
			package_weight_kg, origin_address, origin_lat, origin_lon, destination_address,
			destination_lat, destination_lon, status, scheduled_at, delivered_at, cancelled_at,
			delivered_by_drone_id, drone_id, withdrawn_at, current_lat, current_lon, current_altitude,
			last_location_update_at, estimated_arrival_at, created_at, updated_at, active
		FROM orders
		WHERE user_id = $1 AND active = TRUE
		ORDER BY created_at DESC`)
	if err != nil {
		return err
	}

	// Get orders by status
	r.getOrdersByStatusStmt, err = r.db.Prepare(`
		SELECT
			id, order_number, user_id, receiver_name, receiver_phone, delivery_note,
			package_weight_kg, origin_address, origin_lat, origin_lon, destination_address,
			destination_lat, destination_lon, status, scheduled_at, delivered_at, cancelled_at,
			delivered_by_drone_id,drone_id, withdrawn_at, current_lat, current_lon, current_altitude,
			last_location_update_at, estimated_arrival_at, created_at, updated_at, active
		FROM orders
		WHERE active = TRUE AND status = $1
		ORDER BY updated_at DESC`)
	if err != nil {
		return err
	}

	// Count orders
	r.countOrdersStmt, err = r.db.Prepare(`
		SELECT COUNT(*)
		FROM orders
		WHERE active = TRUE
		AND ($1::VARCHAR IS NULL OR status = $1)
		AND ($2::UUID IS NULL OR user_id = $2)`)
	if err != nil {
		return err
	}

	return nil
}

func (r *OrdersRepositoryImpl) Close() error {
	if r.createStmt != nil {
		if err := r.createStmt.Close(); err != nil {
			r.logger.Error("Failed to close createStmt", "error", err)
			return err
		}
	}
	if r.updateStmt != nil {
		if err := r.updateStmt.Close(); err != nil {
			r.logger.Error("Failed to close updateStmt", "error", err)
			return err
		}
	}
	if r.deleteStmt != nil {
		if err := r.deleteStmt.Close(); err != nil {
			r.logger.Error("Failed to close deleteStmt", "error", err)
			return err
		}
	}
	if r.getByOrderNumberStmt != nil {
		if err := r.getByOrderNumberStmt.Close(); err != nil {
			r.logger.Error("Failed to close getByOrderNumberStmt", "error", err)
			return err
		}
	}
	if r.updateStatusStmt != nil {
		if err := r.updateStatusStmt.Close(); err != nil {
			r.logger.Error("Failed to close updateStatusStmt", "error", err)
			return err
		}
	}
	if r.getByUserIDStmt != nil {
		if err := r.getByUserIDStmt.Close(); err != nil {
			r.logger.Error("Failed to close getByUserIDStmt", "error", err)
			return err
		}
	}
	if r.getOrdersByStatusStmt != nil {
		if err := r.getOrdersByStatusStmt.Close(); err != nil {
			r.logger.Error("Failed to close getOrdersByStatusStmt", "error", err)
			return err
		}
	}
	if r.countOrdersStmt != nil {
		if err := r.countOrdersStmt.Close(); err != nil {
			r.logger.Error("Failed to close countOrdersStmt", "error", err)
			return err
		}
	}
	return nil
}

func (r *OrdersRepositoryImpl) GetDB() *sql.DB {
	return r.db
}

// scanOrder scans a row into an Order struct
func (r *OrdersRepositoryImpl) scanOrder(scanner interface {
	Scan(dest ...interface{}) error
}) (*domain.Order, error) {
	var order domain.Order
	err := scanner.Scan(
		&order.ID,
		&order.OrderNumber,
		&order.UserID,
		&order.ReceiverName,
		&order.ReceiverPhone,
		&order.DeliveryNote,
		&order.PackageWeightKg,
		&order.OriginAddress,
		&order.OriginLat,
		&order.OriginLon,
		&order.DestinationAddress,
		&order.DestinationLat,
		&order.DestinationLon,
		&order.Status,
		&order.ScheduledAt,
		&order.DeliveredAt,
		&order.CancelledAt,
		&order.DeliveredByDroneID,
		&order.DroneID,
		&order.WithdrawnAt,
		&order.CurrentLat,
		&order.CurrentLon,
		&order.CurrentAltitude,
		&order.LastLocationUpdateAt,
		&order.EstimatedArrivalAt,
		&order.CreatedAt,
		&order.UpdatedAt,
		&order.Active,
	)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// CreateOrder creates a new order in the repository
func (r *OrdersRepositoryImpl) CreateOrder(ctx context.Context, userID string, createOrder *domain.CreateOrderRequest) (*domain.Order, error) {
	var order *domain.Order
	var err error

	if r.createStmt != nil {
		order, err = r.scanOrder(r.createStmt.QueryRowContext(ctx,
			userID,
			createOrder.ReceiverName,
			createOrder.ReceiverPhone,
			createOrder.DeliveryNote,
			createOrder.PackageWeightKg,
			createOrder.OriginAddress,
			createOrder.OriginLat,
			createOrder.OriginLon,
			createOrder.DestinationAddress,
			createOrder.DestinationLat,
			createOrder.DestinationLon,
			createOrder.ScheduledAt,
			userID,
		))
	} else {
		order, err = r.scanOrder(r.db.QueryRowContext(ctx, `
			INSERT INTO orders (
				user_id, receiver_name, receiver_phone, delivery_note,
				package_weight_kg, origin_address, origin_lat, origin_lon, destination_address,
				destination_lat, destination_lon, scheduled_at, created_by_id
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
			RETURNING
				id, order_number, user_id, receiver_name, receiver_phone, delivery_note,
				package_weight_kg, origin_address, origin_lat, origin_lon, destination_address,
				destination_lat, destination_lon, status, scheduled_at, delivered_at, cancelled_at,
				delivered_by_drone_id, drone_id, withdrawn_at, current_lat, current_lon, current_altitude,
				last_location_update_at, estimated_arrival_at, created_at, updated_at, active`,
			userID,
			createOrder.ReceiverName,
			createOrder.ReceiverPhone,
			createOrder.DeliveryNote,
			createOrder.PackageWeightKg,
			createOrder.OriginAddress,
			createOrder.OriginLat,
			createOrder.OriginLon,
			createOrder.DestinationAddress,
			createOrder.DestinationLat,
			createOrder.DestinationLon,
			createOrder.ScheduledAt,
			userID,
		))
	}

	if err != nil {
		r.logger.Error("Failed to create order", "error", err)
		return nil, err
	}

	return order, nil
}

// GetOrderByID retrieves an order by its ID
func (r *OrdersRepositoryImpl) GetOrderByID(ctx context.Context, orderID string, options domain.OrderFilter) (*domain.Order, error) {
	// Build base query with order ID
	baseQuery := `
		SELECT
			id, order_number, user_id, receiver_name, receiver_phone, delivery_note,
			package_weight_kg, origin_address, origin_lat, origin_lon, destination_address,
			destination_lat, destination_lon, status, scheduled_at, delivered_at, cancelled_at,
			delivered_by_drone_id, drone_id, withdrawn_at, current_lat, current_lon, current_altitude,
			last_location_update_at, estimated_arrival_at, created_at, updated_at, active
		FROM orders
		WHERE id = $1 AND active = TRUE`

	// Apply additional filters
	query, filterArgs, _ := r.applyOrderFilters(baseQuery, &options, 1)
	args := append([]interface{}{orderID}, filterArgs...)

	r.logger.Info("query", "query", query, "OrderId", orderID, "DroneID", options.DroneID)

	order, err := r.scanOrder(r.db.QueryRowContext(ctx, query, args...))

	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn("Order not found", "orderID", orderID, "filters", options)
			return nil, domain.ErrOrderNotFound
		}
		r.logger.Error("Failed to get order by ID", "orderID", orderID, "filters", options, "error", err)
		return nil, err
	}

	return order, nil
}

// UpdateOrder updates an existing order in the repository
func (r *OrdersRepositoryImpl) UpdateOrder(ctx context.Context, orderID string, update *domain.UpdateOrderRequest) (*domain.Order, error) {
	var err error

	r.logger.Info("Updating order", "orderID", orderID, "update", update)

	var order *domain.Order
	if r.updateStmt != nil {
		order, err = r.scanOrder(r.updateStmt.QueryRowContext(ctx,
			orderID,
			update.ReceiverName,
			update.ReceiverPhone,
			update.DeliveryNote,
			update.PackageWeightKg,
			update.OriginAddress,
			update.OriginLat,
			update.OriginLon,
			update.DestinationAddress,
			update.DestinationLat,
			update.DestinationLon,
			update.Status,
			update.ScheduledAt,
			update.DeliveredAt,
			update.CancelledAt,
			update.DeliveredByDroneID,
			update.UpdatedByID,
			update.WithdrawnAt,
			update.CurrentLat,
			update.CurrentLon,
			update.CurrentAltitude,
			update.LastLocationUpdateAt,
			update.EstimatedArrivalAt,
			update.DroneID,
		))
	} else {
		// Fallback to regular query with COALESCE for null handling
		order, err = r.scanOrder(r.db.QueryRowContext(ctx, `
			UPDATE orders SET
				receiver_name = COALESCE($2, receiver_name),
				receiver_phone = COALESCE($3, receiver_phone),
				delivery_note = COALESCE($4, delivery_note),
				package_weight_kg = COALESCE($5, package_weight_kg),
				origin_address = COALESCE($6, origin_address),
				origin_lat = COALESCE($7, origin_lat),
				origin_lon = COALESCE($8, origin_lon),
				destination_address = COALESCE($9, destination_address),
				destination_lat = COALESCE($10, destination_lat),
				destination_lon = COALESCE($11, destination_lon),
				status = COALESCE($12, status),
				scheduled_at = COALESCE($13, scheduled_at),
				delivered_at = COALESCE($14, delivered_at),
				cancelled_at = COALESCE($15, cancelled_at),
				delivered_by_drone_id = COALESCE($16, delivered_by_drone_id),
				updated_by_id = COALESCE($17, updated_by_id),
				withdrawn_at = COALESCE($18, withdrawn_at),
				current_lat = COALESCE($19, current_lat),
				current_lon = COALESCE($20, current_lon),
				current_altitude = COALESCE($21, current_altitude),
				last_location_update_at = COALESCE($22, last_location_update_at),
				estimated_arrival_at = COALESCE($23, estimated_arrival_at),
				drone_id = COALESCE($24, drone_id),
				updated_at = NOW()
			WHERE id = $1 AND active = TRUE
			RETURNING
				id, order_number, user_id, receiver_name, receiver_phone, delivery_note,
				package_weight_kg, origin_address, origin_lat, origin_lon, destination_address,
				destination_lat, destination_lon, status, scheduled_at, delivered_at, cancelled_at,
				delivered_by_drone_id, drone_id, withdrawn_at, current_lat, current_lon, current_altitude,
				last_location_update_at, estimated_arrival_at, created_at, updated_at, active`,
			orderID,
			update.ReceiverName,
			update.ReceiverPhone,
			update.DeliveryNote,
			update.PackageWeightKg,
			update.OriginAddress,
			update.OriginLat,
			update.OriginLon,
			update.DestinationAddress,
			update.DestinationLat,
			update.DestinationLon,
			update.Status,
			update.ScheduledAt,
			update.DeliveredAt,
			update.CancelledAt,
			update.DeliveredByDroneID,
			update.UpdatedByID,
			update.WithdrawnAt,
			update.CurrentLat,
			update.CurrentLon,
			update.CurrentAltitude,
			update.LastLocationUpdateAt,
			update.EstimatedArrivalAt,
			update.DroneID,
		))
	}

	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn("Order not found for update", "orderID", orderID)
			return nil, domain.ErrOrderNotFound
		}
		r.logger.Error("Failed to update order", "orderID", orderID, "error", err)
		return nil, err
	}

	return order, nil
}

// DeleteOrder permanently deletes an order from the repository
func (r *OrdersRepositoryImpl) DeleteOrder(ctx context.Context, orderID string) error {
	var result sql.Result
	var err error

	if r.deleteStmt != nil {
		result, err = r.deleteStmt.ExecContext(ctx, orderID)
	} else {
		result, err = r.db.ExecContext(ctx, `DELETE FROM orders WHERE id = $1`, orderID)
	}

	if err != nil {
		r.logger.Error("Failed to delete order", "orderID", orderID, "error", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("Failed to get rows affected", "orderID", orderID, "error", err)
		return err
	}

	if rowsAffected == 0 {
		r.logger.Warn("Order not found for deletion", "orderID", orderID)
		return domain.ErrOrderNotFound
	}

	return nil
}

// applyOrderFilters applies order filters to a query and returns the updated query string and arguments
func (r *OrdersRepositoryImpl) applyOrderFilters(baseQuery string, filter *domain.OrderFilter, startParamCount int) (string, []interface{}, int) {
	query := baseQuery
	args := []interface{}{}
	paramCount := startParamCount

	if filter == nil {
		return query, args, paramCount
	}

	if filter.UserID != nil && *filter.UserID != "" {
		paramCount++
		query += fmt.Sprintf(" AND user_id = $%d", paramCount)
		args = append(args, *filter.UserID)
	}

	if filter.Status != nil && *filter.Status != "" {
		paramCount++
		query += fmt.Sprintf(" AND status = $%d", paramCount)
		args = append(args, *filter.Status)
	}

	if filter.DroneID != nil && *filter.DroneID != "" && filter.DeliveredByDroneID != nil && *filter.DeliveredByDroneID != "" && filter.DeliveredByDroneID == filter.DroneID {
		// Combaine two statements with ord
		paramCount++
		query += fmt.Sprintf(" AND (drone_id = $%d OR delivered_by_drone_id = $%d)", paramCount, paramCount)
		args = append(args, *filter.DroneID)

	} else {
		if filter.DroneID != nil && *filter.DroneID != "" {
			paramCount++
			query += fmt.Sprintf(" AND drone_id = $%d", paramCount)
			args = append(args, *filter.DroneID)
		}

		if filter.DeliveredByDroneID != nil && *filter.DeliveredByDroneID != "" {
			paramCount++
			query += fmt.Sprintf(" AND delivered_by_drone_id = $%d", paramCount)
			args = append(args, *filter.DeliveredByDroneID)
		}
	}

	if filter.ReceiverName != nil && *filter.ReceiverName != "" {
		paramCount++
		query += fmt.Sprintf(" AND LOWER(receiver_name) LIKE LOWER($%d)", paramCount)
		args = append(args, "%"+*filter.ReceiverName+"%")
	}

	if filter.ReceiverPhone != nil && *filter.ReceiverPhone != "" {
		paramCount++
		query += fmt.Sprintf(" AND receiver_phone = $%d", paramCount)
		args = append(args, *filter.ReceiverPhone)
	}

	if filter.OriginAddress != nil && *filter.OriginAddress != "" {
		paramCount++
		query += fmt.Sprintf(" AND LOWER(origin_address) LIKE LOWER($%d)", paramCount)
		args = append(args, "%"+*filter.OriginAddress+"%")
	}

	if filter.DestinationAddress != nil && *filter.DestinationAddress != "" {
		paramCount++
		query += fmt.Sprintf(" AND LOWER(destination_address) LIKE LOWER($%d)", paramCount)
		args = append(args, "%"+*filter.DestinationAddress+"%")
	}

	if filter.CreatedAtFrom != nil {
		paramCount++
		query += fmt.Sprintf(" AND created_at >= $%d", paramCount)
		args = append(args, *filter.CreatedAtFrom)
	}

	if filter.CreatedAtTo != nil {
		paramCount++
		query += fmt.Sprintf(" AND created_at <= $%d", paramCount)
		args = append(args, *filter.CreatedAtTo)
	}

	if filter.ScheduledAtFrom != nil {
		paramCount++
		query += fmt.Sprintf(" AND scheduled_at >= $%d", paramCount)
		args = append(args, *filter.ScheduledAtFrom)
	}

	if filter.ScheduledAtTo != nil {
		paramCount++
		query += fmt.Sprintf(" AND scheduled_at <= $%d", paramCount)
		args = append(args, *filter.ScheduledAtTo)
	}

	if filter.MinWeight != nil {
		paramCount++
		query += fmt.Sprintf(" AND package_weight_kg >= $%d", paramCount)
		args = append(args, *filter.MinWeight)
	}

	if filter.MaxWeight != nil {
		paramCount++
		query += fmt.Sprintf(" AND package_weight_kg <= $%d", paramCount)
		args = append(args, *filter.MaxWeight)
	}

	return query, args, paramCount
}

// ListOrders retrieves orders with filtering and pagination
func (r *OrdersRepositoryImpl) ListOrders(ctx context.Context, options domain.PaginationOption[domain.OrderFilter]) (*domain.Pagination[domain.OrderDTO], error) {
	filter := options.Filter
	offset := options.Offset
	limit := options.Limit

	// Build dynamic count query with all available filters
	countQuery, countArgs, _ := r.applyOrderFilters(`SELECT COUNT(*) FROM orders WHERE active = TRUE`, filter, 0)

	// Count total records
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		r.logger.Error("Failed to count orders", "filter", filter, "error", err)
		return nil, err
	}

	// Build dynamic query with all available filters
	query, args, paramCount := r.applyOrderFilters(`
		SELECT
			id, order_number, user_id, receiver_name, receiver_phone, delivery_note,
			package_weight_kg, origin_address, origin_lat, origin_lon, destination_address,
			destination_lat, destination_lon, status, scheduled_at, delivered_at, cancelled_at,
			delivered_by_drone_id, drone_id, withdrawn_at, current_lat, current_lon, current_altitude,
			last_location_update_at, estimated_arrival_at, created_at, updated_at, active
		FROM orders
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
		r.logger.Error("Failed to list orders", "filter", filter, "limit", limit, "offset", offset, "error", err)
		return nil, err
	}
	defer rows.Close()

	var orders []*domain.OrderDTO
	for rows.Next() {
		order, err := r.scanOrder(rows)
		if err != nil {
			r.logger.Error("Failed to scan order row", "error", err)
			continue
		}
		orders = append(orders, order.ToDTO())
	}

	if err = rows.Err(); err != nil {
		r.logger.Error("Failed to iterate order rows", "error", err)
		return nil, err
	}

	return &domain.Pagination[domain.OrderDTO]{
		Data:       orders,
		Total:      total,
		Page:       offset/limit + 1,
		PageSize:   limit,
		TotalPages: (total + limit - 1) / limit,
	}, nil
} // Additional specialized methods

// GetByOrderNumber retrieves an order by its order number
func (r *OrdersRepositoryImpl) GetByOrderNumber(ctx context.Context, orderNumber string) (*domain.Order, error) {
	var order *domain.Order
	var err error

	if r.getByOrderNumberStmt != nil {
		order, err = r.scanOrder(r.getByOrderNumberStmt.QueryRowContext(ctx, orderNumber))
	} else {
		order, err = r.scanOrder(r.db.QueryRowContext(ctx, `
			SELECT
				id, order_number, user_id, receiver_name, receiver_phone, delivery_note,
				package_weight_kg, origin_address, origin_lat, origin_lon, destination_address,
				destination_lat, destination_lon, status, scheduled_at, delivered_at, cancelled_at,
				delivered_by_drone_id, drone_id, withdrawn_at, current_lat, current_lon, current_altitude,
				last_location_update_at, estimated_arrival_at, created_at, updated_at, active
			FROM orders
			WHERE order_number = $1 AND active = TRUE`, orderNumber))
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrOrderNotFound
		}
		r.logger.Error("Failed to get order by order number", "orderNumber", orderNumber, "error", err)
		return nil, err
	}

	return order, nil
}

// UpdateStatus updates only the status of an order
func (r *OrdersRepositoryImpl) UpdateStatus(ctx context.Context, orderID string, status string, updatedByID string) (*domain.Order, error) {
	var order *domain.Order
	var err error

	if r.updateStatusStmt != nil {
		order, err = r.scanOrder(r.updateStatusStmt.QueryRowContext(ctx, orderID, status, updatedByID))
	} else {
		order, err = r.scanOrder(r.db.QueryRowContext(ctx, `
			UPDATE orders SET
				status = $2,
				updated_by_id = $3,
				updated_at = NOW()
			WHERE id = $1 AND active = TRUE
			RETURNING
				id, order_number, user_id, receiver_name, receiver_phone, delivery_note,
				package_weight_kg, origin_address, origin_lat, origin_lon, destination_address,
				destination_lat, destination_lon, status, scheduled_at, delivered_at, cancelled_at,
				delivered_by_drone_id, drone_id, withdrawn_at, current_lat, current_lon, current_altitude,
				last_location_update_at, estimated_arrival_at, created_at, updated_at, active`, orderID, status, updatedByID))
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrOrderNotFound
		}
		r.logger.Error("Failed to update order status", "orderID", orderID, "error", err)
		return nil, err
	}

	return order, nil
}

// GetByUserID retrieves all orders for a specific user
func (r *OrdersRepositoryImpl) GetByUserID(ctx context.Context, userID string) ([]*domain.Order, error) {
	var rows *sql.Rows
	var err error

	if r.getByUserIDStmt != nil {
		rows, err = r.getByUserIDStmt.QueryContext(ctx, userID)
	} else {
		rows, err = r.db.QueryContext(ctx, `
			SELECT
				id, order_number, user_id, receiver_name, receiver_phone, delivery_note,
				package_weight_kg, origin_address, origin_lat, origin_lon, destination_address,
				destination_lat, destination_lon, status, scheduled_at, delivered_at, cancelled_at,
				delivered_by_drone_id, drone_id, withdrawn_at, current_lat, current_lon, current_altitude,
				last_location_update_at, estimated_arrival_at, created_at, updated_at, active
			FROM orders
			WHERE user_id = $1 AND active = TRUE
			ORDER BY created_at DESC`, userID)
	}
	if err != nil {
		r.logger.Error("Failed to get orders by user ID", "userID", userID, "error", err)
		return nil, err
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		order, err := r.scanOrder(rows)
		if err != nil {
			r.logger.Error("Failed to scan order by user ID row", "error", err)
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, rows.Err()
}

// GetOrdersByStatus retrieves all orders with a specific status
func (r *OrdersRepositoryImpl) GetOrdersByStatus(ctx context.Context, status string) ([]*domain.Order, error) {
	var rows *sql.Rows
	var err error

	if r.getOrdersByStatusStmt != nil {
		rows, err = r.getOrdersByStatusStmt.QueryContext(ctx, status)
	} else {
		rows, err = r.db.QueryContext(ctx, `
			SELECT
				id, order_number, user_id, receiver_name, receiver_phone, delivery_note,
				package_weight_kg, origin_address, origin_lat, origin_lon, destination_address,
				destination_lat, destination_lon, status, scheduled_at, delivered_at, cancelled_at,
				delivered_by_drone_id, drone_id, withdrawn_at, current_lat, current_lon, current_altitude,
				last_location_update_at, estimated_arrival_at, created_at, updated_at, active
			FROM orders
			WHERE active = TRUE AND status = $1
			ORDER BY updated_at DESC`, status)
	}
	if err != nil {
		r.logger.Error("Failed to get orders by status", "status", status, "error", err)
		return nil, err
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		order, err := r.scanOrder(rows)
		if err != nil {
			r.logger.Error("Failed to scan order by status row", "error", err)
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, rows.Err()
}

// CountOrders counts the total number of orders with optional filters
func (r *OrdersRepositoryImpl) CountOrders(ctx context.Context, status *string, userID *string) (int64, error) {
	var count int64
	var statusParam interface{}
	if status != nil {
		statusParam = *status
	}

	var userIDParam interface{}
	if userID != nil {
		userIDParam = *userID
	}

	var err error
	if r.countOrdersStmt != nil {
		err = r.countOrdersStmt.QueryRowContext(ctx, statusParam, userIDParam).Scan(&count)
	} else {
		err = r.db.QueryRowContext(ctx, `
			SELECT COUNT(*)
			FROM orders
			WHERE active = TRUE
			AND ($1::VARCHAR IS NULL OR status = $1)
			AND ($2::UUID IS NULL OR user_id = $2)`, statusParam, userIDParam).Scan(&count)
	}
	if err != nil {
		r.logger.Error("Failed to count orders", "error", err)
		return 0, err
	}

	return count, nil
}

func (r *OrdersRepositoryImpl) GetOrderByFilter(ctx context.Context, options domain.OrderFilter) (*domain.Order, error) {
	baseQuery := `
		SELECT
			id, order_number, user_id, receiver_name, receiver_phone, delivery_note,
			package_weight_kg, origin_address, origin_lat, origin_lon, destination_address,
			destination_lat, destination_lon, status, scheduled_at, delivered_at, cancelled_at,
			delivered_by_drone_id, drone_id, withdrawn_at, current_lat, current_lon, current_altitude,
			last_location_update_at, estimated_arrival_at, created_at, updated_at, active
		FROM orders
		WHERE active = TRUE`

	query, args, _ := r.applyOrderFilters(baseQuery, &options, 0)
	query += " LIMIT 1"

	order, err := r.scanOrder(r.db.QueryRowContext(ctx, query, args...))

	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn("Order not found with given filters", "filters", options)
			return nil, domain.ErrOrderNotFound
		}
		r.logger.Error("Failed to get order by filters", "filters", options, "error", err)
		return nil, err
	}

	return order, nil
}

// Use transactions to ensure atomicity
// update order status
// if status is 'reserved' mark drone as loading
// if status is 'in_transit' mark drone as delivering
// if status is 'delivered' mark drone as returning
// if status is 'drone_failed' mark drone as broken
func (r *OrdersRepositoryImpl) UpdateOrderStatus(ctx context.Context, orderID string, options domain.UpdateStatusRequest) (*domain.Order, error) {
	// Begin transaction
	droneID := options.DroneID
	status := options.Status
	updatedByID := options.UpdatedByID
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Error("Failed to begin transaction", "error", err)
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Update order status
	order, err := r.scanOrder(tx.QueryRowContext(ctx, `
		UPDATE orders SET
			status = $2::VARCHAR,
			updated_by_id = $3,
			drone_id = CASE 
				WHEN $2::VARCHAR IN ('delivered') THEN NULL
				WHEN $2::VARCHAR IN ('reserved') AND $4::UUID IS NOT NULL THEN $4
				ELSE COALESCE($4, drone_id)
			END,
			delivered_by_drone_id = CASE
				WHEN $2::VARCHAR = 'delivered' THEN COALESCE($4, delivered_by_drone_id)
				ELSE delivered_by_drone_id
			END,
			updated_at = NOW()
		WHERE id = $1 AND active = TRUE
		RETURNING
			id, order_number, user_id, receiver_name, receiver_phone, delivery_note,
			package_weight_kg, origin_address, origin_lat, origin_lon, destination_address,
			destination_lat, destination_lon, status, scheduled_at, delivered_at, cancelled_at,
			delivered_by_drone_id, drone_id, withdrawn_at, current_lat, current_lon, current_altitude,
			last_location_update_at, estimated_arrival_at, created_at, updated_at, active`,
		orderID, status, updatedByID, droneID))

	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn("Order not found for status update", "orderID", orderID)
			return nil, domain.ErrOrderNotFound
		}
		r.logger.Error("Failed to update order status", "orderID", orderID, "error", err)
		return nil, err
	}

	// Determine drone status based on order status
	var droneStatus domain.DroneStatus
	shouldUpdateDrone := true

	switch status {
	case domain.OrderStatusReserved:
		droneStatus = domain.DroneStatusLoading
	case domain.OrderStatusInTransit:
		droneStatus = domain.DroneStatusDelivering
	case domain.OrderStatusDelivered:
		droneStatus = domain.DroneStatusReturning
	case domain.OrderStatusFailed:
		droneStatus = domain.DroneStatusReturning
	default:
		shouldUpdateDrone = false
	}

	// Update drone status if needed
	if shouldUpdateDrone && droneID != "" {
		_, err = tx.ExecContext(ctx, `
			UPDATE drones SET
				status = $2,
				updated_by_id = $3,
				updated_at = NOW()
			WHERE id = $1 AND active = TRUE`,
			droneID, droneStatus, updatedByID)

		if err != nil {
			r.logger.Error("Failed to update drone status", "droneID", droneID, "status", droneStatus, "error", err)
			return nil, err
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		r.logger.Error("Failed to commit transaction", "error", err)
		return nil, err
	}

	return order, nil
}
