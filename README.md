## Intro
เริ่มต้น concept ของเราคือการพยายามเอาความลับทั้งหมดไปเก็บอยู่ที่ vault server เเล้วส่งมาให้ตัว app ของเราเรียกใช้งานในตัวอย่างนี้เราจะเอา config env มาจาก vailt server เเล้วเอามายัดใส่ struct config ของเราเพื่อเรียกใช้งานใน application ของเรา

## Work flow
เรามาดูการทำงานของโปรเเกรมที่เราจะทำกันบ้างนะครับ <br>

<i> app -> vault agent -> vault server # สำหรับ env อย่างเดียว</i> 
<br>
<i> app -> vault server # สำหรับเรียกขอข้อมูล secret อย่างอื่น</i>

** เหตุผลที่เราต้องการใช้งาน vault agent คือเราต้องการใช้งานความสามารถในการ renew token, cache เเละการทำ auth กับ vault server ส่วนอื่นๆผมจะให้ app คุยกับ vault server โดยตรงเลยเพื่อความง่ายในการใช้งานเเละการ implement



## How to use
### Run command
```
docker compose up vault -d

docker exec -it vault sh
```

### In vault container
```
# setup VAULT_ADDR for connect to server default is https:localhost:8200
export VAULT_ADDR="http://localhost:8200"

# login to vault server
vault login <your_root_token>

# go to vault server ui and create secret, policies and approle

# bind policies to approle path
vault write auth/approle/role/<your_role> \
  token_policies="<your_policies>" \
  token_ttl=1h \
  token_max_ttl=4h

# create role_id and secret_id for generate token access
vault read auth/approle/role/<your_role>/role-id

# -f if your don't setting config
vault write -f auth/approle/role/<your_role>/secret-id 

# **paste the role-id and secret-id to 
# - /vault_configs/role_id
# - /vault_configs/secret_id
```