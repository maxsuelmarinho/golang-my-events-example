FROM alpine:3.5

COPY events-service /app/app-service
RUN chmod +x /app/app-service

EXPOSE 8181

ENTRYPOINT /app/app-service