# Example

## Tasks

### generate-go-from-css

```
go run ../cmd/. generate --file-name=./css/styles.css --package="css" > ./css/styles.go
```

### run

```
templ generate --watch --cmd="go run ." --proxy="http://localhost:8080"
```
