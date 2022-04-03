import {generate} from "generate-password";
import * as  pulumi from "@pulumi/pulumi";

export function getOrCreateSecret(config: pulumi.Config, key: string) {
    if (config.getSecret(key)) {
        return config.getSecret(key)
    } else {
        return pulumi.secret(generate({length: 25, numbers: true}))
    }
}
