baseDir=/etc/simple-webhook/pki
mkdir -p $baseDir
#Generate the webhook admission ca key
openssl genrsa -out $baseDir/ca.key 2048

# Generate csr and sing it
openssl req -new -key $baseDir/ca.key -subj "/CN=SIMPLE-WEBHOOK-CA" -out $baseDir/ca.csr 
openssl x509 -req -in $baseDir/ca.csr -singkey $baseDir/ca.key -out $baseDir/ca.crt

#generate server cert
openssl genrsa -out $baseDir/server.key 2048
openssl req -new -key $baseDir/server.key -subj "/CN=simple-webhook-svc.simple-webhook.svc" -out $baseDir/server.csr 
openssl x509 -req -in $baseDir/server.csr -CAkey $baseDir/ca.key -CA $baseDir/ca.crt -out $baseDir/server.crt

## Check if the secret already exist and delete it 
secretCount=$(kubectl get secret webhook-tls -n simple-webhook --no-header | wc -l)
if (( secretCount > 0 ))
then
    kubectl delete -n simple-webhook secret webhook-tls
fi

kubectl create secret tls webhook-tls --cert=$basedir/server.crt --key=$basedir/server.key -n simple-webhook
caBase64=$(cat $basedir/ca.crt | base64 | tr -d "\n")
sed -i "s/.*caBundle:.*/      caBundle: $caBase64/" mutate-webhook.yaml
sed -i "s/.*caBundle:.*/      caBundle: $caBase64/" validate-webhook.yaml

cd deployment
kubectl apply -f ./