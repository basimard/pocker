
## API Reference

#### Create a new deck

```http
  POST /v1/create-deck?${shuffle}&{cards}
```

| Parameter | Type     | Usage                |
| :-------- | :------- | :------------------------- |
| `shuffle` | `string` | `true or false` |
| `cards` | `string` | `AS,2S` |

#### Open a deck

```http
  GET /v1/open-deck/${deck_id}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `deck_id`      | `string` | `uuid deck id` |

#### Draw cards from deck
```http
  POST /v1/draw-cards/${deck_id}&${count}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `deck_id`      | `string` | `uuid deck id` |
  `count`      | `int` | `1,2,3` |



## Run Service

Navigate to the root directory of the cloned repository where the file "**main.go**" is located.

**RUN** 
    
    
    go run .


