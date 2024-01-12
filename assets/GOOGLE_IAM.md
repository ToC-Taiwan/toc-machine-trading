# GOOGLE IAM

## Login gcloud

```bash
gcloud auth application-default login
```

## List all service accounts in Golang

```go
func list() {
 ctx := context.Background()

 iamService, err := iam.NewService(ctx)
 if err != nil {
  log.Fatal(err)
 }

 // Required. The resource name of the project associated with the service
 // accounts, such as `projects/my-project-123`.
 name := "projects/my-project-123" //

 req := iamService.Projects.ServiceAccounts.List(name)
 if err := req.Pages(ctx, func(page *iam.ListServiceAccountsResponse) error {
  for _, serviceAccount := range page.Accounts {
   // TODO: Change code below to process each `serviceAccount` resource:
   fmt.Printf("%#v\n", serviceAccount)
  }
  return nil
 }); err != nil {
  log.Fatal(err)
 }
}
```

## Create a service account key in Golang

```go
func create() {
 fileWriter, err := os.Create("configs/service_account.json")
 if err != nil {
  log.Fatalf("error creating file: %v\n", err)
 }
 defer fileWriter.Close()

 _, err = createKey(fileWriter, "tmt-331@trade-agent-87e47.iam.gserviceaccount.com")
 if err != nil {
  log.Fatalf("error creating key: %v\n", err)
 }
}
```

```go
func createKey(w io.Writer, serviceAccountEmail string) (*iam.ServiceAccountKey, error) {
 ctx := context.Background()
 service, err := iam.NewService(ctx)
 if err != nil {
  return nil, fmt.Errorf("iam.NewService: %w", err)
 }

 resource := "projects/-/serviceAccounts/" + serviceAccountEmail
 request := &iam.CreateServiceAccountKeyRequest{}
 key, err := service.Projects.ServiceAccounts.Keys.Create(resource, request).Do()
 if err != nil {
  return nil, fmt.Errorf("Projects.ServiceAccounts.Keys.Create: %w", err)
 }
 // The PrivateKeyData field contains the base64-encoded service account key
 // in JSON format.
 // TODO(Developer): Save the below key (jsonKeyFile) to a secure location.
 // You cannot download it later.
 jsonKeyFile, _ := base64.StdEncoding.DecodeString(key.PrivateKeyData)
 fmt.Fprint(w, string(jsonKeyFile))
 return key, nil
}
```
