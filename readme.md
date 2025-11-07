# FitWell API

FitWell is a comprehensive fitness and healthcare platform built with Go, combining AI-powered workout generation with medical consultation services.

## Features

### Fitness Module
- **AI-Powered Workout Generation** - Generate personalized workout plans using Google Gemini AI
- **Exercise Management** - Create, update, and track exercises
- **Exercise Logging** - Log workout activities with sets, reps, and weight
- **BMI Calculator** - Calculate BMI using RapidAPI with fallback mechanism

### Medical Module
- **Doctor Registration & Authentication** - Separate authentication for medical professionals
- **Appointment System** - Book appointments with doctors
- **Diagnosis Management** - Create and update patient diagnoses
- **Medication Prescription** - Prescribe medications with automatic cost calculation
- **Billing System** - Automatic billing generation based on consultation and medication costs
- **Invoice Generation** - Auto-create invoices when payment is confirmed

## Tech Stack

- **Language**: Go 1.21+
- **Framework**: Echo v4
- **Database**: PostgreSQL with GORM
- **Authentication**: JWT
- **Testing**: Go testing with mockgen (go.uber.org/mock)
- **External APIs**: 
  - Google Gemini AI (Workout Generation)
  - RapidAPI BMI Calculator

## Prerequisites

- Go 1.21 or higher
- PostgreSQL
- Google Gemini API Key
- RapidAPI Key (for BMI calculator)

## Installation

1. Clone the repository
```bash
git clone https://github.com/Dedenruslan19/med-project.git
cd med-project
```

2. Install dependencies
```bash
go mod download
```

3. Set up environment variables
```bash
cp .env.example .env
```

JWT_SECRET=your_jwt_secret_key
GEMINI_API_KEY=your_gemini_api_key
RAPIDAPI_KEY=your_rapidapi_key
```

4. Set up the database
```bash
# Create database
createdb fitconnect_db

# Run migrations (the application will auto-migrate on start)
go run cmd/echo-server/main.go
```

## Running the Application

### Development
```bash
go run cmd/echo-server/main.go
```

The server will start on `http://localhost:8080`

### Production Build
```bash
go build -o fitconnect-api cmd/echo-server/main.go
./fitconnect-api
```

## Project Structure

```
.
├── cmd/
│   └── echo-server/
│       ├── main.go              # Application entry point
│       ├── controller/          # HTTP handlers
│       └── middleware/          # JWT & validation middleware
├── service/                     # Business logic layer
│   ├── appointments/
│   ├── billings/
│   ├── diagnoses/
│   ├── doctors/
│   ├── exercises/
│   ├── invoices/
│   ├── logs/
│   ├── users/
│   └── workouts/
├── repository/                  # Data access layer
│   ├── appointment/
│   ├── billing/
│   ├── diagnose/
│   ├── doctor/
│   ├── exercise/
│   ├── gemini/                 # Gemini AI integration
│   ├── invoice/
│   ├── logs/
│   ├── rapidAPI/               # RapidAPI BMI integration
│   ├── user/
│   └── workout/
├── util/                       # Utility functions
│   └── jwt.go                  # JWT helper
├── ddl.sql                     # Database schema
├── .yaml                       # OpenAPI specification
├── diagrams.md                 # PlantUML diagrams
└── coverage.html               # Test coverage report
```

## Database Schema

The database uses a normalized 3NF structure. See [ddl.sql](./ddl.sql) for the complete schema.

Key tables:
- `users` - User accounts
- `doctors` - Doctor accounts
- `workouts` - Workout plans
- `exercises` - Exercise details
- `exercise_logs` - Exercise activity tracking
- `appointments` - Medical appointments
- `diagnoses` - Medical diagnoses
- `billings` - Billing information
- `invoices` - Invoice details

## Business Process Flow

### 1. Fitness Flow
```
User Registration → BMI Calculation → AI Workout Generation (Gemini) 
→ Save Workout → Log Activities
```

### 2. Medical Flow
```
User Books Appointment → Doctor Creates Diagnosis → Auto-generate Billing 
→ User Pays → Auto-create Invoice
```

### 3. Payment Calculation
- **Consultation Fee**: Rp 200,000 (fixed)
- **Medication Cost**: Rp 50,000 per medication
- **Total**: Consultation + (Medication Count × 50,000)

## Key Features Implementation

### Auto-Invoice Creation
When billing payment status changes to "paid", the system automatically:
1. Creates an invoice record
2. Calculates consultation fee (Rp 200,000)
3. Calculates medication costs from diagnosis
4. Stores complete invoice with user, doctor, and appointment details

### AI Workout Generation
Uses Google Gemini AI to generate 3-5 exercises based on:
- Workout name/target
- User goals
- Available equipment

### BMI Calculation
Dual-strategy approach:
1. Primary: RapidAPI BMI Calculator
2. Fallback: Local calculation (weight / height²)

## Architecture Patterns

- **Layered Architecture**: Controller → Service → Repository
- **Dependency Injection**: All dependencies injected via constructors
- **Repository Pattern**: Abstract data access layer
- **Black-box Testing**: Tests in `_test` packages using only public APIs
- **Mock Generation**: Using go.uber.org/mock for unit testing

## Testing Guidelines

- Write tests for all new features
- Maintain minimum 25% code coverage
- Use black-box testing approach (`_test` package)
- Generate mocks for external dependencies
- Run tests before submitting PR

## License

This project is created for Hacktiv8 FTGO Phase 2 final project.

## Authors

- Deden Ruslan (@Dedenruslan19)

## Acknowledgments

- Echo Framework - High performance Go web framework
- GORM - Go ORM library
- Google Gemini AI - Workout generation
- RapidAPI - BMI calculation service
- Hacktiv8 - FTGO Phase 2 curriculum

## Support

For questions or issues, please open an issue on GitHub or contact the development team.

---

- API DOC : https://documenter.getpostman.com/view/45800381/2sB3WsPKeW
- DEPLOYMENT URL : https://med-project-production.up.railway.app/
