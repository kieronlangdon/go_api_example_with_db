### Local postgres db setup
```
docker pull postgres   
docker pull postgres:[tag_you_want]   
mkdir -p $HOME/docker/volumes/postgres   
docker run --rm   --name pg-docker -e POSTGRES_USER=docker,POSTGRES_PASSWORD=docker -d -p 5432:5432 -v $HOME/docker/volumes/postgres:/var/lib/postgresql/data  postgres   
```
Enter running container (`docker ps -a` will show running containers)   
```
docker exec -it 708cb73e233c bash   
su root   
su postgres   
psql -h localhost   
ALTER USER postgres with password 'docker';   
```
Exit out of container   

### Alt steps for running minikube and to avoid going down a rabbit hole of nonsense   
**Build image**   
`docker build -t my-postgres-app .`   
**Tag image**   
`docker tag my-go-app:latest localhost:5000/my-postgres-app`   
**Point your terminal to use the docker daemon inside minikube**   
`eval $(minikube docker-env)`   
**Check to see all images**      
`docker images`   
Notice that the image we just created and tagged is not here? We're using minikube's docker daemon now   
**Push to cache in minikube**   
`minikube cache add localhost:5000/my-postgres-app`   
This will fail as the image doesn't exist in this daemon   
So in another terminal OR after undoing `eval $(minikube docker-env)`   
Re-run `minikube cache add localhost:5000/my-postgres-app`   
**Check image is in minikube cache**   
`minikube cache list`   
We will see an image with no repo & tag, check the imageId, it's the same as the local image in non minikube docker   
**Tag image in minikube docker daemon**   
`docker tag c22dbba37091 localhost:5000/my-postgres-app`   
**Create deployment in minikube**   
`minikube kubectl create deployment testdev -- --image=localhost:5000/my-postgres-app`  
### Alter deployment yaml to pull from minikube local   
Get deployment yaml   
`minikube kubectl get deploy testdev -- -o yaml --export >> testdev.yaml`   
In testdev.yaml   
`imagePullPolicy: Always`   
change to   
`imagePullPolicy: Never`   
Apply yaml   
`minikube kubectl apply -- -f testdev.yaml`   
### Expose container port using service   
Set deployment to expose deployment of type node port   
`minikube kubectl expose deployment testdev -- --type=NodePort --port=5432`   
Get port that has been exposed externally   
` minikube kubectl get svc testdev`   
Get minikube ip   
`minikube ip`    
