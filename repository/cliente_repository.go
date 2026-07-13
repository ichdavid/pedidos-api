package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ichdavid/pedidos-api/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ClienteRepository é responsável por gerenciar as operações de banco de dados relacionadas aos clientes.
type ClienteRepository struct {
	BD *pgxpool.Pool
}

// NewClienteRepository cria uma nova instância de ClienteRepository.
func NewClienteRepository(BD *pgxpool.Pool) *ClienteRepository {
	return &ClienteRepository{BD: BD}
}

// Função create do tipo clienteRepository que cria um novo cliente no banco
// passando nome, email e senha para o bando de dados
func (r *ClienteRepository) Create(ctx context.Context, cliente *model.Cliente) error {

	query := `
        INSERT INTO clientes (name, email, password_hash)
        VALUES ($1, $2, $3)
        RETURNING id, created_at
    `
	// Executa a query de inserção no banco de dados, passando os valores do cliente e obtendo o
	// id e data de criação do postgres para a struct cliente.
	err := r.BD.QueryRow(ctx, query,
		cliente.Name,
		cliente.Email,
		cliente.PasswordHash,
	).Scan(&cliente.ID, &cliente.CreatedAt)

	// caso ocorra algum erro de duplicidade de email (já que é uma restrição única),
	// retorna uma mensagem de erro informando que o email já está cadastrado,
	// caso ocorra algum outro erro, retorna uma mensagem de erro genérica.
	if err != nil {
		if isUniqueViolation(err) {
			return fmt.Errorf("email já cadastrado")
		}
		return fmt.Errorf("erro ao criar cliente: %w", err)
	}

	return nil
}

// metodo do tipo clienteRepository que recebe um contexto e retorna uma lista de todos os clientes cadastrados no banco de dados e um erro.
func (r *ClienteRepository) FindAll(ctx context.Context) ([]model.Cliente, error) {
	//essa query utiliza o select para selecionar os atributos da busca(id, name, email, created_at) da tabela clientes e ordena os resultados pela data de criação em ordem decrescente.
	query := `
        SELECT id, name, email, created_at
        FROM clientes
        ORDER BY created_at DESC
    `

	//executa a query e retorna em rows todas as linhas encontradas e um erro caso ocorra algum problema na execução da query.
	rows, err := r.BD.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar clientes: %w", err)
	}
	defer rows.Close()

	//cria um slice de cliente vazio, com base no model cliente e uma variavel do tipo cliente para rodar todos os resultados do rows
	//e add no slice de clientes desde que não haja erro na leitura dos dados.
	var clientes []model.Cliente
	for rows.Next() {
		var c model.Cliente
		err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear cliente: %w", err)
		}
		clientes = append(clientes, c)
	}

	//retorna o slice de clientes preenchido e um erro nulo.
	return clientes, nil
}

// Metodo semelhante ao anterior, porem vai buscar apenas um cliente baseado no id dele.
func (r *ClienteRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Cliente, error) {
	query := `
        SELECT id, name, email, created_at
        FROM clientes
        WHERE id = $1
    `

	var c model.Cliente
	err := r.BD.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.Name, &c.Email, &c.CreatedAt,
	)

	//pgx.ErrorNoRows é um erro da biblioteca pgx, quando não ocorre erro na query,
	// mas também não é encontrado nenhum resultado.
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("cliente não encontrado")
	}
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar cliente: %w", err)
	}

	// O retorno é o cliente encontrado e um erro nulo, em caso de sucesso.
	return &c, nil
}
