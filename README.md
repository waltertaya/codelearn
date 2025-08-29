# CodeLearn Platform

**Code Locally, Learn Globally** - A revolutionary programming learning platform that allows developers to code in their favorite IDE while submitting solutions via CLI and receiving instant feedback.

## 🌟 Overview

CodeLearn is a comprehensive programming education platform consisting of:

- **Backend API**: Go-based REST API with authentication and challenge management
- **CLI Client**: Python-based command-line tool for seamless code submission
- **Auto-Grading System**: Automated code evaluation with detailed feedback



## 🏗️ Architecture



### Backend (Go + Gin + SQLite)
- RESTful API with JWT authentication
- SQLite database for data persistence
- CORS enabled for frontend integration
- Comprehensive error handling

### CLI Client (Python)
- Cross-platform command-line interface
- Automatic language detection from file extensions
- Secure token-based authentication
- Rich terminal output with emojis and formatting

## 🛠️ Technology Stack

### Backend
- **Language**: Go 1.21+
- **Framework**: Gin Web Framework
- **Database**: SQLite
- **Authentication**: JWT tokens
- **Password Hashing**: bcrypt
- **CORS**: gin-contrib/cors

### CLI
- **Language**: Python 3.8+
- **HTTP Client**: requests
- **Configuration**: JSON file storage
- **Cross-platform**: Works on Windows, macOS, Linux

## 🚀 Quick Start

### Prerequisites
- Go 1.21+ (for backend development)
- Python 3.8+ (for CLI usage)
- Node.js 18+ (for frontend development)

### Backend Setup
```bash
cd codelearn-backend
go mod tidy
# go build -o codelearn-backend
# ./codelearn-backend
# go run .
go run migrate/migrate.go
go run cmd/main.go
```

### CLI Usage
```bash
cd codelearn-cli
pip install .

# Register a new account
codelearn register username email@example.com password123

# List available challenges
codelearn challenges

# Submit a solution
codelearn submit 1 solution.py

# View your submissions
codelearn submissions

# Check leaderboard
codelearn leaderboard
```



## 📊 API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Refresh token

### Protected Endpoints (require JWT token)
- `GET /api/v1/profile` - Get user profile
- `PUT /api/v1/profile` - Update user profile
- `GET /api/v1/challenges` - List challenges
- `GET /api/v1/challenges/:id` - Get specific challenge
- `POST /api/v1/challenges/:id/submit` - Submit solution
- `GET /api/v1/submissions` - List user submissions
- `GET /api/v1/submissions/:id` - Get specific submission
- `GET /api/v1/leaderboard` - Get leaderboard
- `POST /api/v1/cli/auth` - CLI authentication

## 🎯 Sample Challenges

The platform comes with 5 pre-loaded challenges:

1. **Two Sum** (Easy, Python) - Array manipulation
2. **Reverse String** (Easy, JavaScript) - String operations
3. **Binary Search** (Medium, Go) - Algorithm implementation
4. **Valid Parentheses** (Easy, Python) - Stack operations
5. **Fibonacci Sequence** (Easy, JavaScript) - Mathematical sequences

## 🔧 CLI Commands

### User Management
```bash
# Register new account
codelearn register <username> <email> <password>

# Login to existing account
codelearn login <username> <password>
```

### Challenge Management
```bash
# List all challenges
codelearn challenges

# Filter by difficulty
codelearn challenges --difficulty Easy

# Filter by language
codelearn challenges --language python
```

### Solution Submission
```bash
# Submit solution (auto-detect language)
codelearn submit <challenge_id> <file_path>

# Submit with specific language
codelearn submit <challenge_id> <file_path> --language python
```

### Progress Tracking
```bash
# View your submissions
codelearn submissions

# View leaderboard
codelearn leaderboard
```

## 🔒 Security Features

- **JWT Authentication**: Secure token-based authentication
- **Password Hashing**: bcrypt for secure password storage
- **CORS Protection**: Configured for frontend integration
- **Input Validation**: Comprehensive request validation
- **SQL Injection Prevention**: Parameterized queries

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.
