# Aliyun SDK for Go

The encapsulation of `Aliyun` HTTP API with limited support at present

**Install:**

    go get -u github.com/jiangshengwu/aliyun-sdk-for-go

**ECS Client Usage:**

    // Initialize client
	cli := &client.EcsClient{}
	cli.AccessKeyId = "Your Access Key Id"
	cli.AccessKeySecret = "Your Access Key Secrect"
	cli.Format = "JSON" //optional
    cli.Init() // Ready to use

    // Perform request
    result, err := cli.Execute("StopInstance", map[string]string{
        "InstanceId": "i-253op6931",
    })
    if (err != nil) {
        log.Fatal(err)
        return
    }
    fmt.Println(result)