# --- Stage 1: Builder Stage ---
# ใช้ Go image ที่มี compiler สำหรับการ build
FROM golang:1.24-alpine AS builder

# ตั้งค่า Working Directory ภายในคอนเทนเนอร์
WORKDIR /app

# Copy go.mod และ go.sum เพื่อดาวน์โหลด dependencies
COPY go.mod go.sum ./
# ดาวน์โหลด Go modules (dependencies)
RUN go mod tidy

# Copy source code ของแอปพลิเคชันทั้งหมด
COPY . .

# Build แอปพลิเคชัน Go
# -o /app/hashicorp-vault: กำหนดชื่อและ path ของ executable file
# ./cmd/hashicorp-vault: path ไปยัง main package ของแอปพลิเคชันคุณ
RUN go build -o /app/hashicorp-vault ./cmd/hashicorp-vault

# --- Stage 2: Runner Stage ---
# ใช้ Alpine Linux เป็น base image สำหรับ runtime เพื่อให้ Image มีขนาดเล็ก
FROM alpine:latest

# ติดตั้งแพ็กเกจที่จำเป็น (ถ้ามี)
# เช่น, ca-certificates หากแอปพลิเคชันของคุณมีการเรียก HTTPS ไปยังภายนอก
RUN apk add --no-cache ca-certificates

# ตั้งค่า Working Directory สำหรับแอปพลิเคชัน
WORKDIR /app

RUN apk add --no-cache ca-certificates curl

# Copy executable file ที่ build ได้จาก builder stage
COPY --from=builder /app/hashicorp-vault .

# --- (Optional) Entrypoint Script สำหรับรอ Vault Agent ---
# นี่เป็นสิ่งสำคัญมาก หากคุณต้องการให้แอปพลิเคชันรอจนกว่า Vault Agent จะพร้อม
# Copy entrypoint script เข้ามาในคอนเทนเนอร์
COPY ./entrypoint.sh .
# ทำให้ script สามารถ execute ได้
RUN chmod +x ./entrypoint.sh

# กำหนด ENTRYPOINT และ CMD
# ENTRYPOINT จะรัน entrypoint.sh ก่อน
ENTRYPOINT ["./entrypoint.sh"]
# CMD จะถูกส่งเป็น arguments ไปยัง ENTRYPOINT script
CMD ["./hashicorp-vault"]

# หรือถ้าคุณไม่ใช้ entrypoint script และแอปพลิเคชันของคุณมี retry logic ในตัวอยู่แล้ว:
# CMD ["./hashicorp-vault"]