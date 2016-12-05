# Go Language - API sample


## Web Framework 

[Gin Web Framework](https://github.com/gin-gonic/gin)


```
go get gopkg.in/gin-gonic/gin.v1
```

## ORM library

[GORM](https://github.com/jinzhu/gorm)

```
go get -u github.com/jinzhu/gorm
```
	
## API
	
Hello world!
```$xslt
http://localhost:8080/
```

All Users
```
http://localhost:8080/users/
```

### Example

```$xslt
{
  "users": [
    {
      "id": 1,
      "name": "ユーザー1",
      "items": [
      {
        "id": 1, 
        "name": "マリオブラザース", 
        "score": 5, 
        "user_id": 1
      }, 
      {
        "id": 2,
        "name": "スーパーマリオブラザース", 
        "score": 0, 
        "user_id": 1
      }, 
      ...          {
      ] 
    }, 
    {
      "id": 2,
      "name": "ユーザー1",
      "items": [
      {
        "id": 9, 
        "name": "マリオブラザース", 
        "score": 1, 
        "user_id": 2
      }, 
      ...
      ] 
    },
    ...
  ]
}
```


User & MatchingUsers
```
http://localhost:8080/users/:id
```

### Example

```
{
  "user": {
    "id": 2,
    "name": "ユーザー2",
    "items": [
      {
        "id": 9, 
        "name": "マリオブラザース", 
        "score": 1, 
        "user_id": 2
      }, 
      {
        "id": 10, 
        "name": "スーパーマリオブラザース", 
        "score": 1, 
        "user_id": 2
      }, 
      {
        "id": 11, 
        "name": "ゼルダの伝説", 
        "score": 3, 
        "user_id": 2
      },
      ...
    ]
  },
  "matching_users": [
    {
      "id": 1, 
      "name": "ユーザー1", 
      "score": 0.10886898139200682
    }, 
    {
      "id": 9, 
      "name": "ユーザー9", 
      "score": 0.1
    }, 
    {
      "id": 7, 
      "name": "ユーザー7", 
      "score": 0.08462632608958592
    }, 
    {
      "id": 4, 
      "name": "ユーザー4", 
      "score": 0.08270931562630669
    },
    ...
  ]
}
```

## Algorithm
The basic scoring algorithm is Euclidean distance

```math
d(p,q)=d(q,p)=\sqrt{(q_1-p_1)^2+(q_2-p_2)^2+\cdots+(q_n+p_n)^2}
```


## Author

**Daisuke Takayama**
* [@webcyou](https://twitter.com/webcyou)
* [@panicdragon](https://twitter.com/panicdragon)
* <https://github.com/webcyou>
* <https://github.com/panicdragon>
* <http://www.webcyou.com/>
