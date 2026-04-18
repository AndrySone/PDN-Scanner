#!/usr/bin/env bash
set -euo pipefail

SUMMARY_JSON="${1:-report.summary.json}"
JSONL="${2:-report.jsonl}"
OUT_MD="${3:-final_report.md}"

TOTAL=$(jq -r '.files_total' "$SUMMARY_JSON")
PARSED=$(jq -r '.files_parsed' "$SUMMARY_JSON")
FAILED=$(jq -r '.files_failed' "$SUMMARY_JSON")
STARTED=$(jq -r '.started_at' "$SUMMARY_JSON")
FINISHED=$(jq -r '.finished_at' "$SUMMARY_JSON")
SCAN_ID=$(jq -r '.scan_id' "$SUMMARY_JSON")

UZ_DIST=$(jq -r '.uz_level' "$JSONL" | sort | uniq -c | sort -nr)
STATUS_DIST=$(jq -r '.status' "$JSONL" | sort | uniq -c | sort -nr)
ERR_CAT_TOP=$(jq -r '.error_categories[]?' "$JSONL" | sort | uniq -c | sort -nr | head -20)
TYPE_TOP=$(jq -r '.findings[]?.type' "$JSONL" | sort | uniq -c | sort -nr | head -20)

cat > "$OUT_MD" <<EOF
# Финальный отчёт PII Scanner

## 1. Общая сводка
- Scan ID: \`$SCAN_ID\`
- Started: \`$STARTED\`
- Finished: \`$FINISHED\`
- Files total: **$TOTAL**
- Parsed (ok + partial): **$PARSED**
- Failed: **$FAILED**

## 2. Распределение по UZ
\`\`\`
$UZ_DIST
\`\`\`

## 3. Распределение по статусам
\`\`\`
$STATUS_DIST
\`\`\`

## 4. Топ категорий ошибок
\`\`\`
$ERR_CAT_TOP
\`\`\`

## 5. Топ типов найденных ПДн
\`\`\`
$TYPE_TOP
\`\`\`

## 6. Примечания
- Статусы:
  - **ok**: без ошибок
  - **partial**: найдены ПДн, но есть warning/ошибки
  - **failed**: нет находок и есть ошибки
- Категории ошибок нормализованы для аналитики и мониторинга.
EOF

echo "Generated: $OUT_MD"