FROM alpine:3.5

COPY events-service /app/app-service
COPY cert.pem /app/cert.pem
COPY key.pem /app/key.pem
RUN chmod +x /app/app-service
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
RUN chown -R appuser:appgroup /app

USER appuser
WORKDIR /app

ENV LISTEN_URL=0.0.0.0:8181
EXPOSE 8181

ENTRYPOINT ["/app/app-service"]