#!/bin/sh

# กำหนดที่อยู่ของ Vault Agent HTTP Listener ภายในคอนเทนเนอร์เดียวกัน
# เนื่องจาก network_mode: service:vault-agent ทำให้ app และ agent แชร์ network namespace เดียวกัน
# Vault Agent จะฟังอยู่ที่ 127.0.0.1:8200 ตามที่เราตั้งค่าใน agent-config.hcl
VAULT_AGENT_ADDR="127.0.0.1:8200"

# กำหนด Path ของความลับใน Vault Server ที่คุณต้องการดึงผ่าน Vault Agent
# นี่คือ Path ที่คุณจะใช้ในโค้ด Go ของคุณด้วย เช่น http://127.0.0.1:8200/v1/secret/data/my-app/config
# ให้ตรงกับ secret engine และ path ที่คุณตั้งค่าใน Vault
SECRET_API_PATH="v1/secret/data/dotenv"

# กำหนดจำนวนครั้งสูงสุดที่จะลองเชื่อมต่อและช่วงเวลาการลองใหม่
MAX_RETRIES=60
RETRY_INTERVAL=1 # วินาที

echo "Waiting for Vault Agent HTTP Listener at ${VAULT_AGENT_ADDR}..."

# วนลูปเพื่อลองเชื่อมต่อไปยัง Vault Agent จนกว่าจะพร้อม
for i in $(seq 1 $MAX_RETRIES); do
  # ใช้ curl เพื่อลองเรียก health check หรือ path ความลับใดๆ จาก Vault Agent
  # -s: Silent mode (ไม่แสดง progress)
  # -o /dev/null: ทิ้ง output ของ response body
  # -w "%{http_code}": แสดงเฉพาะ HTTP status code
  # grep "200": ตรวจสอบว่า status code เป็น 200 (OK)
  if curl -s -o /dev/null -w "%{http_code}" http://${VAULT_AGENT_ADDR}/${SECRET_API_PATH} | grep "200" > /dev/null; then
    echo "Vault Agent is ready. Starting application."
    # ใช้ exec "$@" เพื่อรันคำสั่งหลักของคอนเทนเนอร์ (ซึ่งกำหนดโดย CMD ใน Dockerfile)
    # exec จะแทนที่ process ของ shell script ด้วย process ของแอปพลิเคชัน Go
    exec "$@"
  fi
  echo "Vault Agent not ready yet. Retrying in ${RETRY_INTERVAL} second(s)... ($i/$MAX_RETRIES)"
  sleep $RETRY_INTERVAL
done

# หากลองเชื่อมต่อครบจำนวนครั้งแล้วแต่ Vault Agent ยังไม่พร้อม ก็ให้จบการทำงานด้วย error
echo "Timeout waiting for Vault Agent. Exiting."
exit 1