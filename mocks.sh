#!/usr/bin/env bash

COMMAND="$(which $0)"
CMD_HOME="$(dirname "$COMMAND")"

cd "${CMD_HOME}"

mockgen -destination=mocks/mock_dao.go -package=mocks -source=storage/dao.go
mockgen -destination=mocks/mock_configuration.go -package=mocks -source=rest_backend/configuration.go
