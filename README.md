# remoteit

`remoteit` is a command line interface to remot3.it.

The remot3.it API allows you to login, list devices, and request a proxy connection. You currently 
cannot add or edit devices from the API.

In order to use `remoteit` you will need a remot3.it account. You can signup for an account at 
[https://app.remote.it/auth/#/sign-up](https://app.remote.it/auth/#/sign-up).

I am in no way afflicated with remote3.it, I just happen to use it to connect to my
raspberry pi and didn't want to have to hit the website everytime I needed a proxy.

## Build

`remoteit` uses go modules. Building with go 1.11 should be as simple as `go build -o remoteit`

If you're using an earlier version of go you'll need to `go get` all the dependencies in
the `go.mod` file manually.

## Tests

Yeah, there aren't any. I hacked this together over a few evenings which is 
obvious from looking at the code. It suits my purposes. Pull requests are welcome.

## Usage

The following environment variables are used but can also be passed through flags
or a prompt in the case of your password.

- `REMOTEIT_APIKEY` - the api/develop key for your remot3.it account
- `REMOTEIT_USERNAME` - your remot3.it username
- `REMOTEIT_PASSWORD` - you guessed it

`remoteit` will create the `~/.remoteit` directory upon first launch. Two files
will be stored here, `login` and `devices`. `login` is a cache of your login token
which is good for one week. `devices` contains a cache of all your devices.

`login` is created and updated anytime you run the `remoteit login` command. `devices`
is created/updated when `remoteit devices` is run. You must be logged in for this 
command to succeed.

To create a proxy connection, run `remoteit connect`. 
You must be logged in for this command to succeed.

All commands return tabular output by default. The `--json` flag can be passed
to get JSON output instead.

#### login
```bash
$ remoteit login
Token             Expiry Unix  Expiry
-----             -----------  ------  
1234567890abcdef  1542507552   2018-11-17 21:19:12 -0500 EST
```

#### devices
```bash
$ ./remoteit devices
Alias        Address                  Service       Last IP
-----        -------                  -------       -------         
ssh-pi       00:00:00:00:00:00:00:00  SSH           127.0.0.1  
```

#### connect
The connect command can return tabular output or JSON for your to parse
yourself however it has some convenience flags also.

Currently only SSH is supported but if you wanted to connect to the `ssh-pi` device
as user `pi` you could run the following.
```bash
$ ssh pi@$(remoteit connect ssh-pi 127.0.0.1 --format ssh)
```

If you want to bypass using the cached devices you can add the `--nocache` flag.
If the device cache doesn't exist (because you have ran `remoteit devices`) running
`remoteit connect` will run as if `--nocache` was passed.

If you do not specify the `hostip` argument, the Last IP used to access the device
will be used, either from the cache or from a fresh API call.