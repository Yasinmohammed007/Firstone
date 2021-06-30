# You dont need to execute below comand and either build the image or expect it available there.
docker build -t app-test:latest .

k3d image import -c paloma app-test:latest

kubectl apply -f resources/test-deployment.yaml             
