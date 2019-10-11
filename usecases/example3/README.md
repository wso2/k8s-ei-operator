# k8s-ei-operator Example 3

## Message Routing Scenario

Let's define a JMS sender receiver scenario using WSO2 Micro Integrator and deploy it on your Kubernetes environment.

Follow the below steps to deploy and run the integration solution on Kubernetes.

1.  Start the Docker daemon in the host machine.
2.  Navigate to the **JMSSampleProject** Maven Multi Module Project given here.

    ```
    cd JMSSampleProject
    ```

3.  Run the following command to build the project. It will create a docker image with the provided target repository and tag once the build is successfull.
    ```bash
    mvn clean install -Dmaven.test.skip=true
    ```
4.  Run the `docker image ls` command to verify whether or not the docker image has been built. 

5.  Navigate to the Kubernetes project inside the MavenParentProject and run the following command to the push docker image to the remote docker registry. Here **username** and **password** are the credentials of the remote Docker registry.
    ```bash
    cd K8sJMSSenderReceiverSample
    mvn dockerfile:push -Ddockerfile.username={username} -Ddockerfile.password={password}
    ``` 
    
6.  [Install](https://github.com/wso2/k8s-ei-operator/blob/master/README.md#install-k8s-ei-operator) the **k8s-ei-operator**, if it is not installed in the Kubernetes environment.

7.  Deploy sample integration to start WSO2 micro integrator runtime which having 'ArithmaticOperationService' proxy-service

    ```
    kubectl apply -f kubernetes_cr.yaml
    ``` 
**Note**- Update above **tcp://localhost:61616** URL with the actual/connecting URL which will reachable from the Kubernetes pod.

8.  It will create a new queue called **firstQueue** in ActiveMQ. Send a message to this queue. Proxy-service will listen to this message and send that message to a new queue called **secondQueue**.  
    
**Note** - Follow the [JMS Sender Receiver Sample]() to implement this from the ground up. 