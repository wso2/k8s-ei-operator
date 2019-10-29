# Developer Guide for the EI Operator

Use the following steps to setup, deploy and run the integration solutions with the EI operator. 

## Prerequisites
The k8s-ei-operator is built with operator-sdk v0.7.0 and supported in the following environment.

-   [Kubernetes](https://kubernetes.io/docs/setup/) cluster and client v1.11+   
-   [Docker](https://docs.docker.com/)
-	[Operator SDK CLI](https://github.com/operator-framework/operator-sdk#quick-start)
-	[MiniKube](https://github.com/kubernetes/minikube#installation)
-	[Go Lang](https://golang.org/doc/install)
-	[Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)

## Setup the EI Operator
1. Start the Kubernetes cluster (MiniKube)
2. Create the folder structure `$GOPATH/src/github.com/wso2` if not available and clone k8s-ei-operator git repo
	```
	git clone https://github.com/wso2/k8s-ei-operator.git
	```
3. Change directory to k8s-ei-operator
	```
	cd $GOPATH/src/github.com/wso2/wso2/k8s-ei-operator
	```
4. Deploy integration CustomResourceDefinition into Kubernetes cluster to understand custom resource type
	```
	kubectl create -f deploy/crds/integration_v1alpha1_integration_crd.yaml
	```
5. Run the operator code as a Go program outside the Kubernetes cluster
	```
	operator-sdk up local
	```
6. Apply configuration for the ingress controller
    ```
    kubectl apply -f deploy/config_map.yaml
    ```   

## Deploy the Integration Solutions with EI Operator
Deploy sample integration to start WSO2 micro integrator runtime which having 'User Info' API
```
kubectl apply -f deploy/crds/user_mgt_demo_integration.yaml
```
List the deployed integration
```
kubectl get integration
```

## Run the Integration Solution

Invoke the User Info API once STATUS becomes **Running** 

### With **Ingress Controller**
i. HTTP Request
```
curl http://wso2/user-mgt-demo-integration-service/userInfo/users
```
ii. HTTPS Request
```
curl https://wso2/user-mgt-demo-integration-service/userInfo/users -k
```

### Without **Ingress Controller**
i. Port forward
```
kubectl port-forward service/user-mgt-demo-integration-service 8290:8290
```
ii. Invoke the API
```
curl http://localhost:8290/userInfo/users
```