# Match-Me GraphQL API Documentation 🚀

This document provides comprehensive information about the GraphQL API implementation in the Match-Me application.

## 🔧 GraphQL Stack

The GraphQL API is built using:

- **[gqlgen](https://github.com/99designs/gqlgen)** v0.17.78 - Go GraphQL server library
- **PostgreSQL + PostGIS** - Database with geospatial support
- **WebSocket Subscriptions** - Real-time features
- **JWT Authentication** - Secure user sessions
- **Cloudinary Integration** - Image upload and management

## 📍 GraphQL Endpoints

### Development
- **GraphQL Playground**: `http://localhost:8088/graphql` (GET)
- **GraphQL API**: `http://localhost:8088/graphql` (POST)
- **WebSocket Subscriptions**: `ws://localhost:8088/graphql-ws` (WebSocket)

### Authentication
Most operations require JWT authentication via the `Authorization` header:
```
Authorization: Bearer <your-jwt-token>
```

## 🔑 Core Schema Types

### User Management
- **User** - Core user profile 
- **Bio** - Extended profile information
- **Profile** - Additional profile metadata


### Matching & Connections
- **Recommendations** - Recommended users profiles
- **ConnectionRequests** - Match requests between users
- **Connection** - Established connections between matched users


## 🚀 Key Features

### 1. Authentication & User Management
To help with the review, there are some querries provided.

```graphql
# Register a new user
mutation RegisterUser{
  registerUser(email:"tiia@example.com", password:"12345a"){
    userID
  }
}

# Login
mutation{
  loginUser(email:"graph1@gmail.com", password:"12345a"){
    token
    user{
      userID
    }
  }
}
```

### 2. User, Profile, Bio & Photo Management
```graphql
For some querries variable "userID" should be provided.
VAR {"userID": "9b6a2aa5-5ea8-4ada-a850-37e849ce498c"}
For all querries need to provide token. 
Header {
  "Authorization": "Bearer <your_token>"
}

query GetProfile($userID: ID!) {
  profile(userID: $userID) {
  userID,
  name,
  about,
  childName,
  ChildInterests,
    lat
  }
}

query GetUser($userID: ID!) {
  user(id: $userID) {
  userID,
  email,
  profilePicture
    profile{
      name,
      languages,
      childName,
      addressCity,
      ChildInterests
    }
    bio{
      childBirthday,
      play_styles
    }
  }
}



query GetBio($userID: ID!) {
  bio(userID: $userID) {
  userID,
  parentGender,
  limitations,
  preferredDistance,
    user{
      email
    }
  }
}

query GetMeProfile {
  myProfile{
  userID,
  name,
  about,
  childName,
  childInterests,
    lat
  }
}


mutation UpdateBio {
  updateBio(
    parentGender: "FEMALE"
    preferredDistance: 15
    childBirthday: "2020-09-17"
    childActivity_level: "HIGH"
    limitations: ["None"]
    allergies: ["Peanuts"]
    play_styles: ["Outdoor", "Creative"]
  ) {
    parentGender
    preferredDistance
    childBirthday
    childGender
    childActivity_level
    limitations
    allergies
    play_styles
  }
}


mutation UpdateProfile {
  updateProfile(
    name: "Irina"
    childInterests: ["reading", "writing"]
  ) {
    userID
    name
    about
    languages
    addressCity
    lat
    lon
    childName
    childAbout
    childInterests
  }
}

mutation DeletePhoto{
  deletePhoto
}

```

### 3. Discovery & Matching
```graphql
# Get personalized recommendations
query GetRec{
  recommendations {
    userID
    bio{
      limitations
      childGender
    }
    profile{
      name
      childAbout
    }
  }
}
```

### 4. Connection Management
```graphql
mutation UpsertReaction{
  upsertReaction(targetedUserID: "9b6a2aa5-5ea8-4ada-a850-37e849ce498c", reaction: like)
}

query GetCon{
  connections {
    userID
    bio{
      preferredDistance
      childBirthday
      childGender
    }
    profile{
      name
      childAbout
      childInterests
    }
  }
}
```

## 📡 Real-time Subscriptions

The API supports WebSocket subscriptions for real-time features:

### Message Subscriptions
```graphql
subscription {
  onNewMessage(chatID: "5493a4e8-7ebe-49b7-b8c9-b3108bf2a542") {
    id
    content
    sender { userID, email }
  }
}

mutation send {
  sendMessage(chatID: "5493a4e8-7ebe-49b7-b8c9-b3108bf2a542", content: "Hello, this is a test!") {
    id
    content
    sender{
      email
    }
  }
}

```

## 🔧 Configuration

### Environment Variables
```bash
# GraphQL-specific settings
APP_ENV=development         # Enables GraphQL Playground
PORT=8088                   # Server port
JWT_SECRET=your-secret-key  # JWT signing key
```

### GraphQL Server Configuration
- **Introspection**: Enabled in development
- **WebSocket Transport**: Real-time subscriptions


## 🧪 Testing & Development

### GraphQL Playground
In development mode, open the interactive GraphQL Playground at http://localhost:8088/playground.
If you're using Altair or Apollo Sandbox, access the endpoint at http://localhost:8088/graphql. There you can:
- Explore the schema
- Test queries and mutations
- View documentation
- Debug subscriptions

### Example Queries
The playground includes example queries for all major operations:
- User registration/login
- Profile, Bio management
- Connection requests
- Recommendations
- Connections
- Real-time subscriptions

## 🔄 Integration with Frontend

The GraphQL API is designed to work seamlessly with the React frontend:


For more information about the overall application architecture, see the main [README.md](./README.md).