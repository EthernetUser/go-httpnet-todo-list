package postgres

import (
	"context"
	"go-httpnet-todo-list/internal/database"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres interface {
	GetTasks(ctx context.Context, userId int) ([]database.Task, error)
	MarkTask(ctx context.Context, taskId int, userId int, done bool) error
	MarkAsDeleted(ctx context.Context, taskId int, userId int) error
	AddTask(
		ctx context.Context,
		userId int,
		title string,
		description string,
	) error
	Close()
}

type postgres struct {
	dbpool *pgxpool.Pool
}

func New(connString string) *postgres {
	dbpool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		panic(err)
	}

	if err := dbpool.Ping(context.Background()); err != nil {
		panic(err)
	}

	_, err = dbpool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS tasks (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			description TEXT,
			status BOOLEAN NOT NULL DEFAULT false,
			deleted_at TIMESTAMP,
			is_deleted BOOLEAN NOT NULL DEFAULT false
		);
	`)
	if err != nil {
		panic(err)
	}

	return &postgres{
		dbpool: dbpool,
	}
}

func (p *postgres) Close() {
	p.dbpool.Close()
}

func (p *postgres) GetTasks(
	ctx context.Context,
	userId int,
) ([]database.Task, error) {
	tx, err := p.dbpool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	query := `
  		SELECT id, title, description, status
  		FROM tasks
  		WHERE user_id = $1 AND is_deleted = false
	`
	rows, err := tx.Query(
		ctx,
		query,
		userId,
	)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	tasks := make([]database.Task, 0)
	for rows.Next() {
		task := database.Task{}
		if err := rows.Scan(&task.Id, &task.Title, &task.Description, &task.Done); err != nil {
			tx.Rollback(ctx)
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		tx.Rollback(ctx)
		return nil, err

	}

	if err := tx.Commit(ctx); err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	return tasks, nil
}

func (p *postgres) MarkTask(
	ctx context.Context,
	taskId int,
	userId int,
	done bool,
) error {
	tx, err := p.dbpool.Begin(ctx)
	if err != nil {
		return err
	}

	query := `
		UPDATE tasks SET status = $1 
		WHERE id = $2 
		AND user_id = $3 
		AND is_deleted = false
	`
	_, err = tx.Exec(
		ctx,
		query,
		done,
		taskId,
		userId,
	)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		tx.Rollback(ctx)
		return err
	}

	return nil
}

func (p *postgres) MarkAsDeleted(
	ctx context.Context,
	taskId int,
	userId int,
) error {
	tx, err := p.dbpool.Begin(ctx)
	if err != nil {
		return err
	}

	query := `
		UPDATE tasks SET deleted_at = now(), is_deleted = true 
		WHERE id = $1 
		AND user_id = $2 
		AND is_deleted = false
	`
	_, err = tx.Exec(
		ctx,
		query,
		taskId,
		userId,
	)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		tx.Rollback(ctx)
		return err
	}

	return nil
}

func (p *postgres) AddTask(ctx context.Context, task database.Task) error {
	tx, err := p.dbpool.Begin(ctx)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO tasks (user_id, title, description) 
		VALUES ($1, $2, $3)
	`
	_, err = tx.Exec(
		ctx,
		query,
		task.UserId,
		task.Title,
		task.Description,
	)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		tx.Rollback(ctx)
		return err
	}

	return nil
}
