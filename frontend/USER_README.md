# Match-Me Children - Playdate Connection Platform

A full-stack web application that helps parents find compatible playdate matches for their children based on interests, location, and preferences.

## 🎯 Overview

Match-Me Children connects parents looking for playdate opportunities. Users create profiles for themselves and their children, receive intelligent recommendations based on matching algorithms, and can chat in real-time with their connections.

## 🛠️ Tech Stack

### Frontend
- **React 18** with **TypeScript**
- **React Router** for navigation
- **Bulma CSS** for styling
- **WebSocket** for real-time chat
- **Vite** as build tool

### Backend
- **Go (Golang)** with **Gin** web framework
- **PostgreSQL** database
- **PostGIS** for proximity-based location filtering
- **JWT** for authentication
- **bcrypt** for password hashing
- **WebSocket (Gorilla)** for real-time features

## 📋 Features

### ✅ Implemented (All Mandatory + Extra)

#### **Authentication & User Management**
- Secure registration with email and password (bcrypt + salt)
- JWT-based session management
- Logout from any page
- Profile completion requirement (100% before recommendations)

#### **User Profiles**
- Parent profile (name, about, languages, location, photo)
- Child profile (name, age, gender, interests, activity level, play styles, allergies)
- Profile picture upload/change/remove
- Modify profile anytime
- Email privacy (never exposed in API)

#### **Matching Algorithm**
- **5+ biographical data points** used for matching:
  1. Child interests
  2. Activity level
  3. Play styles
  4. Age compatibility
  5. Location proximity
- Customizable matching weights
- Proximity-based filtering using PostGIS (extra feature)
- Prioritized recommendations (best matches first)
- Dismiss functionality (dismissed users never shown again)
- Maximum 10 recommendations at a time

#### **Connections**
- Send connection requests
- Accept/reject incoming requests
- View all connections
- Disconnect option

#### **Real-Time Chat** (No Polling!)
- WebSocket-based messaging
- Chats ordered by most recent activity
- Message timestamps
- Unread message indicators (real-time)
- Chat pagination
- Start chat from connected user's profile
- **Online/offline indicator** (extra feature)
- **Typing indicator** (extra feature)

#### **RESTful API**
- `/recommendations` - returns only IDs
- `/connections` - returns only IDs
- `/users/:id` - name and profile picture
- `/users/:id/profile` - full profile
- `/users/:id/bio` - biographical data
- `/me/*` - shortcuts for authenticated user
- HTTP 404 for not found/no access

#### **Responsive Design**
- Mobile and desktop optimized
- Clean, modern UI

## 🚀 Installation & Setup

### Prerequisites
- Go 1.21+
- Node.js 18+
- PostgreSQL 14+ with PostGIS extension
- Git

### 1. Clone the Repository
```bash
git clone <repository-url>
cd match-me
```

### 2. Database Setup
```bash
# Create database
psql -U postgres
CREATE DATABASE matchme;
\q

# Run schema
psql -U postgres -d matchme -f server/queries.txt

# Create PostGIS extension (if not exists)
psql -U postgres -d matchme
CREATE EXTENSION IF NOT EXISTS postgis;
\q
```

### 3. Backend Setup
```bash
cd server

# Create .env file
cat > .env << EOF
DATABASE_URL=postgres://postgres:yourpassword@localhost:5432/matchme
JWT_SECRET=your-secret-key-min-32-characters
CLOUDINARY_CLOUD_NAME=your-cloudinary-name
CLOUDINARY_API_KEY=your-api-key
CLOUDINARY_API_SECRET=your-api-secret
EOF

# Install dependencies
go mod download

# Run backend
go run main.go
```

Backend will run on `http://localhost:8088`

### 4. Frontend Setup
```bash
cd ../frontend

# Install dependencies
npm install

# Create .env file
cat > .env << EOF
VITE_API_BASE_URL=http://localhost:8088
EOF

# Run frontend
npm run dev
```

Frontend will run on `http://localhost:5173`

## 🧪 Testing Setup

### Load 100 Test Users

We've created 100 fictitious users with diverse profiles for testing:
```bash
cd server
go run scripts/seed_users.go
```

**Test User Credentials:**
- **Emails:** `test1@example.com` through `test100@example.com`
- **Password:** `password123` (all users)

**Test Users Include:**
- Various locations across Finnish cities
- Different child ages (1-10 years)
- Diverse interests and activity levels
- Mixed languages (Finnish, English, Swedish)

### Reset Database

To drop and recreate:
```bash
# Drop database
psql -U postgres -c "DROP DATABASE IF EXISTS matchme;"

# Recreate
psql -U postgres -c "CREATE DATABASE matchme;"

# Run schema
psql -U postgres -d matchme -f server/queries.txt

# Reload test users
cd server && go run scripts/seed_users.go
```

## 📝 Testing Checklist

### Authentication
- [ ] Register with email and password
- [ ] Login with credentials
- [ ] Logout from any page
- [ ] JWT persists across page refreshes

### Profile Management
- [ ] Cannot see recommendations until profile 100% complete
- [ ] Can set/change/remove profile picture
- [ ] Can edit parent profile (name, about, languages, location)
- [ ] Can edit child profile (name, age, interests, etc.)
- [ ] Email not visible to other users

### Matching & Recommendations
- [ ] Only see users from same location (if not using proximity)
- [ ] Get maximum 10 recommendations at a time
- [ ] Recommendations prioritized by score
- [ ] Dismissed users don't reappear
- [ ] Poor matches not recommended
- [ ] Good matches appear in recommendations

### Connections
- [ ] Can send connection request from recommendation
- [ ] Can view incoming requests
- [ ] Can accept connection request
- [ ] Can reject connection request
- [ ] Can disconnect from connected user
- [ ] Profile visible only when: recommended, request pending, or connected

### Chat (Real-Time)
- [ ] Can only chat with connected users
- [ ] Start chat from connected user's profile
- [ ] Messages appear instantly (real-time)
- [ ] Chats ordered by most recent activity
- [ ] Unread indicator appears in real-time
- [ ] Message timestamps displayed
- [ ] Chat history paginated
- [ ] Both users see same history
- [ ] **No polling** - WebSocket only
- [ ] Online/offline indicator works
- [ ] Typing indicator appears when user types

### API Endpoints
- [ ] `/recommendations` returns only IDs
- [ ] `/connections` returns only IDs
- [ ] `/users/:id` returns name and photo only
- [ ] `/users/:id/profile` returns full profile
- [ ] `/users/:id/bio` returns bio data
- [ ] `/me` shortcuts work
- [ ] HTTP 404 for not found/no permission
- [ ] Email never returned in API

### Security
- [ ] Passwords hashed with bcrypt
- [ ] JWT required for protected routes
- [ ] Cannot access other users' private data
- [ ] Cannot see profiles without permission

### Responsive Design
- [ ] Works on desktop browsers
- [ ] Works on mobile browsers
- [ ] UI adapts to screen size


## 🎯 Key Implementation Details

### Matching Algorithm
- Uses **weighted scoring system**
- Users can customize weights for different factors
- **PostGIS** for efficient proximity filtering
- Prevents weak recommendations (<50% match score)

### Real-Time Chat
- **WebSocket** implementation (Gorilla WebSocket library)
- Hub pattern for managing connections
- Separate channels per user
- Broadcasts to specific recipients only

### Database
- **PostgreSQL** with **PostGIS** extension
- Spatial indexes for fast location queries
- Triggers for automatic chat ordering
- CASCADE deletes for data integrity

### Security
- **bcrypt** password hashing (cost 10)
- **JWT** with 24-hour expiration
- Route-level authentication middleware
- Email privacy protection

## 🐛 Troubleshooting

### Backend won't start
```bash
# Check if port 8088 is available
lsof -i :8088

# Check database connection
psql -U postgres -d matchme -c "SELECT 1;"
```

### Frontend won't start
```bash
# Clear node_modules and reinstall
rm -rf node_modules package-lock.json
npm install
```

### WebSocket not connecting
- Ensure backend is running
- Check that `VITE_API_BASE_URL` is correct
- Look for CORS errors in browser console

### No recommendations showing
- Ensure user profile is 100% complete
- Check there are other users in the database
- Verify location matching is correct

## 📊 Performance Notes

- **PostGIS spatial indexing** enables fast proximity queries
- **WebSocket** reduces server load vs polling (no constant HTTP requests)
- **Pagination** prevents large payload transfers
- **JWT** eliminates database session lookups





**For Testing:** Use test accounts `test1@example.com` through `test100@example.com`, password: `password123`

