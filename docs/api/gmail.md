# Gmail
- Type = `gmails`

### Data
Adding gmail to scenario. There is just one sub type yet.
- Data request to create:
    ```json
    {
        "parentId": null,
        "name": "mygmail",
        "type": "gmails",
        "subType": "sendEmail",
        "position": {
            "x": 100,
            "y": 100
        },
        "data": {
                "body":"hello",
                "subject":"test",
                "to":["mansourikhahreza@gmail.com"],
                "email":"rezamansourikhah@gmail.com",
                "password":"XXXXXXXXXXXXXXXX"
        }
    }
    ```

- Data response
    ```json
   {
        "message": "node is created successfully.",
        "node": {
            "id": 116,
            "parentId": null,
            "name": "mygithub",
            "type": "gmails",
            "subType": "sendEmail",
            "position": {
                "x": 100,
                "y": 100
            },
            "data": {
                "body": "hello",
                "email": "rezamansourikhah@gmail.com",
                "password": "XXXXXXXXXXXXXXXX",
                "subject": "test",
                "to": [
                    "mansourikhahreza@gmail.com"
                ]
            }
        }
    }
    ```

