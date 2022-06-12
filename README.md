# Depositaja: Wallet deposit service using event-driven architecture

In this project, I build a fictitious wallet deposit service. Our main goal is to explore implementation of event-driven architecture using Kafka as the message broker.
Several assumptions were made to pertain simplicity of the project.

The service offers only two endpoints:

1. `localhost:8080/deposit` HTTP POST endpoint for user to deposit money.
2. `localhost:8080/check/{wallet_id}` HTTP GET endpoint to get the balance of the wallet, also a flag whether the wallet has ever done one or more deposits with amounts more than 10,000 within a single 2-minute window (rolling-period).