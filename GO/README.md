# SQLite Cloud GO Client

## Get started

### Fetching the code:
```console
cd
mkdir test
cd test
git clone https://github.com/sqlitecloud/sdk
User: <github user name>
Password: <access-token or password>

```

### Setting up the development environment
```console
go env -w GO111MODULE=off
cd sdk/GO
export GOPATH=`pwd`
echo $GOPATH

```

### Install and compile the pre-requirements
First, you will have to install all the pre-requirements on your machine: `make install-prerequirements`. This will also patch the source code of the linenoise package...

### Synthesizing the C Proxy
The new GO lines are generated out of the C SDK. Every time you change something in the C SDK, you have to go to the SqliteCloud/sdk/GO (`cd $GOPATH`) directory and there you must enter: `make proxy`. This will generate a fresh GO source file out of the C files. Everything gets embedded and then, when you compile a GO program, there will be no external dependencies

### Precompile the static csdk library
If you like the pre-compiled C SDK more than the embedded C SDK, you can enter `make csdk`

### Building the test programs
If you want to build the Test programs: `make test`

### Building the CLI App
To build the CLI App (Warning: not fully functional, this is officially Step 1), you have to enter: `make cli`

### Build all at the same time:
If you want to do all at the same time: `make all`

## Documentation
If you want to see the Documentation: `make doc` - Warning: A browser window will open and display the documentation to you. The Documentation is updated live while coding. To stop the live mode, press CRTL-C on the command line.

## Development helpers
- Open the repo in github: `make github`.
- See changes: `make diff`
- Clean dependencies and precompiled code: `make clean`
