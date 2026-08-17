package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/JBK2116/buildon/internal/model"
)

// Repository provides functionality for interacting with the database.
type Repository struct {
	db     *sql.DB
	logger *slog.Logger
}

// NewRepository returns a Repository backed by the provided database connection and logger.
func NewRepository(db *sql.DB, logger *slog.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// CreateProject inserts a new project into the database.
func (r *Repository) CreateProject(title string, description string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDB)
	defer cancel()
	query := `
	INSERT INTO projects (title, description)
	VALUES (?, ?)
	`
	_, err := r.db.ExecContext(ctx, query, title, description)
	if err != nil {
		r.logger.Error("failed to insert project into database", "error", err)
		return err
	}
	return nil
}

// CreateProblem inserts a new problem into the database.
func (r *Repository) CreateProblem(projectID int, title string, content string, solved bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDB)
	defer cancel()
	query := `
	INSERT INTO problems (project_id, title, content, solved)
	VALUES (?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query, projectID, title, content, solved)
	if err != nil {
		r.logger.Error("failed to insert problem into database", "error", err)
		return err
	}
	return nil
}

// GetProjects queries all projects from the database.
func (r *Repository) GetProjects() ([]model.Project, error) {
	query := `
	SELECT p.id, p.title, p.description, p.created_at, p.updated_at,
	       pr.id, pr.project_id, pr.title, pr.content, pr.solved, pr.created_at, pr.updated_at
	FROM projects p
	LEFT JOIN problems pr ON pr.project_id = p.id
	ORDER BY p.created_at DESC, pr.created_at DESC
	`
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDB)
	defer cancel()
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			r.logger.Error("potential error in closing rows", "error", err)
		}
	}(rows)

	projects := make([]model.Project, 0)
	index := make(map[int]int)

	for rows.Next() {
		var p model.Project
		var probID, probProjectID sql.NullInt64
		var probTitle, probContent sql.NullString
		var probSolved sql.NullBool
		var probCreatedAt, probUpdatedAt sql.NullTime

		if scanErr := rows.Scan(
			&p.ID, &p.Title, &p.Description, &p.CreatedAt, &p.UpdatedAt,
			&probID, &probProjectID, &probTitle, &probContent, &probSolved, &probCreatedAt, &probUpdatedAt,
		); scanErr != nil {
			return nil, scanErr
		}

		i, ok := index[p.ID]
		if !ok {
			p.Problems = make([]model.Problem, 0)
			index[p.ID] = len(projects)
			projects = append(projects, p)
			i = index[p.ID]
		}

		if probID.Valid {
			projects[i].Problems = append(projects[i].Problems, model.Problem{
				ID:        int(probID.Int64),
				ProjectID: int(probProjectID.Int64),
				Title:     probTitle.String,
				Content:   probContent.String,
				Solved:    probSolved.Bool,
				CreatedAt: probCreatedAt.Time,
				UpdatedAt: probUpdatedAt.Time,
			})
		}
	}
	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, rowsErr
	}
	return projects, nil
}

// GetProject queries a project with the matching id from the database.
func (r *Repository) GetProject(id int) (*model.Project, error) {
	query := `
	SELECT p.id, p.title, p.description, p.created_at, p.updated_at,
	       pr.id, pr.project_id, pr.title, pr.content, pr.solved, pr.created_at, pr.updated_at
	FROM projects p
	LEFT JOIN problems pr ON pr.project_id = p.id
	WHERE p.id = ?
	ORDER BY pr.created_at DESC
	`
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDB)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			r.logger.Error("potential error in closing rows", "error", err)
		}
	}(rows)

	var project *model.Project

	for rows.Next() {
		var p model.Project
		var probID, probProjectID sql.NullInt64
		var probTitle, probContent sql.NullString
		var probSolved sql.NullBool
		var probCreatedAt, probUpdatedAt sql.NullTime

		if scanErr := rows.Scan(
			&p.ID, &p.Title, &p.Description, &p.CreatedAt, &p.UpdatedAt,
			&probID, &probProjectID, &probTitle, &probContent, &probSolved, &probCreatedAt, &probUpdatedAt,
		); scanErr != nil {
			return nil, scanErr
		}

		if project == nil {
			p.Problems = make([]model.Problem, 0)
			project = &p
		}

		if probID.Valid {
			project.Problems = append(project.Problems, model.Problem{
				ID:        int(probID.Int64),
				ProjectID: int(probProjectID.Int64),
				Title:     probTitle.String,
				Content:   probContent.String,
				Solved:    probSolved.Bool,
				CreatedAt: probCreatedAt.Time,
				UpdatedAt: probUpdatedAt.Time,
			})
		}
	}
	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, rowsErr
	}
	if project == nil {
		return nil, sql.ErrNoRows
	}
	return project, nil
}

// GetProblems queries all problems from the database.
func (r *Repository) GetProblems() ([]model.Problem, error) {
	query := `
	SELECT id, project_id, title, content, solved, created_at, updated_at
	FROM problems
	ORDER BY created_at DESC
	`
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDB)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			r.logger.Error("potential error in closing rows", "error", err)
		}
	}(rows)

	problems := make([]model.Problem, 0)
	for rows.Next() {
		var p model.Problem
		if scanErr := rows.Scan(
			&p.ID,
			&p.ProjectID,
			&p.Title,
			&p.Content,
			&p.Solved,
			&p.CreatedAt,
			&p.UpdatedAt,
		); scanErr != nil {
			return nil, scanErr
		}
		problems = append(problems, p)
	}
	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, rowsErr
	}
	return problems, nil
}

// GetProblem queries a problem from the database with the matching id.
func (r *Repository) GetProblem(id int) (*model.Problem, error) {
	query := `
	SELECT id, project_id, title, content, solved, created_at, updated_at
	FROM problems
	WHERE id = ?
	`
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDB)
	defer cancel()

	var p model.Problem
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.ProjectID, &p.Title, &p.Content, &p.Solved, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// UpdateProject updates the project with the provided id, title and description
// in the database.
func (r *Repository) UpdateProject(id int, title string, description string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDB)
	defer cancel()

	query := `
	UPDATE projects SET title = ?, description = ?, updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = ?
	`
	result, err := r.db.ExecContext(ctx, query, title, description, id)
	if err != nil {
		r.logger.Error("failed to update project in database", "error", err)
		return err
	}
	rowCount, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("failed to confirm update project success", "error", err)
		return err
	}
	if rowCount == 0 {
		return fmt.Errorf("no rows affected: project with id %d not found", id)
	}
	return nil
}

// UpdateProblem updates the problem with the provided id, title, content and
// solved status in the database.
func (r *Repository) UpdateProblem(id int, title string, content string, solved bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDB)
	defer cancel()

	query := `
	UPDATE problems SET title = ?, content = ?, solved = ?, updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = ?
	`
	result, err := r.db.ExecContext(ctx, query, title, content, solved, id)
	if err != nil {
		r.logger.Error("failed to update project in database", "error", err)
		return err
	}

	rowCount, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("failed to confirm update problem success", "error", err)
		return err
	}
	if rowCount == 0 {
		return fmt.Errorf("no rows affected: problem with id %d not found", id)
	}
	return nil
}

// DeleteProject deletes the project with the matching id in the database.
func (r *Repository) DeleteProject(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDB)
	defer cancel()

	query := `
	DELETE FROM projects WHERE id = ?
	`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete project in database", "error", err)
		return err
	}

	rowCount, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("failed to confirm delete project success", "error", err)
		return err
	}

	if rowCount == 0 {
		return fmt.Errorf("no rows affected: project with id %d not found", id)
	}
	return nil
}

// DeleteProblem deletes the problem with the matching id in the database.
func (r *Repository) DeleteProblem(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDB)
	defer cancel()

	query := `
	DELETE FROM problems WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete problem in database", "error", err)
		return err
	}

	rowCount, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("failed to confirm delete problem success", "error", err)
		return err
	}

	if rowCount == 0 {
		return fmt.Errorf("no rows affected: project with id %d not found", id)
	}
	return nil
}
