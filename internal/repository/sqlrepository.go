package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"task-tracker/internal/models"
)

type PostgresTaskRepository struct {
	db *pgxpool.Pool
}

func NewPostgresTaskRepository(db *pgxpool.Pool) *PostgresTaskRepository {
	return &PostgresTaskRepository{
		db: db,
	}
}

func CreateDBPool(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("error parsing config: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("error connecting to db: %v", err)
	}

	return pool, nil
}

func (repo *PostgresTaskRepository) Add(ctx context.Context, task *models.Task) error {
	query := `INSERT INTO tasks (id, title, description, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := repo.db.Exec(
		ctx,
		query,
		task.ID,
		task.Title,
		task.Description,
		task.Status,
		task.CreatedAt,
		task.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error adding task: %v", err)
	}

	return nil
}

func (repo *PostgresTaskRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM tasks WHERE id=$1`
	_, err := repo.db.Exec(ctx, query, id)

	if err != nil {
		return fmt.Errorf("error deleting task: %v", err)
	}

	return nil
}

func (repo *PostgresTaskRepository) Exists(ctx context.Context, id string) (bool, error) {
	var exists bool

	query := `SELECT EXISTS(SELECT 1 FROM tasks WHERE id=$1)`
	err := repo.db.QueryRow(ctx, query, id).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("error checking if task exists: %v", err)
	}

	return exists, nil
}

func (repo *PostgresTaskRepository) Get(ctx context.Context, id string) (models.Task, error) {
	var task models.Task

	query := `SELECT id, title, description, status, created_at, updated_at FROM tasks WHERE id=$1`
	err := repo.db.QueryRow(ctx, query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		return models.Task{}, models.ErrTaskNotFound
	}

	return task, nil
}

func (repo *PostgresTaskRepository) GetAll(ctx context.Context) ([]models.Task, error) {
	query := `SELECT id, title, description, status, created_at, updated_at FROM tasks`
	rows, err := repo.db.Query(ctx, query)

	if err != nil {
		return []models.Task{}, fmt.Errorf("error getting tasks: %v", err)
	}

	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {
		var task models.Task
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.CreatedAt,
			&task.UpdatedAt,
		)

		if err != nil {
			return []models.Task{}, fmt.Errorf("error scanning row: %v", err)
		}

		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return []models.Task{}, fmt.Errorf("error iterating rows: %w", err)
	}

	return tasks, nil
}

func (repo *PostgresTaskRepository) Update(ctx context.Context, updatedTask *models.Task) error {
	query := `UPDATE tasks SET title=$1, description=$2, status=$3, updated_at=$4 WHERE id=$5`
	updatedTask.UpdatedAt = time.Now().Format(time.RFC3339Nano)
	_, err := repo.db.Exec(
		ctx,
		query,
		updatedTask.Title,
		updatedTask.Description,
		updatedTask.Status,
		updatedTask.UpdatedAt,
		updatedTask.ID,
	)

	if err != nil {
		return fmt.Errorf("error updating task: %v", err)
	}

	return nil
}
