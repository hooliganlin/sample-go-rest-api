# Simple Go Rest API

This is a simple web server to demonstrate creating a simple API server and http client which hits 
https://jsonplaceholder.typicode.com.

## Build

```shell
$ go build -o simple-go-rest-api
```

## Run
```shell
$ ./simple-go-rest-api 
```
The server by default will run on port `8080` on host `127.0.0.1`. These can overriden with the follwing environment variables:

|Environment Variable | Default Value|
| ------ | ------ |
| MYAPP_SERVER_HOST |  127.0.0.1 |
| MYAPP_SERVER_PORT |  8080 | 

```shell
$ curl -s  "http://localhost:8080/v1/user-posts/1" 
```

This should yield a user JSON response with a user's associated posts.
```json
{
  "id": 1,
  "userInfo": {
    "name": "Leanne Graham",
    "username": "Bret",
    "email": "Sincere@april.biz"
  },
  "posts": [
    {
      "id": 1,
      "title": "sunt aut facere repellat provident occaecati excepturi optio reprehenderit",
      "body": "quia et suscipit\nsuscipit recusandae consequuntur expedita et cum\nreprehenderit molestiae ut ut quas totam\nnostrum rerum est autem sunt rem eveniet architecto"
    }
  ]
}
```