# Pingly 📡

**Pingly** é uma plataforma de monitorização de infraestrutura leve e eficiente desenvolvida em Go. Permite acompanhar a disponibilidade de websites, integridade de registos DNS e conectividade de portas TCP, enviando alertas em tempo real.

## 🚀 Funcionalidades

* **Monitorização HTTP(S)**: Verifica o status code (2xx-5xx) e latência.
* **Monitorização DNS**: Deteta alterações não autorizadas ou falhas em registos A, AAAA, MX, NS, TXT e CNAME.
* **Monitorização TCP/Ping**: Testa a conectividade de portas (TCP Handshake) em qualquer IP ou Host.
* **Notificações Multi-canal**:
    * 📧 E-mail (via SMTP).
    * ✈️ Telegram (Mensagens instantâneas).
* **Arquitetura Robusta**: Separação entre API e Worker, garantindo escalabilidade.

## 🛠 Tech Stack

* **Linguagem**: Go 1.25+
* **Framework Web**: Echo v4
* **Base de Dados**: PostgreSQL (pgx/v5)
* **Cache/Sessão**: Redis
* **Infraestrutura**: Docker & Docker Compose
