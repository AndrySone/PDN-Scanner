# Финальный отчёт PII Scanner

## 1. Общая сводка
- Scan ID: `scan-1776516083`
- Started: `2026-04-18T15:41:23.337055109+03:00`
- Finished: `2026-04-18T15:45:09.767716811+03:00`
- Files total: **3220**
- Parsed (ok + partial): **3015**
- Failed: **205**

## 2. Распределение по UZ
```
   2517 NO_PD
    413 UZ-3
    211 UZ-4
     79 UZ-2
```

## 3. Распределение по статусам
```
   2503 ok
    512 
    205 failed
```

## 4. Топ категорий ошибок
```
    196 file_open_error
      5 ocr_error
      4 unsupported_format
```

## 5. Топ типов найденных ПДн
```
    445 passport_rf
    438 inn
    263 email
    225 phone
    118 card_pan
      8 bank_account
      6 snils
      3 bik
      1 mrz
      1 birth_date
```

## 6. Примечания
- Статусы:
  - **ok**: без ошибок
  - **partial**: найдены ПДн, но есть warning/ошибки
  - **failed**: нет находок и есть ошибки
- Категории ошибок нормализованы для аналитики и мониторинга.
