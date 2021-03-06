#!/bin/bash

# MySQL driver
go get -u github.com/go-sql-driver/mysql

# HTTP routing
go get -u github.com/gorilla/mux

# HTTP secure cookies
go get -u github.com/gorilla/securecookie

# Password hashing
go get -u golang.org/x/crypto/bcrypt

# Logging
go get -u go.uber.org/zap

# UUID
go get github.com/google/uuid

# unit testing
go get github.com/golang/mock/gomock
go get github.com/golang/mock/mockgen
go get github.com/DATA-DOG/go-sqlmock
