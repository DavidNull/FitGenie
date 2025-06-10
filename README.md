# FitGenie - AI Fashion Stylist

FitGenie is an AI-powered fashion styling application that helps users create perfect outfits based on color theory, style analysis, and personal preferences.

## Features

- **Color Theory Analysis**: Advanced color harmony analysis and seasonal color recommendations
- **Style Analysis**: Intelligent style categorization and coherence scoring
- **AI Outfit Recommendations**: Smart outfit combinations based on user preferences
- **Personal Style Profiles**: Customizable style and color profiles
- **Wardrobe Management**: Digital wardrobe with clothing item analysis

## Quick Start (Local Development)

### Prerequisites

- Go 1.23 or higher
- Docker (for PostgreSQL database)
- Git

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd FitGenie
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Start the application**
   ```bash
   ./run-local.sh
   ```

   This script will:
   - Start PostgreSQL in a Docker container
   - Wait for the database to be ready
   - Run the Go application

4. **Access the API**
   - API will be available at: `http://localhost:8080`
   - Health check: `http://localhost:8080/api/v1/health`

### Stopping the Application

```bash
./stop-local.sh
```

## API Endpoints

### Health Check
- `GET /api/v1/health` - Service health status

### Users
- `POST /api/v1/users` - Create a new user
- `GET /api/v1/users` - List all users
- `GET /api/v1/users/:id` - Get user by ID
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

### User Profiles
- `POST /api/v1/users/:id/style-profile` - Create/update style profile
- `GET /api/v1/users/:id/style-profile` - Get style profile
- `POST /api/v1/users/:id/color-profile` - Create/update color profile
- `GET /api/v1/users/:id/color-profile` - Get color profile

### Clothing Items
- `POST /api/v1/users/:userId/clothing` - Add clothing item
- `GET /api/v1/users/:userId/clothing` - List clothing items
- `GET /api/v1/users/:userId/clothing/:id` - Get clothing item
- `PUT /api/v1/users/:userId/clothing/:id` - Update clothing item
- `DELETE /api/v1/users/:userId/clothing/:id` - Delete clothing item

### Outfits
- `POST /api/v1/users/:userId/outfits` - Create outfit
- `GET /api/v1/users/:userId/outfits` - List outfits
- `GET /api/v1/users/:userId/outfits/:id` - Get outfit
- `DELETE /api/v1/users/:userId/outfits/:id` - Delete outfit

### AI Recommendations
- `POST /api/v1/users/:userId/recommendations/outfits` - Get outfit recommendations
- `GET /api/v1/users/:userId/recommendations/personalized` - Get personalized recommendations

## Configuration

The application uses environment variables for configuration. Copy `.env.example` to `.env` and modify as needed:

```bash
cp .env.example .env
```

Key configuration options:
- `DATABASE_URL`: PostgreSQL connection string
- `PORT`: Server port (default: 8080)
- `GIN_MODE`: Gin framework mode (debug/release)

## Development

### Project Structure

```
FitGenie/
├── cmd/                    # Application entry points
├── internal/
│   ├── api/               # API routes and handlers
│   ├── config/            # Configuration management
│   ├── database/          # Database connection and migrations
│   ├── models/            # Data models
│   └── services/          # Business logic services
├── docker-compose.yml     # Docker services (optional)
├── Dockerfile            # Container definition
├── go.mod               # Go module definition
├── main.go              # Application entry point
├── run-local.sh         # Local development script
└── stop-local.sh        # Stop local environment script
```

### Services

1. **Color Theory Service**: Handles color analysis, seasonal color determination, and color harmony calculations
2. **Style Service**: Manages style categorization, body type recommendations, and style coherence analysis
3. **AI Service**: Combines color and style analysis for intelligent outfit recommendations

### Database

The application uses PostgreSQL with GORM for database operations. Migrations run automatically on startup.

## Testing the API

You can test the API using curl or any HTTP client:

```bash
# Health check
curl http://localhost:8080/api/v1/health

# Create a user
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com"}'

# Get available style categories
curl http://localhost:8080/api/v1/style-analysis/categories

# Get color seasons
curl http://localhost:8080/api/v1/color-theory/seasons
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## Support

For support and questions:
- Create an issue in the repository
- Contact the development team
- Check the documentation wiki

## Roadmap

### Upcoming Features
- [ ] Machine learning model training pipeline
- [ ] Advanced image recognition for clothing items
- [ ] Social features and outfit sharing
- [ ] Mobile app integration
- [ ] Weather-based recommendations
- [ ] Shopping integration and price tracking
- [ ] Virtual wardrobe visualization
- [ ] Style trend analysis

### Performance Improvements
- [ ] Redis caching implementation
- [ ] Database query optimization
- [ ] CDN integration for images
- [ ] API rate limiting
- [ ] Monitoring and alerting

---

**FitGenie** - Making fashion choices smarter with AI 🤖👗✨
