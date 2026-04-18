# PDN-Scanner

``` bash
uvicorn python.app.main:app --host 0.0.0.0 --port 8081
```

```bash
go run ./cmd/cli \
  -root ./ПДнDataset/share \
  -out report \
  -py-url http://127.0.0.1:8081 \
  --shm-path /tmp/pii_shm.dat \
  --shm-size-mb 256
```

