# 📦 API Testing Guide — Go E-commerce (ລາວ)

**Base URL:** `http://localhost:3000/api`  
**Content-Type:** `application/json`  
**Authorization:** `Bearer <access_token>` (route ທີ່ຕ້ອງການ Auth)

---

## 🔄 ຂັ້ນຕອນການ Test ແບບ Flow ສົມບູນ

```
1. Register / Login          →  ເກັບ access_token
2. (Admin) ສ້າງ Category     →  ເກັບ category_id
3. (Admin) ສ້າງ Seller       →  Login ດ້ວຍ seller account  →  ເກັບ seller_token
4. (Seller) ສ້າງ Product     →  ເກັບ product_id
5. (User) Add to Cart
6. (User) Create Order       →  ໃສ່ຂໍ້ມູນຈັດສົ່ງ  →  ເກັບ order_id
7. (User) Create Payment     →  ເກັບ payment_id
8. (Admin) Confirm Payment
9. (Admin) Update Order Status → confirmed
10.(Admin) Create Shipment   →  ເກັບ shipment_id
11.(Admin) Update Shipment   → shipped → delivered
```

---

## 🗂️ สรุป Route ทั้งหมด

| Method | Path | Auth | Role | ໜ້າທີ |
|--------|------|------|------|--------|
| POST | `/api/register` | ❌ | - | ສະໝັກສະມາຊິກ |
| POST | `/api/login` | ❌ | - | ເຂົ້າສູ່ລະບົບ |
| POST | `/api/refresh` | ❌ | - | ຕໍ່ Token |
| GET | `/api/products` | ❌ | - | ດູສິນຄ້າທັງໝົດ |
| GET | `/api/products/:id` | ❌ | - | ດູສິນຄ້າຕາມ ID |
| GET | `/api/products/seller/:id` | ❌ | - | ດູສິນຄ້າຂອງ Seller |
| GET | `/api/categories` | ❌ | - | ດູ Category ທັງໝົດ |
| GET | `/api/categories/:id` | ❌ | - | ດູ Category ຕາມ ID |
| GET | `/api/user/cart` | ✅ | user | ດູຕະກ້ານ |
| POST | `/api/user/cart/items` | ✅ | user | ເພີ່ມສິນຄ້າ |
| PUT | `/api/user/cart/items` | ✅ | user | ອັບເດດຈຳນວນ |
| DELETE | `/api/user/cart/items/:productId` | ✅ | user | ລຶບສິນຄ້າ |
| DELETE | `/api/user/cart` | ✅ | user | ລ້າງຕະກ້ານ |
| POST | `/api/user/orders` | ✅ | user | ສ້າງ Order |
| GET | `/api/user/orders` | ✅ | user | ດູ Order ຂອງຕົນ |
| GET | `/api/user/orders/:id` | ✅ | user | ດູ Order ຕາມ ID |
| POST | `/api/user/payments` | ✅ | user | ສ້າງ Payment |
| GET | `/api/user/payments/order/:orderId` | ✅ | user | ດູ Payment ຂອງ Order |
| GET | `/api/user/shipments/order/:orderId` | ✅ | user | ຕິດຕາມການຈັດສົ່ງ |
| POST | `/api/seller/products` | ✅ | seller | ສ້າງສິນຄ້າ |
| PUT | `/api/seller/products/:id` | ✅ | seller | ແກ້ໄຂສິນຄ້າ |
| DELETE | `/api/seller/products/:id` | ✅ | seller | ລຶບສິນຄ້າ |
| POST | `/api/admin/users` | ✅ | admin | ສ້າງ User |
| GET | `/api/admin/users` | ✅ | admin | ດູ User ທັງໝົດ |
| GET | `/api/admin/users/:id` | ✅ | admin | ດູ User ຕາມ ID |
| DELETE | `/api/admin/users/:id` | ✅ | admin | ລຶບ User |
| POST | `/api/admin/sellers` | ✅ | admin | ສ້າງ Seller |
| GET | `/api/admin/sellers` | ✅ | admin | ດູ Seller ທັງໝົດ |
| GET | `/api/admin/sellers/:id` | ✅ | admin | ດູ Seller ຕາມ ID |
| PUT | `/api/admin/sellers/:id` | ✅ | admin | ອັບເດດ Seller |
| DELETE | `/api/admin/sellers/:id` | ✅ | admin | ລຶບ Seller |
| POST | `/api/admin/categories` | ✅ | admin | ສ້າງ Category |
| PUT | `/api/admin/categories/:id` | ✅ | admin | ອັບເດດ Category |
| DELETE | `/api/admin/categories/:id` | ✅ | admin | ລຶບ Category |
| GET | `/api/admin/orders` | ✅ | admin | ດູ Order ທັງໝົດ |
| PATCH | `/api/admin/orders/:id/status` | ✅ | admin | ອັບເດດ Order Status |
| PATCH | `/api/admin/payments/:id/confirm` | ✅ | admin | ຢືນຢັນ Payment |
| POST | `/api/admin/shipments` | ✅ | admin | ສ້າງ Shipment |
| PATCH | `/api/admin/shipments/:id/status` | ✅ | admin | ອັບເດດ Shipment Status |
| PATCH | `/api/admin/shipments/:id/tracking` | ✅ | admin | ອັບເດດ Tracking |

---

## 1. 🔑 Auth (Public)

### POST /api/register — ສະໝັກສະມາຊິກໃໝ່

```http
POST http://localhost:3000/api/register
Content-Type: application/json

{
  "email": "somchai@example.com",
  "password": "secret123"
}
```

**Response 200:**
```json
{ "message": "User registered successfully" }
```

---

### POST /api/login — ເຂົ້າສູ່ລະບົບ

> ⭐ **ເກັບ `access_token` ເອົາໄວ້ໃຊ້ທຸກ request ທີ່ຕ້ອງການ Auth**

```http
POST http://localhost:3000/api/login
Content-Type: application/json

{
  "email": "somchai@example.com",
  "password": "secret123"
}
```

**Response 200:**
```json
{
  "access_token": "eyJhbGci...",
  "refresh_token": "eyJhbGci..."
}
```

---

### POST /api/refresh — ຂໍ Token ໃໝ່

```http
POST http://localhost:3000/api/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGci..."
}
```

**Response 200:**
```json
{
  "access_token": "eyJhbGci...(ໃໝ່)",
  "refresh_token": "eyJhbGci...(ໃໝ່)"
}
```

---

## 2. 🛍️ Products (Public — ບໍ່ຕ້ອງ Login)

### GET /api/products — ດູສິນຄ້າທັງໝົດ

```http
GET http://localhost:3000/api/products
```

**Response 200:**
```json
[
  {
    "id": "uuid",
    "seller_id": "uuid",
    "name": "ໂທລະສັບ Samsung A55",
    "description": "ສະມາດໂຟນ 6.4 ນິ້ວ RAM 8GB",
    "price": 3500000,
    "stock": 20,
    "status": "active",
    "images": [],
    "categories": [],
    "created_at": "2026-04-20T..."
  }
]
```

---

### GET /api/products/:id — ດູສິນຄ້າຕາມ ID

```http
GET http://localhost:3000/api/products/550e8400-e29b-41d4-a716-446655440000
```

---

### GET /api/products/seller/:id — ດູສິນຄ້າຂອງ Seller

```http
GET http://localhost:3000/api/products/seller/550e8400-e29b-41d4-a716-446655440000
```

---

## 3. 🏷️ Categories (Public)

### GET /api/categories — ດູ Category ທັງໝົດ

```http
GET http://localhost:3000/api/categories
```

**Response 200:**
```json
[
  {
    "id": "uuid",
    "name": "ອີເລັກໂຕຣນິກ",
    "parent_id": null,
    "created_at": "2026-04-20T..."
  }
]
```

---

### GET /api/categories/:id — ດູ Category ຕາມ ID

```http
GET http://localhost:3000/api/categories/550e8400-e29b-41d4-a716-446655440000
```

---

## 4. 🛒 Cart (🔒 ຕ້ອງ Login)

> Header ທຸກ Request: `Authorization: Bearer <access_token>`

### GET /api/user/cart — ດູຕະກ້ານ

> ສ້າງໃໝ່ອັດຕະໂນມັດຖ້າຍັງບໍ່ມີ

```http
GET http://localhost:3000/api/user/cart
Authorization: Bearer eyJhbGci...
```

**Response 200:**
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "items": [
    {
      "cart_id": "uuid",
      "product_id": "uuid",
      "quantity": 2
    }
  ],
  "created_at": "2026-04-20T..."
}
```

---

### POST /api/user/cart/items — ເພີ່ມສິນຄ້າເຂົ້າຕະກ້ານ

> ຖ້າສິນຄ້ານັ້ນມີຢູ່ແລ້ວໃນຕະກ້ານ ຈະບວກ quantity ໃຫ້ອັດຕະໂນມັດ

```http
POST http://localhost:3000/api/user/cart/items
Authorization: Bearer eyJhbGci...
Content-Type: application/json

{
  "product_id": "550e8400-e29b-41d4-a716-446655440000",
  "quantity": 2
}
```

**Response 200:**
```json
{ "message": "Item added to cart" }
```

---

### PUT /api/user/cart/items — ອັບເດດຈຳນວນ

> ໃສ່ `quantity: 0` ເພື່ອລຶບສິນຄ້ານັ້ນອອກ

```http
PUT http://localhost:3000/api/user/cart/items
Authorization: Bearer eyJhbGci...
Content-Type: application/json

{
  "product_id": "550e8400-e29b-41d4-a716-446655440000",
  "quantity": 5
}
```

---

### DELETE /api/user/cart/items/:productId — ລຶບສິນຄ້າອອກຈາກຕະກ້ານ

```http
DELETE http://localhost:3000/api/user/cart/items/550e8400-e29b-41d4-a716-446655440000
Authorization: Bearer eyJhbGci...
```

---

### DELETE /api/user/cart — ລ້າງຕະກ້ານທັງໝົດ

```http
DELETE http://localhost:3000/api/user/cart
Authorization: Bearer eyJhbGci...
```

---

## 5. 📋 Orders (🔒 ຕ້ອງ Login)

> ⚠️ **ຕ້ອງມີສິນຄ້າໃນຕະກ້ານກ່ອນ ຈຶ່ງ Order ໄດ້**

### POST /api/user/orders — ສ້າງ Order ຈາກ Cart

```http
POST http://localhost:3000/api/user/orders
Authorization: Bearer eyJhbGci...
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
|-------|----------|---------|
| `receiver_name` | ✅ | ຊື່ຜູ້ຮັບ |
| `phone` | ✅ | ເບີໂທ |
| `province` | ✅ | ແຂວງ |
| `district` | ✅ | ເມືອງ |
| `logistic` | ❌ | ບໍລິສັດຂົນສົ່ງ (`ans`, `mpsl`, `ems`) |
| `logistic_branch` | ❌ | ສາຂາຂົນສົ່ງ |

**Response 201:**
```json
{
  "order_id": "uuid",
  "message": "Order created successfully"
}
```

> ⭐ ຫລັງ Order ສຳເລັດ ຕະກ້ານຈະຖືກລ້າງອັດຕະໂນມັດ

---

### GET /api/user/orders — ດູ Order ຂອງຕົນເອງ

```http
GET http://localhost:3000/api/user/orders
Authorization: Bearer eyJhbGci...
```

**Response 200:**
```json
[
  {
    "id": "uuid",
    "user_id": "uuid",
    "total_price": 7000000,
    "status": "pending",
    "receiver_name":   "ສົມຊາຍ ໃຈດີ",
    "phone":           "20 5555 6666",
    "province":        "ໄຊຍະບູລີ",
    "district":        "ປາກລາຍ",
    "logistic":        "ans",
    "logistic_branch": "ສາຂາຫັດດາຍ",
    "order_items": [
      {
        "id": "uuid",
        "order_id": "uuid",
        "product_id": "uuid",
        "quantity": 2,
        "price": 3500000
      }
    ],
    "created_at": "2026-04-20T...",
    "updated_at": "2026-04-20T..."
  }
]
```

---

### GET /api/user/orders/:id — ດູ Order ຕາມ ID

```http
GET http://localhost:3000/api/user/orders/550e8400-e29b-41d4-a716-446655440000
Authorization: Bearer eyJhbGci...
```

---

## 6. 💳 Payments (🔒 ຕ້ອງ Login)

### POST /api/user/payments — ສ້າງ Payment

```http
POST http://localhost:3000/api/user/payments
Authorization: Bearer eyJhbGci...
Content-Type: application/json

{
  "order_id": "550e8400-e29b-41d4-a716-446655440000",
  "method": "bank_transfer",
  "amount": 7000000
}
```

| Field | ຄ່າທີ່ຮອງຮັບ |
|-------|------------|
| `method` | `bank_transfer` \| `bcel_one` \| `cash` |

**Response 201:**
```json
{ "payment_id": "uuid", "message": "Payment created" }
```

---

### GET /api/user/payments/order/:orderId — ດູ Payment ຂອງ Order

```http
GET http://localhost:3000/api/user/payments/order/550e8400-e29b-41d4-a716-446655440000
Authorization: Bearer eyJhbGci...
```

**Response 200:**
```json
{
  "id": "uuid",
  "order_id": "uuid",
  "method": "bank_transfer",
  "status": "pending",
  "amount": 7000000,
  "transaction_id": "",
  "paid_at": null,
  "created_at": "2026-04-20T..."
}
```

---

## 7. 🚚 Shipments — User (🔒 ຕ້ອງ Login)

### GET /api/user/shipments/order/:orderId — ຕິດຕາມການຈັດສົ່ງ

```http
GET http://localhost:3000/api/user/shipments/order/550e8400-e29b-41d4-a716-446655440000
Authorization: Bearer eyJhbGci...
```

**Response 200:**
```json
{
  "id": "uuid",
  "order_id": "uuid",
  "provider": "ans",
  "tracking_number": "ANS123456789",
  "status": "shipped",
  "shipped_at": "2026-04-20T...",
  "delivered_at": null,
  "created_at": "2026-04-20T..."
}
```

---

## 8. 🏪 Seller Routes (🔒 Auth + Role: seller)

> Header ທຸກ Request: `Authorization: Bearer <seller_token>`

### POST /api/seller/products — ສ້າງສິນຄ້າໃໝ່

```http
POST http://localhost:3000/api/seller/products
Authorization: Bearer <seller_token>
Content-Type: application/json

{
  "name": "ໂທລະສັບ Samsung A55",
  "description": "ສະມາດໂຟນ 6.4 ນິ້ວ RAM 8GB",
  "price": 3500000,
  "stock": 20,
  "status": "active",
  "image_urls": [
    "https://example.com/images/samsung-a55-1.jpg"
  ],
  "category_ids": [
    "550e8400-e29b-41d4-a716-446655440000"
  ]
}
```

| Field | Required | ຄວາມໝາຍ |
|-------|----------|---------|
| `name` | ✅ | ຊື່ສິນຄ້າ |
| `description` | ❌ | ລາຍລະອຽດ |
| `price` | ✅ | ລາຄາ (ກີບ) |
| `stock` | ✅ | ຈຳນວນສິນຄ້າ |
| `status` | ✅ | `active` \| `inactive` |
| `image_urls` | ❌ | URL ຮູບສິນຄ້າ |
| `category_ids` | ❌ | UUID ຂອງ Category |

**Response 201:**
```json
{ "product_id": "uuid", "message": "Product created successfully" }
```

---

### PUT /api/seller/products/:id — ແກ້ໄຂສິນຄ້າ

```http
PUT http://localhost:3000/api/seller/products/550e8400-e29b-41d4-a716-446655440000
Authorization: Bearer <seller_token>
Content-Type: application/json

{
  "name": "ໂທລະສັບ Samsung A55 (Updated)",
  "description": "ອັບເດດລາຍລະອຽດ",
  "price": 3200000,
  "stock": 15,
  "status": "active"
}
```

---

### DELETE /api/seller/products/:id — ລຶບສິນຄ້າ

```http
DELETE http://localhost:3000/api/seller/products/550e8400-e29b-41d4-a716-446655440000
Authorization: Bearer <seller_token>
```

---

## 9. 👑 Admin Routes (🔒 Auth + Role: admin)

> Header ທຸກ Request: `Authorization: Bearer <admin_token>`

---

### 👤 User Management

#### POST /api/admin/users — ສ້າງ User + ກຳນົດ Role

```http
POST http://localhost:3000/api/admin/users
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "email": "seller@example.com",
  "password": "secret123",
  "roles": ["seller"]
}
```

**Roles ທີ່ໃຊ້ໄດ້:** `user` | `seller` | `admin`

---

#### GET /api/admin/users — ດູ User ທັງໝົດ

```http
GET http://localhost:3000/api/admin/users
Authorization: Bearer <admin_token>
```

---

#### GET /api/admin/users/:id — ດູ User ຕາມ ID

```http
GET http://localhost:3000/api/admin/users/550e8400-e29b-41d4-a716-446655440000
Authorization: Bearer <admin_token>
```

---

#### DELETE /api/admin/users/:id — ລຶບ User

```http
DELETE http://localhost:3000/api/admin/users/550e8400-e29b-41d4-a716-446655440000
Authorization: Bearer <admin_token>
```

---

### 🏪 Seller Management

#### POST /api/admin/sellers — ສ້າງ Seller (User + Seller profile ພ້ອມກັນ)

```http
POST http://localhost:3000/api/admin/sellers
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "email": "myshop@example.com",
  "password": "secret123",
  "roles": ["seller"],
  "store_name": "ຮ້ານ My Shop",
  "description": "ຮ້ານຂາຍສິນຄ້າເອເລັກໂຕຣນິກ"
}
```

**Response 201:**
```json
{ "seller_id": "uuid", "message": "Seller created successfully" }
```

---

#### GET /api/admin/sellers — ດູ Seller ທັງໝົດ

```http
GET http://localhost:3000/api/admin/sellers
Authorization: Bearer <admin_token>
```

---

#### GET /api/admin/sellers/:id — ດູ Seller ຕາມ ID

```http
GET http://localhost:3000/api/admin/sellers/550e8400-e29b-41d4-a716-446655440000
Authorization: Bearer <admin_token>
```

---

#### PUT /api/admin/sellers/:id — ອັບເດດ Seller

```http
PUT http://localhost:3000/api/admin/sellers/550e8400-e29b-41d4-a716-446655440000
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "store_name": "ຮ້ານ My Shop (ໃໝ່)",
  "description": "ອັບເດດ"
}
```

---

#### DELETE /api/admin/sellers/:id — ລຶບ Seller

```http
DELETE http://localhost:3000/api/admin/sellers/550e8400-e29b-41d4-a716-446655440000
Authorization: Bearer <admin_token>
```

---

### 🏷️ Category Management

#### POST /api/admin/categories — ສ້າງ Category ຫລັກ

```http
POST http://localhost:3000/api/admin/categories
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "name": "ອີເລັກໂຕຣນິກ",
  "parent_id": null
}
```

**Response 201:**
```json
{ "id": "uuid", "message": "Category created successfully" }
```

---

#### POST /api/admin/categories — ສ້າງ Sub-category

```http
POST http://localhost:3000/api/admin/categories
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "name": "ໂທລະສັບ",
  "parent_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

---

#### PUT /api/admin/categories/:id — ອັບເດດ Category

```http
PUT http://localhost:3000/api/admin/categories/550e8400-e29b-41d4-a716-446655440000
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "name": "ອີເລັກໂຕຣນິກ & ກາດເຈັດ",
  "parent_id": null
}
```

---

#### DELETE /api/admin/categories/:id — ລຶບ Category

```http
DELETE http://localhost:3000/api/admin/categories/550e8400-e29b-41d4-a716-446655440000
Authorization: Bearer <admin_token>
```

---

### 📋 Order Management

#### GET /api/admin/orders — ດູ Order ທັງໝົດ

```http
GET http://localhost:3000/api/admin/orders
Authorization: Bearer <admin_token>
```

---

#### PATCH /api/admin/orders/:id/status — ອັບເດດ Order Status

```http
PATCH http://localhost:3000/api/admin/orders/550e8400-e29b-41d4-a716-446655440000/status
Authorization: Bearer <admin_token>
Content-Type: application/json

{ "status": "confirmed" }
```

**Status Flow:**

```
pending → confirmed → shipped → delivered
                             ↘ cancelled
```

| Status | ຄວາມໝາຍ |
|--------|---------|
| `pending` | ລໍຖ້າ (ຄ່າເລີ່ມຕົ້ນ) |
| `confirmed` | ຢືນຢັນ Order ແລ້ວ |
| `shipped` | ສົ່ງສິນຄ້າແລ້ວ |
| `delivered` | ສົ່ງຮອດແລ້ວ |
| `cancelled` | ຍົກເລີກ |

---

### 💳 Payment Management

#### PATCH /api/admin/payments/:id/confirm — ຢືນຢັນ Payment

```http
PATCH http://localhost:3000/api/admin/payments/550e8400-e29b-41d4-a716-446655440000/confirm
Authorization: Bearer <admin_token>
Content-Type: application/json

{ "transaction_id": "TXN-20260420-001" }
```

---

### 🚚 Shipment Management

#### POST /api/admin/shipments — ສ້າງ Shipment

```http
POST http://localhost:3000/api/admin/shipments
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "order_id": "550e8400-e29b-41d4-a716-446655440000",
  "provider": "ans",
  "tracking_number": "ANS123456789"
}
```

**Provider ທີ່ໃຊ້ໄດ້:** `ans` | `mpsl` | `ems`

---

#### PATCH /api/admin/shipments/:id/status — ອັບເດດ Shipment Status

```http
PATCH http://localhost:3000/api/admin/shipments/550e8400-e29b-41d4-a716-446655440000/status
Authorization: Bearer <admin_token>
Content-Type: application/json

{ "status": "shipped" }
```

> ລະບົບ set `shipped_at` / `delivered_at` ອັດຕະໂນມັດ

---

#### PATCH /api/admin/shipments/:id/tracking — ອັບເດດ Tracking Number

```http
PATCH http://localhost:3000/api/admin/shipments/550e8400-e29b-41d4-a716-446655440000/tracking
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "provider": "mpsl",
  "tracking_number": "MPSL987654321"
}
```

---

## 📊 HTTP Status Codes

| Status | ຄວາມໝາຍ |
|--------|---------|
| `200` | ສຳເລັດ |
| `201` | ສ້າງຂໍ້ມູນສຳເລັດ |
| `400` | ຂໍ້ມູນບໍ່ຖືກຕ້ອງ / ຂໍ້ມູນຂາດ |
| `401` | ບໍ່ມີ Token / Token ໝົດອາຍຸ |
| `403` | ບໍ່ມີສິດ (Role ບໍ່ຕົງ) |
| `404` | ບໍ່ພົບຂໍ້ມູນ |
| `500` | Server Error |

---

## 🧪 curl — Full Flow Example

```bash
# ─── 1. Register ───────────────────────────────────────────────────
curl -X POST http://localhost:3000/api/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secret123"}'

# ─── 2. Login & Save Token ─────────────────────────────────────────
TOKEN=$(curl -s -X POST http://localhost:3000/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secret123"}' \
  | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)

echo "TOKEN: $TOKEN"

# ─── 3. ດູສິນຄ້າ ────────────────────────────────────────────────────
curl http://localhost:3000/api/products

# ─── 4. ເພີ່ມສິນຄ້າເຂົ້າຕະກ້ານ ──────────────────────────────────────
curl -X POST http://localhost:3000/api/user/cart/items \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"product_id":"<product_uuid>","quantity":2}'

# ─── 5. ສ້າງ Order ─────────────────────────────────────────────────
ORDER=$(curl -s -X POST http://localhost:3000/api/user/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "receiver_name": "ສົມຊາຍ ໃຈດີ",
    "phone": "20 5555 6666",
    "province": "ໄຊຍະບູລີ",
    "district": "ປາກລາຍ",
    "logistic": "ans",
    "logistic_branch": "ສາຂາຫັດດາຍ"
  }')
echo "$ORDER"
ORDER_ID=$(echo "$ORDER" | grep -o '"order_id":"[^"]*"' | cut -d'"' -f4)

# ─── 6. ສ້າງ Payment ───────────────────────────────────────────────
curl -X POST http://localhost:3000/api/user/payments \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{\"order_id\":\"$ORDER_ID\",\"method\":\"bank_transfer\",\"amount\":7000000}"
```

---

## 🛠️ ການ Test ດ້ວຍ Postman

### ຕັ້ງຄ່າ Environment Variables ໃນ Postman

| Variable | Value |
|----------|-------|
| `base_url` | `http://localhost:3000/api` |
| `access_token` | *(ໃສ່ຄ່າຫລັງ Login)* |
| `admin_token` | *(ໃສ່ຄ່າຫລັງ Admin Login)* |
| `seller_token` | *(ໃສ່ຄ່າຫລັງ Seller Login)* |
| `product_id` | *(ໃສ່ຄ່າຫລັງສ້າງ Product)* |
| `order_id` | *(ໃສ່ຄ່າຫລັງສ້າງ Order)* |
| `payment_id` | *(ໃສ່ຄ່າຫລັງສ້າງ Payment)* |
| `shipment_id` | *(ໃສ່ຄ່າຫລັງສ້າງ Shipment)* |

### Auto-save Token (Postman Pre-request Script)

ໃຊ້ script ນີ້ໃນ **Tests tab** ຂອງ Login request:

```javascript
const res = pm.response.json();
if (res.access_token) {
    pm.environment.set("access_token", res.access_token);
    pm.environment.set("refresh_token", res.refresh_token);
    console.log("✅ Token saved:", res.access_token.substring(0, 30) + "...");
}
```

### Authorization Header (ໃຊ້ Environment Variable)

ໃນ Header ຂອງ Request ທີ່ຕ້ອງ Auth:
```
Authorization: Bearer {{access_token}}
```

---

## ⚠️ Common Errors

| Error | ສາເຫດ | ວິທີແກ້ |
|-------|-------|--------|
| `401 Unauthorized` | Token ໝົດ / ບໍ່ມີ Token | Login ໃໝ່ ຫລື Refresh Token |
| `403 Forbidden` | Role ບໍ່ຖືກຕ້ອງ | ໃຊ້ Token ທີ່ຖືກ Role |
| `400 Bad Request` | ຂໍ້ມູນບໍ່ຄົບ / ຮູບແບບຜິດ | ກວດ Request Body |
| `404 Not Found` | ບໍ່ພົບ ID ທີ່ລະບຸ | ກວດ UUID ທີ່ໃຊ້ |
| `500 Internal Server Error` | Server/DB Error | ກວດ log ໃນ terminal |
