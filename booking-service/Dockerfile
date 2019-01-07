FROM alpine:3.5

COPY booking-service /app/app-service
RUN chmod +x /app/app-service
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
RUN chown -R appuser:appgroup /app

USER appuser
WORKDIR /app

ENV LISTEN_URL=0.0.0.0:8282
EXPOSE 8282

ENTRYPOINT ["/app/app-service"]