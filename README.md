# Nervatura server-side Go components

### Benefits of server-side components

Component based development is an approach to software development that focuses on the design and development
of reusable components. Server components are also reusable bits of code, but are compiled into HTML before
the browser sees them. The server-side components tend to perform better. When the page that a browser receives
contains everything it needs for presentation, it’s going to be able to deliver that presentation to the user
much quicker.

- The development of a client-side application and component takes place in a very complex ecosystem. An average node_modules size can be hundreds of MB and contain hundreds or even over a thousand different packages. Each of these also means external dependencies of varying quality and reliability, which is also a big security risk. In addition, the constant updating and tracking of these different packages and the dozens of frameworks and technologies based on them requires a lot of resources.
A server-side component has <span style='color: green;font-weight: bold;'>no external dependencies</span>. They can be easily created within the technical capabilities of a given server-side language. Their maintenance needs are limited to their actual code, which is very small and much safer due to the lack of external dependencies.

- The language of the client-side components is basically javascript, but most are server-side languages, so go is a much more efficient and safer programming language. JavaScript is originally an add-on to html code and browsers, which was originally created to increase the efficiency of the user interface and not to develop the codebase of complex programs. It is possible to partially replace it during development with, for example, the typescript language, but this also means additional dependencies and the complexity of the development ecosystem, the end result of which will be a javascript code base. This practically means that code written in a programming language is translated into the code of another language and the content to be displayed is created during its execution. There are many intermediate steps, used resources, potential for errors, security risks and uncertainties in the process. With the server-side components, it is possible to simply <span style='color: green;font-weight: bold;'>write the program code in an easy-to-use and safe language</span>, the end result of which is the html content to be displayed.

- Client-side components usually communicate with the server using a JSON-based REST API and receive the data to be displayed. This also means that the data retrieval must adapt to the data structure of the REST API, so the database data must first be converted to this structure, and then reprocessed on the client side for final display. In addition to possible changes to the data structure, this also means JSON encoding and decoding in all cases. The server-side components <span style='color: green;font-weight: bold;'>can directly access the database</span> and use the data immediately in the data structure to be displayed. This also means <span style='color: green;font-weight: bold;'>faster rendering and better resource management</span> for the server-side components.

### Nervatura components

Server components can be written in any server-side language. This enables you to write your client in the
same language as your server application’s logic.
On the user side, an application that is loaded in the browser in html syntax is a set of components
that are hierarchically related to each other. Any component of the application may be able to send a
request to the server, and depending on the processing of the request, any part of the application may
change. The entire page is not replaced or reloaded in the browser, only the required parts of the application.
The components do not use json data format to send data, all data is sent in URL-encoded form. All data of the
application is stored on the server, and the components do not contain javascript code.

- Nervatura components use the htmx library for direct communication with the server. Htmx is small (~14k),
dependency-free, browser-oriented javascript library that allows you to access modern browser
features directly from HTML, rather than using javascript. The server-side components use only a small part
of the possibilities of htmx. More information about htmx can be found on the https://htmx.org link.
The *Application* component contains and automatically loads the appropriate version of htmx when used.

- <span style='color: green;font-weight: bold;'>Nervatura components are not a framework, they use only the built-in packages of go and have no external
dependencies.</span> It is a library of components whose elements can be freely combined with each other and can
be easily further developed. A Nervatura component is actually just <span style='color: green;font-weight: bold;'>a code implementation proposal that
anyone can easily create a server-side component in any program language</span>.

### Documentation and examples

[![GoDoc](https://godoc.org/github.com/nervatura/component?status.svg)](https://pkg.go.dev/github.com/nervatura/component/pkg/component)

### Quick start (demo application)

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
session files:
- open the http://localhost:5000/ (memory session) 
- or http://localhost:5000/session (file session)
