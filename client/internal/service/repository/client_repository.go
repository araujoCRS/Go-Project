package repository

import (
	shared "client/shared/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ClientRepository interface {
	RepositoryBase[shared.Client]
}

type clientRepository struct {
	db *pgxpool.Pool
}

func NewClientRepository(db any) ClientRepository {
	return &clientRepository{db: (db).(*pgxpool.Pool)}
}

func scanClient(dest *shared.Client, scanFn func(dest ...any) error) error {
	return scanFn(
		&dest.Id,
		&dest.Name,
		&dest.Sobrenome,
		&dest.Contato,
		&dest.Cpf,
		&dest.Endereco,
		&dest.DataNascimento,
	)
}

func queryRowClient(ctx context.Context, client *shared.Client, db *pgxpool.Pool, query string, args ...any) (*shared.Client, error) {
	row := db.QueryRow(ctx, query, args...)
	if err := scanClient(client, row.Scan); err != nil {
		return nil, err
	}
	return client, nil
}

func (r *clientRepository) Get(ctx context.Context, id int) (*shared.Client, error) {
	var client shared.Client
	query := `SELECT id, nome, sobrenome, contato, cpf, endereco, data_nascimento FROM cliente WHERE id = $1`
	return queryRowClient(ctx, &client, r.db, query, id)
}

func (r *clientRepository) Create(ctx context.Context, client *shared.Client) (*shared.Client, error) {
	query := `INSERT INTO cliente (nome, sobrenome, contato, cpf, endereco, data_nascimento) 
        VALUES ($1, $2, $3, $4, $5, $6) 
        RETURNING id, nome, sobrenome, contato, cpf, endereco, data_nascimento`
	return queryRowClient(ctx, client, r.db, query,
		client.Name,
		client.Sobrenome,
		client.Contato,
		client.Cpf,
		client.Endereco,
		client.DataNascimento,
	)
}

func (r *clientRepository) Update(ctx context.Context, client *shared.Client) (*shared.Client, error) {
	query := `UPDATE cliente SET nome = $1, sobrenome = $2, contato = $3, cpf = $4, endereco = $5, data_nascimento = $6 WHERE id = $7`
	result, err := r.db.Exec(ctx, query,
		client.Name,
		client.Sobrenome,
		client.Contato,
		client.Cpf,
		client.Endereco,
		client.DataNascimento,
		client.Id,
	)

	if err != nil || result.RowsAffected() == 0 {
		return client, fmt.Errorf("error: %w", err)
	}

	return client, nil

}

func (r *clientRepository) Delete(ctx context.Context, id int) (bool, error) {
	query := `DELETE FROM cliente WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil || result.RowsAffected() == 0 {
		return false, fmt.Errorf("error: %w", err)
	}

	return true, nil
}
