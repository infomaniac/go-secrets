# SecretManager

The goal of this library is to simplify using the Google Secrets Manager.
It provides a simple and effective way of getting and setting secrets within a given project.

The client is initialized with `New(..)` and closed with `Close()`.  
Secrets can then handled with `Get(..)`, `Set(..)` or `Delete(..)` in a simple manor.  
The client always returns the value of the latest version, so you don't have to deal with the different versions.

For more possibilities and to use the full feature set of the service, please refer to the [official client libraries](https://pkg.go.dev/cloud.google.com/go/secretmanager/apiv1#DefaultAuthScopes).

Like all of the provided official client libraries, the `application-default` authentication is used.  
If you get some *permission denied* or similar error, please make sure you're authenticated with GCP and have the correct access for the project you are trying to use. (Execute `gcloud auth application-default login` if unsure).


## Example Usage
```go
func main() {
    // create a context
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    s, err := New(ctx, "<my-google-project>", "europe-west1", "europe-west6")
    if err != nil {
        log.Fatal(err)
    }
    // don't forget to close the client
    defer s.Close()

    // set the value of a secret
    err = s.Set("<my-secret-id>", []byte("my secret value"))
    if err != nil {
        log.Fatal(err)
    }

    // get the latest value of a secret
    secret, err := s.Get("<my-secret-id>")
    if err != nil {
        log.Fatal(err)
    }

    // delete secret
    err = s.Delete("<my-secret-id>")
    if err != nil {
        log.Fatal(err)
    }	
}
```

# Work in progress
This library is work in progress. Feel free to open requests or bugs in the project's issue tracking. Thanks!