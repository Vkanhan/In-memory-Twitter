## In-Memory-Twitter

### Create Users
```bash
curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"username":"elonmusk"}'
curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"username":"jackdorsey"}'
```

### Follow (user 1 follows user 2)
```bash
curl -X POST http://localhost:8080/users/2/follow -H "X-User-ID: 1" -H "Content-Type: application/json"
```

### Create Tweets
```bash
curl -X POST http://localhost:8080/tweets -H "X-User-ID: 1" -H "Content-Type: application/json" -d '{"content":"Just bought Twitter!"}'
curl -X POST http://localhost:8080/tweets -H "X-User-ID: 2" -H "Content-Type: application/json" -d '{"content":"Welcome to the team!"}'
```

### Like a Tweet
```bash
curl -X POST http://localhost:8080/tweets/1/like -H "X-User-ID: 2"
```

### Get Feeds
```bash
curl -X GET "http://localhost:8080/feeds/1"
curl -X GET "http://localhost:8080/feeds/2"
```

### Get User Profiles
```bash
curl -X GET http://localhost:8080/users/1
curl -X GET http://localhost:8080/users/2
```

### Cleanup (Delete Tweet)
```bash
curl -X DELETE http://localhost:8080/tweets/1 -H "X-User-ID: 1"
```
