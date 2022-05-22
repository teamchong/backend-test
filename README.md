[![godoc for benthosdev/benthos][godoc-badge]][godoc-url]
![Coverage](https://img.shields.io/badge/Coverage-75.6%25-brightgreen)
[![Build Status][actions-badge]][actions-url]

# ETL (extract, transform and load)

![ETL](https://www.benthos.dev/img/what-is-blob.svg)
Use [Benthos](https://www.benthos.dev) to solve the [ETL challenge](https://databricks.com/glossary/extract-transform-load)

---

## Source

```yaml
label: "random_user"
http_client:
    url: "https://random-data-api.com/api/users/random_user"
    headers:
        User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:100.0) Gecko/20100101 Firefox/100.0
        Accept: application/json
    tls:
        enabled: true
        skip_cert_verify: false
    rate_limit: random
```

[Other build-in sources][sources-link]

## Benthos rate limits plugin

### random

The random rate limit is X every (Y1 - Y2)ms rate limit that can be shared across any number of components within the pipeline but does not support distributed rate limits across multiple running instances of Benthos.

```yaml
# Config fields, showing default values
label: ""
random:
  count: 1
  min_interval: 250ms
  max_interval: 720ms
```

#### Fields

##### count

The maximum number of requests to allow for a given period of time.

Type: int
Default: `1`

##### min_interval

The minmum time window to limit requests by.

Type string
Default: `250ms`

##### max_interval

The maxmum time window to limit requests by.

Type string
Default: `750ms`

---

## Transform

```yaml
pipeline:
    threads: 4
    processors:
        - label: transform
          bloblang: |
              root = {}
              root.id = this.id
              root.first_name = this.first_name
              root.last_name = this.last_name
              root.date_of_birth = this.date_of_birth
              root.city = this.address.city
              root.street_name = this.address.street_name
              root.street_address = this.address.street_address
              root.zip_code = this.address.zip_code
              root.state = this.address.state
              root.country = this.address.country
              root.lat = this.address.coordinates.lat.number()
              root.lng = this.address.coordinates.lng.number()
```

[Other build-in processors][processors-link]

---

## Sink

```yaml
output:
    label: postgres
    sql_raw:
        driver: postgres
        dsn: ${BENTHOS_DSN}
        query: >
            INSERT INTO random_user (id, first_name, last_name, date_of_birth, city, street_name, street_address, zip_code, state, country, lat, lng)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
            ON CONFLICT(id) DO UPDATE SET
              first_name=$2, last_name=$3, date_of_birth=$4, city=$5, street_name=$6, street_address=$7, zip_code=$8, state=$9, country=$10, lat=$11, lng=$12;
        args_mapping: |
            root = [
              this.id,
              this.first_name,
              this.last_name,
              this.date_of_birth,
              this.city,
              this.street_name,
              this.street_address,
              this.zip_code,
              this.state,
              this.country,
              this.lat,
              this.lng
            ]
        batching:
            period: 1s
```

[Other build-in sinks][sinks-link]

## Author

[Steven Chong](https://github.com/teamchong)

[godoc-badge]: https://pkg.go.dev/badge/github.com/benthosdev/benthos/v4/public
[godoc-url]: https://pkg.go.dev/github.com/benthosdev/benthos/v4/public
[actions-badge]: https://github.com/teamchong/backend-test/actions/workflows/test.yaml/badge.svg
[actions-url]: https://github.com/teamchong/backend-test/actions/workflows/test.yaml
[sources-link]: https://www.benthos.dev/docs/components/inputs/about
[processors-link]: https://www.benthos.dev/docs/components/processors/about
[sinks-link]: https://www.benthos.dev/docs/components/outputs/about
