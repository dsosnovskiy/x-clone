# 🔥 X-CLONE - A simplified copy of the social network X (Twitter)

## 🛠 Technologies

- **Programming language**: `Golang`
- **Architecture**: `Layered`
- **Web framework**: `net/http` (HTTP server), `chi` (HTTP router)
- **Configuration**: `cleanenv`
- **Logging**: `logrus`
- **Validation**: `validator`
- **Database**: `PostgreSQL` (using GORM)
- **Authentication**: `JWT` (Access only)

## 🎯 Goals

- **Secure**: `Refresh Token`
- **Caching**: `Redis`
- **New functionality**:
  - News feed
  - Notifications
- **Containerization**: `Docker`

# 🔍 API Description

## ⚙️ Setup .env file

```
APP_ENV=local/dev/prod

DB_NAME=...
DB_USER=...
DB_PASSWORD=...
JWT_SECRET=...
```

# 🚀 Endpoints

## 🔐 Authentication

**Upon registration/authentication, the user receives a `JWT token` (Access only, `access_token_ttl: 2h`).**

**All other application routes are protected!**

**Header**:

- _Key_: Authorization
- _Value_: Bearer your_token

## **/auth/register {POST}**

**Description**: Register new user

**Request Body Schema**:

```json
{
  "username": "string",
  "password": "string",
  "first_name": "string",
  "last_name": "string",
  "birthday": "string",
  "bio": "string"
}
```

**Registration Fields**:

| Field        | Type   | Required | Limits               | Example              |
| ------------ | ------ | -------- | -------------------- | -------------------- |
| `username`   | string | Yes      | 6-20                 | `john_doe22`         |
| `password`   | string | Yes      | 7-32                 | `qwerty123`          |
| `first_name` | string | Yes      | 2-32                 | `John`               |
| `last_name`  | string | Yes      | 2-32                 | `Doe`                |
| `birthday`   | string | No       | Format: `YYYY-MM-DD` | `1990-05-15`         |
| `bio`        | string | No       | 1-300                | `Software Developer` |

**Response Body Schema**:

```json
{
  "access_token": "string"
}
```

---

## **/auth/login {POST}**

**Description**: Authenticate user

**Request Body Schema**:

```json
{
  "username": "string",
  "password": "string"
}
```

**Response Body Schema**:

```json
{
  "access_token": "string"
}
```

# 👤 User

## **/settings/profile {PATCH}**

**Description**: Change profile

**Request Body Schema**:

```json
{
  "username": "string",
  "first_name": "string",
  "last_name": "string",
  "birthday": "string",
  "bio": "string"
}
```

**Changed Fields**:

| Field        | Type   | Required | Limits               | Example              |
| ------------ | ------ | -------- | -------------------- | -------------------- |
| `username`   | string | No       | 6-20                 | `john_doe22`         |
| `first_name` | string | No       | 2-32                 | `John`               |
| `last_name`  | string | No       | 2-32                 | `Doe`                |
| `birthday`   | string | No       | Format: `YYYY-MM-DD` | `1990-05-15`         |
| `bio`        | string | No       | 1-300                | `Software Developer` |

**Response Body Schema**:

```json
{
  "user_id": "int",
  "username": "string",
  "first_name": "string",
  "last_name": "string",
  "birthday": "string",
  "bio": "string",
  "created_at": "string",
  "followers": "int",
  "following": "int"
}
```

## **/settings/password {PATCH}**

**Description**: Change password

**Response Body Schema**:

```json
{
  "old_password": "string",
  "new_password": "string",
  "confirm_password": "string"
}
```

**Changed Fields**:

| Field              | Type   | Required | Limits | Example     |
| ------------------ | ------ | -------- | ------ | ----------- |
| `old_password`     | string | Yes      | 7-32   | `qwerty123` |
| `new_password`     | string | Yes      | 7-32   | `zxcvbn456` |
| `confirm_password` | string | Yes      | 7-32   | `zxcvbn456` |

**Response Body Schema**:

```json
{
  "message": "successfully changed password"
}
```

## **/{username} {GET}**

**Description**: Get information about the user

**Response Body Schema**:

```json
{
  "user_id": "int",
  "username": "string",
  "first_name": "string",
  "last_name": "string",
  "birthday": "string",
  "bio": "string",
  "created_at": "string",
  "followers": "int",
  "following": "int"
}
```

## **/{username}/follow {POST}**

**Description**: Follow another user

**Response Body Schema**:

```json
{
  "message": "successfully followed the user",
  "username": "string"
}
```

## **/{username}/follow {DELETE}**

**Description**: Stop following another user

**Response Body Schema**:

```json
{
  "message": "successfully stop following the user",
  "username": "string"
}
```

## **/{username}/followers {GET}**

**Description**: Get the user's followers

**Response Body Schema**:

```json
[
  {
    "user_id": "int",
    "username": "string",
    "first_name": "string",
    "last_name": "string",
    "birthday": "string",
    "bio": "string",
    "created_at": "string",
    "followers": "int",
    "following": "int"
  },
  {
    "user_id": "int",
    "username": "string",
    "first_name": "string",
    "last_name": "string",
    "birthday": "string",
    "bio": "string",
    "created_at": "string",
    "followers": "int",
    "following": "int"
  }
]
```

## **/{username}/following {GET}**

**Description**: Get the user's following

**Response Body Schema**:

```json
[
  {
    "user_id": "int",
    "username": "string",
    "first_name": "string",
    "last_name": "string",
    "birthday": "string",
    "bio": "string",
    "created_at": "string",
    "followers": "int",
    "following": "int"
  },
  {
    "user_id": "int",
    "username": "string",
    "first_name": "string",
    "last_name": "string",
    "birthday": "string",
    "bio": "string",
    "created_at": "string",
    "followers": "int",
    "following": "int"
  }
]
```

# 📝 Post

## **/compose/post {POST}**

**Description**: Create post

**Request Body Schema**:

```json
{
  "content": "string"
}
```

**Content Fields**:

| Field     | Type   | Required | Limits | Example                        |
| --------- | ------ | -------- | ------ | ------------------------------ |
| `content` | string | Yes      | 1-1000 | `Hi, this is my first post :)` |

**Response Body Schema**:

```json
{
  "post_id": "int",
  "user_id": "int",
  "content": "string",
  "likes": "int",
  "reposts": "int",
  "created_at": "string",
  "original_post_id": null,
  "original_post": null
}
```

## **/{username}/posts {GET}**

**Description**: Get the user's posts

**Response Body Schema**:

```json
[
  {
    "post_id": "int",
    "user_id": "int",
    "content": "string",
    "likes": "int",
    "reposts": "int",
    "created_at": "string",
    "original_post_id": null,
    "original_post": null
  },
  {
    "post_id": "int",
    "user_id": "int",
    "content": "string",
    "likes": "int",
    "reposts": "int",
    "created_at": "string",
    "original_post_id": null,
    "original_post": null
  }
]
```

## **/{username}/posts/{post_id} {GET}**

**Description**: Get the user's post by ID

**Response Body Schema**:

```json
{
  "post_id": "int",
  "user_id": "int",
  "content": "string",
  "likes": "int",
  "reposts": "int",
  "created_at": "string",
  "original_post_id": null,
  "original_post": null
}
```

## **/{username}/posts/{post_id} {PATCH}**

**Description**: Change the content of the post by ID

**Request Body Schema**:

```json
{
  "content": "string"
}
```

**Content Fields**:

| Field     | Type   | Required | Limits | Example                           |
| --------- | ------ | -------- | ------ | --------------------------------- |
| `content` | string | Yes      | 1-1000 | `Hi, this is my modified post :)` |

**Response Body Schema**:

```json
{
  "post_id": "int",
  "user_id": "int",
  "content": "string",
  "likes": "int",
  "reposts": "int",
  "created_at": "string",
  "original_post_id": null,
  "original_post": null
}
```

## **/{username}/posts/{post_id} {DELETE}**

**Description**: Delete post by ID

**Response Body Schema**:

```json
{
  "message": "successfully deleted the post"
}
```

## **/{username}/posts/{post_id}/like {POST}**

**Description**: Like the post by ID

**Response Body Schema**:

```json
{
  "message": "successfully liked the post",
  "post_id": "int"
}
```

## **/{username}/posts/{post_id}/like {DELETE}**

**Description**: Unlike the post by ID

**Response Body Schema**:

```json
{
  "message": "successfully unliked the post",
  "post_id": "int"
}
```

## **/{username}/reposts {GET}**

**Description**: Get user reposts

**Response Body Schema**:

```json
[
  {
    "post_id": "int",
    "user_id": "int",
    "content": "string",
    "likes": "int",
    "reposts": "int",
    "created_at": "string",
    "original_post_id": null,
    "original_post": null
  },
  {
    "post_id": "int",
    "user_id": "int",
    "content": "string",
    "likes": "int",
    "reposts": "int",
    "created_at": "string",
    "original_post_id": null,
    "original_post": null
  }
]
```

## **/{username}/posts/{post_id}/repost {POST}**

**Description**: Repost another user's post by ID

**Response Body Schema**:

```json
{
  "message": "successfully reposted the post",
  "post_id": "int"
}
```

## **/{username}/posts/{post_id}/repost {DELETE}**

**Description**: Undo repost another user's post by ID

**Response Body Schema**:

```json
{
  "message": "successfully cancelled repost the post",
  "post_id": "int"
}
```

## **/{username}/posts/{post_id}/quote {POST}**

**Description**: Quote post by ID

**Request Body Schema**:

```json
{
  "content": "string"
}
```

| Field     | Type   | Required | Limits | Example    |
| --------- | ------ | -------- | ------ | ---------- |
| `content` | string | Yes      | 1-1000 | `My quote` |

**Response Body Schema**:

```json
{
  "post_id": "int",
  "user_id": "int",
  "content": "string",
  "likes": "int",
  "reposts": "int",
  "created_at": "string",
  "original_post_id": "int",
  "original_post": {
    "post_id": "int",
    "user_id": "int",
    "content": "string",
    "likes": "int",
    "reposts": "int",
    "created_at": "string",
    "original_post_id": null,
    "original_post": null
  }
}
```

# 🚧 Planned

## **/auth/refresh {POST}**

**Description**: Get the new Access token

## **/auth/logout {POST}**

**Description**: Log out of the account

## **/feed {GET}**

**Description**: News feed (following users' posts)

## **/notifications {GET}**

**Description**: User notifications
