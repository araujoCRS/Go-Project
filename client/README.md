# Client API & Worker - GoLang

Este projeto contém uma API REST para cadastro de clientes e um worker para processamento assíncrono, ambos escritos em Go. O sistema utiliza PostgreSQL como banco de dados e RabbitMQ como broker de mensagens.

---

## Sumário

- [Pré-requisitos](#pré-requisitos)
- [Configuração do ambiente Kubernetes](#configuração-do-ambiente-kubernetes)
  - [`init_cluster.yaml`](#init_clusteryaml)
  - [`postgres.yaml`](#postgtesyaml)
  - [`rabbitMQ.yaml`](#rabbitmqyaml)
- [Configuração da aplicação](#configuração-da-aplicação)
  - [`config.yaml`](#configyaml)
- [Build e execução](#build-e-execução)
  - [`Build.ps1`](#buildps1)
- [Sobre a fila RabbitMQ](#sobre-a-fila-rabbitmq)
- [Swagger](#swagger)

---

## Pré-requisitos

- Go 1.20+
- Docker e Kubernetes (opcional, para rodar via manifestos)
- PowerShell (para rodar o script Build.ps1 no Windows)

---

## Configuração do ambiente Kubernetes

Os arquivos em `infra/` permitem subir o ambiente completo no Kubernetes.

### `init_cluster.yaml`

Cria o namespace `golang-apps` para isolar os recursos do projeto:

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: golang-apps