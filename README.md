# Nervatura server-side Go components
![Coverage](https://img.shields.io/badge/Coverage-100.0%25-brightgreen)
[![Go Report Card](https://goreportcard.com/badge/github.com/nervatura/component)](https://goreportcard.com/report/github.com/nervatura/component)
[![GoDoc](https://godoc.org/github.com/nervatura/component?status.svg)](https://pkg.go.dev/github.com/nervatura/component/pkg/component)
[![Release](https://img.shields.io/github/v/release/nervatura/component)](https://github.com/nervatura/component/releases)

An easy way to create a server-side component in any programming language

### Documentation

- [Benefits of server-side components](https://nervatura.github.io/component/#benefits)
- [HTTP request management](https://nervatura.github.io/component/#request_management)
- [Server-side component events](https://nervatura.github.io/component/#events)
- [Creating a server-side component](https://nervatura.github.io/component/#creating)
- [Examples and demo application](https://nervatura.github.io/component/#examples)
- Go package documentation:  
[![GoDoc](https://godoc.org/github.com/nervatura/component?status.svg)](https://pkg.go.dev/github.com/nervatura/component/pkg/component)

### Quick start (demo application)

The demo application displays all components with their test data. Applications can store component data in memory, but they can save it anywhere in json format and load it back. The demo application can store session data in memory and as session files. The source code of the example application also contains an example of using a session database (sqlite3, postgres, mysql, mssql). If you want to use a database session, uncomment before importing the database driver you want to use.

**1. Prebuild binaries**
- [Linux x64](https://github.com/nervatura/component/releases/latest/download/component_linux_x86_64.tar.gz), 
[Linux arm](https://github.com/nervatura/component/releases/latest/download/component_linux_arm64.tar.gz)
- [Windows x64](https://github.com/nervatura/component/releases/latest/download/component_windows_x86_64.zip)
- [MacOS x64](https://github.com/nervatura/component/releases/latest/download/component_darwin_x86_64.tar.gz), 
[MacOS arm](https://github.com/nervatura/component/releases/latest/download/component_darwin_arm64.tar.gz)

**2. Docker file**
- Clone the repository: 
  ```bash
  git clone https://github.com/nervatura/component.git
  ```
  ```bash
  cd component
  ```
- Docker build
  ```bash
  docker build -t component .
  ```
- Run the demo application
  ```bash
  docker run -i -t --rm --name component -p 5000:5000 -v $(pwd)/session:/session component:latest
  ```

**3. Build the project**
- Clone the repository: 
  ```bash
  git clone https://github.com/nervatura/component.git
  ```
  ```bash
  cd component
  ```
- Ensure that you have Golang installed on your system. If not, please follow the [official installation guide](https://golang.org/doc/install).
- Build the project:
  ```bash
  go build -ldflags="-w -s -X main.version=demo" -o ./component main.go
  ```
- Run the demo application
  ```bash
  ./component 5000
  ```

The demo application can store session data in memory and as
session files or session database:
- open the http://localhost:5000/ (memory session) 
- or http://localhost:5000/session (file or database session)

The [Nervatura Client](https://github.com/nervatura/nervatura) interface of the application is another example of the use of server-side components (session and JWT token, database session and more).