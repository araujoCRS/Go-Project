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

## Passo a passo para executar o projeto

1. **Clone o repositório e acesse a pasta do projeto:**
   ```sh
   git clone <url-do-repositorio>
   cd Go-Project/client
   ```

2. **Suba as dependências no Kubernetes:**
   - Certifique-se de que o Kubernetes está rodando (pode usar Rancher Desktop).
   - Aplique os manifestos para criar o namespace, PostgreSQL e RabbitMQ:
     ```sh
     kubectl apply -f infra/init_cluster.yaml
     kubectl apply -f infra/postgres.yaml
     kubectl apply -f infra/rabbitMQ.yaml
     ```

3. **Faça o port-forward para acessar PostgreSQL e RabbitMQ localmente:**
   - Em dois terminais separados, execute:
     ```sh
     kubectl port-forward svc/postgres 5432:5432 -n golang-apps
     kubectl port-forward svc/rabbitmq 5672:5672 15672:15672 -n golang-apps
     ```

4. **Configure o arquivo `configs/config.yaml` se necessário:**
   - Verifique se os parâmetros de conexão estão de acordo com seu ambiente.

5. **Compile e execute a API e o worker:**
   - No Windows, basta rodar:
     ```powershell
     .\build.ps1
     ```
   - Isso irá compilar e iniciar tanto a API quanto o worker.

6. **Acesse a API:**
    - A API estará disponível em `http://localhost:8080`. A API expõe os seguintes endpoints para gerenciamento de clientes:

    - `POST /api/client`  
      Cria um novo cliente (os dados devem ser enviados no corpo da requisição).

    - `PUT /api/client`  
      Atualiza um cliente existente (os dados devem ser enviados no corpo da requisição).

    - `GET /api/client/:id`  
      Busca um cliente pelo seu ID.

    - `DELETE /api/client/:id`  
      Remove um cliente pelo seu ID.

    Além disso, a documentação interativa da API (Swagger) está disponível em:  
    `GET /swagger/index.html`
7. **Acesse a API:**
  Para criar ou atualizar um cliente via API, envie um JSON no seguinte formato no corpo da requisição:

  ```json
  {
      "id": 42,
      "name": "Roberto",
      "sobrenome": "Araujo",
      "contato": "roberto@hotmail.com",
      "cpf": "75735806145",
      "endereco": "Rua 6 Jd Goias",
      "nascimento": "1994-07-25T00:00:00Z"
  }
```
Utilize este modelo ao fazer requisições para os endpoints `POST /api/client` e `PUT /api/client`.

---
---

## Pré-requisitos

- Go 1.20+
- Docker e Kubernetes (opcional, para rodar via manifestos)
- PowerShell (para rodar o script Build.ps1 no Windows)

---

## Configuração do ambiente Kubernetes

Os arquivos em `infra/` permitem subir o ambiente completo no Kubernetes.

## `Build.ps1`

O script `Build.ps1` automatiza o processo de build e execução dos binários da aplicação no Windows. Ele realiza as seguintes etapas:

- Compila o worker (`cmd/worker/main.go`) e gera o executável em `pkg/worker.exe`.
- Compila a API (`cmd/api/main.go`) e gera o executável em `pkg/api.exe`.
- Inicia ambos os executáveis em processos separados.

Exemplo de uso:

```powershell
.\build.ps1
```

Esse script requer o PowerShell e o Go instalados no sistema.

## Dica: Gerenciamento visual com Rancher Desktop e OpenLens

Para facilitar o gerenciamento dos containers do banco de dados PostgreSQL e do RabbitMQ em ambientes locais ou de desenvolvimento, recomenda-se o uso de ferramentas como **Rancher Desktop** e **OpenLens**. 

- **Rancher Desktop**: Permite criar, visualizar e gerenciar clusters Kubernetes de forma simples, além de facilitar o controle dos recursos implantados, como pods, serviços e volumes.
- **OpenLens**: Fornece uma interface gráfica avançada para monitorar e administrar clusters Kubernetes, facilitando a visualização do status dos containers, logs, métricas e eventos.

## Observação sobre portas e conexões

Tanto a API quanto o worker foram desenvolvidos para se conectar aos serviços utilizando as portas e parâmetros definidos nos arquivos `configs/config.yaml` e `infra/rabbitMQ.yaml`. Por padrão, a aplicação espera:

- **PostgreSQL** disponível em `localhost:5432`
- **RabbitMQ** disponível em `localhost:5672` (AMQP) e `localhost:15672` (interface de gerenciamento)
- **API** escutando na porta `8080`

Essas configurações são refletidas tanto no arquivo de configuração da aplicação (`config.yaml`) quanto nos manifestos do Kubernetes (`rabbitMQ.yaml` e `postgres.yaml`).  
Caso seja necessário alterar as portas ou endereços, lembre-se de atualizar ambos os arquivos para garantir o correto funcionamento da comunicação entre os componentes do sistema.

## Acesso externo ao RabbitMQ e PostgreSQL no Kubernetes

Neste projeto, **não foi criado um recurso Ingress no Kubernetes** para expor o RabbitMQ e o PostgreSQL externamente. Por questões de simplicidade e segurança em ambientes de desenvolvimento, recomenda-se utilizar o comando `kubectl port-forward` para acessar esses serviços a partir da sua máquina local.

Exemplo de uso para PostgreSQL:

```sh
kubectl port-forward postgres-xxxx 5432:5432 -n golang-apps
```

Exemplo de uso para RabbitMQ (AMQP e interface web):

```sh
kubectl port-forward rabbitmq-xxxx 5672:5672 15672:15672 -n golang-apps
```

Dessa forma, os serviços ficam acessíveis em `localhost` nas portas padrão, conforme esperado pela aplicação.

Além disso, **não foram criados Dockerfiles específicos para a API e o worker**. Como o foco é o uso local com port-forwarding e o build dos binários é feito diretamente via script PowerShell, a criação de imagens customizadas não é necessária neste cenário. Caso deseje rodar os componentes em containers, será preciso criar os Dockerfiles posteriormente.

### `init_cluster.yaml`

Cria o namespace `golang-apps` para isolar os recursos do projeto:

## `config.yaml`

O arquivo `config.yaml` centraliza as configurações da aplicação, facilitando a customização dos parâmetros sem necessidade de alterar o código-fonte. Ele está localizado em `configs/config.yaml` e possui as seguintes seções principais:

- **api**: Configurações da API, como porta de escuta (`port`) e nomes das filas utilizadas.
- **database**: Parâmetros de conexão com o banco PostgreSQL, incluindo host, porta, usuário, senha, nome do banco, modo SSL e timeout.
- **rabbitmq**: Configurações do broker RabbitMQ, como host, porta, credenciais e propriedades das filas.

### Sobre o tipo de fila configurada em `rabbitmq.queues`

No exemplo do projeto, a fila `clients` está configurada como uma fila do tipo **quorum**, a configuração inclui parâmetros para dead-letter exchange (DLX), permitindo o redirecionamento de mensagens que não puderam ser processadas após um determinado número de tentativas (`x-delivery-limit`). 

## `postgres.yaml`

O arquivo `postgres.yaml` define todos os recursos necessários para executar uma instância do PostgreSQL no Kubernetes dentro do namespace `golang-apps`. Ele inclui:

- **Deployment**: Cria e gerencia o pod do PostgreSQL, especificando a imagem utilizada, variáveis de ambiente para configuração do banco (usuário, senha e nome do banco), portas expostas e o volume onde os dados serão persistidos.
- **Service**: Expõe o PostgreSQL internamente no cluster Kubernetes, permitindo que outros pods acessem o banco de dados pela porta 5432.
- **PersistentVolumeClaim (PVC)**: Garante que os dados do banco sejam armazenados de forma persistente, mesmo que o pod seja reiniciado ou recriado.

Esse arquivo facilita a implantação do PostgreSQL de forma automatizada e consistente no ambiente Kubernetes do projeto.

## Arquitetura do Projeto

A arquitetura deste projeto segue uma separação clara de responsabilidades, utilizando o padrão de camadas para facilitar a manutenção, testes e evolução do sistema. Abaixo está um resumo dos principais componentes e suas funções, com base na estrutura de pastas e interfaces encontradas:

## Arquitetura do Projeto

A arquitetura do projeto é baseada em camadas bem definidas e utiliza abstrações para promover baixo acoplamento, reutilização de código e facilidade de testes. Veja como os principais componentes e modelos de abstração estão organizados:

### Estrutura de Pastas

- **cmd/**: Pontos de entrada da aplicação (API e worker).
- **configs/**: Arquivos de configuração.
- **infra/**: Manifestos Kubernetes para infraestrutura.
- **internal/service/**: Lógica de negócio e serviços.
  - **repository/**: Repositórios para acesso a dados.
- **shared/models/**: Modelos de dados compartilhados.
- **worker/**: Lógica do worker para processamento assíncrono.

### Abstrações e Interfaces

O projeto faz uso extensivo de interfaces genéricas para desacoplar as camadas e facilitar a manutenção:

- **ServiceBase[T] (`internal/service/service_base.go`)**  
  Interface genérica para serviços de domínio, definindo operações básicas como `Get`, `Save` e `Delete`. Permite que diferentes entidades (como Cliente) implementem serviços reutilizando a mesma assinatura.

- **RepositoryBase[T] (`internal/service/repository/repository_base.go`)**  
  Interface genérica para repositórios, padronizando métodos de acesso ao banco de dados: `Get`, `Create`, `Update` e `Delete`. Isso facilita a implementação de repositórios para diferentes entidades com o mesmo contrato.

- **Worker (`worker/worker.go`)**  
  Interface que define o contrato para workers assíncronos, com métodos `Run` e `Handler`.  
  - **BaseWorker**: Estrutura base que implementa a lógica comum de consumo de mensagens da fila RabbitMQ e delega o processamento para o método `Handler` de cada worker específico.

### Exemplo de Implementação

- **ClientService** implementa `ServiceBase[Client]`, orquestrando regras de negócio e persistência.
- **ClientRepository** implementa `RepositoryBase[Client]`, lidando com operações no banco de dados.
- **WorkerClient** implementa a interface `Worker`, processando mensagens da fila e utilizando o serviço de clientes.

### Benefícios do Modelo de Abstração

- **Reutilização de código**: Interfaces genéricas evitam duplicidade e facilitam a criação de novos serviços e repositórios.
- **Baixo acoplamento**: As camadas de serviço, repositório e worker são independentes, facilitando manutenção e testes.
- **Testabilidade**: Interfaces permitem a criação de mocks para testes unitários.
- **Escalabilidade**: O padrão worker permite processar mensagens de forma assíncrona e escalável.

> **Nota:**  
> Este projeto foi beneficiado pelo uso de Inteligência Artificial generativa para acelerar tanto o desenvolvimento do código quanto a produção da documentação. Ferramentas de IA auxiliaram na geração de exemplos, explicações técnicas e revisão, contribuindo para maior produtividade e clareza na entrega final.

> **Nota sobre a arquitetura:**  
> A separação em camadas e o uso de abstrações neste projeto foram inspirados nas boas práticas do desenvolvimento em C#, especialmente no uso de interfaces e padrões de projeto amplamente adotados no ecossistema .NET, como o padrão MVC (Model-View-Controller).  
> Assim como no .NET, as interfaces são utilizadas para definir contratos claros entre as camadas de serviço, repositório e processamento assíncrono, promovendo baixo acoplamento, testabilidade e facilidade de manutenção. Essa abordagem facilita a evolução do sistema e a adoção de boas práticas de engenharia de software, tornando o projeto mais robusto e escalável.