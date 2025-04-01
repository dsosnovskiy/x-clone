# 🔥 X-CLONE - A simplified copy of the social network X (Twitter)

## 🛠 Technologies

- **Programming language**: `Golang`
- **Architecture**: `Layered`
- **Web framework**: `net/http` (HTTP server), `chi` (HTTP router)
- **Configuration**: `cleanenv`
- **Logging**: `logrus`
- **Database**: `PostgreSQL` (using GORM)
- **Authentication**: `JWT` (Access only)

## 🎯 Goals

- **Caching**: `Redis`
- **Secure**: `Refresh Token`
- **New functionality**:
  - Edit profile
  - News feed
  - Notifications
- **Containerization**: `Docker`

# 🔍 API Description

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

# 👤 User Endpoints

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

```
successfully follow user: {username}
```

## **/{username}/follow {DELETE}**

**Description**: Stop following another user

**Response Body Schema**:

```
successfully stop following user: {username}
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

# 📝 Post Endpoints

## **/compose/post {POST}**

**Description**: Create post

**Response Body Schema**:

```json
{
  "post_id": "int",
  "user_id": "int",
  "content": "string",
  "likes": "int",
  "reposts": "int",
  "created_at": "string",
  "updated_at": "string",
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
    "updated_at": "string",
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
    "updated_at": "string",
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
  "updated_at": "string",
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

**Response Body Schema**:

```json
{
  "post_id": "int",
  "user_id": "int",
  "content": "string",
  "likes": "int",
  "reposts": "int",
  "created_at": "string",
  "updated_at": "string",
  "original_post_id": null,
  "original_post": null
}
```

## **/{username}/posts/{post_id} {DELETE}**

**Description**: Delete post by ID

**Response Body Schema**:

```
successful deletion of the post
```

## **/{username}/posts/{post_id}/like {POST}**

**Description**: Like the post by ID

**Response Body Schema**:

```
successful liking of the post: {post_id}
```

## **/{username}/posts/{post_id}/like {DELETE}**

**Description**: Unlike the post by ID

**Response Body Schema**:

```
successful unliking of the post: {post_id}
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
    "updated_at": "string",
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
    "updated_at": "string",
    "original_post_id": null,
    "original_post": null
  }
]
```

## **/{username}/posts/{post_id}/repost {POST}**

**Description**: Repost another user's post by ID

**Response Body Schema**:

```
successful repost of the post: {post_id}
```

## **/{username}/posts/{post_id}/repost {DELETE}**

**Description**: Undo repost another user's post by ID

**Response Body Schema**:

```
successful undo repost of the post: {post_id}
```

## **/{username}/posts/{post_id}/quote {POST}**

**Description**: Quote post by ID

**Response Body Schema**:

```json
{
  "post_id": "int",
  "user_id": "int",
  "content": "string",
  "likes": "int",
  "reposts": "int",
  "created_at": "string",
  "updated_at": "string",
  "original_post_id": "int",
  "original_post": {
    "post_id": "int",
    "user_id": "int",
    "content": "string",
    "likes": "int",
    "reposts": "int",
    "created_at": "string",
    "updated_at": "string",
    "original_post_id": null,
    "original_post": null
  }
}
```

# 🚧 Planned Endpoints

## **/auth/refresh {POST}**

**Description**: Get the new Access token

## **/auth/logout {POST}**

**Description**: Log out of the account

## **/setting/profile {PATCH}**

**Description**: Update user profile

## **/feed {GET}**

**Description**: News feed (following users' posts)

## **/notifications {GET}**

**Description**: User notifications
