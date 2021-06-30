# execute below steps in component source code.

docker build -t $1 .

k3d image import -c paloma $1

kubectl apply -f $2
