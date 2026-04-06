#!/bin/bash
# Acest script asigură declanșarea a 3 alerte distincte, din 3 Pod-uri diferite,
# dar toate având exact aceeași imagine (nginx:1.25.0).
# Acest volum de "păcate" pentru o singură imagine va activa logica de "Repeat Offender".

NAMESPACE="default"
IMAGE="nginx:1.25.0"

for i in {1..3}
do
  POD_NAME="repeat-offender-$i"
  echo ">>> Creare pod: $POD_NAME (Imagine: $IMAGE)"
  
  kubectl run $POD_NAME --image=$IMAGE --restart=Never --namespace=$NAMESPACE
  
  echo ">>> Așteptăm ca podul să pornească..."
  kubectl wait --for=condition=Ready pod/$POD_NAME --namespace=$NAMESPACE --timeout=60s
  
  echo ">>> Rulăm o comandă suspectă (shell în container) pentru a genera o alertă Falco..."
  kubectl exec $POD_NAME --namespace=$NAMESPACE -- /bin/bash -c "apt-get update"
  
  echo "---"
done

echo ">>> Alertele au fost generate. Analizează UI-ul și cardul imaginii $IMAGE"
