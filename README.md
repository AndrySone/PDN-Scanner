## Запуск проекта

Все зависимости присутствует в файле `requirements.txt`
Датасет уже находится в корне проекта, поэтому перед запуском убедитесь, что путь к нему корректный: `./ПДнDataset/share`.

### 1) Запустить Python API (PDN-Scanner)

```bash
uvicorn python.app.main:app --host 0.0.0.0 --port 8081
```

API будет доступен по адресу: `http://127.0.0.1:8081`.

### 2) Запустить Go CLI-сканер

В новом терминале выполните:

```bash
go run ./cmd/cli \
  -root ./ПДнDataset/share \
  -out report \
  -py-url http://127.0.0.1:8081 \
  --shm-path /tmp/pii_shm.dat \
  --shm-size-mb 256
```

### Что делают параметры CLI

- `-root` — путь к директории с датасетом.
- `-out` — имя/путь выходной директории (например, `report`).
- `-py-url` — URL запущенного Python API.
- `--shm-path` — путь к файлу shared memory.
- `--shm-size-mb` — размер shared memory в МБ.

### Ожидаемый результат

После завершения работы CLI в директории `report` будет сформирован отчёт по результатам сканирования.
