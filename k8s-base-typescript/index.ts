// Example inspired by:
// https://www.digitalocean.com/community/tutorials/how-to-set-up-an-nginx-ingress-on-digitalocean-kubernetes-using-helm
// Note that, as described in that article, you will need to configure DNS for hw1/hw2.your_domain_name.

import * as k8s from "@pulumi/kubernetes";

const apps = [];
const appBase = "hello-k8s";
const appNames = [`${appBase}-first`, `${appBase}-second`];
for (const appName of appNames) {
    const appSvc = new k8s.core.v1.Service(`${appName}-svc`, {
        metadata: {name: appName},
        spec: {
            type: "ClusterIP",
            ports: [{port: 80, targetPort: 8080}],
            selector: {app: appName},
        },
    });
    const appDep = new k8s.apps.v1.Deployment(`${appName}-dep`, {
        metadata: {name: appName},
        spec: {
            replicas: 3,
            selector: {
                matchLabels: {app: appName},
            },
            template: {
                metadata: {
                    labels: {app: appName},
                },
                spec: {
                    containers: [{
                        name: appName,
                        image: "paulbouwer/hello-kubernetes:1.8",
                        ports: [{containerPort: 8080}],
                        env: [{name: "MESSAGE", value: "Hello K8s!"}],
                    }],
                },
            },
        },
    });
    apps.push(appSvc.status);
}

const appIngress = new k8s.networking.v1.Ingress(`${appBase}-ingress`, {
    metadata: {
        name: "hello-k8s-ingress",
        annotations: {
            "traefik.ingress.kubernetes.io/router.entrypoints": "web, websecure",
            "traefik.ingress.kubernetes.io/router.middlewares": "default-nginx-middleware@kubernetescrd",
            "traefik.ingress.kubernetes.io/router.tls.certresolver": "traefikresolver",

        },
    },
    spec: {
        rules: [
            {
                host: "test.4l5.de",
                http: {
                    paths: [{
                        pathType: "Prefix",
                        path: "/",
                        backend: {
                            service: {
                                name: appNames[0],
                                port: {number: 80},
                            },
                        },
                    }],
                },
            },
            {
                host: "nginx.4l5.de",
                http: {
                    paths: [{
                        pathType: "Prefix",
                        path: "/",
                        backend: {
                            service: {
                                name: appNames[1],
                                port: {number: 80},
                            },
                        },
                    }],
                },
            },
        ],
    },
});

export const appStatuses = apps;
// export const controllerStatus = ctrl.status;