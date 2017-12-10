# Kubernetes Autoscaler with AutoScalr Support


This repository contains a version of the Kubernetes Cluster AutoScaler component that works with AutoScalr to reduce the cost of running your k8s cluster by 50%-75%.

## Usage

If you have a working kubernetes/autoscaler deployment all you have to do is update 3 items in your existing yaml file:

- container image to: autoscalr/k8s_autoscalr
- cloud-provider to: autoscalr
- set desired parameter(s) in env section of yaml file

To illustrate, here are the changes required to the 1 ASG example yaml provided [here](https://github.com/kubernetes/autoscaler/blob/master/cluster-autoscaler/cloudprovider/aws/README.md#1-asg-setup-min-1-max-10-asg-name-k8s-worker-asg-1)


```yaml

---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: cluster-autoscaler
  namespace: kube-system
  labels:
    app: cluster-autoscaler
spec:
  replicas: 1
  selector:
    matchLabels:
      app: cluster-autoscaler
  template:
    metadata:
      labels:
        app: cluster-autoscaler
    spec:
      containers:
-        - image: gcr.io/google_containers/cluster-autoscaler:v0.6.0        
+        - image: autoscalr/k8s_autoscalr:v0.3.1
          name: cluster-autoscaler
          resources:
            limits:
              cpu: 100m
              memory: 300Mi
            requests:
              cpu: 100m
              memory: 300Mi
          command:
            - ./cluster-autoscaler
            - --v=4
            - --stderrthreshold=info
-            - --cloud-provider=aws
+            - --cloud-provider=autoscalr
            - --skip-nodes-with-local-storage=false
            - --nodes=1:10:k8s-worker-asg-1
          env:
            - name: AWS_REGION
              value: us-east-1
+            - name: AUTOSCALR_API_KEY
+              value: myApiKeyFromAutoScalr
+            - name: DISPLAY_NAME
+              value: nameToDisplayInUI
+            - name: MAX_SPOT_PERCENT_TOTAL
+              value: 90
+            - name: MAX_SPOT_PERCENT_ONE_MARKET
+              value: 20
+            - name: INSTANCE_TYPES
+              value: m4.large,c4.large,r4.large
          volumeMounts:
            - name: ssl-certs
              mountPath: /etc/ssl/certs/ca-certificates.crt
              readOnly: true
          imagePullPolicy: "Always"
      volumes:
        - name: ssl-certs
          hostPath:
            path: "/etc/ssl/certs/ca-certificates.crt"
```

## Contact Info

support@autoscalr.com

or via chat interface on website: www.autoscalr.com




