Пример сборки и запуска:
```(bash)
GROUP_UDP_ADDR="224.0.0.1:9000" docker compose up --build -d
```

Подключиться к прокси для тестирования:
```(bash)
nc -u localhost 2000
```

Прокси обрабатывает два типа команд:
```(bash)
get_result <native/xml/json/protobuf/avro/yaml/message_pack>
```
```(bash)
get_result all
```

Внутренние детали:
* образы контейнеров скачиваются из докерхаба
* контейнер прокси слушает 2000 порт, который прокинут на хост. После получения команды, прокси запрашивает бенчмарки и результат возвращает на клиент
* контейнеры, которые проводят бенчмарк, прослушивают два адреса: уникальный адрес в сети и мультикаст адрес. Для обработки двух адресов запускаются две горутины
* все контейнеры объединены в одну bridge сеть(все они могут общаться друг с другом по <имя сервиса>:<прослушиваемый порт>)
* есть возможность задать мультикаст адрес через переменную окружения GROUP_UDP_ADDR. При инициализации контейнера, проверяется корректность адреса.
* есть возможность для контейнеров выбрать уровень логирования, указав значение APP_ENVIRONMENT в файле docker-compose.yaml:
  * если указать debug и testing, то логи будут выводиться до уровня DEBUG
  * если указать pre_prod и prod, то логи будут выводиться до уровня INFO

Иллюстрация подключения к прокси и получения замеров сериализации и десериализации:
```(bash)
❯ nc -u localhost 2000
get_result xml
xml:
	size of serialized data  	225 bytes
	serialization duration   	52µs
	deserialization duration 	32µs
get_result all
xml:
	size of serialized data  	225 bytes
	serialization duration   	49µs
	deserialization duration 	41µs
avro:
	size of serialized data  	44 bytes
	serialization duration   	72µs
	deserialization duration 	49µs
message_pack:
	size of serialized data  	84 bytes
	serialization duration   	23µs
	deserialization duration 	18µs
json:
	size of serialized data  	115 bytes
	serialization duration   	77µs
	deserialization duration 	27µs
protobuf:
	size of serialized data  	132 bytes
	serialization duration   	269µs
	deserialization duration 	18µs
native:
	size of serialized data  	147 bytes
	serialization duration   	182µs
	deserialization duration 	75µs
yaml:
	size of serialized data  	131 bytes
	serialization duration   	106µs
	deserialization duration 	105µs
```
