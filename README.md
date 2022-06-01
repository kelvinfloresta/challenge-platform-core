# Conformity Core

...

![Deploy Production](https://github.com/venturaac/conformity-core/actions/workflows/deploy.yml/badge.svg)

# Useful commands

## Watch mode

Install CompileDaemon

```sh
go get github.com/githubnemo/CompileDaemon
```

Running CompileDaemon

```sh
CompileDaemon -command="./conformity-core.exe"
```

## Running tests

```sh
go clean -testcache && go test ./... -p 1
```

## Add collors to test output

```sh
go test -v ./... -p 1 |
grep -v RUN |
grep -v time= |
sed ''/^PASS/s//$(printf "")/'' |
sed ''/PASS/s//$(printf "\033[32mPASS\033[0m")/'' |
sed ''/FAIL/s//$(printf "\033[31mFAIL\033[0m")/''
```

## Generate Sigle Sign On (SAML) keys

```sh
openssl req -x509 -newkey rsa:2048 -keyout conformity-core.key -out conformity-core.cert -days 365 -nodes -subj "/CN=myservice.example.com"
```
