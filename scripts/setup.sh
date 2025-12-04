#!/bin/bash

# Hospital EMR System - Setup Script
# This script sets up the development environment

set -e

echo "=========================================="
echo "Hospital EMR System - Setup"
echo "=========================================="
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed${NC}"
    echo "Please install Go 1.21 or higher from https://golang.org/dl/"
    exit 1
fi

echo -e "${GREEN}✓ Go is installed: $(go version)${NC}"

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo -e "${YELLOW}Warning: Docker is not installed${NC}"
    echo "Docker is recommended for running PostgreSQL, Redis, and NATS"
    echo "Install from https://www.docker.com/get-started"
else
    echo -e "${GREEN}✓ Docker is installed: $(docker --version)${NC}"
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo -e "${YELLOW}Warning: Docker Compose is not installed${NC}"
else
    echo -e "${GREEN}✓ Docker Compose is installed: $(docker-compose --version)${NC}"
fi

echo ""
echo "Setting up the project..."
echo ""

# Install Go dependencies
echo "Installing Go dependencies..."
go mod download
go mod tidy
echo -e "${GREEN}✓ Dependencies installed${NC}"

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo "Creating .env file..."
    cp .env.example .env
    
    # Generate random JWT secret
    JWT_SECRET=$(openssl rand -base64 32)
    sed -i.bak "s|JWT_SECRET=your_jwt_secret_key|JWT_SECRET=${JWT_SECRET}|g" .env
    
    # Generate random encryption key
    ENCRYPTION_KEY=$(openssl rand -base64 32)
    sed -i.bak "s|ENCRYPTION_KEY=your_32_byte_encryption_key_here|ENCRYPTION_KEY=${ENCRYPTION_KEY}|g" .env
    
    rm .env.bak
    echo -e "${GREEN}✓ .env file created with secure keys${NC}"
else
    echo -e "${YELLOW}⚠ .env file already exists, skipping...${NC}"
fi

# Create necessary directories
echo "Creating directories..."
mkdir -p uploads
mkdir -p logs
mkdir -p tmp
echo -e "${GREEN}✓ Directories created${NC}"

# Install development tools
echo ""
echo "Installing development tools..."
echo "This may take a few minutes..."

# Air (hot reload)
if ! command -v air &> /dev/null; then
    go install github.com/cosmtrek/air@latest
    echo -e "${GREEN}✓ Air installed${NC}"
else
    echo -e "${GREEN}✓ Air already installed${NC}"
fi

# Swag (Swagger documentation)
if ! command -v swag &> /dev/null; then
    go install github.com/swaggo/swag/cmd/swag@latest
    echo -e "${GREEN}✓ Swag installed${NC}"
else
    echo -e "${GREEN}✓ Swag already installed${NC}"
fi

# golangci-lint (linter)
if ! command -v golangci-lint &> /dev/null; then
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
    echo -e "${GREEN}✓ golangci-lint installed${NC}"
else
    echo -e "${GREEN}✓ golangci-lint already installed${NC}"
fi

echo ""
echo "=========================================="
echo "Setup Complete!"
echo "=========================================="
echo ""
echo "Next steps:"
echo ""
echo "1. Start Docker services:"
echo "   ${GREEN}docker-compose up -d${NC}"
echo ""
echo "2. Run database migrations:"
echo "   ${GREEN}make migrate-up${NC}"
echo ""
echo "3. Seed initial data:"
echo "   ${GREEN}make seed${NC}"
echo ""
echo "4. Start the application:"
echo "   ${GREEN}make run${NC}"
echo ""
echo "5. Access the API at:"
echo "   ${GREEN}http://localhost:8080${NC}"
echo ""
echo "Default credentials:"
echo "   Admin: admin@hospital-emr.com / admin123"
echo "   Doctor: doctor@hospital-emr.com / doctor123"
echo ""
echo -e "${YELLOW}⚠ Remember to change default passwords!${NC}"
echo ""
