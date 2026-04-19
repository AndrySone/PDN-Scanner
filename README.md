# PDN-Scanner v4.0
Сканер ищет файлы с персональными данными (ПДн) и сохраняет результат в result.csv.

## Запуск проекта

Все зависимости присутствует в файле `requirements.txt`

### Запустить Python (PDN-Scanner)

```bash
ln -s /absolute/path/to/your/data ./DATA
python scanner.py
```

## Настройки
В коде можно изменить:

ROOT_DIR = Path("./DATA") — директория с данными
OUTPUT_CSV = Path("result.csv") — выходной файл
MAX_WORKERS = 8 — число потоков

## Типичный вывод
количество найденных файлов с ПДн
уровни критичности (УЗ-1..УЗ-4)
путь к сохранённому result.csv
