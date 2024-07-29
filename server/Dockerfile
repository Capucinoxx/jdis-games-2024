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
ARG API_URL
ARG TYPE

WORKDIR /app

COPY interface/package.json .
COPY interface/package-lock.json .
COPY interface/tsconfig.json .
COPY interface/index.html .
COPY interface/style.css .
COPY interface/src/ src/

COPY --from=wasm /app/wasm_exec.js src/wasm/wasm_exec.js
RUN sed -i "s|const URL_BASE = '';|const URL_BASE = '$API_URL';|g" src/config.ts
RUN sed -i "s|let TYPE = '';|let TYPE = '$TYPE';|g" src/config.ts

RUN npm install
RUN npm run build -- --public-url "https://$API_URL"

# ------------------------------
FROM golang:1.21.6 as server

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY --from=build /app/dist /app/dist
COPY --from=wasm /app/lib.wasm /app/dist/lib.wasm
COPY interface/assets /app/dist/assets

RUN go build -o server /app/main.go

ENTRYPOINT ["./server"]