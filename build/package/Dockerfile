FROM python:3.8-alpine

RUN pip3 install dwdweather2 pika

RUN addgroup -S producer && adduser -S producer -G producer
USER producer

WORKDIR /app
COPY build/package/entrypoint.py /app/

CMD python /app/entrypoint.py
