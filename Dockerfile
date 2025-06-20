# Tahap pertama: Pembangunan
FROM golang:latest AS builder

# Set lingkungan kerja
WORKDIR /app

# Mengcopy go.mod dan go.sum agar dependensi dapat di-cache
COPY go.mod go.sum ./

# Mengunduh dependensi
RUN go mod download

# Mengcopy kode sumber ke dalam kontainer
COPY . .

# Membuat binary aplikasi
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/myapp main.go


# Tahap kedua: Produksi
FROM alpine:latest

# Instal tzdata untuk dukungan zona waktu
RUN apk --no-cache add tzdata

WORKDIR /app

# Mengcopy binary dari tahap pembangunan
COPY --from=builder /app/myapp .
COPY --from=builder /app/.env .

# Konfigurasi zona waktu
ENV TZ=Asia/Jakarta
RUN cp /usr/share/zoneinfo/Asia/Jakarta /etc/localtime && echo "Asia/Jakarta" > /etc/timezone

# Membuat direktori /assets/images dan mengatur izin
RUN mkdir -p /app/assets/images && chmod 777 /app/assets/images

# Menambahkan logging izin direktori
RUN ls -l /app/assets && ls -l /app/assets/images

# Menjalankan aplikasi
EXPOSE 3000

CMD ["./myapp"]