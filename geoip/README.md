ip query, supported lite2mmdb(default), ipcc, iptaobao, ipapi, ipinfo or ip2location

# Build

go build -o iploc main.go

# Usage

## Help

```shell
Usage of ./iploc:
      --disable-auto-merge-cells   disable auto merge cells
      --disable-my-external        disable show my external ip address
  -g, --geo string                 use sourcees for ip geo, eg: lite2mmdb,ipcc,iptaobao,ipapi,ipinfo or ip2location (default "lite2mmdb")
      --json                       output with json format, default table
      --row-line                   show row line (default true)
      --set-align-center           set column align center
pflag: help requested

```

## lite2mmdb

```shell
./iploc 202.106.0.20
+--------------+--------------+---------+-----------+------+--------------------------------+-----+----------------+
|      IP      | COUNTRY CODE | COUNTRY | CONTINENT | CITY |       TIMEZONE/LOCATION        | ZIP |    COMMENT     |
+--------------+--------------+---------+-----------+------+--------------------------------+-----+----------------+
| 46.3.240.205 | RU           | Russia  | Europe    |      | Europe/Moscow(55.7386,37.6068) |     | My External IP |
+--------------+--------------+---------+-----------+------+--------------------------------+-----+----------------+
| 202.106.0.20 | CN           | China   | Asia      |      | Asia/Shanghai(34.7732,113.722) |     |                |
+--------------+--------------+---------+-----------+------+--------------------------------+-----+----------------+
|                                                                                                      TOTAL: 2    |
+--------------+--------------+---------+-----------+------+--------------------------------+-----+----------------+
```

## ipinfo

```shell
./iploc 202.106.0.20 --geo ipinfo
+--------------+--------------------+---------+-----------+-----------+--------------------------------+----------------------------------+-----+----------------+
|      IP      |      HOSTNAME      | COUNTRY |  REGION   |   CITY    |              ORG               |        TIMEZONE/LOCATION         | ZIP |    COMMENT     |
+--------------+--------------------+---------+-----------+-----------+--------------------------------+----------------------------------+-----+----------------+
| 46.3.240.205 |                    | HK      | Hong Kong | Hong Kong | AS38136 Akari Networks         | Asia/Hong_Kong(22.2783,114.1747) |     | My External IP |
+--------------+--------------------+---------+-----------+-----------+--------------------------------+----------------------------------+-----+----------------+
| 202.106.0.20 | gjjline.bta.net.cn | CN      | Beijing   | Beijing   | AS4808 China Unicom Beijing    | Asia/Shanghai(39.9075,116.3972)  |     |                |
|              |                    |         |           |           | Province Network               |                                  |     |                |
+--------------+--------------------+---------+-----------+-----------+--------------------------------+----------------------------------+-----+----------------+
|                                                                                                                                                    TOTAL: 2    |
+--------------+--------------------+---------+-----------+-----------+--------------------------------+----------------------------------+-----+----------------+
```
