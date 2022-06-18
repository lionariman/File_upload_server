# File_upload_server

## Description:

Данный сервер принимает файл, делит его на части по мегабайту
и сохраняет у себя. Если какая-то часть получается меньше 1 мб,
то дописывает туда байты (нулями), чтобы все части файла
не превышали 1 мб, даже если сам файл весит несколько КБ,
и сохраняет у себя.
Когда пользователь запрашивает ранее загруженный файл, 
сервер собирает его по частям (по мегабайтам) и возвращает его целым,
в том же размере в котором был получаен,
то есть убирает ранее дописанные байты (нули).

### Endpoints:

```/upload_file``` <br>
принимает файл

```/get_filfe:название файла``` <br>
возвращает файл по названию

```/delete_file:название файла``` <br>
удаляет файл на сервере по названию

```/delete_all_files``` <br>
удаляет все загруженные файлы с сервера

### Examples:

#### Run server:

```port :8080``` <br>
```go run main.go```

#### /upload_file:
```curl -X POST http://localhost:8080/upload_file -F "file=@fileName.txt"```

#### /get_file
```curl http://localhost:8080/get_file:fileName > newFileName.txt```

#### /delete_file
```curl http://localhost:8080/delete_file:fileName```

#### /delete_all_files
```curl http://localhost:8080/delete_all_files```