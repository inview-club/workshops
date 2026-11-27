FROM python:3.13-alpine

WORKDIR /app
COPY server.py .

ENV PYTHONUNBUFFERED=1
USER 10001:10001
EXPOSE 8080

CMD ["python", "server.py"]
