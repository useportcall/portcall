#!/usr/bin/env bash

# Portcall Kubernetes Health Check Script
# Tests all services and infrastructure components

set -e

NAMESPACE="portcall-dev"
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "🔍 Portcall Kubernetes Health Check"
echo "===================================="
echo ""

# Function to check if namespace exists
check_namespace() {
    echo "📦 Checking namespace..."
    if kubectl get namespace $NAMESPACE &> /dev/null; then
        echo -e "${GREEN}✓${NC} Namespace '$NAMESPACE' exists"
        return 0
    else
        echo -e "${RED}✗${NC} Namespace '$NAMESPACE' not found"
        echo "   Run: helm install portcall ./k8s/portcall-chart -n $NAMESPACE --create-namespace"
        return 1
    fi
}

# Function to check pod status
check_pods() {
    echo ""
    echo "🚀 Checking pod status..."
    
    # Get all pods
    PODS=$(kubectl get pods -n $NAMESPACE --no-headers 2>/dev/null)
    
    if [ -z "$PODS" ]; then
        echo -e "${RED}✗${NC} No pods found in namespace $NAMESPACE"
        return 1
    fi
    
    # Count pods by status
    TOTAL=$(echo "$PODS" | wc -l | tr -d ' ')
    RUNNING=$(echo "$PODS" | grep "Running" | wc -l | tr -d ' ')
    COMPLETED=$(echo "$PODS" | grep "Completed" | wc -l | tr -d ' ')
    PENDING=$(echo "$PODS" | grep "Pending" | wc -l | tr -d ' ')
    ERROR=$(echo "$PODS" | grep -E "Error|CrashLoopBackOff|ImagePullBackOff" | wc -l | tr -d ' ')
    
    echo "   Total pods: $TOTAL"
    echo -e "   ${GREEN}Running:${NC} $RUNNING"
    echo -e "   ${GREEN}Completed:${NC} $COMPLETED (jobs)"
    
    if [ "$PENDING" -gt 0 ]; then
        echo -e "   ${YELLOW}Pending:${NC} $PENDING"
    fi
    
    if [ "$ERROR" -gt 0 ]; then
        echo -e "   ${RED}Errors:${NC} $ERROR"
        echo ""
        echo "   Pods with errors:"
        echo "$PODS" | grep -E "Error|CrashLoopBackOff|ImagePullBackOff" | awk '{print "   - " $1 " (" $3 ")"}'
        return 1
    fi
    
    echo ""
    kubectl get pods -n $NAMESPACE
    return 0
}

# Function to check stateful sets
check_statefulsets() {
    echo ""
    echo "💾 Checking StatefulSets..."
    
    for sset in postgres redis minio; do
        if kubectl get statefulset $sset -n $NAMESPACE &> /dev/null; then
            READY=$(kubectl get statefulset $sset -n $NAMESPACE -o jsonpath='{.status.readyReplicas}')
            DESIRED=$(kubectl get statefulset $sset -n $NAMESPACE -o jsonpath='{.status.replicas}')
            
            if [ "$READY" == "$DESIRED" ] && [ "$READY" != "" ]; then
                echo -e "${GREEN}✓${NC} $sset is ready ($READY/$DESIRED)"
            else
                echo -e "${RED}✗${NC} $sset is not ready ($READY/$DESIRED)"
            fi
        else
            echo -e "${YELLOW}⊘${NC} $sset not deployed"
        fi
    done
}

# Function to check services
check_services() {
    echo ""
    echo "🌐 Checking Services..."
    
    SERVICES=$(kubectl get svc -n $NAMESPACE --no-headers 2>/dev/null)
    
    if [ -z "$SERVICES" ]; then
        echo -e "${RED}✗${NC} No services found"
        return 1
    fi
    
    SERVICE_COUNT=$(echo "$SERVICES" | wc -l | tr -d ' ')
    echo -e "${GREEN}✓${NC} Found $SERVICE_COUNT services"
    echo ""
    kubectl get svc -n $NAMESPACE
}

# Function to test database connectivity
test_postgres() {
    echo ""
    echo "🐘 Testing PostgreSQL..."
    
    if kubectl get pod postgres-0 -n $NAMESPACE &> /dev/null; then
        if kubectl exec -n $NAMESPACE postgres-0 -- pg_isready -U admin &> /dev/null; then
            echo -e "${GREEN}✓${NC} PostgreSQL is accepting connections"
            
            # Try to connect and list databases
            DBS=$(kubectl exec -n $NAMESPACE postgres-0 -- psql -U admin -d main_portcall_db -tAc "SELECT 1" 2>&1)
            if [ "$?" -eq 0 ]; then
                echo -e "${GREEN}✓${NC} Can query main database"
            else
                echo -e "${YELLOW}⚠${NC} Database might still be initializing"
            fi
        else
            echo -e "${RED}✗${NC} PostgreSQL not ready"
            return 1
        fi
    else
        echo -e "${RED}✗${NC} PostgreSQL pod not found"
        return 1
    fi
}

# Function to test Redis
test_redis() {
    echo ""
    echo "📮 Testing Redis..."
    
    if kubectl get pod redis-0 -n $NAMESPACE &> /dev/null; then
        if kubectl exec -n $NAMESPACE redis-0 -- redis-cli ping 2>&1 | grep -q "PONG"; then
            echo -e "${GREEN}✓${NC} Redis is responding"
        else
            echo -e "${RED}✗${NC} Redis not responding"
            return 1
        fi
    else
        echo -e "${RED}✗${NC} Redis pod not found"
        return 1
    fi
}

# Function to test MinIO
test_minio() {
    echo ""
    echo "📦 Testing MinIO..."
    
    if kubectl get pod minio-0 -n $NAMESPACE &> /dev/null; then
        # Check if MinIO is responding
        MINIO_STATUS=$(kubectl exec -n $NAMESPACE minio-0 -- sh -c "curl -s -o /dev/null -w '%{http_code}' http://localhost:9000/minio/health/live" 2>/dev/null)
        
        if [ "$MINIO_STATUS" == "200" ]; then
            echo -e "${GREEN}✓${NC} MinIO is healthy"
        else
            echo -e "${RED}✗${NC} MinIO health check failed (HTTP $MINIO_STATUS)"
            return 1
        fi
    else
        echo -e "${RED}✗${NC} MinIO pod not found"
        return 1
    fi
}

# Function to test API endpoints
test_api_endpoints() {
    echo ""
    echo "🔌 Testing API Endpoints..."
    
    # Wait a moment for port-forwarding if needed
    sleep 2
    
    # Test API
    if kubectl get svc api -n $NAMESPACE &> /dev/null; then
        NODE_PORT=$(kubectl get svc api -n $NAMESPACE -o jsonpath='{.spec.ports[0].nodePort}')
        echo "   API: http://localhost:$NODE_PORT"
        
        # Try to access health endpoint
        HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:$NODE_PORT/health 2>/dev/null || echo "000")
        if [ "$HTTP_CODE" == "200" ]; then
            echo -e "   ${GREEN}✓${NC} API health check passed"
        elif [ "$HTTP_CODE" == "404" ]; then
            echo -e "   ${YELLOW}⚠${NC} API responding but /health endpoint not found (HTTP 404)"
        elif [ "$HTTP_CODE" == "000" ]; then
            echo -e "   ${YELLOW}⚠${NC} API not accessible (connection refused - pod may still be starting)"
        else
            echo -e "   ${YELLOW}⚠${NC} API returned HTTP $HTTP_CODE"
        fi
    fi
    
    # Test Dashboard
    if kubectl get svc dashboard -n $NAMESPACE &> /dev/null; then
        NODE_PORT=$(kubectl get svc dashboard -n $NAMESPACE -o jsonpath='{.spec.ports[0].nodePort}')
        echo "   Dashboard: http://localhost:$NODE_PORT"
    fi
    
    # Test Keycloak
    if kubectl get svc keycloak -n $NAMESPACE &> /dev/null; then
        NODE_PORT=$(kubectl get svc keycloak -n $NAMESPACE -o jsonpath='{.spec.ports[0].nodePort}')
        echo "   Keycloak: http://localhost:$NODE_PORT"
    fi
}

# Function to show access URLs
show_access_urls() {
    echo ""
    echo "🌍 Service Access URLs"
    echo "======================"
    
    # List of services to check (service_name display_name)
    services=(
        "dashboard:Dashboard"
        "api:API"
        "admin:Admin API"
        "checkout:Checkout"
        "file-api:File API"
        "quote:Quote Service"
        "keycloak:Keycloak"
        "minio-console:MinIO Console"
        "mailpit:Mailpit (Email Testing)"
        "redis-insight:Redis Insight"
    )
    
    for entry in "${services[@]}"; do
        svc_name="${entry%%:*}"
        display_name="${entry#*:}"
        
        if kubectl get svc $svc_name -n $NAMESPACE &> /dev/null; then
            NODE_PORT=$(kubectl get svc $svc_name -n $NAMESPACE -o jsonpath='{.spec.ports[?(@.nodePort)].nodePort}')
            if [ ! -z "$NODE_PORT" ]; then
                printf "%-25s http://localhost:%s\n" "$display_name:" "$NODE_PORT"
            fi
        fi
    done
}

# Function to show helpful commands
show_helpful_commands() {
    echo ""
    echo "📝 Helpful Commands"
    echo "==================="
    echo "View all pods:              kubectl get pods -n $NAMESPACE"
    echo "Watch pod status:           kubectl get pods -n $NAMESPACE -w"
    echo "View pod logs:              kubectl logs -n $NAMESPACE <pod-name> -f"
    echo "Describe pod:               kubectl describe pod -n $NAMESPACE <pod-name>"
    echo "Check Helm release:         helm status portcall -n $NAMESPACE"
    echo "Port-forward a service:     kubectl port-forward -n $NAMESPACE svc/dashboard 8082:8082"
    echo "Restart a deployment:       kubectl rollout restart -n $NAMESPACE deployment/api"
    echo "Delete everything:          helm uninstall portcall -n $NAMESPACE"
}

# Main execution
main() {
    if ! check_namespace; then
        exit 1
    fi
    
    check_pods
    POD_STATUS=$?
    
    check_statefulsets
    check_services
    
    if [ $POD_STATUS -eq 0 ]; then
        test_postgres
        test_redis
        test_minio
        test_api_endpoints
    else
        echo ""
        echo -e "${YELLOW}⚠${NC} Skipping service tests due to pod errors"
        echo ""
        echo "Check pod logs for errors:"
        kubectl get pods -n $NAMESPACE --no-headers | grep -v "Running\|Completed" | awk '{print $1}' | while read pod; do
            echo "  kubectl logs -n $NAMESPACE $pod"
        done
    fi
    
    show_access_urls
    show_helpful_commands
    
    echo ""
    if [ $POD_STATUS -eq 0 ]; then
        echo -e "${GREEN}✅ Portcall cluster is healthy!${NC}"
    else
        echo -e "${YELLOW}⚠️  Some issues detected - check logs above${NC}"
    fi
}

# Run main function
main
