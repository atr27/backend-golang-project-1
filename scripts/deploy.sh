#!/bin/bash

# Hospital EMR System - Deployment Script
# Automates deployment to various environments

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Configuration
ENVIRONMENT=${1:-staging}
VERSION=${2:-latest}
IMAGE_NAME="hospital-emr-backend"
REGISTRY=${DOCKER_REGISTRY:-""}

echo "=========================================="
echo "Hospital EMR System - Deployment"
echo "=========================================="
echo ""
echo "Environment: ${ENVIRONMENT}"
echo "Version: ${VERSION}"
echo ""

# Validate environment
if [[ ! "$ENVIRONMENT" =~ ^(development|staging|production)$ ]]; then
    echo -e "${RED}Error: Invalid environment '${ENVIRONMENT}'${NC}"
    echo "Valid environments: development, staging, production"
    exit 1
fi

# Confirm production deployment
if [ "$ENVIRONMENT" == "production" ]; then
    echo -e "${YELLOW}WARNING: You are about to deploy to PRODUCTION${NC}"
    read -p "Are you sure you want to continue? (yes/no): " confirm
    if [ "$confirm" != "yes" ]; then
        echo "Deployment cancelled"
        exit 0
    fi
fi

# Run tests before deployment
echo -e "${BLUE}Running tests...${NC}"
./scripts/test.sh || {
    echo -e "${RED}Tests failed! Deployment aborted.${NC}"
    exit 1
}

# Build Docker image
echo -e "${BLUE}Building Docker image...${NC}"
docker build -t ${IMAGE_NAME}:${VERSION} .
echo -e "${GREEN}✓ Docker image built${NC}"

# Tag image
if [ -n "$REGISTRY" ]; then
    echo -e "${BLUE}Tagging image for registry...${NC}"
    docker tag ${IMAGE_NAME}:${VERSION} ${REGISTRY}/${IMAGE_NAME}:${VERSION}
    docker tag ${IMAGE_NAME}:${VERSION} ${REGISTRY}/${IMAGE_NAME}:${ENVIRONMENT}
    echo -e "${GREEN}✓ Image tagged${NC}"
    
    # Push to registry
    echo -e "${BLUE}Pushing image to registry...${NC}"
    docker push ${REGISTRY}/${IMAGE_NAME}:${VERSION}
    docker push ${REGISTRY}/${IMAGE_NAME}:${ENVIRONMENT}
    echo -e "${GREEN}✓ Image pushed to registry${NC}"
fi

# Deploy based on environment
case $ENVIRONMENT in
    development)
        echo -e "${BLUE}Deploying to development environment...${NC}"
        docker-compose -f docker-compose.yml up -d
        ;;
    
    staging)
        echo -e "${BLUE}Deploying to staging environment...${NC}"
        kubectl set image deployment/emr-api \
            emr-api=${REGISTRY}/${IMAGE_NAME}:${VERSION} \
            -n emr-staging
        kubectl rollout status deployment/emr-api -n emr-staging
        ;;
    
    production)
        echo -e "${BLUE}Deploying to production environment...${NC}"
        
        # Create backup before deployment
        echo -e "${BLUE}Creating backup...${NC}"
        ./scripts/backup.sh
        
        # Deploy with blue-green strategy
        kubectl set image deployment/emr-api \
            emr-api=${REGISTRY}/${IMAGE_NAME}:${VERSION} \
            -n emr-production
        
        # Wait for rollout
        kubectl rollout status deployment/emr-api -n emr-production
        
        # Run smoke tests
        echo -e "${BLUE}Running smoke tests...${NC}"
        sleep 10
        curl -f http://production-api.hospital-emr.com/health || {
            echo -e "${RED}Smoke tests failed! Rolling back...${NC}"
            kubectl rollout undo deployment/emr-api -n emr-production
            exit 1
        }
        ;;
esac

echo ""
echo -e "${GREEN}=========================================="
echo "Deployment Successful!"
echo "==========================================${NC}"
echo ""
echo "Environment: ${ENVIRONMENT}"
echo "Version: ${VERSION}"
echo "Time: $(date)"
echo ""

# Post-deployment checks
echo "Post-deployment checks:"
echo "1. Monitor logs: kubectl logs -f deployment/emr-api -n emr-${ENVIRONMENT}"
echo "2. Check metrics: http://monitoring.hospital-emr.com"
echo "3. Run manual tests: ./scripts/smoke-test.sh ${ENVIRONMENT}"
echo ""
