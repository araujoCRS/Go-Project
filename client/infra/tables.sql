CREATE TABLE cliente (
    id SERIAL PRIMARY KEY,
    nome VARCHAR(100) NOT NULL,
    sobrenome VARCHAR(100) NOT NULL,
    contato VARCHAR(50),
    endereco TEXT,
    data_nascimento DATE,
    cpf CHAR(11) UNIQUE NOT NULL
);

-- Índice para buscas rápidas pelo CPF
CREATE INDEX idx_client_cpf ON cliente (cpf);