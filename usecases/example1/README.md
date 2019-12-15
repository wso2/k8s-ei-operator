# k8s-ei-operator Example 1

## Hello World Scenario

Let's define a basic Hello World scenario using WSO2 Micro Integrator and deploy it on your Kubernetes environment.

Follow the below steps to deploy and run the integration solution on Kubernetes.

1.  Start the Docker daemon in the host machine.
2.  Navigate to the **HelloWorldSampleProject** Maven Multi Module Project given here.

    ```
    cd HelloWorldSampleProject
    ```

3.  Run the following command to build the project. It will create a docker image with the provided target repository and tag once the build is successfull.
    ```bash
    mvn clean install -Dmaven.test.skip=true
    ```
4.  Run the `docker image ls` command to verify whether or not the docker image has been built. 

5.  Navigate to the Kubernetes project inside the MavenParentProject and run the following command to the push docker image to the remote docker registry. Here **username** and **password** are the credentials of the remote Docker registry.
    ```bash
    cd K8sHelloWorldSample
    mvn dockerfile:push -Ddockerfile.username={username} -Ddockerfile.password={password}
    ``` 
    
6.  [Install](https://github.com/wso2/k8s-ei-operator/blob/master/README.md#install-k8s-ei-operator) the **k8s-ei-operator**, if it is not installed in the Kubernetes environment.

7.  Deploy sample integration to start WSO2 micro integrator runtime which having 'HelloWorld' proxy-service

    ```
    kubectl apply -f kubernetes_cr.yaml
    ```

8.  Port forward to expose the cluster port:

    ```
    kubectl port-forward service/hello-world-service 8290:8290
    ```
9. Invoke the service as follows:

    ```bash
    curl http://localhost:8290/services/HelloWorld
    ```  
    
**Note** - Follow the [Hello World Sample](https://ei.docs.wso2.com/en/latest/micro-integrator/setup/deployment/k8s-samples/hello-world/) to implement this from the ground up. 