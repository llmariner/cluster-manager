# cluster-manager
Cluster Manager


## Running Locally

```bash
make build-server
./bin/server run --config config.yaml
```

`config.yaml` has the following content:

```yaml
httpPort: 8080
grpcPort: 8081
internalGrpcPort: 8082

debug:
  standalone: true
  sqlitePath: /tmp/cluster_manager.db
```

You can then connect to the DB.

```bash
sqlite3 /tmp/cluster_manager.db
```

You can then hit the endpoint.

```bash
curl http://localhost:8080/v1/clusters

curl -X POST http://localhost:8080/v1/clusters -d '{"name": "my-cluster"}'
```


## Running with Docker Compose

Run the following command:

```bash
docker-compose build
docker-compose up
```
