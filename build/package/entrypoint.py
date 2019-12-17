#!/bin/python
import os
import dwdweather as dwd
import pika
import dataclasses
import json
import datetime


@dataclasses.dataclass
class SensorValues:
    timestamp:  str
    station:    str
    temperature: float
    humidity:    float
    pressure:    float


def get_environment(variable: str, default: str) -> str:
    retrieved = os.environ.get(variable)
    if retrieved is None:
        return default

    return retrieved


def retrieve_credentials(secret_path) -> pika.PlainCredentials:
    file_username = open(secret_path + "/username", "r")
    file_password = open(secret_path + "/password", "r")
    credentials = pika.PlainCredentials(file_username.read().rstrip(), file_password.read().rstrip())

    file_username.close()
    file_password.close()

    return credentials


def retrieve_data() -> list:
    client = dwd.DwdWeather(resolution="10_minutes")

    # Create timestamp and floor to multiples of 10 min.
    # The time is decreased by an hour to make sure data is available.
    time = datetime.datetime.utcnow()
    time = datetime.datetime(time.year, time.month, time.day, time.hour - 2, time.minute - (time.minute % 10))

    results = []
    result = client.query(station_id=433, timestamp=time)
    while result is not None:
        results.append(result)
        delta = datetime.timedelta(minutes=10)
        time += delta
        result = client.query(station_id=433, timestamp=time)

    return results


def parse_values(values: dict):
    if values is None:
        return None

    # timestamp = datetime.datetime.strptime(str(body["datetime"]), "%Y%m%d%H%M")
    return SensorValues(timestamp=values["datetime"],
                        station=values["station_id"],
                        temperature=values["airtemp_temperature_200"],
                        humidity=values["airtemp_humidity"],
                        pressure=values["airtemp_pressure_station"])


def publish(rabbit_channel, exchange: str, value: SensorValues):
    if value is None:
        return

    body_json = json.dumps(dataclasses.asdict(value))
    rabbit_channel.basic_publish(exchange=exchange, routing_key="sensor", body=body_json)


def main():
    rabbit_url = get_environment("RABBITMQ_SERVICE_URL", "rabbitmq-ha.default.svc.cluster.local")
    rabbit_port = int(get_environment("RABBITMQ_SERVICE_PORT", "5672"))
    rabbit_queue = get_environment("RABBITMQ_QUEUE", "sensor")
    rabbit_exchange = get_environment("RABBITMQ_EXCHANGE", "sensor")
    rabbit_secret = get_environment("RABBITMQ_SECRET_PATH", "/credentials/rabbitmq")

    parameters = pika.ConnectionParameters(host=rabbit_url,
                                           port=rabbit_port,
                                           credentials=retrieve_credentials(rabbit_secret))
    connection = pika.BlockingConnection(parameters=parameters)
    channel = connection.channel()

    for measurement in retrieve_data():
        publish(channel, rabbit_exchange, parse_values(measurement))

    connection.close()


if __name__ == "__main__":
    main()
