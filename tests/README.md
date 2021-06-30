## If you want to mock the component(s) with which your component is communcating then 
* create .proto files in side resources folder: Ref --> [here](../.grpc-scripts/protos)
* provide mock data. Ref --> [here](../.grpc-scripts/mock-data)

## genereate the client side code using: 
Go-to your test project
``` python3 -m grpc_tools.protoc -I./resources/protos --python_out=./test --grpc_python_out=./test ./resources/protos/* ```


## Way-1 If you want to test in local machine without k3s setup
* If your component needs to communicate with mock server, 
    * start the mock server using command : ``` docker run -p 4770:4770 -p 4771:4771 -v $(pwd)/resources/protos/:/proto tkpd/gripmock /proto/Comp-b.proto /proto/Comp-c.proto```
    * Go to your application you want to test and change the ip:port to localhost:4770 of components its communicating.
<!-- * Go to your application you want to test , start docker container and expose port so that test case can communicate with server -->
* Go to your application you want to test , start the service
* Chnage the test cases as per your need and execute your test case : e.g., python3 test.py

## Way-2 If you want to test on local k3s cluster
* Setup k3d if not done: https://k3d.io/
* Go-to your test project, Create cluster: sh resources/cluster-setup.sh
* Go-to your application you want to test deploy your application either using deployment.yaml or helm chart
* Go-to your test project
* execute ``` sh resources/test-component.sh ```

