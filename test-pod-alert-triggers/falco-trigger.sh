#!/bin/bash

# Configuration
TEST_NS="falco-test-ns"
PROD_NS="prod"
POD_NAME="falco-trigger-bot"

echo "=========================================================="
echo "🛡️  FALCO & TALON AUTOMATED THREAT GENERATION SCRIPT       🛡"
echo "=========================================================="
echo "Creating dummy namespaces..."
kubectl create namespace $TEST_NS --dry-run=client -o yaml | kubectl apply -f -
kubectl create namespace $PROD_NS --dry-run=client -o yaml | kubectl apply -f -

echo "=========================================================="
echo "🚀 1. Spawning dummy testing pod in '$TEST_NS' namespace..."
kubectl run $POD_NAME -n $TEST_NS --image=ubuntu --restart=Never -- sleep 3600
kubectl wait --for=condition=ready pod/$POD_NAME -n $TEST_NS --timeout=30s
echo ""

echo "=========================================================="
echo "🚀 2. Spawning dummy testing pod in '$PROD_NS' namespace..."
kubectl run $POD_NAME -n $PROD_NS --image=ubuntu --restart=Never -- sleep 3600
kubectl wait --for=condition=ready pod/$POD_NAME -n $PROD_NS --timeout=30s
echo ""

echo "=========================================================="
echo "🚨 TRIGGER 1: Terminal shell in container (Exec Shell Bash in TEST)"
echo "Action: Executing /bin/bash..."
kubectl exec -n $TEST_NS $POD_NAME -- /bin/bash -c "echo 'Hello from malicious bash!'"
sleep 2

echo "=========================================================="
echo "🚨 TRIGGER 2: Read Sensitive File (Read Shadow File in PROD)"
echo "Action: Accessing /etc/shadow..."
kubectl exec -n $PROD_NS $POD_NAME -- cat /etc/shadow > /dev/null 2>&1
sleep 2

echo "=========================================================="
echo "🚨 TRIGGER 3: Network Tool Launched (Curl or Wget in PROD)"
echo "Action: Running wget..."
kubectl exec -n $PROD_NS $POD_NAME -- apt-get update -y > /dev/null 2>&1
kubectl exec -n $PROD_NS $POD_NAME -- apt-get install -y wget > /dev/null 2>&1
kubectl exec -n $PROD_NS $POD_NAME -- wget -qO- http://google.com > /dev/null 2>&1
sleep 2

echo "=========================================================="
echo "🚨 TRIGGER 4: Network Access from Shell (Reverse Shell Connect in TEST)"
echo "Action: Executing reverse shell connection via /dev/tcp output..."
# Falco detects `sh` connecting to random IPs on the raw socket descriptors
kubectl exec -n $TEST_NS $POD_NAME -- bash -c "exec 3<>/dev/tcp/8.8.8.8/80; echo -e 'GET / HTTP/1.0\r\n' >&3; cat <&3 | head -n 1"
sleep 2

echo "=========================================================="
echo "🚨 TRIGGER 5: Any Execution in Prod Namespace (Exec in PROD)"
echo "Action: Executing benign command in strictly monitored PROD environment..."
kubectl exec -n $PROD_NS $POD_NAME -- ls -la /tmp
sleep 2

echo "=========================================================="
echo "✅ All tests fired! Go inspect your FalcoSight Dashboard at http://falcosight.ghosty-dashboard.com/"
echo "NOTE: Some Talon (SOAR) actions like 'isolate-namespace' might have altered your environment (NetworkPolicies)."
echo ""
echo "Cleanup: Do you want to delete the testing namespaces? (y/n)"
read cleanup
if [ "$cleanup" = "y" ]; then
    kubectl delete namespace $TEST_NS
    kubectl delete namespace $PROD_NS
    echo "Cleanup complete."
fi
