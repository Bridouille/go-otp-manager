## Otp manager in Golang ##

Web application to manage your Time-Based One Time Password.
The TOTP generated are compatible with Google Authenticator.

How to install : 

    go get github.com/Bridouille/go-otp-manager
    cd $GOPATH/src/github.com/Bridouille/go-otp-manager
    go build && ./go-otp-manager

That's it, the server is running on the port 3000.
To change the running port :

    export PORT=xxxx

Modify the configuration file at your convenience, the environment chosen will be the one in the OTP_MANAGER_ENV environment variable, development if not set : 

    ./config/config.json

You'll need to provide the `credentials` to access to the panel.

## What it look like ? ##

![enter image description here](http://image.noelshack.com/fichiers/2015/20/1431866875-capture-d-ecran-2015-05-17-a-14-47-32.png)

I used [Angular Material](https://material.angularjs.org/) for the front-end part, the TOTPs freshen themselves once arrived to expire.
