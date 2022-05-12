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

HTTP response status: **201**

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
