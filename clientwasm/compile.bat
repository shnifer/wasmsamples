SET GOOS=js
SET GOARCH=wasm
go build -o main.wasm
copy /Y  main.wasm ..\server
