package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/unifocus/backend/internal/domain"
)

// OpportunityRepository handles opportunity data access operations
type OpportunityRepository struct {
	db *DB
}

// NewOpportunityRepository creates a new opportunity repository
func NewOpportunityRepository(db *DB) *OpportunityRepository {
	return &OpportunityRepository{db: db}
}

// Create creates a new opportunity
func (r *OpportunityRepository) Create(ctx context.Context, opp *domain.Opportunity) error {
	query := `
		INSERT INTO opportunities (
			title, type, description, source_url, source_type,
			competition_level, certification_type, organizer, organizer_type, award_level, points_value, is_official,
			start_date, deadline, event_date, location,
			requirements, eligibility_rules, target_majors,
			tags, attachments, description_vector, is_active
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
		opp.Title,
		opp.Type,
		opp.Description,
		opp.SourceURL,
		opp.SourceType,
		opp.CompetitionLevel,
		opp.CertificationType,
		opp.Organizer,
		opp.OrganizerType,
		opp.AwardLevel,
		opp.PointsValue,
		opp.IsOfficial,
		opp.StartDate,
		opp.Deadline,
		opp.EventDate,
		opp.Location,
		opp.Requirements,
		opp.EligibilityRules,
		pq.Array(opp.TargetMajors),
		pq.Array(opp.Tags),
		opp.Attachments,
		pq.Array(opp.DescriptionVector),
		opp.IsActive,
	).Scan(&opp.ID, &opp.CreatedAt, &opp.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

// GetByID retrieves an opportunity by ID
func (r *OpportunityRepository) GetByID(ctx context.Context, id int64) (*domain.Opportunity, error) {
	query := `
		SELECT id, title, type, description, source_url, source_type,
			competition_level, certification_type, organizer, organizer_type, award_level, points_value, is_official,
			start_date, deadline, event_date, location,
			requirements, eligibility_rules, target_majors,
			tags, attachments, description_vector, is_active, view_count, save_count,
			created_at, updated_at
		FROM opportunities
		WHERE id = $1
	`

	opp := &domain.Opportunity{}
	var targetMajors, tags []string
	var descriptionVector []float32

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&opp.ID,
		&opp.Title,
		&opp.Type,
		&opp.Description,
		&opp.SourceURL,
		&opp.SourceType,
		&opp.CompetitionLevel,
		&opp.CertificationType,
		&opp.Organizer,
		&opp.OrganizerType,
		&opp.AwardLevel,
		&opp.PointsValue,
		&opp.IsOfficial,
		&opp.StartDate,
		&opp.Deadline,
		&opp.EventDate,
		&opp.Location,
		&opp.Requirements,
		&opp.EligibilityRules,
		pq.Array(&targetMajors),
		pq.Array(&tags),
		&opp.Attachments,
		pq.Array(&descriptionVector),
		&opp.IsActive,
		&opp.ViewCount,
		&opp.SaveCount,
		&opp.CreatedAt,
		&opp.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("opportunity not found")
		}
		return nil, err
	}

	opp.TargetMajors = targetMajors
	opp.Tags = tags
	opp.DescriptionVector = descriptionVector

	return opp, nil
}

// List retrieves opportunities with filtering and pagination
func (r *OpportunityRepository) List(ctx context.Context, filter *domain.OpportunityFilter) ([]*domain.Opportunity, int64, error) {
	var conditions []string
	var args []interface{}
	argPos := 1

	// Build WHERE clause
	if filter.Type != "" {
		conditions = append(conditions, fmt.Sprintf("type = $%d", argPos))
		args = append(args, filter.Type)
		argPos++
	}

	if filter.CompetitionLevel != "" {
		conditions = append(conditions, fmt.Sprintf("competition_level = $%d", argPos))
		args = append(args, filter.CompetitionLevel)
		argPos++
	}

	if filter.Major != "" {
		conditions = append(conditions, fmt.Sprintf("$%d = ANY(target_majors)", argPos))
		args = append(args, filter.Major)
		argPos++
	}

	if filter.DeadlineAfter != nil {
		conditions = append(conditions, fmt.Sprintf("deadline >= $%d", argPos))
		args = append(args, *filter.DeadlineAfter)
		argPos++
	}

	if filter.DeadlineBefore != nil {
		conditions = append(conditions, fmt.Sprintf("deadline <= $%d", argPos))
		args = append(args, *filter.DeadlineBefore)
		argPos++
	}

	if filter.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", argPos))
		args = append(args, *filter.IsActive)
		argPos++
	}

	if len(filter.Tags) > 0 {
		conditions = append(conditions, fmt.Sprintf("tags && $%d", argPos))
		args = append(args, pq.Array(filter.Tags))
		argPos++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM opportunities %s", whereClause)
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Set default pagination
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	// Select opportunities
	query := fmt.Sprintf(`
		SELECT id, title, type, description, source_url, source_type,
			competition_level, certification_type, organizer, organizer_type, award_level, points_value, is_official,
			start_date, deadline, event_date, location,
			requirements, eligibility_rules, target_majors,
			tags, attachments, description_vector, is_active, view_count, save_count,
			created_at, updated_at
		FROM opportunities
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argPos, argPos+1)

	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var opportunities []*domain.Opportunity
	for rows.Next() {
		opp := &domain.Opportunity{}
		var targetMajors, tags []string
		var descriptionVector []float32

		err := rows.Scan(
			&opp.ID,
			&opp.Title,
			&opp.Type,
			&opp.Description,
			&opp.SourceURL,
			&opp.SourceType,
			&opp.CompetitionLevel,
			&opp.CertificationType,
			&opp.Organizer,
			&opp.OrganizerType,
			&opp.AwardLevel,
			&opp.PointsValue,
			&opp.IsOfficial,
			&opp.StartDate,
			&opp.Deadline,
			&opp.EventDate,
			&opp.Location,
			&opp.Requirements,
			&opp.EligibilityRules,
			pq.Array(&targetMajors),
			pq.Array(&tags),
			&opp.Attachments,
			pq.Array(&descriptionVector),
			&opp.IsActive,
			&opp.ViewCount,
			&opp.SaveCount,
			&opp.CreatedAt,
			&opp.UpdatedAt,
		)

		if err != nil {
			return nil, 0, err
		}

		opp.TargetMajors = targetMajors
		opp.Tags = tags
		opp.DescriptionVector = descriptionVector

		opportunities = append(opportunities, opp)
	}

	return opportunities, total, nil
}

// Update updates an opportunity
func (r *OpportunityRepository) Update(ctx context.Context, opp *domain.Opportunity) error {
	query := `
		UPDATE opportunities
		SET title = $1, type = $2, description = $3, source_url = $4, source_type = $5,
			competition_level = $6, certification_type = $7, organizer = $8, organizer_type = $9, 
			award_level = $10, points_value = $11, is_official = $12,
			start_date = $13, deadline = $14, event_date = $15, location = $16,
			requirements = $17, eligibility_rules = $18, target_majors = $19,
			tags = $20, attachments = $21, description_vector = $22, is_active = $23,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $24
		RETURNING updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
		opp.Title,
		opp.Type,
		opp.Description,
		opp.SourceURL,
		opp.SourceType,
		opp.CompetitionLevel,
		opp.CertificationType,
		opp.Organizer,
		opp.OrganizerType,
		opp.AwardLevel,
		opp.PointsValue,
		opp.IsOfficial,
		opp.StartDate,
		opp.Deadline,
		opp.EventDate,
		opp.Location,
		opp.Requirements,
		opp.EligibilityRules,
		pq.Array(opp.TargetMajors),
		pq.Array(opp.Tags),
		opp.Attachments,
		pq.Array(opp.DescriptionVector),
		opp.IsActive,
		opp.ID,
	).Scan(&opp.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("opportunity not found")
		}
		return err
	}

	return nil
}

// Delete soft deletes an opportunity (sets is_active = false)
func (r *OpportunityRepository) Delete(ctx context.Context, id int64) error {
	query := `
		UPDATE opportunities
		SET is_active = false, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("opportunity not found")
	}

	return nil
}

// IncrementViewCount increments the view count for an opportunity
func (r *OpportunityRepository) IncrementViewCount(ctx context.Context, id int64) error {
	query := `UPDATE opportunities SET view_count = view_count + 1 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// IncrementSaveCount increments the save count for an opportunity
func (r *OpportunityRepository) IncrementSaveCount(ctx context.Context, id int64) error {
	query := `UPDATE opportunities SET save_count = save_count + 1 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
