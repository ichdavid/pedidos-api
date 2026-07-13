package controller

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/ichdavid/pedidos-api/model"
	"github.com/ichdavid/pedidos-api/service"
)

type ClienteController struct {
	service *service.ClienteService
}

func NewClienteController(service *service.ClienteService) *ClienteController {
	return &ClienteController{service: service}
}

// Metodo create do controller que recebe uma requisição e um response HTTP e
// no final atualiza se tiver sucesso o cliente no service (que atualiza o BD via repository) de acordo com o Json recebido e
// retorna um Json em resposta da criação do cliente no BD.
func (c *ClienteController) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateClienteRequest

	//verifica o corpo do Json e atualiza a struct createClienteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "dados inválidos")
		return
	}

	//Cria um cliente para atualizar no service, baseado no Json recebido
	cliente, err := c.service.Create(r.Context(), req)
	if err != nil {
		if err.Error() == "email já cadastrado" {
			writeError(w, http.StatusConflict, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, cliente)
}

// atraves da requisição no postman vai puxar o FindAll no service e no repository e retornar a lista de clientes e um status.
func (c *ClienteController) FindAll(w http.ResponseWriter, r *http.Request) {
	clientes, err := c.service.FindAll(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro ao listar clientes")
		return
	}

	writeJSON(w, http.StatusOK, clientes)
}

// atraves da requisição no postman vai puxar o FindByID no service e no repository e retornar a lista de clientes e um status.
func (c *ClienteController) FindByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		writeError(w, http.StatusBadRequest, "id inválido")
		return
	}

	cliente, err := c.service.FindByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "cliente não encontrado")
		return
	}

	writeJSON(w, http.StatusOK, cliente)
}
