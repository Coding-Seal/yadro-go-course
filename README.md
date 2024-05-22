# yadro-go-course
Course from Yadro on microservices with Go
## Пользователи
Есть три пользователя:
* login: admin password: admin
* login: alice password: alice
* login: bob password: bob
## cURL
Получить токен
```
curl -v -X POST http://localhost:8080/login?login=<login>&password=<password>
```
Использовать токен
```
curl -v -H "Authorization: <token>" -X POST http://localhost:8080/update

```
