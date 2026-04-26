# API Testing Guide — Go E-commerce (ລາວ)

**Base URL:** `http://localhost:3000/api`  
**Port:** `3000` (ตั้งใน .env)

---

## 🔑 Auth Headers

| Route | Header |
|-------|--------|
| Public | ไม่ต้อง |
| User | `Authorization: Bearer <access_token>` |
| Seller | `Authorization: Bearer <seller_token>` (role: seller) |
| Admin | `Authorization: Bearer <admin_token>` (role: admin) |

---

## 🔄 Full Test Flow

```
 1. Register → Login → เก็บ token
 2. (Admin) สร้าง Category
 3. (Admin) สร้าง Seller account
 4. (Seller) Login → สร้าง Product
 5. (User) เพิ่มสินค้าเข้า Cart
 6. (User) Create Order → ใส่ข้อมูลจัดส่ง
 7. (User) Create Payment → ได้ Phajay payment_url
 8. ลูกค้าเปิด payment_url → scan QR จ่ายเงิน
 9. Phajay webhook callback → order confirmed อัตโนมัติ
10. (Admin) สร้าง Shipment → อัปเดต tracking
```

---

## 1. 🔑 Auth (Public)

### POST /api/register
```http
POST http://localhost:3000/api/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "secret123"
}
```
✅ `200` → `{ "message": "User registered successfully" }`

---

### POST /api/login
> ⭐ เก็บ `access_token` ไว้ใช้ทุก request

```http
POST http://localhost:3000/api/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "secret123"
}
```
✅ `200`:
```json
{
  "access_token": "eyJhbGci...",
  "refresh_token": "eyJhbGci..."
}
```

---

### POST /api/refresh
```http
POST http://localhost:3000/api/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGci..."
}
```
✅ `200` → token ชุดใหม่

---

## 2. 🏷️ Categories (Public)

### GET /api/categories
```http
GET http://localhost:3000/api/categories
```
✅ `200`:
```json
[{ "id": "uuid", "name": "ອີເລັກໂຕຣນິກ", "created_at": "..." }]
```

---

### GET /api/categories/:id
```http
GET http://localhost:3000/api/categories/<category_id>
```

---

## 3. 🛍️ Products (Public)

### GET /api/products
```http
GET http://localhost:3000/api/products
```
✅ `200`:
```json
[
  {
    "id": "uuid",
    "seller_id": "uuid",
    "name": "Samsung A55",
    "description": "ສະມາດໂຟນ",
    "price": 3500000,
    "stock": 20,
    "status": "active",
    "images": [],
    "created_at": "..."
  }
]
```

---

### GET /api/products/:id
```http
GET http://localhost:3000/api/products/<product_id>
```

---

### GET /api/products/category/:id — ສິນຄ້າຕາມ Category
```http
GET http://localhost:3000/api/products/category/<category_id>
```
✅ `200` → list ສິນຄ້າ (ถ้าไม่มี return `[]`)

---

### GET /api/products/seller/:id — ສິນຄ້າຂອງ Seller
```http
GET http://localhost:3000/api/products/seller/<seller_id>
```

---

## 4. 🛒 Cart (🔒 User Auth)

### GET /api/user/cart
> Cart ສ້າງອັດຕະໂນມັດຖ້າຍັງບໍ່ມີ

```http
GET http://localhost:3000/api/user/cart
Authorization: Bearer <access_token>
```
✅ `200`:
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "items": [
    { "cart_id": "uuid", "product_id": "uuid", "quantity": 2 }
  ],
  "created_at": "..."
}
```

---

### POST /api/user/cart/items — ເພີ່ມສິນຄ້າ
```http
POST http://localhost:3000/api/user/cart/items
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "product_id": "<product_id>",
  "quantity": 2
}
```
✅ `200` → `{ "message": "Item added to cart" }`  
❌ `400` → `product not found` | `quantity must be greater than zero`

---

### PUT /api/user/cart/items — ອັບເດດ quantity
```http
PUT http://localhost:3000/api/user/cart/items
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "product_id": "<product_id>",
  "quantity": 5
}
```

---

### DELETE /api/user/cart/items/:productId
```http
DELETE http://localhost:3000/api/user/cart/items/<product_id>
Authorization: Bearer <access_token>
```

---

### DELETE /api/user/cart — ລ້າງທັງໝົດ
```http
DELETE http://localhost:3000/api/user/cart
Authorization: Bearer <access_token>
```

---

## 5. 📋 Orders (🔒 User Auth)

> ⚠️ ຕ້ອງມີສິນຄ້າໃນຕະກ້ານກ່ອນ

### POST /api/user/orders — ສັ່ງຊື້
```http
POST http://localhost:3000/api/user/orders
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "receiver_name":   "ສົມຊາຍ ໃຈດີ",
  "phone":           "20 5555 6666",
  "province":        "ໄຊຍະບູລີ",
  "district":        "ປາກລາຍ",
  "logistic":        "ans",
  "logistic_branch": "ສາຂາຫັດດາຍ"
}
```

| Field | Required | ຄວາມໝາຍ |
|-------|:--------:|---------|
| `receiver_name` | ✅ | ຊື່ຜູ້ຮັບ |
| `phone` | ✅ | ເບີໂທ |
| `province` | ✅ | ແຂວງ |
| `district` | ✅ | ເມືອງ |
| `logistic` | ❌ | ບໍລິສັດຂົນສົ່ງ (`ans`, `mpsl`, `ems`) |
| `logistic_branch` | ❌ | ສາຂາ |

✅ `201`:
```json
{ "order_id": "uuid", "message": "Order created successfully" }
```
❌ `400` → `cart is empty` | `insufficient stock for product: xxx`

> ⭐ ຕະກ້ານຈະຖືກລ້າງອັດຕະໂນມັດຫລັງ Order ສຳເລັດ

---

### GET /api/user/orders
```http
GET http://localhost:3000/api/user/orders
Authorization: Bearer <access_token>
```
✅ `200`:
```json
[
  {
    "id": "uuid",
    "user_id": "uuid",
    "total_price": 7000000,
    "status": "pending",
    "receiver_name": "ສົມຊາຍ ໃຈດີ",
    "phone": "20 5555 6666",
    "province": "ໄຊຍະບູລີ",
    "district": "ປາກລາຍ",
    "logistic": "ans",
    "logistic_branch": "ສາຂາຫັດດາຍ",
    "order_items": [
      { "id": "uuid", "product_id": "uuid", "quantity": 2, "price": 3500000 }
    ],
    "created_at": "...",
    "updated_at": "..."
  }
]
```

---

### GET /api/user/orders/:id
```http
GET http://localhost:3000/api/user/orders/<order_id>
Authorization: Bearer <access_token>
```

---

## 6. 💳 Payments — Phajay (🔒 User Auth)

### POST /api/user/payments — ສ້າງ Payment (ໄດ້ Phajay Link ກັບ)

> ⭐ ລະບົບຈະເອີ້ນ Phajay API ໃຫ້ອັດຕະໂນມັດ → return `payment_url` ໃຫ້ redirect ລູກຄ້າ

```http
POST http://localhost:3000/api/user/payments
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "order_id": "<order_id>",
  "amount": 7000000
}
```

| Field | Required | ຄວາມໝາຍ |
|-------|:--------:|---------|
| `order_id` | ✅ | ID ຂອງ order ທີ່ຕ້ອງການຈ່າຍ |
| `amount` | ✅ | ຈຳນວນເງິນ (ຕ້ອງ > 0) |

✅ `201`:
```json
{
  "payment_id": "uuid",
  "payment_url": "https://payment-gateway.phajay.co/pay/abc123",
  "message": "Payment created — redirect customer to payment_url"
}
```

> **ຂັ້ນຕອນຕໍ່ໄປ:** Frontend ເອົາ `payment_url` ໄປເປີດໃຫ້ລູກຄ້າ → ລູກຄ້າ scan QR ຈ່າຍເງິນຜ່ານ BCEL/JDB/LDB

❌ `400` → `order not found` | `order_id and amount are required`

---

### GET /api/user/payments/order/:orderId — ດູ Payment ຂອງ Order
```http
GET http://localhost:3000/api/user/payments/order/<order_id>
Authorization: Bearer <access_token>
```
✅ `200`:
```json
{
  "id": "uuid",
  "order_id": "uuid",
  "method": "phajay",
  "status": "pending",
  "amount": 7000000,
  "transaction_id": "",
  "payment_url": "https://payment-gateway.phajay.co/pay/abc123",
  "paid_at": null,
  "created_at": "..."
}
```

**Payment Status Flow:**
```
pending → paid (ຈ່າຍສຳເລັດ)
pending → failed (ຈ່າຍບໍ່ສຳເລັດ/ຍົກເລີກ)
```

---

## 7. 🔔 Phajay Webhook (Public — Phajay ເອີ້ນ endpoint ນີ້)

### POST /api/webhooks/phajay

> ⚠️ Route ນີ້ Frontend ບໍ່ຕ້ອງເອີ້ນ — Phajay server ຈະ callback ມາອັດຕະໂນມັດ

```http
POST http://localhost:3000/api/webhooks/phajay
Content-Type: application/json

{
  "orderNo": "<order_id>",
  "transactionId": "TXN-20260424-001",
  "amount": 7000000,
  "status": "success"
}
```

| Status | ຜົນ |
|--------|-----|
| `success` | payment → `paid`, order → `confirmed` (ອັດຕະໂນມັດ) |
| `failed` | payment → `failed` |
| `cancelled` | payment → `failed` |

✅ `200` → `{ "message": "webhook processed" }`

> ⭐ ຕ້ອງ config Webhook URL ໃນ [portal.phajay.co](https://portal.phajay.co) ໃຫ້ point ມາທີ່:  
> `https://your-domain.com/api/webhooks/phajay`

---

## 8. 🚚 Shipments (🔒 User — ດູໄດ້)

### GET /api/user/shipments/order/:orderId
```http
GET http://localhost:3000/api/user/shipments/order/<order_id>
Authorization: Bearer <access_token>
```
✅ `200`:
```json
{
  "id": "uuid",
  "order_id": "uuid",
  "provider": "ans",
  "tracking_number": "ANS123456789",
  "status": "shipped",
  "shipped_at": "...",
  "delivered_at": null,
  "created_at": "..."
}
```

---

## 9. 🏪 Seller (🔒 Auth + Role: seller)

### POST /api/seller/products/create
```http
POST http://localhost:3000/api/seller/products/create
Authorization: Bearer <seller_token>
Content-Type: application/json

{
  "name": "Samsung A55",
  "description": "ສະມາດໂຟນ 6.4 ນິ້ວ",
  "price": 3500000,
  "stock": 20,
  "status": "active",
  "image_urls": ["https://example.com/img/samsung.jpg"],
  "category_ids": ["<category_id>"]
}
```
✅ `201` → `{ "product_id": "uuid", "message": "Product created successfully" }`

---

### PUT /api/seller/products/update/:id
```http
PUT http://localhost:3000/api/seller/products/update/<product_id>
Authorization: Bearer <seller_token>
Content-Type: application/json

{
  "name": "Samsung A55 (Updated)",
  "description": "ອັບເດດ",
  "price": 3200000,
  "stock": 15,
  "status": "active"
}
```

---

### DELETE /api/seller/products/delete/:id
```http
DELETE http://localhost:3000/api/seller/products/delete/<product_id>
Authorization: Bearer <seller_token>
```

---

## 10. 👑 Admin (🔒 Auth + Role: admin)

### 👤 Users

#### POST /api/admin/users
```http
POST http://localhost:3000/api/admin/users
Authorization: Bearer <admin_token>
Content-Type: application/json

{ "email": "newseller@example.com", "password": "secret123", "roles": ["seller"] }
```

#### GET /api/admin/users
```http
GET http://localhost:3000/api/admin/users
Authorization: Bearer <admin_token>
```

#### GET /api/admin/users/:id
```http
GET http://localhost:3000/api/admin/users/<user_id>
Authorization: Bearer <admin_token>
```

#### DELETE /api/admin/users/:id
```http
DELETE http://localhost:3000/api/admin/users/<user_id>
Authorization: Bearer <admin_token>
```

---

### 🏪 Sellers

#### POST /api/admin/sellers
```http
POST http://localhost:3000/api/admin/sellers
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "email": "shop@example.com",
  "password": "secret123",
  "roles": ["seller"],
  "store_name": "ຮ້ານ ABC",
  "description": "ຮ້ານຂາຍເຄື່ອງໄຟຟ້າ"
}
```

#### GET /api/admin/sellers
```http
GET http://localhost:3000/api/admin/sellers
Authorization: Bearer <admin_token>
```

#### GET /api/admin/sellers/:id
```http
GET http://localhost:3000/api/admin/sellers/<seller_id>
Authorization: Bearer <admin_token>
```

#### PUT /api/admin/sellers/:id
```http
PUT http://localhost:3000/api/admin/sellers/<seller_id>
Authorization: Bearer <admin_token>
Content-Type: application/json

{ "store_name": "ຮ້ານ ABC (ໃໝ່)", "description": "ອັບເດດ" }
```

#### DELETE /api/admin/sellers/:id
```http
DELETE http://localhost:3000/api/admin/sellers/<seller_id>
Authorization: Bearer <admin_token>
```

---

### 🏷️ Categories

#### POST /api/admin/categories/create
```http
POST http://localhost:3000/api/admin/categories/create
Authorization: Bearer <admin_token>
Content-Type: application/json

{ "name": "ອີເລັກໂຕຣນິກ" }
```
✅ `201` → `{ "id": "uuid", "message": "Category created successfully" }`

#### PUT /api/admin/categories/update/:id
```http
PUT http://localhost:3000/api/admin/categories/update/<category_id>
Authorization: Bearer <admin_token>
Content-Type: application/json

{ "name": "ອີເລັກໂຕຣນິກ (ໃໝ່)" }
```

#### DELETE /api/admin/categories/delete/:id
```http
DELETE http://localhost:3000/api/admin/categories/delete/<category_id>
Authorization: Bearer <admin_token>
```

---

### 📋 Orders

#### GET /api/admin/orders
```http
GET http://localhost:3000/api/admin/orders
Authorization: Bearer <admin_token>
```

#### PATCH /api/admin/orders/:id/status
```http
PATCH http://localhost:3000/api/admin/orders/<order_id>/status
Authorization: Bearer <admin_token>
Content-Type: application/json

{ "status": "confirmed" }
```
**Status flow:** `pending` → `confirmed` → `shipped` → `delivered` | `cancelled`

---

### 💳 Payments (Admin — ยืนยัน manual)

#### PATCH /api/admin/payments/:id/confirm
```http
PATCH http://localhost:3000/api/admin/payments/<payment_id>/confirm
Authorization: Bearer <admin_token>
Content-Type: application/json

{ "transaction_id": "TXN-MANUAL-001" }
```

---

### 🚚 Shipments

#### POST /api/admin/shipments/create
```http
POST http://localhost:3000/api/admin/shipments/create
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "order_id": "<order_id>",
  "provider": "ans",
  "tracking_number": "ANS123456789"
}
```

#### PATCH /api/admin/shipments/:id/status
```http
PATCH http://localhost:3000/api/admin/shipments/<shipment_id>/status
Authorization: Bearer <admin_token>
Content-Type: application/json

{ "status": "shipped" }
```
**Status flow:** `pending` → `shipped` → `delivered`

#### PATCH /api/admin/shipments/:id/tracking
```http
PATCH http://localhost:3000/api/admin/shipments/<shipment_id>/tracking
Authorization: Bearer <admin_token>
Content-Type: application/json

{ "provider": "mpsl", "tracking_number": "MPSL987654321" }
```

---

## 📊 HTTP Status Codes

| Code | ຄວາມໝາຍ |
|------|---------|
| `200` | ສຳເລັດ |
| `201` | ສ້າງສຳເລັດ |
| `400` | ຂໍ້ມູນບໍ່ຖືກ / business logic error |
| `401` | ບໍ່ມີ token / token ໝົດອາຍຸ |
| `403` | Role ບໍ່ຕົງ |
| `404` | ບໍ່ພົບຂໍ້ມູນ |
| `500` | Server error |

---

## ⚙️ Phajay Setup (ຕັ້ງຄ່າ)

1. ສະໝັກບັນຊີທີ່ [portal.phajay.co](https://portal.phajay.co)
2. ເອົາ **Secret Key** ມາໃສ່ `.env`:
   ```
   PHAJAY_SECRET_KEY=pk_live_xxxxxxxxxxxxxxxx
   ```
3. ຕັ້ງ **Webhook URL** ໃນ portal:
   ```
   https://your-domain.com/api/webhooks/phajay
   ```
4. ຕັ້ງ **Success URL** ແລະ **Cancel URL** ໃນ portal (ສຳລັບ redirect ຫລັງຈ່າຍ)

---

## 🔄 Payment Flow Diagram

```
ລູກຄ້າ                    Server                     Phajay
  │                          │                          │
  │  POST /user/payments     │                          │
  │ {order_id, amount}       │                          │
  │─────────────────────────▶│                          │
  │                          │  POST /v1/api/link/      │
  │                          │  payment-link             │
  │                          │─────────────────────────▶│
  │                          │                          │
  │                          │  { paymentUrl: "..." }   │
  │                          │◀─────────────────────────│
  │                          │                          │
  │  { payment_url: "..." }  │                          │
  │◀─────────────────────────│                          │
  │                          │                          │
  │  ເປີດ payment_url        │                          │
  │  scan QR ຈ່າຍເງິນ       │                          │
  │─────────────────────────────────────────────────────▶│
  │                          │                          │
  │                          │  POST /webhooks/phajay   │
  │                          │  {status: "success"}     │
  │                          │◀─────────────────────────│
  │                          │                          │
  │                          │  payment→paid            │
  │                          │  order→confirmed         │
  │                          │                          │
```
