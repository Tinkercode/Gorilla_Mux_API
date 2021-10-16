# Gorilla_Mux_API

For Generating base64 encoded protobuf string:

```
request := &assignmentpb.PatchRequest{UserId: 823740, Email: "tarun@gmail.com"}
req, err := proto.Marshal(request)
if err != nil {
    log.Fatalf("Unable to marshal request : %v", err)
}
bsedat := base64.StdEncoding.EncodeToString(req)
fmt.Println("bsedat", bsedat)
```
