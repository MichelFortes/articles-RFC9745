
# Build the plugin

[ More Info](https://www.krakend.io/docs/extending/writing-plugins/#compile-the-plugin)

## AMD64
```bash
docker run -it -v "$PWD:/app" -w /app krakend/builder:2.9.3 go build -buildmode=plugin -o deprecated-header.so .
```

## ARM64
```bash
docker run -it -v "$PWD:/app" -w /app \
-e "CGO_ENABLED=1" \
-e "CC=aarch64-linux-musl-gcc" \
-e "GOARCH=arm64" \
-e "GOHOSTARCH=amd64" \
krakend/builder:2.9.3 \
go build -ldflags='-extldflags=-fuse-ld=bfd -extld=aarch64-linux-musl-gcc' \
-buildmode=plugin -o deprecated-header.so .
```
