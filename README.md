# AlGoRhythm

| Tests |
|:-----:|
| [![Tests](https://github.com/markusewalker/algorhythm/actions/workflows/run-tests.yml/badge.svg?branch=main)](https://github.com/markusewalker/algorhythm/actions/workflows/run-tests.yml) |

AlGoRhythm is a Golang-based web application that tracks the user's top artists and top tracks of the month.

## Table of Contents
1. [Prerequisites](#Prerequisites)
2. [Getting Started](#Getting-Started)
3. [Testing](#Testing)

## Prerequisites
In order to properly run this, you will need the following tools installed and configured on your client machine:

- `go`
- `npm`
- Spotify developer account

Once you have your Spotify developer account, you will need to ensure that you have an app specific to run this project. Once created, take note of the client secret and client ID as you will need it later.

## Getting Started
Before running the command, you will need to create an `.env` file at the root directory. The file needs to have the following:

```
CLIENT_ID=""
CLIENT_SECRET=""
STATE=""
```

For `STATE`, this value is up to the user, but it is recommended to use an encrypted password.

Once you have followed the prerequisites and installed all the needed dependencies, simply run the following command: `go run main.go`. From there, your default browser should kick off the web application against `http://localhost:8080`.

## Running Program
To run the program, you have three options: GUI mode, interactive CLI, or command line arguments.

### Testing
There is frontend and backend testing for the web application. For the backend changes, simply run the following command at the root directory: `go test ./... -v`.

For the frontend changes, once you have npm properly configured with Jest, simply run the following command: `npm test`.