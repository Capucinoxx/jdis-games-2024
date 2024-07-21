FROM golang:1.21.6 as wasm

WORKDIR /app

COPY go.mod .
COPY go.sum .
COPY pkg/ pkg/
COPY consts/ consts/
COPY internal/ internal/
COPY interface/wasm/main.go wasm/main.go

RUN GOARCH=wasm GOOS=js go build -o lib.wasm wasm/main.go
RUN cp "$(go env GOROOT)"/misc/wasm/wasm_exec.js .

# ------------------------------
FROM node:18 as build

WORKDIR /app

COPY interface/package.json .
COPY interface/package-lock.json .
COPY interface/tsconfig.json .
COPY interface/index.html .
COPY interface/style.css .
COPY interface/src/ src/

COPY --from=wasm /app/wasm_exec.js src/wasm/wasm_exec.js

RUN npm install
RUN npm run build

# ------------------------------
FROM golang:1.21.6 as server

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY --from=build /app/dist /app/dist
COPY --from=wasm /app/lib.wasm /app/dist/lib.wasm
COPY interface/assets /app/dist/assets

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

RUN go build -o server /app/main.go

EXPOSE 443

ENTRYPOINT ["/entrypoint.sh"]