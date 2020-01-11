#!/bin/python
import os
import dwdweather as dwd
import pika
import dataclasses
import json
import datetime
import requests
import sys


# url = "https://apiserver.lab.oliviermichaelis.dev/measurements/climate"
url = "http://localhost:8080/measurements/climate"


@dataclasses.dataclass
class SensorValues:
    timestamp:      str
    station:        str
    temperature:    float
    humidity:       float
    pressure:       float


class GoodEncoder(json.JSONEncoder):
    def default(self, o):
        if dataclasses.is_dataclass(o):
            converted = dataclasses.asdict(o)
            converted["station"] = str(converted["station"])
            converted["timestamp"] = str(converted["timestamp"])
            return converted
        return super().default(o)


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
    delta = datetime.timedelta(hours=2)
    time = datetime.datetime(time.year, time.month, time.day, time.hour, time.minute - (time.minute % 10)) - delta

    results = []
    result = client.query(station_id=433, timestamp=time)
    while result is not None:
        results.append(result)
        delta = datetime.timedelta(minutes=10)
        time += delta
        result = client.query(station_id=433, timestamp=time)

    print("Retrieved " + str(len(results)) + " measurements")

    return results


def parse_values(values: dict):
    if values is None:
        return None

    return SensorValues(timestamp=str(values["datetime"]) + "00",    # this is needed for minutes
                        station=values["station_id"],
                        temperature=values["airtemp_temperature_200"],
                        humidity=values["airtemp_humidity"],
                        pressure=values["airtemp_pressure_station"])


def transmit(measurement: SensorValues):
    if not bool(measurement):
        return

    body_json = json.dumps(measurement, cls=GoodEncoder)

    response = None
    try:
        response = requests.post(url=url, data=body_json)
        if response.status_code != 200:
            print("Didn't get http status code 200 as response")
            sys.exit(2)
    finally:
        response.close()


def main():
    for measurement in retrieve_data():
        transmit(parse_values(measurement))


if __name__ == "__main__":
    main()
