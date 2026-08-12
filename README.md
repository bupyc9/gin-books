# Books

REST API сервис книг и авторов, написан на golang.

Библиотеки:
* http framework - [gin](https://github.com/gin-gonic/gin)
* ORM - [gorm](https://github.com/go-gorm/gorm)

Реализовано:
* сущности author и book
* RESTful API:
  * CRUD для author
  * CRUD для book

Команда для запуска тестов `DB_HOST=database-test go test ./...`. Выполнять в докер контейнере `backend`. 