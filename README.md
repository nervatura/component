# Nervatura server-side Go components

### Quick start (demo application)

1. 💻 Ensure that you have Golang installed on your system. If not, please follow the [official installation guide](https://golang.org/doc/install).
    
2. 📦 Clone the repository:
    
```bash
git clone https://github.com/nervatura/component.git
```

3. 📂 Change into the project directory:

```bash
cd component
```

4. 🔨 Build the project:

```bash
go build -ldflags="-w -s -X main.version=demo" -o ./demo_app main.go
```

5. 🌍 Run the demo application:

```bash
./demo_app
```
and open the http://localhost:5000/

More documentation coming soon...