# Authentication

### /api/auth/signup `POST`
To create a new user.

- Request
    ```json
    {
        "username": "jafar",
        "email": "jafar@gmail.com",
        "password": "123123123"
    }
    ```
- Response
    ```json
    {
        "message": "user created successfully."
    }
    ```

### /api/auth/login `POST`
To login and get the access token.
- Request
    ```json
    {
        "username": "legato",
        "password": "1234qwer"
    }
    ```
- Response
    ```json
    {
        "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImxlZ2F0byIsImV4cCI6MTYxNjAxODYwMH0.FX_zlYPGn-ypy2KPVmgj-oG2Hx-LGluDF_0fi_fXJkQ"
    }
    ```

### /api/auth/refresh `POST`
To refresh the access token by set it in the request's header.
- Header
    - `Authorization` = `access_token`
- Response
    ```json
    {
        "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InJlemEiLCJleHAiOjE2MTYwMTgwODF9.sMYMNz0Pskr1cfOk19Dimdz6ZAuVbrKjHbLodB8pvPU"
    }
    ```

### /api/auth/user `GET`
To get logged in user.
- Header
    - `Authorization` = `access_token`
- Response
    ```json
    {
        "user": {
            "email": "legato@gmail.com",
            "username": "legato"
        }
    }
    ```


### /api/auth/protected `GET`
To test that only authorized users could see this page.
> This api is for testing purposes, and it returns all of existing users.
- Header
    - `Authorization` = `access_token`
 - Response
    ```json
    {
       "users": [
            {
                "email": "legato@gmail.com",
                "username": "legato"
            },
            {
                "email": "sedaqi@gmail.com",
                "username": "insomnia"
            },
            {
                "email": "amin@gmail.com",
                "username": "amin"
            },
            {
                "email": "amir@gmail.com",
                "username": "amir"
            },
            {
                "email": "ahmad@gmail.com",
                "username": "ahmad"
            }
        ]
    }
    ```