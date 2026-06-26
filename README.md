# 💛 Match-Me: Find the Perfect Playmate

**Match-Me** a new demo app designed for parents who want to find great play partners for their kids.

As parents, we all know how hard it can be to meet other families whose children share similar energy, interests, and schedules. **Match-Me**  helps make that easier — by connecting nearby parents based on children’s personalities, play styles, and activity levels, using a bit of smart matching magic under the hood.

Whether you’re new in town or just looking to make new friends for your child, this app helps parents discover like-minded families and chat in real time. 

![alt text](image.png)

> 🧸 In this demo, each parent can currently add **one child** profile to try out the experience.

### How Matching Works (User Perspective)

When you look for new playmates in **Match-Me**, the app doesn’t just show random families nearby — it actually analyzes how well your child’s personality, interests, and play style fit with others.

Every family creates a **Parent Profile**  and one **Child Profile**.  
The app then calculates **compatibility scores** between families.

The system:
- Shows up to 10 recommendations at a time
- Ranks matches by score (highest first)
- Filters out already dismissed profiles
- Only shows users from your selected location

#### Behind the scenes:

1. 🗺️ **Nearby families are pre-filtered** by your preferred distance (via PostGIS).  
2. 💬 **Your preferences** — such as how important shared interests or play style are — are applied as **weights**.  
3. 💡 Each potential match gets a **score**:
   - Shared interests and play styles → higher score  
   - Similar activity levels → higher score  
   - Smaller age gap → higher score  
   - Conflicting allergies or limitations → lower score
4. 🤝 **Mutual compatibility** is checked both ways — you must match *each other’s* preferences.
5. 🔢 **Results are sorted** by final compatibility.

> Example:  
> If your 5-year-old loves building Lego and has medium energy, you’ll likely see families with similar-aged, creative, and moderately active kids near you.


### 🧩 Matching Algorithm (Developer Perspective)

Matching logic lives in `server/services/matching.go` and works in two main stages:  
**prefiltering (SQL/PostGIS)** and **scoring (Go logic)**.

#### 1. Prefiltering (SQL)

Selects only eligible candidates near the user who haven’t been liked/disliked before:

```sql
WITH pf AS (
  SELECT candidate_user_id
  FROM prefilter_candidates_postgis($1::uuid, 200, 0)
)
SELECT ...
FROM pf
JOIN parent_profiles pp ON pp.user_id = pf.candidate_user_id
JOIN children c ON c.user_id = pp.user_id
WHERE pp.user_id <> $1::uuid
  AND pp.user_id NOT IN (
    SELECT target_user_id FROM user_reactions WHERE user_id = $1::uuid
    UNION
    SELECT user_id FROM user_reactions WHERE target_user_id = $1::uuid
  );
```

  
#### 2. Compatibility Scoring (Go)

Each candidate is evaluated across multiple dimensions.
The main scoring function is CalculateMatchingScores().

Factor	Function	Description	Weight Source
🎨 Interests	CalculateArrayOverlap()	Shared hobbies/interests	prefs.InterestsWeight
⚡ Activity Level	CalculateActivityCompatibility()	Energy match (low/medium/high)	prefs.ActivityLevelWeight
💬 Limitations	CalculateArrayOverlap()	Similar needs or conditions	prefs.LimitationsWeight
🌿 Allergies	CalculateAllergiesCompatibility()	Checks overlap	prefs.AllergiesWeight
🧸 Play Styles	CalculateArrayOverlap()	Creative / outdoor / quiet, etc.	prefs.PlayStylesWeight
⏳ Age	CalculateAgeCompatibility()	Penalizes large age gaps	Always active

Final score:

finalScore = totalWeightedScore / totalWeight

Then mutual scoring:

score1 := CalculateCompatibilityScore(currentUser, candidate, currentUserPrefs)
score2 := CalculateCompatibilityScore(candidate, currentUser, candidatePrefs)
finalScore := (score1 + score2) / 2

The result is normalized to 0.0 – 1.0.

---

## 🚀 Getting Started

Follow these instructions to get a local copy up and running for development and testing.

### Prerequisites

Ensure you have the following installed on your system:

- **Node.js**: `v18` or later
- **Go**: `v1.24` or later
- **PostgreSQL**

### Installation & Setup

#### 1. **Clone the Repository**

    ```bash
    git clone https://gitea.kood.tech/irinapanivenko/match-me.git
    cd match-me
    ```

#### 2.  **Set Up the Database**
    This project uses PostgreSQL with the PostGIS extension for location services.

    ```bash
    # On macOS (using Homebrew)
    brew install postgresql postgis

    # On Debian/Ubuntu
    sudo apt-get update
    sudo apt-get install postgresql postgresql-contrib postgis

    # Create the database and enable the extension
    createdb matchme
    ```

#### 3.  **Configure Environment Variables**
    Create a `.env` file inside the `server/` directory and populate it with your configuration.

    ```bash
    # server/.env
    # Enables GraphQL Playground
    APP_ENV=development         
    # Server port
    PORT=8088

    # PostgreSQL connection string
    # Replace username, password, dbname with your own
    DATABASE_URL=postgres://YOUR_USER:YOUR_PASSWORD@localhost:5432/matchme?

    # Secret for signing JWT tokens (choose something random & long)
    JWT_SECRET=supersecretkey123

    # Allowed CORS origin (your React app)
    CORS_ORIGIN=http://localhost:5173

    CLOUDINARY_CLOUD_NAME=
    CLOUDINARY_API_KEY=
    CLOUDINARY_API_SECRET=
    ```

    Create a `.env` file inside the `frontend/` directory and add:

    ```bash
    VITE_API_BASE_URL=http://localhost:8088

    #You need to obtain an API key from https://www.geoapify.com/ 
    VITE_GEOAPIFY_KEY=
    ```
    > Remember that the port set in both the client and server must match, for example if you change the `PORT` in `server/.env` to 3000 then the client should be http://localhost:3000

#### 4.  **Install Dependencies and Run**

    ```bash
    # Install dependencies for both client and server
    make dev-setup
    ``` 
    Run the frontend and backend servers in **two separate terminal windows**:

    ```bash
    # In terminal 1: Start the client dev server
    make dev-client
    # Frontend will be available at http://localhost:5173
    ```

    ```bash
    # In terminal 2: Start the backend server
    make run-server
    # Backend will be running at http://localhost:8088
    ```

#### 5.  if you do not have make installed

    ```bash
    # In terminal 1: Install the client and start the frontend dev server
    cd frontend
    npm i && npm run dev
    # Frontend will be available at http://localhost:5173
    ```
> 📖 **For detailed GraphQL API documentation, see [GraphQL API documentation](./GraphQL-README.md)**

### Database Management

    The backend includes helpful commands for managing the database during development or testing.

    ```bash
    # Navigate to the server directory
    cd server

    # Seed the database with 100 test users(password "password123")
    go run . -seed

    # Completely reset the database (drop all data)
    go run . -drop
    ```
## ⚙️ Usage

### Makefile Commands

A `Makefile` at the root of the project simplifies common tasks.

| Command | Description |
| :--- | :--- |
| `make help` | Displays a list of all available commands. |
| `make dev-setup` | Installs all dependencies for both client and server. |
| `make dev-client` | Starts the frontend development server with hot-reloading. |
| `make run-server`| Starts the backend API server. |
| `make build-all` | Creates production-ready builds for both client and server. |

-----

## ✨ Features

  * **✅ Secure User Authentication**: JWT-based authentication for secure sessions and profile management.
  * **📍 Geospatial Matching**: Utilizes PostGIS to discover and connect with potential matches nearby.
  * **💬 Real-Time Chat**: Instant messaging between connected users, powered by WebSockets for a fluid conversation experience.
  * **📸 Cloud-Based Image Handling**: Efficient and secure photo uploads and storage managed via Cloudinary.
  * **🤝 Connection Management**: A complete system to send, accept, and manage connection requests.
  * **📱 Fully Responsive Design**: A beautiful and intuitive interface that works flawlessly on both desktop and mobile devices.

-----

## 🛠️ Tech Stack

The project is built with a modern and robust technology stack, separating concerns between a client-side application and a server-side API.

| **Component** | **Technology** | **Purpose** |
| :--- | :--- | :--- |
| **Frontend** | React 19 (TypeScript) | UI development |
| | Vite | Build tooling & dev server |
| | Zustand & React Query | State management & server-state synchronization |
| | Bulma | Styling and responsive layout |
| | WebSocket API | Real-time communication |
| **Backend** | Go 1.24+ | Core application logic |
| | Gin | High-performance HTTP web framework |
| **Database** | PostgreSQL + PostGIS | Relational data and geospatial queries |
| **Infrastructure** | JWT | Authentication |
| | Cloudinary | Cloud-based image storage |


-----

## 📁 Project Structure

The repository is organized into two main parts: a `frontend` directory for the frontend and a `server` directory for the backend.

```bash
match-me
├── server
│   ├── database       # Database connection & initialization
│   ├── endpoints      # API route definitions
│   ├── handlers       # Request handlers & logic
│   ├── helpers        # Utility functions (tokens, validation, etc.)
│   ├── internal       # Internal reusable logic
│   ├── middleware     # Auth, CORS, and logging middleware
│   ├── services       # Core services (matching, Cloudinary, etc.)
│   ├── sqlfiles       # SQL schema & seed scripts
│   ├── structs        # Data models and DTOs
│   └── websocket      # Real-time chat & notifications
│
└── frontend
    ├── node_modules   # Dependencies
    ├── public         # Static assets (index.html, icons)
    └── src
        ├── api        # API calls & request helpers
        ├── assets     # Images and media
        ├── auth       # JWT logic & route protection
        ├── components # Reusable UI components
        ├── hooks      # Custom React hooks
        ├── pages      # Main application views
        ├── styles     # Bulma & custom CSS
        └── types      # TypeScript interfaces & types
```
