terraform-provider-sealedsecrets
===
`sealedsecrets` provider checks if the sealed secret YAML file exists. If the secret file exists and hasn't been updated, it will skip generating and applying the secret. Otherwise, a new secret will be generated and applied based on the the input file.

### Installation
Please download it from [here v0.1.0](https://github.com/rockyhmchen/terraform-provider-sealedsecrets/releases/download/v0.1.0/terraform-provider-sealedsecrets_v0.1.0) and place it to `~/.terraform.d/plugins/`, and then make it executable by the command `chmod +x {file name}` before running `terraform init`.

### Example Usage
```HCL
resource "sealedsecrets_secret" "some_secret" {
  name      = "k8s-secret-name"
  namespace = "default"
  type      = "Opaque"

  secret_source        = "${var.keys_directory}/key.pem"
  sealed_secret_source = "${var.secret_directory}/some-secret.yaml"
  certificate          = ".tmp/sealed-secrets-controller.crt"

  depends_on = [var.sealed_secrets_controller_id]
}
```
### Requirements
- `kubectl` is a command line interface for running commands against Kubernetes clusters
- `kubeseal` utility uses asymmetric crypto to encrypt secrets that only the controller can decrypt. [install via Homebrew](https://github.com/bitnami-labs/sealed-secrets#homebrew)

### Argument Reference
The following arguments are supported:
- `name` - Name of the secret, must be unique.
- `namespace` - Namespace defines the space within which name of the secret must be unique.
- `type` -  The secret type. ex: `Opaque`
- `secret_source` - Path to the file containing the secret data. The file must be existing.
- `sealed_secret_source` - Path to the sealed secret YAML file. If the sealed secret doesn't exist, a new one will be generated at runtime. Otherwise, the file will be applied when it is new or has been updated.
- `certificate` - Path to the public key / certificate of [Sealed Secrets Controller](https://github.com/bitnami-labs/sealed-secrets) that is used for sealing secrets.
- `depends_on` - For specifying hidden dependencies.

*NOTE: All the arguments above are required*

### Use cases
- First deployment
  1. Make sure all input files/keys have been created in the `.keys` directory. The scripts within the `scripts` directory can be used to create these files.
  2. Run `terraform apply`
  
- Add a new secret
  1. Place input file(s)/key(s) in the corresponding directory (default: `.keys`).
  2. Declare the resource as specified in the **Example Usage** section above
  3. Run `terraform apply`

- Replace the current secrets
  1. Place input file(s)/key(s) in the corresponding directory (default: `.keys`).
  2. Delete the current sealed secret YAML(s) in the secret directory (default: `secrets/`)
  3. Run `terraform apply`

- Delete a secret
  1. Removed the sealed secret resource declaration and the sealed secret YAML file
  2. Run `terraform apply`

### How it works
#### When create
Checks if `sealed_secret_source` exists. If so, applies `sealed_secret_source` with `kubectl apply -f`, and then stores the file content in Terraform state file in `SHA256` hash. If not, firstly requires `secret_source` and `certificate` (certificate/public key of SealedSecrets Controller) to generate the sealed secret YAML file as the value of `sealed_secret_source` by the command below, and then applies the sealed secret with `kubectl apply -f`. After successfully applying it, stores the file content in Terraform state file in `SHA256` hash. 
```bash
    kubectl create secret generic {sealedsecrets_secret.[resource].name} \    
    --namespace={sealedsecrets_secret.[resource].namespace} \    
    --type={sealedsecrets_secret.[resource].type} \    
    --from-file={filename(sealedsecrets_secret.[resource].secret_source)}={sealedsecrets_secret.[resource].secret_source} \    
    --dry-run \    
    --output=yaml | \    
    kubeseal > {sealedsecrets_secret.[resource].sealed_secret_source} \    
    --cert ${sealedsecrets_secret.[resource].certificate} \    
    --format yaml
```

#### When delete
First, deletes the secret and the sealed secret in Kubernetes by the commands `kubectl delete secret {sealedsecrets_secret.[resource].name} -n {sealedsecrets_secret.[resource].namespace}` and `kubectl delete SealedSecret {sealedsecrets_secret.[resource].name} -n {sealedsecrets_secret.[resource].namespace}`.
Second, deletes the physical file of `sealed_secret_source`.
Last, removes the stored state of `sealedsecrets` resource from Terraform state file

#### When update
First, runs `delete` process
And then, runs `create` process

#### When read
Compares the file content of `sealed_secret_source` in `SHA256` hash with the stored state of `sealedsecrets` resource in Terraform state file