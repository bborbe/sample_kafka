# Run Samples

## Setup Kafka

```bash
echo "127.0.0.1 kafka" >> /etc/hosts
```

```bash
docker network create kafka
```

```bash 
docker kill zookeeper
docker rm zookeeper
docker run \
--net=kafka \
--name=zookeeper \
-e ZOOKEEPER_CLIENT_PORT=2181 \
confluentinc/cp-zookeeper:5.0.0
```

```bash 
docker kill kafka
docker rm kafka
docker run \
--net=kafka \
--name=kafka \
-p 9092:9092 \
-e KAFKA_ZOOKEEPER_CONNECT=zookeeper:2181 \
-e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://kafka:9092 \
-e KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1 \
confluentinc/cp-kafka:5.0.0
```

## Run Avro Producer

```
go run avro/avro_producer/avro_producer.go
```

open http://127.0.0.1:8001

## Run Avro Consumer

```
go run avro/avro_consumer/avro_consumer.go
```

open http://127.0.0.1:8002
