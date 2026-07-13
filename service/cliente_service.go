package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ichdavid/pedidos-api/model"
	"github.com/ichdavid/pedidos-api/repository"
)

// struct que possui o repositório de clientes como dependência,
// para acesso de forma indireta e segura ao banco de dados.
type ClienteService struct {
	repo *repository.ClienteRepository
}

// construtor que recebe o repositorio de clientes como parâmetro e
// retorna uma nova service de cliente com atualização do repositório.
func NewClienteService(repo *repository.ClienteRepository) *ClienteService {
	return &ClienteService{repo: repo}
}

// metodo create do tipo clienteService que recebe um contexto e uma struct createClienteRequest so com os dados passados pelo usuario,
// e retorna um cliente e um erro.
func (s *ClienteService) Create(ctx context.Context, req model.CreateClienteRequest) (*model.Cliente, error) {
	// verifica se nenhum campo esta vazio.
	if req.Name == "" || req.Email == "" || req.Password == "" {
		return nil, fmt.Errorf("name, email e password são obrigatórios")
	}

	// pega o password e gera o hashPassaword para ser enviado pro BD.
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar hash da senha: %w", err)
	}

	// Cria um cliente com os dados recebidos no Json e salva o hash na senha.
	cliente := &model.Cliente{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hash),
	}

	// Salva o cliente criado no BD pelo repository.
	if err := s.repo.Create(ctx, cliente); err != nil {
		return nil, err
	}

	// se finalizado com sucesso, retorna o cliente criado e um erro nulo.
	return cliente, nil
}

// metodo que atraves do repository busca e retorna a lista de clientes do BD
func (s *ClienteService) FindAll(ctx context.Context) ([]model.Cliente, error) {
	return s.repo.FindAll(ctx)
}

// metodo que atraves do repository busca e retorna o cliente atraves pelo ID no BD
func (s *ClienteService) FindByID(ctx context.Context, id uuid.UUID) (*model.Cliente, error) {
	return s.repo.FindByID(ctx, id)
}

// acabam sendo metodos repetitivos, mas que em questão de responsabilidades são necessario, caso no futuro prescisem de atualizações para validação de alguma variavel.
