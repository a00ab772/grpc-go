FROM confluentinc/cp-kafka-connect:7.5.0

# Install connectors via confluent-hub
RUN confluent-hub install --no-prompt confluentinc/kafka-connect-jdbc:latest && \
    confluent-hub install --no-prompt confluentinc/kafka-connect-elasticsearch:latest && \
    confluent-hub install --no-prompt confluentinc/kafka-connect-http:latest && \
    confluent-hub install --no-prompt debezium/debezium-connector-mongodb:2.4.2 && \
    confluent-hub install --no-prompt debezium/debezium-connector-mysql:2.4.2 && \
    confluent-hub install --no-prompt debezium/debezium-connector-postgresql:2.4.2 && \
    confluent-hub install --no-prompt redis/redis-kafka-connect:latest

# Install MongoDB Kafka connector manually
RUN mkdir -p /usr/share/confluent-hub-components/mongodb-kafka-connector/lib && \
    curl -fsSL "https://repo1.maven.org/maven2/org/mongodb/kafka/mongo-kafka-connect/1.13.0/mongo-kafka-connect-1.13.0-all.jar" \
      -o /usr/share/confluent-hub-components/mongodb-kafka-connector/lib/mongo-kafka-connect.jar && \
    echo '{"name":"mongodb-kafka-connector","version":"1.13.0","title":"MongoDB Kafka Connector","description":"MongoDB sink and source connector"}' \
      > /usr/share/confluent-hub-components/mongodb-kafka-connector/manifest.json

# Install Redis Kafka connector with ALL required dependencies
# Using a consistent Netty version that matches what Lettuce expects
RUN mkdir -p /usr/share/confluent-hub-components/kafka-connect-redis/lib

RUN REDIS_LIB=/usr/share/confluent-hub-components/kafka-connect-redis/lib && \
    MAVEN=https://repo1.maven.org/maven2 && \
    NETTY_VERSION=4.1.94.Final && \
    \
    # Main connector
    curl -fsSL "$MAVEN/io/github/jaredpetersen/kafka-connect-redis/1.2.3/kafka-connect-redis-1.2.3.jar" -o "$REDIS_LIB/kafka-connect-redis.jar" && \
    \
    # Lettuce Redis client
    curl -fsSL "$MAVEN/io/lettuce/lettuce-core/6.2.6.RELEASE/lettuce-core-6.2.6.RELEASE.jar" -o "$REDIS_LIB/lettuce-core.jar" && \
    \
    # Reactor
    curl -fsSL "$MAVEN/io/projectreactor/reactor-core/3.5.10/reactor-core-3.5.10.jar" -o "$REDIS_LIB/reactor-core.jar" && \
    curl -fsSL "$MAVEN/org/reactivestreams/reactive-streams/1.0.4/reactive-streams-1.0.4.jar" -o "$REDIS_LIB/reactive-streams.jar" && \
    \
    # Netty - consistent version set
    curl -fsSL "$MAVEN/io/netty/netty-common/$NETTY_VERSION/netty-common-$NETTY_VERSION.jar" -o "$REDIS_LIB/netty-common.jar" && \
    curl -fsSL "$MAVEN/io/netty/netty-buffer/$NETTY_VERSION/netty-buffer-$NETTY_VERSION.jar" -o "$REDIS_LIB/netty-buffer.jar" && \
    curl -fsSL "$MAVEN/io/netty/netty-transport/$NETTY_VERSION/netty-transport-$NETTY_VERSION.jar" -o "$REDIS_LIB/netty-transport.jar" && \
    curl -fsSL "$MAVEN/io/netty/netty-codec/$NETTY_VERSION/netty-codec-$NETTY_VERSION.jar" -o "$REDIS_LIB/netty-codec.jar" && \
    curl -fsSL "$MAVEN/io/netty/netty-handler/$NETTY_VERSION/netty-handler-$NETTY_VERSION.jar" -o "$REDIS_LIB/netty-handler.jar" && \
    curl -fsSL "$MAVEN/io/netty/netty-resolver/$NETTY_VERSION/netty-resolver-$NETTY_VERSION.jar" -o "$REDIS_LIB/netty-resolver.jar" && \
    curl -fsSL "$MAVEN/io/netty/netty-resolver-dns/$NETTY_VERSION/netty-resolver-dns-$NETTY_VERSION.jar" -o "$REDIS_LIB/netty-resolver-dns.jar" && \
    curl -fsSL "$MAVEN/io/netty/netty-transport-native-unix-common/$NETTY_VERSION/netty-transport-native-unix-common-$NETTY_VERSION.jar" -o "$REDIS_LIB/netty-transport-native-unix-common.jar" && \
    curl -fsSL "$MAVEN/io/netty/netty-transport-native-epoll/$NETTY_VERSION/netty-transport-native-epoll-$NETTY_VERSION.jar" -o "$REDIS_LIB/netty-transport-native-epoll.jar" && \
    \
    echo '{"name":"kafka-connect-redis","version":"1.2.3","title":"Kafka Connect Redis","description":"Redis sink connector"}' \
      > /usr/share/confluent-hub-components/kafka-connect-redis/manifest.json

RUN ls -lh /usr/share/confluent-hub-components/kafka-connect-redis/lib/