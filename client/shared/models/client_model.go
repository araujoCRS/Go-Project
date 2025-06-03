package shared

import (
	"time"
)

type Client struct {
	ModelBase
	Name           string    `json:"name" example:"João Silva"`
	Sobrenome      string    `json:"sobrenome" example:"da Silva"`
	Contato        string    `json:"contato" example:"joao.silva@example.com"`
	Cpf            string    `json:"cpf" example:"75735506147"`
	Endereco       string    `json:"endereco" example:"Rua Exemplo, 123, Bairro, Cidade, Estado, 12345-678"`
	DataNascimento time.Time `json:"nascimento" example:"1990-12-31T00:00:00Z"`
}
