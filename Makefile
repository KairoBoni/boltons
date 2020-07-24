kafka:
	MY_IP=172.17.0.1 docker-compose -f docker-compose-kafka.yml up --build

workers:
	MY_IP=172.17.0.1 docker-compose up --build

topic:
	docker run --net=host --rm confluentinc/cp-kafka:5.0.0 kafka-topics --create --topic foo --partitions 4 --replication-factor 2 --if-not-exists --zookeeper localhost:32181

kafcat:
	kafkacat -C -b localhost:19092,localhost:29092,localhost:39092 -t worker -p 0
