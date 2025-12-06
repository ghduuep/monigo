# Monigo 📡

O **Monigo** é uma ferramenta robusta de monitorização de websites e registos DNS desenvolvida em Go. O sistema permite registar websites para verificação periódica de disponibilidade (HTTP) e monitorizar alterações críticas em registos DNS, notificando os utilizadores por e-mail sempre que ocorrem incidentes ou alterações inesperadas.

## 🚀 Funcionalidades

* **Monitorização HTTP**: Verificação periódica de estado (UP/DOWN), medição de latência e análise de códigos de resposta HTTP.
* **Monitorização de DNS Inteligente**:
    * Suporte para registos **A**, **AAAA**, **MX** e **NS**.
    * **Auto-Discovery**: Se não for fornecido um valor esperado, o sistema aprende automaticamente o valor atual do DNS na primeira verificação e passa a monitorizar alterações baseadas nesse valor.
* **Sistema de Notificações**: Envio automático de e-mails via SMTP para:
    * Falhas de disponibilidade (Site Down).
    * Recuperação de serviços (Site Up).
    * Falhas na resolução de DNS.
    * Alterações não autorizadas em registos DNS (Hijacking alerts).
* **Arquitetura Worker-Pool**: Separação clara entre a API (gestão de dados) e o Worker (processamento em *background*) para garantir performance e escalabilidade sem bloquear pedidos HTTP.
* **API REST**: Interface JSON moderna construída com o framework Echo para gestão de utilizadores e monitores.

## 🛠 Tech Stack

* **Linguagem**: [Go 1.25+](https://go.dev/)
* **Web Framework**: [Echo v4](https://echo.labstack.com/) (High performance, extensible, minimalist Go web framework).
* **Base de Dados**: PostgreSQL
* **Driver BD**: [pgx/v5](https://github.com/jackc/pgx) (Driver PostgreSQL de alta performance).
* **Infraestrutura**: Docker & Docker Compose (Builds *multi-stage* otimizados com Alpine Linux).

## 📂 Estrutura do Projeto

A estrutura segue os padrões modernos de projetos Go (Go Standard Project Layout):

* `cmd/api`: Ponto de entrada (`main.go`) para o servidor da API REST.
* `cmd/worker`: Ponto de entrada (`main.go`) para o serviço de monitorização em background.
* `internal/api`: Definição de rotas, handlers e lógica HTTP.
* `internal/database`: Repositórios, migrações e interação direta com o PostgreSQL.
* `internal/models`: Definições das estruturas de dados (`User`, `Monitor`, `CheckResult`).
* `internal/monitor`: Motores de verificação ("Checkers") para HTTP e DNS, e gestor de rotinas.
* `internal/notification`: Serviço de envio de e-mails e templates HTML responsivos.
