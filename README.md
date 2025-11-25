# **Pingly**

O **Pingly** é uma ferramenta robusta de monitorização de websites e registos DNS escrita em Go. O sistema permite que os utilizadores registem websites para verificação periódica de disponibilidade (HTTP) e monitorizem alterações em registos DNS (A, AAAA, MX, NS), enviando notificações por email sempre que ocorrem mudanças de estado ou configuração.

## **🚀 Funcionalidades**

* **Monitorização HTTP**: Verifica periodicamente o estado de websites (UP/DOWN).  
* **Monitorização de DNS**: Acompanha alterações nos registos A, AAAA, MX e NS de domínios.  
* **Notificações**: Envio automático de emails ao detetar falhas no website ou alterações no DNS.  
* **API REST**: Gestão de utilizadores e monitores através de uma API segura.  
* **Autenticação JWT**: Proteção de rotas e gestão de sessões de utilizador.  
* **Worker Dedicado**: Processamento em *background* para verificações contínuas sem bloquear a API.

## **🛠 Tecnologias Utilizadas**

* **Linguagem**: [Go](https://go.dev/) (Golang)  
* **Base de Dados**: [PostgreSQL](https://www.postgresql.org/)  
* **Driver BD**: [pgx/v5](https://github.com/jackc/pgx)  
* **Router HTTP**: [chi](https://github.com/go-chi/chi)  
* **Autenticação**: [jwtauth](https://github.com/go-chi/jwtauth)  
* **Containerização**: [Docker](https://www.docker.com/) (para a base de dados)

## **📂 Estrutura do Projeto**

* `cmd/api`: Ponto de entrada para o servidor da API REST.  
* `cmd/worker`: Ponto de entrada para o worker de monitorização em background.  
* `internal/api`: Definição de rotas, handlers e middleware.  
* `internal/database`: Lógica de interação com o PostgreSQL.  
* `internal/models`: Estruturas de dados (Users, Websites, DNSMonitors).  
* `internal/monitor`: Lógica principal de verificação HTTP e DNS.  
* `internal/notification`: Serviço de envio de emails (SMTP).