package shared

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

// Valida os campos obrigatórios do Client.
// Retorna nil se válido, ou um erro descritivo se inválido.
func (c *Client) Validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return errors.New("nome é obrigatório")
	}
	if strings.TrimSpace(c.Sobrenome) == "" {
		return errors.New("sobrenome é obrigatório")
	}
	if strings.TrimSpace(c.Contato) == "" {
		return errors.New("contato é obrigatório")
	}
	if strings.TrimSpace(c.Cpf) == "" {
		return errors.New("CPF é obrigatório")
	}
	if !isValidCPF(c.Cpf) {
		return errors.New("CPF inválido")
	}
	if strings.TrimSpace(c.Endereco) == "" {
		return errors.New("endereço é obrigatório")
	}
	if c.DataNascimento.IsZero() {
		return errors.New("data de nascimento é obrigatória")
	}
	// Opcional: validar se DataNascimento não é futura
	if c.DataNascimento.After(time.Now()) {
		return errors.New("data de nascimento não pode ser futura")
	}
	return nil
}

// Validação simples de CPF (apenas tamanho e dígitos)
func isValidCPF(cpf string) bool {
	re := regexp.MustCompile(`^\d{11}$`)
	return re.MatchString(cpf)
}
