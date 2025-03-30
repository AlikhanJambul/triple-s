# Triple-S

Triple-S (Simple Storage Service) — это REST API, реализующее функционал объектного хранилища по аналогии с Amazon S3.

## Функциональность
- Управление бакетами (создание, удаление, список)
- Загрузка, получение и удаление объектов
- Хранение метаданных в CSV
- Обработка HTTP-запросов в XML-формате

## Требования
- Go 1.22+

## Установка и запуск
```sh
# Клонирование репозитория
git clone https://github.com/AlikhanJambul/triple-s.git
cd triple-s

# Сборка
go build -o triple-s

# Запуск
./triple-s
```

## Использование
### Создание бакета
```sh
curl -X PUT http://localhost:8080/bucket/my-bucket
```

### Загрузка объекта
```sh
curl -X PUT --data-binary @file.txt http://localhost:8080/bucket/my-bucket/object/file.txt
```

### Получение объекта
```sh
curl -X GET http://localhost:8080/bucket/my-bucket/object/file.txt -o file.txt
```

### Удаление объекта
```sh
curl -X DELETE http://localhost:8080/bucket/my-bucket/object/file.txt
```

### Удаление бакета
```sh
curl -X DELETE http://localhost:8080/bucket/my-bucket
```

## Формат запросов и ответов
- Все ответы формируются в XML-формате.
- Метаданные объектов хранятся в CSV-файле.

## Лицензия
MIT

