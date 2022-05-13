# auction-api doc:

This is a CRUD based Web services implementation with following assumed domain requirement:

The goal of the platform is to provide its users ability to auction their stuff to anyone interested in buying that on the same platform. 
Anyone could add the item they want to auction and then other people could start seeing what has been added and can bid on the price they would like to buy a particular item.

as its auction based multiple users can bid for an item on the platform.

### APIs

- `POST /api/users/register`

users can register themselves through this endpoint. it accepts the following request payload format:

```
{
    "firstName": "John",
    "lastName": "Doe",
    "email": "johndoe23@abc.com",
    "password": "secret12"
}
```

there are some basic validations like, `email` must match email format and `password` must be atleast length 8.

Which should result in following response when success:

HTTP response status: **201 Created**

```
{
    "id": 1,
    "uuid": "74c2cb29-8441-4f0f-8138-f3bddc4bf3c3",
    "createdAt": "2022-05-12T06:31:10.921Z",
    "updatedAt": "2022-05-12T06:31:10.921Z",
    "deletedAt": null,
    "version": 1,
    "firstName": "John",
    "lastName": "Doe",
    "email": "johndoe23@abc.com",
    "isActive": true,
    "createdBy": null,
    "updatedBy": null,
    "deletedBy": null
}
```

- `POST /api/users/login`

Users should be able to login through this endpoint. It accepts the following request payload:

```
{
    "email": "johndoe23@abc.com",
    "password": "secret12"
}
```

Which should result in the following response when success:

HTTP response status code: **200 OK**
```
{
    "id": 1,
    "uuid": "74c2cb29-8441-4f0f-8138-f3bddc4bf3c3",
    "createdAt": "2022-05-12T06:31:10.921Z",
    "updatedAt": "2022-05-12T06:31:10.921Z",
    "deletedAt": null,
    "version": 1,
    "firstName": "John",
    "lastName": "Doe",
    "email": "johndoe23@abc.com",
    "isActive": true,
    "createdBy": null,
    "updatedBy": null,
    "deletedBy": null,
    "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCIsInVzZXJuYW1lIjoiam9obmRvZTIzQGFiYy5jb20ifQ.eyJpc3MiOiJqb2huZG9lMjNAYWJjLmNvbSIsInN1YiI6IkFDQ0VTU19UT0tFTiIsImV4cCI6MTY1MjM0MDkyMywibmJmIjoxNjUyMzM3MzIzLCJpYXQiOjE2NTIzMzczMjMsImp0aSI6IjgwYTFlNGU4LWNhM2YtNGRiMS1hZjc2LTA3NmVkYjk5MzdlNiJ9.J5kwv8-ALlijcL71RjYRcm0rq26V1WxuL7dUXIHXZfY",
    "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCIsInVzZXJuYW1lIjoiam9obmRvZTIzQGFiYy5jb20ifQ.eyJpc3MiOiJqb2huZG9lMjNAYWJjLmNvbSIsInN1YiI6IlJFRlJFU0hfVE9LRU4iLCJleHAiOjE2NTc1MjEzMjMsIm5iZiI6MTY1MjMzNzMyMywiaWF0IjoxNjUyMzM3MzIzLCJqdGkiOiI1MjUwNmQxMS04Y2ZhLTQ1NjUtOWI1ZC1iOWFmYTIzOWEyMTYifQ.6_qiAMF6h6AnJQ00aK2SfzNdV8KmQEQySsar9MNpP1s"
}
```

If invalid credentials were provided, then this should be the error response:

HTTP response status code: **401 Unauthorized**

```
{
    "error": "Invalid credentials were found."
}
```

- `POST /api/items`

This endpoint can be used to add a new item on the platform. this is an authenticated endpoint.

To authorize a request use Auth header:(This is generic throughout all authenticated endpoints)

`Authorization: <token from login endpoint response>`

It accepts the following request payload:

```
{
    "name": "Test Item 1",
    "description": "Item description goes here.",
    "category": 0,
    "brandName": "ABC Brand",
    "marketValue": 100,
    "lastBidDate": "2023-01-02T15:04:05Z"
}
```

Which should result in following success response:

HTTP response status code: **201 Created**

```
{
    "id": 1,
    "uuid": "44627618-c9a1-48e9-b027-a50d15ecf861",
    "createdAt": "2022-05-12T06:57:20.663Z",
    "updatedAt": "2022-05-12T06:57:20.663Z",
    "deletedAt": null,
    "version": 1,
    "name": "Test Item 1",
    "description": "Item description goes here.",
    "category": 0,
    "brandName": "ABC Brand",
    "marketValue": 100,
    "lastBidDate": "2023-01-02T15:04:05Z",
    "isOffBid": false,
    "Bids": null,
    "itemImages": null,
    "reactions": null
}
```

- `POST /api/items/:itemId/images` (limited to item's author)

This endpoint allows to add images for an auction item. it accepts a multipart form-data content request with **images** param containing the list of item files.

Here is sample cURL request example:

```
curl --location --request POST 'localhost:8080/api/items/1/images' \
--header 'Authorization: token' \
--form 'images=@"/path/to/image/file.png"'
```

which when successful would result in sample success response:

HTTP response status code: **201 Created**

```
[
    {
        "id": 1,
        "uuid": "c3e0e077-68ee-4afe-89e4-8f82ff10f1b8",
        "createdAt": "2022-05-12T07:10:04.493Z",
        "updatedAt": "2022-05-12T07:10:04.493Z",
        "deletedAt": null,
        "version": 1,
        "item": null,
        "name": "gopher_772cf3ec-9d57-4ef5-a000-fb1a7a4ff1a4.png",
        "isThumbnail": true
    }
]
```

There are basic validations for images, it allows max 5 images upload and max image size can be 2 MBs.
Only accepted file types at this point are **jpeg, png**

**Thumbnail:** The first uploaded item image is automatically tagged as Thumbnail to show on items list.
(there are APIs that allow editing this as well later in the document.)

The list of images can also be overridden, by the use of `removeExisting=true` query param in this way while submitting:

```
POST /api/items/1/images?removeExisting=true
```


- `GET /api/items`

This anonymous endpoint is to list the items on platform. this doesn't require any request body content.

It should result in following response message format.

HTTP response status code: **200 OK**

```
{
    "items": [
        {
            "id": 2,
            "uuid": "8e2f54e5-c07d-4236-b91e-5a228e250490",
            "createdAt": "2022-05-12T07:01:08.436Z",
            "updatedAt": "2022-05-12T07:01:08.436Z",
            "deletedAt": null,
            "version": 1,
            "createdBy": {
                "id": 1,
                "uuid": "74c2cb29-8441-4f0f-8138-f3bddc4bf3c3",
                "createdAt": "2022-05-12T06:31:10.921Z",
                "updatedAt": "2022-05-12T06:31:10.921Z",
                "deletedAt": null,
                "version": 1,
                "firstName": "John",
                "lastName": "Doe",
                "email": "johndoe23@abc.com",
                "isActive": true,
                "createdBy": null,
                "updatedBy": null,
                "deletedBy": null
            },
            "updatedBy": {
                "id": 1,
                "uuid": "74c2cb29-8441-4f0f-8138-f3bddc4bf3c3",
                "createdAt": "2022-05-12T06:31:10.921Z",
                "updatedAt": "2022-05-12T06:31:10.921Z",
                "deletedAt": null,
                "version": 1,
                "firstName": "John",
                "lastName": "Doe",
                "email": "johndoe23@abc.com",
                "isActive": true,
                "createdBy": null,
                "updatedBy": null,
                "deletedBy": null
            },
            "name": "Test Item 1",
            "description": "Item description goes here.",
            "category": 0,
            "brandName": "ABC Brand",
            "marketValue": 100,
            "lastBidDate": "2023-01-02T00:00:00Z",
            "isOffBid": false,
            "Bids": null,
            "itemImages": [],
            "reactions": null
        },
        {
            "id": 1,
            "uuid": "44627618-c9a1-48e9-b027-a50d15ecf861",
            "createdAt": "2022-05-12T06:57:20.663Z",
            "updatedAt": "2022-05-12T06:57:20.663Z",
            "deletedAt": null,
            "version": 1,
            "createdBy": {
                "id": 1,
                "uuid": "74c2cb29-8441-4f0f-8138-f3bddc4bf3c3",
                "createdAt": "2022-05-12T06:31:10.921Z",
                "updatedAt": "2022-05-12T06:31:10.921Z",
                "deletedAt": null,
                "version": 1,
                "firstName": "John",
                "lastName": "Doe",
                "email": "johndoe23@abc.com",
                "isActive": true,
                "createdBy": null,
                "updatedBy": null,
                "deletedBy": null
            },
            "updatedBy": {
                "id": 1,
                "uuid": "74c2cb29-8441-4f0f-8138-f3bddc4bf3c3",
                "createdAt": "2022-05-12T06:31:10.921Z",
                "updatedAt": "2022-05-12T06:31:10.921Z",
                "deletedAt": null,
                "version": 1,
                "firstName": "John",
                "lastName": "Doe",
                "email": "johndoe23@abc.com",
                "isActive": true,
                "createdBy": null,
                "updatedBy": null,
                "deletedBy": null
            },
            "name": "Test Item 1",
            "description": "Item description goes here.",
            "category": 0,
            "brandName": "ABC Brand",
            "marketValue": 100,
            "lastBidDate": "2023-01-02T00:00:00Z",
            "isOffBid": false,
            "Bids": null,
            "itemImages": [
                {
                    "id": 1,
                    "uuid": "c3e0e077-68ee-4afe-89e4-8f82ff10f1b8",
                    "createdAt": "2022-05-12T07:10:04.493Z",
                    "updatedAt": "2022-05-12T07:10:04.493Z",
                    "deletedAt": null,
                    "version": 1,
                    "item": null,
                    "name": "gopher_772cf3ec-9d57-4ef5-a000-fb1a7a4ff1a4.png",
                    "isThumbnail": true
                }
            ],
            "reactions": null
        }
    ],
    "page": 0,
    "size": 10,
    "max_page": 1,
    "total_pages": 1,
    "total": 2,
    "last": false,
    "first": true,
    "visible": 2
}
```

- `GET /api/items/:itemId/images/:imageId`

This anonymous endpoint allows to fetch an item's image using **itemId** and **imageId**

which should send back its corresponding image file back to response stream.

- `PATCH  /api/items/:itemId` (limited to item's author)

This authorized endpoint allows item authors to edit the item details. It accepts the following request payload:

its same payload format as `POST /api/items` just with updated details.

```
{
    "name": "Test Item 1",
    "description": "Item's updated description goes here.",
    "category": 0,
    "brandName": "ABC Brand updated",
    "marketValue": 110,
    "lastBidDate": "2022-10-02T15:04:05Z"
}
```

Which on success should return HTTP response status: **204 No Content** with empty response message body content.

- `PUT /api/items/:itemId/mark-off-bid` (limited to item's author)

once an item is added on platform, it can't be deleted but can be put off bid. which can be done using this above endpoint.
Which can only be done by item's author. once an item is marked/ put off bid then other users won't be able to bid on this
particular item.

It doesn't take any request body content and on success it would return HTTP response status: **204 No Content** with empty response message body content.

- `DELETE /api/items/:itemId/images` (limited to item's author)

This endpoint allows deleting/clearing all images from an item. this endpoint doesn't take any request body content.
On success, it should return HTTP response status: **204 No Content** with empty response message body content.

- `DELETE /api/items/:itemId/images/:imageId` (limited to item's author)

This endpoint allows deleting/clearing a particular image from an item. this endpoint doesn't take any request body content.
On success, it should return HTTP response status: **204 No Content** with empty response message body content.

- `PATCH /api/items/:itemId/images/:imageId/make-thumbnail` (limited to item's author)

This endpoint allows marking an item's image thumbnail, it doesn't take any request body content and on success,
it returns HTTP response status: **204 No Content** with empty response message body content.

- `DELETE /api/items/:itemId/images/remove-thumbnail` (limited to item's author)

This endpoint allows removing/unsetting an item's image thumbnail, it doesn't take any request body content and on success,
it returns HTTP response status: **204 No Content** with empty response message body content.

**Note:** It doesn't delete the existing item's thumbnail image just takes thumbnail tag off of it.

- `POST /api/items/:itemId/bid`

This endpoint allows other people(Non-Item authors) to put their bid on an item. this API can be used to add a new bid by a 
user or update an existing bid placed by a user.

This accepts the following request body content:

```
{
    "bidValue": 12
}
```

Which on successful operation should return the following sample response message:

HTTP response status code: **201 Created**

with body content:
```
{
    "id": 1,
    "uuid": "58cfd6f1-4034-4f85-b3c5-1f178575251d",
    "createdAt": "2022-05-12T08:19:49.053Z",
    "updatedAt": "2022-05-12T08:19:49.053Z",
    "deletedAt": null,
    "version": 1,
    "item": null,
    "bidValue": 12
}
```

- `GET /api/items/:itemId/bids`

This is an anonymous endpoint which allows anyone on the platform to see the list of bids for an item. It doesn't 
require any request payload content.

It returns the following sample response message:

HTTP response status code: **200 OK**

with body content:
```
{
    "items": [
        {
            "id": 1,
            "uuid": "58cfd6f1-4034-4f85-b3c5-1f178575251d",
            "createdAt": "2022-05-12T08:19:49.053Z",
            "updatedAt": "2022-05-12T08:19:49.053Z",
            "deletedAt": null,
            "version": 1,
            "createdBy": {
                "id": 2,
                "uuid": "ffde17f1-603b-466d-9ba2-99440ac173fa",
                "createdAt": "2022-05-12T08:19:09.962Z",
                "updatedAt": "2022-05-12T08:19:09.962Z",
                "deletedAt": null,
                "version": 1,
                "firstName": "John",
                "lastName": "Doe",
                "email": "johndoe24@abc.com",
                "isActive": true,
                "createdBy": null,
                "updatedBy": null,
                "deletedBy": null
            },
            "item": {
                "id": 1,
                "uuid": "44627618-c9a1-48e9-b027-a50d15ecf861",
                "createdAt": "2022-05-12T06:57:20.663Z",
                "updatedAt": "2022-05-12T07:52:29.06Z",
                "deletedAt": null,
                "version": 1,
                "name": "Test Item 1",
                "description": "Item's updated description goes here.",
                "category": 0,
                "brandName": "ABC Brand updated",
                "marketValue": 110,
                "lastBidDate": "2022-10-02T00:00:00Z",
                "isOffBid": false,
                "Bids": null,
                "itemImages": null,
                "reactions": null
            },
            "bidValue": 12
        }
    ],
    "page": 0,
    "size": 10,
    "max_page": 1,
    "total_pages": 1,
    "total": 1,
    "last": false,
    "first": true,
    "visible": 1
}
```

- `POST /api/items/:itemId/reaction`

This endpoint allows authenticated users to react to items (at this point it only supports Like and Dislike).

It accepts the following request format:

```
{
  "reactionType": 0
}
```

`reactionType -> 0 (for Like reaction), reactionType -> 1 (for Dislike reaction)`

on success, it should return following sample response format:

HTTP response status code: **201 Created**

```
{
    "id": 1,
    "uuid": "dae06008-a928-4b28-8601-da373dc5b44d",
    "createdAt": "2022-05-13T06:48:51.539Z",
    "updatedAt": "2022-05-13T06:48:51.539Z",
    "deletedAt": null,
    "version": 1,
    "item": null,
    "type": 0
}
```

The same endpoint can also be used to update existing reaction on item by a user.

- `DELETE /api/items/:itemId/reaction`

This endpoint can be used to delete an existing user's reaction from an item. This API doesn't require any request body payload.

on success, it should return, HTTP response status: **204 No Content** with empty response message body content.

- `POST /api/items/:itemId/comment`

This endpoint allows authenticated users to add a comment on an Item. one user can add many comments on an item as well.

It accepts the following request payload format:

```
{
	"comment": "this is a test item comment."
}
```

which on success, should return following sample response message:

HTTP response status code: **201 Created**

```
{
    "id": 1,
    "uuid": "da9f49ba-7d2c-4481-9b71-592e972b7249",
    "createdAt": "2022-05-13T07:31:39.912Z",
    "updatedAt": "2022-05-13T07:31:39.912Z",
    "deletedAt": null,
    "version": 1,
    "description": "this is a test item comment.",
    "item": null
}
```

- `PATCH /api/items/:itemId/comment/:commentId`

This endpoint allows a comment author to edit item comment. It accepts the following request
payload format:

```
{
	"comment": "this is an updated test item comment."
}
```

on success, it should return, HTTP response status: **204 No Content** with empty response message body content.

- `DELETE /api/items/:itemId/comment/:commentId`

This endpoint allows a comment author to delete an item comment. It doesn't require any request 
body payload.

on success, it should return, HTTP response status: **204 No Content** with empty response message body content.

