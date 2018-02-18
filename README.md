# House Party


### Build instructions

```bash
# ws server
dep ensue
go run main.go

# react server
cd web/
npm i
npm run build
npm i -g serve
serve -s build/
```