# GRPC_Calc
### tg @leanq_ha
## Чтобы запустить:
```go run cmd/main.go --config=config/local.yaml```

Сервер работает на порту 8080.

Проверить работоспособность можно отправляя запросы можно через postman, загрузив в него прото файл
```proto/protos.proto```
Или через тесты **на момент написания работают не все**

## Как все работает
Для регистрации /register {email: , password: } в ответ user_id
Для логина /login {email: , password: } в ответ jwt-token
Для вычислений /calculate {expr: , uid: } в ответ id вычисления **Не реализовано до конца**
