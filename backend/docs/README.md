## Getting started

1. Add comments to your API source code, See [Declarative Comments Format]
2. Download swag by using: `go install github.com/swaggo/swag/cmd/swag@latest`
3. If the default terminal is zsh, open the file `~/.zshrc` and add the line `export PATH=$PATH:$HOME/go/bin` to the end of the file.
4. Run `swag init` in the project's root folder which contains the `main.go` file. This will parse your comments and generate the required files (`docs` folder and `docs/docs.go`).
5. (optional) Use `swag fmt` format the SWAG comment. (Please upgrade to the latest version)
6. Open Swagger : `http://localhost:8080/swagger/index.html#/`
