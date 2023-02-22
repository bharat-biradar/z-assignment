# Simple crud api

A simple api to manage a list of generic items (books, movies) etc... using go and mongodb.
## api interaction
### User interaction
- POST `api/user/signup` - Send a json in body to register as new user.
```json
{
    "username" : "some_username",
    "password" : "some_password",
    "firstName" : "firstName",
    "lastName" : "lastName"
}
```
- POST `api/user/login` - Login to existing account and recieve an api key for further requests.
  sample login request
```json
  {
    "username" : "some_username",
    "password" : "some_password",
    "firstName" : "firstName",
    "lastName" : "lastName"
}
```
 sample response
 ```json
 {
    "statusCode": 200,
    "status": "success",
    "data": {
        "token": "0d171765-ec88-45ff-81c9-c95c968029ef"
    }
}
 ```
- POST `/api/user/delete` - Delete account using username and password similar to login
  
### Managing items
Add api_key header in all the requests to authenticate access to items.
- GET `/api/items` - Retrieve all items
- POST `/api/items` - Add a new items. Data should be a valid json format.
- Get `/api/item/:id` - Get a specific item.
- Patch `/api/item/:id` - Update the item info.
- Delete `/api/item/:id` - Delete item if exists