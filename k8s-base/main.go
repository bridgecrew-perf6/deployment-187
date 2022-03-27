package main

import (
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	v1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/networking/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		appName := "nginx"
		appLabels := pulumi.StringMap{
			"app": pulumi.String(appName),
		}

		deployment, err := appsv1.NewDeployment(ctx, appName, &appsv1.DeploymentArgs{
			Spec: appsv1.DeploymentSpecArgs{
				Selector: &metav1.LabelSelectorArgs{
					MatchLabels: appLabels,
				},
				Replicas: pulumi.Int(2),
				Template: &corev1.PodTemplateSpecArgs{
					Metadata: &metav1.ObjectMetaArgs{
						Labels: appLabels,
					},
					Spec: &corev1.PodSpecArgs{
						Containers: corev1.ContainerArray{
							corev1.ContainerArgs{
								Name:  pulumi.String(appName),
								Image: pulumi.String("nginx"),
								Ports: &corev1.ContainerPortArray{
									&corev1.ContainerPortArgs{
										ContainerPort: pulumi.Int(80),
									},
								},
							}},
					},
				},
			},
		})
		if err != nil {
			return err
		}

		template := deployment.Spec.ApplyT(func(v *appsv1.DeploymentSpec) *corev1.PodTemplateSpec {
			return &v.Template
		}).(corev1.PodTemplateSpecPtrOutput)

		meta := template.ApplyT(func(v *corev1.PodTemplateSpec) *metav1.ObjectMeta { return v.Metadata }).(metav1.ObjectMetaPtrOutput)

		_, err = corev1.NewService(ctx, appName+"-service", &corev1.ServiceArgs{
			Metadata: meta,
			Spec: &corev1.ServiceSpecArgs{
				Ports: &corev1.ServicePortArray{
					&corev1.ServicePortArgs{
						Name:       pulumi.String("http"),
						Port:       pulumi.Int(80),
						TargetPort: pulumi.Int(80),
						Protocol:   pulumi.String("TCP"),
					},
				},
				Selector: appLabels,
			},
		})
		if err != nil {
			return err
		}

		host := "nginx.4l5.de"

		_, err = v1.NewIngress(ctx, appName+"-ingress", &v1.IngressArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Annotations: pulumi.StringMap{
					"traefik.ingress.kubernetes.io/router.entrypoints":      pulumi.String("web, websecure"),
					"traefik.ingress.kubernetes.io/router.middlewares":      pulumi.String("default-nginx-middleware@kubernetescrd"),
					"traefik.ingress.kubernetes.io/router.tls.certresolver": pulumi.String("traefikresolver"),
				},
			},
			Spec: &v1.IngressSpecArgs{
				Rules: v1.IngressRuleArray{
					&v1.IngressRuleArgs{
						Host: pulumi.String(host),
						Http: &v1.HTTPIngressRuleValueArgs{
							Paths: v1.HTTPIngressPathArray{
								&v1.HTTPIngressPathArgs{
									Path:     pulumi.String("/"),
									PathType: pulumi.String("Prefix"),
									Backend: &v1.IngressBackendArgs{
										Service: &v1.IngressServiceBackendArgs{
											// connection to endpoint not working, as name is changed.
											Name: pulumi.String(appName),
											Port: &v1.ServiceBackendPortArgs{
												Number: pulumi.Int(80),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		})

		if err != nil {
			return err
		}

		return nil
	})
}
