FROM confluentinc/cp-kafka-connect:7.5.0

# Install connectors via confluent-hub
RUN confluent-hub install --no-prompt confluentinc/kafka-connect-jdbc:latest && \
    confluent-hub install --no-prompt confluentinc/kafka-connect-elasticsearch:latest && \
    confluent-hub install --no-prompt confluentinc/kafka-connect-http:latest && \
    confluent-hub install --no-prompt debezium/debezium-connector-mongodb:2.4.2 && \
    confluent-hub install --no-prompt debezium/debezium-connector-mysql:2.4.2 && \
    confluent-hub install --no-prompt debezium/debezium-connector-postgresql:2.4.2

# Install MongoDB Kafka connector manually
RUN mkdir -p /usr/share/confluent-hub-components/mongodb-kafka-connector/lib && \
    curl -fsSL "https://repo1.maven.org/maven2/org/mongodb/kafka/mongo-kafka-connect/1.13.0/mongo-kafka-connect-1.13.0-all.jar" \
      -o /usr/share/confluent-hub-components/mongodb-kafka-connector/lib/mongo-kafka-connect.jar && \
    echo '{"name":"mongodb-kafka-connector","version":"1.13.0","title":"MongoDB Kafka Connector","description":"MongoDB sink and source connector"}' \
      > /usr/share/confluent-hub-components/mongodb-kafka-connector/manifest.json

# Install Redis Kafka connector with ALL required dependencies
RUN mkdir -p /usr/share/confluent-hub-components/kafka-connect-redis/lib && \
    # Main connector JAR
    curl -fsSL "https://repo1.maven.org/maven2/io/github/jaredpetersen/kafka-connect-redis/1.2.3/kafka-connect-redis-1.2.3.jar" \
      -o /usr/share/confluent-hub-components/kafka-connect-redis/lib/kafka-connect-redis.jar && \
    # Lettuce Redis client
    curl -fsSL "https://repo1.maven.org/maven2/io/lettuce/lettuce-core/6.2.4.RELEASE/lettuce-core-6.2.4.RELEASE.jar" \
      -o /usr/share/confluent-hub-components/kafka-connect-redis/lib/lettuce-core.jar && \
    # Reactor Core (required by Lettuce)
    curl -fsSL "https://repo1.maven.org/maven2/io/projectreactor/reactor-core/3.5.4/reactor-core-3.5.4.jar" \
      -o /usr/share/confluent-hub-components/kafka-connect-redis/lib/reactor-core.jar && \
    # Reactive Streams (required by Reactor)
    curl -fsSL "https://repo1.maven.org/maven2/org/reactivestreams/reactive-streams/1.0.4/reactive-streams-1.0.4.jar" \
      -o /usr/share/confluent-hub-components/kafka-connect-redis/lib/reactive-streams.jar && \
    # Netty (required by Lettuce) - all required modules
    curl -fsSL "https://repo1.maven.org/maven2/io/netty/netty-common/4.1.90.Final/netty-common-4.1.90.Final.jar" \
      -o /usr/share/confluent-hub-components/kafka-connect-redis/lib/netty-common.jar && \
    curl -fsSL "https://repo1.maven.org/maven2/io/netty/netty-buffer/4.1.90.Final/netty-buffer-4.1.90.Final.jar" \
      -o /usr/share/confluent-hub-components/kafka-connect-redis/lib/netty-buffer.jar && \
    curl -fsSL "https://repo1.maven.org/maven2/io/netty/netty-transport/4.1.90.Final/netty-transport-4.1.90.Final.jar" \
      -o /usr/share/confluent-hub-components/kafka-connect-redis/lib/netty-transport.jar && \
    curl -fsSL "https://repo1.maven.org/maven2/io/netty/netty-codec/4.1.90.Final/netty-codec-4.1.90.Final.jar" \
      -o /usr/share/confluent-hub-components/kafka-connect-redis/lib/netty-codec.jar && \
    curl -fsSL "https://repo1.maven.org/maven2/io/netty/netty-handler/4.1.90.Final/netty-handler-4.1.90.Final.jar" \
      -o /usr/share/confluent-hub-components/kafka-connect-redis/lib/netty-handler.jar && \
    curl -fsSL "https://repo1.maven.org/maven2/io/netty/netty-resolver/4.1.90.Final/netty-resolver-4.1.90.Final.jar" \
      -o /usr/share/confluent-hub-components/kafka-connect-redis/lib/netty-resolver.jar && \
    curl -fsSL "https://repo1.maven.org/maven2/io/netty/netty-transport-native-unix-common/4.1.90.Final/netty-transport-native-unix-common-4.1.90.Final.jar" \
      -o /usr/share/confluent-hub-components/kafka-connect-redis/lib/netty-transport-native-unix-common.jar && \
    # Manifest
    echo '{"name":"kafka-connect-redis","version":"1.2.3","title":"Kafka Connect Redis","description":"Redis sink connector"}' \
      > /usr/share/confluent-hub-components/kafka-connect-redis/manifest.json