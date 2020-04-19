apiVersion: apps/v1
kind: Deployment
metadata:
  name: xrootd-operator
spec:
  replicas: 1
  selector:
    matchLabels:
      name: xrootd-operator
  template:
    metadata:
      labels:
        name: xrootd-operator
    spec:
      serviceAccountName: xrootd-operator
      containers:
        - name: xrootd-operator
          # Replace this with the built image name
          image: REPLACE_IMAGE
          command:
          - xrootd-operator
          imagePullPolicy: Always
          env:
            - name: WATCH_NAMESPACE
              valueFrom:
                fieldRef:
                  fieldPath: metadata.namespace
            - name: POD_NAME
              valueFrom:
                fieldRef:
                  fieldPath: metadata.name
            - name: OPERATOR_NAME
              value: "xrootd-operator"
