# User-Guide
Saúl Eduardo Zepeda de la Torre | 0214016

Karla Sofía González Rodríguez | 0214774

## Installation

First of all, it is important to run the next commands to install certain features and upgrade the full system

`$ export GO111MODULE=on`

> (Be able to run the commands)

`$ sudo pacman -Syu --noconfirm`

`$ sudo pacman -Syu --noconfirm protobuf`

`$ go version`  

> (Check that the version is go1.16.4)

`$ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26`

`$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1`

`$ go env GOPATH ` 

> (Check thath the path is /home/cs-user/go)

`$ ls /home/cs-user/go/bin`

`$ export PATH="$PATH:$(go env GOPATH)/bin"`

`$ export PATH="$PATH:$(go env GOPATH)/bin"`

For the next installation, it is needed to be in the next path

`$ cd /home/cs-user/go/src/github.com/<GITHUB USERNAME>/dc-final`

Now we can do the installation, with the following command

`$ protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/helloworld.proto`

## Run the code

To run the code, it is important to be in the next path 

`$ cd /home/cs-user/go/src/github.com/<GITHUB USERNAME>/dc-final`

To start running any part of the code please run the next command

`$ export GO111MODULE=off`

#### Controller, Scheduler, API

Once you are here you can run the next commad to start the Controller, Scheduler and API

`$ go run main.go`

#### API
If you want to start using the API, you need to open another terminal. It also needs to be in the same location of /home/cs-user/go/src/github.com/<GITHUB USERNAME>/dc-final, so run again in the new terminal the command:

`$ cd /home/cs-user/go/src/github.com/<GITHUB USERNAME>/dc-final`

Remember to run the next command in the new terminal

`$ export GO111MODULE=off`

##### >Login

There is a default account, which is username, so it is time to login. Run the next command:

`curl -u username:password localhost:8080/login`


>If the user and the password are correct, the following message will appear and it will generate a secret token that you will need to remember in order to do certain things in the program: { "message": "Hi username, welcome to the DPIP System", "token" <-TOKEN-> }

  

>If the user is not correct or it is not registered, it will display: { "message": "The username is not registered" }

  

>Or if the password is incorrect, this message will appear: { "message": "The password is incorrect" }

  

##### >Status

If you logged in but you want to doble check that you logged in correctly, you can see that with the next line (remember to have your token) :

`curl -H "Authorization: Bearer <TOKEN>" localhost:8080/status`

  

>If the token exists and there are no mistakes, the following message will be shown: { "message": "Hi username, the DPIP System is Up and Running" "time": <-DATE_TIME-> }

  

>If the token is incorrect and there is no user with that token, then it will display this message { "message": "Error, you have to login." }

##### >Create new workload

You can also create some workloads with the next command(remember to have your token and also a new token will be generated, please save the workload id):

`curl -X POST -H "Authorization: Bearer <ACCESS_TOKEN>" http://localhost:8080/workloads`

##### >Details of workload

You can also check some details of the workloads that were created previously with the next command(remember to have the inicial token and the workload id):

`curl -H "Authorization: Bearer <ACCESS_TOKEN>" http://localhost:8080/workloads/<workload_id>`

##### >Images

Once you logged in and there are no problems with your status, then you can upload an image using the following command (remember to have your token and the workload id):

`$ curl -F 'data=@test.jpg' -F 'workload-id=0001' -F 'type=original' -H "Authorization: Bearer Nf4KsSiA" localhost:8080/images`


  

##### >Logout

If you want to logout you can do that, however your token that was generated in the login is going to be needed. Run the command:

`curl -H "Authorization: Bearer <TOKEN>" localhost:8080/logout`

  

>With this action the information of the user will be deleted, and it will show the next line: { "message": "Bye username, your token has been revoked" }

#### Workers

If you want to start using the Workers, you need to open another terminal. It also needs to be in the same location of /home/cs-user/go/src/github.com/<GITHUB USERNAME>/dc-final, so run again in the new terminal the command:

`$ cd /home/cs-user/go/src/github.com/<GITHUB USERNAME>/dc-final`

Remember to run the next command in the new terminal

`$ export GO111MODULE=off`

Once you are done with that you can create a worker with the next command, the worker will be a self-running component that will do the real work of filtering images.

`$ go run main.go --controller <host>:<port> --worker-name <worker_name>`

>You can run as many workers you want, but remember to open a new terminal and follow the previous steps of Workers
