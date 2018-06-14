# Kubernetes Autoscaler with AutoScalr Support


This repository contains a version of the Kubernetes Cluster AutoScaler component that works with AutoScalr to reduce the cost of running your k8s cluster by 50%-75% by changing one yaml file.

## Usage

If you have a working kubernetes/autoscaler deployment all you have to do enable AutoScalr and start saving money 
is to update 3 items in your existing yaml file:

- container image to: autoscalr/k8s_autoscalr
- cloud-provider to: autoscalr
- set desired parameter(s) in env section of yaml file

To illustrate, here are the changes in diff format required to the 1 ASG example yaml:


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
-       - image: gcr.io/google_containers/cluster-autoscaler:v0.6.0        
+       - image: autoscalr/k8s_autoscalr:v0.6.0
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
-           - --cloud-provider=aws
+           - --cloud-provider=autoscalr
            - --skip-nodes-with-local-storage=false
            - --nodes=1:10:k8s-worker-asg-1
          env:
            - name: AWS_REGION
              value: us-east-1
+           - name: AUTOSCALR_API_KEY
+             value: myApiKeyFromAutoScalr
+           - name: DISPLAY_NAME
+             value: nameToDisplayInUI
+           - name: INSTANCE_TYPES
+             value: m4.large,c4.large,r4.large
+           - name: MAX_SPOT_PERCENT_TOTAL
+             value: 90
+           - name: MAX_SPOT_PERCENT_ONE_MARKET
+             value: 20
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

## AUTOSCALR_API_KEY

If have an AutoScalr account, you can login at [autoscalr.com](https://app.autoscalr.com) to get your api key.

If you do not have an AutoScalr account, you can signup for a free 14 day trial [here](https://aws.amazon.com/marketplace/pp/B074N1N5QM).

## Environment Variable Reference

The following environment variables are supported:

* `AWS_REGION` - (Required) AWS Region the k8s cluster is running in
* `AUTOSCALR_API_KEY` - (Required) The api key provided by AutoScalr when you signup
* `DISPLAY_NAME` - (Required) Short name to be used in AutoScalr web UI display
* `INSTANCE_TYPES` - (Required) Comma delimited list of instance types to use
* `MAX_HOURS_INSTANCE_AGE` - (Optional, Default: off) When set, AutoScalr will schedule instance replacement if an instance's age exceeds this setting
* `MAX_SPOT_PERCENT_TOTAL` - (Optional, Default: 80) Maximum percentage of capacity to allow in Spot instances
* `MAX_SPOT_PERCENT_ONE_MARKET` - (Optional, Default: 20) Maximum percentage of capacity to allow in a single Spot market
* `OS_FAMILY` - (Optional, Default: Linux/UNIX) Options: Linux/Unix, SUSE Linux, Windows
* `TARGET_SPARE_CPU_PERCENT` - (Optional, Default: 20) Target spare cpu percentage to scale to, e.g. 20% spare capacity = 80% cpu utilization
* `TARGET_SPARE_MEMORY_PERCENT` - (Optional, Default: 20) Target spare memory percentage to scale to, e.g. 20% spare capacity = 80% memory utilization

## Supported Versions

The referenced image is built against Kubernetes 1.7 using upstream cluster-autoscaler 0.6.0

For use with other versions, check other branches or rebuild of the image from source.

Contact us if you need assistance with another version.

## Contact Info

support@autoscalr.com

or via chat interface on website: www.autoscalr.com




