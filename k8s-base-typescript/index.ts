import * as k8s from "@pulumi/kubernetes";
import * as  pulumi from "@pulumi/pulumi";

const config = new pulumi.Config();
const nextcloudAdminPw = config.requireSecret("nextcloudAdminPw");

const nextcloud = new k8s.helm.v3.Chart("nextcloud", {
    chart: "nextcloud",
    version: "2.13.2",
    fetchOpts: {
        repo: "https://nextcloud.github.io/helm/",
    },
    values: {
        ingress: {
            enabled: true,
            annotations: {
                "traefik.ingress.kubernetes.io/router.entrypoints": "web, websecure",
                "traefik.ingress.kubernetes.io/router.middlewares": "default-redirect-to-https@kubernetescrd",
                "traefik.ingress.kubernetes.io/router.tls.certresolver": "traefikresolver",
            }
        },
        nextcloud: {
            host: "nextcloud.alex-stadler.com",
            password: nextcloudAdminPw
        },
        mariadb: {
            enabled: true
        },
        persistence: {
            enabled: true,
            size: "1Gi"
        }
    },
});
