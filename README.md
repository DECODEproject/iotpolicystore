# iotpolicystore

Implementation of the policy store component for the DECODE IoT pilot.

## Building

Run `make` or `make build` to build our binary compiled for `linux/amd64`
with the current directory volume mounted into place. This will store
incremental state for the fastest possible build. To build for `arm` or
`arm64` you can use: `make build ARCH=arm` or `make build ARCH=arm64`. To
build all architectures you can run `make all-build`.

Run `make container` to package the binary inside a container. It will
calculate the image tag based on the current VERSION (calculated from git tag
or commit - see `make version` to view the current version). To build
containers for the other supported architectures you can run
`make container ARCH=arm` or `make container ARCH=arm64`. To make all
containers run `make all-container`.

Run `make push` to push the container image to `REGISTRY`, and similarly you
can run `make push ARCH=arm` or `make push ARCH=arm64` to push different
architecture containers. To push all containers run `make all-push`.

Run `make clean` to clean up.

## Testing

To run the test suite, use the make task: `test`. This will run all testcases
inside a containerized environment, generating a test coverage report that
can be found in `.coverage/coverage.html`.

## Using this image

As the server requires access to a postgres DB instance to persist data the
simplest way to run the image locally is via docker-compose. An example
compose file is shown below:

```yaml
version: '3'

services:
    postgres:
        image: postgres:10-alpine
        ports:
          - "5432"
        restart: always
        volumes:
          - postgres_vol: /var/lib/postgresql/data
        environment:
          - POSTGRES_PASSWORD=password
          - POSTGRES_USER=decode
          - POSTGRES_DB=postgres

    policystore:
        image: thingful/policystore-amd64:v0.0.1
        ports:
          - "8082:8082"
        restart: always
        environment:
          - POLICYSTORE_DATABASE_URL=postgres://decode:password@postgres:5432/postgres?sslmode=disable
          - POLICYSTORE_ENCRYPTION_PASSWORD=secret-password-changeme
          - POLICYSTORE_HASHID_SALT=hashid-salt-changeme
        depends_on:
          - postgres
        command: [ "server", "--verbose" ]

volumes:
  postgres_vol:
```

## License

In accordance with the rules defined for the DECODE consortium this project
is licensed under the terms of the GNU Affero General Public License. Please
see the LICENSE file in the repository root for details.