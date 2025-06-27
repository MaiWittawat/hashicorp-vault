## Intro
หลักจากที่ลองใช้งานตัว main(v1) จะเห็นว่าไฟล์ทั้งหมดถูกเขียนลงใน disk ซึ่งเป็นวิธีที่อาจจะยังไม่ปลอดภัยมากพอสำหรับบางระบบ v2 จึงเกิดขึ้นโดยเราจะย้ายการเก็บข้อมูล secret ทั้งหมดไปเก็บลงใน memmory เเทนส่งผลให้ความปลอดภัยของเรานั้นเพิ่มสูงขึ้นเเน่นอนทุกตรั้งที่มีการ down service หรือ restart มันจะลบข้อมูล secret ทั้งหมดออกไป 


## How to use

### Run Command
```
docker compose up vault -d

docker exec -it vault sh # เข้าไปใน vault ก่อนเพื่อตั่งค่าต่างๆ
```

### In vault container sh
```
export VAULT_ADDR="http://localhost:8200"

vault login <your_root_token>

# go to your vault ui at localhost:8200
# add you secret and setting the approle, policies

# Bind policies to approle path
vault write auth/approle/role/<your_role> \
  token_policies="<your_policies>" \
  token_ttl=1h \
  token_max_ttl=4h

# Create the role_id and secret id
# Save role_id and secret_id and plase it it ./vault_configs/role_id, ./vault_configs/secret_id
vault read auth/approle/role/<your-role>/role-id 
vault write -f auth/approle/role/admin/secret-id

# exit container
exit
```
### In root dir
```
docker compose up -d --build
```