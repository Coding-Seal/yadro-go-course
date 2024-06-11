# yadro-go-course
Course from Yadro on microservices with Go
![video](example.GIF)
## Пользователи
Есть три пользователя:
* login: admin password: admin
* login: alice password: alice
* login: bob password: bob
## cURL
Получить токен
```
curl -X POST -v --data '{"login":"bob", "password":"bob"}' localhost:8080/api/login
```
Использовать токен
```
curl -v -H "Authorization: <token>" -X POST http://localhost:8080/api/update

```
## API
Все эндпоинты для работы с REST API имеют префикс /api
Специально для 9 задания

## Test Coverage
![Cover.svg](test/out/cover.svg)
